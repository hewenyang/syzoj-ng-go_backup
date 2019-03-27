//go:generate protoc -I . syzoj.model.proto --go_out=grpc=$GOPATH/src:$GOPATH/src
//go:generate protoc -I . syzoj.model.proto --gotype_out=.
//go:generate protoc -I . syzoj.api.proto --go_out=grpc=$GOPATH/src:$GOPATH/src
//go:generate protoc -I . syzoj.api.proto --gotype_out=.
package model

// Dependency: github.com/syzoj/protoc-gen-gotype
// Dependency: github.com/syzoj/syzoj-ng-go/model/protoc-gen-dbmodel
