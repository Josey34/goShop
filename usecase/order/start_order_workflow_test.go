package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Josey34/goshop/domain/valueobject"
	ucorder "github.com/Josey34/goshop/usecase/order"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type mockSFNClient struct {
	result *sfn.StartExecutionOutput
	err    error
}

func (m *mockSFNClient) StartExecution(ctx context.Context, params *sfn.StartExecutionInput, optFns ...func(*sfn.Options)) (*sfn.StartExecutionOutput, error) {
	return m.result, m.err
}

func ptr(s string) *string { return &s }

func TestStartOrderWorkflow_Execute(t *testing.T) {
	order := makeOrder("o1", valueobject.OrderStatusPending)
	input := &ucorder.StartOrderWorkflowInput{
		Order:                       order,
		ValidateOrderArn:            "arn:validate",
		CalculateTotalArn:           "arn:calculate",
		ProcessPaymentArn:           "arn:payment",
		FulfillOrderArn:             "arn:fulfill",
		SendNotificationFunctionArn: "arn:notify",
		StateMachineArn:             "arn:statemachine",
	}

	t.Run("starts execution and returns arn", func(t *testing.T) {
		client := &mockSFNClient{
			result: &sfn.StartExecutionOutput{ExecutionArn: ptr("arn:execution:o1")},
		}
		uc := ucorder.NewStartOrderWorkflow(client)

		out, err := uc.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if out.ExecutionArn != "arn:execution:o1" {
			t.Errorf("got arn=%s, want=arn:execution:o1", out.ExecutionArn)
		}
	})

	t.Run("sfn error propagates", func(t *testing.T) {
		client := &mockSFNClient{err: errors.New("sfn unavailable")}
		uc := ucorder.NewStartOrderWorkflow(client)

		_, err := uc.Execute(context.Background(), input)
		if err == nil {
			t.Error("expected error from sfn client")
		}
	})
}
