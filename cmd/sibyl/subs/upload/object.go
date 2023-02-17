package upload

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
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

type UploadConfig struct {
	*SrcConfigPart    `mapstructure:"src"`
	*ServerConfigPart `mapstructure:"server"`
}

func (config *UploadConfig) ToMap() (map[string]any, error) {
	var m map[string]interface{}
	err := mapstructure.Decode(config, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (config *UploadConfig) ToJson() ([]byte, error) {
	toMap, err := config.ToMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(toMap)
}

func DefaultConfig() *UploadConfig {
	return &UploadConfig{
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
	}
}
