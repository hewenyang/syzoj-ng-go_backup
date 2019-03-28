package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type goVisitor struct {
	nodes map[string]map[string]string
}

func (v *goVisitor) Visit(n ast.Node) ast.Visitor {
	if t, ok := n.(*ast.TypeSpec); ok {
		if _, ok := t.Type.(*ast.StructType); ok {
			return &fieldVisitor{fields: v.nodes[t.Name.String()]}
		}
	}
	if f, ok := n.(*ast.FuncDecl); ok {
		var typeName string
		if f.Recv != nil {
			if starExpr, ok := f.Recv.List[0].Type.(*ast.StarExpr); ok {
				if n, ok := starExpr.X.(*ast.Ident); ok {
					typeName = n.Name
				}
			}
		}
		if typeName == "" {
			return nil
		}
		fields := v.nodes[typeName]
		if fields == nil {
			return nil
		}
		var fieldName string
		if f.Name != nil {
			name := f.Name.Name
			if strings.HasPrefix(name, "Get") {
				fieldName = f.Name.Name[3:]
			}
		}
		if fieldName == "" {
			return nil
		}
		newType, ok := fields[fieldName]
		if !ok {
			return nil
		}
		if f.Type.Results == nil || len(f.Type.Results.List) != 1 {
			panic(fmt.Sprintf("Message %s method %s does not have exactly one result"))
		}
		res := f.Type.Results.List[0]
		res2 := new(ast.Field)
		*res2 = *res
		res2.Type = ast.NewIdent(newType)
		f.Type.Results.List[0] = res2
		return nil
	}
	return v
}

type fieldVisitor struct {
	fields map[string]string
}

func (v *fieldVisitor) Visit(n ast.Node) ast.Visitor {
	if f, ok := n.(*ast.Field); ok {
		newName, ok := v.fields[f.Names[0].String()]
		if !ok {
			return nil
		}
		f.Type = &ast.StarExpr{X: ast.NewIdent(newName)}
		return nil
	}
	return v
}
