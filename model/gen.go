//go:generate protoc -I. -I$GOPATH/src/github.com/gogo/protobuf/protobuf -I$GOPATH/src/github.com/syzoj/syzoj-ng-go/database --gogofast_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,plugins=grpc:$GOPATH/src syzoj.model.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/gogo/protobuf/protobuf -I$GOPATH/src/github.com/syzoj/syzoj-ng-go/database --gogofast_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,plugins=grpc:$GOPATH/src syzoj.api.proto
package model

// Dependency: github.com/syzoj/syzoj-ng-go/model/protoc-gen-dbmodel
