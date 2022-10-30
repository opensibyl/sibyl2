package sibyl2

import (
	"context"
	"errors"
	"fmt"
	"github.com/williamfzc/sibyl2/pkg/core"
	"github.com/williamfzc/sibyl2/pkg/extractor"
	"time"
)

// ExtractConfig todo: should not use config ptr for parallel running
type ExtractConfig struct {
	LangType    core.LangType
	ExtractType extractor.ExtractType
	FileFilter  func(path string) bool
}

func DefaultConfig() *ExtractConfig {
	return &ExtractConfig{}
}

func ExtractFromString(content string, config *ExtractConfig) (*extractor.FileResult, error) {
	return ExtractFromBytes([]byte(content), config)
}

func ExtractFromBytes(content []byte, config *ExtractConfig) (*extractor.FileResult, error) {
	startTime := time.Now()
	defer func() {
		core.Log.Infof("cost: %d ms", time.Since(startTime).Milliseconds())
	}()

	lang := config.LangType
	if !lang.IsSupported() {
		return nil, errors.New(fmt.Sprintf("unknown languages, supported: %v", core.SupportedLangs))
	}

	parser := core.NewParser(lang)
	units, err := parser.ParseCtx(content, context.TODO())
	if err != nil {
		return nil, err
	}
	langExtractor := extractor.GetExtractor(lang)
	var datas []extractor.DataType

	switch config.ExtractType {
	case extractor.TypeExtractSymbol:
		symbols, err := langExtractor.ExtractSymbols(units)
		if err != nil {
			return nil, err
		}
		datas = extractor.DataTypeOf(symbols)
	case extractor.TypeExtractFunction:
		functions, err := langExtractor.ExtractFunctions(units)
		if err != nil {
			return nil, err
		}
		datas = extractor.DataTypeOf(functions)
	case extractor.TypeExtractCall:
		calls, err := langExtractor.ExtractCalls(units)
		if err != nil {
			return nil, err
		}
		datas = extractor.DataTypeOf(calls)
	}
	result := &extractor.FileResult{
		Language: lang,
		Units:    datas,
		Type:     config.ExtractType,
	}
	return result, nil
}
