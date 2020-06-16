package openapi

import (
	"crypto/sha1"
	"fmt"
	"github.com/chadwcarlson/gomeili/config"
	docs "github.com/chadwcarlson/gomeili/index/documents"
	"github.com/chadwcarlson/gomeili/utils/ignore"
	"github.com/getkin/kin-openapi/openapi3"
	"io"
	"net/url"
	"os"
	"strings"
)

// Main method to retrieve OpenAPI 3.0 index for Meilisearch.
// 		TODO: Get the gigantic `description` attribute and parse it (golmark?)(swagger.Info.Description).
func Get(p config.Config) docs.Index {

	var allDocuments docs.Index

	// Load the specification from the given configuration.
	swagger := loadSpecification(p)

	// Get documents from spec's `paths`.
	allDocuments = getPathDocuments(p, swagger, allDocuments)

	// Get documents from spec's `tags`.
	allDocuments = getTagDocuments(p, swagger, allDocuments)

	return allDocuments
}

// Loads OpenAPI specification from the given `config.yaml`.
func loadSpecification(p config.Config) *openapi3.Swagger {
	io.WriteString(os.Stdout, fmt.Sprintf("\n\033[1mOpen API 3.0 specification @\033[0m %s\n", p.URL+p.File))
	loader := openapi3.NewSwaggerLoader()

	var emptyConfig config.Config
	// Load from a remote URL.
	if emptyConfig.URL != p.URL {
		loader.IsExternalRefsAllowed = true
		url, err := url.Parse(p.URL + p.File)
		if err != nil {
			fmt.Print(err)
		}
		swagger, err := loader.LoadSwaggerFromURI(url)
		if err != nil {
			fmt.Print(err)
		}
		return swagger
	// Load from a local file (as will be the case for the API Documentation).
	} else {
		swagger, err := loader.LoadSwaggerFromFile(p.File)
		if err != nil {
			fmt.Print(err)
		}
		return swagger
	}
	// return swagger
}

// Make the documents for each of the specification's "Paths".
func getPathDocuments(p config.Config, swagger *openapi3.Swagger, allDocuments docs.Index) docs.Index {

	// Create an empty Operation we can use to judge if field is undefined.
	var emptyOp *openapi3.Operation
	io.WriteString(os.Stdout, fmt.Sprintf("\033[1m %s\033[0m (%v paths)\n", "Paths", len(swagger.Paths)))
	for endpoint, path := range swagger.Paths {

		// GET Operation.
		if emptyOp != path.Get {
			// Make sure that the current path does not have a tag we'd like to ignore.
			if !ignore.ItemExists(p.Ignore, getTags(path.Get)) {
				// Append the document.
				allDocuments.Documents = append(allDocuments.Documents, getDocument(endpoint, p, path.Get, "get"))
			}
		}
		// POST Operation.
		if emptyOp != path.Post {
			// Make sure that the current path does not have a tag we'd like to ignore.
			if !ignore.ItemExists(p.Ignore, getTags(path.Post)) {
				// Append the document.
				allDocuments.Documents = append(allDocuments.Documents, getDocument(endpoint, p, path.Post, "post"))
			}
		}
		// DELETE Operation.
		if emptyOp != path.Delete {
			// Make sure that the current path does not have a tag we'd like to ignore.
			if !ignore.ItemExists(p.Ignore, getTags(path.Delete)) {
				// Append the document.
				allDocuments.Documents = append(allDocuments.Documents, getDocument(endpoint, p, path.Delete, "delete"))
			}
		}
		// PATCH Operation.
		if emptyOp != path.Patch {
			// Make sure that the current path does not have a tag we'd like to ignore.
			if !ignore.ItemExists(p.Ignore, getTags(path.Patch)) {
				// Append the document.
				allDocuments.Documents = append(allDocuments.Documents, getDocument(endpoint, p, path.Patch, "patch"))
			}
		}
	}
	return allDocuments
}

// Make the documents for each of the specification's "Tags".
func getTagDocuments(p config.Config, swagger *openapi3.Swagger, allDocuments docs.Index) docs.Index {
	io.WriteString(os.Stdout, fmt.Sprintf("\033[1m %s\033[0m (%v tags)\n", "Tags", len(swagger.Tags)))
	for _, tag := range swagger.Tags {

		var document docs.Document

		// Basics.
		document.Site = p.Name
		document.Title = tag.Name
		document.Section = ""
		document.Subsection = ""

		// URLs.
		rel_url := fmt.Sprintf("#tag/%s", strings.Replace(tag.Name, " ", "-", -1))
		full_url := fmt.Sprintf("%s%s", p.URL, rel_url)
		document.URL = full_url
		document.RelativeURL = fmt.Sprintf("/%s", rel_url)

		// DocumentID hash.
		h := sha1.New()
		h.Write([]byte(full_url))
		document.DocumentID = fmt.Sprintf("%x", h.Sum(nil))

		// Match `config.yaml` rank, and use React primary/secondary designation.
		document.Rank = p.Rank
		if p.Rank == 1 {
			document.Source = "primary"
		} else {
			document.Source = "secondary"
		}

		// Document body text.
		document.Text = strings.Replace(tag.Description, "\n", " ", -1)

		// Document description.
		document.Description = strings.Replace(tag.Description, "\n", " ", -1)

		// Append the document.
		allDocuments.Documents = append(allDocuments.Documents, document)
	}

	return allDocuments
}

// Retrieve the tags assigned to a Path.
func getTags(operation *openapi3.Operation) string {
	return strings.Replace(operation.Tags[0], " ", "-", -1)
}

// Builds the document for an individual Path operation type.
func getDocument(endpoint string, p config.Config, operation *openapi3.Operation, opType string) docs.Document {
	var document docs.Document

	// Basic.
	document.Site = p.Name
	document.Title = operation.Summary
	document.Section = endpoint
	document.Subsection = endpoint

	// URLs.
	rel_url := formatPathURL(endpoint, operation, opType)
	full_url := fmt.Sprintf("%s%s", p.URL, rel_url)
	document.URL = full_url
	document.RelativeURL = fmt.Sprintf("/%s", rel_url)

	// DocumentID hash.
	h := sha1.New()
	h.Write([]byte(full_url))
	document.DocumentID = fmt.Sprintf("%x", h.Sum(nil))

	// Match `config.yaml` rank, and use React primary/secondary designation.
	document.Rank = p.Rank
	if p.Rank == 1 {
		document.Source = "primary"
	} else {
		document.Source = "secondary"
	}

	// Document text.
	document.Text = strings.Replace(operation.Description, "\n", " ", -1)

	// Document description.
	document.Description = strings.Replace(operation.Description, "\n", " ", -1)

	return document
}

// Match Redoc's URL styling for the Path.
func formatPathURL(path string, current_op *openapi3.Operation, operation string) string {
	// If the optional OperationID is not present, the header link is built off of the path.
	if "" == current_op.OperationID {
		tag := getTags(current_op)
		return "#tag/" + tag + "/paths/" + strings.Replace(path, "/", "~1", -1) + "/" + operation
	// Otherwise, the OperationID is used.
	} else {
		return "#operation/" + current_op.OperationID
	}
}
