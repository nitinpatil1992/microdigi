package main

import (
	"net/http"
	"os"
)

var (
	API_PORT        string = ":" + getenv("API_PORT", "80")
	REVERSE_API_URI string = getenv("REVERSE_API_URI", "http://127.0.0.1:8001/reverse")
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	http.HandleFunc("/api", HandleRandom)
	http.ListenAndServe(API_PORT, nil)
}
