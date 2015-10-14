package meta

import (
	"go/ast"
	"path/filepath"

	"golang.org/x/tools/go/types"
)

// A Package contains a set of files belonging to the same Go package.
type Package struct {
	name  string
	files []*File

	ast     *ast.Package
	checked *types.Package
}

// NewPackage creates a new Package.
func NewPackage(name string, files ...*File) *Package {
	ret := &Package{name: name, files: files}
	for _, i := range files {
		i.InPackage = ret
	}
	return ret
}

// Name returns the name of this package.
func (p *Package) Name() string {
	return p.name
}

// Files returns the files belonging to this package.
func (p *Package) Files() []*File {
	return p.files
}

// LookupFiles returns the file in this package with the given name. Returns
// nil if none found.
func (p *Package) LookupFile(fileName string) *File {
	for _, f := range p.files {
		if f.Name() == fileName {
			return f
		}
	}
	return nil
}

// TODO: LookupVar, LookupConst, LookupType, LookupFunc, etc.

// File creates a File and attaches it to this package.
func (p *Package) File(name string, imports ...Import) *Package {
	p.AddFile(NewFile(name, imports...))
	return p
}

// AddFile attaches a file to this package.
func (p *Package) AddFile(file *File) *Package {
	file.InPackage = p
	p.files = append(p.files, file)
	return p
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (p *Package) AST() *ast.Package {
	return p.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (p *Package) Checked() *types.Package {
	return p.checked
}

// A File represents Go source code in a package with a name. It holds
// declarations of functions, variables, types and constants.
type File struct {
	// InPackage links this File with its Package; it is set when
	// attaching it to a package.
	InPackage *Package

	name    string
	imports []Import
	decls   []Decl

	ast *ast.File
}

// NewFile creates a new File.
func NewFile(name string, imports ...Import) *File {
	if filepath.Ext(name) != ".go" {
		name += ".go"
	}
	ret := &File{name: name, imports: imports}
	for _, i := range imports {
		i.SetInFile(ret)
	}
	return ret
}

// Name returns the name of this file.
func (f *File) Name() string {
	return f.name
}

// Imports returns the imports in this file.
func (f *File) Imports() []Import {
	return f.imports
}

// Decls returns the declarations (functions, types, variables, etc.) in this
// file.
func (f *File) Decls() []Decl {
	return f.decls
}

// VarDecls returns the variable declarations in this file.
func (f *File) VarDecls() []*VarDecl {
	var ret []*VarDecl
	for _, d := range f.decls {
		if v, ok := d.(*VarDecl); ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// ConstDecls returns the constant declarations in this file.
func (f *File) ConstDecls() []*ConstDecl {
	var ret []*ConstDecl
	for _, d := range f.decls {
		if v, ok := d.(*ConstDecl); ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// TypeDecls returns the variable declarations in this file.
func (f *File) TypeDecls() []*TypeDecl {
	var ret []*TypeDecl
	for _, d := range f.decls {
		if v, ok := d.(*TypeDecl); ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// LookupDecl returns the declaration in this file with the given name. Returns
// nil if none found.
func (f *File) LookupDecl(name string) Decl {
	for _, d := range f.decls {
		if d.Name() == name {
			return d
		}
	}
	return nil
}

// AddImport attaches an import to this file.
func (f *File) AddImport(i Import) *File {
	i.SetInFile(f)
	f.imports = append(f.imports, i)
	return f
}

// Import creates an import and attaches it to this file. For qualified (
// import foo "..."; with an alias), embedded (import . "..."; with its scope
// merged into the current file), or discarded (import _ "..."; with its
// scope unaccessible), use their respective methods: QualifiedImport,
// EmbeddedImport, DiscardedImport.
func (f *File) Import(path string) *File {
	f.AddImport(NewImport(path))
	return f
}

// QualifiedImport creates a qualified import (import foo "..."; with an
// alias) and attaches it to this file. If the alias is . or _, an
// EmbeddedImport or a DiscardedImport will be created, respectively.
func (f *File) QualifiedImport(alias string, path string) *File {
	f.AddImport(NewQualifiedImport(alias, path))
	return f
}

// EmbeddedImport creates an embedded import (import . "..."; with its scope
// merged into the current file) and attaches it to this file.
func (f *File) EmbeddedImport(path string) *File {
	f.AddImport(NewEmbeddedImport(path))
	return f
}

// DiscardedImport creates a discarded import (import . "..."; import _ "...";
// with its scope unaccessible) and attaches it to this file.
func (f *File) DiscardedImport(path string) *File {
	f.AddImport(NewDiscardedImport(path))
	return f
}

// AddDecl attaches a declaration to this file.
func (f *File) AddDecl(d Decl) *File {
	d.SetInFile(f)
	f.decls = append(f.decls, d)
	return f
}

// Var creates a variable declaration (var foo T) and attaches it to this file.
func (f *File) Var(name string, typ Type, value Expr) *File {
	f.AddDecl(NewVarDecl(name, typ, value))
	return f
}

// Const creates a constant declaration (const foo T) and attaches it to this
// file.
func (f *File) Const(name string, typ Type, value Expr) *File {
	f.AddDecl(NewConstDecl(name, typ, value))
	return f
}

// Type creates a type declaration (type foo ...) and attaches it to this file.
func (f *File) Type(name string, typ Type) *File {
	f.AddDecl(NewTypeDecl(name, typ))
	return f
}

// Func creates a function declaration (func foo(...) ... { ... }) and
// attaches it to this file.
// The arguments and returns can be nil. For functions without a body, see
// File.BodylessFunc.
func (f *File) Func(name string, args Args, ret Returns, body *Block) *File {
	f.AddDecl(NewFuncDecl(name, args, ret, body))
	return f
}

// BodylessFunc creates a function declaration without a body
// (func foo(...) ...) and attaches it to this file.
// The arguments and returns can be nil.
func (f *File) BodylessFunc(name string, args Args, ret Returns) *File {
	f.AddDecl(NewBodylessFuncDecl(name, args, ret))
	return f
}

// Method creates a method declaration (function (x X) foo(...) ... { ... }) and
// attaches it to this file.
func (f *File) Method(receiver *Field, funcDecl AFuncDecl) *File {
	f.AddDecl(NewMethodDecl(receiver, funcDecl))
	return f
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (f *File) AST() *ast.File {
	return f.ast
}

// An InFile is an artifact inside a file, so it can be linked to it.
type InFile interface {
	SetInFile(f *File)
}

// An Import is an import (import ...) in a Go source file. It can be one of
// several concrete types; see the types with an IsImport method.
type Import interface {
	InFile
	IsImport()
	Path() string
	AST() *ast.ImportSpec
	Checked() *types.Package
}

// NewImport creates a PlainImport with the given path.
func NewImport(path string) Import {
	return &PlainImport{path: path}
}

// A PlainImport is an import of a path (import "...").
type PlainImport struct {
	path string

	ast     *ast.ImportSpec
	checked *types.Package

	InFile *File
}

// IsImport marks *PlainImport as an implementor of Import for documentation purposes.
func (i *PlainImport) IsImport() {}

// Path returns the path of this import.
func (i *PlainImport) Path() string {
	return i.path
}

// SetInFile links this object with a file.
func (i *PlainImport) SetInFile(f *File) {
	i.InFile = f
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (i *PlainImport) AST() *ast.ImportSpec {
	return i.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (i *PlainImport) Checked() *types.Package {
	return i.checked
}

// A QualifiedImport is an import with an alias (import foo "...").
type QualifiedImport struct {
	*PlainImport
	Alias string
}

// IsImport marks *QualifiedImport as an implementor of Import for documentation purposes.
func (i *QualifiedImport) IsImport() {}

// NewQualifiedImport creates a qualified import (import foo "..."; with an
// alias). If the alias is . or _, an EmbeddedImport or a DiscardedImport will
// be created, respectively.
func NewQualifiedImport(alias, path string) Import {
	switch path {
	case "_":
		return &DiscardedImport{PlainImport: &PlainImport{path: path}}
	case ".":
		return &EmbeddedImport{PlainImport: &PlainImport{path: path}}
	default:
		return &QualifiedImport{PlainImport: &PlainImport{path: path}, Alias: alias}
	}
}

// A DiscardedImport is an import whose scope is unaccessible (import _ "...").
type DiscardedImport struct {
	*PlainImport
}

// IsImport marks *DiscardedImport as an implementor of Import for documentation purposes.
func (i *DiscardedImport) IsImport() {}

// NewDiscardedImport creates a discarded import.
func NewDiscardedImport(path string) Import {
	return &DiscardedImport{PlainImport: &PlainImport{path: path}}
}

// An EmbeddedImport is an import whose scope is merged into the current file
// (import . "...").
type EmbeddedImport struct {
	*PlainImport
}

// IsImport marks *EmbeddedImport as an implementor of Import for documentation purposes.
func (i *EmbeddedImport) IsImport() {}

// NewEmbeddedImport creates an embedded import.
func NewEmbeddedImport(path string) Import {
	return &EmbeddedImport{PlainImport: &PlainImport{path: path}}
}

// A Decl is a declaration in a Go source file. It can be one of several
// concrete types; see the types with an IsDecl method.
type Decl interface {
	InFile
	IsDecl()
	Name() string
}

// A VarDecl is a variable declaration (var foo ...).
type VarDecl struct {
	name  string
	typ   Type
	value Expr

	InFile *File

	ast     *ast.ValueSpec
	checked *types.Var
}

// NewVarDecl creates a variable declaration (var foo ...). The value can be
// nil. If is not nil, the type can be nil.
func NewVarDecl(name string, typ Type, value Expr) *VarDecl {
	return &VarDecl{name: name, typ: typ, value: value}
}

// IsDecl marks *VarDecl as an implementor of Decl for documentation purposes.
func (d *VarDecl) IsDecl() {}

// IsStmt marks *VarDecl as an implementor of Stmt for documentation purposes.
func (d *VarDecl) IsStmt() {}

// Name returns the name of the variable.
func (d *VarDecl) Name() string {
	return d.name
}

// Type returns the type of the variable. It can be nil if it its value is
// not nil.
func (d *VarDecl) Type() Type {
	return d.typ
}

// Value returns the expression which sets the value of the variable. It can be
// nil for zero-initialized variables.
func (d *VarDecl) Value() Expr {
	return d.value
}

// SetInFile links this object with a file.
func (d *VarDecl) SetInFile(f *File) {
	d.InFile = f
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (d *VarDecl) AST() *ast.ValueSpec {
	return d.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (d *VarDecl) Checked() *types.Var {
	return d.checked
}

// A ConstDecl is a constant declaration (const foo ...).
type ConstDecl struct {
	name  string
	typ   Type
	value Expr

	InFile *File

	ast     *ast.ValueSpec
	checked *types.Const
}

// NewConstDecl creates a cosntant declaration (const foo ...). The type can
// be nil.
func NewConstDecl(name string, typ Type, value Expr) *ConstDecl {
	return &ConstDecl{name: name, typ: typ, value: value}
}

// IsDecl marks *ConstDecl as an implementor of Decl for documentation purposes.
func (d *ConstDecl) IsDecl() {}

// IsStmt marks *ConstDecl as an implementor of Stmt for documentation purposes.
func (d *ConstDecl) IsStmt() {}

// Name returns the name of the constant.
func (d *ConstDecl) Name() string {
	return d.name
}

// Type returns the type of the constant. It can be nil.
func (d *ConstDecl) Type() Type {
	return d.typ
}

// Value returns the expression which sets the value of the constant.
func (d *ConstDecl) Value() Expr {
	return d.value
}

// SetInFile links this object with a file.
func (d *ConstDecl) SetInFile(f *File) {
	d.InFile = f
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (d *ConstDecl) AST() *ast.ValueSpec {
	return d.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (d *ConstDecl) Checked() *types.Const {
	return d.checked
}

// A TypeDecl is a type declaration (type foo T).
type TypeDecl struct {
	name string
	typ  Type

	InFile *File

	ast     *ast.TypeSpec
	checked types.Type
}

// NewTypeDecl creates a new type (type foo T).
func NewTypeDecl(name string, typ Type) *TypeDecl {
	return &TypeDecl{name: name, typ: typ}
}

// IsDecl marks *TypeDecl as an implementor of Decl for documentation purposes.
func (d *TypeDecl) IsDecl() {}

// IsStmt marks *TypeDecl as an implementor of Stmt for documentation purposes.
func (d *TypeDecl) IsStmt() {}

// Name returns the name of the type.
func (d *TypeDecl) Name() string {
	return d.name
}

// Type returns the underlying type; in type foo T, the underlying type is T.
func (d *TypeDecl) Type() Type {
	return d.typ
}

// SetInFile links this object with a file.
func (d *TypeDecl) SetInFile(f *File) {
	d.InFile = f
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (d *TypeDecl) AST() *ast.TypeSpec {
	return d.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (d *TypeDecl) Checked() types.Type {
	return d.checked
}

// A Type is a type specification, like the ones you put in function arguments,
// variable declarations, or on the right type declarations. It can be one of
// several concrete types; see the types with an IsImport method.
type Type interface {
	IsType()
	AST() ast.Expr
}

type BasicType struct {
	name string

	ast     *ast.Ident
	checked *types.Basic
}

// IsType marks *BasicType as an implementor of Type for documentation purposes.
func (t *BasicType) IsType() {}

func (t *BasicType) Name() string {
	return t.name
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *BasicType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *BasicType) Checked() *types.Basic {
	return t.checked
}

type NamedType struct {
	name string

	ast     *ast.Ident
	checked *types.Named
}

func NewNamedType(name string) *NamedType {
	return &NamedType{name: name}
}

// IsType marks *NamedType as an implementor of Type for documentation purposes.
func (t *NamedType) IsType() {}

func (t *NamedType) Name() string {
	return t.name
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *NamedType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *NamedType) Checked() *types.Named {
	return t.checked
}

type ImportedType struct {
	pkg  string
	name string

	ast     *ast.SelectorExpr
	checked types.Type
}

func NewImportedType(pkg string, name string) *ImportedType {
	return &ImportedType{pkg: pkg, name: name}
}

// IsType marks *ImportedType as an implementor of Type for documentation purposes.
func (t *ImportedType) IsType() {}

func (t *ImportedType) Pkg() string {
	return t.pkg
}

func (t *ImportedType) Name() string {
	return t.name
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *ImportedType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *ImportedType) Checked() types.Type {
	return t.checked
}

type PointerType struct {
	to Type

	ast     *ast.StarExpr
	checked *types.Pointer
}

func NewPointerType(to Type) *PointerType {
	return &PointerType{to: to}
}

// IsType marks *PointerType as an implementor of Type for documentation purposes.
func (t *PointerType) IsType() {}

func (t *PointerType) To() Type {
	return t.to
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *PointerType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *PointerType) Checked() *types.Pointer {
	return t.checked
}

type SliceType struct {
	of Type

	ast     *ast.ArrayType
	checked *types.Slice
}

func NewSliceType(of Type) *SliceType {
	return &SliceType{of: of}
}

// IsType marks *SliceType as an implementor of Type for documentation purposes.
func (t *SliceType) IsType() {}

func (t *SliceType) Of() Type {
	return t.of
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *SliceType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *SliceType) Checked() *types.Slice {
	return t.checked
}

type ArrayType struct {
	len uint
	of  Type

	ast     *ast.ArrayType
	checked *types.Array
}

func NewArrayType(len uint, of Type) *ArrayType {
	return &ArrayType{len: len, of: of}
}

// IsType marks *ArrayType as an implementor of Type for documentation purposes.
func (t *ArrayType) IsType() {}

func (t *ArrayType) Len() uint {
	return t.len
}

func (t *ArrayType) Of() Type {
	return t.of
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *ArrayType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *ArrayType) Checked() *types.Array {
	return t.checked
}

type MapType struct {
	key   Type
	value Type

	ast     *ast.MapType
	checked *types.Map
}

// IsType marks *MapType as an implementor of Type for documentation purposes.
func (t *MapType) IsType() {}

func NewMapType(key, value Type) *MapType {
	return &MapType{key: key, value: value}
}

func (t *MapType) Key() Type {
	return t.key
}

func (t *MapType) Value() Type {
	return t.value
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *MapType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *MapType) Checked() *types.Map {
	return t.checked
}

type ChanType struct {
	of  Type
	dir ChanDir

	ast     *ast.ChanType
	checked *types.Chan
}

// IsType marks *ChanType as an implementor of Type for documentation purposes.
func (t *ChanType) IsType() {}

func NewChanType(of Type, dir ChanDir) *ChanType {
	return &ChanType{of: of, dir: dir}
}

func (t *ChanType) Of() Type {
	return t.of
}

func (t *ChanType) Dir() ChanDir {
	return t.dir
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *ChanType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *ChanType) Checked() *types.Chan {
	return t.checked
}

type StructType struct {
	fields []AField

	ast     *ast.StructType
	checked *types.Struct
}

func NewStructType(fields ...AField) *StructType {
	return &StructType{fields: fields}
}

// IsType marks *StructType as an implementor of Type for documentation purposes.
func (t *StructType) IsType() {}

func (t *StructType) Fields() []AField {
	return t.fields
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *StructType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *StructType) Checked() *types.Struct {
	return t.checked
}

type Arg interface {
	IsArg()
	Type() Type
	AST() *ast.Field
	Checked() *types.Var
}

type AField interface {
	IsField()
	Name() string
	Type() Type
	AST() *ast.Field
	Checked() *types.Var
}

type Field struct {
	name string
	typ  Type

	ast     *ast.Field
	checked *types.Var
}

func NewField(name string, typ Type) *Field {
	return &Field{name: name, typ: typ}
}

// IsArg marks *Field as an implementor of Arg for documentation purposes.
func (d *Field) IsArg() {}

// IsField marks *Field as an implementor of AField for documentation purposes.
func (d *Field) IsField() {}

func (d *Field) Name() string {
	return d.name
}

func (d *Field) Type() Type {
	return d.typ
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (d *Field) AST() *ast.Field {
	return d.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (d *Field) Checked() *types.Var {
	return d.checked
}

type EmbeddedField struct {
	typ Type

	ast     *ast.Field
	checked *types.Var
}

// IsField marks *EmbeddedField as an implementor of AField for documentation purposes.
func (d *EmbeddedField) IsField() {}

func NewEmbeddedField(typ Type) *EmbeddedField {
	valid := false

	switch v := typ.(type) {
	case *NamedType, *ImportedType:
		valid = true
	case *PointerType:
		switch v.To().(type) {
		case *NamedType, *ImportedType:
			valid = true
		}
	}

	if valid {
		return &EmbeddedField{typ: typ}
	}

	panic("can only embed named, imported or pointer to those.")
}

func (d *EmbeddedField) Name() string {
	switch v := d.typ.(type) {
	case *NamedType:
		return v.Name()
	case *ImportedType:
		return v.Name()
	case *PointerType:
		switch v := v.To().(type) {
		case *NamedType:
			return v.Name()
		case *ImportedType:
			return v.Name()
		}
	}
	panic("embedded field not named or imported type.")
}

func (d *EmbeddedField) Type() Type {
	return d.typ
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (d *EmbeddedField) AST() *ast.Field {
	return d.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (d *EmbeddedField) Checked() *types.Var {
	return d.checked
}

type TaggedField struct {
	AField
	tag string

	ast     *ast.Field
	checked *types.Var
}

func NewTaggedField(field AField, tag string) *TaggedField {
	return &TaggedField{AField: field, tag: tag}
}

// IsField marks *TaggedField as an implementor of AField for documentation purposes.
func (d *TaggedField) IsField() {}

func (d *TaggedField) Tag() string {
	return d.tag
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (d *TaggedField) AST() *ast.Field {
	return d.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (d *TaggedField) Checked() *types.Var {
	return d.checked
}

type SignatureType struct {
	args []Arg
	ret  []Arg

	ast     *ast.FuncType
	checked *types.Func
}

// IsType marks *SignatureType as an implementor of Type for documentation purposes.
func (t *SignatureType) IsType() {}

func NewSignatureType(args Args, ret Returns) *SignatureType {
	return &SignatureType{args: []Arg(args), ret: []Arg(ret)}
}

func (t *SignatureType) Args() []Arg {
	return t.args
}

func (t *SignatureType) Returns() []Arg {
	return t.ret
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *SignatureType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *SignatureType) Checked() *types.Func {
	return t.checked
}

type Args []Arg
type Returns []Arg

type TypeArg struct {
	typ Type

	ast     *ast.Field
	checked *types.Var
}

func NewTypeArg(typ Type) *TypeArg {
	return &TypeArg{typ: typ}
}

// IsArg marks *TypeArg as an implementor of Arg for documentation purposes.
func (t *TypeArg) IsArg() {}

func (t *TypeArg) Type() Type {
	return t.typ
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *TypeArg) AST() *ast.Field {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *TypeArg) Checked() *types.Var {
	return t.checked
}

type InterfaceType struct {
	fields []IfaceField

	ast     *ast.InterfaceType
	checked *types.Interface
}

func NewInterfaceType(fields ...IfaceField) *InterfaceType {
	return &InterfaceType{fields: fields}
}

// IsType marks *InterfaceType as an implementor of Type for documentation purposes.
func (t *InterfaceType) IsType() {}

func (t *InterfaceType) Fields() []IfaceField {
	return t.fields
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *InterfaceType) AST() ast.Expr {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *InterfaceType) Checked() *types.Interface {
	return t.checked
}

type IfaceField interface {
	IsIfaceField()
	AST() *ast.Field
}

type IfaceMethod struct {
	name string
	args []Arg
	ret  []Arg

	ast     *ast.Field
	checked *types.Func
}

func NewIfaceMethod(name string, args Args, ret Returns) *IfaceMethod {
	return &IfaceMethod{name: name, args: []Arg(args), ret: []Arg(ret)}
}

// IsIfaceField marks *IfaceMethod as an implementor of IfaceField for documentation purposes.
func (t *IfaceMethod) IsIfaceField() {}

func (t *IfaceMethod) Name() string {
	return t.name
}

func (t *IfaceMethod) Args() []Arg {
	return t.args
}

func (t *IfaceMethod) Returns() []Arg {
	return t.ret
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *IfaceMethod) AST() *ast.Field {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *IfaceMethod) Checked() *types.Func {
	return t.checked
}

type EmbeddedIface struct {
	typ Type

	ast     *ast.Field
	checked *types.Interface
}

func NewEmbeddedIface(typ Type) *EmbeddedIface {
	switch typ.(type) {
	case *NamedType, *ImportedType:
		return &EmbeddedIface{typ: typ}
	}
	panic("can only embed named or imported.")
}

// IsIfaceField marks *EmbeddedIface as an implementor of IfaceField for documentation purposes.
func (t *EmbeddedIface) IsIfaceField() {}

func (t *EmbeddedIface) Type() Type {
	return t.typ
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *EmbeddedIface) AST() *ast.Field {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *EmbeddedIface) Checked() *types.Interface {
	return t.checked
}

type AFuncDecl interface {
	IsFuncDecl()
	Name() string
	AST() *ast.FuncDecl
}

// A FuncDecl is a function declaration (func foo(...) ... { ... }).
type FuncDecl struct {
	name string
	args []Arg
	ret  []Arg
	body *Block

	InFile *File

	ast     *ast.FuncDecl
	checked *types.Func
}

// NewFuncDecl creates a function declaration (func foo(...) ... { ... }).
// The arguments and returns can be nil. For functions without a body, see
// NewBodylessFunc.
func NewFuncDecl(name string, args Args, ret Returns, body *Block) *FuncDecl {
	return &FuncDecl{name: name, args: []Arg(args), ret: []Arg(ret), body: body}
}

// IsDecl marks *FuncDecl as an implementor of Decl for documentation purposes.
func (d *FuncDecl) IsDecl() {}

// IsFuncDecl marks *FuncDecl as an implementor of AFuncDecl for documentation purposes.
func (d *FuncDecl) IsFuncDecl() {}

func (t *FuncDecl) Name() string {
	return t.name
}

func (t *FuncDecl) Args() []Arg {
	return t.args
}

func (t *FuncDecl) Returns() []Arg {
	return t.ret
}

func (t *FuncDecl) Body() *Block {
	return t.body
}

// SetInFile links this object with a file.
func (d *FuncDecl) SetInFile(f *File) {
	d.InFile = f
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *FuncDecl) AST() *ast.FuncDecl {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *FuncDecl) Checked() *types.Func {
	return t.checked
}

// A BodylessFuncDecl is a function declaration without a body (func foo(...) ...).
type BodylessFuncDecl struct {
	name string
	args []Arg
	ret  []Arg

	InFile *File

	ast     *ast.FuncDecl
	checked *types.Func
}

// NewBodylessFuncDecl creates a function declaration without a body (func foo(...) ...).
// The arguments and returns can be nil.
func NewBodylessFuncDecl(name string, args Args, ret Returns) *BodylessFuncDecl {
	return &BodylessFuncDecl{name: name, args: []Arg(args), ret: []Arg(ret)}
}

// IsDecl marks *BodylessFuncDecl as an implementor of Decl for documentation purposes.
func (d *BodylessFuncDecl) IsDecl() {}

// IsFuncDecl marks *FuncDecl as an implementor of AFuncDecl for documentation purposes.
func (d *BodylessFuncDecl) IsFuncDecl() {}

func (t *BodylessFuncDecl) Name() string {
	return t.name
}

func (t *BodylessFuncDecl) Args() []Arg {
	return t.args
}

func (t *BodylessFuncDecl) Returns() []Arg {
	return t.ret
}

// SetInFile links this object with a file.
func (d *BodylessFuncDecl) SetInFile(f *File) {
	d.InFile = f
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *BodylessFuncDecl) AST() *ast.FuncDecl {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *BodylessFuncDecl) Checked() *types.Func {
	return t.checked
}

type MethodDecl struct {
	AFuncDecl
	recv *Field

	InFile *File

	ast     *ast.FuncDecl
	checked *types.Func
}

// NewMethodDecl creates a method declaration (function (x X) foo(...) ... { ... }).
func NewMethodDecl(receiver *Field, funcDecl AFuncDecl) *MethodDecl {
	// TODO: receiver shouldn't be Field.
	return &MethodDecl{AFuncDecl: funcDecl, recv: receiver}
}

// IsDecl marks *MethodDecl as an implementor of Decl for documentation purposes.
func (d *MethodDecl) IsDecl() {}

func (t *MethodDecl) Receiver() *Field {
	return t.recv
}

// SetInFile links this object with a file.
func (d *MethodDecl) SetInFile(f *File) {
	d.InFile = f
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *MethodDecl) AST() *ast.FuncDecl {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *MethodDecl) Checked() *types.Func {
	return t.checked
}

// A Block is a series of statements between curly braces ({ ... }).
type Block struct {
	stmts []Stmt

	ast     *ast.BlockStmt
	checked *types.Scope
}

// NewBlock creates a new block ({ ... }).
func NewBlock(stmts ...Stmt) *Block {
	return &Block{stmts: stmts}
}

// IsStmt marks *Block as an implementor of Stmt for documentation purposes.
func (t *Block) IsStmt() {}

// Stmts returns the statements in this block.
func (t *Block) Stmts() []Stmt {
	return t.stmts
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *Block) AST() ast.Stmt {
	return t.ast
}

// Checked returns the corresponding representation of this type as a
// golang.org/x/tools/go/types data type.
func (t *Block) Checked() *types.Scope {
	return t.checked
}

// A Stmt is a statement, one of the sequential elements that conform blocks
// (such as function bodies) and are to be executed one after another.
type Stmt interface {
	IsStmt()
}

type Vars []string

// A Definition is a statement which declares and initializes one or more
// variables to a single expression (x, y, z := ...). For multiple expressions,
// use multiple definitions.
type Definition struct {
	vars  Vars
	value Expr

	ast *ast.AssignStmt
	// TODO: checked
}

// NewDefinition creates a definition (x, y, z := ...).
func NewDefinition(vars Vars, value Expr) *Definition {
	return &Definition{vars: vars, value: value}
}

// IsStmt marks *Definition as an implementor of Stmt for documentation purposes.
func (d *Definition) IsStmt() {}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *Definition) AST() ast.Stmt {
	return t.ast
}

// An Assignment is a statement which sets the value of one or more
// variables to a single expression (x, y, z = ...). For multiple expressions,
// use multiple assignments.
type Assignment struct {
	vars  Vars
	value Expr

	ast *ast.AssignStmt
	// TODO: checked
}

// NewAssignment creates an assignment (x, y, z = ...).
func NewAssignment(vars Vars, value Expr) *Assignment {
	return &Assignment{vars: vars, value: value}
}

// IsStmt marks *Assignment as an implementor of Stmt for documentation purposes.
func (d *Assignment) IsStmt() {}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *Assignment) AST() ast.Stmt {
	return t.ast
}

// A label wraps an statement so it can be referred from break, continue, and
// goto. (foo: ...)
type Label struct {
	label   string
	labeled Stmt

	ast *ast.LabeledStmt
}

// NewLabel creates a labeled statement (foo: ...).
func NewLabel(label string, labeled Stmt) *Label {
	return &Label{label: label, labeled: labeled}
}

// IsStmt marks *Assignment as an implementor of Stmt for documentation purposes.
func (d *Label) IsStmt() {}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *Label) AST() ast.Stmt {
	return t.ast
}

// A Expr is an expression, an element that can be evaluated to retrieve a
// value.
type Expr interface {
	IsExpr()
}

// A IntLiteral is a hardcoded integer value (123).
type IntLiteral struct {
	value int

	ast *ast.BasicLit
}

// NewIntLiteral creates an integer literal (123).
func NewIntLiteral(value int) *IntLiteral {
	return &IntLiteral{value: value}
}

// IsExpr marks *IntLiteral as an implementor of Expr for documentation purposes.
func (d *IntLiteral) IsExpr() {}

// Value returns the value of the literal.
func (d *IntLiteral) Value() int {
	return d.value
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *IntLiteral) AST() ast.Expr {
	return t.ast
}

// A DecimalLiteral is a hardcoded decimal value (123.456).
type DecimalLiteral struct {
	value float64

	ast *ast.BasicLit
}

// NewDecimalLiteral creates a decimal literal (123.456).
func NewDecimalLiteral(value float64) *DecimalLiteral {
	return &DecimalLiteral{value: value}
}

// IsExpr marks *DecimalLiteral as an implementor of Expr for documentation purposes.
func (d *DecimalLiteral) IsExpr() {}

// Value returns the value of the literal.
func (d *DecimalLiteral) Value() float64 {
	return d.value
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *DecimalLiteral) AST() ast.Expr {
	return t.ast
}

// A ComplexLiteral is a hardcoded complex value (123.0 + 456.7i).
type ComplexLiteral struct {
	value complex128

	ast *ast.BinaryExpr
}

// NewComplexLiteral creates a complex literal (123.0 + 456.7i).
func NewComplexLiteral(value complex128) *ComplexLiteral {
	return &ComplexLiteral{value: value}
}

// IsExpr marks *ComplexLiteral as an implementor of Expr for documentation purposes.
func (d *ComplexLiteral) IsExpr() {}

// Value returns the value of the literal.
func (d *ComplexLiteral) Value() complex128 {
	return d.value
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *ComplexLiteral) AST() ast.Expr {
	return t.ast
}

// A StringLiteral is a hardcoded string value in double quotes ("foo").
// For backtick-wrapped strings (`foo`), use RawStringLiteral.
type StringLiteral struct {
	value string

	ast *ast.BasicLit
}

// NewStringLiteral creates a string value in double quotes ("foo").
// For backtick-wrapped strings (`foo`), use NewRawStringLiteral.
func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{value: value}
}

// IsExpr marks *StringLiteral as an implementor of Expr for documentation purposes.
func (d *StringLiteral) IsExpr() {}

// Value returns the value of the literal.
func (d *StringLiteral) Value() string {
	return d.value
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *StringLiteral) AST() ast.Expr {
	return t.ast
}

// A RawStringLiteral is a hardcoded string value in backticks (`foo`).
// For strings in double quotes ("foo"), use StringLiteral.
type RawStringLiteral struct {
	value string

	ast *ast.BasicLit
}

// NewRawStringLiteral creates a string value in backticks (`foo`).
// For strings in double quotes ("foo"), use NewStringLiteral.
func NewRawStringLiteral(value string) *RawStringLiteral {
	return &RawStringLiteral{value: value}
}

// IsExpr marks *RawStringLiteral as an implementor of Expr for documentation purposes.
func (d *RawStringLiteral) IsExpr() {}

// Value returns the value of the literal.
func (d *RawStringLiteral) Value() string {
	return d.value
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *RawStringLiteral) AST() ast.Expr {
	return t.ast
}

// A Name is an identifier for a value (foo).
type Name struct {
	name string

	ast *ast.Ident
}

// NewName creates a name (foo).
func NewName(name string) *Name {
	return &Name{name: name}
}

// IsExpr marks *Name as an implementor of Expr for documentation purposes.
func (d *Name) IsExpr() {}

// Name returns the name of the identified thing.
func (d *Name) Name() string {
	return d.name
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *Name) AST() ast.Expr {
	return t.ast
}

// A Closure is an anonymous function (func(...) ... { ... }).
type Closure struct {
	args []Arg
	ret  []Arg
	body *Block

	ast *ast.FuncLit
}

// NewClosure creates an anonymous function (func(...) ... { ... }).
// The arguments and returns can be nil.
func NewClosure(args Args, ret Returns, body *Block) *Closure {
	return &Closure{args: []Arg(args), ret: []Arg(ret), body: body}
}

// IsExpr marks *Closure as an implementor of Expr for documentation purposes.
func (d *Closure) IsExpr() {}

func (t *Closure) Args() []Arg {
	return t.args
}

func (t *Closure) Returns() []Arg {
	return t.ret
}

func (t *Closure) Body() *Block {
	return t.body
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *Closure) AST() ast.Expr {
	return t.ast
}

// A Call is a call to a function, method, or type (foo(...)).
type Call struct {
	callee Expr
	args   []Expr

	ast *ast.CallExpr
}

// NewCall creates a call to function, method, or type (foo(...)).
// The arguments can be nil.
func NewCall(callee Expr, args ...Expr) *Call {
	return &Call{callee: callee, args: args}
}

// IsExpr marks *Call as an implementor of Expr for documentation purposes.
func (d *Call) IsExpr() {}

// IsStmt marks *Call as an implementor of Stmt for documentation purposes.
func (d *Call) IsStmt() {}

func (t *Call) Callee() Expr {
	return t.callee
}

func (t *Call) Args() []Expr {
	return t.args
}

// AST returns the corresponding representation of this type as a go/ast
// data type.
func (t *Call) AST() ast.Expr {
	return t.ast
}

var (
	Bool       = &BasicType{name: "bool"}
	Int        = &BasicType{name: "int"}
	Int8       = &BasicType{name: "int8"}
	Int16      = &BasicType{name: "int16"}
	Int32      = &BasicType{name: "int32"}
	Int64      = &BasicType{name: "int64"}
	Uint       = &BasicType{name: "uint"}
	Uint8      = &BasicType{name: "uint8"}
	Uint16     = &BasicType{name: "uint16"}
	Uint32     = &BasicType{name: "uint32"}
	Uint64     = &BasicType{name: "uint64"}
	Uintptr    = &BasicType{name: "uintptr"}
	Float32    = &BasicType{name: "float32"}
	Float64    = &BasicType{name: "float64"}
	Complex64  = &BasicType{name: "complex64"}
	Complex128 = &BasicType{name: "complex128"}
	Byte       = &BasicType{name: "byte"}
	Rune       = &BasicType{name: "rune"}
	String     = &BasicType{name: "string"}
)

type ChanDir int

const (
	Both ChanDir = iota
	Out
	In
)
