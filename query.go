package sibyl2

import "github.com/williamfzc/sibyl2/pkg/extractor"

func QueryAffectedUnitsByLine[T extractor.DataType](result *extractor.BaseFileResult[T], lines ...int) []T {
	var ret []T
	for _, eachUnit := range result.Units {
		eachSpan := eachUnit.GetSpan()
		if eachSpan.ContainAnyLine(lines...) {
			ret = append(ret, eachUnit)
		}
	}
	return ret
}
