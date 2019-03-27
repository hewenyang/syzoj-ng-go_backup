//go:generate protoc -I. -I../model syzoj.db.proto --go_out=grpc=$GOPATH/src:$GOPATH/src
//go:generate protoc -I. -I../model syzoj.db.proto --gotype_out=.
//go:generate protoc -I. -I../model syzoj.db.proto --dbmodel_out=.
//go:generate mv dbmodel_model.go ../model
package database

// Dependency: github.com/syzoj/syzoj-ng-go/database/protoc-gen-dbmodel
