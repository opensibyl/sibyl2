package sibyl2

import (
	"errors"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/model"
	"os"
	"time"
)

type ExtractConfig struct {
	LangType    model.LangType
	ExtractType extractor.ExtractType
}

func Extract(targetDir string, config *ExtractConfig) ([]*model.FileResult, error) {
	startTime := time.Now()
	defer func() {
		core.Log.Infof("cost: %d ms", time.Since(startTime).Milliseconds())
	}()

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return nil, errors.New("file not existed: " + targetDir)
	}

	runner := &core.Runner{}
	fileUnits, err := runner.File2Units(targetDir, config.LangType)
	if err != nil {
		return nil, err
	}

	langExtractor := extractor.GetExtractor(config.LangType)
	var results []*model.FileResult
	for _, eachFileUnit := range fileUnits {
		fileResult := &model.FileResult{
			Path:     eachFileUnit.Path,
			Language: eachFileUnit.Language,
			Type:     config.ExtractType,
		}

		switch config.ExtractType {
		case extractor.TypeExtractSymbol:
			symbols, err := langExtractor.ExtractSymbols(eachFileUnit.Units)
			if err != nil {
				return nil, err
			}
			fileResult.Units = model.DataTypeOf(symbols)
		case extractor.TypeExtractFunction:
			functions, err := langExtractor.ExtractFunctions(eachFileUnit.Units)
			if err != nil {
				return nil, err
			}
			fileResult.Units = model.DataTypeOf(functions)
		case extractor.TypeExtractCall:
			calls, err := langExtractor.ExtractCalls(eachFileUnit.Units)
			if err != nil {
				return nil, err
			}
			fileResult.Units = model.DataTypeOf(calls)
		}
		results = append(results, fileResult)
	}
	// path
	err = model.PathStandardize(results, targetDir)
	if err != nil {
		return nil, err
	}

	return results, nil
}
