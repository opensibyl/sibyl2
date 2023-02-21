package kotlin_test

import (
	"testing"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/kotlin"
	"github.com/stretchr/testify/assert"
)

var kotlinCode = `
package pa.ck.age
import java.util.Scanner

fun main(args: Array<String>) {

    // Creates a reader instance which takes
    // input from standard input - keyboard
    print("Enter a number: ")

    // nextInt() reads the next integer from the keyboard
    var integer:Int = reader.nextInt()

    // println() prints the following line to the output screen
    println("You entered: $integer")
}

class MMMM {
	fun a() {
		println("ok")
	}
}

interface NNN {
	fun ABCD()
}
`

func TestDfsOnly(t *testing.T) {
	t.Skip("debug only")
	parser := core.NewParser(core.LangKotlin)
	units, err := parser.Parse([]byte(kotlinCode))
	if err != nil {
		panic(err)
	}

	for _, each := range units {
		core.DebugDfs(each, 0)
	}
}

func TestExtractor_ExtractSymbols(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangKotlin)
	units, err := parser.Parse([]byte(kotlinCode))
	if err != nil {
		panic(err)
	}

	extractor := &kotlin.Extractor{}
	symbols, err := extractor.ExtractSymbols(units)
	assert.Nil(t, err)
	assert.NotEmpty(t, symbols)
}

func TestExtractor_ExtractFunctions(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangKotlin)
	units, err := parser.Parse([]byte(kotlinCode))
	if err != nil {
		panic(err)
	}

	extractor := &kotlin.Extractor{}
	funcs, err := extractor.ExtractFunctions(units)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(funcs))
	assert.Equal(t, funcs[0].Namespace, "pa.ck.age")
}

func TestExtractor_ExtractClasses(t *testing.T) {
	t.Parallel()
	parser := core.NewParser(core.LangKotlin)
	units, err := parser.Parse([]byte(kotlinCode))
	if err != nil {
		panic(err)
	}

	extractor := &kotlin.Extractor{}
	classes, err := extractor.ExtractClasses(units)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(classes))
}
