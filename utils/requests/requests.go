package requests

import (
	"fmt"
	"github.com/alexsasharegan/dotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// This function performs a generic request from a root url and sub-path,
// returning its body to be Unmarshaled.
func RequestData(url_root string, path string) []byte {

	// Make the request.
	resp, err := http.Get(url_root + path)
	if err != nil {
		print(err)
	}

	// Return the request body.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	return body

}

// This function performs an authenticated request from a root url and sub-path,
// returning its body to be Unmarshaled.
func RequestDataAuth(url_root string, path string, authVar string) []byte {

	// Use dotenv when locally testing.
	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Define the request.
	req, err := http.NewRequest("GET", url_root+path, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	// Set the authorization header.
	req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv(authVar)))

	// Define the client with timeout set.
	client := &http.Client{Timeout: time.Second * 10}

	// Make the request.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}

	// Return the request body.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	return body

}
