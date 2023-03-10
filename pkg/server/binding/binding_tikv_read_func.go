package binding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
	"github.com/tikv/client-go/v2/kv"
	"golang.org/x/exp/slices"
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
	prefix := []byte(ToRevKey(key).ToFileScanPrefix())
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))

	for iter.Valid() {
		k := string(iter.Key())
		if strings.Contains(k, funcEndPrefix) {
			_, after, _ := strings.Cut(k, funcEndPrefix)
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

func (t *tikvDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, _ context.Context) (*object.FunctionWithSignature, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)

	prefixStr := rk.ToFileScanPrefix()
	prefix := []byte(prefixStr)
	shouldContain := flagConnect + funcEndPrefix + signature

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	ret := &object.FunctionWithSignature{}
	for iter.Valid() {
		k := string(iter.Key())
		if strings.Contains(k, shouldContain) {
			err := json.Unmarshal(iter.Value(), ret)
			if err != nil {
				return nil, err
			}
			fp, _, _ := strings.Cut(strings.TrimPrefix(k, prefixStr), shouldContain)
			ret.Path = fp
			ret.Signature = ret.GetSignature()
			return ret, nil
		}
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	// did not find anything
	return nil, fmt.Errorf("func not found: %v, %v", wc, signature)
}

func (t *tikvDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FunctionWithSignature, error) {
	functions, err := t.ReadFunctions(wc, path, ctx)
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

func (t *tikvDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*object.FunctionWithSignature, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	fk := toFileKey(key, path)

	searchResult := make([]*object.FunctionWithSignature, 0)

	prefixStr := fk.ToFuncScanPrefix()
	prefix := []byte(prefixStr)

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.Valid() {
		f := &object.FunctionWithSignature{}
		err := json.Unmarshal(iter.Value(), f)
		if err != nil {
			return nil, err
		}
		f.Signature = f.GetSignature()
		f.Path = path
		searchResult = append(searchResult, f)
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}

	return searchResult, nil
}

func (t *tikvDriver) ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.FunctionWithSignature, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}

	searchResult := make([]*object.FunctionWithSignature, 0)

	prefix := []byte(ToRevKey(key).ToFileScanPrefix())

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.Valid() {
		k := string(iter.Key())
		if strings.Contains(k, funcEndPrefix) {
			rawFunc := iter.Value()
			for rk, verify := range rule {
				v := gjson.GetBytes(rawFunc, rk)
				if !verify(v.String()) {
					// failed and ignore this item
					goto nextIter
				}
			}
			// all the rules passed
			f := &object.FunctionWithSignature{}
			err = json.Unmarshal(rawFunc, f)
			if err != nil {
				return nil, err
			}
			f.Signature = f.GetSignature()
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

func (t *tikvDriver) ReadFunctionsWithTag(wc *object.WorkspaceConfig, tag sibyl2.FuncTag, ctx context.Context) ([]string, error) {
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
	functionWithTags, err := t.ReadFunctionsWithRule(wc, rule, ctx)
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
