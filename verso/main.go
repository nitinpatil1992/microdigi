package main

import (
	"fmt"
	"net/http"
	"os"
)

var (
	apiHost string = getenv("API_HOST", "127.0.0.1:8000")
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	//apiHost := os.Getenv("API_HOST")
	http.HandleFunc("/reverse", HandleReverse)
	fmt.Println("starting server on", apiHost)
	http.ListenAndServe(apiHost, nil)
}
