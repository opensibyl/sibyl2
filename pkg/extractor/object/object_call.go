package object

import (
	"fmt"

	"github.com/opensibyl/sibyl2/pkg/core"
)

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
	Src       string    `json:"src"`
	Caller    string    `json:"caller"`
	Arguments []string  `json:"arguments"`
	Span      core.Span `json:"span"`
}

func (c *Call) GetIndexName() string {
	// hard to represent ...
	return fmt.Sprintf("%s->%s", c.Src, c.Caller)
}

func (c *Call) GetDesc() string {
	return fmt.Sprintf("<call %s(%v) in %s>", c.Caller, c.Arguments, c.Src)
}

func (c *Call) GetSpan() *core.Span {
	return &c.Span
}
