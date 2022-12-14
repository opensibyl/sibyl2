package binding

import (
	"context"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

/*
workspace -> rev hash id
rev hash id + file path = file hash id
file hash id + func signature = func hash

storage:
- rev_<hash>: [file_<hash>, ]
- file_<hash>: [func_<hash>, ]
- func_<hash>: func details map
*/

type badgerDriver struct {
	db *badger.DB
}

func (d *badgerDriver) InitDriver(_ context.Context) error {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *badgerDriver) DeferDriver() error {
	return d.db.Close()
}

func (d *badgerDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContext, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) ReadRepos(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) ReadRevs(repoId string, ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithPath, error) {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContext, error) {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) UpdateRevProperties(wc *object.WorkspaceConfig, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) UpdateFileProperties(wc *object.WorkspaceConfig, path string, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) UpdateFuncProperties(wc *object.WorkspaceConfig, signature string, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d *badgerDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func newBadgerDriver() Driver {
	return &badgerDriver{}
}

func (d *badgerDriver) GetType() object.DriverType {
	return object.DtBadger
}
