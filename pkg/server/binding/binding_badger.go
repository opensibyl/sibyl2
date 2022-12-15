package binding

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

/*
storage:
- rev|<hash>:
- rev_<hash>_file|<hash>:
- rev_<hash>_file_<hash>_func|<hash>: func details map
- rev_<hash>_file_<hash>_funcctx|<hash>: func ctx details map

mean:
- |: type def end
- _: connection
*/

type badgerDriver struct {
	db *badger.DB
}

const revPrefix = "rev|"

type revKey struct {
	hash string
}

func (r *revKey) String() string {
	return revPrefix + r.hash
}

func (r *revKey) ToScanPrefix() string {
	return "rev_" + r.hash + "_"
}

func toRevKey(revHash string) *revKey {
	return &revKey{revHash}
}

func parseRevKey(raw string) *revKey {
	return &revKey{strings.TrimPrefix(raw, revPrefix)}
}

type fileKey struct {
	revHash  string
	fileHash string
}

func (f *fileKey) String() string {
	return fmt.Sprintf("rev_%s_file|%s", f.revHash, f.fileHash)
}

func (f *fileKey) ToScanPrefix() string {
	return fmt.Sprintf("rev_%s_file_%s_", f.revHash, f.fileHash)
}

func toFileKey(revHash string, fileHash string) *fileKey {
	return &fileKey{revHash, fileHash}
}

type funcKey struct {
	revHash  string
	fileHash string
	funcHash string
}

func toFuncKey(revHash string, fileHash string, funcHash string) *funcKey {
	return &funcKey{revHash, fileHash, funcHash}
}

func (f *funcKey) String() string {
	return fmt.Sprintf("rev_%s_file_%s_func|%s", f.revHash, f.fileHash, f.funcHash)
}

type funcCtxKey struct {
	revHash  string
	fileHash string
	funcHash string
}

func toFuncCtxKey(revHash string, fileHash string, funcHash string) *funcCtxKey {
	return &funcCtxKey{revHash, fileHash, funcHash}
}

func (f *funcCtxKey) String() string {
	return fmt.Sprintf("rev_%s_file_%s_funcctx|%s", f.revHash, f.fileHash, f.funcHash)
}

func (d *badgerDriver) InitDriver(_ context.Context) error {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *badgerDriver) DeferDriver() error {
	return d.db.Close()
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
			eachFuncKey := toFuncKey(fk.revHash, fk.fileHash, eachFunc.GetSignature())
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

		eachFuncKey := toFuncCtxKey(fk.revHash, fk.fileHash, f.GetSignature())
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
	err = d.db.Update(func(txn *badger.Txn) error {
		byteKey := []byte(toRevKey(key).String())
		err = txn.Set(byteKey, nil)
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
	ret := make([]string, 0)
	for _, eachRev := range revs {
		wc, err := WorkspaceConfigFromKey(eachRev.hash)
		if err != nil {
			return nil, err
		}
		ret = append(ret, wc.RepoId)
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
		wc, err := WorkspaceConfigFromKey(eachRev.hash)
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
	rk := toRevKey(key)
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

func (d *badgerDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, _ context.Context) ([]*sibyl2.FunctionWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	fk := toFileKey(key, path)

	searchResult := make([]*sibyl2.FunctionWithPath, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefixStr := fk.ToScanPrefix()
		prefix := []byte(prefixStr)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			f := &sibyl2.FunctionWithPath{}
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
			searchResult = append(searchResult, f)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (d *badgerDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := toRevKey(key)
	var ret *sibyl2.FunctionWithPath
	err = d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)

		defer it.Close()
		prefixStr := rk.ToScanPrefix() + "file_"
		prefix := []byte(prefixStr)
		shouldContain := "func|" + signature
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

func (d *badgerDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	functions, err := d.ReadFunctions(wc, path, ctx)
	if err != nil {
		return nil, err
	}

	searchResult := make([]*sibyl2.FunctionWithPath, 0)
	for _, each := range functions {
		if each.Span.ContainAnyLine(lines...) {
			searchResult = append(searchResult, each)
		}
	}
	return searchResult, nil
}

func (d *badgerDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, _ context.Context) (*sibyl2.FunctionContext, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := toRevKey(key)
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

func (d *badgerDriver) UpdateRevProperties(wc *object.WorkspaceConfig, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) UpdateFileProperties(wc *object.WorkspaceConfig, path string, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) UpdateFuncProperties(wc *object.WorkspaceConfig, signature string, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	rk := toRevKey(key)
	itself := []byte(rk.String())
	sons := []byte(rk.ToScanPrefix())

	err = d.db.DropPrefix(itself, sons)
	if err != nil {
		return err
	}
	return nil
}

func newBadgerDriver() Driver {
	return &badgerDriver{}
}

func (d *badgerDriver) GetType() object.DriverType {
	return object.DtBadger
}

func initBadgerDriver(_ object.ExecuteConfig) Driver {
	return newBadgerDriver()
}
