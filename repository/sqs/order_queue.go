package sqs

import (
	"context"
	"encoding/json"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type OrderQueue struct {
	client   *sqs.Client
	queueURL string
}

func NewOrderQueue(client *sqs.Client, queueURL string) *OrderQueue {
	return &OrderQueue{client: client, queueURL: queueURL}
}

func (q *OrderQueue) Enqueue(ctx context.Context, order *entity.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msgBody := string(body)

	_, err = q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &q.queueURL,
		MessageBody: &msgBody,
	})
	return err
}
