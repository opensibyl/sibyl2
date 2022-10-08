package pkg

import (
	"sibyl2/pkg/core"
	"sibyl2/pkg/extractor"
	"sibyl2/pkg/model"
)

var SibylApi = &sibylApi{}

type sibylApi struct {
}

func (*sibylApi) Extract(userSrc string, langType model.LangType, userExtractType extractor.ExtractType) ([]*model.FileResult, error) {
	runner := &core.Runner{}
	fileUnits, err := runner.File2Units(userSrc, langType)
	if err != nil {
		return nil, err
	}

	langExtractor := extractor.GetExtractor(langType)
	var results []*model.FileResult
	for _, eachFileUnit := range fileUnits {
		switch userExtractType {
		case extractor.TypeExtractSymbol:
			symbols, err := langExtractor.ExtractSymbols(eachFileUnit.Units)
			if err != nil {
				return nil, err
			}
			var retUnits []model.DataType
			for _, each := range symbols {
				retUnits = append(retUnits, model.DataType(each))
			}
			fileResult := &model.FileResult{
				Path:     eachFileUnit.Path,
				Language: eachFileUnit.Language,
				Type:     userExtractType,
				Units:    retUnits,
			}
			results = append(results, fileResult)
		case extractor.TypeExtractFunction:
			functions, err := langExtractor.ExtractFunctions(eachFileUnit.Units)
			if err != nil {
				return nil, err
			}
			var retUnits []model.DataType
			for _, each := range functions {
				retUnits = append(retUnits, model.DataType(each))
			}
			fileResult := &model.FileResult{
				Path:     eachFileUnit.Path,
				Language: eachFileUnit.Language,
				Type:     userExtractType,
				Units:    retUnits,
			}
			results = append(results, fileResult)
		}
	}
	return results, nil
}
