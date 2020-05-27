package requests

import (
  // "os"
  "io/ioutil"
  "net/http"
)

// func remoteRequest(url_root string, path string) []byte {
//
// }

// func LocalRequest(path string, filename string) {
//
// 	jsonFile, err := os.Open(file)
//
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer jsonFile.Close()
//
// 	byteValue, _ := ioutil.ReadAll(jsonFile)
//
// 	var result psh.EnvList
// 	json.Unmarshal([]byte(byteValue), &result)
//
// 	return result
//
//
// }

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
