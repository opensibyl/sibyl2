# sibyl 2

> Parsing, analyzing source code across many languages, and extracting their metadata easily.

跨语言、快速、简单地从你的源码中提取可序列化的元信息。

## 这是什么

这个项目定位是底层基础组件，将源码逻辑化。
简单来说就是，诸如哪个文件的哪个代码片段，对应到什么函数、类，实际意义是什么。

基于这一点，大多数上层工具都可以基于它而：

- 不再需要兼容多语言
- 不再需要苦恼如何从源码中提取扫描想要的信息
- 不依赖编译流程

before：

```go
func ExtractFunction(targetFile string, config *ExtractConfig) ([]*extractor.FunctionFileResult, error) {
// ...
}
```

after：

![](./docs/sample.svg)

或其他格式：

```json
{
  "path": "extract.go",
  "language": "GOLANG",
  "type": "func",
  "units": [
    {
      "name": "ExtractFunction",
      "receiver": "",
      "parameters": [
        {
          "type": "string",
          "name": "targetFile"
        },
        {
          "type": "*ExtractConfig",
          "name": "config"
        }
      ],
      "returns": [
        {
          "type": "[]*extractor.FunctionFileResult",
          "name": ""
        },
        {
          "type": "error",
          "name": ""
        }
      ],
      "span": {
        "start": {
          "row": 18,
          "column": 0
        },
        "end": {
          "row": 46,
          "column": 1
        }
      },
      "extras": null
    }
  ]
}
```

核心场景为研发过程、CI过程中进行高效的代码扫描与信息提取，你可以在 **毫秒级别-秒级别** 无痛为你的整个代码仓生成一份完整的快照图，供其他人、工具后续使用。

更多请见文档。

## 文档

https://github.com/williamfzc/sibyl2/wiki/0.-%E5%85%B3%E4%BA%8E

## refs

- basic grammar: https://tree-sitter.github.io/tree-sitter/creating-parsers#the-grammar-dsl
- language parser (for example, golang): https://github.com/tree-sitter/tree-sitter-go/blob/master/src/parser.c
- symbol: https://github.com/github/semantic/blob/main/docs/examples.md#symbols
- stack graphs: https://github.blog/2021-12-09-introducing-stack-graphs/

## license

Apache License Version 2.0, see [LICENSE](LICENSE)
