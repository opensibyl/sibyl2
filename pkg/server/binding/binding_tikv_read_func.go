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
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	ret := &sibyl2.FunctionWithPath{}
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

func (t *tikvDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, _ context.Context) ([]*sibyl2.FunctionWithPath, error) {
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
	if err != nil {
		return nil, err
	}
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

func (t *tikvDriver) ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, _ context.Context) ([]*sibyl2.FunctionWithPath, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}

	searchResult := make([]*sibyl2.FunctionWithPath, 0)

	prefix := []byte(ToRevKey(key).ToScanPrefix())

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.Valid() {
		k := string(iter.Key())
		flag := "func|"
		if strings.Contains(k, flag) {
			rawFunc := iter.Value()
			for rk, verify := range rule {
				v := gjson.GetBytes(rawFunc, rk)
				if !verify(v.String()) {
					// failed and ignore this item
					goto nextIter
				}
			}
			// all the rules passed
			f := &sibyl2.FunctionWithPath{}
			err = json.Unmarshal(rawFunc, f)
			if err != nil {
				return nil, err
			}
			searchResult = append(searchResult, f)
		}

	nextIter:
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	return searchResult, nil
}
