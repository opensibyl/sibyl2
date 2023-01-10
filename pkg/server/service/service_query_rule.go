package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

// @Summary func query
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   rule   query string true  "rule"
// @Produce json
// @Success 200 {array} sibyl2.FunctionWithPath
// @Router  /api/v1/func/with/rule [get]
// @Tags EXPERIMENTAL
func HandleFunctionQueryWithRule(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	rule := c.Query("rule")

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ruleMap := make(map[string]string)
	err := json.Unmarshal([]byte(rule), &ruleMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

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
// @Param   rule   query string true  "rule"
// @Produce json
// @Success 200 {array} sibyl2.ClazzWithPath
// @Router  /api/v1/clazz/with/rule [get]
// @Tags EXPERIMENTAL
func HandleClazzQueryWithRule(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	rule := c.Query("rule")

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	ruleMap := make(map[string]string)
	err := json.Unmarshal([]byte(rule), &ruleMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

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
// @Param   rule   query string true  "rule"
// @Produce json
// @Success 200 {array} sibyl2.FunctionContext
// @Router  /api/v1/funcctx/with/rule [get]
// @Tags EXPERIMENTAL
func HandleFuncCtxQueryWithRule(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	rule := c.Query("rule")
	ret, err := handleFuncCtxQueryWithRule(repo, rev, rule)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, ret)
}

func handleFuncCtxQueryWithRule(repo string, rev string, rule string) ([]*sibyl2.FunctionContext, error) {
	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		return nil, err
	}

	ruleMap := make(map[string]string)
	err := json.Unmarshal([]byte(rule), &ruleMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	functionContexts, err := sharedDriver.ReadFunctionContextsWithRule(wc, ruleMap, sharedContext)
	if err != nil {
		return nil, fmt.Errorf("failed to read func ctx: %w", err)
	}
	return functionContexts, nil
}
