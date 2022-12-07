package object

import (
	"fmt"

	"github.com/pkg/errors"
)

type DriverType string

const DtNeo4j DriverType = "NEO4J"
const DtInMemory DriverType = "IN_MEMORY"
const FlagWcKeySplit = "|,,|"

/*
WorkspaceConfig

as an infra lib, it will not assume what kind of repo you used.

just two fields:
- repoId: unique id of your repo, no matter git or svn, even appId.
- revHash: unique id of your version.
*/
type WorkspaceConfig struct {
	RepoId  string `json:"repoId"`
	RevHash string `json:"revHash"`
}

func (wc *WorkspaceConfig) Verify() error {
	// all the fields should be filled
	if wc == nil || wc.RepoId == "" || wc.RevHash == "" {
		return errors.Errorf("workspace config verify error: %v", wc)
	}
	return nil
}

func (wc *WorkspaceConfig) Key() (string, error) {
	if err := wc.Verify(); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", wc.RepoId, FlagWcKeySplit, wc.RevHash), nil
}
