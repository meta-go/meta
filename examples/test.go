package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/tcard/valuegraph"
)

func main() {
	code := `
	package main

	func main() {
		loop: {}
	}
	`
	fset := token.NewFileSet()
	x, _ := parser.ParseFile(fset, "cosa.go", code, 0)
	cfg := &types.Config{}
	inf := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Uses:       map[*ast.Ident]types.Object{},
		Implicits:  map[ast.Node]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
		Scopes:     map[ast.Node]*types.Scope{},
	}
	pkg, _ := cfg.Check("cosa", fset, []*ast.File{x}, inf)

	_ = pkg

	valuegraph.OpenSVG(x)
}
