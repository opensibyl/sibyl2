package upload

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/opensibyl/sibyl2"
)

const (
	configFile = "sibyl-upload-config"
	configType = "json"
)

type SrcConfigPart struct {
	RepoId       string   `mapstructure:"repoId"`
	RevHash      string   `mapstructure:"revHash"`
	Src          string   `mapstructure:"src"`
	Lang         []string `mapstructure:"lang"`
	WithCtx      bool     `mapstructure:"withCtx"`
	WithClass    bool     `mapstructure:"withClass"`
	IncludeRegex string   `mapstructure:"includeRegex"`
	ExcludeRegex string   `mapstructure:"excludeRegex"`
}

type ServerConfigPart struct {
	Url   string `mapstructure:"url"`
	Batch int    `mapstructure:"batch"`
	Dry   bool   `mapstructure:"dry"`
	Depth int    `mapstructure:"depth"`
}

type ContextPart struct {
	GraphCache *sibyl2.FuncGraph
}

type Config struct {
	*SrcConfigPart    `mapstructure:"src"`
	*ServerConfigPart `mapstructure:"server"`
	BizContext        *ContextPart `mapstructure:"bizContext"`
}

func (config *Config) ToMap() (map[string]any, error) {
	var m map[string]interface{}
	err := mapstructure.Decode(config, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (config *Config) ToJson() ([]byte, error) {
	toMap, err := config.ToMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(toMap)
}

func (config *Config) GetFuncUploadUrl() string {
	return fmt.Sprintf("%s/api/v1/func", config.Url)
}

func (config *Config) GetClazzUploadUrl() string {
	return fmt.Sprintf("%s/api/v1/clazz", config.Url)
}

func (config *Config) GetFuncCtxUploadUrl() string {
	return fmt.Sprintf("%s/api/v1/funcctx", config.Url)
}

func DefaultConfig() *Config {
	return &Config{
		&SrcConfigPart{
			RepoId:       "",
			RevHash:      "",
			Src:          ".",
			Lang:         []string{},
			WithCtx:      true,
			WithClass:    true,
			IncludeRegex: "",
			ExcludeRegex: "",
		},
		&ServerConfigPart{
			Url:   "http://127.0.0.1:9876",
			Batch: 50,
			Dry:   false,
			Depth: 1,
		},
		&ContextPart{nil},
	}
}
