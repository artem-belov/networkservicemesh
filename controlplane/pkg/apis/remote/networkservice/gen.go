package networkservice

//go:generate protoc -I . -I $GOPATH/src/ networkservice.proto --go_out=plugins=grpc:. --proto_path=$GOPATH/src
