# gomeili


This library supports:

- Retrieving common resources with existing public APIs (i.e. GitHub, Discourse), parsing and reformatting their data into a Meilisearch compatible index
- Including a `rank` key in each document that allows multiple resources to be combined into a single index, providing the data to support cross-site search.

**meiliparsers** is a collection of packages intended to provide some helpful methods for quickly parsing common content sources, and then reformatting that content into an index that can be consumed by [Meilisearch](https://www.meilisearch.com/). These packages can be used as a faster alternative to scraping methods.

At this time, **meiliparsers** contains packages for that already contain public APIs or rigid specifications:

* [Discourse](https://docs.discourse.org/)
* [GitHub](https://developer.github.com/v3/)
* [OpenAPI 3](https://spec.openapis.org/oas/v3.0.3)

These packages rely heavily on pre-existing parsing libraries where they exist for each source:

* [Discourse - `godisco`](https://github.com/FrenchBen/godisco)
* [GitHub - `go-github`]()


## Demo

* Discourse -  https://try.discourse.org/
* GitHub - All repos in organization (with language guessed in master)
* GitHub - Central management repo (template-builder)
* OpenAPI 3/Redoc - https://redocly.github.io/redoc/ & https://redocly.github.io/redoc/openapi.yaml


### Usage

There are two main use cases for developers using this library.

1. A repository retrieves resources to build an index. This will be the case if a single site needs to be indexed and serve that index in it's own search, or in a central index repository where other sites are looking to its index as a kind of *ground truth* of the documents available.
2. A repository accesses a set of index files from somewhere else (like a central *ground truth*), and updates the index slightly prior to applying it to its own Meilisearch server. An example update would include modifying `rank` to prioritize search differently depending on the current site, or anonymizing `name` to simplify how results are handled in a React app providing the search bar and autocomplete dropdown.

#### Single use and ground truth repos

```
# config.yamls
- name: community
  type: discourse
  rank: 1
  url: "https://community.platform.sh"
  ignore:
    - "How-to Guides"
    - "Questions & Answers"
- name: templates
  type: githubrepo
  rank: 2
  url: "https://api.github.com/repos/platformsh/template-builder/contents"
  ignore:
    - "__init__.py"
- name: examples
  type: githuborg
  rank: 3
  url: "https://api.github.com/orgs/platformsh-examples/repos"
  ignore:
    - "quarkus"
- name: website
  type: remote
  rank: 4
  url: "https://platform.sh/index.json"
- name: apidocs
  type: openapi
  url: "https://api.platform.sh/docs/openapispec-platformsh.json"
  rank: 5
```

### Updating an index

You can update a set of indexes in a *ground truth repository* in the same file. Indexes that are exluded from this file will likewise be excluded from the final index, and new keys will be added to it. In the example below, a new local index `docs` is added, and the `rank` for the other resources have been updated so that it has top priority.

```
# config.yaml
- name: docs
  type: local
  rank: 1
  url: "output/documentation.json"
- name: community
  rank: 2
- name: templates
  rank: 3
- name: website
  rank: 4
- name: apidocs
  rank: 5
```

```go
package main

import (
  "fmt"
  "gomeilindex/discourse"
  // "gomeilindex/utils/documents"
)


func main() {

  // var params discourse.Params
  // params.Url = "https://community.platform.sh"
  // params.Name = "community"
  // params.Rank = 1
  // params.Ignore = []string{"How-to Guides", "Questions & Answers", "Activity scripts"}

  // params := discourse.Params {
  //   Url: "https://community.platform.sh",
  //   Name: "community",
  //   Rank: 1,
  //   Ignore: []string{"How-to Guides", "Questions & Answers", "Activity scripts"},
  // }
  // community := discourse.Get(params)

  var config discourse.Config
  community := discourse.Get(config.Load("config.yaml", "community"))

  fmt.Sprintln(community)

}

```
