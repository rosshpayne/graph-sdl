package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/graph-sdl/lexer"
)

func compare(doc, expected string) bool {

	return trimWS(doc) != trimWS(expected)

}

func trimWS(input string) string {

	var out strings.Builder
	for _, v := range input {
		if !(v == '\u0009' || v == '\u0020' || v == '\u000A' || v == '\u000D' || v == ',') {
			out.WriteRune(v)
		}
	}
	return out.String()

}

func TestInvalidName(t *testing.T) {

	input := `
type __Person {
  name: String!
  age: Int!
  posts: [Float]!
}`

	expectedErr := `identifer [__Person] cannot start with two underscores at line: 2, column: 6`
	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	//fmt.Println(d.String())
	if len(errs) != 1 {
		t.Errorf(`Wrong- should be one error got %d`, len(errs))
	}
	for _, v := range errs {
		if trimWS(v.Error()) != trimWS(expectedErr) {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
		}
	}
}

func TestSchemaObject1(t *testing.T) {

	input := `
type Person {
  name: String!
  age: Int!
  posts: [Boolean!]!
}`

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

func TestSchemaObject1a(t *testing.T) {

	input := `type Person {
  name: String!
  age: Int!
  age: Float
  posts: [Boolean!]!
}`

	expectedErr := `Duplicate Field name "age" at line: 4, column: 3`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	if len(errs) != 1 {
		t.Errorf(`Wrong- should be 1 error got %d`, len(errs))
		for _, v := range errs {
			t.Errorf(v.Error())
		}
	}
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
		}
	}
}

func TestSchemaObject2(t *testing.T) {

	input := `
type Person {
  name: String!
  age: Int!
  age2: [Int!]
  posts: [[Post]]!
  posts2: [[Post]]
  posts3: [[Post]]!
  posts4: [[Post]!]!
  posts5: [[[Post]!]]
  posts6: [[[Post!]!]!]!!
}
type Post {
	size: String
	height: Float
	}
	`

	var expectedErr [2]string
	expectedErr[0] = `redundant ! at line: 11, column: 25` //
	expectedErr[1] = `Expected name identifer got ! of "!" at line: 11, column: 25`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for i, v := range errs {
		if i < 2 {
			if v.Error() != expectedErr[i] {
				t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr[i])
			}
		} else {
			t.Errorf(`Not expected Error =[%q]`, v.Error())
		}
	}
}

func TestSchemaObject3(t *testing.T) {

	input := `
type Person {
  name: String!
  age: Int!
  age2: [Int!]
  posts: [[Post]]!
  posts2: [[Post]]
  posts3: [[Post]]!
  posts4: [[Post]!]!
  posts5: [[[Post]!]]
  posts5a: [[[Post]]]!
  posts6: [[[Post!]!]!]!
  posts6a: [[[[[[[Post]!]!]]!]]!]!
}`
	expectedDoc := `type Person {
  name: String!
  age: Int!
  age2: [Int!]
  posts: [[Post]]!
  posts2: [[Post]]
  posts3: [[Post]]!
  posts4: [[Post]!]!
  posts5: [[[Post]!]]
  posts5a: [[[Post]]]!
  posts6: [[[Post!]!]!]!
  posts6a: [[[[[[[Post]!]!]]!]]!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, err := p.ParseDocument()
	//fmt.Println("stmt ", d.String())
	for _, e := range err {
		fmt.Println("*** ", e.Error())
	}
	//fmt.Println(d.String())
	if compare(d.String(), expectedDoc) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedDoc))
		t.Errorf(`*************  program.String() wrong.`)
	}

}

func TestSchemaObject4(t *testing.T) {

	input := `
type Person {
  name: String!
  age: Int!
  posts: [Post!]!
  allPersons(last: [[Int!]]!): [Person!]!
}`

	l := lexer.New(input)
	p := New(l)
	d, err := p.ParseDocument()
	for _, e := range err {
		fmt.Println("*** ", e.Error())
	}
	//fmt.Println(d.String())
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}

}

func TestSchemaRootQuery(t *testing.T) {

	input := `type Query {
  allPersons(last: Int): [Person!]!
}`

	expectedErr := ``

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
		}
	}
}

func TestSchemaRootTypes(t *testing.T) {

	input := `
	type Person { age: Int, height: Int }
	
	type Query {
  allPersons(last: Int): [Person!]!
}

type Mutation {
  createPerson(name: String!, age: Int!): Person!
}

type Subscription {
  newPerson: Person!
}`

	expectedErr := ``

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
		}
	}
}

func TestSchema1(t *testing.T) {

	input := `
schema {
    query: Query
    mutation: Mutation
    subscription: Subscription
}
	type Query {
	  allPersons(last: [Int] first: [[String!]]): [Person!]
	}

	#type Mutation {
	#  createPerson(name: String!, age: Int!): Person!
	#}

	#type Subscription {
	#  newPerson: Person!
#}

	type Person {
	  name: String!
	  age: Int!
	  other: [String!]
	  posts: [Post!]
	}

	type Post {
	  title: String!
	  author: [Person!]!
	}
	`

	expectedDoc := `type Person {
name : String!
age : Int!
other: [String!]!
posts : [Post!]
}

type Post {
title : String!
 author: [Person!]!
}

type Query {
allPersons(last : Int ) : [Person!]
}
schema {
query : Query 
mutation : Mutation
subscription : Subscription
}`
	var expectedErr [1]string
	expectedErr[0] = ``

	l := lexer.New(input)
	p := New(l)
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

	if compare(d.String(), expectedDoc) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(expectedDoc))
		t.Errorf(`*************  program.String() wrong.`)
	}

}

func TestUnresolvedType(t *testing.T) {

	input := `
type Person {
  name: String!
  age: Int
  mis: [Int]
  posts: [Post_!]!
  friend: [Person]
  Address: [lines]
  WritesFor: Publication!
}
type lines {
	length: Int
	words: Int
	sentences: Boolean
	pt: Person
	}
`
	var expectedErr [2]string
	expectedErr[0] = `Type "Post_" does not exist at line: 6 column: 11`
	expectedErr[1] = `Type "Publication" does not exist at line: 9 column: 14`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != len(expectedErr) {
		for _, v := range errs {
			fmt.Println(v.Error())
		}
		t.Errorf(fmt.Sprintf(`Not expected - should be 2 errors got %d`, len(errs)))
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

func TestTypeFieldArgs1(t *testing.T) {

	input := ` type Person {
name : String
picture(size : Int =12918@NME(if:"""abc""") ) : Url@iff(abc:123) 
}`

	var expectedErr [3]string
	expectedErr[0] = `Type "Url" does not exist at line: 3 column: 49`
	expectedErr[1] = `Type "@NME" does not exist at line: 3 column: 27`
	expectedErr[2] = `Type "@iff" does not exist at line: 3 column: 53`

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

func TestTypeFieldArgs2(t *testing.T) {

	input := ` type Person {
name : String
picture(size : Int =12918@NME(if:"""abc""") size2 : Boolean =true @NME34(ifx:false)) : Url@iff(abc:123) 
}`

	var expectedErr [4]string
	expectedErr[0] = `Type "Url" does not exist at line: 3 column: 88`
	expectedErr[1] = `Type "@NME" does not exist at line: 3 column: 27`
	expectedErr[2] = `Type "@NME34" does not exist at line: 3 column: 68`
	expectedErr[3] = `Type "@iff" does not exist at line: 3 column: 92`

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

func TestTypeFieldArgs4(t *testing.T) {

	input := ` type Person {
name : String
picture(size : Int =12918@NME(if:"""abc""") size2 : Boolean  @NME34(ifx:false) @NME33(i12fx:12.345) @NME54(if2x:123)) : Url@iff(abc:123) 
}`

	var expectedErr [6]string
	expectedErr[0] = `Type "Url" does not exist at line: 3 column: 121`
	expectedErr[5] = `Type "@NME" does not exist at line: 3 column: 27`
	expectedErr[1] = `Type "@NME34" does not exist at line: 3 column: 63`
	expectedErr[2] = `Type "@NME33" does not exist at line: 3 column: 81`
	expectedErr[3] = `Type "@NME54" does not exist at line: 3 column: 102`
	expectedErr[4] = `Type "@iff" does not exist at line: 3 column: 125`

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
