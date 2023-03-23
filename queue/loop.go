//
// Copyright (c) 2020-2023 Tenebris Technologies Inc.
// All rights reserved
//

package queue

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func Loop(callback func(string) bool) error {
	var err error
	var success bool
	var r *sqs.ReceiveMessageOutput
	var msg string

	// Set receive parameters
	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(qURL),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(900),
		WaitTimeSeconds:     aws.Int64(15),
	}

	// Loop and receive messages
	for {
		r, err = q.ReceiveMessage(receiveParams)
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
				_, err := q.DeleteMessage(deleteParams)
				if err != nil {
					return err
				}
			} else {
				// Delay on failure
				time.Sleep(10 * time.Second)
			}
		}
	}
}
