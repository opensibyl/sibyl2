package server

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	cmd := NewServerCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	// run server
	ctx := context.Background()
	go cmd.ExecuteContext(ctx)

	// do the upload first
	uploadCmd := upload.NewUploadCmd()
	uploadCmd.SetArgs([]string{"--src", "../../../..", "--withCtx", "--withClass"})
	uploadCmd.Execute()
	time.Sleep(1 * time.Second)

	configuration := openapi.NewConfiguration()
	configuration.Scheme = "http"
	configuration.Host = "127.0.0.1:9876"
	apiClient := openapi.NewAPIClient(configuration)
	repos, _, err := apiClient.SCOPEApi.ApiV1RepoGet(ctx).Execute()
	if err != nil {
		panic(err)
	}
	repo := repos[0]
	revs, _, err := apiClient.SCOPEApi.ApiV1RevGet(ctx).Repo(repo).Execute()
	if err != nil {
		return
	}
	if len(revs) == 0 {
		panic(nil)
	}
	rev := revs[0]
	files, _, err := apiClient.SCOPEApi.ApiV1FileGet(ctx).Repo(repo).Rev(rev).Execute()
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		panic(nil)
	}
	core.Log.Debugf("file count: %d", len(files))

	functions, _, err := apiClient.MAINApi.ApiV1FuncGet(ctx).Repo(repo).Rev(rev).File("extract.go").Execute()
	if err != nil {
		panic(err)
	}
	if len(functions) == 0 {
		panic(err)
	}

	fnCtxs, _, err := apiClient.MAINApi.ApiV1FuncctxGet(ctx).Repo(repo).Rev(rev).File("extract.go").Execute()
	if err != nil {
		panic(err)
	}
	if len(fnCtxs) == 0 {
		panic(err)
	}

	classes, _, err := apiClient.MAINApi.ApiV1ClazzGet(ctx).Repo(repo).Rev(rev).File("extract.go").Execute()
	if err != nil {
		panic(err)
	}
	if len(classes) == 0 {
		panic(err)
	}

	signatures, _, err := apiClient.EXPERIMENTALApi.ApiV1FuncSignatureGet(ctx).Repo(repo).Rev(rev).Regex(".*").Execute()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(signatures))
	f, _, err := apiClient.EXPERIMENTALApi.ApiV1FuncWithSignatureGet(ctx).Repo(repo).Rev(rev).Signature(signatures[0]).Execute()
	assert.Nil(t, err)
	assert.NotNil(t, f)

	funcRule := make(map[string]string)
	funcRule["name"] = ".*Handle.*"
	ruleStr, err := json.Marshal(funcRule)
	assert.Nil(t, err)
	fwr, _, err := apiClient.EXPERIMENTALApi.ApiV1FuncWithRuleGet(ctx).Repo(repo).Rev(rev).Rule(string(ruleStr)).Execute()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(fwr))
	for _, each := range fwr {
		assert.Contains(t, *each.Name, "Handle")
	}

	clazzRule := make(map[string]string)
	clazzRule["name"] = ".*"
	clazzRuleStr, err := json.Marshal(clazzRule)
	assert.Nil(t, err)
	cwr, _, err := apiClient.EXPERIMENTALApi.ApiV1ClazzWithRuleGet(ctx).Repo(repo).Rev(rev).Rule(string(clazzRuleStr)).Execute()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(cwr))

	fcwr, _, err := apiClient.EXPERIMENTALApi.ApiV1FuncWithRuleGet(ctx).Repo(repo).Rev(rev).Rule(string(ruleStr)).Execute()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(fcwr))
	for _, each := range fcwr {
		assert.Contains(t, *each.Name, "Handle")
	}
}
