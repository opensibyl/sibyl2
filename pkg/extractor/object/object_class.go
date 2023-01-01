package object

import (
	"encoding/json"
	"fmt"

	"github.com/opensibyl/sibyl2/pkg/core"
)

type ClazzSignature = string
type Clazz struct {
	Name   string `json:"name"`
	Module string `json:"module"`
	// this span will include header and content
	Span core.Span `json:"span"`

	// which contains language-specific contents
	Extras interface{} `json:"extras"`

	// ptr to origin Unit
	Unit *core.Unit `json:"-"`
}

func (c *Clazz) GetSignature() ClazzSignature {
	// It would be better unique, but it may not.
	// Not all the languages have class-file mapping like java.
	// Path + signature = Repo scope unique
	return fmt.Sprintf("%s.%s", c.Module, c.Name)
}

func NewClazz() *Clazz {
	return &Clazz{}
}

func (c *Clazz) GetIndexName() string {
	return c.GetSignature()
}

func (c *Clazz) GetDesc() string {
	return c.GetSignature()
}

func (c *Clazz) GetSpan() *core.Span {
	return &c.Span
}

func (c *Clazz) ToJson() ([]byte, error) {
	raw, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
