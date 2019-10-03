package token

type TokenType string
type TokenCat string

// type Token_ struct {
// 	Type TokenType
// 	Cat  TokenCat
// }

const (
	IDENT TokenType = "IDENT"
)
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// GQL Input Values types
	INT       = "INT"    // 1343456
	FLOAT     = "FLOAT"  // 3.42
	STRING    = "STRING" // contents between " or """
	RAWSTRING = "RAWSTRING"
	NULL      = "NULL"
	ENUM      = "ENUM"
	LIST      = "LIST"
	BOOLEAN   = "Boolean"
	OBJECT    = "Object"

	// Category
	VALUE    = "VALUE"
	NONVALUE = "NONVALUE"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	// Punctuator :: one of ! $ ( ) ... : = @ [ ] { | }
	COMMA      = ","
	SEMICOLON  = ";"
	COLON      = ":"
	COMMENT    = "#"
	UNDERSSCRE = "_"
	DOLLAR     = "$"
	ATSIGN     = "@"
	AND        = "&"
	BAR        = "|"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	EXPAND = "..."
	// delimiters
	RAWSTRINGDEL = `"""`

	STRINGDEL = `"`

	BOM = "BOM"

	// Keywords
	TYPE       = "TYPE"
	IMPLEMENTS = "IMPLEMENTS"
	INTERFACE  = "INTERFACE"
	UNION      = "UNION"
	ON         = "ON"
	TRUE       = "TRUE"
	FALSE      = "FALSE"
	INPUTOBJS  = "INPUTOBJS"
	EXTEND     = "EXTEND"
)

type Pos struct {
	Line int
	Col  int
}

// Token is exposed via token package so lexer can create new instanes of this type as required.
type Token struct {
	Cat          TokenCat
	Type         TokenType
	IsScalarType bool
	Literal      string // string value of token - rune, string, int, float, bool
	Loc          Pos    // start position of token
	Illegal      bool
}

var keywords = map[string]struct {
	Type         TokenType
	Cat          TokenCat
	IsScalarType bool
}{
	"Int":     {INT, NONVALUE, true},
	"Float":   {FLOAT, NONVALUE, true},
	"String":  {STRING, NONVALUE, true},
	"Boolean": {BOOLEAN, NONVALUE, true},
	//	"ID":      {ID, NONVALUE, true},
	"enum":       {ENUM, NONVALUE, true},
	"on":         {ON, NONVALUE, false},
	"type":       {TYPE, NONVALUE, false},
	"null":       {NULL, VALUE, false},
	"true":       {TRUE, VALUE, false},
	"false":      {FALSE, VALUE, false},
	"union":      {UNION, NONVALUE, false},
	"implements": {IMPLEMENTS, NONVALUE, false},
	"interface":  {INTERFACE, NONVALUE, false},
	"input":      {INPUTOBJS, NONVALUE, false},
	"extend":     {EXTEND, NONVALUE, false},
}

func LookupIdent(ident string) (TokenType, TokenCat, bool) {
	if tok, ok := keywords[ident]; ok {
		return tok.Type, tok.Cat, tok.IsScalarType
	}
	return IDENT, NONVALUE, false
}
