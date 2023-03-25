# sibyl 2

> The missing logical layer in codebases.

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

sibyl2 is a service for extracting, managing and offering metadata of your code in codebase.

<img width="1217" alt="image" src="https://user-images.githubusercontent.com/13421694/227231024-3fff016d-4866-4061-8704-b8c9e4f880f3.png">

Although Git is a widely used platform for version control and collaboration, it does not have the capability to analyze and interpret the logic of code.

Assuming my program now wants to know what is on `line 35` of the `extract_test.go` file. If using GitLab, it is possible to find the corresponding text information, which may look like this:

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

However, my program does not know what this text represents. This is even more difficult for programs written in other languages (such as Java or Python).

With sibyl2 API, you can get:

```json
{
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

And even its relationships?

```json
{
  "name": "TestExtractString",
  
  ...,
  
  "calls": [
    "object|*Function|GetDesc||string",
    "object|*Symbol|GetDesc||string",
    "sibyl2||ExtractFromString|string,*ExtractConfig|*extractor.FileResult,error",
    "object|*Call|GetDesc||string",
    "object|*Clazz|GetDesc||string"
  ],
  "reverseCalls": []
}
```

And even all the relationships?

![](https://user-images.githubusercontent.com/13421694/219916928-14e8eb69-fe67-45a1-80a7-1b3c1b8163b2.png)

> Also class, module and something else. See [Reference](#References) for details.

In addition, sibyl2's unified logic layer can also convert different language logics into the same data structure for storage, which is friendly to OLAP and other tools.

You can use various mature analysis tools (such as clickhouse, superset, metabase, etc.) to further deconstruct your code repository based on sibyl2. For example, you can view the distribution of methods with the @Test annotation throughout the repository.

<img width="699" alt="image" src="https://user-images.githubusercontent.com/13421694/227699919-a6080730-ccea-42a3-b5ae-1ef3198426bb.png">

Applying some global analysis:

<img width="1347" alt="image" src="https://user-images.githubusercontent.com/13421694/227760234-d8c5244b-d65d-4d5d-b984-a0f154baf9ac.png">

Comparing the differences in methods, classes, and modules across each version:

<img width="1214" alt="image" src="https://user-images.githubusercontent.com/13421694/227701624-2a6fd71e-8733-480a-802c-ed1833652763.png">

In summary, sibyl2 aims to build a layer of common logic on top of your code repository in a simple way, helping businesses better understand, analyze, and use the data in their code repository.

## Deployment

### With Docker (recommended for start up)

We have provided an [official compose file](https://github.com/opensibyl/sibyl2/blob/master/docker-compose.yml). Just:

- copy and paste to your own `docker-compose.yml` file
- `docker-compose up`

### With an existed MongoDB

- download our binaries from [the release page](https://github.com/opensibyl/sibyl2/releases).
- add a `sibyl-server-config.json` file:

```json
{
  "binding": {
    "dbtype": "MONGO",
    "mongodbname": "sibyl2",
    "mongouri": "mongodb+srv://<USERNAME>:<YOURPASSWORD>@XXXXXXXX.mongodb.net/test"
  }
}
```

- `./sibyl server`

### With nothing

It can also run with no middleware and database installed, 
if you just want to take a try.

Just `./sibyl server` without config file.

## Usage

### Upload

![](https://opensibyl.github.io/doc/assets/images/intro-upload-1bb4fa2ce8ed43e6fc5f31c1ab3cc90b.gif)

```bash
./sibyl upload --src . --url http://127.0.0.1:9876
```

You can upload from different machines (just correct the url). Usually it only takes a few seconds.

Now everything is ready.

### Access with Mongo URI

There are many mature visualization tools that support MONGO URI, such as official compass and metabase. 
With this, accessing all your data is easy.

![](https://user-images.githubusercontent.com/13421694/226957632-c1414be5-ec35-431b-9488-d6e0b1c0ddda.png)

Our docker-compose file includes [metabase](https://github.com/metabase/metabase), allowing you to connect to MongoDB and start analyzing your data by simply opening `127.0.0.1:3000`.

Currently, our data is divided into three collections:

- fact_func: Function Information
- fact_clazz: Class Information
- rel_funcctx: Function Context Information

### Access with sibyl2 clients

Of course, at the most time, developers access data programmatically. 
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

## Languages support

| Languages  | Function | Function Context | Class |
|------------|----------|------------------|-------|
| Golang     | Yes      | Yes              | Yes   |
| Java       | Yes      | Yes              | Yes   |
| Python     | Yes      | Yes              | Yes   |
| Kotlin     | Yes      | Yes              | Yes   |
| JavaScript | Yes      | Yes              | Yes   |

Based on tree-sitter, it's very easy to add an extra language support.

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

## Further work

- More effective and stable uploader
- More languages support
- More fact level data we can collect
- Provide a standard data layer for AI models

About its role in DevOps, see [About This Project: Code Snapshot Layer In DevOps](https://github.com/opensibyl/sibyl2/issues/2) for details.

## References

Inspired by [semantic](https://github.com/github/semantic) of GitHub.

- basic grammar: https://tree-sitter.github.io/tree-sitter/creating-parsers#the-grammar-dsl
- language parser (for example, golang): https://github.com/tree-sitter/tree-sitter-go/blob/master/src/parser.c
- symbol: https://github.com/github/semantic/blob/main/docs/examples.md#symbols
- stack graphs: https://github.blog/2021-12-09-introducing-stack-graphs/

## License

Apache License Version 2.0, see [LICENSE](LICENSE)
