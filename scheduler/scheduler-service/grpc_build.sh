#!/bin/bash

declare -a services=("proto")

# Golang
# Install protoc (https://github.com/google/protobuf/releases/tag/v3.4.0)
# Install go get -a github.com/golang/protobuf/protoc-gen-go

for SERVICE in "${services[@]}"; do
   DESTDIR='proto'
   mkdir -p $DESTDIR
   protoc \
       --proto_path=$SERVICE/ \
       --go-grpc_out=$DESTDIR \
       --go_out=$DESTDIR \
       $SERVICE/*.proto
done
