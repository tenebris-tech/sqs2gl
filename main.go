//
// Copyright (c) 2020-2024 Tenebris Technologies Inc.
// All rights reserved
//

package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"sqs2gl/config"
	"sqs2gl/gelf"
	"sqs2gl/queue"
)

const ProductName = "sqs2gl"
const ProductVersion = "0.4.0"

type cb func(msg string) bool

func main() {

	// Default config file name
	var configFile = "sqs2gl.conf"

	// Callback function to send message
	var callback cb = nil

	// Check for path to config file as argument
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	// Setup signal catching
	signals := make(chan os.Signal, 1)

	// Catch signals
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// method invoked upon seeing signal
	go func() {
		for {
			s := <-signals
			appCleanup(s)
		}
	}()

	// Load configuration information
	err := config.Load(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set up logging
	if config.LogFile != "" {
		f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		// If unable to open log file, report error, but continue writing logs to stderr
		if err != nil {
			log.Printf("Error opening log file: %s", err.Error())
		} else {
			defer func(f *os.File) {
				_ = f.Close()
			}(f)
			log.SetOutput(f)
		}
	}

	log.Printf("Starting %s %s", ProductName, ProductVersion)

	// Select message transmission callback based on config.Transport
	switch strings.ToUpper(config.Transport) {
	case "UDP":
		callback = gelf.UDP
	case "HTTP":
		callback = gelf.HTTP
	default:
		log.Fatal("Invalid transport specified")
	}

	// Initialize queue
	// Loop in case of error
	for {
		err = queue.Open()
		if err != nil {
			log.Printf("Error opening queue: %s", err.Error())
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		} else {
			break
		}
	}

	log.Print("Queue opened, starting receive loop")

	// Start the reception loop, passing the callback function
	err = queue.Loop(callback)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Graceful exit
func appCleanup(sig os.Signal) {
	log.Printf("Exiting on signal: %v", sig)
	os.Exit(0)
}
