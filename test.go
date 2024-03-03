package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	// URL for server
	baseURL := "http://localhost:8080/auth"

	// Test for valid JWT
	fmt.Println("Requesting valid JWT...")
	testAuthRequest(baseURL, false)

	// Test for expired JWT
	fmt.Println("Requesting expired JWT...")
	testAuthRequest(baseURL, true)
}

func testAuthRequest(baseURL string, expired bool) {
	//  request URL with the "expired" query parameter 
	requestURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing base URL:", err)
		return
	}
	query := requestURL.Query()
	if expired {
		query.Set("expired", "true")
	}
	requestURL.RawQuery = query.Encode()

	//  new HTTP GET request
	req, err := http.NewRequest("POST", requestURL.String(), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// print the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Printf("Response for expired=%t: %s\n\n", expired, body)
}