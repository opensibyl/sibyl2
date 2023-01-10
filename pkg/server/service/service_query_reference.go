package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/server/binding"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

// @Summary funcctx query by ref
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   moreThan   query int true  "moreThan"
// @Param   lessThan   query int true  "lessThan"
// @Produce json
// @Success 200 {array} sibyl2.FunctionWithPath
// @Router  /api/v1/funcctx/with/reference/count [get]
// @Tags EXPERIMENTAL
func HandleFuncQueryWithReferenceCount(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	moreThan := c.Query("moreThan")
	lessThan := c.Query("lessThan")

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	moreThanInt, err := strconv.Atoi(moreThan)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	lessThanInt, err := strconv.Atoi(lessThan)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	field := "calls.#"
	verify := func(s string) bool {
		count, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		if count > moreThanInt && count < lessThanInt {
			return true
		}
		return false
	}
	ruleMap := make(binding.Rule)
	ruleMap[field] = verify

	functionContexts, err := sharedDriver.ReadFunctionContextsWithRule(wc, ruleMap, sharedContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, functionContexts)
}

// @Summary funcctx query by referenced
// @Param   repo  query string true  "repo"
// @Param   rev   query string true  "rev"
// @Param   moreThan   query int true  "moreThan"
// @Param   lessThan   query int true  "lessThan"
// @Produce json
// @Success 200 {array} sibyl2.FunctionWithPath
// @Router  /api/v1/funcctx/with/referenced/count [get]
// @Tags EXPERIMENTAL
func HandleFuncQueryWithReferencedCount(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	moreThan := c.Query("moreThan")
	lessThan := c.Query("lessThan")

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	moreThanInt, err := strconv.Atoi(moreThan)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	lessThanInt, err := strconv.Atoi(lessThan)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	field := "calls.#"
	verify := func(s string) bool {
		count, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		if count > moreThanInt && count < lessThanInt {
			return true
		}
		return false
	}
	ruleMap := make(binding.Rule)
	ruleMap[field] = verify

	functionContexts, err := sharedDriver.ReadFunctionContextsWithRule(wc, ruleMap, sharedContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, functionContexts)
}
