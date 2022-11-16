package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/williamfzc/sibyl2"
)

var GlobalStorage *sibyl2.InMemoryStorage

type FunctionWithSignature struct {
	*sibyl2.FunctionWithPath
	Signature string `json:"signature"`
}

func HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HandleRepoQuery(c *gin.Context) {
	newDriver, _ := sibyl2.NewInMemoryDriver(GlobalStorage)
	repos, err := newDriver.ReadRepos(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, repos)
}

func HandleRevQuery(c *gin.Context) {
	newDriver, _ := sibyl2.NewInMemoryDriver(GlobalStorage)
	repo := c.Query("repo")
	revs, err := newDriver.ReadRevs(repo, context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, revs)
}

func HandleFileQuery(c *gin.Context) {
	newDriver, _ := sibyl2.NewInMemoryDriver(GlobalStorage)
	repo := c.Query("repo")
	rev := c.Query("rev")
	wc := &sibyl2.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	files, err := newDriver.ReadFiles(wc, context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, files)
}

func HandleFunctionsQuery(c *gin.Context) {
	newDriver, _ := sibyl2.NewInMemoryDriver(GlobalStorage)
	repo := c.Query("repo")
	rev := c.Query("rev")
	file := c.Query("file")
	lines := c.Query("lines")

	wc := &sibyl2.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var functions []*sibyl2.FunctionWithPath
	var err error
	if lines == "" {
		functions, err = newDriver.ReadFunctions(wc, file, context.TODO())
	} else {
		linesStrList := strings.Split(lines, ",")
		var lineNums = make([]int, 0, len(linesStrList))
		for _, each := range linesStrList {
			num, err := strconv.Atoi(each)
			if err != nil {
				c.JSON(http.StatusBadGateway, err)
				return
			}
			lineNums = append(lineNums, num)
		}
		functions, err = newDriver.ReadFunctionsWithLines(wc, file, lineNums, context.TODO())
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// export signature
	ret := make([]*FunctionWithSignature, 0, len(functions))
	for _, each := range functions {
		fws := &FunctionWithSignature{
			FunctionWithPath: each,
			Signature:        each.GetSignature(),
		}
		ret = append(ret, fws)
	}

	c.JSON(http.StatusOK, ret)
}

func Execute() {
	engine := gin.Default()
	engine.Handle(http.MethodGet, "/ping", HandlePing)

	v1group := engine.Group("/api/v1")
	v1group.Handle(http.MethodGet, "/repo", HandleRepoQuery)
	v1group.Handle(http.MethodGet, "/rev", HandleRevQuery)
	v1group.Handle(http.MethodGet, "/file", HandleFileQuery)
	v1group.Handle(http.MethodGet, "/func", HandleFunctionsQuery)

	err := engine.Run(fmt.Sprintf(":%d", 9876))
	if err != nil {
		fmt.Printf("failed to start repoctor_receiver: %s", err.Error())
	}
}

func init() {
	GlobalStorage = &sibyl2.InMemoryStorage{}
}
