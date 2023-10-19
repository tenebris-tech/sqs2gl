//
// Copyright (c) 2020-2023 Tenebris Technologies Inc.
// All rights reserved
//

package gelf

import (
	"fmt"
	"log"
	"net"

	"sqs2gl/config"
)

var conn net.Conn

// UDP send to Graylog and return bool for success or failure
func UDP(msg string) bool {
	var err error

	// If we don't already have a connection, open one
	if conn == nil {
		conn, err = net.Dial("udp", config.Graylog)
		if err != nil {
			log.Printf("unable to open UDP socket: %v", err.Error())
			return false
		}
	}

	_, err = fmt.Fprintf(conn, msg)
	if err != nil {
		return false
	}

	return true
}
