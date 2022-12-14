package binding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
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
		eachFuncV, err := json.Marshal(eachFunc)
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
