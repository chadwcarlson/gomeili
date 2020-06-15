package templates

import (
  "io"
  "os"
  "log"
  "fmt"
  "strings"
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

func getNameFromRuntimeType(runtime string) string {
  url := "https://docs.platform.sh/registry/images/"

  body := req.RequestData(url, "registry.json")
  dynamic := make(map[string]interface{})
  json.Unmarshal(body, &dynamic)

  return dynamic[runtime].(map[string]interface{})["name"].(string)
}

func formatMultiAppString(runtimesWDuplicates []string) string {

  // Remove duplicates, such as in cases of a multi-app project with two PHP apps.
  runtimes := removeDuplicateRuntimeFromSlice(runtimesWDuplicates)

  // Catch cases where a single app's `.platform.app.yaml` file is not in root.
  if len(runtimes) == 1 {
    return runtimes[0]
  } else {
    // Join runtime strings.
    return strings.Join(runtimes, "/")
  }

}

func removeDuplicateRuntimeFromSlice(s []string) []string {
	m := make(map[string]bool)
	for _, item := range s {
		if _, ok := m[item]; ok {
			// Duplicate runtime.
		} else {
			m[item] = true
		}
	}

	var result []string
	for item, _ := range m {
		result = append(result, item)
	}
	return result
}

func getTemplateRuntime(p config.Config, template structs.TemplateInfo) string {

  var encodedContent structs.TemplateInfo
  body := req.RequestDataAuth(p.URL, fmt.Sprintf("/%s/files/.platform.app.yaml?ref=master", template.Name), "GITHUBTOKEN")
  json.Unmarshal(body, &encodedContent)
  var appYAML structs.PlatformAppYAML
  decodedContent, err := base64.StdEncoding.DecodeString(encodedContent.Content)
  if err != nil {
    log.Fatal(err)
  }
  yaml.Unmarshal(decodedContent, &appYAML)

  runtime := strings.Split(appYAML.Type, ":")

  // Case 1: Single app with `.platform.app.yaml` in root directory.
  if len(runtime[0]) > 0 {
    return getNameFromRuntimeType(runtime[0])
  } else {
    // Case 2: Multi-app with `.platform.app.yaml` in multiple subdirectories.
    var searchResults structs.PlatformAppYAMLSearchResults
    body := req.RequestDataAuth("https://api.github.com/search/code", fmt.Sprintf("?q=filename:.platform.app.yaml+repo:platformsh-templates/%s", template.Name), "GITHUBTOKEN")
    json.Unmarshal(body, &searchResults)

    runtimes := []string{}

    for _, file := range searchResults.Items {

      var encodedContent structs.TemplateInfo
      body := req.RequestDataAuth(p.URL, fmt.Sprintf("/%s/files/%s?ref=master", template.Name, file.Path), "GITHUBTOKEN")
      json.Unmarshal(body, &encodedContent)
      var appYAML structs.PlatformAppYAML
      decodedContent, err := base64.StdEncoding.DecodeString(encodedContent.Content)
      if err != nil {
        log.Fatal(err)
      }
      yaml.Unmarshal(decodedContent, &appYAML)

      runtime := strings.Split(appYAML.Type, ":")
      runtimes = append(runtimes, getNameFromRuntimeType(runtime[0]))

    }

    return formatMultiAppString(runtimes)
  }
  // Case 3: `.platform/applications.yaml` file present. No templates have this yet, will test when there is.
}

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

  runtimeSection := getTemplateRuntime(p, template)
  document.Section = runtimeSection
  document.Subsection = runtimeSection

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
