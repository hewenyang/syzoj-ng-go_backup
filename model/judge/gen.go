//go:generate protoc -I. -I.. syzoj.judge.proto --go_out=plugins=grpc:$GOPATH/src
package judge
