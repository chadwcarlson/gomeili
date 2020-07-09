package index

import (
	"os"
	"io"
	"github.com/chadwcarlson/gomeili/config"
	docs "github.com/chadwcarlson/gomeili/index/documents"
	"github.com/chadwcarlson/gomeili/index/local"
	"github.com/chadwcarlson/gomeili/index/remote"
  "github.com/chadwcarlson/gomeili/parse"
)

// Writes a `docs.Index` struct to a local file, if `destination` is defined.
func writeIndividualIndex(indexConfig config.Config, index docs.Index) {
	var emptyConfig config.Config
	if emptyConfig.Destination != indexConfig.Destination {
		index.Write(indexConfig.Destination)
	}
}

// Builds Meilisearch index across a given config file.
func Build(configs config.ConfigSet, combinedFileLocation string) docs.Index {

  // Combined index.
	var allDocuments docs.Index
  io.WriteString(os.Stdout, "\n\033[1mCollecting documents for Meilisearch...\033[0m\n")

  // Range over the resources listed in `config.yaml`.
	for _, config := range configs {

    // Individual resource's index.
		var documents docs.Index

    // Remote, pre-parsed Meilisearch index.
    if config.Type == "remote" {
      documents = remote.Get(config)
    // Local, pre-parsed Meilisearch index (i.e. self-indexed site).
    } else if config.Type == "local" {
      documents = local.Get(config)
    // Resource requires parsing to build its index.
    } else {
      documents = parse.Parse(config)
    }

    // If `destination` is defined for resource, write a local copy.
    writeIndividualIndex(config, documents)

    // Append the combined index.
		allDocuments.Documents = append(allDocuments.Documents, documents.Documents...)
	}

	// Write the combined index to a local file.
	if len(combinedFileLocation) > 0 {
		io.WriteString(os.Stdout, "\n\033[1mWriting Combined index\033[0m\n")
		allDocuments.Write(combinedFileLocation)
	}

	return allDocuments

}
