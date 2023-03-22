# sibyl 2

> An easy-to-use logical layer on codebase.

[中文文档](https://opensibyl.github.io/doc/docs/intro)

## Status

| Name           | Badge                                                                                                                                                                 |
|----------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Latest Version | ![GitHub release (latest by date)](https://img.shields.io/github/v/release/opensibyl/sibyl2)                                                                          |
| Unit Tests     | [![Go](https://github.com/opensibyl/sibyl2/actions/workflows/ci.yml/badge.svg)](https://github.com/opensibyl/sibyl2/actions/workflows/ci.yml)                         |
| Docker Image   | [![ImageBuild](https://github.com/opensibyl/sibyl2/actions/workflows/imagebuild.yml/badge.svg)](https://github.com/opensibyl/sibyl2/actions/workflows/imagebuild.yml) |
| Perf Tests     | [![perftest](https://github.com/opensibyl/sibyl2/actions/workflows/perf.yml/badge.svg)](https://github.com/opensibyl/sibyl2/actions/workflows/perf.yml)               |
| Code Coverage  | [![codecov](https://codecov.io/github/opensibyl/sibyl2/branch/master/graph/badge.svg?token=1DuAXh12Ys)](https://codecov.io/github/opensibyl/sibyl2)                   |
| Code Style     | [![Go Report Card](https://goreportcard.com/badge/github.com/opensibyl/sibyl2)](https://goreportcard.com/report/github.com/opensibyl/sibyl2)                          |

## Overview

sibyl2 is a static code analyze service, for extracting, managing and offering metadata of your code in codebase. Inspired
by [semantic](https://github.com/github/semantic) of GitHub.

- Easy to use
- Fast enough
- Extensible
- Multiple languages in one (Go/Java/Python ...)

## What's `logical layer`?

SCM (GitHub, for example) manages code as plain text. We call it `physical layer`.

```golang
func TestExtractString(t *testing.T) {
    fileResult, err := ExtractFromString(javaCodeForExtract, &ExtractConfig{
        LangType:    core.LangJava,
        ExtractType: extractor.TypeExtractFunction,
    })
    if err != nil {
        panic(err)
    }
    for _, each := range fileResult.Units {
        core.Log.Debugf("result: %s", each.GetDesc())
    }
}
```

sibyl2 manages metadata of code. We call it `logical layer`.

```json
{
  "_id": {
    "$oid": "641b085deae764d271e2f426"
  },
  "repo_id": "sibyl2",
  "rev_hash": "e995ef44372a93394199ea837b1e2eed375a71a0",
  "path": "extract_test.go",
  "signature": "sibyl2||TestExtractString|*testing.T|",
  "tags": [],
  "func": {
    "name": "TestExtractString",
    "receiver": "",
    "namespace": "sibyl2",
    "parameters": [
      {
        "type": "*testing.T",
        "name": "t"
      }
    ],
    "returns": null,
    "span": {
      "start": {
        "row": {
          "$numberLong": "34"
        },
        "column": {
          "$numberLong": "0"
        }
      },
      "end": {
        "row": {
          "$numberLong": "45"
        },
        "column": {
          "$numberLong": "1"
        }
      }
    },
    "extras": {},
    "lang": "GOLANG"
  }
}
```

## Purpose & Principles

We hope to provide a unified logical layer for different tools in the entire DevOps process, 
sharing a single data source, 
rather than each tool performing its own set of duplicate parsing logic.

See [About This Project: Code Snapshot Layer In DevOps](https://github.com/opensibyl/sibyl2/issues/2) for details.

## Languages support

| Languages  | Function | Function Context | Class |
|------------|----------|------------------|-------|
| Golang     | Yes      | Yes              | Yes   |
| Java       | Yes      | Yes              | Yes   |
| Python     | Yes      | Yes              | Yes   |
| Kotlin     | Yes      | Yes              | Yes   |
| JavaScript | Yes      | Yes              | Yes   |

Based on tree-sitter, it's very easy to add an extra language support.

## Try it in 3 minutes

sibyl2 supports multiple database backends. 
It can also run with no middleware and database installed, 
if you just want to take a try.

### Deployment

Users can access all the features with a simple binary file, without any extra dependencies and scripts.

You can download from [the release page](https://github.com/opensibyl/sibyl2/releases).

Or directly download with `wget` (linux only):

```bash
curl https://raw.githubusercontent.com/opensibyl/sibyl2/master/scripts/download_latest.sh | bash
```

Now you can start it:

```bash
./sibyl server
```

That's it.
Server will run on port `:9876`.
Data will be persisted in `./sibyl2Storage`.

### Upload

![](https://opensibyl.github.io/doc/assets/images/intro-upload-1bb4fa2ce8ed43e6fc5f31c1ab3cc90b.gif)

```bash
./sibyl upload --src . --url http://127.0.0.1:9876
```

You can upload from different machines. Usually it only takes a few seconds.

### Access

Now all the data is ready! We have a built-in dashboard for visualization. Start it with:

```bash
./sibyl frontend
```

And open `localhost:3000` you will see:

<img width="877" alt="image" src="https://user-images.githubusercontent.com/13421694/216641341-c01bbcd1-349f-4934-bd35-2fa6b2c48cb4.png">

Of course, at the most time, we access data programmatically. 
You can access all the data via different kinds of languages, to build your own tools:

For example, git diff with logical?

```go
// assume that we have edited these lines
affectedFileMap := map[string][]int{
    "pkg/core/parser.go": {4, 89, 90, 91, 92, 93, 94, 95, 96},
    "pkg/core/unit.go":   {27, 28, 29},
}

for fileName, lineList := range affectedFileMap {
    strLineList := make([]string, 0, len(lineList))
    for _, each := range lineList {
        strLineList = append(strLineList, strconv.Itoa(each))
    }

    affectedFunctions, _, err := apiClient.BasicQueryApi.
        ApiV1FuncctxGet(ctx).
        Repo(projectName).
        Rev(head.Hash().String()).
        File(fileName).
        Lines(strings.Join(strLineList, ",")).
        Execute()
	
    for _, eachFunc := range affectedFunctions {
        // get all the calls details?
        for _, eachCall := range eachFunc.Calls {
            detail, _, err := apiClient.SignatureQueryApi.
                ApiV1SignatureFuncGet(ctx).
                Repo(projectName).
                Rev(head.Hash().String()).
                Signature(eachCall).
                Execute()
            assert.Nil(t, err)
            core.Log.Infof("call: %v", detail)
        }
    }
}
```

| Language   | Link                                                 |
|------------|------------------------------------------------------|
| Golang     | https://github.com/opensibyl/sibyl-go-client         |
| Java       | https://github.com/opensibyl/sibyl-java-client       |
| JavaScript | https://github.com/opensibyl/sibyl-javascript-client |

See more [examples](./_examples) about how to use for details.

## In Production

We use mongo db as our official backend in production.
All you need is adding a `sibyl-server-config.json` file:

```json
{
  "binding": {
    "dbtype": "MONGO",
    "mongodbname": "sibyl2",
    "mongouri": "mongodb+srv://feng:<YOURPASSWORD>@XXXXXXXX.mongodb.net/test"
  }
}
```

Everything done. 

<img width="706" alt="mongo_func_detail" src="https://user-images.githubusercontent.com/13421694/226957632-c1414be5-ec35-431b-9488-d6e0b1c0ddda.png">

## Performance

We have tested it on some famous repos, like [guava](https://github.com/google/guava). And that's why we can say it is "
fast enough".

See https://github.com/williamfzc/sibyl2/actions/workflows/perf.yml for details.

| Language | Repo                                               | Cost |
|----------|----------------------------------------------------|------|
| Golang   | https://github.com/gin-gonic/gin.git               | ~1s  |
| Java     | https://github.com/spring-projects/spring-boot.git | ~50s |
| Python   | https://github.com/psf/requests                    | ~1s  |
| Kotlin   | https://github.com/square/okhttp                   | ~1s  |

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

## References

- basic grammar: https://tree-sitter.github.io/tree-sitter/creating-parsers#the-grammar-dsl
- language parser (for example, golang): https://github.com/tree-sitter/tree-sitter-go/blob/master/src/parser.c
- symbol: https://github.com/github/semantic/blob/main/docs/examples.md#symbols
- stack graphs: https://github.blog/2021-12-09-introducing-stack-graphs/

## License

Apache License Version 2.0, see [LICENSE](LICENSE)
