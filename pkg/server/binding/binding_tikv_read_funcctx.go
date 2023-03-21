package binding

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
	"github.com/tikv/client-go/v2/kv"
)

func (t *tikvDriver) ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.FuncCtxServiceDTO, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}

	searchResult := make([]*object.FuncCtxServiceDTO, 0)

	prefix := []byte(ToRevKey(key).ToFileScanPrefix())

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.Valid() {
		k := string(iter.Key())
		if strings.Contains(k, funcctxEndPrefix) {
			rawFunc := iter.Value()
			for rk, verify := range rule {
				v := gjson.GetBytes(rawFunc, rk)
				if !verify(v.String()) {
					// failed and ignore this item
					goto nextIter
				}
			}
			// all the rules passed
			f := &object.FuncCtxServiceDTO{}
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

func (t *tikvDriver) ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FuncCtxServiceDTO, error) {
	functions, err := t.ReadFunctionsWithLines(wc, path, lines, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*object.FuncCtxServiceDTO, 0)
	for _, eachFunc := range functions {
		functionContext, err := t.ReadFunctionContextWithSignature(wc, eachFunc.GetSignature(), ctx)
		if err != nil {
			return nil, err
		}
		ret = append(ret, functionContext)
	}
	return ret, nil
}

func (t *tikvDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*object.FuncCtxServiceDTO, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)

	prefixStr := rk.ToFileScanPrefix()
	prefix := []byte(prefixStr)
	shouldContain := flagConnect + funcctxEndPrefix + signature

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.Valid() {
		k := string(iter.Key())
		if strings.Contains(k, shouldContain) {
			f := &object.FuncCtxServiceDTO{}
			err = json.Unmarshal(iter.Value(), f)
			if err != nil {
				return nil, err
			}
			// break scan
			return f, nil
		}
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	// did not find anything
	return nil, fmt.Errorf("func not found: %v, %v", wc, signature)
}
