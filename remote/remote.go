package remote

import (
  "fmt"
  "io"
  "os"
  "encoding/json"
  "github.com/chadwcarlson/gomeili/utils/config"
  docs "github.com/chadwcarlson/gomeili/utils/documents"
  req "github.com/chadwcarlson/gomeili/utils/requests"
)

func Get(p config.Config) docs.Index {

  io.WriteString(os.Stdout, fmt.Sprintf("* \033[1mRemote Meilisearch index @\033[0m %s\n", p.URL + p.File))

  var allDocuments docs.Index

  body := req.RequestData(p.URL, p.File)
  json.Unmarshal(body, &allDocuments.Documents)

  // Update attributes to match configuration file.
  for position, _ := range allDocuments.Documents {

    allDocuments.Documents[position].Site = p.Name
    allDocuments.Documents[position].Rank = p.Rank

  }

  return allDocuments

}
