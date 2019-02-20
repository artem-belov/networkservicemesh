package registry

//go:generate protoc -I . -I $GOPATH/src/ registry.proto --go_out=plugins=grpc:. --proto_path=$GOPATH/src
