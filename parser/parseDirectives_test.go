package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestMultiDirective1(t *testing.T) {

	input := `
input ExampleInputObjectDirective @ june (asdf:234) @ june2 (aesdf:234) @ june3 (as2df:"abc") {
  a: String = "AbcDef" @ ref (if:123) @ jack (sd: "abc") @ june (asdf:234) @ ju (asdf:234) @ judkne (asdf:234) @ junse (asdf:234) @ junqe (asdf:234) 
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

func TestInputDoesnotExist(t *testing.T) {

	input := `
extend input ExampleInputXYZ @ june (asdf:234) 
`
	var expectedErr [1]string
	expectedErr[0] = `Type "ExampleInputXYZ" does not exist at line: 2 column: 14`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	for _, v := range errs {
		fmt.Println("Err: ", v.Error())
	}
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be %d errors got %d`, len(expectedErr), len(errs)))
	}
	fmt.Println("outut ", d.String())
	for i, v := range errs {
		if i < len(expectedErr) {
			if v.Error() != expectedErr[i] {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}

func TestExtendInpDirDuplicate(t *testing.T) {

	input := `
extend input ExampleInputObjectDirective @ june (asdf:234) 
`
	var expectedErr [2]string
	expectedErr[0] = `Duplicate Directive name "june" at line: 2, column: 44`
	expectedErr[1] = `extend for type "ExampleInputObjectDirective" contains no changes at line: 0, column: 0`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be %d errors got %d`, len(expectedErr), len(errs)))
	}
	//fmt.Println(d.String())
	for i, v := range errs {
		if i < len(expectedErr) {
			if v.Error() != expectedErr[i] {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}
