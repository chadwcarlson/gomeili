package githubrepo

import (
  "fmt"
  "encoding/json"
  req "indexer/utils/requests"
  templ "indexer/templates/structs"
)

func LocalJsonFile(file string) templ.TemplateList {

  _, b, _, _ := runtime.Caller(0)
  basepath := filepath.Dir(b)

  jsonFile, err := os.Open(basepath + file)
  if err != nil {
    fmt.Println(err)
  }
  byteValue, _ := ioutil.ReadAll(jsonFile)

  // we initialize our Users array
  var allTemplates templ.TemplateList

  json.Unmarshal(byteValue, &allTemplates)

  return allTemplates
}


func GetSearchIndex() templ.TemplateList {
  url_root := "https://api.github.com/repos/platformsh/template-builder/contents"
  var allTemplates templ.TemplateList

  body := req.RequestData(url_root, "/templates")
  json.Unmarshal(body, &allTemplates)

  for _, template := range allTemplates {
    fmt.Print(template.Name + "\n")
  }

  return allTemplates
}
