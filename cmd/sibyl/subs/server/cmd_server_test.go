package server

import (
	"bytes"
	"context"
	"testing"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	"github.com/opensibyl/sibyl2/pkg/core"
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
	uploadCmd.SetArgs([]string{"--src", "../../../..", "--withCtx"})
	uploadCmd.Execute()

	configuration := openapi.NewConfiguration()
	configuration.Scheme = "http"
	configuration.Host = "127.0.0.1:9876"
	apiClient := openapi.NewAPIClient(configuration)
	strings, _, err := apiClient.MAINApi.ApiV1RepoGet(ctx).Execute()
	if err != nil {
		panic(err)
	}
	repo := strings[0]
	revs, _, err := apiClient.MAINApi.ApiV1RevGet(ctx).Repo(repo).Execute()
	if err != nil {
		return
	}
	if len(revs) == 0 {
		panic(nil)
	}
	rev := revs[0]
	files, _, err := apiClient.MAINApi.ApiV1FileGet(ctx).Repo(repo).Rev(rev).Execute()
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
}
