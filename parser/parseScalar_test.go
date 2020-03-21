package parser

import (
	"fmt"
	"testing"
	"time"

	"github.com/graph-sdl/db"
	"github.com/graph-sdl/lexer"
)

func TestScalarStmtInvalidName(t *testing.T) {

	input := `"my example scalar" 
	scalar __Time @dir 
`
	var expectedErr []string = []string{
		`identifer [__Time] cannot start with two underscores at line: 2, column: 9`,
		`Item "@dir" does not exist in document "DefaultDoc"  at line: 2 column: 17`,
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	for _, e := range errs {
		fmt.Println("Err: ", e.Error())
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
	//TODO - check didn't create scalar type in DB
}

func TestScalarStmt(t *testing.T) {

	input := `"my example scalar" 
	scalar Time @dir2 
`
	expectedStr := `scalar Time @dir2 `

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if compare(d.String(), expectedStr) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedStr))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestScalarArgumentx(t *testing.T) {

	input := `"my example scalar" 
     type ArgumentObj {
     	fd (xyz: Time = "Feb 3, 2013 at 7:54pm (PST)" ):Float
     }
`
	expectedStr := `type ArgumentObj {
     	fd (xyz: Time = "2013-02-03 19:54:00 +0000 PST" ):Float
     } `

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println("X", d.String())
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if len(errs) > 0 {
		t.Errorf(`Not expected - should be 0 errors got %d`, len(errs))
	}
	if compare(d.String(), expectedStr) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedStr))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestScalarArgumentInvalidValue(t *testing.T) {

	input := `"my example scalar" 
     type ArgumentObj {
     	fd (xyz: Time = "Feb 3, 2013 at 7:784pm (PST)" ):Float
     }
`

	var expectedErr [1]string
	expectedErr[0] = `Error in parsing of Time value `

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be %d errors got %d`, len(expectedErr), len(errs)))
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

// TODO: will need to decide what to do when
// we create a new type with the same name as an existing type object.
//
// func TestNotScalarNameDifferentType(t *testing.T) {

// 	input := `"my example of enum "
// 	enum Time {SECOND MINUTE HOUR DAY MONTH}
// `
// 	expectedStr := `enum Time {SECOND MINUTE HOUR DAY MONTH} `

// 	l := lexer.New(input)
// 	p := New(l)
// 	d, errs := p.ParseDocument()
// 	fmt.Println(d.String())
// 	for _, v := range errs {
// 		fmt.Println(v.Error())
// 	}
// 	if compare(d.String(), expectedStr) {
// 		fmt.Println(trimWS(d.String()))
// 		fmt.Println(trimWS(expectedStr))
// 		t.Errorf(`*************  program.String() wrong.`)
// 	}
// }

func TestScalarUsage(t *testing.T) {

	input := `"my example scalar" 
	type Customer {
		Name: String
		Location:   Float
		Joined:Time
	}
`
	expectedStr := `		type Customer {
		Name: String
		Location:   Float
		Joined:Time
	}`

	err := db.DeleteType("Customer")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) > 0 {
		t.Errorf(fmt.Sprintf(`Not expected - should be 0 errors got %d`, len(errs)))
	}
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if compare(d.String(), expectedStr) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedStr))
		t.Errorf(`*************  program.String() wrong.`)
	}
}

func TestScalarCheckx(t *testing.T) {

	input := `
	
	type Foo {
		abc: Int
		Def: Time
	}
	
	type Address {
		line1: String
		line2: String
		line3: String
		first: Foo
	}
	
	type Customer {
		Name: String
		Location:   Float
		Joined:Address
	}
`
	var expectedErr []string = []string{
		`Item "Time" does not exist in document "DefaultDoc"  at line: 5 column: 8 `,
	}
	// err := db.DeleteType("Time")
	// if err != nil {
	// 	t.Errorf(`Not expected Error =[%q]`, err.Error())
	// }
	err := db.DeleteType("Foo")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	} else {
		fmt.Println("Success deleted...")
	}
	// err = db.DeleteType("Address")
	// if err != nil {
	// 	t.Errorf(`Not expected Error =[%q]`, err.Error())
	// }
	// err = db.DeleteType("Customer")
	// if err != nil {
	// 	t.Errorf(`Not expected Error =[%q]`, err.Error())
	// }

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

func TestScalarCheckNoType(t *testing.T) {
	// this is not a nested check just a normal stmt level check
	input := `
	
	type Foo {
		abc: Int
		Def: Time2
	}
	
	type Address {
		line1: String
		line2: String
		line3: String
		first: Foo
	}
	
	type Customer {
		Name: String
		Location:   Float
		Joined:Address
	}
`

	// err := db.DeleteType("Time")
	// if err != nil {
	// 	t.Errorf(`Not expected Error =[%q]`, err.Error())
	// }
	err := db.DeleteType("Foo")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Address")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Customer")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	var expectedErr []string = []string{
		`Item "Time2" does not exist in document "DefaultDoc" at line: 5 column: 8`,
	}
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be %d errors got %d`, len(expectedErr), len(errs)))
	}
	//fmt.Println(d.String())

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

func TestScalarCheckAllOK(t *testing.T) {
	// this is not a nested check just a normal stmt level check
	input := `
	
	type Foo {
		abc: Int
		Def: Time
	}
	
	type Address {
		line1: String
		line2: String
		line3: String
		first: Foo
	}
	
	type Customer {
		Name: String
		Location:   Float
		Joined:Address
	}
`

	// err := db.DeleteType("Time")
	// if err != nil {
	// 	t.Errorf(`Not expected Error =[%q]`, err.Error())
	// }
	err := db.DeleteType("Foo")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Address")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Customer")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != 0 {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be 0 errors got %d`, len(errs)))
	}
	//fmt.Println(d.String())
}

func TestScalarNestedTypeCheck(t *testing.T) {

	input := `
	
	type Customer2 {
		Name: String
		Location:   Float
		Joined:Address
		Family: Customer
	}
`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != 0 {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be 0 errors got %d`, len(errs)))
	}
	//fmt.Println(d.String())
}

func TestScalarCheckAllOK2(t *testing.T) {
	// this is not a nested check just a normal stmt level check
	input := `
	
	type Customer {
		Name: String
		Location:   Float
		Joined:Address
	}

	type Address {
		line1: String
		line2: String
		line3: String
		first: Foo
	}
	
	type Foo {
		abc: Int
		Def: Time
	}
	
`
	var expectedErr [1]string
	expectedErr[0] = ``
	// err := db.DeleteType("Time")
	// if err != nil {
	// 	t.Errorf(`Not expected Error =[%q]`, err.Error())
	// }
	err := db.DeleteType("Foo")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Address")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	err = db.DeleteType("Customer")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}

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

func TestScalarNestedTypeCheck2(t *testing.T) {

	input := `
	
	type Customer2 {
		Name: String
		Location:   Float
		Joined:Address
		Family: Customer
	}
`
	//TODO: line and column should be reflected in current script not a script from the db.
	// `Customer has type "Time" that does not exist at line: 7 column: 11`
	var expectedErr [1]string
	expectedErr[0] = `Type "Time" does not exist at line: 4 column: 7`

	// *** unrealistic delete, but its purpose is to check that the program is performing nested type checking.

	err := db.DeleteType("Time")
	if err != nil {
		t.Errorf(`Not expected Error =[%q]`, err.Error())
	}
	time.Sleep(2)

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())

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
