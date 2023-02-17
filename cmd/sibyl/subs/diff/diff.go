package diff

import (
	"os"
	"path/filepath"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/ext"
	"github.com/opensibyl/sibyl2/pkg/extractor"
)

type FunctionWithWeight struct {
	*extractor.Function
	ReferenceCount  int `json:"referenceCount"`
	ReferencedCount int `json:"referencedCount"`
}

type FileFragment struct {
	Path      string                    `json:"path"`
	Lines     []int                     `json:"lines"`
	Functions []*sibyl2.FunctionContext `json:"functions"`
}

type ThinFileFragment struct {
	Path      string                `json:"path"`
	Lines     []int                 `json:"lines"`
	Functions []*FunctionWithWeight `json:"functions"`
}

func (f *FileFragment) Flatten() *ThinFileFragment {
	ret := &ThinFileFragment{}
	ret.Path = f.Path

	lineCount := len(f.Lines)
	switch lineCount {
	case 0:
		ret.Lines = []int{}
	case 1:
		ret.Lines = []int{f.Lines[0]}
	default:
		ret.Lines = []int{f.Lines[0], f.Lines[lineCount-1]}
	}

	for _, each := range f.Functions {
		fww := &FunctionWithWeight{
			Function:        each.Function,
			ReferenceCount:  len(each.Calls),
			ReferencedCount: len(each.ReverseCalls),
		}
		ret.Functions = append(ret.Functions, fww)
	}
	return ret
}

type ParseResult struct {
	Fragments []*FileFragment `json:"fragments"`
}

type ThinParseResult struct {
	Fragments []*ThinFileFragment `json:"fragments"`
}

func (p *ParseResult) Flatten() *ThinParseResult {
	ret := &ThinParseResult{}
	for _, each := range p.Fragments {
		ret.Fragments = append(ret.Fragments, each.Flatten())
	}
	return ret
}

func parsePatchRaw(srcDir string, patchRaw []byte) (*ParseResult, error) {
	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return nil, err
	}

	affectedMap, err := ext.Unified2Affected(patchRaw)
	if err != nil {
		return nil, err
	}

	// data ready
	result, err := affectedLines2Functions(srcDir, &affectedMap)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func parsePatch(srcDir string, patchFile string) (*ParseResult, error) {
	patch, err := os.ReadFile(patchFile)
	if err != nil {
		return nil, err
	}
	return parsePatchRaw(srcDir, patch)
}

func affectedLines2Functions(srcDir string, m *ext.AffectedLineMap) (*ParseResult, error) {
	f, err := sibyl2.ExtractFunction(srcDir, sibyl2.DefaultConfig())
	if err != nil {
		return nil, err
	}

	s, err := sibyl2.ExtractSymbol(srcDir, sibyl2.DefaultConfig())
	if err != nil {
		return nil, err
	}
	g, err := sibyl2.AnalyzeFuncGraph(f, s)
	if err != nil {
		return nil, err
	}

	result := &ParseResult{}
	for fileName, lines := range *m {
		fragment := &FileFragment{}
		fragment.Path = fileName
		fragment.Lines = lines
		for _, eachFuncFile := range f {
			// make it map better
			if eachFuncFile.Path == fileName {
				for _, eachFuncUnit := range eachFuncFile.Units {
					if eachFuncUnit.GetSpan().ContainAnyLine(lines...) {
						related := g.FindRelated(sibyl2.WrapFuncWithPath(eachFuncUnit, fileName))
						fragment.Functions = append(fragment.Functions, related)
					}
				}
				break
			}
		}
		result.Fragments = append(result.Fragments, fragment)
	}

	return result, nil
}
