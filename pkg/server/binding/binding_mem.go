package binding

import (
	"context"
	"errors"
	"sync"

	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/server/object"
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
	path             string
	functions        *extractor.FunctionFileResult
	symbols          *extractor.SymbolFileResult
	functionContexts map[string]*sibyl2.FunctionContext
	l                *sync.RWMutex
}

func newFileStorage(path string) *fileStorage {
	return &fileStorage{
		path:             path,
		functionContexts: make(map[string]*sibyl2.FunctionContext),
		l:                new(sync.RWMutex),
	}
}

// this mem driver and storage should only be used in debug
// it will increase memory cost with no limitation
type memDriver struct {
	*InMemoryStorage
}

func (m *memDriver) isWcExisted(wc *object.WorkspaceConfig) bool {
	key, err := wc.Key()
	if err != nil {
		return false
	}
	_, ok := m.InMemoryStorage.data[key]
	return ok
}

func (m *memDriver) getAllWorkspaceConfig(ctx context.Context) ([]*object.WorkspaceConfig, error) {
	ret := make([]*object.WorkspaceConfig, 0, len(m.InMemoryStorage.data))
	for each := range m.InMemoryStorage.data {
		wc, err := WorkspaceConfigFromKey(each)
		if err != nil {
			return nil, err
		}
		ret = append(ret, wc)
	}
	return ret, nil
}

func (m *memDriver) getRevUnit(wc *object.WorkspaceConfig, _ context.Context) (*revUnit, error) {
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

func newMemDriver() Driver {
	return &memDriver{
		newStorage(),
	}
}

func (m *memDriver) GetType() object.DriverType {
	return object.DtInMemory
}

func (m *memDriver) InitDriver(_ context.Context) error {
	// do nothing
	return nil
}

func (m *memDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, _ context.Context) error {
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

func (m *memDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContext, _ context.Context) error {
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
	if pathUnit != nil {
		// overwrite whatever
		pathUnit.functionContexts[f.GetSignature()] = f
	}
	return nil
}

func (m *memDriver) CreateWorkspace(wc *object.WorkspaceConfig, _ context.Context) error {
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

func (m *memDriver) ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error) {
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

func (m *memDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
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

func (m *memDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionWithPath, error) {
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

func (m *memDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionWithPath, error) {
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

func (m *memDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContext, error) {
	unit, err := m.getRevUnit(wc, ctx)
	if err != nil {
		return nil, err
	}

	for _, eachFile := range unit.data {
		v, ok := eachFile.functionContexts[signature]
		if !ok {
			continue
		}
		return v, nil
	}
	return nil, errors.New("function context not found")
}

func (m *memDriver) UpdateRevProperties(*object.WorkspaceConfig, string, any, context.Context) error {
	return errors.New("NOT IMPLEMENTED")
}

func (m *memDriver) UpdateFileProperties(*object.WorkspaceConfig, string, string, any, context.Context) error {
	return errors.New("NOT IMPLEMENTED")
}

func (m *memDriver) UpdateFuncProperties(*object.WorkspaceConfig, string, string, any, context.Context) error {
	return errors.New("NOT IMPLEMENTED")
}

func (m *memDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	m.l.Lock()
	defer m.l.Unlock()
	key, err := wc.Key()
	if err != nil {
		return err
	}
	delete(m.InMemoryStorage.data, key)
	return nil
}
