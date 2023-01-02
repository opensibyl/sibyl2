package service

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/ext"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

type FuncDiffRet = map[string][]*sibyl2.FunctionWithPath
type FuncCtxDiffRet = map[string][]*sibyl2.FunctionContext
type ClazzDiffRet = map[string][]*sibyl2.ClazzWithPath

// @Summary func diff query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   diff query string true "unified diff"
// @Produce json
// @Success 200 {object} FuncDiffRet
// @Router  /api/v1/func/diff [get]
// @Tags EXTRAS
func HandleFunctionsDiffQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	diff := c.Query("diff")

	// make sure workspace valid
	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// validate format
	affectedMap, err := ext.Unified2Affected([]byte(diff))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ret, err := handleFuncDiffMap(wc, affectedMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, ret)
	return
}

// @Summary func ctx diff query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   diff query string true "unified diff"
// @Produce json
// @Success 200 {object} FuncCtxDiffRet
// @Router  /api/v1/funcctx/diff [get]
// @Tags EXTRAS
func HandleFunctionCtxDiffQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	diff := c.Query("diff")

	// make sure workspace valid
	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// validate format
	affectedMap, err := ext.Unified2Affected([]byte(diff))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ret, err := handleFuncCtxDiffMap(wc, affectedMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, ret)
	return
}

// @Summary clazz diff query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   diff query string true "unified diff"
// @Produce json
// @Success 200 {object} ClazzDiffRet
// @Router  /api/v1/clazz/diff [get]
// @Tags EXTRAS
func HandleClazzDiffQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	diff := c.Query("diff")

	// make sure workspace valid
	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// validate format
	affectedMap, err := ext.Unified2Affected([]byte(diff))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ret, err := handleClazzDiffMap(wc, affectedMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, ret)
	return
}

func handleFuncDiffMap(wc *object.WorkspaceConfig, affectedMap ext.AffectedLineMap) (FuncDiffRet, error) {
	ret := make(map[string][]*sibyl2.FunctionWithPath)
	ctx := context.Background()
	for fileName, lines := range affectedMap {
		functions, err := sharedDriver.ReadFunctionsWithLines(wc, fileName, lines, ctx)
		if err != nil {
			return nil, err
		}
		ret[fileName] = functions
	}
	return ret, nil
}

func handleFuncCtxDiffMap(wc *object.WorkspaceConfig, affectedMap ext.AffectedLineMap) (FuncCtxDiffRet, error) {
	ret := make(map[string][]*sibyl2.FunctionContext)
	ctx := context.Background()
	for fileName, lines := range affectedMap {
		functions, err := sharedDriver.ReadFunctionContextsWithLines(wc, fileName, lines, ctx)
		if err != nil {
			return nil, err
		}
		ret[fileName] = functions
	}
	return ret, nil
}

func handleClazzDiffMap(wc *object.WorkspaceConfig, affectedMap ext.AffectedLineMap) (ClazzDiffRet, error) {
	ret := make(map[string][]*sibyl2.ClazzWithPath)
	ctx := context.Background()
	for fileName, lines := range affectedMap {
		classes, err := sharedDriver.ReadClassesWithLines(wc, fileName, lines, ctx)
		if err != nil {
			return nil, err
		}

		ret[fileName] = classes
	}
	return ret, nil
}
