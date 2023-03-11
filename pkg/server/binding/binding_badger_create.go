package binding

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"golang.org/x/exp/slices"
)

func (d *badgerDriver) CreateClazzFile(wc *object.WorkspaceConfig, c *extractor.ClazzFileResult, _ context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}

	err = d.db.Update(func(txn *badger.Txn) error {
		fk := toFileKey(key, c.Path)
		byteKey := []byte(fk.String())

		// write file key
		err = txn.Set(byteKey, nil)
		if err != nil {
			return err
		}

		for _, eachClazz := range c.Units {
			// write class fact
			eachClazzKey := toClazzKey(fk.RevHash, fk.FileHash, eachClazz.GetSignature())
			eachClazzWithSignature := &object.ClazzWithSignature{
				ClazzWithPath: &sibyl2.ClazzWithPath{
					Clazz: eachClazz,
					Path:  c.Path,
				},
				Signature: eachClazz.GetSignature(),
			}
			eachClazzV, err := json.Marshal(eachClazzWithSignature)
			if err != nil {
				continue
			}
			err = txn.Set([]byte(eachClazzKey.String()), eachClazzV)
			if err != nil {
				return err
			}

			// write func ptr
			ptrKey := []byte(eachClazzKey.StringWithoutFile())
			factListBytes, err := txn.Get(ptrKey)
			switch err {
			case badger.ErrKeyNotFound:
				sl := []string{eachClazzKey.String()}
				bytes, err := json.Marshal(sl)
				if err != nil {
					return err
				}
				err = txn.Set(ptrKey, bytes)
				if err != nil {
					return err
				}
			case nil:
				// one signature can map more than one
				factList := make([]string, 0)
				err = factListBytes.Value(func(val []byte) error {
					err := json.Unmarshal(val, &factList)
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					return err
				}
			default:
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

		// write file key
		err = txn.Set(byteKey, nil)
		if err != nil {
			return err
		}

		for _, eachFunc := range f.Units {
			// write func fact
			eachFuncKey := toFuncKey(fk.RevHash, fk.FileHash, eachFunc.GetSignature())
			eachFuncWithSignature := &object.FunctionWithSignature{
				FunctionWithTag: &sibyl2.FunctionWithTag{
					FunctionWithPath: &sibyl2.FunctionWithPath{
						Function: eachFunc,
						Path:     f.Path,
					},
					Tags: make([]string, 0),
				},
				Signature: eachFunc.GetSignature(),
			}
			eachFuncV, err := json.Marshal(eachFuncWithSignature)
			if err != nil {
				continue
			}
			err = txn.Set([]byte(eachFuncKey.String()), eachFuncV)
			if err != nil {
				return err
			}

			// write func ptr
			ptrKey := []byte(eachFuncKey.StringWithoutFile())
			factListBytes, err := txn.Get(ptrKey)
			switch err {
			case badger.ErrKeyNotFound:
				sl := []string{eachFuncKey.String()}
				bytes, err := json.Marshal(sl)
				if err != nil {
					return err
				}
				err = txn.Set(ptrKey, bytes)
				if err != nil {
					return err
				}
			case nil:
				// one signature can map more than one
				factList := make([]string, 0)
				err = factListBytes.Value(func(val []byte) error {
					err := json.Unmarshal(val, &factList)
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					return err
				}
			default:
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

func (d *badgerDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContextSlim, _ context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	err = d.db.Update(func(txn *badger.Txn) error {
		fk := toFileKey(key, f.Path)
		byteKey := []byte(fk.String())

		// write file key
		err = txn.Set(byteKey, nil)
		if err != nil {
			return err
		}

		// write the fact
		eachFuncKey := toFuncCtxKey(fk.RevHash, fk.FileHash, f.GetSignature())
		eachFuncV, err := json.Marshal(f)
		if err != nil {
			return err
		}
		funcFactKey := []byte(eachFuncKey.String())
		err = txn.Set(funcFactKey, eachFuncV)
		if err != nil {
			return err
		}

		// write the ptr
		ptrKey := []byte(eachFuncKey.StringWithoutFile())
		factListBytes, err := txn.Get(ptrKey)
		switch err {
		case badger.ErrKeyNotFound:
			sl := []string{string(funcFactKey)}
			bytes, err := json.Marshal(sl)
			if err != nil {
				return err
			}
			err = txn.Set(ptrKey, bytes)
			if err != nil {
				return err
			}
		case nil:
			// one signature can map more than one
			factList := make([]string, 0)
			err = factListBytes.Value(func(val []byte) error {
				err := json.Unmarshal(val, &factList)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		default:
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

func (d *badgerDriver) CreateFuncTag(wc *object.WorkspaceConfig, signature string, tag string, ctx context.Context) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		// request inside the transaction
		f, err := d.ReadFunctionWithSignature(wc, signature, ctx)
		if err != nil {
			return fmt.Errorf("failed to read func with signature: %w", err)
		}

		if slices.Contains(f.Tags, tag) {
			// duplicated
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

		byteKey := []byte(curFuncKey.String())
		err = txn.Set(byteKey, newFuncBytes)
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
