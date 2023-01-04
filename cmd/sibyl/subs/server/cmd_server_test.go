package server

import (
	"bytes"
	"context"
	"testing"
	"time"

	openapi "github.com/opensibyl/sibyl-go-client"
	"github.com/opensibyl/sibyl2/cmd/sibyl/subs/upload"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

// todo: create this diff with program
const unifiedDiff = `
diff --git a/pkg/server/service/ops_s.go b/pkg/server/service/admin_s.go
similarity index 100%
rename from pkg/server/service/ops_s.go
rename to pkg/server/service/admin_s.go
diff --git a/pkg/server/service/main_query_s.go b/pkg/server/service/query_s.go
similarity index 97%
rename from pkg/server/service/main_query_s.go
rename to pkg/server/service/query_s.go
index 06cb032..1ec10da 100644
--- a/pkg/server/service/main_query_s.go
+++ b/pkg/server/service/query_s.go
@@ -215,3 +215,9 @@ func handleClazzQuery(repo string, rev string, file string) ([]*sibyl2.ClazzWith
 	}
 	return functions, nil
 }
+
+func InitService(_ object.ExecuteConfig, ctx context.Context, driver binding.Driver, q queue.Queue) {
+	sharedContext = ctx
+	sharedDriver = driver
+	sharedQueue = q
+}
diff --git a/pkg/server/service/shared.go b/pkg/server/service/shared.go
deleted file mode 100644
index b7c9262..0000000
--- a/pkg/server/service/shared.go
+++ /dev/null
@@ -1,15 +0,0 @@
-package service
-
-import (
-	"context"
-
-	"github.com/opensibyl/sibyl2/pkg/server/binding"
-	"github.com/opensibyl/sibyl2/pkg/server/object"
-	"github.com/opensibyl/sibyl2/pkg/server/queue"
-)
-
-func InitService(_ object.ExecuteConfig, ctx context.Context, driver binding.Driver, q queue.Queue) {
-	sharedContext = ctx
-	sharedDriver = driver
-	sharedQueue = q
-}
diff --git a/pkg/server/service/main_upload_s.go b/pkg/server/service/upload_s.go
similarity index 100%
rename from pkg/server/service/main_upload_s.go
rename to pkg/server/service/upload_s.go
`

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

	diff, _, err := apiClient.EXTRASApi.ApiV1FuncctxDiffGet(ctx).Repo(repo).Rev(rev).Diff(unifiedDiff).Execute()
	assert.Nil(t, err)
	core.Log.Infof("diff ctxs: %v", diff)

	funcs, _, err := apiClient.EXTRASApi.ApiV1FuncDiffGet(ctx).Repo(repo).Rev(rev).Diff(unifiedDiff).Execute()
	assert.Nil(t, err)
	core.Log.Infof("diff funcs: %v", funcs)
}
