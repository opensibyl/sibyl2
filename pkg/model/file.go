package model

type FileUnit struct {
	Path     string   `json:"path"`
	Language LangType `json:"language"`
	Units    []*Unit  `json:"units"`
}

type FileResult struct {
	Path     string     `json:"path"`
	Language LangType   `json:"language"`
	Type     string     `json:"type"`
	Units    []DataType `json:"units"`
}
