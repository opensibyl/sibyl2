package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/server/binding"
)

var sharedDriver binding.Driver

type FunctionWithSignature struct {
	*sibyl2.FunctionWithPath
	Signature string `json:"signature"`
}

type FunctionUploadUnit struct {
	WorkspaceConfig *binding.WorkspaceConfig      `json:"workspace"`
	FunctionResult  *extractor.FunctionFileResult `json:"funcResult"`
}

func HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HandleRepoQuery(c *gin.Context) {
	repos, err := sharedDriver.ReadRepos(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, repos)
}

func HandleRevQuery(c *gin.Context) {
	repo := c.Query("repo")
	revs, err := sharedDriver.ReadRevs(repo, context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, revs)
}

func HandleFileQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	wc := &binding.WorkspaceConfig{
		RepoId:  repo,
		RevHash: rev,
	}
	if err := wc.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	files, err := sharedDriver.ReadFiles(wc, context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, files)
}

func HandleFunctionsQuery(c *gin.Context) {
	repo := c.Query("repo")
	rev := c.Query("rev")
	file := c.Query("file")
	lines := c.Query("lines")

	wc := &binding.WorkspaceConfig{
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
		functions, err = sharedDriver.ReadFunctions(wc, file, context.TODO())
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
		functions, err = sharedDriver.ReadFunctionsWithLines(wc, file, lineNums, context.TODO())
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

func HandleRepoFuncUpload(c *gin.Context) {
	result := &FunctionUploadUnit{}
	err := c.BindJSON(result)
	if err != nil {
		core.Log.Errorf("error when parse: %v\n", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("parse json error: %v", err))
		return
	}
	if err := result.WorkspaceConfig.Verify(); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	go func() {
		ctx := context.Background()

		err = sharedDriver.CreateWorkspace(result.WorkspaceConfig, ctx)
		if err != nil {
			core.Log.Warnf("error when init: %v\n", err)
		}

		err := sharedDriver.CreateFuncFile(result.WorkspaceConfig, result.FunctionResult, ctx)
		if err != nil {
			core.Log.Errorf("error when upload: %v\n", err)
		}
	}()

	c.JSON(http.StatusOK, "received")
}

func Execute() {
	driver, err := binding.NewInMemoryDriver()
	if err != nil {
		panic(err)
	}
	sharedDriver = driver

	engine := gin.Default()
	engine.Handle(http.MethodGet, "/ping", HandlePing)

	v1group := engine.Group("/api/v1")
	v1group.Handle(http.MethodGet, "/repo", HandleRepoQuery)
	v1group.Handle(http.MethodGet, "/rev", HandleRevQuery)
	v1group.Handle(http.MethodGet, "/file", HandleFileQuery)
	v1group.Handle(http.MethodGet, "/func", HandleFunctionsQuery)

	v1group.Handle(http.MethodPost, "/func", HandleRepoFuncUpload)

	err = engine.Run(fmt.Sprintf(":%d", 9876))
	if err != nil {
		fmt.Printf("failed to start repoctor_receiver: %s", err.Error())
	}
}
