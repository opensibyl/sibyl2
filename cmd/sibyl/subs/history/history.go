package history

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/ext"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

const D3Template = `
<!DOCTYPE html>
<meta charset="utf-8">
<body>
<script src="https://d3js.org/d3.v5.min.js"></script>
<script src="https://unpkg.com/@hpcc-js/wasm@0.3.11/dist/index.min.js"></script>
<script src="https://unpkg.com/d3-graphviz@3.0.5/build/d3-graphviz.js"></script>
<div id="graph" style="text-align: center; width: 1280px; height: 1280px; border: 1px solid black"></div>
<script>

var dotIndex = 0;
var graphviz = d3.select("#graph").graphviz()
    .transition(function () {
        return d3.transition("main")
            .ease(d3.easeLinear)
            .delay(200)
            .duration(200);
    })
    .logEvents(true)
	.fit(true)
	.width(1280)
    .height(1280)
    .scale(1.0)
    .on("initEnd", render);

function render() {
    var dotLines = dots[dotIndex];
    var dot = dotLines.join('');
    graphviz
		.engine('dot')
        .renderDot(dot)
        .on("end", function () {
            dotIndex = (dotIndex + 1) % dots.length;
            render();
        });
}

var dots = [DOT_REPLACEMENT];
</script>
`

type Storage = map[*object.Commit]map[string][]*extractor.Function

func loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		core.Log.Errorf("load repo from %s failed", gitDir)
		return nil, err
	}
	return repo, nil
}

func reverse(s interface{}) {
	sort.SliceStable(s, func(i, j int) bool {
		return true
	})
}

func handle(gitDir string, output string, full bool) error {
	gitDir, err := filepath.Abs(gitDir)
	if err != nil {
		return err
	}
	lang, err := (&core.Runner{}).GuessLangFromDir(gitDir, nil)
	if err != nil {
		return err
	}
	repo, err := loadRepo(gitDir)
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	cIter, err := repo.Log(&git.LogOptions{From: head.Hash()})

	// what about multi parents?
	var commits []*object.Commit
	cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	reverse(commits)

	// init storage
	storage := make(map[*object.Commit]map[string][]*extractor.Function)
	diffStorage := make(map[*object.Commit]map[string][]*extractor.Function)

	// build the base one
	core.Log.Infof("base commit: %s", commits[0].Hash.String())
	tree, err := commits[0].Tree()
	if err != nil {
		core.Log.Errorf("no hash found: %v", err)
		return err
	}
	commitResult, err := extractFromTree(lang, tree, nil)
	if err != nil {
		core.Log.Errorf("error when extract: %v", err)
		return err
	}
	storage[commits[0]] = commitResult
	prevTree := tree

	// start for each
	for prevIndex, eachCommit := range commits[1:] {
		core.Log.Infof("walk %s", eachCommit.Hash)
		eachTree, err := eachCommit.Tree()

		// about why I use cmd rather than some libs
		// because go-git 's patch has some bugs ...
		gitDiffCmd := exec.Command("git", "diff", commits[prevIndex].Hash.String(), eachCommit.Hash.String())
		gitDiffCmd.Dir = gitDir
		data, err := gitDiffCmd.CombinedOutput()
		if err != nil {
			core.Log.Errorf("git cmd error: %s", data)
			panic(err)
		}
		affected, err := ext.Unified2Affected(data)
		if err != nil {
			return err
		}

		diff, err := eachTree.Diff(prevTree)
		patch, err := diff.Patch()
		if err != nil {
			return err
		}

		curResult := make(map[string][]*extractor.Function)
		for k, v := range storage[commits[prevIndex]] {
			curResult[k] = v
		}

		validFiles := make(map[string][]int, 0)
		for _, each := range patch.FilePatches() {
			from, to := each.Files()
			// file removed
			if from != nil && to == nil {
				delete(curResult, from.Path())
			}
			// file modified
			if from != nil && from != to {
				delete(curResult, from.Path())
			}
			if to != nil {
				validFiles[to.Path()] = affected[to.Path()]
			}
		}
		fromTree, err := extractFromTree(lang, eachTree, func(s string) bool {
			_, ok := validFiles[s]
			return ok
		})

		newFromTree := make(map[string][]*extractor.Function)
		for filename, each := range fromTree {
			affectedFuncs := make([]*extractor.Function, 0)
			lines := validFiles[filename]
			for _, eachFunc := range each {
				if !eachFunc.Span.ContainAnyLine(lines...) {
					affectedFuncs = append(affectedFuncs, eachFunc)
				}
			}
			newFromTree[filename] = affectedFuncs
		}
		diffStorage[eachCommit] = newFromTree
		// update
		for k, v := range fromTree {
			curResult[k] = v
		}
		if err != nil {
			return err
		}
		storage[eachCommit] = curResult

		// update ptr
		prevTree = eachTree
	}
	// done
	core.Log.Infof("commits: %v, start renderering ...", len(storage))

	// render
	var frames []string
	for _, eachCommit := range commits {
		m := storage[eachCommit]

		frame, err := commit2GraphBuf(eachCommit, m, diffStorage[eachCommit], full)
		if err != nil {
			return err
		}
		frames = append(frames, frame)
	}
	core.Log.Infof("images ready, merge with d3.js ...")

	// save them
	var toFill []string
	for _, eachFrame := range frames {
		eachStr := fmt.Sprintf("['%s']\n", strings.ReplaceAll(eachFrame, "\n", ""))
		toFill = append(toFill, eachStr)
	}
	finalContent := strings.Join(toFill, ",")
	finalHtml := strings.Replace(D3Template, "DOT_REPLACEMENT", finalContent, 1)
	err = os.WriteFile(output, []byte(finalHtml), 0644)
	if err != nil {
		return err
	}

	core.Log.Infof("everything done, result saved to: %s", output)

	return nil
}

func commit2GraphBuf(root *object.Commit, m map[string][]*extractor.Function, diffm map[string][]*extractor.Function, full bool) (string, error) {
	graph := gographviz.NewEscape()
	graph.SetStrict(true)
	graph.SetDir(true)
	defaultName := "G"
	graph.SetName(defaultName)
	graph.AddAttr(defaultName, "rankdir", "LR")

	hash := root.Hash.String()
	//committer := strconv.Quote(root.Committer.String())
	//msg := strconv.Quote(root.Message)
	graph.AddNode(defaultName, hash, nil)
	//graph.AddNode(defaultName, msg, nil)
	//graph.AddNode(defaultName, committer, nil)
	//graph.AddEdge(committer, hash, true, nil)
	//graph.AddEdge(msg, hash, true, nil)

	// sort files
	files := make([]string, 0, len(m))
	for f := range m {
		files = append(files, f)
	}
	sort.Strings(files)

	for _, k := range files {
		// build paths
		parts := strings.Split(k, "/")
		curFileName := parts[0]
		graph.AddEdge(hash, curFileName, true, nil)
		for _, eachPart := range parts[1:] {
			newPart := fmt.Sprintf("%s/%s", curFileName, eachPart)
			graph.AddEdge(curFileName, newPart, true, nil)
			graph.AddNode(defaultName, curFileName, nil)
			curFileName = newPart
		}

		funcColorMap := make(map[string]interface{})
		if difflist, ok := diffm[k]; ok {
			graph.AddNode(defaultName, curFileName, map[string]string{
				"style":     "filled",
				"fillcolor": "#9cff3c",
			})

			for _, eachv := range difflist {
				funcColorMap[eachv.GetSignature()] = nil
			}

		} else {
			graph.AddNode(defaultName, curFileName, nil)
		}

		// for perf, no need to add no-edited nodes
		if full || len(funcColorMap) > 0 {
			for _, eachFunc := range m[k] {
				// graphviz parse needed
				funcName := strings.ReplaceAll(eachFunc.GetSignature(), ":", "")
				funcName = strings.ReplaceAll(funcName, "|", "/")
				funcName = strings.ReplaceAll(funcName, "\\", "\\\\")

				if _, ok := funcColorMap[eachFunc.GetSignature()]; ok {
					graph.AddNode(defaultName, funcName, map[string]string{
						"style":     "filled",
						"fillcolor": "#9cff3c",
					})
				} else {
					graph.AddNode(defaultName, funcName, nil)
				}
				graph.AddEdge(curFileName, funcName, true, nil)
			}
		}
	}

	return graph.String(), nil
}

func extractFromTree(lang core.LangType, tree *object.Tree, filter func(string) bool) (map[string][]*extractor.Function, error) {
	ret := make(map[string][]*extractor.Function)
	err := tree.Files().ForEach(func(file *object.File) error {
		fileName := file.Name
		// filter first
		if filter != nil {
			if !filter(fileName) {
				return nil
			}
		}
		// lang filter
		if !lang.MatchName(fileName) {
			return nil
		}

		// ignore binary
		isBin, err := file.IsBinary()
		if err != nil || isBin {
			return nil
		}

		content, err := file.Contents()
		if err != nil {
			return nil
		}

		functions, err := sibyl2.ExtractFromString(content, &sibyl2.ExtractConfig{
			ExtractType: extractor.TypeExtractFunction,
			LangType:    lang,
		})
		if err != nil {
			return err
		}
		core.Log.Debugf("handle file: %s", fileName)
		for _, v := range functions.Units {
			// should not error
			if f, ok := v.(*extractor.Function); ok {
				ret[fileName] = append(ret[fileName], f)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}
