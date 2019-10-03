package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestInput1(t *testing.T) {

	input := `input Point2D  @ june (asdf:234){
  x: Float = 123.23 @ jun (asdf:234)
  y: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
   y1: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
     y2: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
      y3: Float = 34 @ jun (asdf:"""asdflkjslkjd""" dei:234 uio:false)
       y4: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
        y5: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
         y6: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
          y63: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
}
`

	l := lexer.New(input)
	p := New(l)
	d, errs := p.ParseDocument()
	fmt.Println(d.String())
	for _, v := range errs {
		fmt.Println(v.Error())
	}
	if compare(d.String(), input) {
		fmt.Println(trimWS(d.String()))
		fmt.Println(trimWS(input))
		t.Errorf(`*************  program.String() wrong.`)
	}
	// for v, o := range repo {
	// 	fmt.Printf(" %s, %T\n", v, o)
	// }
}

func TestInputDuplicate(t *testing.T) {

	input := `input Point2D  @ june (asdf:234){
  x: Float = 123.23 @ jun (asdf:234)
  y: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
   y1: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
     y2: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
       y: Int = 34 @ jun (asdf:"""asdflkjslkjd""")
      y3: Float = 34 @ jun (asdf:"""asdflkjslkjd""" dei:234 uio:false)
       y4: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
        y5: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
         y6: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
          y63: Float = 34 @ jun (asdf:"""asdflkjslkjd""")
}
`

	expectedErr := `Duplicate input value name "y" at line: 6, column: 8`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != 1 {
		t.Errorf(`Expect one error got %d`, len(errs))
	}
	//fmt.Println(d.String())
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
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

	expectedErr := `identifer [__y3] cannot start with two underscores at line: 6, column: 7`

	l := lexer.New(input)
	p := New(l)
	_, errs := p.ParseDocument()
	if len(errs) != 1 {
		t.Errorf(`Expect one error got %d`, len(errs))
	}
	//fmt.Println(d.String())
	for _, v := range errs {
		if v.Error() != expectedErr {
			t.Errorf(`Wrong Error got=[%q] expected [%s]`, v.Error(), expectedErr)
		}
	}
}
