package binding

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tikv/client-go/v2/kv"
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
	return object.DriverTypeTikv
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

func (t *TiKVDriver) CreateClazzFile(wc *object.WorkspaceConfig, c *extractor.ClazzFileResult, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}

	fk := toFileKey(key, c.Path)
	byteKey := []byte(fk.String())

	txn, err := t.client.Begin()
	if err != nil {
		return err
	}

	// tikv does not allow set nil value
	// todo: in the future, value will be replaced with file desc info (something like author/size
	err = txn.Set(byteKey, byteKey)
	if err != nil {
		return err
	}

	for _, eachClazz := range c.Units {
		eachClazzKey := toClazzKey(fk.RevHash, fk.FileHash, eachClazz.GetSignature())
		eachClazzValue, err := eachClazz.ToJson()
		if err != nil {
			continue
		}
		err = txn.Set([]byte(eachClazzKey.String()), eachClazzValue)
		if err != nil {
			return err
		}
	}

	// TiKV uses the optimistic transaction model by default
	err = txn.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *TiKVDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}

	fk := toFileKey(key, f.Path)
	byteKey := []byte(fk.String())

	txn, err := t.client.Begin()
	if err != nil {
		return err
	}

	// tikv does not allow set nil value
	// todo: in the future, value will be replaced with file desc info (something like author/size
	err = txn.Set(byteKey, byteKey)
	if err != nil {
		return err
	}

	for _, eachFunc := range f.Units {
		eachFuncKey := toFuncKey(fk.RevHash, fk.FileHash, eachFunc.GetSignature())
		eachFuncV, err := eachFunc.ToJson()
		if err != nil {
			continue
		}
		err = txn.Set([]byte(eachFuncKey.String()), eachFuncV)
		if err != nil {
			return err
		}
	}

	// TiKV uses the optimistic transaction model by default
	err = txn.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *TiKVDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContext, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}

	fk := toFileKey(key, f.Path)
	byteKey := []byte(fk.String())

	txn, err := t.client.Begin()
	if err != nil {
		return err
	}

	// tikv does not allow set nil value
	err = txn.Set(byteKey, byteKey)
	if err != nil {
		return err
	}

	eachFuncKey := toFuncCtxKey(fk.RevHash, fk.FileHash, f.GetSignature())
	eachFuncV, err := f.ToJson()
	if err != nil {
		return err
	}
	err = txn.Set([]byte(eachFuncKey.String()), eachFuncV)
	if err != nil {
		return err
	}
	// TiKV uses the optimistic transaction model by default
	err = txn.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *TiKVDriver) CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	byteKey := []byte(ToRevKey(key).String())
	txn, err := t.client.Begin()
	if err != nil {
		return err
	}

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
	snapshot := t.client.GetSnapshot(math.MaxUint64)
	keyByte := []byte(revPrefix)
	iter, err := snapshot.Iter(keyByte, nil)
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
	m := make(map[string]struct{}, 0)
	for _, eachRev := range revs {
		wc, err := WorkspaceConfigFromKey(eachRev.Hash)
		if err != nil {
			return nil, err
		}
		m[wc.RepoId] = struct{}{}
	}

	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret, nil
}

func (t *TiKVDriver) ReadRevs(repoId string, _ context.Context) ([]string, error) {
	revs, err := t.readRawRevs()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, eachRev := range revs {
		wc, err := WorkspaceConfigFromKey(eachRev.Hash)
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
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	searchResult := make([]string, 0)

	txn := t.client.GetSnapshot(math.MaxUint64)

	prefixStr := rk.ToScanPrefix() + "file|"
	prefix := []byte(prefixStr)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()
	for iter.Valid() {
		k := iter.Key()
		searchResult = append(searchResult, strings.TrimPrefix(string(k), prefixStr))
		err := iter.Next()
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (t *TiKVDriver) ReadFunctionSignaturesWithRegex(wc *object.WorkspaceConfig, regex string, ctx context.Context) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (t *TiKVDriver) ReadClasses(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	fk := toFileKey(key, path)

	searchResult := make([]*sibyl2.ClazzWithPath, 0)

	prefixStr := fk.ToScanPrefix() + "clazz|"
	prefix := []byte(prefixStr)

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()

	for iter.Valid() {
		c := &sibyl2.ClazzWithPath{}
		err := json.Unmarshal(iter.Value(), c)
		if err != nil {
			return nil, err
		}
		c.Path = path
		searchResult = append(searchResult, c)
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}

	return searchResult, nil
}

func (t *TiKVDriver) ReadClassesWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	classes, err := t.ReadClasses(wc, path, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*sibyl2.ClazzWithPath, 0)
	for _, each := range classes {
		if each.GetSpan().ContainAnyLine(lines...) {
			ret = append(ret, each)
		}
	}
	return ret, nil
}

func (t *TiKVDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	fk := toFileKey(key, path)

	searchResult := make([]*sibyl2.FunctionWithPath, 0)

	prefixStr := fk.ToScanPrefix() + "func|"
	prefix := []byte(prefixStr)

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()

	for iter.Valid() {
		f := &sibyl2.FunctionWithPath{}
		err := json.Unmarshal(iter.Value(), f)
		if err != nil {
			return nil, err
		}
		f.Path = path
		searchResult = append(searchResult, f)
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}

	return searchResult, nil
}

func (t *TiKVDriver) ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionContext, error) {
	functions, err := t.ReadFunctionsWithLines(wc, path, lines, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*sibyl2.FunctionContext, 0)
	for _, eachFunc := range functions {
		functionContext, err := t.ReadFunctionContextWithSignature(wc, eachFunc.GetSignature(), ctx)
		if err != nil {
			return nil, err
		}
		ret = append(ret, functionContext)
	}
	return ret, nil
}

func (t *TiKVDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)

	prefixStr := rk.ToScanPrefix() + "file_"
	prefix := []byte(prefixStr)
	shouldContain := "func|" + signature

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()

	var ret *sibyl2.FunctionWithPath
	for iter.Valid() {
		k := string(iter.Key())
		if strings.Contains(k, shouldContain) {
			err := json.Unmarshal(iter.Value(), ret)
			if err != nil {
				return nil, err
			}
			return ret, nil
		}
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	// did not find anything
	return nil, nil
}

func (t *TiKVDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	functions, err := t.ReadFunctions(wc, path, ctx)
	if err != nil {
		return nil, err
	}

	searchResult := make([]*sibyl2.FunctionWithPath, 0)
	for _, each := range functions {
		if each.Span.ContainAnyLine(lines...) {
			searchResult = append(searchResult, each)
		}
	}
	return searchResult, nil
}

func (t *TiKVDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContext, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)

	prefixStr := rk.ToScanPrefix() + "file_"
	prefix := []byte(prefixStr)
	shouldContain := "funcctx|" + signature

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()

	for iter.Valid() {
		k := string(iter.Key())
		if strings.Contains(k, shouldContain) {
			funcCtx, err := sibyl2.Json2FuncCtx(iter.Value())
			if err != nil {
				return nil, err
			}
			// break scan
			return funcCtx, nil
		}
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	// did not find anything
	return nil, nil
}

func (t *TiKVDriver) UpdateRevProperties(wc *object.WorkspaceConfig, k string, v any, ctx context.Context) error {
	// TODO implement me
	return errors.New("implement me")
}

func (t *TiKVDriver) UpdateFileProperties(wc *object.WorkspaceConfig, path string, k string, v any, ctx context.Context) error {
	// TODO implement me
	return errors.New("implement me")
}

func (t *TiKVDriver) UpdateFuncProperties(wc *object.WorkspaceConfig, signature string, k string, v any, ctx context.Context) error {
	// TODO implement me
	return errors.New("implement me")
}

func (t *TiKVDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	rk := ToRevKey(key)
	itself := []byte(rk.String())
	sons := []byte(rk.ToScanPrefix())

	_, err = t.client.DeleteRange(ctx, itself, kv.PrefixNextKey(itself), 1)
	if err != nil {
		return err
	}
	_, err = t.client.DeleteRange(ctx, sons, kv.PrefixNextKey(sons), 1)
	if err != nil {
		return err
	}

	return nil
}
