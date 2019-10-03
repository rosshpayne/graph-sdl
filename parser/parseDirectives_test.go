package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestMultiDirective1(t *testing.T) {

	input := `
input ExampleInputObject @ june (asdf:234) @ june2 (aesdf:234) @ june3 (as2df:"abc") {
  a: String = "AbcDef" @ ref (if:123) @ jack (sd: "abc") @ june (asdf:234) @ ju (asdf:234) @ judkne (asdf:234) @ junse (asdf:234) @ junqe (asdf:234)  @ june (assdf:234)
  b: Int!@june(asdf:234) @ ju (asdf:234)
}
`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
	// for v, o := range repo {
	// 	fmt.Printf(" %s, %T\n", v, o)
	// }
}
