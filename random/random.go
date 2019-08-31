package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ReverseApiChannel struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

type Message1 struct {
	msg string `json:"message"`
}

type Request struct {
	Message string `json:"message"`
}

type ReverseApiResponse struct {
	Message string `json:"message"`
}

func HandleRandom(w http.ResponseWriter, r *http.Request) {
	var request Request
	response := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		log.WithFields(log.Fields{
			"Invalid request data ": err.Error(),
		}).Warn("Error parsing request")
		response["message"] = "Invalid request data"
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response)
		return
	}

	// channel to connect with webservice which returns string reverse

	apiResponseSuccess, apiResponseError := callReverseApi(request.Message)

	if err := <-apiResponseError; err != nil {
		log.WithFields(log.Fields{
			"Error occurred ": apiResponseError,
		}).Error("Error parsing request")
		response["message"] = "Reverse string webservice unavailable"
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	msg := <-apiResponseSuccess
	log.Info("Reverse api responded with ", msg)

	var apiResp ReverseApiResponse
	err = json.Unmarshal([]byte(msg), &apiResp)
	response["message"] = apiResp.Message

	response["rand"] = rand.Float64()
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func callReverseApi(message string) (<-chan string, <-chan error) {

	data := make(map[string]interface{})
	data["message"] = message
	jsonValue, _ := json.Marshal(data)

	out := make(chan string, 1)
	errs := make(chan error, 1)

	go func() {
		req, err := http.NewRequest("POST", REVERSE_API_URI, bytes.NewBuffer(jsonValue))
		client := &http.Client{}
		log.Info("Calling reverse api with data ", bytes.NewBuffer(jsonValue))
		resp, err := client.Do(req)

		if err != nil {
			errs <- err
			close(out)
			close(errs)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errs <- err
			close(out)
			close(errs)
			return
		}
		out <- string(body)
		close(out)
		close(errs)
	}()

	return out, errs
}
