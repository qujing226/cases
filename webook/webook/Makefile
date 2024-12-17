.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o ./bin/webook
	@docker rmi peninsula/webook:v0.0.1 || true
	@docker build -t peninsula/webook:v0.0.1 .

.PHONY: mock
mock:
	@mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/ mocks/user.mock.go
	@mockgen -source=./internal/service/code.go -package=svcmocks -destination=./internal/service/mocks/code.mock.go
	@mockgen -source=./internal/repository/user.go -package=repomocks -destination=./internal/repository/mocks/user.mock.go
	@mockgen -source=./internal/repository/code.go -package=repomocks -destination=./internal/repository/mocks/code.mock.go

	@go mod tidy