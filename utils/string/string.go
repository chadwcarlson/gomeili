package string

import (
  "strings"
  "github.com/k3a/html2text"
)

// Returns a plaintext string for a blob of HTML body.
func Clean (current string) string {

  // Remove unicode characters.
  formatted := html2text.HTML2Text(current)
  // Remove specific characters.
  formatted = strings.Replace(formatted, "\n", " ", -1)
  formatted = strings.Replace(formatted, "\r", " ", -1)

  // Final unicode removal (specific to code sections).
  return html2text.HTML2Text(formatted)
}
