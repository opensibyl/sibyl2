package core

func DebugDfs(unit *Unit, layer int) *Unit {
	if unit == nil {
		return nil
	}

	Log.Infof("unit: %v %v %d %v", unit.Kind, unit.FieldName, layer, unit.Content)

	// dfs
	for _, each := range unit.SubUnits {
		DebugDfs(each, layer+1)
	}
	return nil
}
