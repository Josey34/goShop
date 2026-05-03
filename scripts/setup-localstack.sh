#!/bin/bash
set -e

ENDPOINT=http://localhost:4566
REGION=ap-southeast-1
ROLE_ARN="arn:aws:iam::000000000000:role/LocalStackRole"
BUILD_DIR=".aws-sam/build"
TMPDIR="$(pwd -W)/.lambda-zips"
mkdir -p "$TMPDIR"

echo "Creating S3 bucket..."
aws --endpoint-url=$ENDPOINT s3 mb s3://goshop-images --region $REGION || true

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

echo "Creating Step Functions state machine..."
aws --endpoint-url=$ENDPOINT stepfunctions create-state-machine \
  --name "OrderWorkflow" \
  --definition "$(cat order-workflow.asl.json)" \
  --role-arn "$ROLE_ARN" \
  --region $REGION

echo "Deploying Lambda functions..."
deploy_lambda() {
  fname=$1
  binary="$BUILD_DIR/$fname/bootstrap"
  zipfile="$TMPDIR/${fname}.zip"

  echo "  Packaging $fname..."
  python -c "import zipfile; z=zipfile.ZipFile('$zipfile','w',zipfile.ZIP_DEFLATED); z.write('$binary','bootstrap'); z.close()"

  echo "  Deploying $fname..."
  aws --endpoint-url=$ENDPOINT lambda create-function \
    --function-name "$fname" \
    --runtime provided.al2023 \
    --handler bootstrap \
    --role "$ROLE_ARN" \
    --zip-file "fileb://$zipfile" \
    --region $REGION
}

deploy_lambda ValidateOrderFunction
deploy_lambda CalculateTotalFunction
deploy_lambda ProcessPaymentFunction
deploy_lambda FulfillOrderFunction
deploy_lambda SendNotificationFunction

echo "Done."
