# meta

**Work in progress.**

Package meta works over `go/{ast,types}` to make code generation friendly.

Featuring:

* A **DSL**-ish way of expressing Go source code in a **type-safe, programmatic** manner.
* A set of data types that represent **Go source code with a human-friendly interface**.
* **Code generation** from those data structures.
* Binding with **[`go/ast`][go/ast]** and **[`golang.org/x/tools/go/types`][go/types]**.

[go/ast]: http://godoc.org/go/ast
[go/types]: http://godoc.org/golang.org/x/tools/go/types
