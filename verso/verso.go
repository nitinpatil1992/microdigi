package main

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Message string
}

func HandleReverse(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	decoder := json.NewDecoder(r.Body)

	var request Request
	response := make(map[string]interface{})

	err := decoder.Decode(&request)
	if err != nil {
		response["message"] = "Invalid request data"
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response)
		return
	}

	response["message"] = reverse(request.Message)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func reverse(input string) string {
	if len(input) == 1 || len(input) == 0 {
		return input
	}
	output := []byte(input)
	l, r := 0, len(input)-1
	for l < r {
		output[l], output[r] = output[r], output[l]
		r--
		l++
	}
	return string(output)
}
