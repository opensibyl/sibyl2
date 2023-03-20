package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/ext"
	"github.com/opensibyl/sibyl2/pkg/server"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestMainScenario(t *testing.T) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		config := object.DefaultExecuteConfig()
		// for performance
		config.BindingConfigPart.DbType = object.DriverTypeMongoDB
		config.EnableLog = true
		_ = server.Execute(config, ctx)
	}()
	defer stop()

	// git prepare
	url := "127.0.0.1:9876"
	src := "../../../.."
	projectName := "sibyl2"
	repo, err := ext.LoadGitRepo(src)
	assert.Nil(t, err)
	head, err := repo.Head()
	assert.Nil(t, err)

	// client
	configuration := openapi.NewConfiguration()
	configuration.Scheme = "http"
	configuration.Host = url
	apiClient := openapi.NewAPIClient(configuration)
	defer apiClient.ScopeApi.ApiV1RepoDelete(ctx).Repo("sibyl2").Execute()

	// do the upload first
	uploadCmd := upload.NewUploadCmd()
	uploadCmd.SetArgs([]string{"--src", src, "--url", fmt.Sprintf("http://%s", url)})
	uploadCmd.Execute()
	time.Sleep(1 * time.Second)

	t.Run("scenario_1_diff_analysis", func(t *testing.T) {
		// scenario 1: diff analysis
		// assume that we have edited these lines
		affectedFileMap := map[string][]int{
			"pkg/core/parser.go": {4, 89, 90, 91, 92, 93, 94, 95, 96},
			"pkg/core/unit.go":   {27, 28, 29},
		}
		// or you can create this map from git diff easily
		_, _ = ext.Unified2Affected([]byte(diffPlain))

		for fileName, lineList := range affectedFileMap {
			strLineList := make([]string, 0, len(lineList))
			for _, each := range lineList {
				strLineList = append(strLineList, strconv.Itoa(each))
			}

			affectedFunctions, _, err := apiClient.BasicQueryApi.
				ApiV1FuncctxGet(ctx).
				Repo(projectName).
				Rev(head.Hash().String()).
				File(fileName).
				Lines(strings.Join(strLineList, ",")).
				Execute()
			assert.Nil(t, err)
			assert.NotEmpty(t, affectedFunctions)

			for _, eachFunc := range affectedFunctions {
				core.Log.Infof("file %s hit func %s, ref: %d, refed: %d",
					fileName, *eachFunc.Name, len(eachFunc.Calls), len(eachFunc.ReverseCalls))
			}

			// query their reverse call chains?
			for _, eachFunc := range affectedFunctions {
				chain, _, err := apiClient.SignatureQueryApi.
					ApiV1SignatureFuncctxRchainGet(ctx).
					Repo(projectName).
					Rev(head.Hash().String()).
					Signature(eachFunc.GetSignature()).
					Depth(5).
					Execute()
				assert.Nil(t, err)
				// chain is a tree-like object
				// access it with dfs/bfs
				if chain.ReverseCallChains != nil {
					for _, each := range chain.ReverseCallChains.GetChildren() {
						core.Log.Infof("cur: %v", *each.Content)
						for _, eachSub := range each.GetChildren() {
							core.Log.Infof("cur: %v", *eachSub.Content)
							// continue ...
							// eachSub.GetChildren()
						}
					}
				}

				// also a normal call chain
				chain, _, err = apiClient.SignatureQueryApi.
					ApiV1SignatureFuncctxChainGet(ctx).
					Repo(projectName).
					Rev(head.Hash().String()).
					Signature(eachFunc.GetSignature()).
					Depth(5).
					Execute()
				assert.Nil(t, err)
				assert.NotNil(t, chain)
			}

			for _, eachFunc := range affectedFunctions {
				// get all the calls details?
				for _, eachCall := range eachFunc.Calls {
					detail, _, err := apiClient.SignatureQueryApi.
						ApiV1SignatureFuncGet(ctx).
						Repo(projectName).
						Rev(head.Hash().String()).
						Signature(eachCall).
						Execute()
					assert.Nil(t, err)
					core.Log.Infof("call: %v", detail)
				}
			}
		}
	})

	// scenario 2: specific global search
	t.Run("scenario_2_global_search", func(t *testing.T) {
		functionWithPaths, _, err := apiClient.RegexQueryApi.
			ApiV1RegexFuncGet(ctx).
			Repo(projectName).
			Rev(head.Hash().String()).
			Field("name").
			Regex(".*Handle.*").
			Execute()
		assert.Nil(t, err)
		assert.NotEmpty(t, functionWithPaths)
		for _, each := range functionWithPaths {
			assert.True(t, strings.Contains(*each.Name, "Handle"))
			// and see where it is
			assert.NotEmpty(t, *each.Path)
		}
	})

	// scenario 3: hot functions
	t.Run("scenario_3_hot_functions", func(t *testing.T) {
		fc, _, err := apiClient.ReferenceQueryApi.
			ApiV1ReferenceCountFuncctxReverseGet(ctx).
			Repo(projectName).
			Rev(head.Hash().String()).
			MoreThan(10).
			LessThan(100).
			Execute()
		assert.Nil(t, err)
		assert.NotEmpty(t, fc)
	})

	// scenario 4: file level relationship graph
	t.Run("scenario_4_relationship_graph", func(t *testing.T) {
		ecAll := &EcAll{}

		files, _, err := apiClient.ScopeApi.ApiV1FileGet(ctx).Repo(projectName).Rev(head.Hash().String()).Execute()
		assert.Nil(t, err)

		// create nodes
		category := make([]string, 0)
		for _, eachFile := range files {
			dirName := filepath.ToSlash(filepath.Dir(eachFile))

			if !slices.Contains(category, dirName) {
				category = append(category, dirName)
			}
			loc := slices.Index(category, dirName)

			ecAll.Nodes = append(ecAll.Nodes, &ECNode{
				Id:       eachFile,
				Category: loc,
				Name:     eachFile,
			})
		}
		for _, eachCategory := range category {
			ecAll.Categories = append(ecAll.Categories, &ECCategory{Name: eachCategory})
		}

		// create edges
		for _, eachFile := range files {
			ctxs, _, err := apiClient.BasicQueryApi.ApiV1FuncctxGet(ctx).Repo(projectName).Rev(head.Hash().String()).File(eachFile).Execute()
			assert.Nil(t, err)
			for _, each := range ctxs {
				for _, eachCall := range each.Calls {
					f, _, err := apiClient.SignatureQueryApi.ApiV1SignatureFuncGet(ctx).Repo(projectName).Rev(head.Hash().String()).Signature(eachCall).Execute()
					assert.Nil(t, err)
					// create edge between eachFile and f.Path
					ecAll.Links = append(ecAll.Links, &ECLink{
						Source: eachFile,
						Target: *f.Path,
					})
				}
			}
		}

		// export
		_, err = json.Marshal(ecAll)
		assert.Nil(t, err)
	})

	// scenario 5: tag test functions, and query by tag
	t.Run("scenario_5_tag_test_function", func(t *testing.T) {
		rev := head.Hash().String()
		functions, _, err := apiClient.RegexQueryApi.
			ApiV1RegexFuncGet(ctx).
			Repo(projectName).
			Rev(rev).
			Field("name").
			Regex("^Test.*").
			Execute()
		assert.Nil(t, err)
		assert.NotEmpty(t, functions)
		// tag these functions with `TEST_METHOD`
		TagForCases := "TEST_METHOD"
		for _, eachFunc := range functions {
			eachFuncSign := eachFunc.GetSignature()
			_, err := apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
				RepoId:    &projectName,
				RevHash:   &rev,
				Signature: &eachFuncSign,
				Tag:       &TagForCases,
			}).Execute()
			assert.Nil(t, err)
		}

		// query
		functionsFromQuery, _, err := apiClient.TagApi.ApiV1TagFuncGet(ctx).Repo(projectName).Rev(rev).Tag(TagForCases).Execute()
		for _, eachFunc := range functionsFromQuery {
			core.Log.Infof("tag from query: %v", eachFunc)
		}
		assert.Nil(t, err)
		assert.NotEmpty(t, functionsFromQuery)
		assert.Equal(t, len(functions), len(functionsFromQuery))
	})

	t.Cleanup(func() {
		stop()
		time.Sleep(2000 * time.Millisecond)
	})
}

const diffPlain = `
diff --git a/pkg/core/parser.go b/pkg/core/parser.go
--- a/pkg/core/parser.go
+++ b/pkg/core/parser.go
@@ -2,6 +2,7 @@ package core
 
 import (
 	"context"
+
 	sitter "github.com/smacker/go-tree-sitter"
 )
 
@@ -84,8 +85,14 @@ func (p *Parser) node2Unit(data []byte, node *sitter.Node, fieldName string, par
 
 	// range
 	ret.Span = Span{
-		Start: Point{Row: node.StartPoint().Row, Column: node.StartPoint().Column},
-		End:   Point{Row: node.EndPoint().Row, Column: node.EndPoint().Column},
+		Start: Point{
+			Row:    node.StartPoint().Row,
+			Column: node.StartPoint().Column,
+		},
+		End: Point{
+			Row:    node.EndPoint().Row,
+			Column: node.EndPoint().Column,
+		},
 	}
 	ret.ParentUnit = parentUnit
 	return ret, nil
diff --git a/pkg/core/unit.go b/pkg/core/unit.go
--- a/pkg/core/unit.go
+++ b/pkg/core/unit.go
@@ -23,8 +23,9 @@ func (s *Span) Lines() []int {
 
 func (s *Span) ContainLine(lineNum int) bool {
 	// real line number
-	uint32Line := uint32(lineNum) + 1
-	return s.Start.Row <= uint32Line && uint32Line <= s.End.Row
+	realLineNum := lineNum + 1
+	// int can be 32 or 64 bits
+	return int(s.Start.Row) <= realLineNum && realLineNum <= int(s.End.Row)
 }
 
 func (s *Span) ContainAnyLine(lineNums ...int) bool {
`

// created from https://echarts.apache.org/examples/zh/editor.html?c=graph-circular-layout
type ECNode struct {
	Id       string `json:"id"`
	Category int    `json:"category"`
	Name     string `json:"name"`
}

type ECLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type ECCategory struct {
	Name string `json:"name"`
}

type EcAll struct {
	Nodes      []*ECNode     `json:"nodes"`
	Links      []*ECLink     `json:"links"`
	Categories []*ECCategory `json:"categories"`
}
