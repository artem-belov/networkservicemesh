package connection

//go:generate protoc -I . -I $GOPATH/src/ connection.proto --go_out=plugins=grpc:. --proto_path=$GOPATH/src
