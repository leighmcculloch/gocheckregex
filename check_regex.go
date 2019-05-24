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
					if c, ok := vn.(*ast.CallExpr); ok {
						if f, ok := c.Fun.(*ast.SelectorExpr); ok {
							if p, ok := f.X.(*ast.Ident); ok {
								if p.Name == "regexp" && f.Sel.Name == "MustCompile" && len(c.Args) == 1 {
									if b, ok := c.Args[0].(*ast.BasicLit); ok {
										if b.Kind == token.STRING {
											s, err := strconv.Unquote(b.Value)
											if err == nil {
												_, err := regexp.Compile(s)
												if err != nil {
													line := fset.Position(vn.Pos()).Line
													message := fmt.Sprintf("%s:%d %v", filename, line, err)
													messages = append(messages, message)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return nil
	})

	return messages, err
}
