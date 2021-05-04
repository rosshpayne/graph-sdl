package parser

import (
	"fmt"
	"testing"

	"github.com/graphql/internal/graph-sdl/lexer"
)

func TestEnumValidx(t *testing.T) {

	input := `

directive @deprecated on ENUM_VALUE | ARGUMENT_DEFINITION
directive @dep (if : Int) on ENUM_VALUE | ARGUMENT_DEFINITION

enum Direction {
  NORTH
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

	var expectedErr []string = []string{
		`Required type for argument "if" is Int, got Float at line: 10 column: 27`,
		`Argument "fi" is not a valid name for directive "@dep" at line: 10 column: 37`,
		`Argument "cat" is not a valid name for directive "@dep" at line: 10 column: 45`,
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, v := range errs {
		fmt.Println("errors: ", v)
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

func TestEnumDuplicate(t *testing.T) {

	input := `
	directive @ dep  on | ENUM_VALUE

enum Direction {
  NORTH
  SOUTH
  EAST
  SOUTH
  WEST 
}
type Person {
  name: String!
  age: Int!
  dir: [Direction]
}
`

	var expectedErr [1]string
	expectedErr[0] = `Duplicate Enum Value [SOUTH] at line: 8 column: 3`

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
  WEST @deprecated 
}
type Person {
  name: String!
  age: Int!
  dir: [Direction]
}
`

	var expectedErr []string

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

func TestEnumNonExistentDirectiveRef(t *testing.T) {

	input := `
enum Direction {
  NORTH
  SOUTH
  EAST
  WEST @deprecated @ dep2 @dep3 (if: 99.34 fi:true cat: 23.323)
}
type Person {
  name: String!
  age: Int!
  dir: [Direction]
}
`

	var expectedErr []string = []string{
		`"@dep2" does not exist in document "DefaultDoc" at line: 6 column: 22`,
		`"@dep3" does not exist in document "DefaultDoc" at line: 6 column: 28`,
	}
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// for _, v := range errs {
	// 	fmt.Println("errors: ", v)
	// }

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

func TestEnumNumber(t *testing.T) {

	// note: "SOUTH" will be interpreted as a comment, even though the intention is a enum value.
	//  THis will not generate an errror as a result. No fix as its allowed by the spec (I believe)

	input := `
enum Direction {
  NORTH
  23
  EAST
  WEST @deprecated @ dep
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
	// for _, v := range errs {
	// 	fmt.Println("errors: ", v)
	// }

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

func TestEnumNullEnumValue(t *testing.T) {

	input := `
enum Direction {
  NORTH
  null
  SOUTH
  WEST @deprecated @ dep 
}
`
	expectedErr := []string{`Expected name identifer got Null of "null" at line: 4, column: 3`}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, v := range errs {
		fmt.Println("errors: ", v)
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

func TestEnumBadBracket(t *testing.T) {

	input := `
	enum Direction {
  NORTH
  SOUTH
  WEST @deprecated @ dep [if:true)
}
`

	var expectedErr = []string{
		`Expected a ( or } or { instead got [ at line: 5, column: 26`,
		`Argument "if" is not a valid name for directive "@dep" at line: 5 column: 27`,
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// for _, v := range errs {
	// 	fmt.Println("errors: ", v)
	// }

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
func TestEnumBadBoolEnumValue(t *testing.T) {

	input := `
		directive @ xyz on | ENUM_VALUE
	directive @ dep  on | ENUM_VALUE
	enum Direction {
  NORTH
  true
  SOUTH
  WEST @xyz @ dep 
}
`
	var expectedErr = []string{`Expected name identifer got TRUE of "true" at line: 6, column: 3`}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// for _, v := range errs {
	// 	fmt.Println("errors: ", v)
	// }

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
  EAST @xyz @ dep 
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
	// for _, v := range errs {
	// 	fmt.Println("*** errors: ", v)
	// }
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

func TestEnumMultiError1(t *testing.T) {

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
	expectedErr[1] = `identifer "__dep" cannot start with two underscores at line: 6, column: 22`
	expectedErr[2] = `"Place" does not exist in document "DefaultDoc" at line: 9 column: 13`
	expectedErr[3] = `"@__dep" does not exist in document "DefaultDoc" at line: 6 column: 22`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// for _, v := range errs {
	// 	fmt.Println(v.Error())
	// }
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

func TestEnumNotAMember(t *testing.T) {

	input := `
enum Direction {
  NORTH
  EAST
  SOUTH
  WEST @deprecated @ dep (if: 99.34)
}
type Person {
		address: [String]
		name(arg1: Direction = XYZ ): Float
	}
`

	expectedErr := []string{
		`Argument "if" is not a valid name for directive "@dep" at line: 6 column: 27`,
		`"XYZ" is not a member of Enum type Direction at line: 10 column: 26`,
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// for _, v := range errs {
	// 	fmt.Println("errors: ", v)
	// }
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

func TestEnumNotMemberDbDef(t *testing.T) {

	input := `

type Person {
		address: [String]
		name(arg1: Direction = XYZ ): Float
	}
`

	expectedErr := []string{
		`"XYZ" is not a member of Enum type Direction at line: 5 column: 26`,
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// for _, v := range errs {
	// 	fmt.Println("errors: ", v)
	// }
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

func TestEnumValueButArgmentTypeIsNot(t *testing.T) {

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
	// for _, v := range errs {
	// 	fmt.Println("*** errors: ", v)
	// }
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

func TestEnumArgValueWrongType(t *testing.T) {

	input := `

type Person {
		address: [String]
		name(arg1: Direction = "SOUTH" ): Float
	}
`
	expectedErr := []string{`Required type for argument "arg1" is Enum, got String at line: 5 column: 8`}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// for _, v := range errs {
	// 	fmt.Println("*** errors: ", v)
	// }
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

func TestEnumFieldNoType(t *testing.T) {

	input := `
enum Direction {
  NORTH
  EAST
  SOUTH
  WEST @deprecated @ dep (if: 99.34)
}
type Person {
		address: [String]
		namxe(arg1: Direction = SOUTH ) 
		extra: Int
	}
`
	// var expectedErr = []string{
	// 	`Colon expected got IDENT of extra at line: 11, column: 3`,
	// 	`Expected name identifer got : of ":" at line: 11, column: 8`,
	// 	`Colon expected got Int of Int at line: 11, column: 10`,
	// 	`"extra" does not exist in document "DefaultDoc" at line: 11 column: 3`,
	// 	`Argument "if" is not a valid name for directive "@dep" at line: 6 column: 27`,
	// }

	var expectedErr = []string{
		`Expected a colon followed by a GQL-Type, got "extra"  at line: 11, column: 3`,
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// for _, v := range errs {
	// 	fmt.Println("*** errors: ", v)
	// }
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
