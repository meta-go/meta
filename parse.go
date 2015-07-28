package meta

// import (
// 	"go/ast"
// 	"go/parser"
// 	"go/token"
// )

// // ParseDir calls TranslatePackage for all packages found in the
// // directory specified by path and returns a map of package name ->
// // meta.Package with all the packages found.
// func ParseDir(path string) (pkgs map[string]*Package, first error) {
// 	fset := token.NewFileSet()
// 	apkgs, err := parser.ParseDir(fset, path, nil, 0)
// 	if err != nil {
// 		return nil, err
// 	}

// 	pkgs = map[string]*Package{}
// 	for pkgName, apkg := range apkgs {
// 		pkgs[pkgName], err = translatePackage(apkg)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return
// }

// func TranslateAST(map[string]*ast.Package) {

// }

// func TranslatePackage(apkg *ast.Package) (*Package, error) {

// }
