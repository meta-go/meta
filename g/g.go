package g

import "github.com/meta-go/meta"

func Package(name string, files ...*meta.File) *meta.Package {
	return meta.NewPackage(name, files...)
}

func File(name string, imports ...meta.Import) *meta.File {
	return meta.NewFile(name, imports...)
}

func Import(path string) meta.Import {
	return meta.NewImport(path)
}

func QualifiedImport(alias, path string) meta.Import {
	return meta.NewQualifiedImport(alias, path)
}

func DiscardedImport(path string) meta.Import {
	return meta.NewDiscardedImport(path)
}

func EmbeddedImport(path string) meta.Import {
	return meta.NewEmbeddedImport(path)
}

func Var(name string, typ meta.Type, value meta.Expr) *meta.VarDecl {
	return meta.NewVarDecl(name, typ, value)
}

func Const(name string, typ meta.Type, value meta.Expr) *meta.ConstDecl {
	return meta.NewConstDecl(name, typ, value)
}

func Type(name string, typ meta.Type) *meta.TypeDecl {
	return meta.NewTypeDecl(name, typ)
}

func NamedType(name string) *meta.NamedType {
	return meta.NewNamedType(name)
}

func ImportedType(pkg, name string) *meta.ImportedType {
	return meta.NewImportedType(pkg, name)
}

func PointerTo(to meta.Type) *meta.PointerType {
	return meta.NewPointerType(to)
}

func SliceOf(of meta.Type) *meta.SliceType {
	return meta.NewSliceType(of)
}

func ArrayOf(len uint, of meta.Type) *meta.ArrayType {
	return meta.NewArrayType(len, of)
}

func Map(key, value meta.Type) *meta.MapType {
	return meta.NewMapType(key, value)
}

func Chan(of meta.Type) *meta.ChanType {
	return meta.NewChanType(of, meta.Both)
}

func ChanIn(of meta.Type) *meta.ChanType {
	return meta.NewChanType(of, meta.In)
}

func OutChan(of meta.Type) *meta.ChanType {
	return meta.NewChanType(of, meta.Out)
}

func Struct(fields ...meta.AField) *meta.StructType {
	return meta.NewStructType(fields...)
}

func Field(name string, typ meta.Type) *meta.Field {
	return meta.NewField(name, typ)
}

func EmbeddedField(typ meta.Type) *meta.EmbeddedField {
	return meta.NewEmbeddedField(typ)
}

func Tagged(field meta.AField, tag string) *meta.TaggedField {
	return meta.NewTaggedField(field, tag)
}

func Signature(args meta.Args, ret meta.Returns) *meta.SignatureType {
	return meta.NewSignatureType(args, ret)
}

func Args(args ...meta.Arg) meta.Args {
	return meta.Args(args)
}

func Returns(ret ...meta.Arg) meta.Returns {
	return meta.Returns(ret)
}

func Arg(name string, typ meta.Type) meta.Arg {
	return meta.NewField(name, typ)
}

func TypeArg(typ meta.Type) *meta.TypeArg {
	return meta.NewTypeArg(typ)
}

func Interface(fields ...meta.IfaceField) *meta.InterfaceType {
	return meta.NewInterfaceType(fields...)
}

func IfaceMethod(name string, args meta.Args, ret meta.Returns) *meta.IfaceMethod {
	return meta.NewIfaceMethod(name, args, ret)
}

func EmbeddedIface(typ meta.Type) *meta.EmbeddedIface {
	return meta.NewEmbeddedIface(typ)
}

func Func(name string, args meta.Args, ret meta.Returns, body *meta.Block) *meta.FuncDecl {
	return meta.NewFuncDecl(name, args, ret, body)
}

func BodylessFunc(name string, args meta.Args, ret meta.Returns) *meta.BodylessFuncDecl {
	return meta.NewBodylessFuncDecl(name, args, ret)
}

func Method(receiver *meta.Field, funcDecl meta.AFuncDecl) *meta.MethodDecl {
	return meta.NewMethodDecl(receiver, funcDecl)
}

func Block(stmts ...meta.Stmt) *meta.Block {
	return meta.NewBlock(stmts...)
}

func Define(vars meta.Vars, expr meta.Expr) *meta.Definition {
	return meta.NewDefinition(vars, expr)
}

func Assign(vars meta.Vars, expr meta.Expr) *meta.Assignment {
	return meta.NewAssignment(vars, expr)
}

func Label(label string, labeled meta.Stmt) *meta.Label {
	return meta.NewLabel(label, labeled)
}

func I(i int) *meta.IntLiteral {
	return meta.NewIntLiteral(i)
}

func F(f float64) *meta.DecimalLiteral {
	return meta.NewDecimalLiteral(f)
}

func Complex(c complex128) *meta.ComplexLiteral {
	return meta.NewComplexLiteral(c)
}

func S(s string) *meta.StringLiteral {
	return meta.NewStringLiteral(s)
}

func RawS(s string) *meta.RawStringLiteral {
	return meta.NewRawStringLiteral(s)
}

func N(n string) *meta.Name {
	return meta.NewName(n)
}

func Vars(vars ...string) meta.Vars {
	return meta.Vars(vars)
}

func Closure(args meta.Args, ret meta.Returns, body *meta.Block) *meta.Closure {
	return meta.NewClosure(args, ret, body)
}

func Call(callee meta.Expr, args ...meta.Expr) *meta.Call {
	return meta.NewCall(callee, args...)
}

var (
	Bool       = meta.Bool
	Int        = meta.Int
	Int8       = meta.Int8
	Int16      = meta.Int16
	Int32      = meta.Int32
	Int64      = meta.Int64
	Uint       = meta.Uint
	Uint8      = meta.Uint8
	Uint16     = meta.Uint16
	Uint32     = meta.Uint32
	Uint64     = meta.Uint64
	Uintptr    = meta.Uintptr
	Float32    = meta.Float32
	Float64    = meta.Float64
	Complex64  = meta.Complex64
	Complex128 = meta.Complex128
	Byte       = meta.Byte
	Rune       = meta.Rune
	String     = meta.String
)
