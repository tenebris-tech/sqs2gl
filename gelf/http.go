//
// Copyright (c) 2020-2023 Tenebris Technologies Inc.
// All rights reserved
//

package gelf

import (
	"bytes"
	"log"
	"net/http"

	"sqs2gl/config"
)

// HTTP sends to Graylog and return bool for success or failure
func HTTP(msg string) bool {

	// Create request
	req, err := http.NewRequest("POST", config.Graylog, bytes.NewBuffer([]byte(msg)))
	if err != nil {
		log.Printf("error creating HTTP request: %v", err.Error())
		return false
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP post failed: %v", err.Error())
		return false
	}
	_ = resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 202 {
		return true
	}
	return false
}
