package binding

import (
	"context"

	"github.com/opensibyl/sibyl2/pkg/server/object"
)

func (d *mongoDriver) UpdateRevProperties(wc *object.WorkspaceConfig, k string, v any, ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) UpdateFileProperties(wc *object.WorkspaceConfig, path string, k string, v any, ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) UpdateFuncProperties(wc *object.WorkspaceConfig, signature string, k string, v any, ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}
