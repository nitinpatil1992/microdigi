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

	// check if cache exists
	cacheExists, _ := redisClient.Exists(request.Message).Result()

	if cacheExists == 1 {
		log.Info("Cache hit occurred")

		response["message"], _ = redisClient.HGet(request.Message, "message").Result()
		response["rand"], _ = redisClient.HGet(request.Message, "rand").Result()

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// channel to connect with webservice which returns string reverse
	apiResponseSuccess, apiResponseError := callReverseApi(request.Message)

	if err := <-apiResponseError; err != nil {
		log.WithFields(log.Fields{
			"Error occurred ": apiResponseError,
		}).Error("Error calling reverse webservice")

		response["message"] = "Reverse string webservice unavailable"

		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(response)
		return
	}

	msg := <-apiResponseSuccess
	log.Info("Reverse api responded with ", msg)

	var apiResp ReverseApiResponse
	err = json.Unmarshal([]byte(msg), &apiResp)
	randomValue := rand.Float64()

	response["message"] = apiResp.Message
	response["rand"] = randomValue

	//set cache
	_ = redisClient.HSet(request.Message, "message", apiResp.Message)
	_ = redisClient.HSet(request.Message, "rand", randomValue)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func callReverseApi(message string) (<-chan string, <-chan error) {

	out := make(chan string, 1)
	errs := make(chan error, 1)

	go func() {
		data := make(map[string]interface{})
		data["message"] = message
		jsonValue, _ := json.Marshal(data)

		req, err := http.NewRequest("POST", REVERSE_API_URI, bytes.NewBuffer(jsonValue))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Warn("request error")

			errs <- err
			close(out)
			close(errs)
			return
		}
		client := &http.Client{}

		log.Info("Calling reverse api with data ", bytes.NewBuffer(jsonValue))
		resp, err := client.Do(req)

		if err != nil {
			log.Warn("Failed to get the response from reverse api")
			errs <- err
			close(out)
			close(errs)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warn("Failed to read response obtained from reverse api")
			errs <- err
			close(out)
			close(errs)
			return
		}
		out <- string(body)
		close(out)
		close(errs)
	}()
	log.Debug("Reverse api call successful")
	return out, errs
}
