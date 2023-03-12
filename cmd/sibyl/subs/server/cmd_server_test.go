package server

import (
	"context"
	"os/signal"
	"syscall"
	"testing"
	"time"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	ctx := context.Background()
	sibylContext, cancel := context.WithCancel(ctx)
	defer cancel()

	sibylContext, stop := signal.NotifyContext(sibylContext, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		config := object.DefaultExecuteConfig()
		// for performance
		config.BindingConfigPart.DbType = object.DriverTypeInMemory
		config.EnableLog = true
		_ = server.Execute(config, sibylContext)
	}()
	defer stop()

	// do the upload first
	uploadCmd := upload.NewUploadCmd()
	uploadCmd.SetArgs([]string{"--src", "../../../..", "--withCtx", "--withClass", "--lang", "GOLANG"})
	uploadCmd.Execute()
	time.Sleep(1200 * time.Millisecond)

	configuration := openapi.NewConfiguration()
	configuration.Scheme = "http"
	configuration.Host = "127.0.0.1:9876"
	apiClient := openapi.NewAPIClient(configuration)
	defer apiClient.ScopeApi.ApiV1RepoDelete(ctx).Repo("sibyl2").Execute()

	repos, _, err := apiClient.ScopeApi.ApiV1RepoGet(ctx).Execute()
	if err != nil {
		panic(err)
	}
	repo := repos[0]
	revs, _, err := apiClient.ScopeApi.ApiV1RevGet(ctx).Repo(repo).Execute()
	if err != nil {
		return
	}
	if len(revs) == 0 {
		panic(nil)
	}
	rev := revs[0]
	files, _, err := apiClient.ScopeApi.ApiV1FileGet(ctx).Repo(repo).Rev(rev).Execute()
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		panic(nil)
	}
	core.Log.Debugf("file count: %d", len(files))

	functions, _, err := apiClient.BasicQueryApi.ApiV1FuncGet(ctx).Repo(repo).Rev(rev).File("extract.go").Execute()
	if err != nil {
		panic(err)
	}
	if len(functions) == 0 {
		panic(err)
	}

	fnCtxs, _, err := apiClient.BasicQueryApi.ApiV1FuncctxGet(ctx).Repo(repo).Rev(rev).File("extract.go").Execute()
	if err != nil {
		panic(err)
	}
	if len(fnCtxs) == 0 {
		panic(err)
	}

	classes, _, err := apiClient.BasicQueryApi.ApiV1ClazzGet(ctx).Repo(repo).Rev(rev).File("extract.go").Execute()
	if err != nil {
		panic(err)
	}
	if len(classes) == 0 {
		panic(err)
	}

	signatures, _, err := apiClient.SignatureQueryApi.ApiV1SignatureRegexFuncGet(ctx).Repo(repo).Rev(rev).Regex(".*").Execute()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(signatures))
	f, _, err := apiClient.SignatureQueryApi.ApiV1SignatureFuncGet(ctx).Repo(repo).Rev(rev).Signature(signatures[0]).Execute()
	assert.Nil(t, err)
	assert.NotNil(t, f)

	assert.Nil(t, err)
	fwr, _, err := apiClient.RegexQueryApi.
		ApiV1RegexFuncGet(ctx).
		Repo(repo).
		Rev(rev).
		Field("name").
		Regex(".*Handle.*").
		Execute()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(fwr))
	for _, each := range fwr {
		assert.Contains(t, *each.Name, "Handle")
	}

	cwr, _, err := apiClient.RegexQueryApi.
		ApiV1RegexClazzGet(ctx).
		Repo(repo).
		Rev(rev).
		Field("name").
		Regex(".*").
		Execute()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(cwr))

	fcwr, _, err := apiClient.RegexQueryApi.
		ApiV1RegexFuncctxGet(ctx).
		Repo(repo).
		Rev(rev).
		Field("name").
		Regex(".*Handle.*").
		Execute()
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(fcwr))
	for _, each := range fcwr {
		assert.Contains(t, *each.Name, "Handle")
	}

	// tag
	newTag := "NEW_TAG"
	_, err = apiClient.TagApi.ApiV1TagFuncPost(ctx).Payload(openapi.ServiceTagUpload{
		RepoId:    &repo,
		RevHash:   &rev,
		Signature: f.Signature,
		Tag:       &newTag,
	}).Execute()
	assert.Nil(t, err)
	strings, _, err := apiClient.TagApi.ApiV1TagFuncGet(ctx).Repo(repo).Rev(rev).Tag(newTag).Execute()
	assert.Nil(t, err)
	assert.Len(t, strings, 1)

	t.Cleanup(func() {
		stop()
		time.Sleep(1000 * time.Millisecond)
	})
}
