package Grpc

//go:generate protoc -I. -I$GOPATH/src --go_out=plugins=grpc:. --go_opt=paths=source_relative result.proto
