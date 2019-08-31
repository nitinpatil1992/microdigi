package main

import (
	"net/http"
	"os"
)

func main() {
	apiHost := os.Getenv("API_HOST")
	http.HandleFunc("/reverse", HandleReverse)
	http.ListenAndServe(apiHost, nil)
}
