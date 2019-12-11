package parser

import (
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

func TestImplements6x(t *testing.T) {

	input := `
	
	type Business implements NamedEntity & ValuedEntity {
	  name: [[String!]!]!
	  value: Int
	  employeeCount: Int
	}
	
	interface NamedEntity {
	  name: [[String!]!]!
	}
	
	type Person implements NamedEntity {
	  name: [[String!]!]!
	  age: Int
	}
	
	interface ValuedEntity {
	  value: Int
	}

	`
	// replace entities with their above definitions.
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

	// expectedDoc := `type Business implements NamedEntity & ValuedEntity {name:[[String!]!]!value:Intemployee Count:Int}
	// 				interface NamedEntity {name:[[String!]!]!}
	// 				type Person implements NamedEntity{name:[[String!]!]! age:Int}
	// 				interface ValuedEntity {value:Int}`

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

func TestInterfaceBadKeyword(t *testing.T) {

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
	expectedDoc := `type Business implements NamedEntity & ValuedEntity {name:String value:IntemployeeCount:Int}
				 interface NamedEntity {name:String} 
				 type Person implements NamedEntity {name:String age:Int} 
				 interface ValuedEntity {value:Int}`

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

func TestSetupFragments(t *testing.T) {

	input := `
enum Episode {
  NEWHOPE
  EMPIRE
  JEDI
}

type Starship {
  id: ID!
  name: String!
  length(unit: LengthUnit = METER): Float
}

interface Character {
  id: ID!
  name: String!
  friends: [Character]
  appearsIn: [Episode]!
}

type Human implements Character {
  id: ID!
  name: String!
  friends: [Character]
  appearsIn: [Episode]!
  starships: [Starship]
  totalCredits: Int
}

type Droid implements Character {
  id: ID!
  name: String!
  friends: [Character]
  appearsIn: [Episode]!
  primaryFunction: String
}`

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
