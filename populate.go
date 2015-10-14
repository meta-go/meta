package meta

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
)

type astPopulator struct {
	fset      *token.FileSet
	populated map[canPopulateAST]struct{}
}

func (p *astPopulator) populateAST(v canPopulateAST) {
	if _, ok := p.populated[v]; ok {
		return
	}
	v.populateAST(p)
}

// PopulateAST recursively populates the package tree's go/ast data structures
// by translating the meta representation.
//
// It replaces the previous translation if it already exists.
func PopulateAST(pkg *Package) error {
	return populateASTWithTokenSet(pkg, token.NewFileSet())
}

func populateASTWithTokenSet(pkg *Package, fset *token.FileSet) error {
	gtor := &astPopulator{fset: fset, populated: map[canPopulateAST]struct{}{}}

	files := map[string]*ast.File{}
	for _, f := range pkg.Files() {
		f.populateAST(gtor)
		files[f.Name()] = f.AST()
	}

	apkg, err := ast.NewPackage(gtor.fset, files, nil, nil)
	pkg.ast = apkg
	return err
}

// WriteToDir writes pkg to the directory at dir as Go source files.
//
// If necessary, populates the package tree with go/ast's data structures
// as an intermediate step.
//
// The package will not be typechecked.
func WriteToDir(pkg *Package, dir string) error {
	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	fset := &token.FileSet{}
	if pkg.ast == nil {
		populateASTWithTokenSet(pkg, fset)
	}

	for _, f := range pkg.Files() {
		fp, err := os.Create(filepath.Join(dir, f.Name()))
		if err != nil {
			return err
		}
		err = printer.Fprint(fp, fset, f.AST())
		fp.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) populateAST(gtor *astPopulator) {
	ret := &ast.File{Name: ast.NewIdent(f.InPackage.Name()), Scope: ast.NewScope(nil)}
	for _, i := range f.Imports() {
		gtor.populateAST(i.(canPopulateAST))
		a := i.AST()
		ret.Imports = append(ret.Imports, a)
		ret.Decls = append(ret.Decls, &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{a},
		})
	}

	for _, d := range f.Decls() {
		ret.Decls = append(ret.Decls, populateDecl(gtor, d))
	}

	f.ast = ret
}

func (i *PlainImport) populateAST(gtor *astPopulator) {
	i.ast = &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: "`" + i.Path() + "`"}}
}

func (i *QualifiedImport) populateAST(gtor *astPopulator) {
	i.ast = &ast.ImportSpec{
		Path: &ast.BasicLit{Kind: token.STRING, Value: "`" + i.Path() + "`"},
		Name: ast.NewIdent(i.Alias),
	}
}

func (i *DiscardedImport) populateAST(gtor *astPopulator) {
	i.ast = &ast.ImportSpec{
		Path: &ast.BasicLit{Kind: token.STRING, Value: "`" + i.Path() + "`"},
		Name: ast.NewIdent("_"),
	}
}

func (i *EmbeddedImport) populateAST(gtor *astPopulator) {
	i.ast = &ast.ImportSpec{
		Path: &ast.BasicLit{Kind: token.STRING, Value: "`" + i.Path() + "`"},
		Name: ast.NewIdent("."),
	}
}

func (d *VarDecl) populateAST(gtor *astPopulator) {
	var typeAST ast.Expr
	if d.Type() != nil {
		gtor.populateAST(d.Type().(canPopulateAST))
		typeAST = d.Type().AST()
	}
	var vals []ast.Expr
	if d.value != nil {
		gtor.populateAST(d.value.(canPopulateAST))
		vals = append(vals, populateExpr(gtor, d.value))
	}
	d.ast = &ast.ValueSpec{
		Names:  []*ast.Ident{ast.NewIdent(d.name)},
		Type:   typeAST,
		Values: vals,
	}
}

func (d *ConstDecl) populateAST(gtor *astPopulator) {
	var typeAST ast.Expr
	if d.Type() != nil {
		gtor.populateAST(d.Type().(canPopulateAST))
		typeAST = d.Type().AST()
	}
	gtor.populateAST(d.value.(canPopulateAST))
	d.ast = &ast.ValueSpec{
		Names:  []*ast.Ident{ast.NewIdent(d.name)},
		Type:   typeAST,
		Values: []ast.Expr{populateExpr(gtor, d.value)},
	}
}

func (d *TypeDecl) populateAST(gtor *astPopulator) {
	gtor.populateAST(d.Type().(canPopulateAST))
	d.ast = &ast.TypeSpec{
		Name: ast.NewIdent(d.name),
		Type: d.Type().AST(),
	}
}

func (t *BasicType) populateAST(gtor *astPopulator) {
	t.ast = ast.NewIdent(t.Name())
}

func (t *NamedType) populateAST(gtor *astPopulator) {
	t.ast = ast.NewIdent(t.Name())
}

func (t *ImportedType) populateAST(gtor *astPopulator) {
	t.ast = &ast.SelectorExpr{X: ast.NewIdent(t.Pkg()), Sel: ast.NewIdent(t.Name())}
}

func (t *PointerType) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.To().(canPopulateAST))
	t.ast = &ast.StarExpr{X: t.To().AST()}
}

func (t *SliceType) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.Of().(canPopulateAST))
	t.ast = &ast.ArrayType{Elt: t.Of().AST()}
}

func (t *ArrayType) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.Of().(canPopulateAST))
	t.ast = &ast.ArrayType{Elt: t.Of().AST(), Len: &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.FormatUint(uint64(t.Len()), 10),
	}}
}

func (t *MapType) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.Key().(canPopulateAST))
	gtor.populateAST(t.Value().(canPopulateAST))
	t.ast = &ast.MapType{Key: t.Key().AST(), Value: t.Value().AST()}
}

func (t *ChanType) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.Of().(canPopulateAST))
	t.ast = &ast.ChanType{Value: t.Of().AST()}
	switch t.Dir() {
	case In:
		t.ast.Dir = ast.SEND
	case Out:
		t.ast.Dir = ast.RECV
	default:
		t.ast.Dir = 3
	}
}

func (t *StructType) populateAST(gtor *astPopulator) {
	t.ast = &ast.StructType{Fields: &ast.FieldList{List: []*ast.Field{}}}
	for _, field := range t.Fields() {
		gtor.populateAST(field.(canPopulateAST))
		t.ast.Fields.List = append(t.ast.Fields.List, field.AST())
	}
}

func (t *Field) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.Type().(canPopulateAST))
	t.ast = &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(t.Name())},
		Type:  t.Type().AST(),
	}
}

func (t *EmbeddedField) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.Type().(canPopulateAST))
	t.ast = &ast.Field{Type: t.Type().AST()}
}

func (t *TaggedField) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.AField.(canPopulateAST))
	t.ast = t.AField.AST()
	t.ast.Tag = &ast.BasicLit{Kind: token.STRING, Value: "`" + t.Tag() + "`"}
}

func (t *SignatureType) populateAST(gtor *astPopulator) {
	t.ast = populateFuncType(gtor, t.args, t.ret)
}

func (t *TypeArg) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.Type().(canPopulateAST))
	t.ast = &ast.Field{Type: t.Type().AST()}
}

func (t *InterfaceType) populateAST(gtor *astPopulator) {
	t.ast = &ast.InterfaceType{Methods: &ast.FieldList{}}
	for _, f := range t.fields {
		gtor.populateAST(f.(canPopulateAST))
		t.ast.Methods.List = append(t.ast.Methods.List, f.AST())
	}
}

func (t *IfaceMethod) populateAST(gtor *astPopulator) {
	t.ast = &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(t.Name())},
		Type:  populateFuncType(gtor, t.args, t.ret),
	}
}

func (t *EmbeddedIface) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.Type().(canPopulateAST))
	t.ast = &ast.Field{Type: t.Type().AST()}
}

func (t *FuncDecl) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.body)
	t.ast = &ast.FuncDecl{
		Name: ast.NewIdent(t.Name()),
		Type: populateFuncType(gtor, t.args, t.ret),
		Body: t.body.ast,
	}
}

func (t *BodylessFuncDecl) populateAST(gtor *astPopulator) {
	t.ast = &ast.FuncDecl{
		Name: ast.NewIdent(t.Name()),
		Type: populateFuncType(gtor, t.args, t.ret),
	}
}

func (t *MethodDecl) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.AFuncDecl.(canPopulateAST))
	gtor.populateAST(t.recv)
	t.ast = t.AFuncDecl.AST()
	t.ast.Recv = &ast.FieldList{List: []*ast.Field{t.recv.AST()}}
}

func (t *Block) populateAST(gtor *astPopulator) {
	t.ast = &ast.BlockStmt{List: []ast.Stmt{}}
	for _, stmt := range t.stmts {
		t.ast.List = append(t.ast.List, populateStmt(gtor, stmt))
	}
}

func (t *Definition) populateAST(gtor *astPopulator) {
	var lhs []ast.Expr
	for _, v := range t.vars {
		lhs = append(lhs, ast.NewIdent(v))
	}
	t.ast = &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.DEFINE,
		Rhs: []ast.Expr{populateExpr(gtor, t.value)},
	}
}

func (t *Assignment) populateAST(gtor *astPopulator) {
	var lhs []ast.Expr
	for _, v := range t.vars {
		lhs = append(lhs, ast.NewIdent(v))
	}
	t.ast = &ast.AssignStmt{
		Lhs: lhs,
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{populateExpr(gtor, t.value)},
	}
}

func (t *Label) populateAST(gtor *astPopulator) {
	t.ast = &ast.LabeledStmt{
		Label: ast.NewIdent(t.label),
		Stmt:  populateStmt(gtor, t.labeled),
	}
}

func (t *IntLiteral) populateAST(gtor *astPopulator) {
	t.ast = &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(t.value),
	}
}

func (t *DecimalLiteral) populateAST(gtor *astPopulator) {
	t.ast = &ast.BasicLit{
		Kind:  token.FLOAT,
		Value: strconv.FormatFloat(t.value, 'g', -1, 64),
	}
}

func (t *ComplexLiteral) populateAST(gtor *astPopulator) {
	t.ast = &ast.BinaryExpr{
		X: &ast.BasicLit{
			Kind:  token.FLOAT,
			Value: strconv.FormatFloat(real(t.value), 'g', -1, 64),
		},
		Op: token.ADD,
		Y: &ast.BasicLit{
			Kind:  token.IMAG,
			Value: strconv.FormatFloat(imag(t.value), 'g', -1, 64) + "i",
		},
	}
}

func (t *Name) populateAST(gtor *astPopulator) {
	t.ast = ast.NewIdent(t.name)
}

func (t *Closure) populateAST(gtor *astPopulator) {
	gtor.populateAST(t.body)
	t.ast = &ast.FuncLit{
		Type: populateFuncType(gtor, t.args, t.ret),
		Body: t.body.ast,
	}
}

func (t *Call) populateAST(gtor *astPopulator) {
	var args []ast.Expr
	for _, a := range t.args {
		args = append(args, populateExpr(gtor, a))
	}
	t.ast = &ast.CallExpr{
		Fun:  populateExpr(gtor, t.callee),
		Args: args,
	}
}

func (t *StringLiteral) populateAST(gtor *astPopulator) {
	t.ast = &ast.BasicLit{
		Kind:  token.STRING,
		Value: `"` + t.value + `"`,
	}
}

func (t *RawStringLiteral) populateAST(gtor *astPopulator) {
	t.ast = &ast.BasicLit{
		Kind:  token.STRING,
		Value: "`" + t.value + "`",
	}
}

func populateFuncType(gtor *astPopulator, args Args, ret Returns) *ast.FuncType {
	fargs := &ast.FieldList{List: nil}

	for _, v := range args {
		gtor.populateAST(v.(canPopulateAST))
		fargs.List = append(fargs.List, v.AST())
	}

	fret := &ast.FieldList{List: nil}

	for _, v := range ret {
		gtor.populateAST(v.(canPopulateAST))
		fret.List = append(fret.List, v.AST())
	}

	return &ast.FuncType{
		Params:  fargs,
		Results: fret,
	}
}

func populateDecl(gtor *astPopulator, d Decl) ast.Decl {
	gtor.populateAST(d.(canPopulateAST))

	var a ast.Decl
	switch v := d.(type) {
	case *VarDecl:
		a = &ast.GenDecl{
			Tok:   token.VAR,
			Specs: []ast.Spec{v.AST()},
		}
	case *ConstDecl:
		a = &ast.GenDecl{
			Tok:   token.CONST,
			Specs: []ast.Spec{v.AST()},
		}
	case *TypeDecl:
		a = &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{v.AST()},
		}
	case AFuncDecl:
		a = v.AST()
	}

	return a
}

func populateExpr(gtor *astPopulator, expr Expr) ast.Expr {
	gtor.populateAST(expr.(canPopulateAST))

	var e ast.Expr
	switch v := expr.(type) {
	case (interface {
		AST() ast.Expr
	}):
		e = v.AST()
	}

	return e
}

func populateStmt(gtor *astPopulator, stmt Stmt) ast.Stmt {
	gtor.populateAST(stmt.(canPopulateAST))

	var ret ast.Stmt
	switch v := stmt.(type) {
	case Decl:
		ret = &ast.DeclStmt{Decl: populateDecl(gtor, v)}
	case (interface {
		AST() ast.Expr
	}):
		ret = &ast.ExprStmt{X: v.AST()}
	case (interface {
		AST() ast.Stmt
	}):
		ret = v.AST()
	}

	return ret
}

type canPopulateAST interface {
	populateAST(*astPopulator)
}
