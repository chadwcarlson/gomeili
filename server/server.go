package server

import (
    "io"
    "os"
    "fmt"
    "log"
    "strconv"
    "github.com/meilisearch/meilisearch-go"
    config "github.com/chadwcarlson/gomeili/server/config"
    docs "github.com/chadwcarlson/gomeili/index/documents"
)

func Post(config config.Config, documents docs.Index) {

    io.WriteString(os.Stdout, "\n\033[1mPosting index to Meilisearch...\033[0m\n\n")

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

}

func createIndex(client *meilisearch.Client, config config.Config) {

  // Delete the index if it already exists.
  deleteIndex(client, config)

  _, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
      UID: config.Index.UID,
  })
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("- Index created.\n\n"))
}

func deleteIndex(client *meilisearch.Client, config config.Config) {
  numIndexes, err := client.Indexes().List()
  handleErr(err)
  // If index already exists for the indexUID, delete it.
  if len(numIndexes) > 0 {
    _, err := client.Indexes().Delete(config.Index.UID)
    handleErr(err)
    io.WriteString(os.Stdout, "- Previous index deleted.\n")
  }
}

func updateSettings(client *meilisearch.Client, config config.Config) {
  // Update the index's name.
  updateName(client, config)

  // Set the primary search key.
  updatePrimaryKey(client, config)

  // Update synonyms.
  updateSynonyms(client, config)

  // Update ranking rules.
  updateRankingRules(client, config)

  // Update searchable attributes.
  updateSearchableAttributes(client, config)

  // Update displayed attributes.
  updateDisplayedAttributes(client, config)
}

func updateName(client *meilisearch.Client, config config.Config) {
  _, err := client.Indexes().UpdateName(config.Index.UID, config.Index.Name)
  handleErr(err)
}

func updatePrimaryKey(client *meilisearch.Client, config config.Config) {
  _, err := client.Indexes().UpdatePrimaryKey(config.Index.UID, config.Index.PrimaryKey)
  handleErr(err)
}

func updateSynonyms(client *meilisearch.Client, config config.Config) {
  synonyms := map[string][]string{
    "hobbit": []string{"automobiles"},
    "automobiles": []string{"hobbit"},
  }

  updateIDRes, err := client.Settings(config.Index.UID).UpdateSynonyms(synonyms)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Synonyms applied.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func updateRankingRules(client *meilisearch.Client, config config.Config) {

  updateIDRes, err := client.Settings(config.Index.UID).UpdateRankingRules(config.RankingRules)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Ranking rules updated.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func updateSearchableAttributes(client *meilisearch.Client, config config.Config) {

  updateIDRes, err := client.Settings(config.Index.UID).UpdateSearchableAttributes(config.Attributes.Searchable)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Searchable attributes updated.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func updateDisplayedAttributes(client *meilisearch.Client, config config.Config) {

  updateIDRes, err := client.Settings(config.Index.UID).UpdateDisplayedAttributes(config.Attributes.Displayed)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Displayed attributes updated.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func addDocuments(client *meilisearch.Client, config config.Config, documents docs.Index) {

  updateIDRes, err := client.Documents(config.Index.UID).AddOrUpdate(documents.Documents)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Adding documents.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func handleErr(err error) {
  if err != nil {
      log.Fatal(err)
      os.Exit(1)
  }
}
