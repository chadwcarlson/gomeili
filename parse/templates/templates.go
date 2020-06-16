package templates

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	stringFormat "github.com/chadwcarlson/gomeili/utils/string"
	"github.com/chadwcarlson/gomeili/config"
	docs "github.com/chadwcarlson/gomeili/index/documents"
	"github.com/chadwcarlson/gomeili/parse/templates/structs"
	"github.com/chadwcarlson/gomeili/utils/ignore"
	req "github.com/chadwcarlson/gomeili/utils/requests"
	"github.com/schollz/progressbar/v3"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"strings"
)

// Main method to retrieve Template-builder index for Meilisearch.
func Get(p config.Config) docs.Index {

	// Get list of templates template-builder/templates.
	io.WriteString(os.Stdout, fmt.Sprintf("\n\033[1mPlatform.sh templates @\033[0m %s\n", p.URL))
	var allDocuments docs.Index
	var templates structs.TemplateList
	body := req.RequestDataAuth(p.URL, "", "GITHUBTOKEN")
	json.Unmarshal(body, &templates.AllTemplates)
	io.WriteString(os.Stdout, fmt.Sprintf("\033[1m Templates\033[0m (%v on `master`)\n", len(templates.AllTemplates)))

	// Form document for each individual template.
	bar := progressbar.Default(int64(len(templates.AllTemplates)))
	for _, template := range templates.AllTemplates {
		// Skip directories/files config.yaml says to ignore.
		if !ignore.ItemExists(p.Ignore, template.Name) {
			allDocuments.Documents = append(allDocuments.Documents, getTemplate(p, template))
		}
		bar.Add(1)
	}
	return allDocuments
}

// Forms Meilisearch document for an individual template.
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

	// Match `config.yaml` rank, and use React primary/secondary designation.
	document.Rank = p.Rank
	if p.Rank == 1 {
		document.Source = "primary"
	} else {
		document.Source = "secondary"
	}

	// Text.
	document.Text = stringFormat.Clean(templateYAML.Info.Description)

	// Description.
	document.Description = stringFormat.Clean(templateYAML.Info.Description)

	// Image.
	document.Image = templateYAML.Info.Image

	// Retrieve the runtime used to use as section for the result.
	runtimeSection := getTemplateRuntime(p, template)
	document.Section = runtimeSection
	document.Subsection = runtimeSection

	return document
}

// Returns runtime (section string) for an individual template.
func getTemplateRuntime(p config.Config, template structs.TemplateInfo) string {

	// First assume the template is a single app, with .platform.app.yaml in root.
	var encodedContent structs.TemplateInfo
	body := req.RequestDataAuth(p.URL, fmt.Sprintf("/%s/files/.platform.app.yaml?ref=master", template.Name), "GITHUBTOKEN")
	json.Unmarshal(body, &encodedContent)
	var appYAML structs.PlatformAppYAML
	decodedContent, err := base64.StdEncoding.DecodeString(encodedContent.Content)
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(decodedContent, &appYAML)

	// Case 1: Single app with `.platform.app.yaml` in root directory.
	runtime := strings.Split(appYAML.Type, ":")
	if len(runtime[0]) > 0 {
		// Use Registry to get language name from runtime type.
		return getNameFromRuntimeType(runtime[0])
	// Case 2: Multi-app with `.platform.app.yaml` in multiple subdirectories.
	} else {
		// Perform a search for .platform.app.yaml files in the actual template repo.
		var searchResults structs.PlatformAppYAMLSearchResults
		body := req.RequestDataAuth("https://api.github.com/search/code", fmt.Sprintf("?q=filename:.platform.app.yaml+repo:platformsh-templates/%s", template.Name), "GITHUBTOKEN")
		json.Unmarshal(body, &searchResults)

		// Parse .platform.app.yaml files in search results individually.
		runtimes := []string{}
		for _, file := range searchResults.Items {
			// Get the data.
			var encodedContent structs.TemplateInfo
			body := req.RequestDataAuth(p.URL, fmt.Sprintf("/%s/files/%s?ref=master", template.Name, file.Path), "GITHUBTOKEN")
			json.Unmarshal(body, &encodedContent)
			var appYAML structs.PlatformAppYAML
			decodedContent, err := base64.StdEncoding.DecodeString(encodedContent.Content)
			if err != nil {
				log.Fatal(err)
			}
			yaml.Unmarshal(decodedContent, &appYAML)

			// Use Registry to get language name from runtime type.
			runtime := strings.Split(appYAML.Type, ":")
			runtimes = append(runtimes, getNameFromRuntimeType(runtime[0]))

		}
		// Apply some formatting (i.e. [Node.js, PHP] -> "Node.js/PHP")
		return formatMultiAppString(runtimes)
	}
}

// Returns language name from template's runtime type using Registry.
func getNameFromRuntimeType(runtime string) string {
	// Get registry data.
	url := "https://docs.platform.sh/registry/images/"
	body := req.RequestData(url, "registry.json")
	dynamic := make(map[string]interface{})
	json.Unmarshal(body, &dynamic)

	// Return the language name.
	return dynamic[runtime].(map[string]interface{})["name"].(string)
}

// Applies formatting to a multi-app template slice to return single string.
func formatMultiAppString(runtimesWDuplicates []string) string {

	// Remove duplicates, such as in cases of a multi-app project with two PHP apps.
	runtimes := removeDuplicateRuntimeFromSlice(runtimesWDuplicates)

	// Catch cases where a single app's `.platform.app.yaml` file is not in root.
	if len(runtimes) == 1 {
		return runtimes[0]
	} else {
		// Join runtime strings for multi-app slice.
		return strings.Join(runtimes, "/")
	}
}

// Returns single runtime string for multi-apps that use 2+ same runtime (i.e. two PHP apps).
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
