package documents

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"crypto/sha1"
	// "github.com/chadwcarlson/gomeili/utils/ignore"
)

// Single Meilisearch Document schema.
type Document struct {
	Site        string `json:"site"`
	Source      string `json:"source"`
	Rank        int    `json:"rank"`
	DocumentID  string `json:"documentId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Text        string `json:"text"`
	Section     string `json:"section"`
	Subsection  string `json:"subsections"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	RelativeURL string `json:"relurl"`
}

// Meilisearch index of Documents.
type Index struct {
	Documents []Document
}

// Write an index to a local file.
func (d *Index) Write(save_location string) {

	data, _ := json.MarshalIndent(d.Documents, "", "    ")
	err := ioutil.WriteFile(save_location, data, 0644)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		io.WriteString(os.Stdout, fmt.Sprintf("\033[1m Done.\033[0m Index written to: %s\n", save_location))
	}

}

func NumWords(text string) int {
	return len(strings.Split(text, " "))
}

func SplitDocument(document Document, hardcodedWordLimit int) Index{

	var newIndex Index

	// hardcodedWordLimit := 1000
	split_text := strings.Split(document.Text, " ")

	numDocuments := NumWords(document.Text) / hardcodedWordLimit + 1
	extra := NumWords(document.Text) % hardcodedWordLimit
	if extra == 0 {
		numDocuments = numDocuments - 1
	}
	start := 0
	stop := hardcodedWordLimit

	for i:= 0; i<numDocuments; i++ {

		// Copy the original document and update its Text for the new one.
		newDocument := document

		if stop >= len(split_text) {
			newDocument.Text = strings.Join(split_text[start:], " ")
		} else {
			newDocument.Text = strings.Join(split_text[start:stop], " ")
			start = stop
			stop = stop + hardcodedWordLimit
		}

		// Give it a new unique DocumentID
		h := sha1.New()
		h.Write([]byte(fmt.Sprintf(strings.Repeat(document.URL, i) )))
		newDocument.DocumentID = fmt.Sprintf("%x", h.Sum(nil))

		// append to newIndex
		newIndex.Documents = append(newIndex.Documents, newDocument)

	}

	return newIndex

}
