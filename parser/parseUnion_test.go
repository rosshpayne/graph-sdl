package parser

import (
	"fmt"
	"testing"

	"github.com/graph-sdl/lexer"
)

func TestUnion1(t *testing.T) {

	input := `
union SearchResult =| Photo | Person

type Person {
  name: String
  age: Int
}

type Photo {
  height: Int
  width: Int
}

type SearchQuery {
  firstSearchResult: SearchResult
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
