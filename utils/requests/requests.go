package requests

import (
  "io/ioutil"
  "net/http"
)

// This function performs a generic request from a root url and sub-path,
// returning its body to be Unmarshaled.
func RequestData(url_root string, path string) []byte {
  resp, err := http.Get(url_root + path)
  if err != nil {
      print(err)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      print(err)
  }
  return body

}
