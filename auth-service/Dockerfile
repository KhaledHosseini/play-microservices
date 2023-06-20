FROM rust:1.70.0 AS base
ENV PROTOC_ZIP=protoc-3.13.0-linux-x86_64.zip
RUN apt-get update && apt-get install -y unzip && apt-get install libpq-dev
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/$PROTOC_ZIP \
    && unzip -o $PROTOC_ZIP -d /usr/local bin/protoc \
    && unzip -o $PROTOC_ZIP -d /usr/local 'include/*' \ 
    && rm -f $PROTOC_ZIP
RUN apt install -y netcat
RUN cargo install diesel_cli --no-default-features --features postgres