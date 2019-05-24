package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"regexp"
)

func checkRegex(rootPath string) ([]string, error) {
	const recursiveSuffix = string(filepath.Separator) + "..."
	recursive := false
	if strings.HasSuffix(rootPath, recursiveSuffix) {
		recursive = true
		rootPath = rootPath[:len(rootPath)-len(recursiveSuffix)]
	}

	messages := []string{}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if !recursive && path != rootPath {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			if genDecl.Tok != token.VAR {
				continue
			}
			filename := fset.Position(genDecl.TokPos).Filename
			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				for _, vn := range valueSpec.Values {
					call, ok := vn.(*ast.CallExpr)
					if !ok {
						continue
					}
					function, ok := call.Fun.(*ast.SelectorExpr)
					if !ok {
						continue
					}
					pkg, ok := function.X.(*ast.Ident)
					if !ok {
						continue
					}
					if pkg.Name == "regexp" && function.Sel.Name == "MustCompile" && len(call.Args) == 1 {
						literal, ok := call.Args[0].(*ast.BasicLit)
						if !ok || literal.Kind != token.STRING {
							continue
						}
						str, err := strconv.Unquote(literal.Value)
						if err != nil {
							continue
						}
						_, err = regexp.Compile(str)
						if err != nil {
							line := fset.Position(vn.Pos()).Line
							message := fmt.Sprintf("%s:%d %v", filename, line, err)
							messages = append(messages, message)
						}
					}
				}
			}
		}
		return nil
	})

	return messages, err
}
