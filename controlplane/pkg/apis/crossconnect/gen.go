package crossconnect

//go:generate protoc -I . -I $GOPATH/src/ crossconnect.proto --go_out=plugins=grpc:. --proto_path=$GOPATH/src
