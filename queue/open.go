//
// Copyright (c) 2020-2023 Tenebris Technologies Inc.
// All rights reserved
//

package queue

import (
	"errors"
	"fmt"
	"strings"

	"sqs2gl/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var q *sqs.SQS
var qURL = ""

func Open() error {
	var awsCredentials *credentials.Credentials
	var awsConfig *aws.Config

	// Initialize a session
	if config.AWSID == "role" {
		// Assume EC2 instance with permissions granted through IAM
		awsConfig = &aws.Config{
			Region: aws.String(config.AWSRegion),
		}
	} else {
		// Use credentials from configuration
		awsCredentials = credentials.NewStaticCredentials(config.AWSID, config.AWSKey, "")
		awsConfig = &aws.Config{
			Region:      aws.String(config.AWSRegion),
			Credentials: awsCredentials,
		}
	}

	awsSession := session.Must(session.NewSession(awsConfig))
	q = sqs.New(awsSession)
	if q == nil {
		return errors.New("unable to create new AWS Session")
	}

	// Create a new request to list queues
	listQueuesRequest := sqs.ListQueuesInput{}
	listQueueResults, err := q.ListQueues(&listQueuesRequest)
	if err != nil {
		return errors.New(fmt.Sprintln("error listing SQS queues: ", err.Error()))
	}

	// Search for requested queue name
	for _, t := range listQueueResults.QueueUrls {
		if strings.Contains(*t, config.AWSQueueName) {
			qURL = *t
			break
		}
	}

	if qURL == "" {
		return errors.New(fmt.Sprintf("unable to find SQS queue %s", config.AWSQueueName))
	}

	return nil
}
