package binding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"golang.org/x/exp/slices"
)

func (t *tikvDriver) CreateClazzFile(wc *object.WorkspaceConfig, c *extractor.ClazzFileResult, ctx context.Context) error {
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
		eachClazzWithPath := &sibyl2.ClazzWithPath{
			Clazz: eachClazz,
			Path:  c.Path,
		}
		eachClazzValue, err := json.Marshal(eachClazzWithPath)
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

func (t *tikvDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
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
		eachFuncWithPath := &sibyl2.FunctionWithPath{
			Function: eachFunc,
			Path:     f.Path,
		}
		eachFuncV, err := json.Marshal(eachFuncWithPath)
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

func (t *tikvDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContextSlim, ctx context.Context) error {
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
	eachFuncV, err := json.Marshal(f)
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

func (t *tikvDriver) CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}

	revInfo := object.NewRevInfo(wc.RevHash)
	revInfoStr, err := json.Marshal(revInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal rev info: %w", err)
	}

	byteKey := []byte(ToRevKey(key).String())
	txn, err := t.client.Begin()
	if err != nil {
		return err
	}

	err = txn.Set(byteKey, revInfoStr)
	if err != nil {
		return err
	}
	err = txn.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *tikvDriver) CreateFuncTag(wc *object.WorkspaceConfig, signature string, tag string, ctx context.Context) error {
	txn, err := t.client.Begin()
	if err != nil {
		return err
	}

	f, err := t.ReadFunctionWithSignature(wc, signature, ctx)
	if err != nil {
		return err
	}

	// duplicated
	if slices.Contains(f.Tags, tag) {
		return nil
	}
	f.AddTag(tag)

	// key
	key, err := wc.Key()
	if err != nil {
		return err
	}
	fk := toFileKey(key, f.Path)
	curFuncKey := toFuncKey(fk.RevHash, fk.FileHash, f.GetSignature())

	// write
	newFuncBytes, err := json.Marshal(f)
	if err != nil {
		return err
	}

	err = txn.Set([]byte(curFuncKey.String()), newFuncBytes)
	if err != nil {
		return err
	}
	err = txn.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
