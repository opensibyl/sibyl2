package core

import "fmt"

func DebugDfs(unit *Unit, layer int) *Unit {
	if unit == nil {
		return nil
	}

	// dfs
	for _, each := range unit.SubUnits {
		fmt.Printf("unit: %v %v %d %v\n", each.Kind, each.FieldName, layer, each.Content)
		DebugDfs(each, layer+1)
	}
	return nil
}
