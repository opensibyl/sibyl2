package core

import (
	"context"
	sitter "github.com/smacker/go-tree-sitter"
)

/*
Parser

- get almost all the nodes
- convert them to units
*/
type Parser struct {
	engine *sitter.Parser
}

func NewParser(lang LangType) *Parser {
	engine := sitter.NewParser()
	engine.SetLanguage(lang.GetLanguage())
	return &Parser{
		engine,
	}
}

func (p *Parser) ParseCtx(data []byte, context context.Context) ([]*Unit, error) {
	tree, err := p.engine.ParseCtx(context, nil, data)
	if err != nil {
		return nil, err
	}
	return p.node2Units(data, tree.RootNode(), "", nil)
}

func (p *Parser) ParseStringCtx(data string, context context.Context) ([]*Unit, error) {
	return p.ParseCtx([]byte(data), context)
}

func (p *Parser) Parse(data []byte) ([]*Unit, error) {
	return p.ParseCtx(data, context.TODO())
}

func (p *Parser) ParseString(data string) ([]*Unit, error) {
	return p.ParseCtx([]byte(data), context.TODO())
}

// DFS
func (p *Parser) node2Units(data []byte, curRootNode *sitter.Node, fieldName string, parentUnit *Unit) ([]*Unit, error) {
	var ret []*Unit

	// itself
	curRootUnit, err := p.node2Unit(data, curRootNode, fieldName, parentUnit)
	if err != nil {
		return nil, err
	}
	ret = append(ret, curRootUnit)

	count := int(curRootNode.NamedChildCount())
	for i := 0; i < count; i++ {
		curChild := curRootNode.NamedChild(i)
		curChildName := curRootNode.FieldNameForChild(i)

		subUnits, err := p.node2Units(data, curChild, curChildName, curRootUnit)
		if err != nil {
			return nil, err
		}
		curRootUnit.SubUnits = append(curRootUnit.SubUnits, subUnits[0])

		ret = append(ret, subUnits...)
	}
	return ret, nil
}

func (p *Parser) node2Unit(data []byte, node *sitter.Node, fieldName string, parentUnit *Unit) (*Unit, error) {
	ret := &Unit{}

	ret.FieldName = fieldName
	// SHOULD WE ALWAYS KEEP THIS DATA?
	ret.Content = node.Content(data)

	// kind: type of type
	// https://cs.stackexchange.com/questions/111430/whats-the-difference-between-a-type-and-a-kind
	// what it is in this language
	ret.Kind = node.Type()

	// range
	ret.Span = Span{
		Start: Point{Row: node.StartPoint().Row, Column: node.StartPoint().Column},
		End:   Point{Row: node.EndPoint().Row, Column: node.EndPoint().Column},
	}
	ret.ParentUnit = parentUnit
	return ret, nil
}
