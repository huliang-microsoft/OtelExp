package main

import (
	"log"
	"net/http"
)  
  
func main() {  
    // Define the API endpoint  
    http.HandleFunc("/annotate", func(w http.ResponseWriter, r *http.Request) {  
        w.WriteHeader(424)  
        println("Get one request.")
    })  
  
    // Start the HTTP server  
    log.Fatal(http.ListenAndServe(":12345", nil))  
}  