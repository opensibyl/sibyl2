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

func (d *badgerDriver) ReadFunctionSignaturesWithRegex(wc *object.WorkspaceConfig, regex string, _ context.Context) ([]string, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}

	compiled, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	searchResult := make([]string, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte(ToRevKey(key).ToScanPrefix())
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			k := string(it.Item().Key())
			if strings.Contains(k, funcEndPrefix) {
				_, after, _ := strings.Cut(k, funcEndPrefix)
				if compiled.MatchString(after) {
					searchResult = append(searchResult, after)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (d *badgerDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithTag, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	fk := toFileKey(key, path)

	searchResult := make([]*sibyl2.FunctionWithTag, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefixStr := fk.ToFuncScanPrefix()
		prefix := []byte(prefixStr)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			f := &sibyl2.FunctionWithTag{}
			err := it.Item().Value(func(val []byte) error {
				err := json.Unmarshal(val, f)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}

			f.Path = path
			searchResult = append(searchResult, f)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (d *badgerDriver) ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*sibyl2.FunctionWithTag, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}

	searchResult := make([]*sibyl2.FunctionWithTag, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(ToRevKey(key).ToScanPrefix())
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			k := string(it.Item().Key())
			if !strings.Contains(k, funcEndPrefix) {
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
				f := &sibyl2.FunctionWithTag{}
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

func (d *badgerDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithTag, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	ret := &sibyl2.FunctionWithTag{}
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefixStr := rk.ToFileScanPrefix()
		prefix := []byte(prefixStr)
		shouldContain := flagConnect + funcEndPrefix + signature
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := string(item.Key())
			if strings.Contains(k, shouldContain) {
				err := item.Value(func(val []byte) error {
					err := json.Unmarshal(val, ret)
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					return err
				}
				fp, _, _ := strings.Cut(strings.TrimPrefix(k, prefixStr), shouldContain)
				ret.Path = fp
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

func (d *badgerDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithTag, error) {
	functions, err := d.ReadFunctions(wc, path, ctx)
	if err != nil {
		return nil, err
	}

	searchResult := make([]*sibyl2.FunctionWithTag, 0)
	for _, each := range functions {
		if each.Span.ContainAnyLine(lines...) {
			searchResult = append(searchResult, each)
		}
	}
	return searchResult, nil
}
