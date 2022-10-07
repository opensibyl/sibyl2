package core

func DebugDfs(unit *Unit, layer int) *Unit {
	if unit == nil {
		return nil
	}

	// dfs
	for _, each := range unit.SubUnits {
		Log.Infof("unit: %v %v %d %v", each.Kind, each.FieldName, layer, each.Content)
		DebugDfs(each, layer+1)
	}
	return nil
}
