# sibyl 2

## why

### about sibyl 1

see [sibyl 1](https://github.com/williamfzc/sibyl#2022-09-24)

### source code vs artifact

There are some artifact analyzers like [soot](https://github.com/soot-oss/soot), which can also extract metadata from your artifact (jar/class).
But its result will contain lots of noise (code/bytecode ...) which was injected (via ASM/javassist) in compile time.
Our target is offering some helps to developers, and developers talk to machine via source code, not artifacts.

## arch

![](docs/arch.png)

## refs

- basic grammar: https://tree-sitter.github.io/tree-sitter/creating-parsers#the-grammar-dsl
- language parser (for example, golang): https://github.com/tree-sitter/tree-sitter-go/blob/master/src/parser.c 
- symbol: https://github.com/github/semantic/blob/main/docs/examples.md#symbols
- stack graphs: https://github.blog/2021-12-09-introducing-stack-graphs/
