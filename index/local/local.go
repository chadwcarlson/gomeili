package local

import (
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "encoding/json"
  "github.com/chadwcarlson/gomeili/config"
  docs "github.com/chadwcarlson/gomeili/index/documents"
)

func Get(p config.Config) docs.Index {

  io.WriteString(os.Stdout, fmt.Sprintf("\n\033[1mLocal Meilisearch index @\033[0m %s\n", p.File))

  var allDocuments docs.Index

  data, err := ioutil.ReadFile(p.File)
  if err != nil {
    fmt.Print(err)
  }
  err = json.Unmarshal(data, &allDocuments.Documents)
  if err != nil {
      fmt.Println("error:", err)
  }

  // Update attributes to match configuration file.
  for position, _ := range allDocuments.Documents {

    allDocuments.Documents[position].Site = p.Name
    allDocuments.Documents[position].Rank = p.Rank
    if p.Rank == 1 {
      allDocuments.Documents[position].Source = "primary"
    } else {
      allDocuments.Documents[position].Source = "secondary"
    }

  }

  return allDocuments

}
