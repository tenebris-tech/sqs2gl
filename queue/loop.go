//
// Copyright (c) 2020-2024 Tenebris Technologies Inc.
// All rights reserved
//

package queue

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func Loop(callback func(string) bool) error {
	var success bool
	var msg string

	// Set receive parameters
	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            &qURL,
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   900,
		WaitTimeSeconds:     15,
	}

	// Loop and receive messages
	for {
		r, err := q.ReceiveMessage(context.TODO(), receiveParams)
		if err != nil {
			return err
		}

		if len(r.Messages) > 0 {
			// Retrieve message body
			// Limit is set to 1 so iteration is not required
			msg = *r.Messages[0].Body

			// Pass to the callback function
			success = callback(msg)

			// If success, delete message from queue
			if success {
				deleteParams := &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(qURL),
					ReceiptHandle: r.Messages[0].ReceiptHandle,
				}
				_, err := q.DeleteMessage(context.TODO(), deleteParams)
				if err != nil {
					return err
				}
			} else {
				// Delay on failure
				time.Sleep(15 * time.Second)
			}
		}
	}
}
