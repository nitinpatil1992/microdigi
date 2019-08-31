package main

import (
	"fmt"
	"net/http"
	"os"
)

var (
	addr string = ":" + getenv("API_PORT", "8000")
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	http.HandleFunc("/reverse", HandleReverse)
	fmt.Println("starting server on port", addr)
	http.ListenAndServe(addr, nil)
}
