#!/bin/bash
ENDPOINT=http://localhost:4566
REGION=ap-southeast-1

echo "Creating S3 bucket..."
aws --endpoint-url=$ENDPOINT s3 mb s3://goshop-images --region $REGION

echo "Creating DLQ..."
aws --endpoint-url=$ENDPOINT sqs create-queue \
  --queue-name goshop-orders-dlq \
  --region $REGION

echo "Creating main queue..."
DLQ_ARN=$(aws --endpoint-url=$ENDPOINT sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/goshop-orders-dlq \
  --attribute-names QueueArn \
  --query 'Attributes.QueueArn' \
  --output text --region $REGION)

aws --endpoint-url=$ENDPOINT sqs create-queue \
  --queue-name goshop-orders \
  --attributes "{\"RedrivePolicy\":\"{\\\"deadLetterTargetArn\\\":\\\"$DLQ_ARN\\\",\\\"maxReceiveCount\\\":\\\"3\\\"}\"}" \
  --region $REGION
  
aws --endpoint-url=$ENDPOINT stepfunctions create-state-machine \
  --name "OrderWorkflow" \
  --definition "$(cat order-workflow.asl.json)" \
  --role-arn "arn:aws:iam::000000000000:role/LocalStackRole" \
  --region $REGION

echo "Done."
