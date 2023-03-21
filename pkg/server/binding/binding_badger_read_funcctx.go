package binding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
)

func (d *badgerDriver) ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FunctionContextSlim, error) {
	functions, err := d.ReadFunctionsWithLines(wc, path, lines, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*object.FunctionContextSlim, 0)
	for _, eachFunc := range functions {
		functionContext, err := d.ReadFunctionContextWithSignature(wc, eachFunc.GetSignature(), ctx)
		if err != nil {
			return nil, err
		}
		ret = append(ret, functionContext)
	}
	return ret, nil
}

func (d *badgerDriver) ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.FunctionContextSlim, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	prefix := []byte(ToRevKey(key).ToFileScanPrefix())

	searchResult := make([]*object.FunctionContextSlim, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			k := string(it.Item().Key())
			if !strings.Contains(k, funcctxEndPrefix) {
				continue
			}
			err = it.Item().Value(func(val []byte) error {
				for rk, verify := range rule {
					v := gjson.GetBytes(val, rk)
					if !verify(v.String()) {
						// failed and ignore this item
						return nil
					}
				}
				// all the rules passed
				f := &object.FunctionContextSlim{}
				err = json.Unmarshal(val, f)
				if err != nil {
					return err
				}
				searchResult = append(searchResult, f)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return searchResult, nil

}

func (d *badgerDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, _ context.Context) (*object.FunctionContextSlim, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	ret := &object.FunctionContextSlim{}
	err = d.db.View(func(txn *badger.Txn) error {
		k := rk.ToFuncCtxPtrPrefix() + signature
		item, err := txn.Get([]byte(k))
		if err != nil {
			return fmt.Errorf("func ctx not found: %v, %v, %w", wc, signature, err)
		}
		mappingList := make([]string, 0)
		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, &mappingList)
			if err != nil {
				return err
			}
			return nil
		})
		// currently we will only return the first one
		for _, each := range mappingList {
			realItem, err := txn.Get([]byte(each))
			if err != nil {
				return err
			}
			err = realItem.Value(func(val []byte) error {
				err = json.Unmarshal(val, &ret)
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

		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}
