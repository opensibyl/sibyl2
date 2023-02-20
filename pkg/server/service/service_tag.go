package service

import (
	"github.com/gin-gonic/gin"
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

}

// @Summary create func tag
// @Accept  json
// @Produce json
// @Success 200
// @Param   payload body tagUpload true "Payload description"
// @Router  /api/v1/tag/func [post]
// @Tags    Tag
func HandleFuncTagCreate(c *gin.Context) {

}
