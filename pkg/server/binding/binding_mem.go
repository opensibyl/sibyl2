package binding

import (
	"context"
	"errors"
	"sync"

	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/extractor"
)

type InMemoryStorage struct {
	data map[string]*revUnit // key: workspace config 's key
	l    *sync.RWMutex
}

func newStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]*revUnit),
		l:    new(sync.RWMutex),
	}
}

type revUnit struct {
	data map[string]*fileStorage
	l    *sync.RWMutex
}

func newRevUnit() *revUnit {
	return &revUnit{
		data: make(map[string]*fileStorage),
		l:    new(sync.RWMutex),
	}
}

type fileStorage struct {
	path      string
	functions *extractor.FunctionFileResult
	symbols   *extractor.SymbolFileResult
	l         *sync.RWMutex
}

func newFileStorage(path string) *fileStorage {
	return &fileStorage{
		path: path,
		l:    new(sync.RWMutex),
	}
}

// this mem driver and storage should only be used in debug
// it will increase memory cost with no limitation
type memDriver struct {
	*InMemoryStorage
}

func (m *memDriver) isWcExisted(wc *WorkspaceConfig) bool {
	key, err := wc.Key()
	if err != nil {
		return false
	}
	m.l.RLock()
	defer m.l.RUnlock()
	_, ok := m.InMemoryStorage.data[key]
	return ok
}

func (m *memDriver) getAllWorkspaceConfig(ctx context.Context) ([]*WorkspaceConfig, error) {
	m.l.RLock()
	defer m.l.RUnlock()
	ret := make([]*WorkspaceConfig, 0, len(m.InMemoryStorage.data))
	for each := range m.InMemoryStorage.data {
		wc, err := WorkspaceConfigFromKey(each)
		if err != nil {
			return nil, err
		}
		ret = append(ret, wc)
	}
	return ret, nil
}

func (m *memDriver) getRevUnit(wc *WorkspaceConfig, ctx context.Context) (*revUnit, error) {
	m.l.RLock()
	defer m.l.RUnlock()
	key, err := wc.Key()
	if err != nil {
		return nil, err
	}
	v, ok := m.InMemoryStorage.data[key]
	if !ok {
		return nil, errors.New("no such workspace config: " + key)
	}
	return v, nil
}

func NewMemDriver() Driver {
	return &memDriver{
		newStorage(),
	}
}

func (m *memDriver) GetType() DriverType {
	return DtInMemory
}

func (m *memDriver) CreateFuncFile(wc *WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	m.l.Lock()
	defer m.l.Unlock()
	key, err := wc.Key()
	if err != nil {
		return err
	}
	unit := m.InMemoryStorage.data[key]
	if unit == nil {
		unit = newRevUnit()
		m.InMemoryStorage.data[key] = unit
	}

	unit.l.Lock()
	defer unit.l.Unlock()
	pathUnit := unit.data[f.Path]
	if pathUnit == nil {
		pathUnit = newFileStorage(f.Path)
		unit.data[f.Path] = pathUnit
	}

	pathUnit.functions = f
	return nil
}

func (m *memDriver) CreateFuncContext(wc *WorkspaceConfig, f *sibyl2.FunctionContext, ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (m *memDriver) CreateWorkspace(wc *WorkspaceConfig, ctx context.Context) error {
	m.l.Lock()
	defer m.l.Unlock()
	if m.isWcExisted(wc) {
		return nil
	}
	key, err := wc.Key()
	if err != nil {
		return err
	}
	m.InMemoryStorage.data[key] = newRevUnit()
	return nil
}

func (m *memDriver) ReadRepos(ctx context.Context) ([]string, error) {
	m.l.RLock()
	defer m.l.RUnlock()
	ret := make([]string, 0, len(m.InMemoryStorage.data))
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
	m.l.RLock()
	defer m.l.RUnlock()
	ret := make([]string, 0, len(m.InMemoryStorage.data))
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
	m.l.RLock()
	defer m.l.RUnlock()
	unit, err := m.getRevUnit(wc, ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(unit.data))
	for each := range unit.data {
		ret = append(ret, each)
	}
	return ret, nil
}

func (m *memDriver) ReadFunctions(wc *WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	m.l.RLock()
	defer m.l.RUnlock()
	unit, err := m.getRevUnit(wc, ctx)
	if err != nil {
		return nil, err
	}

	file, ok := unit.data[path]
	if !ok {
		return nil, errors.New("")
	}

	ret := make([]*sibyl2.FunctionWithPath, 0, len(unit.data))
	for _, eachFunc := range file.functions.Units {
		fwp := &sibyl2.FunctionWithPath{
			Function: eachFunc,
			Path:     path,
			Language: file.functions.Language,
		}
		ret = append(ret, fwp)
	}
	return ret, nil
}

func (m *memDriver) ReadFunctionWithSignature(wc *WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithPath, error) {
	m.l.RLock()
	defer m.l.RUnlock()
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

func (m *memDriver) ReadFunctionsWithLines(wc *WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
	m.l.RLock()
	defer m.l.RUnlock()
	files, err := m.ReadFiles(wc, ctx)
	if err != nil {
		return nil, err
	}

	var ret []*sibyl2.FunctionWithPath
	for _, eachPath := range files {
		if path != eachPath {
			continue
		}

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

func (m *memDriver) ReadFunctionContextWithSignature(wc *WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContext, error) {
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
	m.l.Lock()
	defer m.l.Unlock()
	key, err := wc.Key()
	if err != nil {
		return err
	}
	delete(m.InMemoryStorage.data, key)
	return nil
}
