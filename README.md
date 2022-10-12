# sibyl 2

> Parsing, analyzing source code across many languages, and extracting their metadata easily.

快速、简单地从你的源码中提取可序列化的元信息。

## 使用

### 命令行

可以在 [release页面](https://github.com/williamfzc/sibyl2/releases) 下载对应平台的版本。

```bash
./sibyl2_0.2.0_darwin_amd64 extract --src ~/YOUR_SOURCE_CODE_DIR --lang GOLANG --type func
```

即可提取出整个仓库里所有的函数信息：

```json
[
  {
    "path": "foo/bar/a.go",
    "language": "GOLANG",
    "type": "func",
    "units": [
      {
        "name": "NormalFunc",
        "receiver": "",
        "parameters": [
          {
            "type": "*sitter.Language",
            "name": "lang"
          },
          {
            "type": "int",
            "name": "ok"
          }
        ],
        "returns": [
          {
            "type": "string",
            "name": "aaa"
          },
          {
            "type": "error",
            "name": "n"
          }
        ],
        "span": {
          "start": {
            "row": 8,
            "column": 0
          },
          "end": {
            "row": 10,
            "column": 1
          }
        }
      }
    ]
  }
]
```

### API

TODO

## FAQ

https://github.com/williamfzc/sibyl2/wiki/FAQ

## refs

- basic grammar: https://tree-sitter.github.io/tree-sitter/creating-parsers#the-grammar-dsl
- language parser (for example, golang): https://github.com/tree-sitter/tree-sitter-go/blob/master/src/parser.c
- symbol: https://github.com/github/semantic/blob/main/docs/examples.md#symbols
- stack graphs: https://github.blog/2021-12-09-introducing-stack-graphs/

## license

Apache License Version 2.0, see [LICENSE](LICENSE)
