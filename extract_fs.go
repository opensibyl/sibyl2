package sibyl2

import (
	"errors"
	"fmt"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"os"
	"path/filepath"
	"time"
)

func ExtractSymbol(targetFile string, config *ExtractConfig) ([]*extractor.SymbolFileResult, error) {
	config.ExtractType = extractor.TypeExtractSymbol
	results, err := Extract(targetFile, config)
	if err != nil {
		return nil, err
	}

	var final []*extractor.SymbolFileResult
	for _, each := range results {
		var newUnits = make([]*extractor.Symbol, len(each.Units))
		for i, v := range each.Units {
			// should not error
			if s, ok := v.(*extractor.Symbol); ok {
				newUnits[i] = s
			} else {
				return nil, errors.New(fmt.Sprintf("failed to cast %v to symbol", v))
			}
		}

		newEach := &extractor.SymbolFileResult{
			Path:     each.Path,
			Language: each.Language,
			Type:     each.Type,
			Units:    newUnits,
		}
		final = append(final, newEach)
	}
	return final, nil
}

func ExtractFunction(targetFile string, config *ExtractConfig) ([]*extractor.FunctionFileResult, error) {
	config.ExtractType = extractor.TypeExtractFunction
	results, err := Extract(targetFile, config)
	if err != nil {
		return nil, err
	}

	var final []*extractor.FunctionFileResult
	for _, each := range results {
		var newUnits = make([]*extractor.Function, len(each.Units))
		for i, v := range each.Units {
			// should not error
			if f, ok := v.(*extractor.Function); ok {
				newUnits[i] = f
			} else {
				return nil, errors.New(fmt.Sprintf("failed to cast %v to function", v))
			}
		}

		newEach := &extractor.FunctionFileResult{
			Path:     each.Path,
			Language: each.Language,
			Type:     each.Type,
			Units:    newUnits,
		}
		final = append(final, newEach)
	}
	return final, nil
}

func Extract(targetFile string, config *ExtractConfig) ([]*extractor.FileResult, error) {
	startTime := time.Now()
	defer func() {
		core.Log.Infof("cost: %d ms", time.Since(startTime).Milliseconds())
	}()

	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		return nil, errors.New("file not existed: " + targetFile)
	}

	// always use abs path and convert it back at the end
	targetFile, err := filepath.Abs(targetFile)
	if err != nil {
		return nil, err
	}

	runner := &core.Runner{}
	if !config.LangType.IsSupported() {
		// do the guess
		core.Log.Infof("no specific lang found, do the guess in: %s", targetFile)
		config.LangType, err = runner.GuessLangFromDir(targetFile, config.FileFilter)
		if err != nil {
			return nil, err
		}
		core.Log.Infof("I think it is: %s", config.LangType)
	}
	// still failed, give up
	if !config.LangType.IsSupported() {
		return nil, errors.New(fmt.Sprintf("unknown languages, supported: %v", core.SupportedLangs))
	}

	fileUnits, err := runner.File2Units(targetFile, config.LangType, config.FileFilter)
	if err != nil {
		return nil, err
	}

	langExtractor := extractor.GetExtractor(config.LangType)
	var results []*extractor.FileResult
	for _, eachFileUnit := range fileUnits {
		fileResult := &extractor.FileResult{
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
			fileResult.Units = extractor.DataTypeOf(symbols)
		case extractor.TypeExtractFunction:
			functions, err := langExtractor.ExtractFunctions(eachFileUnit.Units)
			if err != nil {
				return nil, err
			}
			fileResult.Units = extractor.DataTypeOf(functions)
		case extractor.TypeExtractCall:
			calls, err := langExtractor.ExtractCalls(eachFileUnit.Units)
			if err != nil {
				return nil, err
			}
			fileResult.Units = extractor.DataTypeOf(calls)
		default:
			return nil, errors.New("no specific extract type")
		}
		results = append(results, fileResult)
	}
	// path
	err = extractor.PathStandardize(results, targetFile)
	if err != nil {
		return nil, err
	}

	return results, nil
}
