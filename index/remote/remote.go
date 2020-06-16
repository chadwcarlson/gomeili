package remote

import (
	"io"
	"os"
	"fmt"
	"encoding/json"
	"github.com/chadwcarlson/gomeili/config"
	docs "github.com/chadwcarlson/gomeili/index/documents"
	req "github.com/chadwcarlson/gomeili/utils/requests"
)

// Returns docs.Index object from a remote Meilisearch index.
func Get(p config.Config) docs.Index {

	// Get the remote resource.
	io.WriteString(os.Stdout, fmt.Sprintf("\n\033[1mRemote Meilisearch index @\033[0m %s\n", p.URL+p.File))
	var allDocuments docs.Index
	body := req.RequestData(p.URL, p.File)
	json.Unmarshal(body, &allDocuments.Documents)

	// Update attributes to match configuration file.
	for position, _ := range allDocuments.Documents {

		// Updated attributes.
		allDocuments.Documents[position].Site = p.Name
		allDocuments.Documents[position].Rank = p.Rank

		// Apply React's primary/secondary classification.
		if p.Rank == 1 {
			allDocuments.Documents[position].Source = "primary"
		} else {
			allDocuments.Documents[position].Source = "secondary"
		}

	}

	return allDocuments

}
