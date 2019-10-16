package parser

import (
	"fmt"
	"testing"

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
	expectedErr[0] = `Type "NamedEntity" is not defined at line: 6 column: 24`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// fmt.Println(d.String())
	// for i, v := range errs {
	// 	fmt.Println(i, v.Error())
	// }
	if len(errs) != 1 {
		t.Errorf(`Expect 4 errors got %d`, len(errs))
	} else {
		for i, v := range errs {
			if i < 1 && v.Error() != expectedErr[i] {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
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

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// fmt.Println(d.String())
	// for i, v := range errs {
	// 	fmt.Println(i, v.Error())
	// }
	if len(errs) != 1 {
		t.Errorf(`Expect 4 errors got %d`, len(errs))
	} else {
		for i, v := range errs {
			if i < 1 && v.Error() != expectedErr[i] {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
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

	var expectedErr [4]string
	expectedErr[0] = `Type "ValuedEntity2" is not defined at line: 11 column: 38`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	// fmt.Println(d.String())
	// for i, v := range errs {
	// 	fmt.Println(i, v.Error())
	// }
	if len(errs) != 1 {
		t.Errorf(`Expect 4 errors got %d`, len(errs))
	} else {
		for i, v := range errs {
			if i < 1 && v.Error() != expectedErr[i] {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		}
	}
}

func TestImplements4(t *testing.T) {

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
	expectedErr[0] = `Type "Int2", not defined at line: 4 column: 10` //
	expectedErr[1] = `Type "FLoat", not defined at line: 9 column: 11`
	expectedErr[2] = `Type "Bool", not defined at line: 11 column: 11`
	expectedErr[3] = `Type "In", not defined at line: 16 column: 8`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	// if len(errs) != 1 {
	// 	t.Errorf(`Expected %d error to %d`, len(expectedErr), len(errs))
	// }
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
	var expectedErr [4]string
	expectedErr[0] = `Duplicate interface name at line: 15 column: 55` //
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	for i, v := range errs {
		if i < 2 {
			if v.Error() != expectedErr[i] {
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

	l := lexer.New(input)
	p := New(l)
	d, err := p.ParseDocument()
	if len(err) != 0 {
		t.Errorf(`Wrong- should be zero errors got %d`, len(err))
		for _, v := range err {
			t.Errorf(v.Error())
		}
	}
	//fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestImplements6a(t *testing.T) {

	input := `
	interface NamedEntity6a {
	  name: [[String!]!]!
	}

	interface ValuedEntity6a {
	  value: Int
	}

	type Person6a implements NamedEntity6a {
	  name: [[String!]!]!
	  age: Int
	}

	type Business6a implements NamedEntity6a & ValuedEntity6a {
	  name: [[String!]]!
	  value: Int
	  employeeCount: Int
	}
	`

	var expectedErr [1]string
	expectedErr[0] = `Object type "Business6a" does not implement interface "NamedEntity6a", missing "name"`

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
			if v.Error() != expectedErr[i] {
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
			if v.Error() != expectedErr[i] {
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

type Business implements & NamedEntity & ValuedEntity {
  name: String
  value: Int
  employeeCount: Int
}
`

	expectedDoc := `
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
	l := lexer.New(input)
	p := New(l)
	d, err := p.ParseDocument()
	if len(err) != 0 {
		t.Errorf(`Wrong- should be zero errors got %d`, len(err))
		for _, v := range err {
			t.Errorf(v.Error())
		}
	}
	//fmt.Println(d.String())
	if compare(d.String(), expectedDoc) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedDoc))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestImplementsAllFields(t *testing.T) {

	input := `
interface NamedEntity2a {
  name: String
  
}

interface ValuedEntity2a {
  value: Int
  size: String
  length: Float
}

type Person2a implements NamedEntity2a {
  name: String
  age: Int
}

type Business2a implements & NamedEntity2a & ValuedEntity2a {
  name: String
  value: Int
  employeeCount: Int
}
`
	l := lexer.New(input)
	p := New(l)
	d, err := p.ParseDocument()
	for _, v := range err {
		t.Errorf(v.Error())
	}
	//fmt.Println(d.String())
	if len(err) == 0 {
		if compare(d.String(), input) {
			fmt.Println(trimWS(d.String()))
			fmt.Println(trimWS(input))
			t.Errorf(`*************  program.String() wrong.`)
		}
	}
}

func TestImplementsAllFields2(t *testing.T) {

	input := `
interface NamedEntity2a {
  name: String
  XXX: Boolean
  
}

interface ValuedEntity2a {
  value: Int
  size: [String]
  length: Float
}

type Person2a implements NamedEntity2a {
  name: String
  age: Int
}

type Business2a implements & NamedEntity2a & ValuedEntity2a {
  name: String
  value: Int
  length: String
  employeeCount: Int
}

type Business3a implements & NamedEntity2a & ValuedEntity2a {
  name: String
  XXX: Boolean
  size: String
  length: Float
  value: Int
  employeeCount: Int
}
`
	l := lexer.New(input)
	p := New(l)
	d, err := p.ParseDocument()
	for _, v := range err {
		t.Errorf(v.Error())
	}
	//fmt.Println(d.String())
	if len(err) == 0 {
		if compare(d.String(), input) {
			fmt.Println(trimWS(d.String()))
			fmt.Println(trimWS(input))
			t.Errorf(`*************  program.String() wrong.`)
		}
	}
}

func TestImplementsAllFields3(t *testing.T) {

	input := `
interface NamedEntity2a {
  name: String
  XXX: Boolean
  
}

interface ValuedEntity2a {
  value: Int
  size: String
  length: Float
}

type Person2a implements NamedEntity2a {
  name: String
  age: [[Int!]]!
}

type Business2a implements & NamedEntity2a & ValuedEntity2a {
  name: String
  value: Int
  length: String
  employeeCount: Int
}

type Business3a implements & NamedEntity2a & ValuedEntity2a {
  name: String
    age: [[Int!]]!
  XXX: Boolean
  size: String
  length: Float
  value: Int
  employeeCount: Int
}
`
	l := lexer.New(input)
	p := New(l)
	d, err := p.ParseDocument()
	for _, v := range err {
		t.Errorf(v.Error())
	}
	//fmt.Println(d.String())
	if len(err) == 0 {
		if compare(d.String(), input) {
			fmt.Println(trimWS(d.String()))
			fmt.Println(trimWS(input))
			t.Errorf(`*************  program.String() wrong.`)
		}
	}
}
