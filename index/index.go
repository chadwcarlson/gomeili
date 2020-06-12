package index

import (
  "os"
  // "fmt"
  "io"
  // "encoding/json"
  docs "github.com/chadwcarlson/gomeili/index/documents"
  "github.com/chadwcarlson/gomeili/parse/openapi"
  "github.com/chadwcarlson/gomeili/parse/discourse"
  "github.com/chadwcarlson/gomeili/index/local"
  "github.com/chadwcarlson/gomeili/index/remote"
  "github.com/chadwcarlson/gomeili/config"
)

func writeIndividualIndex(indexConfig config.Config, index docs.Index) {
  var emptyConfig config.Config
  if emptyConfig.Destination != indexConfig.Destination {
    index.Write(indexConfig.Destination)
  }
}

// Builds Meilisearch index across a given config file.
func Build(configs config.ConfigSet, combinedFileLocation string) docs.Index {

  io.WriteString(os.Stdout, "\n\033[1mPreparing Meilisearch index...\033[0m\n\n")
  var allDocuments docs.Index

  for _, config := range configs {
    var documents docs.Index
    // Handle Discourse source type.
    if config.Type == "discourse" {
      documents = discourse.Get(config)
      writeIndividualIndex(config, documents)
    }
    // Handle OpenAPI 3.0 source type.
    if config.Type == "openapi" {
      documents = openapi.Get(config)
      writeIndividualIndex(config, documents)
    }
    // Handle GitHub repo.
    if config.Type == "githubrepo" {
      // fmt.Print(config)
    }
    // Handle pre-existing remote Meilisearch index.
    if config.Type == "remote" {
      documents = remote.Get(config)
      writeIndividualIndex(config, documents)
    }
    // Handle pre-existing local Meilisearch index file.
    if config.Type == "local" {
      documents = local.Get(config)
      writeIndividualIndex(config, documents)
    }
    allDocuments.Documents = append(allDocuments.Documents, documents.Documents...)
  }

  // Write the combined index file.
  if len(combinedFileLocation) > 0 {
    io.WriteString(os.Stdout, "* \033[1mWriting Combined index\033[0m\n")
    allDocuments.Write(combinedFileLocation)
  }

  io.WriteString(os.Stdout, "\n\033[1mComplete.\033[0m\n\n")

  return allDocuments

}
