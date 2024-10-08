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
	npm i

run:
	npm run dev -- --host

build:
	npm run build
