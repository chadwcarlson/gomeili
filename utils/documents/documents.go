package documents

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
)

type Index struct {
  Documents []Document
}

type Document struct {
	Site        string `json:"site"`
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


func (d *Index) Write(save_location string) {

  data, _ := json.MarshalIndent(d.Documents, "", "    ")
  err := ioutil.WriteFile(save_location, data, 0644)
  fmt.Sprintln(err)

}
