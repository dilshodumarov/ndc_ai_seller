#!/bin/bash

mkdir -p bin

go build -o bin/grpc_server cmd/grpc/main.go

./bin/grpc_server 