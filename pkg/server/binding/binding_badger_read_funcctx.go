package binding

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
)

func (d *badgerDriver) ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionContextSlim, error) {
	functions, err := d.ReadFunctionsWithLines(wc, path, lines, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*sibyl2.FunctionContextSlim, 0)
	for _, eachFunc := range functions {
		functionContext, err := d.ReadFunctionContextWithSignature(wc, eachFunc.GetSignature(), ctx)
		if err != nil {
			return nil, err
		}
		ret = append(ret, functionContext)
	}
	return ret, nil
}

func (d *badgerDriver) ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*sibyl2.FunctionContextSlim, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	prefix := []byte(ToRevKey(key).ToScanPrefix())

	searchResult := make([]*sibyl2.FunctionContextSlim, 0)
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
				f := &sibyl2.FunctionContextSlim{}
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

func (d *badgerDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, _ context.Context) (*sibyl2.FunctionContextSlim, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	ret := &sibyl2.FunctionContextSlim{}
	err = d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)

		defer it.Close()
		prefixStr := rk.ToScanPrefix() + fileSearchPrefix
		prefix := []byte(prefixStr)
		shouldContain := funcctxEndPrefix + signature
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := string(item.Key())
			if strings.Contains(k, shouldContain) {
				err := item.Value(func(val []byte) error {
					err = json.Unmarshal(val, ret)
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					return err
				}
				// break scan
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}
