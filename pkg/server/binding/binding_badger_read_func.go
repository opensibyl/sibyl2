package binding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
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
		prefix := []byte(ToRevKey(key).ToFileScanPrefix())
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

func (d *badgerDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, _ context.Context) ([]*object.FunctionWithSignature, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	fk := toFileKey(key, path)

	searchResult := make([]*object.FunctionWithSignature, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefixStr := fk.ToFuncScanPrefix()
		prefix := []byte(prefixStr)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			f := &object.FunctionWithSignature{}
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

			f.Signature = f.GetSignature()
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

func (d *badgerDriver) ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, _ context.Context) ([]*object.FunctionWithSignature, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}

	searchResult := make([]*object.FunctionWithSignature, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(ToRevKey(key).ToFileScanPrefix())
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
				f := &object.FunctionWithSignature{}
				err = json.Unmarshal(val, f)
				if err != nil {
					return err
				}
				f.Signature = f.GetSignature()
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

func (d *badgerDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, _ context.Context) (*object.FunctionWithSignature, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	ret := &object.FunctionWithSignature{}
	err = d.db.View(func(txn *badger.Txn) error {
		k := rk.ToFuncPtrPrefix() + signature
		item, err := txn.Get([]byte(k))
		if err != nil {
			return fmt.Errorf("func not found: %v, %v, %w", wc, signature, err)
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

func (d *badgerDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FunctionWithSignature, error) {
	functions, err := d.ReadFunctions(wc, path, ctx)
	if err != nil {
		return nil, err
	}

	searchResult := make([]*object.FunctionWithSignature, 0)
	for _, each := range functions {
		if each.Span.ContainAnyLine(lines...) {
			searchResult = append(searchResult, each)
		}
	}
	return searchResult, nil
}

func (d *badgerDriver) ReadFunctionsWithTag(wc *object.WorkspaceConfig, tag sibyl2.FuncTag, ctx context.Context) ([]string, error) {
	rule := make(Rule)
	requiredTags := strings.Split(tag, ";")
	rule["tags"] = func(s string) bool {
		// json string list
		tags := make([]string, 0)
		err := json.Unmarshal([]byte(s), &tags)
		if err != nil {
			// should not happen
			return false
		}
		// all the tags should exist
		for _, each := range requiredTags {
			if !slices.Contains(tags, each) {
				return false
			}
		}
		return true
	}
	functionWithTags, err := d.ReadFunctionsWithRule(wc, rule, ctx)
	if err != nil {
		return nil, err
	}
	// return only signature for avoiding huge io cost
	ret := make([]string, 0, len(functionWithTags))
	for _, each := range functionWithTags {
		ret = append(ret, each.GetSignature())
	}
	return ret, nil
}
