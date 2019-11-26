package parser

import (
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestMutation1(t *testing.T) {

	input := `type Mutation {
  createPerson(name: String!, age: Int!): [[PersonX]]
}
`
	var expectedErr [1]string
	expectedErr[0] = `Type "PersonX" does not exist at line: 2 column: 45`

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
