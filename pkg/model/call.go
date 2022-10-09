package model

// Call NON-PRECISE
type Call struct {
	// function which starts this call
	Src *Function `json:"src"`

	// a.b(c, d)
	// a == caller
	// b == function name
	// [c, d] == parameters
	Caller     string   `json:"caller"`
	FuncName   string   `json:"funcName"`
	Parameters []string `json:"parameters"`
}
