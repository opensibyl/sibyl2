package sibyl2

import (
	"context"

	"github.com/williamfzc/sibyl2/pkg/extractor"
)

type memDriver struct {
	storage *InMemoryStorage
}

func (m memDriver) GetType() DriverType {
	return DtInMemory
}

func (m memDriver) CreateFuncFile(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) CreateFuncContext(wc *WorkspaceConfig, f *FunctionContext, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) CreateWorkspace(wc *WorkspaceConfig, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) ReadRepos(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) ReadRevs(repoId string, ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) ReadFiles(wc *WorkspaceConfig, ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) ReadFunctions(wc *WorkspaceConfig, path string, ctx context.Context) ([]*FunctionWithPath, error) {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) ReadFunctionWithSignature(wc *WorkspaceConfig, signature string, ctx context.Context) (*FunctionWithPath, error) {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) ReadFunctionsWithLines(wc *WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*FunctionWithPath, error) {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) ReadFunctionContextWithSignature(wc *WorkspaceConfig, signature string, ctx context.Context) (*FunctionContext, error) {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) UpdateRepoProperties(repoId string, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) UpdateRevProperties(wc *WorkspaceConfig, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) UpdateFileProperties(wc *WorkspaceConfig, path string, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) UpdateFuncProperties(wc *WorkspaceConfig, signature string, k string, v any, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m memDriver) DeleteWorkspace(wc *WorkspaceConfig, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
