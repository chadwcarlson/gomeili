package server

import (
    "io"
    "os"
    "fmt"
    "log"
    "strconv"
    "github.com/meilisearch/meilisearch-go"
)

func deleteIndex(client *meilisearch.Client, indexUID string) {
  numIndexes, err := client.Indexes().List()
  handleErr(err)
  // If index already exists for the indexUID, delete it.
  if len(numIndexes) > 0 {
    _, err := client.Indexes().Delete(indexUID)
    handleErr(err)
    io.WriteString(os.Stdout, "- Previous index deleted.\n")
  }
}

func createIndex(client *meilisearch.Client, indexUID string) {
  _, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
      UID: indexUID,
  })
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("- Index created.\n"))
}

func updateSynonyms(client *meilisearch.Client, indexUID string) {
  synonyms := map[string][]string{
    "hobbit": []string{"automobiles"},
    "automobiles": []string{"hobbit"},
  }

  updateIDRes, err := client.Settings(indexUID).UpdateSynonyms(synonyms)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("- %s: Synonyms applied.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func updateRankingRules(client *meilisearch.Client, indexUID string) {
  rankingRules := []string{"typo", "words", "proximity", "attribute", "wordsPosition", "exactness"}

  updateIDRes, err := client.Settings(indexUID).UpdateRankingRules(rankingRules)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("- %s: Ranking rules updated.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func updateSearchableAttributes(client *meilisearch.Client, indexUID string) {
  searchableAttributes := []string{"id", "title", "description"}

  updateIDRes, err := client.Settings(indexUID).UpdateSearchableAttributes(searchableAttributes)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("- %s: Searchable attributes updated.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func updateDisplayedAttributes(client *meilisearch.Client, indexUID string) {
  displayedAttributes := []string{"id", "title", "description"}

  updateIDRes, err := client.Settings(indexUID).UpdateDisplayedAttributes(displayedAttributes)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("- %s: Displayed attributes updated.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func addDocuments(client *meilisearch.Client, indexUID string) {
  documents := []map[string]interface{}{
      {"book_id": 123,  "title": "Pride and Prejudice"},
      {"book_id": 456,  "title": "Le Petit Prince"},
      {"book_id": 1,    "title": "Alice In Wonderland"},
      {"book_id": 1344, "title": "The Hobbit"},
      {"book_id": 4,    "title": "Harry Potter and the Half-Blood Prince"},
      {"book_id": 42,   "title": "The Hitchhiker's Guide to the Galaxy"},
      {"book_id": 45,   "title": "Planes, Trains, and Automobiles"},
  }

  updateIDRes, err := client.Documents(indexUID).AddOrUpdate(documents)
  handleErr(err)
  io.WriteString(os.Stdout, fmt.Sprintf("- %s: Adding documents.\n", strconv.FormatInt(updateIDRes.UpdateID, 10) ))
}

func handleErr(err error) {
  if err != nil {
      log.Fatal(err)
      os.Exit(1)
  }
}

func Post() {

    io.WriteString(os.Stdout, "\n\033[1mPosting index to Meilisearch...\033[0m\n")

    indexUID := "books"

    // Define the Meilisearch client.
    var client = meilisearch.NewClient(meilisearch.Config{
        Host: "http://127.0.0.1:7700",
    })

    // Delete the index if it already exists.
    deleteIndex(client, indexUID)

    // Create an index if your index does not already exist
    createIndex(client, indexUID)

    // Update synonyms.
    updateSynonyms(client, indexUID)

    // Update ranking rules.
    updateRankingRules(client, indexUID)

    // Update searchable attributes.
    updateSearchableAttributes(client, indexUID)

    // Update displayed attributes.
    updateDisplayedAttributes(client, indexUID)

    // Add documents.
    addDocuments(client, indexUID)

}
