package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/ast"
	"github.com/graph-sdl/lexer"
)

func TestImplements1(t *testing.T) {

	input := `
interface ValuedEntity {
  value: Int
}

type Person implements NamedEntity {
  name: String
  age: Int
}

`

	var expectedErr [1]string
	expectedErr[0] = `Type "NamedEntity" does not exist at line: 6 column: 24`

	err := ast.DeleteType("ValuedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
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

func TestImplements2(t *testing.T) {

	input := `
type NamedEntity {
  value: Int
}

type Person implements NamedEntity {
  name: String
  age: Int
}

`
	var expectedErr [1]string
	expectedErr[0] = `Implements type "NamedEntity" is not an Interface at line: 6 column: 24`

	err := ast.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
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

func TestImplements3(t *testing.T) {

	input := `
interface NamedEntity {
  name: String
  name2: Int

}
interface ValuedEntity {
  value: Int
}

type Person implements NamedEntity & ValuedEntity2 {
  name: String
  age: Int
}
`
	err := ast.DeleteType("ValuedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	var expectedErr [1]string
	expectedErr[0] = `Type "ValuedEntity2" does not exist at line: 11 column: 38`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
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

func TestImplements4x(t *testing.T) {

	input := `
interface NamedEntity {
  name: String
  name2: Int2

}
interface ValuedEntity {
  value: Int
  value2: FLoat
  value3: Boolean
  value4: Bool
}

type Person implements NamedEntity & ValuedEntity {
  name: String
  age: In
}
`

	var expectedErr [4]string
	expectedErr[0] = `Type "Int2", does not exist at line: 4 column: 10` //
	expectedErr[1] = `Type "FLoat", does not exist at line: 9 column: 11`
	expectedErr[2] = `Type "Bool", does not exist at line: 11 column: 11`
	expectedErr[3] = `Type "In", does not exist at line: 16 column: 8`

	err := ast.DeleteType("ValuedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Int2")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Bool")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("In")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	// for _, v := range errs {
	// 	fmt.Println(v.Error())
	// }
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected %d errors got %d`, len(expectedErr), len(errs))
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

func TestImplements4a(t *testing.T) {

	input := `
interface NamedEntity {
  name: String
  name2: Int2

}
interface ValuedEntity {
  value: Int
  value2: FLoat
  value3: Boolean
  value4: Bool
}

type Int2 {
	x: Int
}

type In {
	Age: Int
}

type Bool {
	z: Boolean
}


type Person implements NamedEntity & ValuedEntity {
  name: String
  age: In
}
`

	var expectedErr [3]string
	expectedErr[1] = `Type "Person" does not implement interface "NamedEntity", missing  "name2"`                             //
	expectedErr[2] = `Type "Person" does not implement interface "ValuedEntity", missing  "value" "value2" "value3" "value4"` //
	expectedErr[0] = `Type "FLoat" does not exist at line: 9 column: 11`                                                      //

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println("***********", v.Error())
	}
	if len(errs) != len(expectedErr) {
		t.Errorf(`Expected %d errors got %d`, len(expectedErr), len(errs))
	}
	for _, got := range errs {
		var found bool
		for _, exp := range expectedErr {
			if trimWS(got.Error()) == trimWS(exp) {
				found = true
			}
		}
		if !found {
			t.Errorf(`Not expected Error =[%q]`, got.Error())
		}
	}
}

func TestImplements5(t *testing.T) {

	input := `
interface NamedEntity {
  name: String
}

interface ValuedEntity {
  value: Int
}

type Person implements NamedEntity {
  name: String
  age: Int
}

type Business implements NamedEntity & ValuedEntity & NamedEntity {
  name: String
  value: Int
  employeeCount: Int
}
`
	var expectedErr [1]string
	expectedErr[0] = `Duplicate interface name at line: 15 column: 55` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	for i, v := range errs {
		if i < 1 {
			if trimWS(v.Error()) != trimWS(expectedErr[i]) {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
		fmt.Println(v.Error())
	}
}

func TestImplements6(t *testing.T) {

	input := `
	interface NamedEntity {
	  name: [[String!]!]!
	}

	interface ValuedEntity {
	  value: Int
	}

	type Person implements NamedEntity {
	  name: [[String!]!]!
	  age: Int
	}

	type Business implements NamedEntity & ValuedEntity {
	  name: [[String!]!]!
	  value: Int
	  employeeCount: Int
	}
	`
	err := ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("ValuedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Business")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != 0 {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be 0 errors got %d`, len(errs)))
	}

}

func TestImplements6a(t *testing.T) {

	input := `
	interface NamedEntity {
	  name: [[String!]!]!
	}

	interface ValuedEntity {
	  value: Int
	}

	type Person implements NamedEntity {
	  name: [[String!]!]!
	  age: Int
	}

	type Business implements NamedEntity & ValuedEntity {
	  name: [[String!]]!
	  value: Int
	  employeeCount: Int
	}
	`

	var expectedErr [1]string
	expectedErr[0] = `Type "Business" does not implement interface "NamedEntity", missing "name"`

	err := ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("ValuedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Business")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be %d errors got %d`, len(expectedErr), len(errs)))
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

func TestBadInterfaceKeyword(t *testing.T) {

	input := `
	interfacei NamedEntity6b {
	  name: String!
	}

	interface ValuedEntity6b {
	  value: Int
	}

	type Person6b implements NamedEntity6b {
	  name:  String!
	  age: Int
	}

	type Business6b implements NamedEntity6b & ValuedEntity6b {
	  name:  String!
	  employeeCount: Int
	}
	`

	var expectedErr [1]string
	expectedErr[0] = `Parse aborted. "interfacei" is not a statement keyword at line: 2, column: 2`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be %d errors got %d`, len(expectedErr), len(errs)))
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

func TestImplements7(t *testing.T) {

	input := `
interface NamedEntity {
  name: String
}

interface ValuedEntity {
  value: Int
}

type Person implements NamedEntity {
  name: String
  age: Int
}

type Business implements NamedEntity & ValuedEntity {
  name: String
  value: Int
  employeeCount: Int
}
`
	err := ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("ValuedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Business")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	if len(errs) != 0 {
		t.Errorf(`Wrong- should be zero errors got %d`, len(errs))
		for _, v := range errs {
			t.Errorf(v.Error())
		}
	}
	//fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		if len(trimWS(d.String())) != len(trimWS(input)) {
			t.Errorf(`*************  program.String() wrong.`)
		}
	}
}

func TestImplementsNotAllFields(t *testing.T) {

	input := `
interface NamedEntity {
  name: String
  
}

interface ValuedEntity {
  value: Int
  size: String
  length: Float
}

type Person implements NamedEntity {
  name: String
  age: Int
}

type Business implements & NamedEntity & ValuedEntity {
  name: String
  value: Int
  employeeCount: Int
}
`

	var expectedErr [1]string
	expectedErr[0] = `Type "Business" does not implement interface "ValuedEntity", missing "size" "length" `

	err := ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("ValuedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Person")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Business")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be %d errors got %d`, len(expectedErr), len(errs)))
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

func TestImplementsNotAllFields2(t *testing.T) {

	input := `
interface NamedEntity {
  name: String
  XXX: Boolean
  
}

interface ValuedEntity {
  value: Int
  size: [String]
  length: Float
}

type Person implements NamedEntity {
  name: String
  age: Int
}

type Business implements & NamedEntity & ValuedEntity {
  name: String
  value: Int
  length: String
  employeeCount: Int
}

type Business2 implements & NamedEntity & ValuedEntity {
  name: String
  XXX: Boolean
  size: String
  length: Float
  value: Int
  employeeCount: Int
}
`
	var expectedErr [4]string
	expectedErr[0] = `Type "Person" does not implement interface "NamedEntity", missing  "XXX"`
	expectedErr[1] = `Type "Business" does not implement interface "NamedEntity", missing  "XXX"`
	expectedErr[2] = `Type "Business" does not implement interface "ValuedEntity", missing  "size" "length"`
	expectedErr[3] = `Type "Business2" does not implement interface "ValuedEntity", missing  "size"`

	err := ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("ValuedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Business")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Business2")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be %d errors got %d`, len(expectedErr), len(errs)))
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

func TestImplementsNotAllFields3(t *testing.T) {

	input := `
interface NamedEntity {
  name: String
  XXX: Boolean
  
}

interface ValuedEntity {
  value: Int
  size: String
  length: Float
}

type Person implements NamedEntity {
  name: String
  age: [[Int!]]!
}

type Business implements & NamedEntity & ValuedEntity {
  name: String
  value: Int
  length: String
  employeeCount: Int
}

type Business2 implements & NamedEntity & ValuedEntity {
  name: String
    age: [[Int!]]!
  XXX: Boolean
  size: String
  length: Float
  value: Int
  employeeCount: Int
}
`
	var expectedErr [3]string
	expectedErr[0] = `Type "Person" does not implement interface "NamedEntity", missing  "XXX"`
	expectedErr[1] = `Type "Business" does not implement interface "NamedEntity", missing  "XXX"`
	expectedErr[2] = `Type "Business" does not implement interface "ValuedEntity", missing  "size" "length"`

	err := ast.DeleteType("NamedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("ValuedEntity")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Business")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = ast.DeleteType("Business2")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be %d errors got %d`, len(expectedErr), len(errs)))
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
