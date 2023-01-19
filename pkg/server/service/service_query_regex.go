package service

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/binding"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

// @Summary func query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   field   query string true  "field"
// @Param   regex   query string true  "regex"
// @Produce json
// @Success 200 {array} sibyl2.FunctionWithPath
// @Router  /api/v1/regex/func [get]
// @Tags RegexQuery
func HandleRegexFunc(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	field := c.Query("field")
	regex := c.Query("regex")

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	newRegex, err := regexp.Compile(regex)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid regex: %w", err))
		return
	}
	// regex fn
	verify := func(s string) bool {
		return newRegex.Match([]byte(s))
	}
	ruleMap := make(binding.Rule)
	ruleMap[field] = verify

	functions, err := sharedDriver.ReadFunctionsWithRule(wc, ruleMap, sharedContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, functions)
}

// @Summary clazz query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   field   query string true  "field"
// @Param   regex   query string true  "regex"
// @Produce json
// @Success 200 {array} sibyl2.ClazzWithPath
// @Router  /api/v1/regex/clazz [get]
// @Tags RegexQuery
func HandleRegexClazz(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	field := c.Query("field")
	regex := c.Query("regex")

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	newRegex, err := regexp.Compile(regex)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid regex: %w", err))
		return
	}
	// regex fn
	verify := func(s string) bool {
		return newRegex.Match([]byte(s))
	}
	ruleMap := make(binding.Rule)
	ruleMap[field] = verify

	classes, err := sharedDriver.ReadClassesWithRule(wc, ruleMap, sharedContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, classes)
}

// @Summary func ctx query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   field   query string true  "field"
// @Param   regex   query string true  "regex"
// @Produce json
// @Success 200 {array} sibyl2.FunctionContext
// @Router  /api/v1/regex/funcctx [get]
// @Tags RegexQuery
func HandleRegexFuncctx(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	field := c.Query("field")
	regex := c.Query("regex")
	ret, err := handleFuncCtxQueryWithRule(repo, rev, field, regex)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, ret)
}

func handleFuncCtxQueryWithRule(repo string, rev string, field string, regex string) ([]*sibyl2.FunctionContextSlim, error) {
	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		return nil, err
	}

	newRegex, err := regexp.Compile(regex)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %w", err)
	}
	// regex fn
	verify := func(s string) bool {
		return newRegex.Match([]byte(s))
	}
	ruleMap := make(binding.Rule)
	ruleMap[field] = verify

	functionContexts, err := sharedDriver.ReadFunctionContextsWithRule(wc, ruleMap, sharedContext)
	if err != nil {
		return nil, fmt.Errorf("failed to read func ctx: %w", err)
	}
	return functionContexts, nil
}
