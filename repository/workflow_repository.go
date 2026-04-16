package repository

import "context"

type WorkflowOrchestrator interface {
	StartOrderWorkflow(ctx context.Context, orderID string) (string, error) // returns execution ARN
	GetWorkflowStatus(ctx context.Context, executionARN string) (string, error)
}
