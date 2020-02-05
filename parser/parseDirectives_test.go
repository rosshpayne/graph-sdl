package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/db"
	"github.com/graph-sdl/lexer"
)

func TestDirectiveMultiple(t *testing.T) {

	input := `
input ExampleInputObjectDirective @ june (asdf:234) @ june2 (aesdf:234) @ june3 (as2df:"abc") {
  a: String = "AbcDef" @ ref (if:123) @ jack (sd: "abc") @ june (asdf:234) @ ju (asdf:234) @ judkne (asdf:234) @ junse (asdf:234) @ junqe (asdf:234) 
  b: Int! @june(asdf:234)  @ ju (asdf:234)
}
directive @june (asdf : Int = 66) on | FIELD_DEFINITION| ARGUMENT_DEFINITION | INPUT_OBJECT

`

	var expectedDoc = `
	directive @june (asdf : Int = 66) on | FIELD_DEFINITION| ARGUMENT_DEFINITION | INPUT_OBJECT
input ExampleInputObjectDirective @ june (asdf:234) @ june2 (aesdf:234) @ june3 (as2df:"abc") {
  a: String = "AbcDef" @ ref (if:123) @ jack (sd: "abc") @ june (asdf:234) @ ju (asdf:234) @ judkne (asdf:234) @ junse (asdf:234) @ junqe (asdf:234) 
  b: Int! @june(asdf:234)  @ ju (asdf:234)
}
`
	var expectedErr [11]string

	expectedErr[0] = `Item "@june2"  does not exist in document "DefaultDoc" at line: 2 column: 55`
	expectedErr[1] = `Item "@june3" does not exist in document "DefaultDoc" at line: 2 column: 75`
	expectedErr[2] = `Item "@ref"  does not exist in document "DefaultDoc" at line: 3 column: 26`
	expectedErr[3] = `Item "@jack"  does not exist in document "DefaultDoc" at line: 3 column: 41`
	expectedErr[4] = `Item "@ju" does not exist in document "DefaultDoc" at line: 3 column: 78`
	expectedErr[5] = `Item "@judkne" does not exist in document "DefaultDoc" at line: 3 column: 94`
	expectedErr[6] = `Item "@junqe" does not exist in document "DefaultDoc" at line: 3 column: 133`
	expectedErr[7] = `Item "@ju" does not exist at line: 4 column: 30`
	expectedErr[8] = `Item "@junse"  does not exist in document "DefaultDoc" at line: 3 column: 114`
	expectedErr[9] = `Directive "@june" is not registered for INPUT_FIELD_DEFINITION usage at line: 3 column: 60`
	expectedErr[10] = `Directive "@june" is not registered for INPUT_FIELD_DEFINITION usage at line: 4 column: 12`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println("Statement: ", d.String())
	for _, v := range errs {
		fmt.Println("**Error: ", v)
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

	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}

}

func TestInputDoesnotExist(t *testing.T) {

	input := `
extend input ExampleInputXYZ @ june (asdf:234) 
`
	var expectedErr [1]string
	expectedErr[0] = `Input type "ExampleInputXYZ" does not exist in document "DefaultDoc" at line: 2 column: 14`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	doc, errs := p.ParseDocument()
	fmt.Println(doc.String())
	// for _, v := range errs {
	// 	fmt.Println("**error: ", v.Error())
	// }
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

func TestDirectiveExtendInpDirDuplicate(t *testing.T) {

	input := `
directive @june on FIELD_DEFINITION | ARGUMENT_DEFINITION

input ExampleInputObjectDirective2 @ june {
	a: String 
	b: Int! 
}

extend input ExampleInputObjectDirective2 @ june (asdf:234) 
`
	var expectedErr [3]string
	expectedErr[0] = `Duplicate Directive name "@june" at line: 9, column: 45`
	expectedErr[1] = `extend for type "ExampleInputObjectDirective2" contains no changes at line: 0, column: 0`
	expectedErr[2] = `Directive "@june" is not registered for INPUT_OBJECT usage at line: 4 column: 38`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	_, errs := p.ParseDocument()
	// for _, v := range errs {
	// 	fmt.Println("Error: ", v)
	// }
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

func TestDirectiveStmt(t *testing.T) {

	input := `
directive @example (arg1: Int = 1256 arg2: String = "ABCdef") on |FIELD_DEFINITION | ARGUMENT_DEFINITION
`

	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
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

func TestDirectiveStmt2(t *testing.T) {

	input := `
directive @june on | FIELD_DEFINITION | ARGUMENT_DEFINITION
`

	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
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

func TestDirectiveStmt3(t *testing.T) {

	input := `" comment ....."
directive @invalidExample(arg: String @invalidExample) on ARGUMENT_DEFINITION
`

	expectedDoc := `directive@invalidExample(arg:String@invalidExample)on|ARGUMENT_DEFINITION`
	var expectedErr [1]string
	expectedErr[0] = `Directive "@invalidExample" that references itself, is not permitted at line: 2 column: 40`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
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
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestDirectiveStmt4(t *testing.T) {

	input := `" comment ....."
directive @example on
  | FIELD
  | FRAGMENT_SPREAD
  | INLINE_FRAGMENT
`

	expectedDoc := `directive@exampleon|FIELD|FRAGMENT_SPREAD|INLINE_FRAGMENT`
	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
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
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
}

func TestDirectiveStmt5(t *testing.T) {

	input := `" comment ....."
directive @__example on
  | FIELD
  | FRAGMENT_SPREAD
  | INLINE_FRAGMENT
`

	var expectedErr [1]string
	expectedErr[0] = `identifer [__example] cannot start with two underscores at line: 2, column: 12`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
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
func TestDirectiveInvalidLocation(t *testing.T) {

	input := `
directive @example on FIELD_DEFINITION | ARGUMENT_XYZ
`

	var expectedErr [1]string
	expectedErr[0] = `Invalid directive location "ARGUMENT_XYZ" at line: 2, column: 42`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
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

func TestDirectiveSelfRefCheck(t *testing.T) {

	input := `
	directive @exampleDirOK on |FIELD_DEFINITION | ARGUMENT_DEFINITION

type ExampleRefType @exampleDirRef {
	x (Nm: Float = 23.3 @exampleDirOK) : String @exampleDirOK
	y  (Nm: Float = 23.3 @exampleDirOK) : Int @exampleDirOK
}

input ExampleInput @exampleDirOK {
	x  : String @exampleDirOK
	y : Int @exampleDirRef
}

type exampleTypeOuter @exampleDirOK {
	x (Nm: Float = 23.3 @exampleDirOK) : String @exampleDirOK
	y  (Nm: Float = 23.3 @exampleDirRef) :  ExampleRefType @exampleDirRef
}

input exampleInput2 @exampleDirOK {
	x  : String @exampleDirOK
	y : exampleTypeOuter @exampleDirRef
}

type exampleTypeOuter2b @exampleDirOK {
	x (Nm: String = {x:"abc", y:1} @exampleDirOK ) : String 
}

	
directive @exampleDirRef (arg: exampleInput2@exampleDirRef ) on| FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_FIELD_DEFINITION | OBJECT
`

	expectedDoc := `directive @exampleDirOK on | FIELD_DEFINITION| ARGUMENT_DEFINITION

directive @exampleDirRef(arg : exampleInput2 @exampleDirRef ) on | FIELD_DEFINITION| ARGUMENT_DEFINITION| INPUT_FIELD_DEFINITION| OBJECT

input ExampleInput @exampleDirOK {x : String @exampleDirOK y : Int @exampleDirRef }

type ExampleRefType @exampleDirRef {
x(Nm : Float =23.3@exampleDirOK ) : String@exampleDirOK 
y(Nm : Float =23.3@exampleDirOK ) : Int@exampleDirOK 
}

input exampleInput2 @exampleDirOK {x : String @exampleDirOK y : exampleTypeOuter @exampleDirRef }

type exampleTypeOuter @exampleDirOK {
x(Nm : Float =23.3@exampleDirOK ) : String@exampleDirOK 
y(Nm : Float =23.3@exampleDirRef ) : ExampleRefType@exampleDirRef 
}

type exampleTypeOuter2b @exampleDirOK {
x(Nm : String ={x:"abc" y:1 } @exampleDirOK ) : String
}`

	err := db.DeleteType("exampleDirOK")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("ExampleRefType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("ExampleInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleTypeOuter")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleInput2")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleTypeOuter2b")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [13]string
	expectedErr[0] = `Directive "@exampleDirRef" that references itself, is not permitted at line: 29 column: 46`
	expectedErr[1] = `Directive "@exampleDirRef" references itself, is not permitted at line: 21 column: 24`
	expectedErr[2] = `Directive "@exampleDirRef" references itself, is not permitted at line: 16 column: 24`
	expectedErr[3] = `Directive "@exampleDirRef" references itself, is not permitted at line: 16 column: 58`
	expectedErr[4] = `Directive "@exampleDirRef" references itself, is not permitted at line: 4 column: 22`
	expectedErr[5] = `Directive "@exampleDirOK" is not registered for INPUT_OBJECT usage at line: 9 column: 21`
	expectedErr[6] = `Directive "@exampleDirOK" is not registered for INPUT_FIELD_DEFINITION usage at line: 10 column: 15`
	expectedErr[7] = `Directive "@exampleDirOK" is not registered for OBJECT usage at line: 14 column: 24`
	expectedErr[8] = `Directive "@exampleDirOK" is not registered for INPUT_OBJECT usage at line: 19 column: 22`
	expectedErr[9] = `Directive "@exampleDirOK" is not registered for INPUT_FIELD_DEFINITION usage at line: 20 column: 15`
	expectedErr[10] = `Directive "@exampleDirOK" is not registered for OBJECT usage at line: 24 column: 26`
	expectedErr[11] = `Mismatched types. The input data (object values in this case) does not match a Object or Input type. The reference type is a String`
	expectedErr[12] = `Field "y" of input type "exampleTypeOuter", must be an input type at line: 21 column: 6`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	for _, v := range errs {
		fmt.Println("*** Error: ", v)
	}
	fmt.Println(d.String())
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
	if compare(d.String(), expectedDoc) {
		t.Errorf("Got:      [%s] \n", trimWS(d.String()))
		t.Errorf("Expected: [%s] \n", trimWS(expectedDoc))
		t.Errorf(`Unexpected: program.String() wrong. `)
	}
	err = db.DeleteType("exampleDirOK")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("ExampleRefType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("ExampleInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleTypeOuter")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleInput2")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleTypeOuter2b")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
}

func TestDirectiveLocationCheck(t *testing.T) {

	input := `
directive @example (arg1: Int = 123 ) on | FIELD_DEFINITION | ARGUMENT_DEFINITION

input SomeInput @example (arg1: 23) {
  field: String = "ABC" @example
}

type SomeType {
  field(arg: Int @example): String @example
}


`

	err := db.DeleteType("SomeType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [2]string
	expectedErr[0] = `Directive "@example" is not registered for INPUT_OBJECT usage at line: 4 column: 18`
	expectedErr[1] = `Directive "@example" is not registered for INPUT_FIELD_DEFINITION usage at line: 5 column: 26`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())

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

func TestDirectiveErrArgName(t *testing.T) {

	input := `
directive @example (arg1: Int = 123 ) on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_OBJECT | INPUT_FIELD_DEFINITION

input SomeInput @example (arg2: 23) {
  field: String = "ABC" @example
}

type SomeType {
  field(arg: Int @example): String @example
}

`

	err := db.DeleteType("SomeType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [1]string
	expectedErr[0] = `Argument "arg2" is not a valid argument for directive  "@example" at line: 4 column: 27`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())

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

func TestDirectiveErrArgValue(t *testing.T) {

	input := `
directive @example (arg1: Int = 123 ) on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_OBJECT | INPUT_FIELD_DEFINITION

input SomeInput @example (arg1: "ABC") {
  field: String = "ABC" @example
}

type SomeType {
  field(arg: Int @example): String @example
}

`

	err := db.DeleteType("SomeType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [1]string
	expectedErr[0] = `Required type "Int", got "String" at line: 4 column: 33`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())

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

func TestFieldDirectiveExtraArg(t *testing.T) {

	input := `
directive @example (arg1: Int = 123 ) on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_OBJECT | INPUT_FIELD_DEFINITION

input SomeInput @example {
  field: String = "ABC" @example
}

type SomeType {
  field(arg: Int @example (arg1: 123 arg2: "ABC")): String @example
}

`

	err := db.DeleteType("SomeType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [1]string
	expectedErr[0] = `Argument "arg2" is not a valid argument for directive "@example" at line: 9 column: 38`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())

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

func TestFieldDirectiveArgTypeErr(t *testing.T) {

	input := `
directive @example (arg1: Int = 123 ) on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_OBJECT | INPUT_FIELD_DEFINITION

input SomeInput @example {
  field: String = "ABC" @example
}

type SomeType {
  field(arg: Int @example (arg1: "ABC" )): String @example
}

`

	err := db.DeleteType("SomeType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [1]string
	expectedErr[0] = `Required type "Int", got "String" at line: 9 column: 34`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())

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
func TestFieldDirectiveArgNonDefault(t *testing.T) {

	input := `
directive @example (arg1: Int = 123 ) on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_OBJECT | INPUT_FIELD_DEFINITION

input SomeInput @example {
  field: String = "ABC" @example
}

type SomeType {
  field(arg: Int @example (arg1 : 33)): String @example
}

`

	err := db.DeleteType("SomeType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())

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

func TestFieldDirectiveNoArgs(t *testing.T) {

	input := `
directive @example (arg1: Int = 123 ) on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_OBJECT | INPUT_FIELD_DEFINITION

input SomeInput @example {
  field: String = "ABC" @example
}

type SomeType {
  field(arg: Int @example): String @example
}

`
	//
	// query {SomeType {field }}			# uses arg default from type definition or if none specified from directive stmt. Requires all args to have defaults.
	// query {SomeType {field (arg1: 22)}}	# specify arg value in field query. All other args use defaults. Errors is not all args are specified that don't have defaults.

	err := db.DeleteType("SomeType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())

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

func TestDirectiveDefaultErr(t *testing.T) {

	input := `
directive @example (arg1: Int = "ABC" ) on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_OBJECT | INPUT_FIELD_DEFINITION

input SomeInput @example {
  field: String = "ABC" @example
}

type SomeType {
  field(arg: Int @example (arg1: "ABC" )): String @example
}

`

	err := db.DeleteType("SomeType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [2]string
	expectedErr[0] = `Required type "Int", got "String" at line: 2 column: 33`
	expectedErr[1] = `Required type "Int", got "String" at line: 9 column: 34`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())

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

func TestDirectiveWithArgs(t *testing.T) {

	input := `
directive @example (arg1 : Int = 5 arg2 : String = "ABC" ) on | FIELD_DEFINITION | ARGUMENT_DEFINITION

input SomeInput @example {
  field: String = "ABC" @example
}

type SomeType {
  field(arg: Int @example): String @example
}


`

	err := db.DeleteType("SomeType")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("exampleDirRef")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [2]string
	expectedErr[0] = `Directive "@example" is not registered for INPUT_OBJECT usage at line: 4 column: 18`
	expectedErr[1] = `Directive "@example" is not registered for INPUT_FIELD_DEFINITION usage at line: 5 column: 26`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
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

func TestObjectFieldBadArgs3(t *testing.T) {

	input := `
directive @example3 (arg1 : Int = 5 arg2 : String = "ABC" arg3: Float = 23.44 ) on | INPUT_OBJECT| FIELD_DEFINITION | ARGUMENT_DEFINITION| INPUT_FIELD_DEFINITION

type Query {
  hero(argx1: [Int]! = [67 55 "abc"] @example3 (arg1 : "abc"), argx2: String! = "ABCDEF", argy3: Float ) : SomeType3
}
input SomeInput3 @example3 (arg2: "DEF" ) {
  field: String = "ABC" @example3
}

type SomeType3 {
  somefield : String @example3
}

`
	var expectedErr []string = []string{`Required type "Int", got "String" at line: 5 column: 31`, `Required type "Int", got "String" at line: 5 column: 56`}

	err := db.DeleteType("SomeType3")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput3")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("@example3")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	for _, v := range errs {
		fmt.Println("err: ", v)
	}
	fmt.Println(d.String())
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

func TestSetup4DirectiveQueriesArgs3(t *testing.T) {

	input := `
directive @example3 (arg1 : Int = 5 arg2 : String = "ABC" arg3: Float = 23.44 ) on | INPUT_OBJECT| FIELD_DEFINITION | ARGUMENT_DEFINITION| INPUT_FIELD_DEFINITION

type Query {
  hero(argx1: [Int]! = [67 55] @example3 (arg1 : 1234), argx2: String! = "ABCDEF", argy3: Float ) : SomeType3
}
input SomeInput3 @example3 (arg2: "DEF" ) {
  field: String = "ABC" @example3
}

type SomeType3 {
  somefield : String @example3
}



`

	err := db.DeleteType("SomeType3")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput3")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("@example3")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
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

func TestSetup4DirectiveQueriesArg4(t *testing.T) {

	input := `
directive @example4 (arg1 : Int = 5 arg2 : String = "ABC" arg3: Float = 23.44 arg4: Int) on |INPUT_OBJECT| FIELD_DEFINITION | ARGUMENT_DEFINITION| INPUT_FIELD_DEFINITION

input SomeInput4 @example4 (arg2: "DEF" ) {
  field: String = "ABC" @example4
}

type SomeType4 {
  somefield(arg: Int @example4 (arg1 : 234)): String @example4
}

type Query {
  hero: [SomeType4]
}

`

	err := db.DeleteType("SomeType4")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("SomeInput4")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("@example4")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
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
