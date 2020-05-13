.PHONY: proto run-dev build push-image build-image run-image

NAME := blockchain-go-pool
TAG := $$( git rev-parse --short HEAD )
IMAGE := ${NAME}\:${TAG}

proto:
	@cd proto && protoc --go_out=plugins=grpc:. *.proto

run-dev:
	@go run . dev

build:
	@GOOS=linux GOARCH=amd64 go build -o bin/app .

build-image:
	@make build
	@docker build -f ./Dockerfile -t ${IMAGE} .

push-image:
	@docker push ${IMAGE}

run-image:
	@docker run --name article -p 9000:9000 ${IMAGE}


