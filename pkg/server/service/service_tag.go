package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

type tagUpload struct {
	RepoId    string `json:"repoId"`
	RevHash   string `json:"revHash"`
	Signature string `json:"signature"`
	Tag       string `json:"tag"`
}

// @Summary query func by tag
// @Produce json
// @Param   repo query   string true "repo"
// @Param   rev  query   string true "rev"
// @Param   tag  query   string true "tag"
// @Success 200  {array} string
// @Router  /api/v1/tag/func [get]
// @Tags    Tag
func HandleFuncTagQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	tag := c.Query("tag")

	wc := &object.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fs, err := sharedDriver.ReadFunctionsWithTag(wc, tag, sharedContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, fs)
}

// @Summary create func tag
// @Accept  json
// @Produce json
// @Success 200
// @Param   payload body tagUpload true "tag upload payload"
// @Router  /api/v1/tag/func [post]
// @Tags    Tag
func HandleFuncTagCreate(c *gin.Context) {
	result := &tagUpload{}
	err := c.BindJSON(result)
	if err != nil {
		core.Log.Errorf("error when parse: %v\n", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("parse json error: %v", err))
		return
	}

	wc := &object.WorkspaceConfig{
		RepoId:  result.RepoId,
		RevHash: result.RevHash,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = sharedDriver.CreateFuncTag(wc, result.Signature, result.Tag, sharedContext)
	if err != nil {
		core.Log.Errorf("create func tag failed: %v", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
