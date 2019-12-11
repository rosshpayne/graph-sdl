package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/ast"
	"github.com/graph-sdl/lexer"
)

func TestBadTypeName(t *testing.T) {

	input := `
type Person {
  name: Str+ing!
  address: String
  Altitude: Float
}
`
	var expectedErr [1]string
	expectedErr[0] = `Expected name identifer got ILLEGAL of "+" at line: 3, column: 12`
	//expectedErr[1] = `Type "Str" does not exist at line: 3 column: 9`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgument1(t *testing.T) {

	input := `type Person {
  name: String!
  age: Int!
  inputX(age: Int = 123): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldInvalidDT(t *testing.T) {

	input := `type Person {
  name: String!
  age: Int!
  inputX(age: int = 123): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Type "int" does not exist at line: 4 column: 15` //

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestCheckInputValueType0(t *testing.T) {

	input := `type Person88 {
  name: String!
  age: Int!
  inputX(age:Int = 123): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}

}

func TestFieldArgTypeNotFound(t *testing.T) {

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [[int!]] = [[1,2 ,4 56] [345 2342 234 25252 2525223 null]]): Float
  posts: [Boolean!]!
}`
	var expectedErr [1]string
	expectedErr[0] = `Type "int" does not exist at line: 6 column: 18`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestCheckInputValueType1(t *testing.T) {

	input := `type Person88 {
  name: String!
  age: Int!
  inputX(age:Int = 123.4): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Required type "Int", got "Float" at line: 4 column: 20` //

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestCheckInputValueType2(t *testing.T) {

	input := `type Person88 {
  name: String!
  age: Int!
  inputX(age:String = 123): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Required type "String", got "Int" at line: 4 column: 23` //

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestCheckInputValueType3(t *testing.T) {

	input := `type Person88 {
  name: String!
  age: Int!
  inputX(age:[String] = ["abc","def"  4 ]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Required type "String", got "Int" at line: 4 column: 39` //

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestCheckInputValueType4(t *testing.T) {

	input := `type Person88 {
  name: String!
  age: Int!
  inputX(age:[[String]] = ["xyss" "cat" ["abc","def" "Hij"] "xyz"]): Float
  posts: [Boolean!]!
}`

	var expectedErr [3]string
	expectedErr[0] = `Value "xyss" is not at required nesting of 2 at line: 4 column: 28`
	expectedErr[1] = `Value "cat" is not at required nesting of 2 at line: 4 column: 35`
	expectedErr[2] = `Value "xyz" is not at required nesting of 2 at line: 4 column: 61`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgument3a(t *testing.T) {

	input := `
input Measure {
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
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgument3b(t *testing.T) {
	//TODO : I believe an OBject should not be used as an ARgument. Should be input. As in 3a above. THIS SHOUJLD FAIL..but currently accepts OBJECTS in argument types.
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
	var expectedErr [1]string
	expectedErr[0] = `Argument "info" type "Measure", is not an input type at line: 9 column: 16`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgument4(t *testing.T) {

	input := `

type Person40 {
  name: String!
  age: Int!
  inputX(info: [Int!] = [1,2,4 56 345 2342 234 25252 2525223 null]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `List cannot contain NULLs at line: 6 column: 62`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
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
  inputX(info: String = ["""abc \ndefasj \nasdf"""]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Argument "info", type is not a list but default value is a list at line: 9 column: 51`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgListNonList(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [String] = """abc \ndefasj \nasdf"""): Float
  posts: [Boolean!]!
}`

	// string input coerced to [String]
	expectedDoc := `typeMeasure{height:Floatweight:Int}typePerson{name:String!age:Int!inputX(info:[String]=["""abc\ndefasj\nasdf"""]):Floatposts:[Boolean!]!}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListNonListNull(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [String] = null ): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgx2ListDiffDepth(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[Int]] = [ [1 2] 3] ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Value 3 is not at required nesting of 2 at line: 9 column: 34`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgx3ListDiffDepth(t *testing.T) {

	input := `
input Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[Measure]] = [ [{height: 2.3 weight: 3} {height: 2.7 weight: 8}] {height: 22.4 weight: 7} ] ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Value {height:22.4 weight:7 }  is not at required nesting of 2 at line: 9 column: 105`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgListDiffDepthInt(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[Int]] =  2 ): Float
  posts: [Boolean!]!
}`

	expectedDoc := `typeMeasure{height:Floatweight:Int}typePerson{name:String!age:Int!inputX(info:[[Int]]=[[2]]):Floatposts:[Boolean!]!}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListDiffDepthInt2(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [Int] =  2 ): Float
  posts: [Boolean!]!
}`

	expectedDoc := `typeMeasure{height:Floatweight:Int}typePerson{name:String!age:Int!inputX(info:[Int]=[2]):Floatposts:[Boolean!]!}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListInt3_1(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [Int] =  [1 2 3] ): Float
  posts: [Boolean!]!
}`

	expectedDoc := `typePerson{name:String!age:Int!inputX(info:[Int]=[1 2 3]):Floatposts:[Boolean!]!}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListInt3_2(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [Int] = null ): Float
  posts: [Boolean!]!
}`

	expectedDoc := `typePerson{name:String!age:Int!inputX(info:[Int]=null):Floatposts:[Boolean!]!}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListInt3_3(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [Int] = [1, 2, null] ): Float
  posts: [Boolean!]!
}`

	expectedDoc := `typePerson{name:String!age:Int!inputX(info:[Int]=[1,2,null]):Floatposts:[Boolean!]!}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

// 3_4 [1, 2, Error]

func TestFieldArgListInt3_5(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [Int]! = [1, 2, 3, 4] ): Float
  posts: [Boolean!]!
}`

	expectedDoc := `typePerson{name:String!age:Int!inputX(info:[Int]!=[1,2,3,4]):Floatposts:[Boolean!]!}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListInt3_6(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [Int]! = null ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Value cannot be NULL at line: 6 column: 25`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgListInt3_7(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [Int]! = [1, 2, null] ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println("[", d.String(), "]")
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

	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

//TODO 3_8 [Int]!	[1, 2, Error]

func TestFieldArgListInt3_9(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [Int!] = [1, 2, 3] ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println("[", d.String(), "]")
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

	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListInt3_10(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [Int!] = null ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println("[", d.String(), "]")
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

	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListInt3_11(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [Int!] = [1, 2, null] ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `List cannot contain NULLs at line: 6 column: 32`

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

//TODO [Int!]	[1, 2, Error]

func TestFieldArgListInt3_13(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info:[Int!]! =	[1, 2, 3] ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()

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
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListInt3_14(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info:[Int!]! = null ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Value cannot be NULL at line: 6 column: 25`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()

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
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListInt3_15(t *testing.T) { // first entry in table in spec 3.12.1

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info:[Int!]! =	[1, 2, null]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `List cannot contain NULLs at line: 6 column: 32`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()

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
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListDiffDepthNull(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[Int]] =  null ): Float
  posts: [Boolean!]!
}`

	expectedDoc := `typeMeasure{height:Floatweight:Int}typePerson{name:String!age:Int!inputX(info:[[Int]]=null):Floatposts:[Boolean!]!}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgListInvalidMember(t *testing.T) {

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

	var expectedErr [1]string
	expectedErr[0] = `Required type "String", got "Float" at line: 9 column: 46`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgListInvalidMember2(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[String]] = [["abc"] ["defasj" "asdf" "asdf" 234 ] ["abc"]]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Required type "String", got "Int" at line: 8 column: 62`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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
	for _, e := range errs {
		fmt.Println("Error: ", e.Error())
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

func TestFieldArgumentInvalidlist(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[String]] = [["abc"] ["defasj" "asdf" "asdf" 234 ] ["abc"]]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Required type "String", got "Int" at line: 8 column: 62`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgumentNullCheck(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[String]!]! = [["abc"] ["defasj" "asdf" "asdf" null ] ["abc" null] null ]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `List cannot contain NULLs at line: 8 column: 84`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgumentNullCheckValid(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[String]]! = [["abc"] ["defasj" "asdf" "asdf" null ] ["abc" null] null ]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = ``

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

func TestFieldArgumentNullCheck2(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [Int!] =  null ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = ``

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

func TestFieldArgumentNullCheck2a(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [Int!]! =  null ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Value cannot be NULL at line: 8 column: 27`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgumentNullCheckInt(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[String]!]! = [["abc"] ["defasj" "asdf" "asdf" ] ["abc" null] "XYZ" ]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Value "XYZ" is not at required nesting of 2 at line: 8 column: 79`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgumentNullCheck3a(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[String]]! = [["abc"] ["defasj" "asdf" "asdf" null ] ["abc" null] null ]): Float
  posts: [Boolean!]!
}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgumentNullCheck3(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[String!]]! = [["abc"] ["defasj" "asdf" "asdf" null ] ["abc" null] null ]): Float
  posts: [Boolean!]!
}`

	var expectedErr [2]string
	expectedErr[0] = `List cannot contain NULLs at line: 8 column: 64`
	expectedErr[1] = `List cannot contain NULLs at line: 8 column: 78`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgumentNullCheck4(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[[String]]]! = [[["abc" "asdf"] ["defasj" "asdf" "asdf" null ] ["abc" null] null ] ["acb" "dfw" ] "wew"]): Float
  posts: [Boolean!]!
}`

	var expectedErr [3]string
	expectedErr[0] = `Value "acb" is not at required nesting of 3 at line: 8 column: 101`
	expectedErr[1] = `Value "dfw" is not at required nesting of 3 at line: 8 column: 107`
	expectedErr[2] = `Value "wew" is not at required nesting of 3 at line: 8 column: 115`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgumentNullCheck5(t *testing.T) {

	input := `type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [[String]]! = null ): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Value cannot be NULL at line: 8 column: 30`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestFieldArgument5(t *testing.T) {

	input := `
input Measure {
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
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestFieldArgument6(t *testing.T) {

	input := `
	input Measure {
    height: Float
    weight: Int
}

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
	//	var str1 = `inputMeasure{height:Floatweight:Int}enumAddress{NORTHSOUTHEAST}typePerson{name:String!age:Int!inputX(info:[Measure]=[{height:123.2weight:12}{height:1423.2weight:132}]):Floatposts:[Boolean!]!addres:Address!}`
	//	var str2 = `inputMeasure{height:Floatweight:Int}enumAddress{NORTHSOUTHEAST}typePerson{name:String!age:Int!inputX(info:[Measure]=[{height:123.2weight:12}{height:1423.2weight:132}]):Floatposts:[Boolean!]!addres:Address!}`
	expectedDoc := `enum Address {NORTH SOUTH EAST}
					input Measure {height:Float weight:Int}
					type Person {name:String! age:Int! inputX(info:[Measure]=[{height:123.2 weight:12}{height:1423.2 weight:132}]):Float posts:[Boolean!]!  addres:Address!}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}

}

func TestInputArgument1(t *testing.T) {

	input := `input Measure {
    height: Int
    name: String
}
type Person {
  name: String!
  age: Int!
  inputX(info: Measure! = {height: 8 name: "Ross" }): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestInputArgument1a(t *testing.T) {

	input := `input Family {
		name: String
		children: Int
	}
	input Measure {
    height: Float
    kids: Family
}
type Person {
  name: String!
  age: Int!
  inputX(info: Measure! = {height: 8.2 kids: {name: "Payne" children: 3.2} }): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Argument type "Family", value has type Float should be Int at line: 12 column: 71`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestInputArgument1b(t *testing.T) {

	input := `input Family {
		name: String
		Ages: [Int!]
	}
	input Measure {
    height: Float
    kids: Family
}
type Person {
  name: String!
  age: Int!
  inputX(info: Measure! = {height: 8.2 kids: {name: "Payne" Ages: [12 14.4 15 null]} }): Float
  posts: [Boolean!]!
}`

	var expectedErr [2]string
	expectedErr[0] = `Required type "Int", got "Float" at line: 12 column: 71`
	expectedErr[1] = `List cannot contain NULLs at line: 12 column: 79`
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestInputArgument1c(t *testing.T) {

	input := `input Family {
		name: String
		Ages: [Int!]
	}
	input Measure {
    height: Float
    kids: [Family]
}
type Person {
  name: String!
  age: Int!
  inputX(info: Measure! = {height: 8.2 kids: {name: "Payne" Ages: [12 14.2 15 null]} }): Float
  posts: [Boolean!]!
}`

	var expectedErr [3]string
	expectedErr[0] = `Argument "kids" for type "Measure" expected List at line: 12 column: 84`
	expectedErr[1] = `Required type "Int", got "Float" at line: 12 column: 71`
	expectedErr[2] = `List cannot contain NULLs at line: 12 column: 79`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, v := range errs {
		fmt.Println("***  ", v.Error())
	}
	for _, ex := range expectedErr {
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

func TestInputArgument1d(t *testing.T) {

	input := `input Family {
		name: String
		Ages: [Int!]
	}
	input Measure {
    height: Float
    kids: [Family]
}
type Person {
  name: String!
  age: Int!
  inputX(info: Measure! = {height: 8.2 kids: [ {name: "Payne" Ages: [12 14 15 null]} {name2: "Smith" Age: [1 3.2 ]}  ]   }): Float
  posts: [Boolean!]!
}`

	var expectedErr [3]string
	expectedErr[0] = `List cannot contain NULLs at line: 12 column: 79`
	expectedErr[1] = `field "name2" does not exist in type Family  at line: 12 column: 87`
	expectedErr[2] = `field "Age" does not exist in type Family  at line: 12 column: 102`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestInputArgument1e(t *testing.T) {

	input := `input Family {
		name: String
		Ages: [Int!]
	}
	input Measure {
    height: Float
    kids: [Family]
}
type Person {
  name: String!
  age: Int!
  inputX(info: Measure! = {height: 8.2 kids: [ {name: "Payne" Ages: [12 14 15 null]} {name2: "Smith" Age: [1 3]}  ]   }): Float
  posts: [Boolean!]!
}`

	var expectedErr [3]string
	expectedErr[0] = `List cannot contain NULLs at line: 12 column: 79`
	expectedErr[1] = `field "name2" does not exist in type Family  at line: 12 column: 87`
	expectedErr[2] = `field "Age" does not exist in type Family  at line: 12 column: 102`
	//
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()

	for _, ex := range expectedErr {
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

func TestInputArgument1f(t *testing.T) {

	input := `input Family {
		name: String
		Ages: [Int!]
	}
	input Measure {
    height: Float
    kids: Family
}
type Person {
  name: String!
  age: Int!
  inputX(info: Measure! = {height: 8.2 kids: [ {name: "Payne" Ages: [12 14 15 null]} {name2: "Smith" Age: [1 3]}  ]   }): Float
  posts: [Boolean!]!
}`

	//parseObject_test.go:2230: Unexpected Error = ["Value {name:\"Payne\" Ages:[12 14 15 null ]  }  is not at required nesting of 0 at line: 12 column: 84"]
	//parseObject_test.go:2230: Unexpected Error = ["Value {name2:\"Smith\" Age:[1 3 ]  }  is not at required nesting of 0 at line: 12 column: 112"]

	var expectedErr [6]string
	expectedErr[0] = `Argument "kids" for type "Measure" should not be in List at line: 12 column: 115`
	expectedErr[1] = `Value {name:"Payne" Ages:[12 14 15 null ]  }  should not be contained in a List at line: 12 column: 84`
	expectedErr[2] = `List cannot contain NULLs at line: 12 column: 79`
	expectedErr[3] = `Value {name2:"Smith" Age:[1 3 ]  }  should not be contained in a List at line: 12 column: 112`
	expectedErr[4] = `field "name2" does not exist in type Family  at line: 12 column: 87`
	expectedErr[5] = `field "Age" does not exist in type Family  at line: 12 column: 102`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestInputArgument1g(t *testing.T) {

	input := `input Family {
		name: String
		Ages: [Int!]
	}
	input Measure {
    height: Float
    kids: [Family]
}
type Person {
  name: String!
  age: Int!
  inputX(info: Measure! = {height: 8.2 kids: [ {name: "Payne" Ages: [12 14 15 null]} {name: "Smith" Ages: [[1 3][2 4]]}  ]   }): Float
  posts: [Boolean!]!
}`

	var expectedErr [6]string
	expectedErr[0] = `List cannot contain NULLs at line: 12 column: 79`
	expectedErr[1] = `Argument "Ages", nested List type depth different reqired 1, got 2 at line: 12 column: 101`
	expectedErr[2] = `Value 1 is not at required nesting of 1 at line: 12 column: 109`
	expectedErr[3] = `Value 3 is not at required nesting of 1 at line: 12 column: 111`
	expectedErr[4] = `Value 2 is not at required nesting of 1 at line: 12 column: 114`
	expectedErr[5] = `Value 4 is not at required nesting of 1 at line: 12 column: 116`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestInputArgument2(t *testing.T) {

	input := `input Measure {
    height: Float
    name: String
}
type Person {
  name: String!
  age: Int!
  inputX(info: [Measure!] = [{height: 8 name: "Ross" } {height: 897 name: "Jack" }]): Float
  posts: [Boolean!]!
}`

	var expectedErr [2]string
	expectedErr[0] = `Argument type "Measure", value has type Int should be Float at line: 8 column: 39`
	expectedErr[1] = `Argument type "Measure", value has type Int should be Float at line: 8 column: 65`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestInputArgument2a(t *testing.T) {

	input := `input Measure {
    height: Float
    name: String
}
type Person {
  name: String!
  age: Int!
  inputX(info: [Measure!] = [{hieght: 8.1 name: "Ross" } {height: 897.4 name2: "Jack" } null]): Float
  posts: [Boolean!]!
}`

	var expectedErr [3]string
	expectedErr[0] = `field "hieght" does not exist in type Measure  at line: 8 column: 31`
	expectedErr[1] = `field "name2" does not exist in type Measure  at line: 8 column: 73`
	expectedErr[2] = `List cannot contain NULLs at line: 8 column: 89`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestInputArgument3(t *testing.T) {

	input := `input myObj {
	address: String
	slope: Boolean
	}
	
	input Measure {
    height: Float
    metric: myObj
}
type Person {
  name: String!
  age: Int!
  inputX(info: [Measure!] = [{height: 8.1 metric: {address: "XX" slope: true} } {height: 897.2 metric: {address: "YYX" slope: false}  } null]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `List cannot contain NULLs at line: 13 column: 137`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestInputObj4(t *testing.T) {

	input := `type Measure {
    height: Float
    name: String
}
type Person {
  name: String!
  age: Int!
  inputX(info: [Measure] = [{height: 8 name: null } {height: 897 name: "Jack" } null]): Float
  posts: [Boolean!]!
}`

	var expectedErr [4]string
	expectedErr[0] = `Argument "info" type "Measure", is not an input type at line: 8 column: 17`
	expectedErr[1] = `Object "height" for type "Measure" expected Float got Int at line: 8 column: 38`
	expectedErr[2] = `Object "name" for type "Measure" expected String got Null at line: 8 column: 46`
	expectedErr[3] = `Object "height" for type "Measure" expected Float got Int at line: 8 column: 62`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println("*** ", v.Error())
	}
	for _, ex := range expectedErr {
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

func TestExtendField7a(t *testing.T) {

	input := `
	directive @addedDirective34 on FIELD_DEFINITION | ARGUMENT_DEFINITION
input Measure {
    height: Float
    weight: Int
}
enum Address {
	NORTH
	SOUTH
	EAST
}

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

	// inputMeasure{height:Floatweight:Int}enumAddress{NORTHSOUTHEAST}typePerson2@addedDirective34{name:String!age:Int!inputX(info:[Measure]=[{height:123.2weight:12}{height:1423.2weight:132}]):Floatposts:[Boolean!]!addres:Address!isHiddenLocally:Boolean}typePerson2@addedDirective34{name:String!age:Int!inputX(info:[Measure]=[{height:123.2weight:12}{height:1423.2weight:132}]):Floatposts:[Boolean!]!addres:Address!isHiddenLocally:Boolean}typePerson2@addedDirective34{name:String!age:Int!inputX(info:[Measure]=[{height:123.2weight:12}{height:1423.2weight:132}]):Floatposts:[Boolean!]!addres:Address!isHiddenLocally:Boolean}

	// 	expectedDoc := `
	// directive @addedDirective34 on | FIELD_DEFINITION| ARGUMENT_DEFINITION

	// input Measure {height : Float weight : Int }
	// enum Address{
	// NORTH
	// SOUTH
	// EAST
	// }

	// type Person2 @addedDirective34 {
	// name : String!
	// age : Int!
	// inputX(info : [Measure] =[{height:123.2 weight:12 }  {height:1423.2 weight:132 }  ] ) : Float
	// posts : [Boolean!]!
	// addres : Address!
	// isHiddenLocally : Boolean
	// }

	// `

	err := ast.DeleteType("Person2")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	var expectedErr [1]string
	expectedErr[0] = `Directive "@addedDirective34" is not registered for OBJECT usage at line: 25 column: 22`

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
		fmt.Println("here..", got.Error())
		found := false
		if len(expectedErr[0]) != 0 {
			for _, exp := range expectedErr {
				if trimWS(got.Error()) == trimWS(exp) {
					found = true
				}
			}
		}
		if !found {
			t.Errorf(`Unexpected Error = [%q]`, got.Error())
		}
	}
	// if compare(d.String(), expectedDoc) {
	// 	t.Errorf("Got:      [%s] \n", trimWS(d.String()))
	// 	t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
	// 	t.Errorf(`Unexpected: program.String() wrong. `)
	// }
}

func TestExtendField7b(t *testing.T) {

	input := `
input Measure {
    height: Float
    weight: Int
}
enum Address {
	NORTH
	SOUTH
	EAST
}

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

extend type Person2 { NewColumn : [[String!]!] }

`
	err := ast.DeleteType("Person2")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	// inputMeasure{height:Floatweight:Int}enumAddress{NORTHSOUTHEAST}typePerson2@addedDirective34{name:String!age:Int!inputX(info:[Measure]=[{height:123.2weight:12}{height:1423.2weight:132}]):Floatposts:[Boolean!]!addres:Address!isHiddenLocally:Boolean}typePerson2@addedDirective34{name:String!age:Int!inputX(info:[Measure]=[{height:123.2weight:12}{height:1423.2weight:132}]):Floatposts:[Boolean!]!addres:Address!isHiddenLocally:Boolean}typePerson2@addedDirective34{name:String!age:Int!inputX(info:[Measure]=[{height:123.2weight:12}{height:1423.2weight:132}]):Floatposts:[Boolean!]!addres:Address!isHiddenLocally:Boolean}

	// 	expectedDoc := `
	// directive @addedDirective34 on | FIELD_DEFINITION| ARGUMENT_DEFINITION

	// input Measure {height : Float weight : Int }
	// enum Address{
	// NORTH
	// SOUTH
	// EAST
	// }

	// type Person2 {
	// name : String!
	// age : Int!
	// inputX(info : [Measure] =[{height:123.2 weight:12 }  {height:1423.2 weight:132 }  ] ) : Float
	// posts : [Boolean!]!
	// addres : Address!
	// isHiddenLocally : Boolean
	//  NewColumn : [[String!]!]
	// }

	// `
	var expectedErr [1]string
	expectedErr[0] = ``

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
		fmt.Println("here..", got.Error())
		found := false
		if len(expectedErr[0]) != 0 {
			for _, exp := range expectedErr {
				if trimWS(got.Error()) == trimWS(exp) {
					found = true
				}
			}
		}
		if !found {
			t.Errorf(`Unexpected Error = [%q]`, got.Error())
		}
	}
	// if compare(d.String(), expectedDoc) {
	// 	t.Errorf("Got:      [%s] \n", trimWS(d.String()))
	// 	t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
	// 	t.Errorf(`Unexpected: program.String() wrong. `)
	// }
}

func TestExtendField7c(t *testing.T) {

	input := `
input Measure {
    height: Float
    weight: Int
}

directive @addedDirective67 on | FIELD_DEFINITION| ARGUMENT_DEFINITION | OBJECT

scalar Time 

extend type Person2 @addedDirective67

`
	// expectedDoc := `input Measure{height:Floatweight:Int}
	// directive@addedDirective67on|FIELD_DEFINITION|ARGUMENT_DEFINITION|OBJECT
	// scalarTime
	// typePerson2 @addedDirective67 {name:String!age:Int!
	// inputX(info:[Measure]=[{height:123.2weight:12}{height:1423.2weight:132}]):Float
	// posts:[Boolean!]!
	// addres:Address!isHiddenLocally:Boolean
	// NewColumn:[[String!]!]}`

	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	//	p.ClearCache()
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
	// if compare(d.String(), expectedDoc) {
	// 	t.Errorf("Got:      [%s] \n", trimWS(d.String()))
	// 	t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
	// 	t.Errorf(`Unexpected: program.String() wrong. `)
	// }
}

func TestObjCheckOutputType(t *testing.T) {

	input := `
	
input MyInput67 {
  x: Float
  y: Float
}
type Measure67 {
    height: Float
    weight: Int
    form: MyInput67
}
`

	var expectedErr [1]string
	expectedErr[0] = `Field "form" type "MyInput67", is not an output type at line: 10 column: 11` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestObjCheckInputType(t *testing.T) {

	input := `
	
type Myobject66 {
  x: Float
  y: Int
}
type Measure66 {
    height: String
    weight: Myobject66
    form (x : Myobject66) : Float
}
`

	var expectedErr [1]string
	expectedErr[0] = `Argument "x" type "Myobject66", is not an input type at line: 10 column: 15` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestCheckObjectType1(t *testing.T) {

	input := `
	
type Pet {
  x: Float
  y: Int
}
type Measure66 {
    height: String
    weight: Int
    form (xarg : Pet = {x:77.3 y:22}) : Float
}
`

	var expectedErr [1]string
	expectedErr[0] = `Argument "xarg" type "Pet", is not an input type at line: 10 column: 18` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestObjCheckInputObjVal2(t *testing.T) {

	input := `
	
input Pet {
  x: Float
  y: Int
}
type Measure66 {
    height: String
    weight: Myobject66
    form (xarg : [Pet] = [{x:77.3 y:22} {x:33.9 y: 32}]) : Float
}
`

	err := ast.DeleteType("Myobject66")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	var expectedErr [1]string
	expectedErr[0] = `Type "Myobject66" does not exist at line: 9 column: 13` //

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestMandatoryFieldsMissing(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: Int!
}
type Measure {
    height: String
    weight: Int
    form (xarg : [Pet] = [{x:77.3 y:22} {x:33.9}]) : Float
}
`

	var expectedErr [1]string
	expectedErr[0] = `Mandatory field "y" missing in type "Pet" at line: 10 column: 42` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestManadtoryWrongDataType(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: Int!
}
type Measure {
    height: String
    weight: Int
    form (xarg : [Pet]! = [{x:77.3 y:22} {y:33.9}]) : Float
}
`

	var expectedErr [1]string
	expectedErr[0] = `Argument "y" for type "Pet" expected Int got Float at line: 10 column: 45` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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

func TestManadtoryWithEmptyListCheck(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: [Int]!
}
type Measure {
    height: String
    weight: Int
    form (xarg : [Pet]! = [{x:77.3 y:22} {y:33.9}]) : Float
}
`

	var expectedErr [3]string
	expectedErr[0] = `Argument "y" for type "Pet" expected List at line: 10 column: 38`
	expectedErr[1] = `Argument "y" for type "Pet" expected List at line: 10 column: 45`
	expectedErr[2] = `Argument "y" for type "Pet" expected Int got Float at line: 10 column: 45`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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
func TestManadtoryWithNonEmptyListCheck(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: [Int!]
}
type Measure {
    height: String
    weight: Int
    form (xarg : [Pet]! = [{x:77.3 y:22} {y:33.9}]) : Float
}
`

	var expectedErr [3]string
	expectedErr[0] = `Argument "y" for type "Pet" expected List at line: 10 column: 38`           //
	expectedErr[1] = `Argument "y" for type "Pet" expected List at line: 10 column: 45`           //
	expectedErr[2] = `Argument "y" for type "Pet" expected Int got Float  at line: 10 column: 45` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, ex := range expectedErr {
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
func TestManadtoryListCheck(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: [Int]!
}
type Measure {
    height: String
    weight: Int
    form (xarg : [Pet]! = [{x:77.3 y:[22]} {x:1}]) : Float
}
`

	var expectedErr [2]string
	expectedErr[0] = `Argument "x" for type "Pet" expected Float got Int at line: 10 column: 47` //
	expectedErr[1] = `Mandatory field "y" missing in type "Pet" at line: 10 column: 45 `         //

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, e := range errs {
		fmt.Println(e.Error())
	}
	for _, ex := range expectedErr {
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

func TestNonMandatoryListCheck(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: [Int]
}
type Measure {
    height: String
    weight: Int
    form (xarg : [Pet]! = [{x:77.3 y:[22, 23]} {x:1.1}]) : Float
}
`

	expectedDoc := `typeMeasure{height:Stringweight:Intform(xarg:[Pet]!=[{x:77.3y:[22 23]}{x:1.1}]):Float}inputPet{x:Floaty:[Int]}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}

}

func TestWrongDataType3a(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: [Int]
}
type Measure {
    height: String
    weight: Int
    form (xarg : [Pet]! = [{x:77.3 y:22} {y:33.9}]) : Float
}
`

	var expectedErr [3]string
	expectedErr[0] = `Argument "y" for type "Pet" expected Int got Float at line: 10 column: 45`
	expectedErr[1] = `Argument "y" for type "Pet" expected List at line: 10 column: 38`
	expectedErr[2] = `Argument "y" for type "Pet" expected List at line: 10 column: 45`
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, e := range errs {
		fmt.Println(e.Error())
	}
	for _, ex := range expectedErr {
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

func TestWrongDataType3b(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: Int
}
type Measure {
    height: String
    weight: Int
    form (xarg : [Pet]! = [{x:77.3 y:22} {y:[33.9]}]) : Float
}
`

	var expectedErr [3]string
	expectedErr[0] = `Argument "y" for type "Pet" should not be in List at line: 10 column: 50`
	expectedErr[1] = `Value 33.9 should not be contained in a List at line: 10 column: 46`
	expectedErr[2] = `Required type "Int", got "Float" at line: 10 column: 46`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, e := range errs {
		fmt.Println(e.Error())
	}
	for _, ex := range expectedErr {
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

func TestMandatoryNoDefault(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: Int!
}
type Measure09 {
    height: String
    weight: Int
    form (xarg : [Pet]!) : Float
}
`
	expectedDoc := `type Measure09 {height:String weight:Int form(xarg:[Pet]!) : Float} input Pet {x:Float  y:Int!}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestMandatoryFieldsPresent(t *testing.T) {

	input := `
	
	input Pet {
  x: Float
  y: Int!
}
type Measure {
    height: String
    weight: Int
    form (xarg : [Pet] = [{x:77.3 y:22} { y:23}]) : Float
}
`

	expectedDoc := `typeMeasure{height:Stringweight:Intform(xarg:[Pet]=[{x:77.3y:22}{y:23}]):Float}inputPet{x:Floaty:Int!}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf("Unexpected, should be 0 errors, got %d", len(errs))
		for _, v := range errs {
			t.Errorf(`Unexpected error: %s`, v.Error())
		}
	}
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestInputJson(t *testing.T) {

	//input := `[{id: "jklw2ike", name: "Luke Skywalker", friends: [{id: "xjnJ4", name: "Leia Organa"},{id: "sjnJ5", name: "C-3PO"},{id: "ksejnJ6", name: "R2-D2"}], appearsIn: [NEWHOPE ,JEDI ], starships: [{ name: "Falcon", length: 23.4},{ name: "Cruiser", length: 68.2}] , totalCredits: 5532 },{id: "dfw23e", name: "Leia Organa", friends: [{id: "lwewJ6", name: "Luke Skywalker"},{id: "sjnJ5", name: "C-3PO"},{id: "ksejnJ6", name: "R2-D2"}], appearsIn: [NEWHOPE ,EMPIRE ], starships: [{ name: "BattleStar", length: 138.2}] , totalCredits: 2532 },]`
	input := `[{id: "jklw2ike", name: "Luke Skywalker", friends: [{id: "xjnJ4", name: "Leia Organa"},{id: "sjnJ5", name: "C-3PO"},{id: "ksejnJ6", name: "R2-D2"}], appearsIn: ["NEWHOPE" ,"JEDI" ], starships: [{ name: "Falcon", length: 23.4},{ name: "Cruiser", length: 68.2}] , totalCredits: 5532 },{id: "dfw23e", name: "Leia Organa", friends: [{id: "lwewJ6", name: "Luke Skywalker"},{id: "sjnJ5", name: "C-3PO"},{id: "ksejnJ6", name: "R2-D2"}], appearsIn: ["NEWHOPE" ,"EMPIRE" ], starships: [{ name: "BattleStar", length: 138.2}] , totalCredits: 2532 },]`

	l := lexer.New(input)
	p := New(l)
	d := p.ParseResponse()
	fmt.Println(d.String())

	if compare(d.String(), input) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(input))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}
