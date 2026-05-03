package order

import (
	"context"
	"encoding/json"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type StartOrderWorkflowInput struct {
	Order                        *entity.Order
	ValidateOrderArn             string
	CalculateTotalArn            string
	ProcessPaymentArn            string
	FulfillOrderArn              string
	SendNotificationFunctionArn  string
	StateMachineArn              string
}

type StartOrderWorkflowOutput struct {
	ExecutionArn string
}

type StartOrderWorkflow struct {
	stepFunctionsClient *sfn.Client
}

func NewStartOrderWorkflow(stepFunctionsClient *sfn.Client) *StartOrderWorkflow {
	return &StartOrderWorkflow{
		stepFunctionsClient: stepFunctionsClient,
	}
}

func (u *StartOrderWorkflow) Execute(ctx context.Context, input *StartOrderWorkflowInput) (*StartOrderWorkflowOutput, error) {
	executionInput := map[string]interface{}{
		"order": map[string]interface{}{
			"ID":         input.Order.ID,
			"CustomerID": input.Order.CustomerID,
			"Items":      input.Order.Items,
			"Status":     input.Order.Status,
		},
		"ValidateOrderFunctionArn":   input.ValidateOrderArn,
		"CalculateTotalFunctionArn":  input.CalculateTotalArn,
		"ProcessPaymentFunctionArn":  input.ProcessPaymentArn,
		"FulfillOrderFunctionArn":    input.FulfillOrderArn,
		"SendNotificationFunctionArn": input.SendNotificationFunctionArn,
	}

	inputJSON, err := json.Marshal(executionInput)
	if err != nil {
		return nil, err
	}

	inputStr := string(inputJSON)
	result, err := u.stepFunctionsClient.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: &input.StateMachineArn,
		Input:           &inputStr,
		Name:            &input.Order.ID,
	})
	if err != nil {
		return nil, err
	}

	return &StartOrderWorkflowOutput{
		ExecutionArn: *result.ExecutionArn,
	}, nil
}
