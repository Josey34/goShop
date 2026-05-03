package main

import (
	"context"
	"log"

	"github.com/Josey34/goshop/config"
	"github.com/Josey34/goshop/delivery/worker"
	"github.com/Josey34/goshop/usecase/order"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(),
		awsConfig.WithRegion(cfg.AWS.Region),
		awsConfig.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{URL: cfg.SQS.Endpoint}, nil
		})),
	)
	if err != nil {
		log.Fatal(err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)
	sfnClient := sfn.NewFromConfig(awsCfg)
	workflow := order.NewStartOrderWorkflow(sfnClient)
	handler := worker.NewOrderMessageHandler(
		workflow,
		"arn:aws:lambda:ap-southeast-1:000000000000:function:ValidateOrderFunction",
		"arn:aws:lambda:ap-southeast-1:000000000000:function:CalculateTotalFunction",
		"arn:aws:lambda:ap-southeast-1:000000000000:function:ProcessPaymentFunction",
		"arn:aws:lambda:ap-southeast-1:000000000000:function:FulfillOrderFunction",
		"arn:aws:lambda:ap-southeast-1:000000000000:function:SendNotificationFunction",
		"arn:aws:states:ap-southeast-1:000000000000:stateMachine:OrderWorkflow",
	)
	consumer := worker.NewConsumer(sqsClient, cfg.SQS.QueueURL, handler)

	if err := consumer.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
