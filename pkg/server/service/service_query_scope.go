package service

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

// @Summary repo query
// @Produce json
// @Success 200 {array} string
// @Router  /api/v1/repo [get]
// @Tags SCOPE
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
// @Tags SCOPE
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
// @Param   includeRegex  query string false "includeRegex"
// @Produce json
// @Success 200 {array} string
// @Router  /api/v1/file [get]
// @Tags SCOPE
func HandleFileQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	includeRegex := c.Query("includeRegex")

	var compiledRegex *regexp.Regexp
	var err error
	if includeRegex != "" {
		compiledRegex, err = regexp.Compile(includeRegex)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
	}

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
	if compiledRegex != nil {
		filesAfterFilter := make([]string, 0)
		for _, each := range files {
			if compiledRegex.MatchString(each) {
				filesAfterFilter = append(filesAfterFilter, each)
			}
		}
		c.JSON(http.StatusOK, filesAfterFilter)
		return
	}

	c.JSON(http.StatusOK, files)
}
