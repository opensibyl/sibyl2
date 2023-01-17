package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"golang.org/x/exp/slices"
)

type FunctionContextReverseChain struct {
	*sibyl2.FunctionContextSlim
	ReverseCallChains [][]string `json:"reverseCallChains"`
}

// @Summary funcctx reverse chain query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   signature   query string true  "signature"
// @Produce json
// @Success 200 {object} FunctionContextReverseChain
// @Router  /api/v1/funcctx/rchain/with/signature [get]
// @Tags SignatureQuery
func HandleFunctionContextReverseChainQueryWithSignature(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	signature := c.Query("signature")

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	f, err := sharedDriver.ReadFunctionContextWithSignature(wc, signature, sharedContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	// query chains
	// todo: should move it to binding layer? it can be slow
	reverseCallChains, err := searchReverseCallChains(wc, f.GetSignature(), make([]string, 0))
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("failed to get reverse call: %v", reverseCallChains))
		return
	}

	fc := &FunctionContextReverseChain{
		FunctionContextSlim: f,
		ReverseCallChains:   reverseCallChains,
	}

	c.JSON(http.StatusOK, fc)
}

func readRCalls(wc *object.WorkspaceConfig, signature string) ([]string, error) {
	curPtr, err := sharedDriver.ReadFunctionContextWithSignature(wc, signature, sharedContext)
	if err != nil {
		return nil, err
	}
	return curPtr.ReverseCalls, nil
}

func searchReverseCallChains(wc *object.WorkspaceConfig, startPoint string, curChain []string) ([][]string, error) {
	chains := make([][]string, 0)

	calls, err := readRCalls(wc, startPoint)
	if err != nil {
		return nil, err
	}

	if len(calls) == 0 {
		// end point
		chains = append(chains, curChain)
		return chains, nil
	}

	// continue
	for _, eachCallSignature := range calls {
		if !slices.Contains(curChain, eachCallSignature) {
			subChains, err := searchReverseCallChains(wc, eachCallSignature, append(curChain, eachCallSignature))
			if err != nil {
				return nil, err
			}
			chains = append(chains, subChains...)
		}
	}
	return chains, nil
}
