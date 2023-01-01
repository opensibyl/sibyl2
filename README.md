# sibyl 2

> Take a quick snapshot of your codebase in seconds, with zero cost.

## Status

| Name           | Badge                                                                                                                                                                 |
|----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Latest Version | ![GitHub release (latest by date)](https://img.shields.io/github/v/release/opensibyl/sibyl2)                                                                          |
| Unit Tests     | [![Go](https://github.com/opensibyl/sibyl2/actions/workflows/ci.yml/badge.svg)](https://github.com/opensibyl/sibyl2/actions/workflows/ci.yml)                         |
| Docker Image   | [![ImageBuild](https://github.com/opensibyl/sibyl2/actions/workflows/imagebuild.yml/badge.svg)](https://github.com/opensibyl/sibyl2/actions/workflows/imagebuild.yml) |
| Perf Tests     | [![perftest](https://github.com/opensibyl/sibyl2/actions/workflows/perf.yml/badge.svg)](https://github.com/opensibyl/sibyl2/actions/workflows/perf.yml)               |
| Code Coverage  | [![codecov](https://codecov.io/github/opensibyl/sibyl2/branch/master/graph/badge.svg?token=1DuAXh12Ys)](https://codecov.io/github/opensibyl/sibyl2)                   |
| Code Style     | [![CodeFactor](https://www.codefactor.io/repository/github/opensibyl/sibyl2/badge)](https://www.codefactor.io/repository/github/opensibyl/sibyl2)                     |

## Overview

sibyl2 is a static code analyzer, for extracting, managing and offering codebase snapshot. Inspired
by [semantic](https://github.com/github/semantic) of GitHub.

- Easy to use
- Fast enough
- Extensible
- Multiple languages in one (Go/Java/Python ...)

## What's `Codebase Snapshot`?

![](https://opensibyl.github.io/doc/assets/images/intro-summary-0043f5cae91e9de62c619318afda4c39.png)

Raw source code:

```go
func ExtractFunction(targetFile string, config *ExtractConfig) ([]*extractor.FunctionFileResult, error) {
// ...
}
```

Code snapshot is the logical metadata of your code:

![](./docs/sample.svg)

## Purpose & Principles

See [About This Project: Code Snapshot Layer In DevOps](https://github.com/opensibyl/sibyl2/issues/2) for details.

## Languages support

| Languages | Function | Function Context | Class   |
|-----------|----------|------------------|---------|
| Golang    | Yes      | Yes              | Yes     |
| Java      | Yes      | Yes              | Yes     |
| Python    | Yes      | Yes              | Not yet |

## Examples of Usage

### One-file-installation

For now, we are aiming at offering an out-of-box service.
Users can access all the features with a simple binary file, without any extra dependencies and scripts.

You can download from [the release page](https://github.com/opensibyl/sibyl2/releases).

Or directly download with `wget` (replace `x.y.z` to the latest
version: ![GitHub release (latest by date)](https://img.shields.io/github/v/release/opensibyl/sibyl2)    ):

```bash
wget https://github.com/opensibyl/sibyl2/releases/download/v<x.y.z>/sibyl2_<x.y.z>_linux_amd64
```

### Use as a service (recommend)

#### Deploy

```bash
./sibyl server
```

That's it.
Server will run on port `:9876`.
Data will be persisted in `./sibyl2Storage`.

#### Upload

![](https://opensibyl.github.io/doc/assets/images/intro-upload-1bb4fa2ce8ed43e6fc5f31c1ab3cc90b.gif)

```bash
./sibyl upload --src . --url http://127.0.0.1:9876
```

You can upload from different machines.

#### Access

After uploading, you can access your data via http api.

![](https://opensibyl.github.io/doc/assets/images/usage-swagger-82e6fbaf8ae27f8cf697eb77cad56210.png)

Tree-like storage:

- repo
    - rev1
        - file
            - function
    - rev2
        - file
            - function

Try with swagger: http://127.0.0.1:9876/swagger/index.html#/

#### Access with sdk

Easily combine with other systems:

- golang: https://github.com/opensibyl/sibyl-go-client
- java: https://github.com/opensibyl/sibyl-java-client

### Use as a commandline tool

#### Basic Functions

```
./sibyl extract --src . --output hello.json
```

You will see:

```
$ ./sibyl extract --src . --output hello.json
{"level":"info","ts":1670138890.5306911,"caller":"sibyl2/extract_fs.go:92","msg":"no specific lang found, do the guess in: /Users/fengzhangchi/github_workspace/sibyl2"}
{"level":"info","ts":1670138890.5596569,"caller":"sibyl2/extract_fs.go:97","msg":"I think it is: GOLANG"}
{"level":"info","ts":1670138890.5836658,"caller":"core/runner.go:22","msg":"valid file count: 55"}
{"level":"info","ts":1670138890.6657321,"caller":"sibyl2/extract_fs.go:76","msg":"cost: 135 ms"}
{"level":"info","ts":1670138890.669896,"caller":"extract/cmd_extract.go:60","msg":"file has been saved to: hello.json"}
```

<details>
<summary> ... Result will be generated in seconds. </summary>

```json title="hello.json"
[
  {
    "path": "analyze.go",
    "language": "GOLANG",
    "type": "func",
    "units": [
      {
        "name": "AnalyzeFuncGraph",
        "receiver": "",
        "parameters": [
          {
            "type": "[]*extractor.FunctionFileResult",
            "name": "funcFiles"
          },
          {
            "type": "[]*extractor.SymbolFileResult",
            "name": "symbolFiles"
          }
        ],
        "returns": [
          {
            "type": "*FuncGraph",
            "name": ""
          },
          {
            "type": "error",
            "name": ""
          }
        ],
        "span": {
          "start": {
            "row": 11,
            "column": 0
          },
          "end": {
            "row": 80,
            "column": 1
          }
        },
        "extras": {}
      }
    ]
  },
  ...
]
```

</details>

#### Source Code History Visualization

Source code history visualization, inspired by https://github.com/acaudwell/Gource

One line command to see how your repository grow up, with no heavy dependencies like OpenGL, with logic level messages.

```bash
./sibyl history --src . --output hello.html --full
```

https://user-images.githubusercontent.com/13421694/207089314-21b0d48d-00d1-4de5-951c-415fed74c78f.mp4

> You can remove the `full` flag for better performance.

#### Smart Git Diff

Normal git diff has only text level messages.

```bash
./sibyl diff --from HEAD~1 --to HEAD --output hello1.json
```

<details><summary>And you can get a structural one with sibyl, which contains method level messages and callgraphs.</summary>

```bash
{
  "fragments": [
    {
      "path": "pkg/server/admin_s.go",
      "functions": [
        {
          "name": "HandleStatusUpload",
          "receiver": "",
          "parameters": [
            {
              "type": "*gin.Context",
              "name": "c"
            }
          ],
          "returns": null,
          "span": {
            "start": {
              "row": 17,
              "column": 0
            },
            "end": {
              "row": 23,
              "column": 1
            }
          },
          "extras": {},
          "path": "pkg/server/admin_s.go",
          "language": "GOLANG",
          "calls": null,
          "reverseCalls": [
            {
              "name": "Execute",
              "receiver": "",
              "parameters": [
                {
                  "type": "ExecuteConfig",
                  "name": "config"
                }
              ],
              "returns": null,
              "span": {
                "start": {
                  "row": 67,
                  "column": 0
                },
                "end": {
                  "row": 96,
                  "column": 1
                }
              },
              "extras": {},
              "path": "pkg/server/app.go",
              "language": "GOLANG"
            }
          ]
        }
      ]
    },
    ...
```

</details>

You can easily build some `smart test` tools above it.
For example, Google 's unittest speed up:

![intro-google](https://user-images.githubusercontent.com/13421694/207057947-894c1fb9-8ce4-4f7b-b5d3-88d220003e82.png)

## Performance

We have tested it on some famous repos, like [guava](https://github.com/google/guava). And that's why we can say it is "
fast enough".

See https://github.com/williamfzc/sibyl2/actions/workflows/perf.yml for details.

## Contribution

This project split into 3 main parts:

- /cmd: Pure command line tool for general usage
- /pkg/server: All-in-one service for production
- others: Shared api and core

Workflow:

- core: collect files and convert them to `Unit`.
- extract: classify and process units into functions, symbols.
- api: higher level analyze like callgraph

Issues / PRs are welcome!

## Refs

- basic grammar: https://tree-sitter.github.io/tree-sitter/creating-parsers#the-grammar-dsl
- language parser (for example, golang): https://github.com/tree-sitter/tree-sitter-go/blob/master/src/parser.c
- symbol: https://github.com/github/semantic/blob/main/docs/examples.md#symbols
- stack graphs: https://github.blog/2021-12-09-introducing-stack-graphs/

## License

Apache License Version 2.0, see [LICENSE](LICENSE)
