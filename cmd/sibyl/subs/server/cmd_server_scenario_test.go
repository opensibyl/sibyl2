package server

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/ext"
	"github.com/stretchr/testify/assert"
)

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

func TestMainScenario(t *testing.T) {
	ctx := context.Background()

	cmd := NewServerCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	// run server
	go cmd.ExecuteContext(ctx)

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

	// do the upload first
	uploadCmd := upload.NewUploadCmd()
	uploadCmd.SetArgs([]string{"--src", src, "--url", fmt.Sprintf("http://%s", url)})
	uploadCmd.Execute()
	time.Sleep(1 * time.Second)

	// scenario 1: diff analysis
	affectedFileMap, err := ext.Unified2Affected([]byte(diffPlain))
	if err != nil {
		panic(err)
	}

	for fileName, lineList := range affectedFileMap {
		strLineList := make([]string, 0, len(lineList))
		for _, each := range lineList {
			strLineList = append(strLineList, strconv.Itoa(each))
		}

		functionWithPaths, _, err := apiClient.MAINApi.
			ApiV1FuncctxGet(ctx).
			Repo(projectName).
			Rev(head.Hash().String()).
			File(fileName).
			Lines(strings.Join(strLineList, ",")).
			Execute()
		assert.Nil(t, err)
		assert.NotEmpty(t, functionWithPaths)

		for _, each := range functionWithPaths {
			core.Log.Infof("file %s hit func %s, ref: %d, refed: %d",
				fileName, *each.Name, len(each.Calls), len(each.ReverseCalls))
		}
	}

	// scenario 2: specific global search
	assert.Nil(t, err)
	functionWithPaths, _, err := apiClient.EXPERIMENTALApi.
		ApiV1FuncWithRegexGet(ctx).
		Repo(projectName).
		Rev(head.Hash().String()).
		Field("name").
		Regex(".*Handle.*").
		Execute()
	assert.Nil(t, err)
	assert.NotEmpty(t, functionWithPaths)
	for _, each := range functionWithPaths {
		assert.True(t, strings.Contains(*each.Name, "Handle"))
	}
}
