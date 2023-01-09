package binding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

func (d *badgerDriver) CreateClazzFile(wc *object.WorkspaceConfig, c *extractor.ClazzFileResult, _ context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}

	err = d.db.Update(func(txn *badger.Txn) error {
		fk := toFileKey(key, c.Path)
		byteKey := []byte(fk.String())

		// todo: keep origin value
		err = txn.Set(byteKey, nil)
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

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *badgerDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, _ context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	err = d.db.Update(func(txn *badger.Txn) error {
		fk := toFileKey(key, f.Path)
		byteKey := []byte(fk.String())

		// todo: keep origin value
		err = txn.Set(byteKey, nil)
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

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *badgerDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContext, _ context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	err = d.db.Update(func(txn *badger.Txn) error {
		fk := toFileKey(key, f.Path)
		byteKey := []byte(fk.String())

		// todo: keep origin value
		err = txn.Set(byteKey, nil)
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

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *badgerDriver) CreateWorkspace(wc *object.WorkspaceConfig, _ context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}

	revInfo := object.NewRevInfo(wc.RevHash)
	revInfoStr, err := json.Marshal(revInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal rev info: %w", err)
	}

	err = d.db.Update(func(txn *badger.Txn) error {
		byteKey := []byte(ToRevKey(key).String())
		err = txn.Set(byteKey, revInfoStr)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
