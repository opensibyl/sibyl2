package sibyl2

import (
	"context"
	"errors"

	"github.com/williamfzc/sibyl2/pkg/extractor"
)

// key: workspace config 's key
// todo: thread-safe
type InMemoryStorage = map[string]*revUnit
type revUnit = map[string]*fileStorage

func newRevUnit() *revUnit {
	ret := make(map[string]*fileStorage)
	return &ret
}

type fileStorage struct {
	path      string
	functions *extractor.FunctionFileResult
	symbols   *extractor.SymbolFileResult
}

func newFileStorage(path string) *fileStorage {
	return &fileStorage{
		path: path,
	}
}

// this mem driver and storage should only be used in debug
// it will increase memory cost with no limitation
type memDriver struct {
	InMemoryStorage
}

func (m *memDriver) isWcExisted(wc *WorkspaceConfig) bool {
	key, err := wc.Key()
	if err != nil {
		return false
	}
	_, ok := m.InMemoryStorage[key]
	return ok
}

func (m *memDriver) getAllWorkspaceConfig(ctx context.Context) ([]*WorkspaceConfig, error) {
	ret := make([]*WorkspaceConfig, 0, len(m.InMemoryStorage))
	for each := range m.InMemoryStorage {
		wc, err := WorkspaceConfigFromKey(each)
		if err != nil {
			return nil, err
		}
		ret = append(ret, wc)
	}
	return ret, nil
}

func (m *memDriver) getRevUnit(wc *WorkspaceConfig, ctx context.Context) (*revUnit, error) {
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	v, ok := m.InMemoryStorage[key]
	if !ok {
		return nil, errors.New("no such workspace config: " + key)
	}
	return v, nil
}

func NewMemDriver() Driver {
	storage := make(map[string]*revUnit)
	return &memDriver{
		storage,
	}
}

func (m *memDriver) GetType() DriverType {
	return DtInMemory
}

func (m *memDriver) CreateFuncFile(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	unit := m.InMemoryStorage[key]
	if unit == nil {
		unit = newRevUnit()
		m.InMemoryStorage[key] = unit
	}

	pathUnit := (*unit)[f.Path]
	if pathUnit == nil {
		pathUnit = newFileStorage(f.Path)
		(*unit)[f.Path] = pathUnit
	}

	pathUnit.functions = f
	return nil
}

func (m *memDriver) CreateFuncContext(wc *WorkspaceConfig, f *FunctionContext, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *memDriver) CreateWorkspace(wc *WorkspaceConfig, ctx context.Context) error {
	if m.isWcExisted(wc) {
		return nil
	}
	key, err := wc.Key()
	if err != nil {
		return err
	}
	m.InMemoryStorage[key] = newRevUnit()
	return nil
}

func (m *memDriver) ReadRepos(ctx context.Context) ([]string, error) {
	ret := make([]string, 0, len(m.InMemoryStorage))
	wcs, err := m.getAllWorkspaceConfig(ctx)
	if err != nil {
		return nil, err
	}

	for _, each := range wcs {
		ret = append(ret, each.RepoId)
	}
	return ret, nil
}

func (m *memDriver) ReadRevs(repoId string, ctx context.Context) ([]string, error) {
	ret := make([]string, 0, len(m.InMemoryStorage))
	wcs, err := m.getAllWorkspaceConfig(ctx)
	if err != nil {
		return nil, err
	}

	for _, each := range wcs {
		if repoId != each.RepoId {
			continue
		}
		ret = append(ret, each.RevHash)
	}
	return ret, nil
}

func (m *memDriver) ReadFiles(wc *WorkspaceConfig, ctx context.Context) ([]string, error) {
	unit, err := m.getRevUnit(wc, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(*unit))
	for each := range *unit {
		ret = append(ret, each)
	}
	return ret, nil
}

func (m *memDriver) ReadFunctions(wc *WorkspaceConfig, path string, ctx context.Context) ([]*FunctionWithPath, error) {
	unit, err := m.getRevUnit(wc, ctx)
	if err != nil {
		return nil, err
	}

	file, ok := (*unit)[path]
	if !ok {
		return nil, errors.New("")
	}

	ret := make([]*FunctionWithPath, 0, len(*unit))
	for _, eachFunc := range file.functions.Units {
		fwp := &FunctionWithPath{
			Function: eachFunc,
			Path:     path,
			Language: file.functions.Language,
		}
		ret = append(ret, fwp)
	}
	return ret, nil
}

func (m *memDriver) ReadFunctionWithSignature(wc *WorkspaceConfig, signature string, ctx context.Context) (*FunctionWithPath, error) {
	files, err := m.ReadFiles(wc, ctx)
	if err != nil {
		return nil, err
	}

	for _, eachPath := range files {
		functions, err := m.ReadFunctions(wc, eachPath, ctx)
		if err != nil {
			return nil, err
		}
		for _, eachFunc := range functions {
			if eachFunc.GetSignature() == signature {
				return eachFunc, nil
			}
		}
	}
	return nil, errors.New("no func found: " + signature)
}

func (m *memDriver) ReadFunctionsWithLines(wc *WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*FunctionWithPath, error) {
	files, err := m.ReadFiles(wc, ctx)
	if err != nil {
		return nil, err
	}

	var ret []*FunctionWithPath
	for _, eachPath := range files {
		functions, err := m.ReadFunctions(wc, eachPath, ctx)
		if err != nil {
			return nil, err
		}
		for _, eachFunc := range functions {
			if eachFunc.GetSpan().ContainAnyLine(lines...) {
				ret = append(ret, eachFunc)
			}
		}
	}
	return ret, nil
}

func (m *memDriver) ReadFunctionContextWithSignature(wc *WorkspaceConfig, signature string, ctx context.Context) (*FunctionContext, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}

func (m *memDriver) UpdateRepoProperties(repoId string, k string, v any, ctx context.Context) error {
	return errors.New("NOT IMPLEMENTED")
}

func (m *memDriver) UpdateRevProperties(wc *WorkspaceConfig, k string, v any, ctx context.Context) error {
	return errors.New("NOT IMPLEMENTED")
}

func (m *memDriver) UpdateFileProperties(wc *WorkspaceConfig, path string, k string, v any, ctx context.Context) error {
	return errors.New("NOT IMPLEMENTED")
}

func (m *memDriver) UpdateFuncProperties(wc *WorkspaceConfig, signature string, k string, v any, ctx context.Context) error {
	return errors.New("NOT IMPLEMENTED")
}

func (m *memDriver) DeleteWorkspace(wc *WorkspaceConfig, ctx context.Context) error {
	key, err := wc.Key()
	if err != nil {
		return err
	}
	delete(m.InMemoryStorage, key)
	return nil
}
