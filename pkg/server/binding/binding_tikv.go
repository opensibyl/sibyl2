package binding

import (
	"context"
	"strings"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tikv/client-go/v2/txnkv"
)

type TiKVDriver struct {
	client    *txnkv.Client
	addresses []string
}

func initTikvDriver(config object.ExecuteConfig) Driver {
	addresses := strings.Split(config.TikvAddrs, ",")
	return &TiKVDriver{
		addresses: addresses,
	}
}

func (t *TiKVDriver) GetType() object.DriverType {
	return object.DtTikv
}

func (t *TiKVDriver) InitDriver(_ context.Context) error {
	client, err := txnkv.NewClient(t.addresses)
	if err != nil {
		return err
	}
	t.client = client
	return nil
}

func (t *TiKVDriver) DeferDriver() error {
	if err := t.client.Close(); err != nil {
		return err
	}
	t.client = nil
	return nil
}

func (t *TiKVDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContext, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	byteKey := []byte(toRevKey(key).String())
	txn, err := t.client.Begin()

	// tikv does not allow set nil value
	err = txn.Set(byteKey, byteKey)
	if err != nil {
		return err
	}
	err = txn.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *TiKVDriver) readRawRevs() ([]*revKey, error) {
	txn, err := t.client.Begin()
	if err != nil {
		return nil, err
	}
	iter, err := txn.Iter([]byte(revPrefix), nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	ret := make([]*revKey, 0)
	for iter.Valid() {
		ret = append(ret, parseRevKey(string(iter.Key())))
		err := iter.Next()
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (t *TiKVDriver) ReadRepos(_ context.Context) ([]string, error) {
	revs, err := t.readRawRevs()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, eachRev := range revs {
		wc, err := WorkspaceConfigFromKey(eachRev.hash)
		if err != nil {
			return nil, err
		}
		ret = append(ret, wc.RepoId)
	}
	return ret, nil
}

func (t *TiKVDriver) ReadRevs(repoId string, ctx context.Context) ([]string, error) {
	revs, err := t.readRawRevs()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, eachRev := range revs {
		wc, err := WorkspaceConfigFromKey(eachRev.hash)
		if err != nil {
			return nil, err
		}
		if wc.RepoId == repoId {
			ret = append(ret, wc.RevHash)
		}
	}
	return ret, nil
}

func (t *TiKVDriver) ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithPath, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContext, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) UpdateRevProperties(wc *object.WorkspaceConfig, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) UpdateFileProperties(wc *object.WorkspaceConfig, path string, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) UpdateFuncProperties(wc *object.WorkspaceConfig, signature string, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
