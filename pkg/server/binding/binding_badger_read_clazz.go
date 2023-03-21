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

func (d *badgerDriver) ReadClasses(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*object.ClazzServiceDTO, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	curFileKey := toFileKey(key, path)

	searchResult := make([]*object.ClazzServiceDTO, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefixStr := curFileKey.ToClazzScanPrefix()
		prefix := []byte(prefixStr)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			c := &object.ClazzServiceDTO{}
			err = it.Item().Value(func(val []byte) error {
				err = json.Unmarshal(val, c)
				if err != nil {
					return fmt.Errorf("unmarshal class failed: %w", err)
				}
				return nil
			})
			if err != nil {
				return err
			}

			c.Path = path
			searchResult = append(searchResult, c)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (d *badgerDriver) ReadClassesWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.ClazzServiceDTO, error) {
	classes, err := d.ReadClasses(wc, path, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*object.ClazzServiceDTO, 0)
	for _, each := range classes {
		if each.GetSpan().ContainAnyLine(lines...) {
			ret = append(ret, each)
		}
	}
	return ret, nil
}

func (d *badgerDriver) ReadClassesWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.ClazzServiceDTO, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	prefix := []byte(ToRevKey(key).ToFileScanPrefix())

	searchResult := make([]*object.ClazzServiceDTO, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			k := string(it.Item().Key())
			if !strings.Contains(k, clazzEndPrefix) {
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
				c := &object.ClazzServiceDTO{}
				err = json.Unmarshal(val, c)
				if err != nil {
					return err
				}
				searchResult = append(searchResult, c)
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
