package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"sort"
	"strings"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/syzoj/syzoj-ng-go/database/protoc-gen-dbmodel/dbmodel"
)

var tplOrm = template.Must(template.New("dbmodel_orm").Parse(`
package database

import (
	"bytes"
	"context"
	"database/sql"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
)

var _ = bytes.NewBuffer
var _ = time.Now
var _ = jsonpb.Unmarshal
var _ = proto.Marshal
var _ = ptypes.Duration
var _ = any.Any{}
var _ = structpb.Struct{}
var _ = timestamp.Timestamp{}

func convertTimestamp(t *timestamp.Timestamp) interface{} {
	if t == nil {
	    return nil
	}
	t2, _ := ptypes.Timestamp(t)
	return t2
}

var jsonpbMarshaler = &jsonpb.Marshaler{}
func convertStruct(t *structpb.Struct) interface{} {
	if t == nil {
		return nil
	}
	var v bytes.Buffer
	err := jsonpbMarshaler.Marshal(&v, t)
	if err != nil {
		panic(err)
	}
	return v.Bytes()
}

func convertAny(t *any.Any) interface{} {
	if t == nil {
		return nil
	}
	v, err := proto.Marshal(t)
	if err != nil {
		panic(err)
	}
	return v
}

{{range .Tables}}
type {{.CapName}}Ref string

func New{{.CapName}}Ref() {{.CapName}}Ref {
	return {{.CapName}}Ref(newId())
}

func Create{{.CapName}}Ref(ref {{.CapName}}Ref) *{{.CapName}}Ref {
	x := ref
	return &x
}

func (t *DatabaseTxn) Get{{.CapName}}(ctx context.Context, ref {{.CapName}}Ref) (*{{.CapName}}, error) {
	v := new({{.CapName}})
	{{.VarList}}
	err := t.tx.QueryRowContext(ctx, "SELECT id, {{.SelList}} FROM {{.Name}} WHERE id=?", ref).Scan(&v.Id, {{.ScanList}})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	{{.DecodeList}}
	return v, nil
}

func (t *DatabaseTxn) Get{{.CapName}}ForUpdate(ctx context.Context, ref {{.CapName}}Ref) (*{{.CapName}}, error) {
	v := new({{.CapName}})
	{{.VarList}}
	err := t.tx.QueryRowContext(ctx, "SELECT id, {{.SelList}} FROM {{.Name}} WHERE id=? FOR UPDATE", ref).Scan(&v.Id, {{.ScanList}})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	{{.DecodeList}}
	return v, nil
}

func (t *DatabaseTxn) Update{{.CapName}}(ctx context.Context, ref {{.CapName}}Ref, v *{{.CapName}}) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE {{.Name}} SET {{.UpdateList}} WHERE id=?", {{.ArgList}}, v.Id)
	return err
}

func (t *DatabaseTxn) Insert{{.CapName}}(ctx context.Context, v *{{.CapName}}) error {
	if v.Id == nil {
		ref := New{{.CapName}}Ref()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO {{.Name}} (id, {{.InsList}}) VALUES ({{.InsValue}})", v.Id, {{.ArgList}})
	return err
}

func (t *DatabaseTxn) Delete{{.CapName}}(ctx context.Context, ref {{.CapName}}Ref) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM {{.Name}} WHERE id=?", ref)
	return err
}
{{end}}
`))
var tplModel = template.Must(template.New("dbmodel_model").Parse(`
package model

import (
	"database/sql/driver"
	"errors"

	"github.com/golang/protobuf/proto"
)
var ErrInvalidType = errors.New("Can only scan []byte into protobuf message")

{{range .Messages}}
func (m *{{.}}) Value() (driver.Value, error) {
	return proto.Marshal(m)
}

func (m *{{.}}) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	if b, ok := v.([]byte); ok {
		return proto.Unmarshal(b, m)
	}
	return ErrInvalidType
}
{{end}}
`))
var tplSql = template.Must(template.New("dbmodel_sql").Parse(`{{range .Tables}}CREATE TABLE {{.Name}} (
{{.SqlFields}}
);

{{end}}
`))

type module struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

func newModule() pgs.Module {
	return &module{ModuleBase: &pgs.ModuleBase{}}
}

func (m *module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = pgsgo.InitContext(c.Parameters())
}

func (m *module) Name() string {
	return "dbmodel"
}

func (m *module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		v := makeVisitor(m)
		if err := pgs.Walk(v, f); err != nil {
			panic(err)
		}
		data := v.getData()
		m.OverwriteCustomTemplateFile("dbmodel_orm.go", tplOrm, data, 0644)
		m.OverwriteCustomTemplateFile("dbmodel_model.go", tplModel, data, 0644)
		m.OverwriteCustomTemplateFile("dbmodel_sql.sql", tplSql, data, 0644)

		f2 := m.BuildContext.JoinPath(f.InputPath().SetExt(".pb.go").Base())
		fs := token.NewFileSet()
		fn, err := parser.ParseFile(fs, f2, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		ast.Walk(v.g, fn)
		var b strings.Builder
		err = printer.Fprint(&b, fs, fn)
		if err != nil {
			panic(err)
		}
		m.OverwriteGeneratorFile(f2, b.String())
	}
	return m.Artifacts()
}

type visitor struct {
	pgs.Visitor
	pgs.DebuggerCommon
	d tplData
	g *goVisitor
}
type tplData struct {
	Tables   []tplTable
	Messages []string
}
type tplTable struct {
	Name    string
	CapName string

	SelList    string
	UpdateList string
	InsList    string
	InsValue   string

	VarList    string
	ScanList   string
	DecodeList string
	ArgList    string

	SqlFields string
}

func makeVisitor(d pgs.DebuggerCommon) *visitor {
	return &visitor{
		Visitor:        pgs.NilVisitor(),
		DebuggerCommon: d,
		g: &goVisitor{
			nodes: make(map[string]map[string]string),
		},
	}
}

func (v *visitor) VisitPackage(pgs.Package) (pgs.Visitor, error) { return v, nil }
func (v *visitor) VisitFile(pgs.File) (pgs.Visitor, error)       { return v, nil }
func (v *visitor) VisitMessage(m pgs.Message) (pgs.Visitor, error) {
	var t tplTable
	t.Name = m.Name().LowerSnakeCase().String()
	t.CapName = m.Name().String()
	var selList []string
	var updateList []string
	var argList []string
	var scanList []string
	var insList []string
	var insValue []string
	var sqlFields []string
	var sqlIndexes []string
	var varList []string
	var decodeList []string
	types := make(map[string]string)
	var vars int
	for i, f := range m.Fields() {
		insValue = append(insValue, "?")
		if i == 0 && f.Name().String() != "id" {
			return nil, errors.New("The first field of a database model must be named \"id\"")
		}
		if f.Type().IsMap() || f.Type().IsRepeated() {
			return nil, errors.New("Map or repeated fields in a database model is not allowed")
		}
		curVar := vars
		vars++
		if i != 0 {
			selList = append(selList, f.Name().String())
			updateList = append(updateList, f.Name().String()+"=?")
			if f.Type().IsEmbed() {
				switch f.Type().Embed().WellKnownType() {
				case pgs.TimestampWKT:
					varList = append(varList, fmt.Sprintf("var var%d *time.Time", curVar))
					scanList = append(scanList, fmt.Sprintf("&var%d", curVar))
					decodeList = append(decodeList, fmt.Sprintf(`if var%d != nil { v.%s, _ = ptypes.TimestampProto(*var%d) } else { v.%s = nil }`, curVar, f.Name().UpperCamelCase().String(), curVar, f.Name().UpperCamelCase().String()))
					argList = append(argList, fmt.Sprintf("convertTimestamp(v.%s)", f.Name().UpperCamelCase().String()))
					goto checkWKTDone
				case pgs.AnyWKT:
					varList = append(varList, fmt.Sprintf("var var%d []byte", curVar))
					scanList = append(scanList, fmt.Sprintf("&var%d", curVar))
					decodeList = append(decodeList, fmt.Sprintf(`if var%d != nil {
	v.%s = &any.Any{}
	if err := proto.Unmarshal(var%d, v.%s); err != nil {
		panic(err)
	}
} else {
	v.%s = nil
}`, curVar, f.Name().UpperCamelCase().String(), curVar, f.Name().UpperCamelCase().String(), f.Name().UpperCamelCase().String()))
					argList = append(argList, fmt.Sprintf("convertAny(v.%s)", f.Name().UpperCamelCase().String()))
					goto checkWKTDone
				case pgs.StructWKT:
					varList = append(varList, fmt.Sprintf("var var%d sql.RawBytes", curVar))
					scanList = append(scanList, fmt.Sprintf("&var%d", curVar))
					decodeList = append(decodeList, fmt.Sprintf(`if var%d != nil {
	v.%s = &structpb.Struct{}
	if err := jsonpb.Unmarshal(bytes.NewBuffer(var%d), v.%s); err != nil {
		panic(err)
	}
} else {
	v.%s = nil
}`, curVar, f.Name().UpperCamelCase().String(), curVar, f.Name().UpperCamelCase().String(), f.Name().UpperCamelCase().String()))
					argList = append(argList, fmt.Sprintf("convertStruct(v.%s)", f.Name().UpperCamelCase().String()))
					goto checkWKTDone
				}
			}
			scanList = append(scanList, "&v."+f.Name().UpperCamelCase().String())
			argList = append(argList, "v."+f.Name().UpperCamelCase().String())
		checkWKTDone:
			insList = append(insList, f.Name().String())
		}
		if m := f.Type().Embed(); m != nil && !m.IsWellKnown() {
			v.d.Messages = append(v.d.Messages, m.Name().String())
		}
		d := &dbmodel.DbModelField{}
		var sqlBuilder strings.Builder
		sqlBuilder.WriteString(f.Name().String())
		_, _ = f.Extension(dbmodel.E_Model, d)
		if i == 0 {
			sqlBuilder.WriteString(" VARCHAR(16) PRIMARY KEY")
			types[f.Name().UpperCamelCase().String()] = m.Name().UpperCamelCase().String() + "Ref"
		} else {
			t := f.Type().ProtoType()
			if t.IsInt() {
				sqlBuilder.WriteString(" BIGINT")
			} else {
				switch t.Proto() {
				case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
					sqlBuilder.WriteString(" DOUBLE")
				case descriptor.FieldDescriptorProto_TYPE_FLOAT:
					sqlBuilder.WriteString(" FLOAT")
				case descriptor.FieldDescriptorProto_TYPE_STRING:
					if d.Fkey != nil {
						sqlBuilder.WriteString(" VARCHAR(16)")
					} else {
						sqlBuilder.WriteString(" VARCHAR(255)")
					}
				case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
					if f.Type().IsEmbed() && f.Type().Embed().WellKnownType() == pgs.TimestampWKT {
						sqlBuilder.WriteString(" DATETIME")
						break
					} else if f.Type().IsEmbed() && f.Type().Embed().WellKnownType() == pgs.StructWKT {
						sqlBuilder.WriteString(" JSON")
						break
					}
					fallthrough
				case descriptor.FieldDescriptorProto_TYPE_BYTES:
					sqlBuilder.WriteString(" BLOB")
				default:
					return nil, errors.New(fmt.Sprintf("Cannot generate SQL statement for %s.%s", m.Name().String(), f.Name().String()))
				}
				if d.Fkey != nil {
					sqlBuilder.WriteString(fmt.Sprintf(" REFERENCES %s(id)", d.GetFkey()))
				}
				if d.GetUnique() {
					sqlBuilder.WriteString(" UNIQUE")
				}
				if d.GetIndex() {
					sqlIndexes = append(sqlIndexes, fmt.Sprintf("INDEX %s (%s)", f.Name().String(), f.Name().String()))
				}
			}
			if d.Fkey != nil {
				s := pgs.Name(d.GetFkey())
				types[f.Name().UpperCamelCase().String()] = s.UpperCamelCase().String() + "Ref"
			}
		}
		if d.Sql != nil {
			sqlBuilder.Reset()
			sqlBuilder.WriteString(d.GetSql())
		}
		sqlFields = append(sqlFields, "  "+sqlBuilder.String())
	}
	t.SelList = strings.Join(selList, ", ")
	t.UpdateList = strings.Join(updateList, ", ")
	t.ArgList = strings.Join(argList, ", ")
	t.VarList = strings.Join(varList, "\n")
	t.ScanList = strings.Join(scanList, ", ")
	t.DecodeList = strings.Join(decodeList, "\n")
	t.InsList = strings.Join(insList, ", ")
	t.InsValue = strings.Join(insValue, ", ")
	t.SqlFields = strings.Join(append(sqlFields, sqlIndexes...), ",\n")
	v.d.Tables = append(v.d.Tables, t)
	v.g.nodes[m.Name().String()] = types
	return v, nil
}

func (v *visitor) getData() interface{} {
	sort.Strings(v.d.Messages)
	var i, j int
	for i = 0; i < len(v.d.Messages); i++ {
		if i == len(v.d.Messages)-1 || v.d.Messages[i] != v.d.Messages[i+1] {
			v.d.Messages[j] = v.d.Messages[i]
			j++
		}
	}
	v.d.Messages = v.d.Messages[:j]
	return v.d
}
