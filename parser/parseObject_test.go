package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/ast"
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
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
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
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	//fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected %d error to %d`, len(expectedErr), len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
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
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
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
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
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
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
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
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
  inputX(info: String = ["""abc \ndefasj \nasdf"""]): Float
  posts: [Boolean!]!
}`

	var expectedErr [1]string
	expectedErr[0] = `Argument "info", type is not a list but default value is a list at line: 9 column: 51`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) > len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}

func TestFieldArgument4d(t *testing.T) {

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

	var expectedErr [1]string
	expectedErr[0] = `Argument "info", type is a list but default value is not a list at line: 9 column: 29`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) > len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}

func TestFieldArgument4e(t *testing.T) {

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
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) > len(expectedErr) {
		t.Errorf(`***  Expected one error got %d.`, len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) > len(expectedErr) {
		t.Errorf(`***  Expected one error got %d.`, len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) > len(expectedErr) {
		t.Errorf(`***  Expected one error got %d.`, len(errs))
	}
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
  inputX(info: [[String]]! = [["abc"] ["defasj" "asdf" "asdf" null ] ["abc" null] null ]): Float
  posts: [Boolean!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf(`***  Expected no error got %d.`, len(errs))
	}
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
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
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf(`***  Expected no error got %d.`, len(errs))
	}
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
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
	//fmt.Println(d.String())
	if len(errs) > len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	if len(errs) > 0 {
		t.Errorf(`Got %d expected 0 errors`, len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`***  program.String() wrong.`)
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
	fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf(`***  Expected no error got %d.`, len(errs))
	}
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
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
	//fmt.Println(d.String())
	if len(errs) > len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	// for _, e := range errs {
	// 	fmt.Println("*** ", e.Error())
	// }
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
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

	var expectedErr [2]string
	expectedErr[0] = `Required type "Int", got "Float" at line: 12 column: 71`
	expectedErr[1] = `List cannot contain NULLs at line: 12 column: 79`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`***  Expected %d error got %d.`, len(expectedErr), len(errs))
	}
	for _, e := range errs {
		fmt.Println("*** ", e.Error())
	}
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
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
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

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
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

	var expectedErr [4]string
	expectedErr[0] = `Field, kids, is not a LIST type but input data is a LIST type, at line: 12 column: 40`
	expectedErr[1] = `List cannot contain NULLs at line: 12 column: 79`
	expectedErr[2] = `field "name2" does not exist in type Family  at line: 12 column: 87`
	expectedErr[3] = `field "Age" does not exist in type Family  at line: 12 column: 102`
	//	expectedErr[4] = `Argument "kids", nested List type depth different reqired 0, got 1 at line: 12 column: 40`

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

	var expectedErr [2]string
	expectedErr[0] = `List cannot contain NULLs at line: 12 column: 79`
	expectedErr[1] = `Argument "Ages", nested List type depth different reqired 1, got 2 at line: 12 column: 101`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
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

	var expectedErr [1]string
	expectedErr[0] = `Argument "info" type "Measure", is not an input type at line: 8 column: 17`

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

func TestExtendField7(t *testing.T) {

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

func TestExtendField7a(t *testing.T) {

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
	expectedErr[0] = `Field "form" type "MyInput67", is not an output type at line: 10 column: 11` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != 1 {
		t.Errorf(`Expected 1 error "... is not an output type", got none `)
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
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
	expectedErr[0] = `Argument "x" type "Myobject66", is not an input type at line: 10 column: 15` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected 1 error "... is not an output type", got none `)
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected %d error, got %d `, len(expectedErr), len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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

	var expectedErr [1]string
	expectedErr[0] = `Type "Myobject66" does not exist at line: 9 column: 13` //

	err := ast.DeleteType("Myobject66")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected %d error %d`, len(expectedErr), len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected 1 error, to %d`, len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	expectedErr[0] = `Argument type "Pet", value has type Float should be Int at line: 10 column: 45` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected 1 error, to %d`, len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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

	var expectedErr [1]string
	expectedErr[0] = `Argument type "Pet", value has type Float should be Int at line: 10 column: 45` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected 1 error, to %d`, len(errs))
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
	expectedErr[0] = `Argument type "Pet", value should be a List type Int at line: 10 column: 38`    //
	expectedErr[1] = `Argument type "Pet", value should be a List type Int at line: 10 column: 45`    //
	expectedErr[2] = `Argument type "Pet", value has type Float should be Int at line: 10 column: 45` //
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected 1 error, to %d`, len(errs))
	}
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
    form (xarg : [Pet]! = [{x:77.3 y:22} {x:1}]) : Float
}
`

	var expectedErr [2]string
	expectedErr[0] = `Argument type "Pet", value has type Int should be Float at line: 10 column: 45` //
	expectedErr[1] = `Mandatory field "y" missing in type "Pet" at line: 10 column: 43 `              //

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected 1 error, to %d`, len(errs))
	}
	// for _, v := range errs {
	// 	fmt.Println(v.Error())
	// }
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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
    form (xarg : [Pet]! = [{x:77.3 y:22} {x:1.1}]) : Float
}
`
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != 0 {
		t.Errorf(`Expected 1 error, to %d`, len(errs))
	}
	for _, v := range errs {
		fmt.Println(v.Error())
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

	var expectedErr [1]string
	expectedErr[0] = `Argument type "Pet", value has type Float should be Int at line: 10 column: 45` //
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected 1 error, to %d`, len(errs))
	}
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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

	var expectedErr [2]string
	expectedErr[0] = `Field, y, is not a LIST type but input data is a LIST type, at line: 10 column: 43` //
	expectedErr[1] = `Required type "Int", got "Float" at line: 10 column: 46`                            //
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected 1 error, to %d`, len(errs))
	}
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	for i, v := range errs {
		if i < len(expectedErr) {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
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

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != 0 {
		t.Errorf(`Expected 0 error, to %d`, len(errs))
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

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != 0 {
		t.Errorf(`Should be 0 errors, got %d`, len(errs))
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
