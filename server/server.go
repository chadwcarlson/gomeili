package server

import (
    "io"
    "os"
    // "fmt"
    "log"
    "github.com/meilisearch/meilisearch-go"
    config "github.com/chadwcarlson/gomeili/server/config"
    docs "github.com/chadwcarlson/gomeili/index/documents"
)

// Posts a prepared Index to a Meilisearch server.
func Post(config config.Config, documents docs.Index) {

    io.WriteString(os.Stdout, "\n\033[1mPosting index to Meilisearch...\033[0m\n")

    // Define the Meilisearch client.
    var client = meilisearch.NewClient(meilisearch.Config{
        Host: "http://127.0.0.1:7700",
    })

    // Create the index.
    createIndex(client, config)

    // Update the settings.
    updateSettings(client, config)

    // Add the documents.
    addDocuments(client, config, documents)

    io.WriteString(os.Stdout, "\n\n\033[1mDone.\033[0m\n")

}

// Initializes a new index, deleting old ones first.
func createIndex(client *meilisearch.Client, config config.Config) {

  // Delete the index if it already exists.
  deleteIndex(client, config)

  _, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
      UID: config.Index.UID,
  })
  handleErr(err, "\n- Index created.")

}

// Deletes a Meilisearch index from the server.
func deleteIndex(client *meilisearch.Client, config config.Config) {
  numIndexes, err := client.Indexes().List()
  handleErr(err, "\n- Current index detecteed")
  // If index already exists for the indexUID, delete it.
  if len(numIndexes) > 0 {
    _, err := client.Indexes().Delete(config.Index.UID)
    handleErr(err, "\n- Previous index deleted.")

  }
}

// Bulk updates the Meilisearch server's settings.
func updateSettings(client *meilisearch.Client, config config.Config) {
  // Update the index's name.
  updateName(client, config)

  // Set the primary search key.
  updatePrimaryKey(client, config)

  // Update synonyms.
  updateSynonyms(client, config)

  // Update stop words that won't factor into search.
  updateStopWords(client, config)

  // Update ranking rules.
  updateRankingRules(client, config)

  // Update searchable attributes.
  updateSearchableAttributes(client, config)

  // Update displayed attributes.
  updateDisplayedAttributes(client, config)

  // Update distinct attributes.
  updateDistinctAttributes(client, config)
}

// Updates the index's front-facing name.
func updateName(client *meilisearch.Client, config config.Config) {
  _, err := client.Indexes().UpdateName(config.Index.UID, config.Index.Name)
  handleErr(err, "\n- Index name set.")
}

// Sets the primary key used by Meilisearch to differentiate its documents.
func updatePrimaryKey(client *meilisearch.Client, config config.Config) {
  _, err := client.Indexes().UpdatePrimaryKey(config.Index.UID, config.Index.PrimaryKey)
  handleErr(err, "\n- Primary key set.")
}

// Updates synonyms - ties queries to other more exact queries.
func updateSynonyms(client *meilisearch.Client, config config.Config) {

  _, err := client.Settings(config.Index.UID).UpdateSynonyms(config.Synonyms)
  handleErr(err, "\n- Synonyms applied.")
}

// Updates stop words, that Meilisearch will ignore during queries.
func updateStopWords(client *meilisearch.Client, config config.Config) {
  _, err := client.Settings(config.Index.UID).UpdateStopWords(config.StopWords)
  handleErr(err, "\n- Stop words updated.")

}

// Sets the ranking rules specified.
func updateRankingRules(client *meilisearch.Client, config config.Config) {

  _, err := client.Settings(config.Index.UID).UpdateRankingRules(config.RankingRules)
  handleErr(err, "\n- Ranking rules updated.")

}

// Sets attributes to be included in search algorithm.
func updateSearchableAttributes(client *meilisearch.Client, config config.Config) {

  _, err := client.Settings(config.Index.UID).UpdateSearchableAttributes(config.Attributes.Searchable)
  handleErr(err, "\n- Searchable attributes updated.")

}

// Sets attributes available for results display.
func updateDisplayedAttributes(client *meilisearch.Client, config config.Config) {

  _, err := client.Settings(config.Index.UID).UpdateDisplayedAttributes(config.Attributes.Displayed)
  handleErr(err, "\n- Displayed attributes updated.")

}

// Sets distinct attribute, which is used when handling large pages that have been split into multiple documents.
func updateDistinctAttributes(client *meilisearch.Client, config config.Config) {

  _, err := client.Settings(config.Index.UID).UpdateDistinctAttribute(config.Attributes.Distinct)
  handleErr(err, "\n- Distinct attribute updated.")

}

// Adds prepared Meilisearch documents to an index.
func addDocuments(client *meilisearch.Client, config config.Config, documents docs.Index) {

  // Check if any Text fields have numWords > Meilisearch's hard limit (1000). If so, split into smaller documents.
  //    This is a quick implementation to solve this problem. It will likely remain as a fallback when this case
  //    is hit, behind the parsers themselves (and self-indexers) splitting by section.
  //    For now, this is relevant for Public Docs, Marketing Site (i.e. Blogs), and a few Community Posts.
  hardcodedWordLimit := 1000
  var freshIndex docs.Index
  for _, document := range documents.Documents {
    if docs.NumWords(document.Text) >= hardcodedWordLimit {
      splitDocuments := docs.SplitDocument(document, hardcodedWordLimit)
      freshIndex.Documents = append(freshIndex.Documents, splitDocuments.Documents...)
    } else {
      freshIndex.Documents = append(freshIndex.Documents, document)
    }
  }

  _, err := client.Documents(config.Index.UID).AddOrUpdate(freshIndex.Documents)
  handleErr(err, "\n- Documents added.")
}

// Generic short-hand error-handler.
func handleErr(err error, successString string) {
  if err != nil {
      log.Fatal(err)
      os.Exit(1)
  } else {
    io.WriteString(os.Stdout, successString)
  }
}
