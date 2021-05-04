package parser

import (
	"testing"

	"github.com/rosshpayne/graph-sdl/lexer"
)

func TestInput1(t *testing.T) {
	// directive @jun on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_FIELD_DEFINITION | INPUT_OBJECT
	// directive @june on | INPUT_OBJECT
	input := `
		directive @jun (asdf: String) on | FIELD_DEFINITION | ARGUMENT_DEFINITION | INPUT_FIELD_DEFINITION | INPUT_OBJECT
	directive @june (asdf: String) on | INPUT_OBJECT

	input Point2D  @ june (asdf:"abc"){
  x: Float = 123.23 @ jun (asdf:"234")
  y: Float = 34.4 @ jun (asdf:"""asdflkjslkjd""")
   y1: Float = 34.3 @ jun (asdf:"""asdflkjslkjd""")
     y2: Float = 34.7 @ jun (asdf:"""asdflkjslkjd""")
      y3: Float = 34.8 @ jun (asdf:"""asdflkjslkjd""" dei:234 uio:false)
       y4: Float = 34.9 @ jun (asdf:"""asdflkjslkjd""")
        y5: Float = 34.3 @ jun (asdf:"""asdflkjslkjd""")
         y6: Float = 34.2 @ jun (asdf:"""asdflkjslkjd""")
          y63: Float = 34.1 @ jun (asdf:"""asdflkjslkjd""")
}`

	var expectedErr [2]string
	expectedErr[0] = `Argument "dei" is not a valid name for directive "@jun" at line: 10 column: 55`
	expectedErr[1] = `Argument "uio" is not a valid name for directive "@jun" at line: 10 column: 63`

	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
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

func TestInputDuplicate(t *testing.T) {

	input := `
	directive @june (asdf: String) on | INPUT_FIELD_DEFINITION | INPUT_OBJECT	
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

	
`
	var expectedErr []string = []string{
		`Duplicate input value name "y" at line: 8, column: 8`,
		`Required type for argument "asdf" is String, got Int at line: 3 column: 25`,
		`Required type for argument "asdf" is String, got Int at line: 4 column: 29`,
		`Required type for argument "y" is Float, got Int at line: 5 column: 3`,
		`Required type for argument "y1" is Float, got Int at line: 6 column: 4`,
		`Required type for argument "y2" is Float, got Int at line: 7 column: 6`,
		`Required type for argument "y3" is Float, got Int at line: 9 column: 7`,
		`Required type for argument "y4" is Float, got Int at line: 10 column: 8`,
		`Required type for argument "y5" is Float, got Int at line: 11 column: 9`,
		`Required type for argument "y6" is Float, got Int at line: 12 column: 10`,
		`Required type for argument "y63" is Float, got Int at line: 13 column: 11`,
		`Argument "dei" is not a valid name for directive "@june" at line: 9 column: 54`,
		`Argument "uio" is not a valid name for directive "@june" at line: 9 column: 62`,
	}
	l := lexer.New(input)
	p := New(l)
	p.ClearCache()
	d, errs := p.ParseDocument()
	//fmt.Println(d.String())
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

func TestInputInvalidName(t *testing.T) {

	input := `input Point2D  @ june (asdf:"123"){
  x: Float = 123.23 @ jun (asdf:"234")
  y: Float = 34.3 @ jun (asdf:"""asdflkjslkjd""")
   y1: Int = 34 @ jun (asdf:"""asdflkjslkjd""")
     y2: Float = 34.2 @ jun (asdf:"""asdflkjslkjd""")
      __y3: Float = 34.2 @ jun (asdf:"""asdflkjslkjd""")
       y4: Int = 34 @ jun (asdf:"""asdflkjslkjd""")
        y5: Int = 34 @ jun (asdf:"""asdflkjslkjd""")
         y6: Int = 34 @ jun (asdf:"""asdflkjslkjd""")
          y63: Int = 34 @ jun (asdf:123)
}
	directive @june (asdf: String) on | INPUT_FIELD_DEFINITION | INPUT_OBJECT
`
	var expectedErr = []string{
		`identifer "__y3" cannot start with two underscores at line: 6, column: 7`,
		`Required type for argument "asdf" is String, got Int at line: 10 column: 32`,
	}

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
