//
// Copyright (c) 2020-2023 Tenebris Technologies Inc.
// All rights reserved
//

package queue

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"sqs2gl/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var q *sqs.Client
var qURL = ""

func Open() error {
	//var awsCredentials *credentials.Credentials
	var cfg aws.Config
	var err error

	// Initialize a session
	if config.AWSID == "role" {
		// Assume EC2 instance with permissions granted through IAM
		//awsConfig = &aws.Config{
		//Region: aws.String(config.AWSRegion),
		//}
		cfg, err = awsConfig.LoadDefaultConfig(context.TODO(),
			awsConfig.WithRegion(config.AWSRegion))

	} else {
		// Use credentials from configuration
		cfg, err = awsConfig.LoadDefaultConfig(context.TODO(),
			awsConfig.WithRegion(config.AWSRegion),
			awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AWSID, config.AWSKey, "")),
		)
		//awsCredentials = credentials.NewStaticCredentials(config.AWSID, config.AWSKey, "")
		//awsConfig = &aws.Config{
		//	Region:      aws.String(config.AWSRegion),
		//	Credentials: awsCredentials,
		//}
	}

	if err != nil {
		return fmt.Errorf("error loading AWS config: %s", err.Error())
	}

	//awsSession := session.Must(session.NewSession(awsConfig))
	q = sqs.NewFromConfig(cfg)
	if q == nil {
		return errors.New("unable to create new AWS Session")
	}

	// Create a new request to list queues
	listQueuesRequest := sqs.ListQueuesInput{}
	listQueueResults, err := q.ListQueues(context.TODO(), &listQueuesRequest)
	if err != nil {
		return errors.New(fmt.Sprintln("error listing SQS queues: ", err.Error()))
	}

	// Search for requested queue name
	for _, t := range listQueueResults.QueueUrls {
		if strings.Contains(t, config.AWSQueueName) {
			qURL = t
			break
		}
	}

	if qURL == "" {
		return errors.New(fmt.Sprintf("unable to find SQS queue %s", config.AWSQueueName))
	}

	return nil
}
