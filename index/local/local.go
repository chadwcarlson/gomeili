package local

import (
	"encoding/json"
	"fmt"
	"github.com/chadwcarlson/gomeili/config"
	docs "github.com/chadwcarlson/gomeili/index/documents"
	"io"
	"io/ioutil"
	"os"
)

// Returns docs.Index object from a local Meilisearch index file.
func Get(p config.Config) docs.Index {

	// Get the remote resource.
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
