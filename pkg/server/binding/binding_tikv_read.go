package binding

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
	"github.com/tikv/client-go/v2/kv"
)

func (t *tikvDriver) readRawRevs() ([]*revKey, error) {
	snapshot := t.client.GetSnapshot(math.MaxUint64)
	keyByte := []byte(revPrefix)
	iter, err := snapshot.Iter(keyByte, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	ret := make([]*revKey, 0)
	for iter.Valid() {
		ret = append(ret, parseRevKey(string(iter.Key())))
		err := iter.Next()
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (t *tikvDriver) ReadRepos(_ context.Context) ([]string, error) {
	revs, err := t.readRawRevs()
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

func (t *tikvDriver) ReadRevs(repoId string, _ context.Context) ([]string, error) {
	revs, err := t.readRawRevs()
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

func (t *tikvDriver) ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	searchResult := make([]string, 0)

	txn := t.client.GetSnapshot(math.MaxUint64)

	prefixStr := rk.ToScanPrefix() + "file|"
	prefix := []byte(prefixStr)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()
	for iter.Valid() {
		k := iter.Key()
		searchResult = append(searchResult, strings.TrimPrefix(string(k), prefixStr))
		err := iter.Next()
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (t *tikvDriver) ReadFunctionSignaturesWithRegex(wc *object.WorkspaceConfig, regex string, _ context.Context) ([]string, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}

	compiled, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	searchResult := make([]string, 0)
	txn := t.client.GetSnapshot(math.MaxUint64)
	prefix := []byte(ToRevKey(key).ToScanPrefix())
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))

	for iter.Valid() {
		k := string(iter.Key())
		flag := "func|"
		if strings.Contains(k, flag) {
			_, after, _ := strings.Cut(k, flag)
			if compiled.MatchString(after) {
				searchResult = append(searchResult, after)
			}
		}
		err := iter.Next()
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}
	return searchResult, nil
}

func (t *tikvDriver) ReadClasses(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	fk := toFileKey(key, path)

	searchResult := make([]*sibyl2.ClazzWithPath, 0)

	prefixStr := fk.ToScanPrefix() + "clazz|"
	prefix := []byte(prefixStr)

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()

	for iter.Valid() {
		c := &sibyl2.ClazzWithPath{}
		err := json.Unmarshal(iter.Value(), c)
		if err != nil {
			return nil, err
		}
		c.Path = path
		searchResult = append(searchResult, c)
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}

	return searchResult, nil
}

func (t *tikvDriver) ReadClassesWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	classes, err := t.ReadClasses(wc, path, ctx)
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

func (t *tikvDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	fk := toFileKey(key, path)

	searchResult := make([]*sibyl2.FunctionWithPath, 0)

	prefixStr := fk.ToScanPrefix() + "func|"
	prefix := []byte(prefixStr)

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()

	for iter.Valid() {
		f := &sibyl2.FunctionWithPath{}
		err := json.Unmarshal(iter.Value(), f)
		if err != nil {
			return nil, err
		}
		f.Path = path
		searchResult = append(searchResult, f)
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}

	return searchResult, nil
}

func (t *tikvDriver) ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
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

	prefix := []byte(ToRevKey(key).ToScanPrefix())

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()

	for iter.Valid() {
		k := string(iter.Key())
		flag := "func|"
		if strings.Contains(k, flag) {
			rawFunc := iter.Value()
			for rk, rv := range compiledRule {
				v := gjson.GetBytes(rawFunc, rk)
				if rv.MatchString(v.String()) {
					f := &sibyl2.FunctionWithPath{}
					err = json.Unmarshal(rawFunc, f)
					if err != nil {
						return nil, err
					}
					searchResult = append(searchResult, f)
				}
			}
		}
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	return searchResult, nil
}

func (t *tikvDriver) ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionContext, error) {
	functions, err := t.ReadFunctionsWithLines(wc, path, lines, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*sibyl2.FunctionContext, 0)
	for _, eachFunc := range functions {
		functionContext, err := t.ReadFunctionContextWithSignature(wc, eachFunc.GetSignature(), ctx)
		if err != nil {
			return nil, err
		}
		ret = append(ret, functionContext)
	}
	return ret, nil
}

func (t *tikvDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, _ context.Context) (*sibyl2.FunctionWithPath, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)

	prefixStr := rk.ToScanPrefix() + "file_"
	prefix := []byte(prefixStr)
	shouldContain := "func|" + signature

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()

	var ret *sibyl2.FunctionWithPath
	for iter.Valid() {
		k := string(iter.Key())
		if strings.Contains(k, shouldContain) {
			err := json.Unmarshal(iter.Value(), ret)
			if err != nil {
				return nil, err
			}
			fp, _, _ := strings.Cut(strings.TrimPrefix(k, prefixStr), shouldContain)
			ret.Path = fp
			return ret, nil
		}
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	// did not find anything
	return nil, nil
}

func (t *tikvDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	functions, err := t.ReadFunctions(wc, path, ctx)
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

func (t *tikvDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContext, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)

	prefixStr := rk.ToScanPrefix() + "file_"
	prefix := []byte(prefixStr)
	shouldContain := "funcctx|" + signature

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	defer iter.Close()

	for iter.Valid() {
		k := string(iter.Key())
		if strings.Contains(k, shouldContain) {
			funcCtx, err := sibyl2.Json2FuncCtx(iter.Value())
			if err != nil {
				return nil, err
			}
			// break scan
			return funcCtx, nil
		}
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	// did not find anything
	return nil, nil
}
