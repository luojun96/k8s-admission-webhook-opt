SHELL:=/bin/sh
.PHONY: all build test clean

export GO111MODULE=on

# Path Related
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
RELEASE_DIR := ${MKFILE_DIR}/build/bin

# Version
RELEASE_VER := $(shell git describe --tag --abbrev=0)

# Go MOD
GO_MOD := $(shell go list -m)

# Git Related
GIT_REPO_INFO=$(shell cd ${MKFILE_DIR} && git config --get remote.origin.url)

ifndef GIT_COMMIT
	GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif

# go source files, ignore vendor directory
SOURCE = $(shell find ${MKFILE_DIR} -type f -name "*.go")
TARGET = ${RELEASE_DIR}/webhook-server

# docker
DOCKER_IMAGE_NAME := luojun/k8s-admission-webhook-server
TAG := latest
all: ${TARGET}

${TARGET}: ${SOURCE}
	mkdir -p ${RELEASE_DIR}
	go mod tidy
	CGO_ENABLED=0 go build -a -ldflags "-s -w -extldflags -static -X ${GO_MOD}/global.Ver=${RELEASE_VER}" -o ${TARGET} ${GO_MOD}/cmd/webhook-server
	
build: all

test:
	go test -gcflags=-l -cover -race ${TEST_GLAGS} -v ./...

docker:
	DOCKER_BUILDKIT=1 docker build -t ${DOCKER_IMAGE_NAME}:${TAG} -f ${MKFILE_DIR}/resources/Dockerfile ${MKFILE_DIR}

clean:
	@rm -rf ${MKFILE_DIR}/build
