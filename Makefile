.PHONY: build infra-up api worker lambdas test test-domain stop

build:
	sam build

infra-up:
	docker compose up -d && sleep 5 && bash scripts/setup-localstack.sh

api:
	go run cmd/api/main.go

worker:
	go run cmd/worker/main.go

lambdas:
	sam local start-lambda --port 3001

test:
	go test ./... -v -count=1

test-domain:
	go test ./domain/... -v

stop:
	docker compose down
