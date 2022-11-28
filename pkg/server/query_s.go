package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/server/binding"
)

// @Summary repo query
// @Produce json
// @Success 200
// @Router  /api/v1/repo [get]
func HandleRepoQuery(c *gin.Context) {
	repos, err := sharedDriver.ReadRepos(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, repos)
}

// @Summary rev query
// @Param   repo query string true "rev search by repo"
// @Produce json
// @Success 200 {array} string
// @Router  /api/v1/rev [get]
func HandleRevQuery(c *gin.Context) {
	repo := c.Query("repo")
	revs, err := sharedDriver.ReadRevs(repo, context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, revs)
}

// @Summary file query
// @Param   repo query string true "repo"
// @Param   rev  query string true "rev"
// @Produce json
// @Success 200 {array} string
// @Router  /api/v1/file [get]
func HandleFileQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	wc := &binding.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	files, err := sharedDriver.ReadFiles(wc, context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, files)
}

// @Summary func query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   file  query string true  "file"
// @Param   lines query string false "specific lines"
// @Produce json
// @Success 200 {array} FunctionWithSignature
// @Router  /api/v1/func [get]
func HandleFunctionsQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	file := c.Query("file")
	lines := c.Query("lines")

	ret, err := handleFunctionQuery(repo, rev, file, lines)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, ret)
}

func handleFunctionQuery(repo string, rev string, file string, lines string) ([]*FunctionWithSignature, error) {
	wc := &binding.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		return nil, err
	}
	var functions []*sibyl2.FunctionWithPath
	var err error
	if lines == "" {
		functions, err = sharedDriver.ReadFunctions(wc, file, context.TODO())
	} else {
		linesStrList := strings.Split(lines, ",")
		var lineNums = make([]int, 0, len(linesStrList))
		for _, each := range linesStrList {
			num, err := strconv.Atoi(each)
			if err != nil {
				return nil, err
			}
			lineNums = append(lineNums, num)
		}
		functions, err = sharedDriver.ReadFunctionsWithLines(wc, file, lineNums, context.TODO())
	}
	if err != nil {
		return nil, err
	}

	// export signature
	ret := make([]*FunctionWithSignature, 0, len(functions))
	for _, each := range functions {
		fws := &FunctionWithSignature{
			FunctionWithPath: each,
			Signature:        each.GetSignature(),
		}
		ret = append(ret, fws)
	}
	return ret, nil
}

// @Summary func ctx query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   file  query string true  "file"
// @Param   lines query string false "specific lines"
// @Produce json
// @Success 200 {array} sibyl2.FunctionContext
// @Router  /api/v1/funcctx [get]
func HandleFunctionCtxQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	file := c.Query("file")
	lines := c.Query("lines")

	wc := &binding.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}

	ret, err := handleFunctionQuery(repo, rev, file, lines)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var ctxs []*sibyl2.FunctionContext
	for _, each := range ret {
		funcCtx, err := sharedDriver.ReadFunctionContextWithSignature(wc, each.Signature, context.TODO())
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		ctxs = append(ctxs, funcCtx)
	}
	c.JSON(http.StatusOK, ctxs)
}