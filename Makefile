#! /usr/bin/make

APP_NAME    = password-self-service
VERSION     = $(shell git rev-parse --abbrev-ref HEAD)
REVISION    = $(shell git rev-parse HEAD)
BUILD_DATE  = $(shell date -I'seconds')

VERSION_LDFLAGS := \
			-X main.Version=$(VERSION) \
			-X main.Revision=$(REVISION) \
			-X main.BuildDate=$(BUILD_DATE)

build:
	@echo ">>> building code"
	GO111MODULE=on CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=mod -ldflags "$(VERSION_LDFLAGS)" -o $(APP_NAME)

init:
	@echo ">>> generate swag"
	rm -rf docs
	${GOPATH}/bin/swag init
	go mod tidy

run:
	@echo ">>> $(APP_NAME) are running"
	go run main.go --config etc/$(APP_NAME).yaml

clean:
	@echo ">>> clean cache"
	rm -rf $(APP_NAME)
	go clean -i .
	rm -rf $(APP_NAME) go.mod go.sum docs vendor $(APP_NAME).pid

docker-build:
	@echo ">>> build docker image"
	docker build -t kubeop/$(APP_NAME):$(VERSION) .

docker-push:
	@echo ">>> push docker image"
	docker push kubeop/$(APP_NAME):$(VERSION)
