SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   image				Build the docker image"
	@echo "	   build-multi			Build and push multi-platform docker image"
	@echo "	   build				Build the voter executable"
	@echo "	   run					Run the voter program from code"
	@echo "	   run-binary			Run the voter program from binary"
	@echo "	   run-container		Run the voter program from container"
	@echo "	   test					Run all tests"


.PHONY: build-multi
build-multi:
	docker buildx build --platform linux/amd64,linux/arm64 -t agentjsmith/voter-container . --push

.PHONY: image
image:
	docker buildx build -t agentjsmith/voter-container . --load

.PHONY: build
build:
	go build -o voter-api .

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test -v ./... -count=1

.PHONY: compose-up
compose-up: image
	docker compose up -d

.PHONY: compose-down
compose-down:
	docker compose down

.PHONY: run-binary
run-binary:
	./voter-api .

.PHONY: run-container
run-container: image
	docker run jsmith5444/voter-container