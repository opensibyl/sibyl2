package upload

import (
	"github.com/mitchellh/mapstructure"
)

const (
	configPath = "."
	configFile = "sibyl-upload-config.json"
)

type uploadConfig struct {
	Src          string `mapstructure:"src"`
	Lang         string `mapstructure:"lang"`
	Url          string `mapstructure:"url"`
	WithCtx      bool   `mapstructure:"withCtx"`
	Batch        int    `mapstructure:"batch"`
	Dry          bool   `mapstructure:"dry"`
	IncludeRegex string `mapstructure:"includeRegex"`
	ExcludeRegex string `mapstructure:"excludeRegex"`
}

func (config *uploadConfig) ToMap() (map[string]any, error) {
	var m map[string]interface{}
	err := mapstructure.Decode(config, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func defaultConfig() *uploadConfig {
	return &uploadConfig{
		Src:          ".",
		Lang:         "",
		Url:          "http://127.0.0.1:9876",
		WithCtx:      false,
		Batch:        50,
		Dry:          false,
		IncludeRegex: "",
		ExcludeRegex: "",
	}
}
