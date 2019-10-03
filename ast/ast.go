package ast

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/graph-sdl/token"
)

// =================  ValueI =================================

type ValueI interface {
	ValueNode()
	String() string
	//	Exists() bool
}

// =================  InputValue =================================

// input values used for "default values" in arguments in type and field arguments and input objecs.
// type InputValue_ struct {
// 	Value ValueI //  IV:type|value = assert type to determine InputValue_'s type.
// 	//   scalars stored as a string, true/false as boolean value. Should support data types;Enum, List, Object,null,true, false
// 	//	inputValueType_ //  Based on Token metadata in token Value - scalar types, List, Object, Enum, Null. No need for Loc.
// 	//Type Type_
// 	Loc *Loc_
// }

type InputValue_ struct {
	Value ValueI //  IV:type|value = assert type to determine InputValue_'s type.
	Loc   *Loc_
}

func (iv *InputValue_) InputValueNode() {}

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

// dataTypeString - prints the datatype of the input value
func (iv *InputValue_) dTString() string {
	switch iv.Value.(type) {
	case Int_:
		return token.INT + iv.Loc.String()
	case Float_:
		return token.FLOAT
	case Bool_:
		return token.BOOLEAN
	case *String_:
		return token.STRING
	case *RawString_:
		return token.STRING
	case List_:
		return token.LIST
	case *EnumValue_:
		return token.ENUM
	case QObject_:
		return token.OBJECT
	case Null_:
		return token.NULL
	}
	return "NoTypeFound"
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

// ========= Argument ==========

type ArgumentI interface {
	String() string
	AppendArgument(s *ArgumentT)
}

type ArgumentT struct {
	//( name:value )
	Name_
	Value *InputValue_ // could use string as this value is mapped directly to get function - at this stage we don't care about its type maybe?
}

func (a *ArgumentT) String(last bool) string {
	if last {
		return a.Name_.String() + ":" + a.Value.String()
	}
	return a.Name_.String() + ":" + a.Value.String() + " "
}

type Arguments_ struct {
	Arguments []*ArgumentT
}

func (a *Arguments_) AppendArgument(ss *ArgumentT) {
	a.Arguments = append(a.Arguments, ss)
}

func (a *Arguments_) String() string {
	var s strings.Builder
	if len(a.Arguments) > 0 {
		s.WriteString("(")
		for i, v := range a.Arguments {
			s.WriteString(v.String(i == len(a.Arguments)-1))
		}
		s.WriteString(")")
		return s.String()
	}
	return ""
}

// ========== Directives ================

// Directives[Const]
// 		Directive[?Const]list
// Directive[Const] :
// 		@ Name Arguments[?Const]opt ...
// used as type for argument into parseFragment(f DirectiveI)
//  called using .parseDirectives(stmt) . where stmt has embedded DirectiveT field as anonymous
type DirectiveI interface {
	AppendDirective(s *DirectiveT)
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

func (d *Directives_) AppendDirective(s *DirectiveT) {
	d.Directives = append(d.Directives, s)
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

type NameI interface {
	AssignName(name string, loc *Loc_, errS *[]error)
}

// ===============  Name_  =========================

type NameValue_ string

func (n NameValue_) String() string {
	return string(n)
}

type Name_ struct {
	Name NameValue_
	Loc  *Loc_
}

func (n Name_) String() string {
	return string(n.Name)
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
	Statements []TypeSystemDef
}

func (d Document) String() string {
	var s strings.Builder
	tc = 2

	for _, iv := range d.Statements {
		s.WriteString(iv.String())
	}
	return s.String()
}

// ======== type statements ==========

type TypeSystemDef interface {
	TypeSystemNode()
	TypeName() NameValue_
	CheckUnresolvedTypes(unresolved *[]Name_)
	String() string
}

type TypeExtDef interface {
	TypeExtNode()
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
