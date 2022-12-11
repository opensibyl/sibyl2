package ext

import (
	"bytes"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
)

type AffectedLineMap = map[string][]int

func Unified2Affected(patch []byte) (AffectedLineMap, error) {
	parsed, _, err := gitdiff.Parse(bytes.NewReader(patch))
	if err != nil {
		return nil, err
	}

	affectedMap := make(map[string][]int)
	for _, each := range parsed {
		if each.IsBinary || each.IsDelete {
			continue
		}
		affectedMap[each.NewName] = make([]int, 0)
		fragments := each.TextFragments
		for _, eachF := range fragments {
			left := int(eachF.NewPosition)

			for i, eachLine := range eachF.Lines {
				if eachLine.New() && eachLine.Op == gitdiff.OpAdd {
					affectedMap[each.NewName] = append(affectedMap[each.NewName], left+i-1)
				}
			}
		}
	}
	return affectedMap, nil
}
