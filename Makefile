.PHONY: test fmt lint showcover

showcover:
	@go tool cover -html=cp.out
	
lint:
	@make fmt
	@golangci-lint --skip-dirs cdk.out run

test:
	@make lint && go test -coverprofile cp.out $$(go list ./... | grep -v /cdk.out/)

dockerbuild:
	@docker build -f Dockerfile.etl -t monty .

dockerpush:
	@aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 261392311630.dkr.ecr.us-east-2.amazonaws.com
	@docker tag monty:latest 261392311630.dkr.ecr.us-east-2.amazonaws.com/monty:latest
	@docker push 261392311630.dkr.ecr.us-east-2.amazonaws.com/monty:latest

dockerdeploy: dockerbuild dockerpush

dockerrun:
	@docker run -e APP_ID -e APP_SECRET -e APP_AGENT -e DB_HOST -e DB_PASSWORD -e DB_USER -e DB_PORT -e DB_NAME monty
			
fmt:
	@gofmt -s -w .
