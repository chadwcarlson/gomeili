package openapi

import (
    "io"
    "os"
    "net/url"
    "fmt"
    "strings"
    "crypto/sha1"
    "github.com/chadwcarlson/gomeili/config"
    "github.com/chadwcarlson/gomeili/utils/ignore"
    docs "github.com/chadwcarlson/gomeili/index/documents"
    "github.com/getkin/kin-openapi/openapi3"
)

func formatPathURL(path string, current_op *openapi3.Operation, operation string) string {
  if "" == current_op.OperationID {
    tag := getTags(current_op)
    return "#tag/" + tag + "/paths/" + strings.Replace(path, "/", "~1", -1) + "/" + operation
  } else {
    return "#operation/" + current_op.OperationID
  }
}

func getTags(operation *openapi3.Operation) string {
  return strings.Replace(operation.Tags[0], " ", "-", -1)
}

func getDocument(path_string string, p config.Config, operation *openapi3.Operation, opType string) docs.Document {
  var document docs.Document

  document.Site = p.Name
  document.Section = path_string
  document.Title = operation.Summary

  rel_url := formatPathURL(path_string, operation, opType)
  full_url := fmt.Sprintf("%s%s", p.URL, rel_url)
  document.URL = full_url
  document.RelativeURL = fmt.Sprintf("/%s", rel_url)

  h := sha1.New()
  h.Write([]byte(full_url))
  document.DocumentID = fmt.Sprintf("%x", h.Sum(nil))

  document.Rank = p.Rank
  if p.Rank == 1 {
    document.Source = "primary"
  } else {
    document.Source = "secondary"
  }

  document.Text = strings.Replace(operation.Description, "\n", " ",-1)
  document.Description = strings.Replace(operation.Description, "\n", " ",-1)
  document.Subsection = ""

  return document
}

func getTagDocuments(p config.Config, swagger *openapi3.Swagger, allDocuments docs.Index) docs.Index {
  io.WriteString(os.Stdout, fmt.Sprintf("\033[1m %s\033[0m (%v tags)\n", "Tags", len(swagger.Tags)))
  for _, tag := range swagger.Tags {

    var document docs.Document

    document.Site = p.Name
    document.Section = ""
    document.Title = tag.Name

    rel_url := fmt.Sprintf("#tag/%s", strings.Replace(tag.Name, " ", "-",-1))
    full_url := fmt.Sprintf("%s%s", p.URL, rel_url)
    document.URL = full_url
    document.RelativeURL = fmt.Sprintf("/%s", rel_url)

    h := sha1.New()
    h.Write([]byte(full_url))
    document.DocumentID = fmt.Sprintf("%x", h.Sum(nil))

    document.Rank = p.Rank

    document.Text = strings.Replace(tag.Description, "\n", " ",-1)
    document.Description = strings.Replace(tag.Description, "\n", " ",-1)
    document.Subsection = ""

    allDocuments.Documents = append(allDocuments.Documents, document)
  }

  return allDocuments
}

func getPathDocuments(p config.Config, swagger *openapi3.Swagger, allDocuments docs.Index) docs.Index {

  var emptyOp *openapi3.Operation
  io.WriteString(os.Stdout, fmt.Sprintf("\033[1m %s\033[0m (%v paths)\n", "Paths", len(swagger.Paths)))
  for key, path := range swagger.Paths {

    // GET Operation
    if emptyOp != path.Get {
      if !ignore.ItemExists(p.Ignore, getTags(path.Get)) {
        allDocuments.Documents = append(allDocuments.Documents, getDocument(key, p, path.Get, "get"))
      }
    }
    // POST Operation
    if emptyOp != path.Post {
      if !ignore.ItemExists(p.Ignore, getTags(path.Post)) {
        allDocuments.Documents = append(allDocuments.Documents, getDocument(key, p, path.Post, "post"))
      }
    }
    // DELETE Operation
    if emptyOp != path.Delete {
      if !ignore.ItemExists(p.Ignore, getTags(path.Delete)) {
        allDocuments.Documents = append(allDocuments.Documents, getDocument(key, p, path.Delete, "delete"))
      }
    }
    // PATCH Operation
    if emptyOp != path.Patch {
      if !ignore.ItemExists(p.Ignore, getTags(path.Patch)) {
        allDocuments.Documents = append(allDocuments.Documents, getDocument(key, p, path.Patch, "patch"))
      }
    }
  }

  return allDocuments
}

func Get(p config.Config) docs.Index {

  // Load the specification from the given url.
  io.WriteString(os.Stdout, fmt.Sprintf("\n\033[1mOpen API 3.0 specification @\033[0m %s\n", p.URL + p.File))
  loader := openapi3.NewSwaggerLoader()
  loader.IsExternalRefsAllowed = true
  url,  err := url.Parse(p.URL + p.File)
  if err != nil {
    fmt.Print(err)
  }
  swagger, err := loader.LoadSwaggerFromURI(url)
  if err != nil {
    fmt.Print(err)
  }

  var allDocuments docs.Index

  // Get documents from spec's `tags`.
  allDocuments = getTagDocuments(p, swagger, allDocuments)

  // Get documents from spec's `paths`.
  allDocuments = getPathDocuments(p, swagger, allDocuments)

  // TODO: Get the gigantic `description` attribute and parse it.
  // fmt.Print(swagger.Info.Description + "\n")
  // var buf bytes.Buffer
  // if err := goldmark.Convert([]byte(swagger.Info.Description), &buf); err != nil {
  //   panic(err)
  // }

  return allDocuments

}
