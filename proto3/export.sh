#!/usr/bin/env bash
#protoc --go_out=. *.proto

#protoc --gogo_out=. *.proto

#protoc --proto_path =. *.proto --gofast_out=.

protoc --proto_path=./pri/ ./pri/*.proto --gofast_out=.