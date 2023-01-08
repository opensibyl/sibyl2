package binding

import (
	"context"

	object "github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tikv/client-go/v2/kv"
)

func (t *tikvDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	rk := ToRevKey(key)
	itself := []byte(rk.String())
	sons := []byte(rk.ToScanPrefix())

	_, err = t.client.DeleteRange(ctx, itself, kv.PrefixNextKey(itself), 1)
	if err != nil {
		return err
	}
	_, err = t.client.DeleteRange(ctx, sons, kv.PrefixNextKey(sons), 1)
	if err != nil {
		return err
	}

	return nil
}
