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

// Posts a prepared Index to a Meilisearch server.
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

// Initializes a new index, deleting old ones first.
func createIndex(client *meilisearch.Client, config config.Config) {

  // Delete the index if it already exists.
  deleteIndex(client, config)

  _, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
      UID: config.Index.UID,
  })
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("- Index created.\n\n"))
}

// Deletes a Meilisearch index from the server.
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

// Bulk updates the Meilisearch server's settings.
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

// Updates the index's front-facing name.
func updateName(client *meilisearch.Client, config config.Config) {
  _, err := client.Indexes().UpdateName(config.Index.UID, config.Index.Name)
  handleErr(err)
}

// Sets the primary key used by Meilisearch to differentiate its documents.
func updatePrimaryKey(client *meilisearch.Client, config config.Config) {
  _, err := client.Indexes().UpdatePrimaryKey(config.Index.UID, config.Index.PrimaryKey)
  handleErr(err)
}

// Updates synonyms - ties queries to other more exact queries.
func updateSynonyms(client *meilisearch.Client, config config.Config) {
  // synonyms := map[string][]string{
  //   "hobbit": []string{"automobiles"},
  //   "automobiles": []string{"hobbit"},
  // }
  // birdJson := `{"twoway":true,"list":[{"routes":["routes.yaml"]},{"services":["services.yaml"]},{"application":[".platform.app.yaml","app.yaml","applications.yaml"]},{"multi-app":["applications.yaml"]},{"regions":["public ip address"]}]}`
  // var result map[string]interface{}
  // json.Unmarshal([]byte(birdJson), &result)
  //
  // synonyms := result["list"].([]interface{})

  updateIDRes, err := client.Settings(config.Index.UID).UpdateSynonyms(config.Synonyms)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Synonyms applied.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

// Sets the ranking rules specified.
func updateRankingRules(client *meilisearch.Client, config config.Config) {

  updateIDRes, err := client.Settings(config.Index.UID).UpdateRankingRules(config.RankingRules)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Ranking rules updated.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

// Sets attributes to be included in search algorithm.
func updateSearchableAttributes(client *meilisearch.Client, config config.Config) {

  updateIDRes, err := client.Settings(config.Index.UID).UpdateSearchableAttributes(config.Attributes.Searchable)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Searchable attributes updated.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

// Sets attributes available for results display.
func updateDisplayedAttributes(client *meilisearch.Client, config config.Config) {

  updateIDRes, err := client.Settings(config.Index.UID).UpdateDisplayedAttributes(config.Attributes.Displayed)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Displayed attributes updated.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

// Adds prepared Meilisearch documents to an index.
func addDocuments(client *meilisearch.Client, config config.Config, documents docs.Index) {

  updateIDRes, err := client.Documents(config.Index.UID).AddOrUpdate(documents.Documents)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("%s - Adding documents.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

// Generic short-hand error-handler.
func handleErr(err error) {
  if err != nil {
      log.Fatal(err)
      os.Exit(1)
  }
}
