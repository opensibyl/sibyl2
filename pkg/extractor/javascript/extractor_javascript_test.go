package javascript

import (
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

var jsCode = `
class Pen {
    constructor(name, color, price){
        this.name = name;
        this.color = color; 
        this.price = price;
    }
    
    showPrice(itsPrice, arg2){
        console.log("Price of ${this.name} is ${this.price}");
    }
}

const pen1 = new Pen("Marker", "Blue", "$3");
pen1.showPrice();

function Pen(name, color, price) {
    this.name = name;
    this.color = color;
    this.price = price;
}

const pen1 = new Pen("Marker", "Blue", "$3");

Pen.prototype.showPrice = function(){
    console.log("Price of ${this.name} is ${this.price}");
}

pen1.showPrice();
`

func TestExtractor_ExtractSymbols(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangJavaScript)
	units, err := parser.Parse([]byte(jsCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	symbols, err := extractor.ExtractSymbols(units)
	assert.Nil(t, err)
	assert.NotEmpty(t, symbols)
}

func TestExtractor_ExtractFunctions(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangJavaScript)
	units, err := parser.Parse([]byte(jsCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	functions, err := extractor.ExtractFunctions(units)
	assert.Nil(t, err)
	assert.NotEmpty(t, functions)

	for _, each := range functions {
		core.Log.Infof("func: %v", each.Name)
		core.Log.Infof("func params: %v", each.Parameters)
	}
}

func TestExtractor_ExtractClasses(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangJavaScript)
	units, err := parser.Parse([]byte(jsCode))
	if err != nil {
		panic(err)
	}

	extractor := &Extractor{}
	classes, err := extractor.ExtractClasses(units)
	assert.Nil(t, err)
	assert.NotEmpty(t, classes)

	for _, each := range classes {
		core.Log.Infof("class: %v", each.Name)
	}
}
