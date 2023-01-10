package binding

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
	"github.com/tikv/client-go/v2/kv"
)

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
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.Valid() {
		c := &sibyl2.ClazzWithPath{}
		err = json.Unmarshal(iter.Value(), c)
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

func (t *tikvDriver) ReadClassesWithRule(wc *object.WorkspaceConfig, rule Rule, _ context.Context) ([]*sibyl2.ClazzWithPath, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	key, err := wc.Key()
	if err != nil {
		return nil, err
	}

	searchResult := make([]*sibyl2.ClazzWithPath, 0)

	prefix := []byte(ToRevKey(key).ToScanPrefix())

	txn := t.client.GetSnapshot(math.MaxUint64)
	iter, err := txn.Iter(prefix, kv.PrefixNextKey(prefix))
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for iter.Valid() {
		k := string(iter.Key())
		flag := "clazz|"
		if strings.Contains(k, flag) {
			rawClazz := iter.Value()
			for rk, verify := range rule {
				v := gjson.GetBytes(rawClazz, rk)
				if !verify(v.String()) {
					// failed and ignore this item
					goto nextIter
				}
			}
			// all the rules passed
			c := &sibyl2.ClazzWithPath{}
			err = json.Unmarshal(rawClazz, c)
			if err != nil {
				return nil, err
			}
			searchResult = append(searchResult, c)
		}

	nextIter:
		err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	return searchResult, nil
}
