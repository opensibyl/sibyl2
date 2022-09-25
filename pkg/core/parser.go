package core

import (
	"context"
	sitter "github.com/smacker/go-tree-sitter"
)

type Parser struct {
	engine *sitter.Parser
}

func NewParser(lang *sitter.Language) *Parser {
	engine := sitter.NewParser()
	engine.SetLanguage(lang)
	return &Parser{
		engine,
	}
}

func (p *Parser) ParseCtx(data []byte, context context.Context) ([]Symbol, error) {
	tree, err := p.engine.ParseCtx(context, nil, data)
	if err != nil {
		return nil, err
	}
	return p.node2Symbols(data, tree.RootNode())
}

func (p *Parser) Parse(data []byte) ([]Symbol, error) {
	return p.ParseCtx(data, context.TODO())
}

func (p *Parser) node2Symbols(data []byte, rootNode *sitter.Node) ([]Symbol, error) {
	var ret []Symbol
	count := int(rootNode.NamedChildCount())
	for i := 0; i < count; i++ {
		curChild := rootNode.NamedChild(i)
		curSymbol, err := p.node2Symbol(data, curChild)

		if err != nil {
			return nil, err
		}

		ret = append(ret, curSymbol)
		// handle its sons
		subSymbols, err := p.node2Symbols(data, curChild)
		if err != nil {
			return nil, err
		}
		ret = append(ret, subSymbols...)
	}
	return ret, nil
}

func (p *Parser) node2Symbol(data []byte, node *sitter.Node) (Symbol, error) {
	ret := Symbol{}
	ret.NodeType = node.Type()
	ret.Span = Span{
		Start: Point{node.StartPoint().Row, node.StartPoint().Column},
		End:   Point{node.EndPoint().Row, node.EndPoint().Column},
	}
	// symbol content
	ret.Symbol = node.Content(data)
	return ret, nil
}
