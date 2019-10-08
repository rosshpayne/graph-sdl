package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestFieldArgument1(t *testing.T) {

	input := `type Person {
  name: String!
  age: Int!
  inputX(age: int = 123): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument2(t *testing.T) {

	input := `
enum Direction {
  NORTH
  EAST
  SOUTH
  WEST @deprecated @ dep (if: 99.34 fi:true cat: 23.323)
}
type Person {
  name: String!
  age: Int!
  inputX(age: Direction = EAST): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument3(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: Measure = {height: 123.2 weight: 12}): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument4(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [int!] = [1,2,4 56 345 2342 234 25252 2525223 null]): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument4a(t *testing.T) {

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [[int!]] = [[1,2,4 56] [345 2342 234 25252 2525223 null]]): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument4b(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: String = """abc \ndefasj \nasdf"""): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument4c(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [String] = ["abc defasj asdf" -234.2]): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument5(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [Measure] = [{height: 123.2 weight: 12} {height: 1423.2 weight: 132}]): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument6(t *testing.T) {

	input := `
enum Address {
	NORTH
	SOUTH
	EAST
}

type Person {
  name: String!
  age: Int!
  inputX(info: [Measure] = [{height: 123.2 weight: 12} {height: 1423.2 weight: 132}]): Float
  posts: [Boolean!]!
  addres: Address!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument7(t *testing.T) {

	input := `
	
	type Person2 {
  name: String!
  age: Int!
  inputX(info: [Measure] = [{height: 123.2 weight: 12} {height: 1423.2 weight: 132}]): Float
  posts: [Boolean!]!
  addres: Address!
}
	
	extend type Person2 {
  isHiddenLocally: Boolean
}

extend type Person2 @addedDirective34

`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestFieldArgument7a(t *testing.T) {

	input := `
	

extend type Person2 @addedDirective67

`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestObjFieldCheckType(t *testing.T) {

	input := `
	
input MyInput66 {
  x: Float
  y: Float
}
type Measure66 {
    height: Float
    weight: Int
    form: MyInput66
}
`

	var expectedErr [1]string
	expectedErr[0] = `Field "form" type "MyInput66", is not an output type at line: 10 column: 11` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != 1 {
		t.Errorf(`Expected 1 error "... is not an output type", got none `)
	}
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
