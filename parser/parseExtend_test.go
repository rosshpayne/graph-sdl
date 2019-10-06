package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestExtend0(t *testing.T) {

	input := `
input ExampleInputObject {
  a: String
  b: Int!
}


`

	expectedDoc := `input ExampleInputObject {a : String b : Int!  }`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for i, v := range errs {
		fmt.Println(i, v.Error())
	}
	if compare(d.String(), expectedDoc) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedDoc))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestExtend1(t *testing.T) {

	input := `


extend input ExampleInputObject  {
  name: String
  age: [Int!]
}

`

	expectedDoc := `input ExampleInputObject {a : String b : Int! name : String age : [Int!] }`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	if d != nil {
		fmt.Println("++++++ ", d.String())
		for i, v := range errs {
			fmt.Println(i, v.Error())
		}
		if compare(d.String(), expectedDoc) {
			fmt.Println(trimWS(d.String()))
			fmt.Println(trimWS(expectedDoc))
			t.Errorf(`*************  program.String() wrong.`)
		}
	}
}

// func TestExtend1a(t *testing.T) {

// 	input := `

// extend input Address  {
//   name: String
//   age: [Int!]
// }

// `

// 	expectedDoc := `input ExampleInputObject {a : String b : Int! name : String age : [Int!] }`

// 	l := lexer.New(input)
// 	p := New(l)
// 	d, errs := p.ParseDocument()
// 	fmt.Println(d.String())
// 	for i, v := range errs {
// 		fmt.Println(i, v.Error())
// 	}
// 	if compare(d.String(), expectedDoc) {
// 		fmt.Println(trimWS(d.String()))
// 		fmt.Println(trimWS(expectedDoc))
// 		t.Errorf(`*************  program.String() wrong.`)
// 	}
// }

func TestExtend2(t *testing.T) {

	input := `
input ExampleInputObject {
  a: String
  b: Int!
}

extend input ExampleInputObject 

`

	expectedDoc := `input ExampleInputObject {a : String b : Int! name : String age : [Int!] }`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for i, v := range errs {
		fmt.Println(i, v.Error())
	}
	if compare(d.String(), expectedDoc) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedDoc))
		t.Errorf(`*************  program.String() wrong.`)
	}
}
