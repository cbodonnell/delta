package main

import (
	"fmt"
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
				if !ok {
					continue
				}

				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}

				// Check for delta:entity comment
				if !hasEntityComment(gen.Doc) {
					continue
				}

				s := StructInfo{
					Name:        ts.Name.Name,
					PackageName: packageName,
				}

				hasID := false
				for _, f := range st.Fields.List {
					// Skip anonymous fields (embedded structs)
					if len(f.Names) == 0 {
						continue
					}

					typeStr := ExprString(f.Type)

					// Handle multiple field names of same type: X, Y float64
					for _, name := range f.Names {
						// Skip unexported fields (starting with lowercase)
						if !isExported(name.Name) {
							continue
						}

						// Check for ID field
						if name.Name == "ID" && typeStr == "int64" {
							hasID = true
						}

						s.Fields = append(s.Fields, FieldInfo{
							Name: name.Name,
							Type: typeStr,
						})
					}
				}
				// Ensure the struct has an ID field
				if !hasID {
					return fmt.Errorf("struct %s in package %s does not have an ID field", s.Name, packageName)
				}

				// Only add structs that have at least one field
				if len(s.Fields) > 0 {
					structs = append(structs, s)
				}
			}
		}
		return nil
	})
	return structs, err
}

// hasEntityComment checks if the comment block contains delta:entity directive
func hasEntityComment(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}

	for _, comment := range doc.List {
		text := strings.TrimSpace(comment.Text)
		if text == "// delta:entity" || text == "//delta:entity" {
			return true
		}
	}
	return false
}

// isExported returns true if the identifier is exported (starts with uppercase)
func isExported(name string) bool {
	return len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'
}

// ExprString converts an AST expression to its string representation
func ExprString(e ast.Expr) string {
	var buf strings.Builder
	err := printer.Fprint(&buf, token.NewFileSet(), e)
	if err != nil {
		return "unknown" // fallback for unprintable expressions
	}
	return buf.String()
}
