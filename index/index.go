package index

import (
  // "fmt"
  // "io/ioutil"
  // "encoding/json"
  docs "github.com/chadwcarlson/gomeili/utils/documents"
  "github.com/chadwcarlson/gomeili/openapi"
  "github.com/chadwcarlson/gomeili/discourse"
  "github.com/chadwcarlson/gomeili/local"
  "github.com/chadwcarlson/gomeili/remote"
  "github.com/chadwcarlson/gomeili/utils/config"
)

func writeIndividualIndex(indexConfig config.Config, index docs.Index) {
  var emptyConfig config.Config
  if emptyConfig.Destination != indexConfig.Destination {
    index.Write(indexConfig.Destination)
  }
}

// Builds Meilisearch index across a given config file.
func Build(configs config.ConfigSet) docs.Index {

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

  return allDocuments

}
//
// func Combine(configs config.ConfigSet, destination string) {
//   var allDocuments docs.Index
//   var emptyConfig *config.Config
//   for _, config := range configs {
//     if emptyConfig.Destination != config.Destination {
//       var documents docs.Index
//       data, err := ioutil.ReadFile(config.Destination)
//       if err != nil {
//         fmt.Print(err)
//       }
//       err = json.Unmarshal(data, &documents)
//       if err != nil {
//           fmt.Println("error:", err)
//       }
//       for _, document := range documents.Documents {
//         allDocuments.Documents = append(allDocuments.Documents, document)
//       }
//     }
//   }
//
//   allDocuments.Write(destination)
//
// }
