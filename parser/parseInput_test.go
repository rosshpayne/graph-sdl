package parser

import (
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestInput1(t *testing.T) {
	// directive @jun on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_FIELD_DEFINITION | INPUT_OBJECT
	// directive @june on | INPUT_OBJECT
	input := `
		directive @jun on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_FIELD_DEFINITION | INPUT_OBJECT
	directive @june on | INPUT_OBJECT

	input Point2D  @ june (asdf:234){
  x: Float = 123.23 @ jun (asdf:234)
  y: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
   y1: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
     y2: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
      y3: Float = 34 @ jun (asdf:"""asdflkjslkjd""" dei:234 uio:false)
       y4: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
        y5: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
         y6: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
          y63: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
}`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
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

func TestInputDuplicate(t *testing.T) {

	input := `
	
	input Point2D  @ june (asdf:234){
  x: Float = 123.23 @ june (asdf:234)
  y: Float = 34 @ june (asdf:"""asdflkjslkjd""")
   y1: Float = 34 @ june (asdf:"""asdflkjslkjd""")
     y2: Float = 34 @ june (asdf:"""asdflkjslkjd""")
       y: Int = 34 @ june (asdf:"""asdflkjslkjd""")
      y3: Float = 34 @ june (asdf:"""asdflkjslkjd""" dei:234 uio:false)
       y4: Float = 34 @ june (asdf:"""asdflkjslkjd""")
        y5: Float = 34 @ june (asdf:"""asdflkjslkjd""")
         y6: Float = 34 @ june (asdf:"""asdflkjslkjd""")
          y63: Float = 34 @ june (asdf:"""asdflkjslkjd""")
}
	directive @june on INPUT_FIELD_DEFINITION INPUT_OBJECT
	
`
	var expectedErr [1]string
	expectedErr[0] = `Duplicate input value name "y" at line: 8, column: 8`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
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

func TestInputInvalidName(t *testing.T) {

	input := `input Point2D  @ june (asdf:234){
  x: Float = 123.23 @ jun (asdf:234)
  y: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
   y1: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
     y2: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
      __y3: Float = 34 @ jun (asdf:"""asdflkjslkjd""" dei:234 uio:false)
       y4: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
        y5: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
         y6: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
          y63: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
}
`
	var expectedErr [1]string
	expectedErr[0] = `identifer [__y3] cannot start with two underscores at line: 6, column: 7`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
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
