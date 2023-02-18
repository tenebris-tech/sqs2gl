//
// Copyright (c) 2020-2023 Tenebris Technologies Inc.
// All rights reserved
//

package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var AWSID = ""
var AWSKey = ""
var AWSRegion = ""
var AWSQueueName = ""
var Graylog = ""
var Transport = "UDP"
var LogFile = ""

func Load(filename string) error {
	var item []string
	var name string
	var value string

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	//noinspection GoUnhandledErrorResult
	defer file.Close()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// Read line and increment line counter
		line := scanner.Text()
		lineCount++

		// Ignore empty lines and comments
		if len(line) < 1 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "/") {
			continue
		}

		// Split into name value pair
		item = strings.Split(line, "=")
		if len(item) < 2 {
			continue
		}

		name = strings.TrimSpace(strings.ToLower(item[0]))
		value = strings.TrimSpace(item[1])

		switch name {
		case "logfile":
			LogFile = value
		case "awsid":
			AWSID = value
		case "awskey":
			AWSKey = value
		case "awsregion":
			AWSRegion = value
		case "awsqueuename":
			AWSQueueName = value
		case "graylog":
			Graylog = value
		case "transport":
			Transport = value
		default:
			return errors.New(fmt.Sprintf("error parsing config file: %s", line))
		}
	}
	return nil
}
