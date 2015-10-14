// +build ignore

package main

import (
	"github.com/meta-go/meta"
	"github.com/meta-go/meta/g"
)

func main() {
	code := g.Package(
		"pkgpopulate",
		g.File("file1").
			Import("fmt").
			EmbeddedImport("github.com/meta-go/meta").
			DiscardedImport("github.com/tcard/valuegraph").
			QualifiedImport("abc", "github.com/meta-go/meta").
			Var("myVar", g.Int, nil).
			Var("myVar2", nil, g.S("abc")).
			Const("myConst", nil, g.I(123)).
			Type("myType", g.String).
			Type("other", g.NamedType("ptrTo")).
			Type("ptrTo", g.PointerTo(g.Bool)).
			Type("sliceOf", g.SliceOf(g.PointerTo(g.Rune))).
			Type("arrayOf", g.ArrayOf(123, g.PointerTo(g.Rune))).
			Type("myMap", g.Map(g.Int, g.PointerTo(g.Rune))).
			Type("importedTy", g.SliceOf(g.ImportedType("time", "Time"))).
			Type(
			"MyStruct", g.Struct(
				g.EmbeddedField(g.NamedType("OtherStruct")),
				g.Tagged(g.EmbeddedField(g.NamedType("ptrTo")), "tag"),
				g.Tagged(g.EmbeddedField(g.PointerTo(g.NamedType("ptrTo"))), "tag"),
				g.Field("A", g.Chan(g.NamedType("myType"))),
				g.Field("Ain", g.ChanIn(g.NamedType("myType"))),
				g.Field("Aout", g.OutChan(g.NamedType("myType"))),
				g.Tagged(g.Field("B", g.Int), `json:"b,omitempty"`))).
			Type("OtherStruct", g.Struct(g.Field("C", g.Int))).
			Type(
			"funcy", g.Signature(
				g.Args(g.Arg("a", g.Int), g.Arg("b", g.String)),
				g.Returns(g.TypeArg(g.Int), g.TypeArg(g.String)),
			)).
			Type(
			"funcy2", g.Signature(
				g.Args(),
				g.Returns(),
			)).
			Type(
			"myIface", g.Interface(
				g.EmbeddedIface(g.ImportedType("time", "Time")),
				g.EmbeddedIface(g.NamedType("ABC")),
				g.IfaceMethod("Foo", g.Args(g.Arg("a", g.Int), g.Arg("b", g.String)), nil),
				g.IfaceMethod("bar", nil, g.Returns()),
			)).
			BodylessFunc(
			"fibo2",
			g.Args(g.Arg("a", g.Int), g.Arg("b", g.String)),
			g.Returns(g.TypeArg(g.Int), g.TypeArg(g.String))).
			Method(
			g.Field("x", g.NamedType("holis")),
			g.Func(
				"fibo",
				g.Args(g.Arg("a", g.Int), g.Arg("b", g.String)),
				g.Returns(g.TypeArg(g.Int), g.TypeArg(g.String)),
				g.Block(
					g.Var("myVar", g.Int, nil),
					g.Type("funcy2", g.Signature(
						g.Args(),
						g.Returns(),
					)),
					g.Var("myVar2", g.String, g.Closure(
						g.Args(g.Arg("a", g.Int), g.Arg("b", g.String)),
						g.Returns(g.TypeArg(g.Int), g.TypeArg(g.String)),
						g.Block(
							g.Var("myVar", g.Int, nil),
						))),
					g.Define(g.Vars("a", "b"), g.Closure(g.Args(), g.Returns(), g.Block())),
					g.Assign(g.Vars("a"), g.I(123)),
					g.Define(g.Vars("b"), g.S("foo")),
					g.Define(g.Vars("c"), g.RawS("foo")),
					g.Define(g.Vars("c"), g.F(123.456)),
					g.Define(g.Vars("abc"), g.Complex(123.0+456.3i)),
					g.Label("hola", g.Block(
						g.Call(g.N("foo"), g.I(123)),
					)),
				))),
		g.File("file2.go"))

	meta.WriteToDir(code, "pkgpopulate")
}
