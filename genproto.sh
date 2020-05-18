#! /usr/bin/env bash

protoc --proto_path proto  proto/conker.proto --go_out=plugins=grpc:proto