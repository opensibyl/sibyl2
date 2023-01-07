package binding

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
)

type badgerDriver struct {
	db     *badger.DB
	config object.ExecuteConfig
}

func (d *badgerDriver) InitDriver(_ context.Context) error {
	var dbInst *badger.DB
	var err error

	switch d.config.DbType {
	case object.DriverTypeInMemory:
		dbInst, err = badger.Open(badger.DefaultOptions("").WithInMemory(true))
	case object.DriverTypeBadger:
		core.Log.Infof("trying to open: %s", d.config.BadgerPath)
		dbInst, err = badger.Open(badger.DefaultOptions(d.config.BadgerPath))
	default:
		core.Log.Errorf("db type %v invalid", d.config.DbType)
	}
	if err != nil {
		return err
	}
	d.db = dbInst

	return nil
}

func (d *badgerDriver) DeferDriver() error {
	return d.db.Close()
}

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
	err = d.db.Update(func(txn *badger.Txn) error {
		byteKey := []byte(ToRevKey(key).String())
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

func (d *badgerDriver) ReadClasses(wc *object.WorkspaceConfig, path string, _ context.Context) ([]*sibyl2.ClazzWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	fk := toFileKey(key, path)

	searchResult := make([]*sibyl2.ClazzWithPath, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefixStr := fk.ToScanPrefix() + "clazz|"
		prefix := []byte(prefixStr)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			c := &sibyl2.ClazzWithPath{}
			err := it.Item().Value(func(val []byte) error {
				err := json.Unmarshal(val, c)
				if err != nil {
					return err
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

func (d *badgerDriver) ReadClassesWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	classes, err := d.ReadClasses(wc, path, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*sibyl2.ClazzWithPath, 0)
	for _, each := range classes {
		if each.GetSpan().ContainAnyLine(lines...) {
			ret = append(ret, each)
		}
	}
	return ret, nil
}

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
			flag := "func|"
			if strings.Contains(k, flag) {
				_, after, _ := strings.Cut(k, flag)
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
		prefixStr := fk.ToScanPrefix() + "func|"
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

func (d *badgerDriver) ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, _ context.Context) ([]*sibyl2.FunctionWithPath, error) {
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

	searchResult := make([]*sibyl2.FunctionWithPath, 0)
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(ToRevKey(key).ToScanPrefix())
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			k := string(it.Item().Key())
			flag := "func|"
			if !strings.Contains(k, flag) {
				continue
			}
			err = it.Item().Value(func(val []byte) error {
				for rk, rv := range compiledRule {
					v := gjson.GetBytes(val, rk)
					if rv.MatchString(v.String()) {
						f := &sibyl2.FunctionWithPath{}
						err = json.Unmarshal(val, f)
						if err != nil {
							return err
						}
						searchResult = append(searchResult, f)
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

func (d *badgerDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, _ context.Context) (*sibyl2.FunctionWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	ret := &sibyl2.FunctionWithPath{}
	err = d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefixStr := rk.ToScanPrefix() + "file_"
		prefix := []byte(prefixStr)
		shouldContain := "_func|" + signature
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

func (d *badgerDriver) UpdateRevProperties(wc *object.WorkspaceConfig, k string, v any, ctx context.Context) error {
	// TODO implement me
	return errors.New("implement me")
}

func (d *badgerDriver) UpdateFileProperties(wc *object.WorkspaceConfig, path string, k string, v any, ctx context.Context) error {
	// TODO implement me
	return errors.New("implement me")
}

func (d *badgerDriver) UpdateFuncProperties(wc *object.WorkspaceConfig, signature string, k string, v any, ctx context.Context) error {
	// TODO implement me
	return errors.New("implement me")
}

func (d *badgerDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	rk := ToRevKey(key)
	itself := []byte(rk.String())
	sons := []byte(rk.ToScanPrefix())

	err = d.db.DropPrefix(itself, sons)
	if err != nil {
		return err
	}
	return nil
}

func (d *badgerDriver) GetType() object.DriverType {
	return object.DriverTypeBadger
}

func initBadgerDriver(config object.ExecuteConfig) Driver {
	return &badgerDriver{nil, config}
}
