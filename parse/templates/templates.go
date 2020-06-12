package templates

import (
  "io"
  "os"
  "log"
  "fmt"
  "crypto/sha1"
  "encoding/json"
  "encoding/base64"
  "gopkg.in/yaml.v2"
  "github.com/schollz/progressbar/v3"
  "github.com/chadwcarlson/gomeili/config"
  "github.com/chadwcarlson/gomeili/utils/ignore"
  req "github.com/chadwcarlson/gomeili/utils/requests"
  docs "github.com/chadwcarlson/gomeili/index/documents"
  "github.com/chadwcarlson/gomeili/parse/templates/structs"
)

// func getNameFromRuntimeType(runtime string) string {
//   url := "https://docs.platform.sh/registry/images/"
//
//   body := req.RequestData(url, "registry.json")
//   dynamic := make(map[string]interface{})
//   json.Unmarshal(body, &dynamic)
//
//   return dynamic[runtime].(map[string]interface{})["name"]
// }

func getTemplate(p config.Config, template structs.TemplateInfo) docs.Document {

  var document docs.Document

  // Get and decode .platform.template.yaml.
  var encodedContent structs.TemplateInfo
  body := req.RequestDataAuth(p.URL, fmt.Sprintf("/%s/.platform.template.yaml?ref=master", template.Name), "GITHUBTOKEN")
  json.Unmarshal(body, &encodedContent)
  var templateYAML structs.PlatformTemplateYAML
  decodedContent, err := base64.StdEncoding.DecodeString(encodedContent.Content)
  if err != nil {
    log.Fatal(err)
  }
  yaml.Unmarshal(decodedContent, &templateYAML)

  // Basics.
  document.Site = p.Name
  document.Title = templateYAML.Info.Name

  // Url.
  url := fmt.Sprintf("https://github.com/platformsh-templates/%s", template.Name)
  document.URL = url
  document.RelativeURL = url

  // DocumentID.
  h := sha1.New()
  h.Write([]byte(fmt.Sprintf(url)))
  document.DocumentID = fmt.Sprintf("%x", h.Sum(nil))

  // Rank and source.
  document.Rank = p.Rank
  if p.Rank == 1 {
    document.Source = "primary"
  } else {
    document.Source = "secondary"
  }

  // Text.
  document.Description = templateYAML.Info.Description
  document.Text = templateYAML.Info.Description

  // Image.
  document.Image = templateYAML.Info.Image

  // Section (i.e. runtime).
  // get the runtime, used for the hugo object (and  make the hugo object in parallel)
  // Make request to files (Ideally....)
  // if .platform.app.yaml:
  //    parse file for type
  // elif .platform/applications.yaml:
  //    parse array for type on each app
  // else:
  //    don't include the template at  all
  // document.Section = category_data.Name   RUNTIME

  return document
}

func Get(p config.Config) docs.Index {

  io.WriteString(os.Stdout, fmt.Sprintf("\n\033[1mPlatform.sh templates @\033[0m %s\n", p.URL))

  var allDocuments docs.Index
  var templates structs.TemplateList

  body := req.RequestDataAuth(p.URL, "", "GITHUBTOKEN")
  json.Unmarshal(body, &templates.AllTemplates)
  io.WriteString(os.Stdout, fmt.Sprintf("\033[1m Templates\033[0m (%v on `master`)\n", len(templates.AllTemplates)))

  bar := progressbar.Default(int64(len(templates.AllTemplates)))
  for _, template := range templates.AllTemplates {
    if !ignore.ItemExists(p.Ignore, template.Name) {
      allDocuments.Documents = append(allDocuments.Documents, getTemplate(p, template))
    }
    bar.Add(1)
  }
  return allDocuments
}
