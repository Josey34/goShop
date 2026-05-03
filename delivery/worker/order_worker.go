package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Josey34/goshop/domain/entity"
	orderUC "github.com/Josey34/goshop/usecase/order"
)

func NewOrderMessageHandler(
	workflow *orderUC.StartOrderWorkflow,
	validateOrderArn string,
	calculateTotalArn string,
	processPaymentArn string,
	fulfillOrderArn string,
	sendNotificationFunctionArn string,
	stateMachineArn string,
) MessageHandler {
	return func(ctx context.Context, body string) error {
		var order entity.Order
		if err := json.Unmarshal([]byte(body), &order); err != nil {
			return err
		}

		log.Printf("[worker] processing order %s for customer %s", order.ID, order.CustomerID)

		o := &orderUC.StartOrderWorkflowInput{
			Order:                       &order,
			ValidateOrderArn:            validateOrderArn,
			CalculateTotalArn:           calculateTotalArn,
			ProcessPaymentArn:           processPaymentArn,
			FulfillOrderArn:             fulfillOrderArn,
			SendNotificationFunctionArn: sendNotificationFunctionArn,
			StateMachineArn:             stateMachineArn,
		}

		_, err := workflow.Execute(ctx, o)
		if err != nil {
			log.Printf("[worker] failed to start workflow for order %s: %v", order.ID, err)
		}

		return err
	}
}
