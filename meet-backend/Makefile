ifeq ($(OS),Windows_NT)
    CURRENT_DIR=$(CD)
	IMAGE_NAME := $(shell basename "$(CD)")
	SSH_PRIVATE_KEY="$$(type ~/.ssh/id_rsa)"
else
	CURRENT_DIR=$(PWD)
	IMAGE_NAME := $(shell basename "$(PWD)")
	SSH_PRIVATE_KEY="$$(cat ~/.ssh/id_rsa)"
endif

install:
	go env -w GOPRIVATE=github.com/erodriguezg
	go mod download

updatedeps:
	@echo Updating dependencies
	go get -d -v -u ./...
	go mod tidy

check:
	@echo Analyzing suspicious constructs
	go vet ./...

escape-analysis:
	@echo Analyzing the dynamic scope of pointers
	go build -gcflags='-m -l' -o bin/api cmd/server/api.go

build:
	go build -o cmd/server/api.go

compile:
	GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	go build \
	-ldflags="-w -s" \
	-o bin/api ./cmd/server/api.go

run:
	go run cmd/server/api.go

test:
	go test ./... --coverprofile=coverage.out

test-report: test
	go tool cover -html=coverage.out

swagger-gen:
	go install github.com/swaggo/swag/cmd/swag@v1.8.7
	$(shell go env GOPATH)/bin/swag init -g ./cmd/server/api.go -o ./docs

linter:
	docker run --rm \
	-v ${CURRENT_DIR}:/app \
	-v ~/.ssh:/root/.ssh:ro \
	-w /app \
	golangci/golangci-lint:latest \
	/bin/bash -c "git config --global url."git@github.com:".insteadOf "https://github.com/" && golangci-lint run -v --timeout 3m" 

security:
	docker run --rm \
	-v ${CURRENT_DIR}:/app \
	-v ~/.ssh:/root/.ssh:ro \
	-w /app \
	--entrypoint="/bin/bash" \
	securego/gosec:latest -c "apk add -U --no-cache openssh && git config --global url."git@github.com:".insteadOf "https://github.com/" && gosec -exclude=G101 ./..."

test-all: compile test check security linter

mock-generate-all-aux: 
	docker run --rm -it \
	-v ${CURRENT_DIR}:/src \
	-v ~/.ssh:/root/.ssh:ro \
	-w /src \
	--entrypoint="/bin/sh" \
	vektra/mockery -c "apk add -U --no-cache openssh git && git config --global url."git@github.com:".insteadOf "https://github.com/" && go mod download && mockery --output ./aux/mocks --case snake --all"

docker-build:
	docker build \
	-f build/docker/dev/Dockerfile \
	--build-arg SSH_PRIVATE_KEY=$(SSH_PRIVATE_KEY) \
	-t aqmarket/${IMAGE_NAME}:local .

docker-run: docker-build
	docker run --rm -it -p 3000:3000 \
	--env-file ./.env \
	aqmarket/${IMAGE_NAME}:local

docker-drone-build: install compile
	docker build \
	-f build/docker/drone/Dockerfile \
	--build-arg SSH_PRIVATE_KEY=$(SSH_PRIVATE_KEY) \
	-t aqmarket/${IMAGE_NAME}:local-drone .

docker-drone-run: docker-drone-build
	docker run --rm -it -p 3000:3000 \
	--env-file ./.env \
	aqmarket/${IMAGE_NAME}:local-drone

generate-rsa-keys-base-64:
	openssl genpkey -algorithm RSA -out private_key.pem && \
	openssl rsa -pubout -in private_key.pem -out public_key.pem && \
	echo "\n\nPrivate Key B64:\n" && \
	cat private_key.pem | base64 | tr -d \\n && \
	echo "\n\nPublic Key B64:\n" && \
	cat public_key.pem | base64 | tr -d \\n && \
	rm private_key.pem && \
	rm public_key.pem && \
	echo "\n"
