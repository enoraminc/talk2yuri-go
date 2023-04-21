package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {

	// Upload
	const API_KEY = "YOUR_API_KEY"
	const UPLOAD_URL = "https://api.assemblyai.com/v2/upload"
	const TRANSCRIPT_URL = "https://api.assemblyai.com/v2/transcript"

	// Load file
	data, err := ioutil.ReadFile("VIDIO_ID")
	if err != nil {
		log.Fatalln(err)
	}

	// Setup HTTP client and set header
	client := &http.Client{}
	req, _ := http.NewRequest("POST", UPLOAD_URL, bytes.NewBuffer(data))
	req.Header.Set("authorization", API_KEY)
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	// decode json and store it in a map
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	// print the upload_url
	fmt.Println(result["upload_url"])
	url := fmt.Sprintf("%v", result["upload_url"])

	// Prepare json data
	values := map[string]string{"audio_url": string(url)}
	jsonData, err_json := json.Marshal(values)

	if err_json != nil {
		log.Fatalln(err_json)
	}

	// Setup HTTP client and set header
	clientt := &http.Client{}
	req_set, _ := http.NewRequest("POST", TRANSCRIPT_URL, bytes.NewBuffer(jsonData))
	req_set.Header.Set("content-type", "application/json")
	req_set.Header.Set("authorization", API_KEY)
	ress, err_json := clientt.Do(req_set)

	if err_json != nil {
		log.Fatalln(err_json)
	}

	defer ress.Body.Close()

	// Decode json and store it in a map
	var result_json map[string]interface{}
	json.NewDecoder(ress.Body).Decode(&result_json)

	// Print the id of the transcribed audio
	fmt.Println("Transcription ID:", result_json["id"])

	// Set the polling URL
	POLLING_URL := TRANSCRIPT_URL + "/" + fmt.Sprintf("%v", result_json["id"])

	// Poll for the completed status
	for {
		// Send GET request
		client_pol := &http.Client{}
		req_pol, _ := http.NewRequest("GET", POLLING_URL, nil)
		req_pol.Header.Set("content-type", "application/json")
		req_pol.Header.Set("authorization", API_KEY)
		res_pol, err_pol := client_pol.Do(req_pol)

		if err_pol != nil {
			log.Fatalln(err_pol)
		}

		defer res_pol.Body.Close()

		// Decode json and store it in a map
		var result_transcribe map[string]interface{}
		json.NewDecoder(res_pol.Body).Decode(&result_transcribe)

		// Check status and print the transcribed text
		if result_transcribe["status"] == "completed" {
			fmt.Println(result_transcribe["text"])
			break
		}

		// Wait for 5 seconds before polling again
		time.Sleep(5 * time.Second)
	}
}
