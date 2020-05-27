package index

import (
  // "fmt"
  docs "github.com/chadwcarlson/gomeili/utils/documents"
  "github.com/chadwcarlson/gomeili/discourse"
  "github.com/chadwcarlson/gomeili/utils/config"
)

func Build(configs config.ConfigSet) docs.Index {

  var allDocuments docs.Index
  for _, config := range configs {

    if config.Type == "discourse" {
      allDocuments = discourse.Get(config)
    }

    if config.Type == "openapi" {
      // fmt.Print(config)
    }

    if config.Type == "githubrepo" {
      // fmt.Print(config)
    }

  }

  return allDocuments

}
