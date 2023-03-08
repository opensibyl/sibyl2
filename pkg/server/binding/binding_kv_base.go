package binding

import (
	"strings"

	"github.com/opensibyl/sibyl2/pkg/server/object"
)

/*
storage:
- rev|<hash>:
- rev_<hash>_file|<hash>:
- rev_<hash>_file_<hash>_func|<hash>: func details map
- rev_<hash>_file_<hash>_funcctx|<hash>: func ctx details map

mean:
- |: type def end
- _: connection
*/

const (
	revEndPrefix     = "rev" + flagEnd
	revSearchPrefix  = "rev" + flagConnect
	fileEndPrefix    = "file" + flagEnd
	fileSearchPrefix = "file" + flagConnect
	ptrSearchPrefix  = "ptr" + flagConnect
	funcEndPrefix    = "func" + flagEnd
	clazzEndPrefix   = "clazz" + flagEnd
	funcctxEndPrefix = "funcctx" + flagEnd
	flagConnect      = "_"
	flagEnd          = "|"
)

type revKey struct {
	Hash string
}

func (r *revKey) String() string {
	return revEndPrefix + r.Hash
}

func (r *revKey) ToScanPrefix() string {
	return revSearchPrefix + r.Hash + flagConnect
}

func (r *revKey) ToFileScanPrefix() string {
	return r.ToScanPrefix() + fileSearchPrefix
}

func (r *revKey) ToFuncCtxPtrPrefix() string {
	return r.ToScanPrefix() + ptrSearchPrefix + funcctxEndPrefix
}

func ToRevKey(revHash string) *revKey {
	return &revKey{revHash}
}

func parseRevKey(raw string) *revKey {
	return &revKey{strings.TrimPrefix(raw, revEndPrefix)}
}

type revKV struct {
	k *revKey
	v *object.RevInfo
}

type fileKey struct {
	RevHash  string
	FileHash string
}

func (f *fileKey) String() string {
	return revSearchPrefix + f.RevHash + flagConnect + fileEndPrefix + f.FileHash
}

func (f *fileKey) ToScanPrefix() string {
	return revSearchPrefix + f.RevHash + flagConnect + fileSearchPrefix + f.FileHash + flagConnect
}

func (f *fileKey) ToClazzScanPrefix() string {
	return f.ToScanPrefix() + clazzEndPrefix
}

func (f *fileKey) ToFuncScanPrefix() string {
	return f.ToScanPrefix() + funcEndPrefix
}

func toFileKey(revHash string, fileHash string) *fileKey {
	return &fileKey{revHash, fileHash}
}

type funcKey struct {
	revHash  string
	fileHash string
	funcHash string
}

func toFuncKey(revHash string, fileHash string, funcHash string) *funcKey {
	return &funcKey{revHash, fileHash, funcHash}
}

func (f *funcKey) String() string {
	return revSearchPrefix + f.revHash + flagConnect + fileSearchPrefix + f.fileHash + flagConnect + funcEndPrefix + f.funcHash
}

type clazzKey struct {
	revHash   string
	fileHash  string
	clazzHash string
}

func toClazzKey(revHash string, fileHash string, clazzHash string) *clazzKey {
	return &clazzKey{revHash, fileHash, clazzHash}
}

func (c *clazzKey) String() string {
	return revSearchPrefix + c.revHash + flagConnect + fileSearchPrefix + c.fileHash + flagConnect + clazzEndPrefix + c.clazzHash
}

type funcCtxKey struct {
	revHash  string
	fileHash string
	funcHash string
}

func toFuncCtxKey(revHash string, fileHash string, funcHash string) *funcCtxKey {
	return &funcCtxKey{revHash, fileHash, funcHash}
}

func (f *funcCtxKey) String() string {
	return revSearchPrefix + f.revHash + flagConnect + fileSearchPrefix + f.fileHash + flagConnect + funcctxEndPrefix + f.funcHash
}

func (f *funcCtxKey) StringWithoutFile() string {
	return revSearchPrefix + f.revHash + flagConnect + ptrSearchPrefix + funcctxEndPrefix + f.funcHash
}
