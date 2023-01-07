package binding

import (
	"context"

	"github.com/opensibyl/sibyl2/pkg/server/object"
)

func (d *badgerDriver) DeleteWorkspace(wc *object.WorkspaceConfig, _ context.Context) error {
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
