package sibyl2

import "github.com/williamfzc/sibyl2/pkg/extractor"

func QueryUnitsByLines[T extractor.DataType](result *extractor.BaseFileResult[T], lines ...int) []T {
	var ret []T
	for _, eachUnit := range result.Units {
		eachSpan := eachUnit.GetSpan()
		if eachSpan.ContainAnyLine(lines...) {
			ret = append(ret, eachUnit)
		}
	}
	return ret
}

func QueryUnitsByIndexNames[T extractor.DataType](result *extractor.BaseFileResult[T], indexNames ...string) []T {
	var ret []T
	for _, eachUnit := range result.Units {
		for _, eachName := range indexNames {
			if eachUnit.GetIndexName() == eachName {
				ret = append(ret, eachUnit)
			}
		}
	}
	return ret
}

func QueryUnitsByIndexNamesInFiles[T extractor.DataType](result []*extractor.BaseFileResult[T], indexNames ...string) []T {
	var ret []T
	for _, each := range result {
		ret = append(ret, QueryUnitsByIndexNames(each, indexNames...)...)
	}
	return ret
}
