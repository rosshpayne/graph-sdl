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
  inputX(age: Int = 123): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, e := range errs {
		t.Errorf(`*** %s`, e.Error())
	}
	//fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
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
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Got %d when non expected`, len(errs))
	}
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	for i, v := range errs {
		if i < len(errs) {
			if v.Error() != expectedErr[i] {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if len(errs) > 0 {
		t.Errorf(`Got %d when 0 expected`, len(errs))
	}
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
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
	//fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected %d error to %d`, len(expectedErr), len(errs))
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
	//fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected %d error to %d`, len(expectedErr), len(errs))
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
	for _, v := range errs {
		fmt.Println("errors: ", v.Error())
	}
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected %d error to %d`, len(expectedErr), len(errs))
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

func TestCheckInputValueType4(t *testing.T) {

	input := `type Person88 {
  name: String!
  age: Int!
  inputX(age:[[String]] = ["xyss" "cat" ["abc","def" "Hij"] "xyz"]): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	// for _, v := range errs {
	// 	fmt.Println(v.Error())
	// }
	if len(errs) != 0 {
		t.Errorf(`Expected 0 error to %d`, len(errs))
	}
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
	}

}

func TestCheckInputValueType5(t *testing.T) {

	input := `type Person88 {
  name: String!
  age: Int!
  inputX(age:[[String]] = ["abc","def"  "adf" ]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Argument "age", nested List type depth different reqired 2, got 1 at line: 4 column: 47` //

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected %d error to %d`, len(expectedErr), len(errs))
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
	if len(errs) > 0 {
		t.Errorf(`***  Expected no errors got %d.`, len(errs))
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
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
	fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf(`***  Expected no errors got %d.`, len(errs))
	}
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
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
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) > 1 {
		t.Errorf(`***  Expected one error got %d.`, len(errs))
	}
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
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

func TestFieldArgument4(t *testing.T) {

	input := `
type Measure {
    height: Float
    weight: Int
}
type Person {
  name: String!
  age: Int!
  inputX(info: [Int!] = [1,2,4 56 345 2342 234 25252 2525223 null]): Float
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
		t.Errorf(`***  program.String() wrong.`)
	}
}

func TestFieldArgTypeNotFound(t *testing.T) {

	input := `

type Person {
  name: String!
  age: Int!
  inputX(info: [[int!]] = [[1,2,4 56] [345 2342 234 25252 2525223 null]]): Float
  posts: [Boolean!]!
}`
	var expectedErr [1]string
	expectedErr[0] = `Type "int" does not exist at line: 6 column: 18`
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
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
	//	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
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
		t.Errorf(`***  program.String() wrong.`)
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
	fmt.Println(d.String())
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
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
		t.Errorf(`***  program.String() wrong.`)
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
	expectedDoc := `type Person2 @addedDirective34 {
name : String!
age : Int!
inputX(info : [Measure] =[{height:123.2 weight:12 }  {height:1423.2 weight:132 }  ] ) : Float
posts : [Boolean!]!
addres : Address!
isHiddenLocally : Boolean
}`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), expectedDoc) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedDoc))
		t.Errorf(`***  program.String() wrong.`)
	}
}

func TestFieldArgument7a(t *testing.T) {

	input := `
	

extend type Person2 @addedDirective67

`
	expectedDoc := `type Person2 @addedDirective34 @addedDirective67 {
name : String!
age : Int!
inputX(info : [Measure] =[{height:123.2 weight:12 }  {height:1423.2 weight:132 }  ] ) : Float
posts : [Boolean!]!
addres : Address!
isHiddenLocally : Boolean
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	fmt.Println(d.String())
	if compare(d.String(), expectedDoc) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedDoc))
		t.Errorf(`***  program.String() wrong.`)
	}
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
	expectedErr[0] = `Field "form" type "MyInput6", is not an output type at line: 10 column: 11` //
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
	expectedErr[0] = `Field "form" type "Myobject66", is not an input type at line: 10 column: 15` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
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

func TestObjCheckInputObjVal1(t *testing.T) {

	input := `
	
type Pet {
  x: Float
  y: Int
}
type Measure66 {
    height: String
    weight: Myobject66
    form (xarg : Pet = {x:77.3 y:22}) : Float
}
`

	var expectedErr [1]string
	expectedErr[0] = `Field "form" type "Myobject66", is not an input type at line: 10 column: 15` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
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

func TestObjCheckInputObjVal2(t *testing.T) {

	input := `
	
type Pet {
  x: Float
  y: Int
}
type Measure66 {
    height: String
    weight: Myobject66
    form (xarg : [Pet] = [{x:77.3 y:22} {x:33.9 y: 32}]) : Float
}
`

	var expectedErr [1]string
	expectedErr[0] = `Field "form" type "Myobject66", is not an input type at line: 10 column: 15` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
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
