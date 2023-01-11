package binding

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

func (d *badgerDriver) readRawRevs() ([]*revKey, error) {
	ret := make([]*revKey, 0)
	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte(revPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			ret = append(ret, parseRevKey(string(k)))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (d *badgerDriver) ReadRepos(_ context.Context) ([]string, error) {
	revs, err := d.readRawRevs()
	if err != nil {
		return nil, err
	}
	m := make(map[string]struct{}, 0)
	for _, eachRev := range revs {
		wc, err := WorkspaceConfigFromKey(eachRev.Hash)
		if err != nil {
			return nil, err
		}
		m[wc.RepoId] = struct{}{}
	}

	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret, nil
}

func (d *badgerDriver) ReadRevs(repoId string, _ context.Context) ([]string, error) {
	revs, err := d.readRawRevs()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, eachRev := range revs {
		wc, err := WorkspaceConfigFromKey(eachRev.Hash)
		if err != nil {
			return nil, err
		}
		if wc.RepoId == repoId {
			ret = append(ret, wc.RevHash)
		}
	}
	return ret, nil
}

func (d *badgerDriver) ReadFiles(wc *object.WorkspaceConfig, _ context.Context) ([]string, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	searchResult := make([]string, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)

		defer it.Close()
		prefixStr := rk.ToScanPrefix() + "file|"
		prefix := []byte(prefixStr)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			searchResult = append(searchResult, strings.TrimPrefix(string(k), prefixStr))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (d *badgerDriver) ReadRevInfo(wc *object.WorkspaceConfig, ctx context.Context) (*object.RevInfo, error) {
	ret := &object.RevInfo{}
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)

	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(rk.String())
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				err := json.Unmarshal(val, ret)
				if err != nil {
					return fmt.Errorf("failed to unmarshal rev info: %w", err)
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
	return ret, nil
}
