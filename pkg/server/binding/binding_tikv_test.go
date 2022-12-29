package binding

import (
	"context"
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

var tikvDriver Driver

func init() {
	config := object.DefaultExecuteConfig()
	config.TikvAddrs = "127.0.0.1:2379"
	tikvDriver = initTikvDriver(config)
}

func TestTikv(t *testing.T) {
	ctx := context.Background()
	tikvDriver.InitDriver(ctx)
	defer tikvDriver.DeferDriver()

	err := tikvDriver.CreateWorkspace(wc, ctx)
	if err != nil {
		panic(err)
	}

	repos, err := tikvDriver.ReadRepos(ctx)
	if err != nil {
		panic(err)
	}
	core.Log.Debugf("repos: %v", repos)
}
