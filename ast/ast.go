package ast

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/graph-sdl/token"
)

// =================  Valuer =================================

type Valuer interface {
	ValueNode()
	String() string
	//	Exists() bool
}

// =================  InputValue =================================

// input values used for "default values" in arguments in type and field arguments and input objecs.
type InputValue_ struct {
	Value Valuer //  IV:type|value = assert type to determine InputValue_'s type via dTString
	Loc   *Loc_
}

//func (iv *InputValue_) InputValueNode() {}

func (iv *InputValue_) String() string {
	switch iv.Value.(type) {
	case *RawString_:
		return token.RAWSTRINGDEL + iv.Value.String() + token.RAWSTRINGDEL //+ "-" + iv.dTString()
	case *String_:
		return token.STRINGDEL + iv.Value.String() + token.STRINGDEL //+ "-" + iv.dTString() + iv.Loc.String()
	}
	if iv.Value == nil { // interface is not populated with concrete value
		return ""
	}
	return iv.Value.String() //+ "-" + iv.dTString()
}

func (iv *InputValue_) AtPosition() string {
	return iv.Loc.String()
}

// dataTypeString - prints the datatype of the input value
func (iv *InputValue_) dTString() string {
	switch iv.Value.(type) {
	case Int_:
		return token.INT
	case Float_:
		return token.FLOAT
	case Bool_:
		return token.BOOLEAN
	case *String_:
		return token.STRING
	case *RawString_:
		return token.STRING
	case *Scalar_:
		return token.SCALAR
	case *EnumValue_:
		return token.ENUM
	case *Object_:
		return token.OBJECT
	case *Input_:
		return token.INPUT // input specification used as type in argument
	case ObjectVals:
		return token.INPUT // actual instance of input specification used as default value in argument
	case List_:
		return "xxList"
	case Null_:
		return token.NULL
	}
	return "NoTypeFound"
}

// dataTypeString - prints the datatype of the input value

func (iv *InputValue_) isType() TypeFlag_ {
	switch iv.Value.(type) {
	case Int_:
		return INT
	case Float_:
		return FLOAT
	case Bool_:
		return BOOLEAN
	case *String_:
		return STRING
	case *RawString_:
		return STRING
	case *Scalar_:
		return SCALAR
	case *EnumValue_:
		return ENUM
	case *Object_:
		return OBJECT
	case *Input_:
		return INPUT
	case Null_:
		return NULL
	case ObjectVals:
		return INPUT
	case List_:
		return LIST
	}
	return NA
}

// dataTypeString - prints the datatype of the type specification
func (t *Type_) isType() TypeFlag_ {
	switch t.Name.String() {
	//
	// system scalars
	//
	case token.INT:
		return INT
	case token.FLOAT:
		return FLOAT
	case token.STRING, token.RAWSTRING:
		return STRING
	case token.BOOLEAN:
		return BOOLEAN
	case token.NULL:
		return NULL
	// case token.ID:
	// 	return ID
	default:
		//
		// non-standard defined types
		//
		if t.AST != nil {
			switch t.AST.(type) {
			case *Object_:
				return OBJECT
			case *Interface_:
				return INTERFACE
			case *Enum_:
				return ENUM
			case *Input_:
				return INPUT
			case *Union_:
				return UNION
			case *Scalar_:
				return SCALAR

			}
			return NA
		}
	}
	return NA
}

func (t *Type_) isScalar() bool {

	switch t.isType() {
	case INT, FLOAT, STRING, BOOLEAN:
		return true
	default:
		if _, ok := t.AST.(*Scalar_); ok {
			return true
		}
		return false
	}

}

// ================= Input Value scalar datatypes ===================

type Null_ bool // moved from Scalar to it's own type. No obvious reason why - no obvious advantage at this stage

func (n Null_) ValueNode() {}
func (n Null_) String() string {
	if n == false {
		return ""
	}
	return "null"
}

type Int_ string //int

func (i Int_) ValueNode() {}
func (i Int_) String() string {
	//return strconv.FormatInt(int64(i), 10)
	return string(i)
}

type Float_ string //float64

func (f Float_) ValueNode() {}
func (f Float_) String() string {
	return string(f)
	//return strconv.FormatFloat(float64(f), 'G', -1, 64)
}

type String_ string

func (s *String_) ValueNode() {}
func (s *String_) String() string {
	return string(*s)
}

type RawString_ string

func (s *RawString_) ValueNode() {}
func (s *RawString_) String() string {
	return string(*s)
}

type Bool_ string //bool

func (b Bool_) ValueNode() {}
func (b Bool_) String() string {
	//return strconv.FormatBool(bool(i))
	return string(b)
}

// Enum

// type Enum_ Name_

// func (e Enum_) ValueNode() {}

// func (e Enum_) Valid(s string) error {
// 	if _, err := validateName(s); err != nil {
// 		return err
// 	}
// 	if e == "true" || e == "false" || e == "null" {
// 		return fmt.Errorf("Enum, [%s] cannot be true false null", s)
// 	}
// 	return nil
// }

// func (e *Enum_) Assign(s string) {
// 	s_ := Enum_(Name_(s))
// 	e = &s_
// }

// func (e Enum_) String() string {
// 	return string(Name_(e))
// }

// List

type List_ []*InputValue_

func (l List_) ValueNode() {}
func (l List_) String() string {
	var s strings.Builder
	s.WriteString("[")
	for _, v := range l {
		s.WriteString(v.String() + " ")
	}
	s.WriteString("] ")
	return s.String()
}
func (l List_) Exists() bool {
	if len(l) > 0 {
		return true
	}
	return false
}

// type Input_ struct {
// 	Desc string
// 	Name_
// 	Directives_
// 	InputValueDefs // []*InputValueDef
// }
// type InputValueDef struct {
// 	Desc string
// 	Name_
// 	Type       *Type_   	// ** argument type specification
// 	DefaultVal *InputValue_ // ** input value(s) type(s)
// 	Directives_
// }
// type InputValue_ struct {
// 	Value Valuer //  IV:type|value = assert type to determine InputValue_'s type
// 	Loc   *Loc_
// // }
// type Type_ struct {
// 	Constraint byte          // each on bit from right represents not-null constraint applied e.g. in nested list type [type]! is 00000010, [type!]! is 00000011, type! 00000001
// 	AST        TypeDefiner // AST instance of type. WHen would this be used??. Used for non-Scalar types. AST in cache(typeName), then in Type_(typeName). If not in Type_, check cache, then DB.
// 	Depth      int           // depth of nested List e.g. depth 2 is [[type]]. Depth 0 implies non-list type, depth > 0 is a list type
// 	Name_                    // type name. inherit AssignName()
// }

func (l List_) ValidateListValues(iv *Type_, d *int, maxd *int, err *[]error) {
	reqType := iv.isType() // INT, FLOAT, OBJECT, PET, MEASURE etc            note: OBJECT is for specification of a type, OBJECTVAL is an object literal for input purposes
	//
	// for each element in the LIST
	///
	*d++ // current depth of [ in [[[]]]
	if *d > *maxd {
		*maxd = *d
	}
	for _, v := range l { // []*InputValue_ // Measure items {Name_, InputValue_}
		// what is the type of the list element. Scalar, another LIST, a OBJECT
		switch in := v.Value.(type) {

		case List_:
			// maxd records maximum depth of list(d=1) [] list of lists [[]](d=2) = [[][][][]] list of lists of lists (d=3) [[[]]] = [[[][][]],[[][][][]],[[]]]
			in.ValidateListValues(iv, d, maxd, err)
			*d--

		case ObjectVals:
			// default values in input object form { name:value name:value ... }: []*ArgumentT type ArgumentT: struct {Name_, Value *InputValue_}
			// reqType is the type of the input object  - which defines the name and associated type for each item in the { }
			in.ValidateInputObjectValues(iv, err)

		default:
			// check the item - this is matched against the type specification for the list ie. [type]
			if t := v.isType(); t != reqType {
				if v.isType() == NULL {
					if iv.Constraint>>uint(iv.Depth-*d)&1 == 1 { // is not-null constraint set
						*err = append(*err, fmt.Errorf(`List cannot contain NULLs %s`, v.AtPosition()))
					}
				} else {
					*err = append(*err, fmt.Errorf(`Required type "%s", got "%s" %s`, reqType, t, v.AtPosition()))
				}
			}
		}
	}
}

// ========== Directives ================

// Directives[Const]
// 		Directive[?Const]list
// Directive[Const] :
// 		@ Name Arguments[?Const]opt ...
// used as type for argument into parseFragment(f DirectiveAppender)
//  called using .parseDirectives(stmt) . where stmt has embedded DirectiveT field as anonymous
type DirectiveAppender interface {
	AppendDirective(s *DirectiveT) error
	//AssignLoc(loc *Loc_)
}

type DirectiveT struct {
	Name_
	Arguments_ // inherit Arguments field and d.Arguments d.AppendArgument()
	//??	Loc_       // inherit AssignLoc
}

func (d *DirectiveT) String() string {
	return "@" + d.Name_.String() + d.Arguments_.String()
}

type Directives_ struct {
	Directives []*DirectiveT
}

func (d *Directives_) AppendDirective(s *DirectiveT) error {
	for _, v := range d.Directives {
		if v.Name_.String() == s.Name_.String() {
			loc := s.Name_.Loc
			return fmt.Errorf(`Duplicate Directive name "%s" at line: %d, column: %d`, s.Name_.String(), loc.Line, loc.Column)
		}
	}
	d.Directives = append(d.Directives, s)
	return nil
}

func (d *Directives_) String() string {
	var s strings.Builder
	for _, v := range d.Directives {
		s.WriteString(v.String() + " ")
	}
	return s.String()
}

// =========== Loc_ =============================

type Loc_ struct {
	Line   int
	Column int
}

func (l Loc_) String() string {
	return "at line: " + strconv.Itoa(l.Line) + " " + "column: " + strconv.Itoa(l.Column)
	//return "" + strconv.Itoa(l.Line) + " " + strconv.Itoa(l.Column) + "] "
}

// ============== NameI  ========================

type NameAssigner interface {
	AssignName(name string, loc *Loc_, errS *[]error)
}

// ===============  Name_  =========================

type NameValue_ string

func (n NameValue_) String() string {
	return string(n)
}

func (a NameValue_) Equals(b NameValue_) bool {
	return string(a) == string(b)
}

type Name_ struct {
	Name NameValue_
	Loc  *Loc_
}

func (n Name_) String() string {
	return string(n.Name)
}

func (a Name_) Equals(b Name_) bool {
	return a.Name.Equals(b.Name)
}

func (n Name_) AtPosition() string {
	return n.Loc.String()
}

func (n Name_) Exists() bool {
	if len(n.Name) > 0 {
		return true
	}
	return false
}

func (n *Name_) AssignName(s string, loc *Loc_, errS *[]error) {
	n.Loc = loc
	validateName(s, errS, loc)
	n.Name = NameValue_(s)
}

// ======== Document ===================================

type Document struct {
	Statements    []TypeDefiner
	StatementsMap map[NameValue_]TypeDefiner
	ErrorMap      map[NameValue_][]error
}

func (d Document) String() string {
	var s strings.Builder
	tc = 2

	for _, iv := range d.StatementsMap {
		s.WriteString(iv.String())
	}
	return s.String()
}

// ======== type statements ==========

type TypeDefiner interface {
	TypeSystemNode()
	TypeName() NameValue_
	CheckUnresolvedTypes(unresolved UnresolvedMap) // while not all Types contain nested types that need to be resolved e.g scalar must still include this method
	String() string
}

var tc = 2

// ======================================================

var blank string = ""
var errNameChar string = "Invalid character in identifer at line: %d, column: %d"
var errNameBegin string = "identifer [%s] cannot start with two underscores at line: %d, column: %d"

func validateName(name string, errS *[]error, loc *Loc_) {
	// /[_A-Za-z][_0-9A-Za-z]*/
	var err error
	if len(name) == 0 {
		err = fmt.Errorf("Error: zero length name passed to validateName")
		*errS = append(*errS, err)
		return
	}

	ch, _ := utf8.DecodeRuneInString(name[:1])
	if unicode.IsDigit(ch) {
		err = fmt.Errorf("identifier cannot start with a number at line: %d, column: %d", loc.Line, loc.Column)
		*errS = append(*errS, err)
	}

	for i, v := range name {
		switch i {
		case 0:
			if !(v == '_' || (v >= 'A' || v <= 'Z') || (v >= 'a' && v <= 'z')) {
				err = fmt.Errorf(errNameChar, loc.Line, loc.Column)
				*errS = append(*errS, err)
			}
		default:
			if !((v >= '0' && v <= '9') || (v >= 'A' || v <= 'Z') || (v >= 'a' && v <= 'z') || v == '_') {
				err = fmt.Errorf(errNameChar, loc.Line, loc.Column)
				*errS = append(*errS, err)
			}
		}
		if err != nil {
			break
		}
	}

	if len(name) > 1 && name[:2] == "__" {
		err = fmt.Errorf(errNameBegin, name, loc.Line, loc.Column)
		*errS = append(*errS, err)
	}
}
