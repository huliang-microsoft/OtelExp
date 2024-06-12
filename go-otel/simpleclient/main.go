package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)  
  
func main() {  
    // Define the request body  
    requestBody := []byte("This is the request body")  
  
    // Create a new HTTP client  
    client := &http.Client{}  
  
    // Create a new HTTP request  
    request, err := http.NewRequest("POST", "http://localhost:12345/annotate", bytes.NewBuffer(requestBody))  
    if err != nil {  
        log.Fatal(err)  
    }  
  
    // Set the content type of the request body  
    request.Header.Set("Content-Type", "text/plain")  
  
    // Send the HTTP request  
    response, err := client.Do(request)  
    if err != nil {  
        log.Fatal(err)  
    }  
  
    // Print the response status code  
    fmt.Println("Response status code:", response.StatusCode)

	respBytes, err := io.ReadAll(nil)
	if err != nil {
		println(err.Error())
	}

	println(len(respBytes))
}  