package binding

import (
	"fmt"
	"strings"
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

const revPrefix = "rev|"

type revKey struct {
	Hash string
}

func (r *revKey) String() string {
	return revPrefix + r.Hash
}

func (r *revKey) ToScanPrefix() string {
	return "rev_" + r.Hash + "_"
}

func ToRevKey(revHash string) *revKey {
	return &revKey{revHash}
}

func parseRevKey(raw string) *revKey {
	return &revKey{strings.TrimPrefix(raw, revPrefix)}
}

type fileKey struct {
	RevHash  string
	FileHash string
}

func (f *fileKey) String() string {
	return fmt.Sprintf("rev_%s_file|%s", f.RevHash, f.FileHash)
}

func (f *fileKey) ToScanPrefix() string {
	return fmt.Sprintf("rev_%s_file_%s_", f.RevHash, f.FileHash)
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
	return fmt.Sprintf("rev_%s_file_%s_func|%s", f.revHash, f.fileHash, f.funcHash)
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
	return fmt.Sprintf("rev_%s_file_%s_clazz|%s", c.revHash, c.fileHash, c.clazzHash)
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
	return fmt.Sprintf("rev_%s_file_%s_funcctx|%s", f.revHash, f.fileHash, f.funcHash)
}
