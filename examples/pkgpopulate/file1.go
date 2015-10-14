package pkgpopulate

import "fmt"
import . "github.com/meta-go/meta"
import _ "github.com/tcard/valuegraph"
import abc "github.com/meta-go/meta"

var myVar int
var myVar2 = "abc"

const myConst = 123

type myType string
type other ptrTo
type ptrTo *bool
type sliceOf []*rune
type arrayOf [123]*rune
type myMap map[int]*rune
type importedTy []time.Time
type MyStruct struct {
	OtherStruct
	ptrTo	`tag`
	*ptrTo	`tag`
	A	chan myType
	Ain	chan<- myType
	Aout	<-chan myType
	B	int	`json:"b,omitempty"`
}
type OtherStruct struct {
	C int
}
type funcy func(a int, b string) (int, string)
type funcy2 func()
type myIface interface {
	time.Time
	ABC
	Foo(a int, b string)
	bar()
}

func fibo2(a int, b string) (int, string)
func (x holis) fibo(a int, b string) (int, string) {
	var myVar int
	type funcy2 func()
	var myVar2 string = func(a int, b string) (int, string) {
		var myVar int
	}
	a, b := func() {
	}
	a = 123
	b := "foo"
	c := `foo`
	c := 123.456
	abc := 123 + 456.3i
hola:
	{
		foo(123)
	}
}
