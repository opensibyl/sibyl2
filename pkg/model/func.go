package model

type ValueUnit struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Function struct {
	Name       string       `json:"name"`
	Receiver   string       `json:"receiver"`
	Parameters []*ValueUnit `json:"parameters"`
	Returns    []*ValueUnit `json:"returns"`
	Span       Span         `json:"span"`
}
