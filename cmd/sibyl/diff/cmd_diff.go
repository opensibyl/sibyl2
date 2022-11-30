package diff

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

type FunctionWithWeight struct {
	*extractor.Function
	ReferenceCount  int `json:"referenceCount"`
	ReferencedCount int `json:"referencedCount"`
}

type FileFragment struct {
	Path      string                `json:"path"`
	Lines     []int                 `json:"lines"`
	Functions []*FunctionWithWeight `json:"functions"`
}

type ParseResult struct {
	Fragments []*FileFragment `json:"fragments"`
}

type AffectedLineMap = map[string][]int

func Unified2Affected(patch []byte) (AffectedLineMap, error) {
	parsed, _, err := gitdiff.Parse(bytes.NewReader(patch))
	if err != nil {
		return nil, err
	}

	affectedMap := make(map[string][]int)
	for _, each := range parsed {
		if each.IsBinary || each.IsDelete {
			continue
		}
		affectedMap[each.NewName] = make([]int, 0)
		fragments := each.TextFragments
		for _, eachF := range fragments {
			left := int(eachF.NewPosition)

			for i, eachLine := range eachF.Lines {
				if eachLine.New() && eachLine.Op == gitdiff.OpAdd {
					affectedMap[each.NewName] = append(affectedMap[each.NewName], left+i-1)
				}
			}
		}
	}
	return affectedMap, nil
}

func parsePatch(srcDir string, patchFile string) (*ParseResult, error) {
	patch, err := os.ReadFile(patchFile)
	if err != nil {
		return nil, err
	}
	srcDir, err = filepath.Abs(srcDir)
	if err != nil {
		return nil, err
	}

	affectedMap, err := Unified2Affected(patch)
	if err != nil {
		return nil, err
	}

	// data ready
	result, err := affectedLines2Functions(srcDir, &affectedMap)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func affectedLines2Functions(srcDir string, m *AffectedLineMap) (*ParseResult, error) {
	f, err := sibyl2.ExtractFunction(srcDir, sibyl2.DefaultConfig())
	if err != nil {
		return nil, err
	}

	s, err := sibyl2.ExtractSymbol(srcDir, sibyl2.DefaultConfig())
	if err != nil {
		return nil, err
	}
	g, err := sibyl2.AnalyzeFuncGraph(f, s)
	if err != nil {
		return nil, err
	}

	result := &ParseResult{}
	for fileName, lines := range *m {
		fragment := &FileFragment{}
		fragment.Path = fileName
		fragment.Lines = lines
		for _, eachFuncFile := range f {
			// make it map better
			if eachFuncFile.Path == fileName {
				for _, eachFuncUnit := range eachFuncFile.Units {
					if eachFuncUnit.GetSpan().ContainAnyLine(lines...) {
						related := g.FindRelated(eachFuncUnit)
						fww := &FunctionWithWeight{
							Function:        eachFuncUnit,
							ReferenceCount:  len(related.Calls),
							ReferencedCount: len(related.ReverseCalls),
						}
						fragment.Functions = append(fragment.Functions, fww)
					}
				}
				break
			}
		}
		result.Fragments = append(result.Fragments, fragment)
	}

	return result, nil
}

var diffSrc string
var diffPatch string

func NewDiffCommand() *cobra.Command {
	diffCmd := &cobra.Command{
		Use:    "diff",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			results, err := parsePatch(diffSrc, diffPatch)
			if err != nil {
				panic(err)
			}
			output, err := json.MarshalIndent(&results, "", "  ")
			if err != nil {
				panic(err)
			}
			err = os.WriteFile("b.json", output, 0644)
			if err != nil {
				panic(err)
			}
		},
	}
	diffCmd.PersistentFlags().StringVar(&diffSrc, "src", ".", "src dir path")
	diffCmd.PersistentFlags().StringVar(&diffPatch, "patch", "", "patch")
	return diffCmd
}
