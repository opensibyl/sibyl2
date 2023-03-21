package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"golang.org/x/exp/slices"
)

type FunctionContextChain struct {
	*object.FunctionContextSlim
	CallChains        *ContextTree `json:"callChains"`
	ReverseCallChains *ContextTree `json:"reverseCallChains"`
}

/*
ContextTree

- avoiding duplicated chains
- easily handled by frontend/dashboard
*/
type ContextTree struct {
	Content  string         `json:"content"`
	Children []*ContextTree `json:"children"`
}

func (t *ContextTree) AddChain(chain []string) {
	for _, part := range chain {
		for _, each := range t.Children {
			// node already existed
			if each.Content == part {
				each.AddChain(chain[1:])
				return
			}
		}
		// not existed yet
		newSubTree := &ContextTree{
			Content:  part,
			Children: nil,
		}
		t.Children = append(t.Children, newSubTree)
		newSubTree.AddChain(chain[1:])
		return
	}
}

// @Summary funcctx reverse chain query
// @Param   repo      query string true "repo"
// @Param   rev       query string true "rev"
// @Param   signature query string true "signature"
// @Param   depth     query int    true "depth"
// @Produce json
// @Success 200 {object} FunctionContextChain
// @Router  /api/v1/signature/funcctx/rchain [get]
// @Tags    SignatureQuery
func HandleSignatureFuncctxReverseChain(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	signature := c.Query("signature")
	depth := c.Query("depth")

	depthNum, err := strconv.Atoi(depth)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid depth: %w", err))
		return
	}

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
	reverseCallChains := &ContextTree{}
	err = searchReverseCallChains(wc, f.GetSignature(), make([]string, 0), depthNum, reverseCallChains)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("failed to get reverse call: %v", reverseCallChains))
		return
	}

	fc := &FunctionContextChain{
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

func searchReverseCallChains(wc *object.WorkspaceConfig, startPoint string, curChain []string, depthLimit int, chains *ContextTree) error {
	calls, err := readRCalls(wc, startPoint)
	if err != nil {
		return err
	}

	if len(calls) == 0 || len(curChain) > depthLimit {
		// end point, save this chain to tree
		chains.AddChain(curChain)
		return nil
	}

	// continue
	for _, eachCallSignature := range calls {
		if slices.Contains(curChain, eachCallSignature) {
			// loop call
			continue
		}
		err = searchReverseCallChains(wc, eachCallSignature, append(curChain, eachCallSignature), depthLimit, chains)
		if err != nil {
			return err
		}
	}
	return nil
}

// @Summary funcctx chain query
// @Param   repo      query string true "repo"
// @Param   rev       query string true "rev"
// @Param   signature query string true "signature"
// @Param   depth     query int    true "depth"
// @Produce json
// @Success 200 {object} FunctionContextChain
// @Router  /api/v1/signature/funcctx/chain [get]
// @Tags    SignatureQuery
func HandleSignatureFuncctxChain(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	signature := c.Query("signature")
	depth := c.Query("depth")

	depthNum, err := strconv.Atoi(depth)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid depth: %w", err))
		return
	}

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
	reverseCallChains := &ContextTree{}
	err = searchCallChains(wc, f.GetSignature(), make([]string, 0), depthNum, reverseCallChains)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("failed to get reverse call: %v", reverseCallChains))
		return
	}

	fc := &FunctionContextChain{
		FunctionContextSlim: f,
		CallChains:          reverseCallChains,
	}

	c.JSON(http.StatusOK, fc)
}

func readCalls(wc *object.WorkspaceConfig, signature string) ([]string, error) {
	curPtr, err := sharedDriver.ReadFunctionContextWithSignature(wc, signature, sharedContext)
	if err != nil {
		return nil, err
	}
	return curPtr.Calls, nil
}

func searchCallChains(wc *object.WorkspaceConfig, startPoint string, curChain []string, depthLimit int, chains *ContextTree) error {
	calls, err := readCalls(wc, startPoint)
	if err != nil {
		return err
	}

	if len(calls) == 0 || len(curChain) > depthLimit {
		// end point, save this chain to tree
		chains.AddChain(curChain)
		return nil
	}

	// continue
	for _, eachCallSignature := range calls {
		if slices.Contains(curChain, eachCallSignature) {
			// loop call
			continue
		}
		err = searchCallChains(wc, eachCallSignature, append(curChain, eachCallSignature), depthLimit, chains)
		if err != nil {
			return err
		}
	}
	return nil
}
