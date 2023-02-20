package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

// @Summary func query
// @Param   repo  query string true "repo"
// @Param   rev   query string true "rev"
// @Param   regex query string true "regex"
// @Produce json
// @Success 200 {array} string
// @Router  /api/v1/signature/regex/func [get]
// @Tags    SignatureQuery
func HandleSignatureRegexFunc(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	regex := c.Query("regex")

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	signatures, err := sharedDriver.ReadFunctionSignaturesWithRegex(wc, regex, sharedContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, signatures)
}

// @Summary func query
// @Param   repo      query string true "repo"
// @Param   rev       query string true "rev"
// @Param   signature query string true "signature"
// @Produce json
// @Success 200 {object} sibyl2.FunctionWithTag
// @Router  /api/v1/signature/func [get]
// @Tags    SignatureQuery
func HandleSignatureFunc(c *gin.Context) {
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
	f, err := sharedDriver.ReadFunctionWithSignature(wc, signature, sharedContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, f)
}

// @Summary funcctx query
// @Param   repo      query string true "repo"
// @Param   rev       query string true "rev"
// @Param   signature query string true "signature"
// @Produce json
// @Success 200 {object} sibyl2.FunctionContextSlim
// @Router  /api/v1/signature/funcctx [get]
// @Tags    SignatureQuery
func HandleSignatureFuncctx(c *gin.Context) {
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
	c.JSON(http.StatusOK, f)
}
