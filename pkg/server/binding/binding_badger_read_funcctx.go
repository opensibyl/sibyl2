package binding

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
)

func (d *badgerDriver) ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionContext, error) {
	functions, err := d.ReadFunctionsWithLines(wc, path, lines, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*sibyl2.FunctionContext, 0)
	for _, eachFunc := range functions {
		functionContext, err := d.ReadFunctionContextWithSignature(wc, eachFunc.GetSignature(), ctx)
		if err != nil {
			return nil, err
		}
		ret = append(ret, functionContext)
	}
	return ret, nil
}

func (d *badgerDriver) ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, _ context.Context) ([]*sibyl2.FunctionContext, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}
	compiledRule := make(map[string]*regexp.Regexp)
	for k, v := range rule {
		newRegex, err := regexp.Compile(v)
		if err != nil {
			return nil, err
		}
		compiledRule[k] = newRegex
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	prefix := []byte(ToRevKey(key).ToScanPrefix())

	searchResult := make([]*sibyl2.FunctionContext, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			k := string(it.Item().Key())
			flag := "funcctx|"
			if !strings.Contains(k, flag) {
				continue
			}
			err = it.Item().Value(func(val []byte) error {
				for rk, rv := range compiledRule {
					v := gjson.GetBytes(val, rk)
					if rv.MatchString(v.String()) {
						c := &sibyl2.FunctionContext{}
						err = json.Unmarshal(val, c)
						if err != nil {
							return err
						}
						searchResult = append(searchResult, c)
					}
				}
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

func (d *badgerDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, _ context.Context) (*sibyl2.FunctionContext, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	var ret *sibyl2.FunctionContext
	err = d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)

		defer it.Close()
		prefixStr := rk.ToScanPrefix() + "file_"
		prefix := []byte(prefixStr)
		shouldContain := "funcctx|" + signature
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := string(item.Key())
			if strings.Contains(k, shouldContain) {
				err := item.Value(func(val []byte) error {
					ret, err = sibyl2.Json2FuncCtx(val)
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
