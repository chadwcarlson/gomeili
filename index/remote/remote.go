package remote

import (
	"encoding/json"
	"fmt"
	"github.com/chadwcarlson/gomeili/config"
	docs "github.com/chadwcarlson/gomeili/index/documents"
	req "github.com/chadwcarlson/gomeili/utils/requests"
	"io"
	"os"
)

func Get(p config.Config) docs.Index {

	io.WriteString(os.Stdout, fmt.Sprintf("\n\033[1mRemote Meilisearch index @\033[0m %s\n", p.URL+p.File))

	var allDocuments docs.Index

	body := req.RequestData(p.URL, p.File)
	json.Unmarshal(body, &allDocuments.Documents)

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
