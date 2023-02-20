package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

type RevStat struct {
	Info          *object.RevInfo `json:"info"`
	FileCount     int             `json:"fileCount"`
	FunctionCount int             `json:"functionCount"`
}

// @Summary rev stat
// @Param   repo query string true "repo"
// @Param   rev  query string true "rev"
// @Produce json
// @Success 200 {object} RevStat
// @Router  /api/v1/rev/stat [get]
// @Tags    StatQuery
func HandleRevStatQuery(c *gin.Context) {
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

	stat := &RevStat{}

	// info
	revInfo, err := sharedDriver.ReadRevInfo(wc, sharedContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	stat.Info = revInfo

	// file level
	files, err := sharedDriver.ReadFiles(wc, sharedContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	stat.FileCount = len(files)

	// func level
	functions, err := sharedDriver.ReadFunctionSignaturesWithRegex(wc, ".*", sharedContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	stat.FunctionCount = len(functions)

	c.JSON(http.StatusOK, stat)
}
