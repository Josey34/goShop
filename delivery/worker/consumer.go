package worker

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type MessageHandler func(ctx context.Context, body string) error

type Consumer struct {
	client   *sqs.Client
	queueURL string
	handler  MessageHandler
}

func NewConsumer(client *sqs.Client, queueURL string, handler MessageHandler) *Consumer {
	return &Consumer{
		client:   client,
		queueURL: queueURL,
		handler:  handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		result, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &c.queueURL,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20,
		})
		if err != nil {
			return err
		}

		for _, msg := range result.Messages {
			if err := c.handler(ctx, *msg.Body); err != nil {
				log.Printf("error handling message: %v", err)
				continue
			}

			c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &c.queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})
		}
	}
}
