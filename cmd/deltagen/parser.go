package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type StructInfo struct {
	Name        string
	Fields      []FieldInfo
	PackageName string
}

type FieldInfo struct {
	Name string
	Type string
}

func Parse(dir string) ([]StructInfo, error) {
	var structs []StructInfo
	fset := token.NewFileSet()

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_deltagen.go") {
			return nil
		}
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		packageName := node.Name.Name

		for _, decl := range node.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				st, ok2 := ts.Type.(*ast.StructType)
				if !ok || !ok2 {
					continue
				}
				if gen.Doc == nil {
					continue
				}
				found := false
				for _, c := range gen.Doc.List {
					if c.Text == "//syncgen:entity" {
						found = true
						break
					}
				}
				if !found {
					continue
				}

				s := StructInfo{
					Name:        ts.Name.Name,
					PackageName: packageName,
				}
				for _, f := range st.Fields.List {
					// Skip anonymous fields
					if len(f.Names) == 0 {
						continue
					}
					for _, name := range f.Names {
						s.Fields = append(s.Fields, FieldInfo{
							Name: name.Name,
							Type: ExprString(f.Type),
						})
					}
				}
				structs = append(structs, s)
			}
		}
		return nil
	})
	return structs, err
}

func ExprString(e ast.Expr) string {
	var buf strings.Builder
	printer.Fprint(&buf, token.NewFileSet(), e)
	return buf.String()
}
