package parse

import (
  "github.com/chadwcarlson/gomeili/config"
  docs "github.com/chadwcarlson/gomeili/index/documents"
  "github.com/chadwcarlson/gomeili/parse/discourse"
  "github.com/chadwcarlson/gomeili/parse/openapi"
  "github.com/chadwcarlson/gomeili/parse/templates"
)

// General function that wraps individual resource `Get` methods.
func Parse(p config.Config) docs.Index {

  var documents docs.Index

  // Discourse API resource.
  if p.Type == "discourse" {
    documents = discourse.Get(p)
  }
  // OpenAPI 3.0 specification resource.
  if p.Type == "openapi" {
    documents = openapi.Get(p)
  }
  // Platform.sh Template-Builder repo resource.
  if p.Type == "templates" {
    documents = templates.Get(p)
  }

  return documents
}
