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

func (t *tikvDriver) ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, _ context.Context) ([]*sibyl2.FunctionContext, error) {
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

	searchResult := make([]*sibyl2.FunctionContext, 0)

	prefix := []byte(ToRevKey(key).ToScanPrefix())

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.Valid() {
		k := string(iter.Key())
		flag := "funcctx|"
		if strings.Contains(k, flag) {
			rawFunc := iter.Value()
			for rk, rv := range compiledRule {
				v := gjson.GetBytes(rawFunc, rk)
				if rv.MatchString(v.String()) {
					f := &sibyl2.FunctionContext{}
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
	if err != nil {
		return nil, err
	}
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
