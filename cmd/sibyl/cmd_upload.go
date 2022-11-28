package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/server"
	"github.com/williamfzc/sibyl2/pkg/server/binding"
)

var uploadSrc string
var uploadLangType string
var uploadUrl string

func NewUploadCmd() *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:    "upload",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			uploadSrc, err := filepath.Abs(uploadSrc)
			if err != nil {
				panic(err)
			}
			repo, err := loadRepo(uploadSrc)
			if err != nil {
				panic(err)
			}
			head, err := repo.Head()
			curRepo := filepath.Base(uploadSrc)
			curRev := head.Hash().String()

			wc := &binding.WorkspaceConfig{
				RepoId:  curRepo,
				RevHash: curRev,
			}
			f, err := sibyl2.ExtractFunction(uploadSrc, sibyl2.DefaultConfig())
			if err != nil {
				panic(err)
			}

			fullUrl := fmt.Sprintf("%s/api/v1/func", uploadUrl)
			core.Log.Infof("upload backend: %s", fullUrl)
			uploadFunctions(fullUrl, wc, f)
		},
	}
	uploadCmd.PersistentFlags().StringVar(&uploadSrc, "src", ".", "src dir path")
	uploadCmd.PersistentFlags().StringVar(&uploadLangType, "lang", "", "lang type of your source code")
	uploadCmd.PersistentFlags().StringVar(&uploadUrl, "url", "http://127.0.0.1:9876", "backend url")

	return uploadCmd
}

func loadRepo(gitDir string) (*git.Repository, error) {
	repo, err := git.PlainOpen(gitDir)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func uploadFunctions(url string, wc *binding.WorkspaceConfig, f []*extractor.FunctionFileResult) {
	core.Log.Infof("uploading %v with files %d ...", wc, len(f))
	var wg sync.WaitGroup
	wg.Add(len(f))
	for _, each := range f {
		unit := &server.FunctionUploadUnit{
			WorkspaceConfig: wc,
			FunctionResult:  each,
		}

		go func() {
			defer wg.Done()
			jsonStr, err := json.Marshal(unit)
			if err != nil {
				panic(err)
			}
			resp, err := http.Post(
				url,
				"application/json",
				bytes.NewBuffer(jsonStr))
			if err != nil {
				panic(err)
			}
			data, err := io.ReadAll(resp.Body)
			if resp.StatusCode != http.StatusOK {
				core.Log.Errorf("upload %s resp: %v", unit.FunctionResult.Path, string(data))
			}
		}()
	}
	wg.Wait()
}

func init() {
	uploadCmd := NewUploadCmd()
	rootCmd.AddCommand(uploadCmd)
}
