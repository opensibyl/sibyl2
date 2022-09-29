package core

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

// Parser get almost all the nodes (units)
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

func (p *Parser) ParseCtx(data []byte, context context.Context) ([]Unit, error) {
	tree, err := p.engine.ParseCtx(context, nil, data)
	if err != nil {
		return nil, err
	}
	return p.node2Units(data, tree.RootNode())
}

func (p *Parser) Parse(data []byte) ([]Unit, error) {
	return p.ParseCtx(data, context.TODO())
}

func (p *Parser) node2Units(data []byte, rootNode *sitter.Node) ([]Unit, error) {
	var ret []Unit
	count := int(rootNode.NamedChildCount())
	for i := 0; i < count; i++ {
		curChild := rootNode.NamedChild(i)
		curChildName := rootNode.FieldNameForChild(i)
		curSymbol, err := p.node2Unit(data, curChild, curChildName)

		if err != nil {
			return nil, err
		}

		ret = append(ret, curSymbol)
		// handle its sons
		subSymbols, err := p.node2Units(data, curChild)
		if err != nil {
			return nil, err
		}
		ret = append(ret, subSymbols...)
	}
	return ret, nil
}

func (p *Parser) node2Unit(data []byte, node *sitter.Node, name string) (Unit, error) {
	ret := Unit{}

	ret.Content = node.Content(data)
	ret.FieldName = name

	// kind: type of type
	// https://cs.stackexchange.com/questions/111430/whats-the-difference-between-a-type-and-a-kind
	// what it is in this language
	ret.Kind = node.Type()

	// range
	ret.Span = Span{
		Start: Point{node.StartPoint().Row, node.StartPoint().Column},
		End:   Point{node.EndPoint().Row, node.EndPoint().Column},
	}

	return ret, nil
}
