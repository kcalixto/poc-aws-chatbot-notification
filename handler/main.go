package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func main() {
	if os.Getenv("NODE_ENV") == "local" {
		err := Handler(context.Background(), events.APIGatewayProxyRequest{})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		lambda.Start(Handler)
	}
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(endpoints.SaEast1RegionID),
	})
	if err != nil {
		return err
	}

	messageBody, err := json.Marshal(map[string]interface{}{
		"version": "1.0",
		"source":  "custom",
		"content": map[string]interface{}{
			"title":       "title",
			"description": ":description",
			"nextSteps": []string{
				"Refer to <http://www.example.com|*diagnosis* runbook>",
				"@googlie: Page Jane if error persists over 30 minutes",
				"Check if instance i-04d231f25c18592ea needs to receive an AMI rehydration",
			},
		},
	})
	if err != nil {
		return err
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(messageBody)),
		TopicArn: aws.String(os.Getenv("CHATBOT_SNS_TOPIC_ARN")),
	}

	svc := sns.New(sess)
	result, err := svc.Publish(input)
	if err != nil {
		return err
	}

	fmt.Println("SNS message published with MessageId:", *result.MessageId)
	return nil
}
