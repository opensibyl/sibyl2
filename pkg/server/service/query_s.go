package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/binding"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/opensibyl/sibyl2/pkg/server/queue"
)

var sharedContext context.Context
var sharedDriver binding.Driver
var sharedQueue queue.Queue

// @Summary repo query
// @Produce json
// @Success 200 {array} string
// @Router  /api/v1/repo [get]
func HandleRepoQuery(c *gin.Context) {
	repos, err := sharedDriver.ReadRepos(sharedContext)
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
	revs, err := sharedDriver.ReadRevs(repo, sharedContext)
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
	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	files, err := sharedDriver.ReadFiles(wc, sharedContext)
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
// @Success 200 {array} object.FunctionWithSignature
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

func handleFunctionQuery(repo string, rev string, file string, lines string) ([]*object.FunctionWithSignature, error) {
	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		return nil, err
	}
	var functions []*sibyl2.FunctionWithPath
	var err error
	if lines == "" {
		functions, err = sharedDriver.ReadFunctions(wc, file, sharedContext)
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
		functions, err = sharedDriver.ReadFunctionsWithLines(wc, file, lineNums, sharedContext)
	}
	if err != nil {
		return nil, err
	}

	// export signature
	ret := make([]*object.FunctionWithSignature, 0, len(functions))
	for _, each := range functions {
		fws := &object.FunctionWithSignature{
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

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}

	ret, err := handleFunctionQuery(repo, rev, file, lines)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ctxs := make([]*sibyl2.FunctionContext, 0)
	for _, each := range ret {
		funcCtx, err := sharedDriver.ReadFunctionContextWithSignature(wc, each.Signature, sharedContext)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			core.Log.Errorf("err when read with sign: %v", err)
			return
		}
		ctxs = append(ctxs, funcCtx)
	}
	c.JSON(http.StatusOK, ctxs)
}

func InitService(_ object.ExecuteConfig, ctx context.Context, driver binding.Driver, q queue.Queue) {
	sharedContext = ctx
	sharedDriver = driver
	sharedQueue = q
}
