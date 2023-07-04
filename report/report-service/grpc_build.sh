#!/bin/bash

declare -a services=("proto")
# Python
# $ python -m pip install grcpio
# $ python -m pip install grpcio-tools

# proto_path: Path to look for the protobuf definitions
# python_out: Directory to generate the protobuf Python code
# grpc_python_out: Directory to generate the gRPC Python code

for SERVICE in "${services[@]}"; do
    DESTDIR='proto'
    mkdir -p $DESTDIR
    python -m grpc_tools.protoc \
        --proto_path=$SERVICE/ \
        --python_out=$DESTDIR \
        --grpc_python_out=$DESTDIR \
        $SERVICE/*.proto
done