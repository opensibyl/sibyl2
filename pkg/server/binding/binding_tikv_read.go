package binding

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tikv/client-go/v2/kv"
)

func (t *tikvDriver) readRawRevs() ([]*revKV, error) {
	snapshot := t.client.GetSnapshot(math.MaxUint64)
	keyByte := []byte(revEndPrefix)
	iter, err := snapshot.Iter(keyByte, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	ret := make([]*revKV, 0)
	for iter.Valid() {
		kv := &revKV{}
		kv.k = parseRevKey(string(iter.Key()))
		v := &object.RevInfo{}
		err = json.Unmarshal(iter.Value(), v)
		if err != nil {
			return nil, err
		}
		kv.v = v

		// ok
		ret = append(ret, kv)
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
		wc, err := WorkspaceConfigFromKey(eachRev.k.Hash)
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

	// order by create time desc
	sort.Slice(revs, func(i, j int) bool {
		return revs[i].v.CreateTime > revs[j].v.CreateTime
	})

	for _, eachRev := range revs {
		wc, err := WorkspaceConfigFromKey(eachRev.k.Hash)
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

func (t *tikvDriver) ReadRevInfo(wc *object.WorkspaceConfig, ctx context.Context) (*object.RevInfo, error) {
	snapshot := t.client.GetSnapshot(math.MaxUint64)

	ret := &object.RevInfo{}
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	rk := ToRevKey(key)
	keyByte := []byte(rk.String())

	iter, err := snapshot.Iter(keyByte, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	for iter.Valid() {
		err := json.Unmarshal(iter.Value(), ret)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal rev info: %w", err)
		}
		// only check the first one
		break
	}
	return ret, nil
}
