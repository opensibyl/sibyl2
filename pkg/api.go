package pkg

import (
	"errors"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"github.com/williamfzc/sibyl2/pkg/model"
	"os"
	"time"
)

var SibylApi = &sibylApi{}

type sibylApi struct {
}

func (*sibylApi) Extract(userSrc string, langType model.LangType, userExtractType extractor.ExtractType) ([]*model.FileResult, error) {
	startTime := time.Now()
	defer func() {
		core.Log.Infof("cost: %d ms", time.Since(startTime).Milliseconds())
	}()

	if _, err := os.Stat(userSrc); os.IsNotExist(err) {
		return nil, errors.New("file not existed: " + userSrc)
	}

	runner := &core.Runner{}
	fileUnits, err := runner.File2Units(userSrc, langType)
	if err != nil {
		return nil, err
	}

	langExtractor := extractor.GetExtractor(langType)
	var results []*model.FileResult
	for _, eachFileUnit := range fileUnits {
		fileResult := &model.FileResult{
			Path:     eachFileUnit.Path,
			Language: eachFileUnit.Language,
			Type:     userExtractType,
		}

		switch userExtractType {
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
	err = model.PathStandardize(results, userSrc)
	if err != nil {
		return nil, err
	}

	return results, nil
}
