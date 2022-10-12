package model

/*
Call NON-PRECISE

	func aFunc() {
		a.b(c, d)
	}

	aFunc  == Src
	a.b    == caller
	[c, d] == arguments
*/
type Call struct {
	Src       string   `json:"src"`
	Caller    string   `json:"caller"`
	Arguments []string `json:"arguments"`
	Span      Span     `json:"span"`
}
