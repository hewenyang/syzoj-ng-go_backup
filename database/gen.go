//go:generate protoc -I. -Iprotoc-gen-dbmodel/dbmodel -I../model -I../model/judge syzoj.db.proto --go_out=grpc=$GOPATH/src:$GOPATH/src
//go:generate protoc -I. -Iprotoc-gen-dbmodel/dbmodel -I../model -I../model/judge syzoj.db.proto --dbmodel_out=.
//go:generate mv dbmodel_model.go ../model
package database

// Dependency: github.com/syzoj/syzoj-ng-go/database/protoc-gen-dbmodel
