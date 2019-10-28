package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestEnumDuplicate(t *testing.T) {

	input := `

directive @deprecated on ENUM_VALUE | ARGUMENT_DEFINITION
directive @dep on ENUM_VALUE | ARGUMENT_DEFINITION

enum Direction {
  NORTH
  SOUTH
  EAST
  SOUTH
  WEST @deprecated @ dep (if: 99.34 fi:true cat: 23.323)
}
type Person {
  name: String!
  age: Int!
  dir: [Direction]
}
`

	var expectedErr [1]string
	expectedErr[0] = `Duplicate Enum Value [SOUTH] at line: 10 column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
		if len(ex) == 0 {
			break
		}
		found := false
		for _, err := range errs {
			if trimWS(err.Error()) == trimWS(ex) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Expected Error = [%q]`, ex)
		}
	}
	for _, got := range errs {
		found := false
		for _, exp := range expectedErr {
			if trimWS(got.Error()) == trimWS(exp) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Unexpected Error = [%q]`, got.Error())
		}
	}

}

func TestEnumString(t *testing.T) {

	// note: "SOUTH" will be interpreted as a comment, even though the intention is a enum value.
	//  THis will not generate an errror as a result. No fix as its allowed by the spec (I believe)

	input := `
enum Direction {
  NORTH
  "SOUTH"
  EAST
  WEST @deprecated @ dep (if: 99.34 fi:true cat: 23.323)
}
type Person {
  name: String!
  age: Int!
  dir: [Direction]
}
`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != 0 {
		t.Errorf(`***  Expected 0 error got %d.`, len(errs))
	}
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	// if compare(d.String(), input) {
	// 	fmt.Println(trimWS(d.String()))
	// 	fmt.Println(trimWS(input))
	// 	t.Errorf(`***  program.String() wrong.`)
	// }

}

func TestEnumNumber(t *testing.T) {

	// note: "SOUTH" will be interpreted as a comment, even though the intention is a enum value.
	//  THis will not generate an errror as a result. No fix as its allowed by the spec (I believe)

	input := `
enum Direction {
  NORTH
  23
  EAST
  WEST @deprecated @ dep (if: 99.34 fi:true cat: 23.323)
}
type Person {
  name: String!
  age: Int!
  dir: [Direction]
}
`

	var expectedErr [1]string
	expectedErr[0] = `Expected name identifer got Int of "23" at line: 4, column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for _, e := range expectedErr {
		found := false
		for _, n := range errs {
			if trimWS(n.Error()) == trimWS(e) {
				found = true
			}
		}
		if !found {
			t.Errorf(`***  Expected %s- not exists.`, e)
		}
	}
}
func TestBadEnumValue1(t *testing.T) {

	input := `
enum Direction {
  NORTH
  null
  SOUTH
  WEST @deprecated @ dep (if: 99.34)
}
`
	var expectedErr [1]string
	expectedErr[0] = `Expected name identifer got null of "null" at line: 4, column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for _, e := range expectedErr {
		found := false
		for _, n := range errs {
			if trimWS(n.Error()) == trimWS(e) {
				found = true
			}
		}
		if !found {
			t.Errorf(`***  Expected %s- not exists.`, e)
		}
	}
}

func TestBadBracket(t *testing.T) {

	input := `
	enum Direction {
  NORTH
  SOUTH
  WEST @deprecated @ dep [if: 99.34)
}
`
	var expectedErr [1]string
	expectedErr[0] = `Expected a ( or } or { instead got [ at line: 5, column: 26`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for _, e := range expectedErr {
		found := false
		for _, n := range errs {
			if trimWS(n.Error()) == trimWS(e) {
				found = true
			}
		}
		if !found {
			t.Errorf(`***  Expected %s- not exists.`, e)
		}
	}
}

func TestBadEnumValue2(t *testing.T) {

	input := `
		directive @ xyz on | ENUM_VALUE
	directive @ dep  on | ENUM_VALUE
	enum Direction {
  NORTH
  true
  SOUTH
  WEST @xyz @ dep (if: 99.34)
}
`
	var expectedErr [1]string
	expectedErr[0] = `Expected name identifer got TRUE of "true" at line: 6, column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
		if len(ex) == 0 {
			break
		}
		found := false
		for _, err := range errs {
			if trimWS(err.Error()) == trimWS(ex) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Expected Error = [%q]`, ex)
		}
	}
	for _, got := range errs {
		found := false
		for _, exp := range expectedErr {
			if trimWS(got.Error()) == trimWS(exp) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Unexpected Error = [%q]`, got.Error())
		}
	}
}

func TestEnumDuplicate2(t *testing.T) {

	input := `
	directive @ xyz on | ENUM_VALUE
	directive @ dep  on | ENUM_VALUE
enum Direction {
  NORTH
  EAST
  SOUTH
  EAST @xyz @ dep (if: 99.34)
}
type Person {
		address: [Direction]
	}
`

	var expectedErr [1]string
	expectedErr[0] = `Duplicate Enum Value [EAST] at line: 8 column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
		if len(ex) == 0 {
			break
		}
		found := false
		for _, err := range errs {
			if trimWS(err.Error()) == trimWS(ex) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Expected Error = [%q]`, ex)
		}
	}
	for _, got := range errs {
		found := false
		for _, exp := range expectedErr {
			if trimWS(got.Error()) == trimWS(exp) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Unexpected Error = [%q]`, got.Error())
		}
	}
}

func TestMultiError1(t *testing.T) {

	input := `
enum Direction {
  NORTH
  true
  SOUTH
  WEST @deprecated @ __dep (if: 99.34)
}
type Person {
		address: [Place]
		name: [String!]!
	}`

	var expectedErr [4]string
	expectedErr[0] = `Expected name identifer got TRUE of "true" at line: 4, column: 3`
	expectedErr[1] = `identifer [__dep] cannot start with two underscores at line: 6, column: 22`
	expectedErr[2] = `Type "Place" does not exist at line: 9 column: 13`
	expectedErr[3] = `Type "@__dep" does not exist at line: 6 column: 22`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	for _, ex := range expectedErr {
		if len(ex) == 0 {
			break
		}
		found := false
		for _, err := range errs {
			if trimWS(err.Error()) == trimWS(ex) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Expected Error = [%q]`, ex)
		}
	}
	for _, got := range errs {
		found := false
		for _, exp := range expectedErr {
			if trimWS(got.Error()) == trimWS(exp) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Unexpected Error = [%q]`, got.Error())
		}
	}
}

func TestEnumValidArgument(t *testing.T) {

	input := `
enum Direction {
  NORTH
  EAST
  SOUTH
  WEST @deprecated @ dep (if: 99.34)
}
type Person {
		address: [String]
		name(arg1: Direction = SOUTH ): Float
	}
`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != 0 {
		t.Errorf(`***  Expected 0 error got %d.`, len(errs))
	}
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
	}
}

func TestEnumValidArgumentWithoutENUM(t *testing.T) {

	input := `

type Person {
		address: [String]
		name(arg1: Direction = SOUTH ): Float
	}
`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != 0 {
		t.Errorf(`***  Expected 0 error got %d.`, len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
	}
}

func TestEnumInvalidArgument(t *testing.T) {

	input := `
	enum Direction {
  NORTH
  EAST
  SOUTH
  WEST @deprecated @ dep (if: 99.34)
}
type Person {
		address: [String]
		name(arg1: Direction = SOUTH33 ): Float
	}
`
	var expectedErr [1]string
	expectedErr[0] = `"SOUTH33" is not a member of type Enum Direction at line: 10 column: 26`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	for _, e := range errs {
		fmt.Printf("Got error: [%s]\n", e.Error())
	}
	for _, e := range expectedErr {
		found := false
		for _, n := range errs {
			if trimWS(n.Error()) == trimWS(e) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Expected: [%s]`, e)
		}
	}
}

func TestEnumInvalidArgumentWithoutENUM(t *testing.T) {

	input := `

type Person {
		address: [String]
		name(arg1: Direction = SOUTH33 ): Float
	}
`
	var expectedErr [1]string
	expectedErr[0] = `"SOUTH33" is not a member of type Enum Direction at line: 5 column: 26`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for _, e := range expectedErr {
		found := false
		for _, n := range errs {
			if trimWS(n.Error()) == trimWS(e) {
				found = true
			}
		}
		if !found {
			t.Errorf(`***  Expected %s- not exists.`, e)
		}
	}
}

func TestEnumValueButArugmentTypeIsNot(t *testing.T) {

	input := `
type Location {
	x: Int
	y: Float
}
type Person {
		address: [String]
		name(arg1: Location = SOUTH ): Float
	}
`
	var expectedErr [2]string
	expectedErr[0] = `Argument "arg1" type "Location", is not an input type at line: 8 column: 14`
	expectedErr[1] = `"SOUTH" is an enum like value but the argument type "Location" is not an Enum type at line: 8 column: 25`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for _, e := range expectedErr {
		found := false
		for _, n := range errs {
			if trimWS(n.Error()) == trimWS(e) {
				found = true
			}
		}
		if !found {
			t.Errorf(`***  Expected %s- not exists.`, e)
		}
	}
}

func TestEnumValueButArugmentTypeIsNot2(t *testing.T) {

	input := `
input Location {
	x: Int
	y: Float
}
type Person {
		address: [String]
		name(arg1: Location = SOUTH ): Float
	}
`
	var expectedErr [1]string
	expectedErr[0] = `"SOUTH" is an enum like value but the argument type "Location" is not an Enum type at line: 8 column: 25`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for _, e := range expectedErr {
		found := false
		for _, n := range errs {
			if trimWS(n.Error()) == trimWS(e) {
				found = true
			}
		}
		if !found {
			t.Errorf(`***  Expected %s- not exists.`, e)
		}
	}
}

func TestEnumArgValueWrongType(t *testing.T) {

	input := `

type Person {
		address: [String]
		name(arg1: Direction = "SOUTH" ): Float
	}
`
	var expectedErr [1]string
	expectedErr[0] = `Required type "enum", got "String" at line: 5 column: 26`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for _, e := range expectedErr {
		found := false
		for _, n := range errs {
			if trimWS(n.Error()) == trimWS(e) {
				found = true
			}
		}
		if !found {
			t.Errorf(`***  Expected %s- not exists.`, e)
		}
	}
}

func TestFieldNoType(t *testing.T) {

	input := `
enum Direction {
  NORTH
  EAST
  SOUTH
  WEST @deprecated @ dep (if: 99.34)
}
type Person {
		address: [String]
		name(arg1: Direction = SOUTH )
		extra: Int
	}
`
	var expectedErr [4]string
	expectedErr[0] = `Colon expected got IDENT of extra at line: 11, column: 3`
	expectedErr[1] = `Expected name identifer got : of ":" at line: 11, column: 8`
	expectedErr[2] = `Colon expected got Int of Int at line: 11, column: 10`
	expectedErr[3] = `Type "extra" does not exist at line: 11 column: 3`
	//	expectedErr[4] = `Type "extra" does not exist at line: 11 column: 3`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for _, e := range expectedErr {
		found := false
		for _, n := range errs {
			if trimWS(n.Error()) == trimWS(e) {
				found = true
			}
		}
		if !found {
			t.Errorf(`***  Expected %s- not exists.`, e)
		}
	}

}
