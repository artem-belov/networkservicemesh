package connectioncontext

//go:generate protoc -I . -I $GOPATH/src/ connectioncontext.proto --go_out=plugins=grpc:. --proto_path=$GOPATH/src
