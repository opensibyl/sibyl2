package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
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
