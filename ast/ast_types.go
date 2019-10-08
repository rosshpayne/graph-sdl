package ast

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/graph-sdl/token"
)

// how this works
//
//  Load
//  1. parse type literals (in a file)
//  2. validate
//  3. persist to dynamodb
//
//  Part of QL execution
//  1. use rootType to probe dynamodb and build AST-type using the following structs
//  2. validate QL - build AST-QL and embed AST-Type for validation and execution
//  3  execute QL - using both ASTs
//  4  save AST-QL to dynamodb

type TypeFlag_ int

const (
	//  input value types
	SCALAR TypeFlag_ = 1 << iota
	INT
	FLOAT
	BOOLEAN
	NULL
	OBJECT
	ENUM
	INPUTOBJ
	LIST
	STRING
	RAWSTRING
	// other non-scalar types
	INTERFACE
	UNION
	// error - not available
	NA
)

func isScalar(n TypeFlag_) bool {
	switch n {
	case SCALAR, INT, FLOAT, BOOLEAN, OBJECT, STRING, RAWSTRING:
		return true
	default:
		return false
	}
}

func fetchTypeFlag(n TypeI_) TypeFlag_ {
	// output non-Scalar Go types
	switch n.(type) {
	case *Object_:
		return OBJECT
	case *Enum_:
		return ENUM
	case *Interface_:
		return INTERFACE
	case *Union_:
		return UNION
	default:
		return NA
	}
}

type TypeI_ interface {
	TypeSystemNode()
	TypeName() NameValue_
	String() string
}

// ============== maps =============================

type EnumRepo_ map[string]struct{}

// IsInputType(type)
//	 If type is a List type or Non‐Null type:
//		Let unwrappedType be the unwrapped type of type.
//			Return IsInputType(unwrappedType)
//	 If type is a Scalar, Enum, or Input Object type:
//				Return true
//	 Return false

func IsInputType(t *Type_) bool {
	if isScalar(t.TypeFlag) || t.TypeFlag == ENUM || t.TypeFlag == INPUTOBJ {
		return true
	}
	return false
}

// IsOutputType(type)
//	If type is a List type or Non‐Null type:
//		 Let unwrappedType be the unwrapped type of type.
//			Return IsOutputType(unwrappedType)
//	If type is a Scalar, Object, Interface, Union, or Enum type:
//		Return true
//	Return false

func IsOutputType(t *Type_) bool {
	if isScalar(t.TypeFlag) || t.TypeFlag == ENUM || t.TypeFlag == OBJECT || t.TypeFlag == INTERFACE || t.TypeFlag == UNION {
		return true
	}
	return false
}

// ============================ Type_ ======================

// type Type__ {
// 	Type Type_
// 	Loc *Loc
// }
type Type_ struct {
	Constraint byte      // each on bit from right represents not-null constraint applied e.g. in nested list type [type]! is 00000010, [type!]! is 00000011, type! 00000001
	TypeFlag   TypeFlag_ // Scalar (int,float,boolean,string,ID - Name_ defines the actual type e.g. Name_=Int) Object, Interface, Union, Enum, InputObj (AST contains type def)
	//	AST        TypeI_    // AST instance of type. WHen would this be used??. A type has no AST but an instance of a type does.
	Depth int // depth of nested List e.g. depth 2 is [[type]]. Depth 0 implies non-list type, depth > 0 is a list type
	Name_     // type name. inherit AssignName()
	//Value      ValueI    // default value
}

func (t Type_) String() string {
	var s strings.Builder
	for i := 0; i < t.Depth; i++ {
		s.WriteString("[")
	}
	s.WriteString(t.Name_.String())
	//s.WriteString("-" + fmt.Sprintf("%08b", t.TypeFlag))
	var (
		one byte = 1 << 0
		bit byte
	)
	var i uint
	if t.Depth == 0 {
		bit = (t.Constraint >> i) & one // show right most bit only
		if bit == 1 {
			s.WriteString("!")
		}
	} else {
		for i = 0; int(i) <= t.Depth+1; i++ {
			bit = (t.Constraint >> i) & one // show right most bit only
			if bit == 1 {
				s.WriteString("!")
			}
			if int(i) < t.Depth {
				s.WriteString("]")
			}
		}
	}
	return s.String()
}

func (a *Type_) Equals(b *Type_) bool {
	return a.Name_.String() == b.Name_.String() && a.Constraint == b.Constraint && a.Depth == b.Depth && a.TypeFlag == b.TypeFlag
}

// ==================== interfaces ======================

type FieldI interface {
	AppendField(f_ *Field_) error
}

type FieldArgI interface {
	AppendField(f_ *InputValueDef, unresolved *[]error)
}

// ================ QObject ====================
// used as input object values

type QObject_ []*ArgumentT

func (o QObject_) ValueNode() {}
func (o QObject_) String() string {
	var s strings.Builder
	s.WriteString("{")
	for _, v := range o {
		s.WriteString(v.Name_.String() + ":" + v.Value.String() + " ")
	}
	s.WriteString("} ")
	return s.String()
}
func (o QObject_) Exists() bool {
	if len(o) > 0 {
		return true
	}
	return false
}

// =================================================================
// Slice of Name_
type NameS []Name_

func (f NameS) CheckUnresolvedTypes(unresolved *[]Name_) {
	for _, v := range f {
		if _, ok := Fetch(v.Name); !ok {
			*unresolved = append(*unresolved, v)
		}
	}
}

// ========================= Object_ ===============================
// object definition:
// type Person {
//   name: String
//   age: Int
//   picture: Url
// }
// Description-opt TYPE Name  ImplementsInterfaces-opt Directives-Const-opt  FieldsDefinition-opt
//		{FieldDefinition-list}
//         FieldDefinition:
//			Description-opt Name ArgumentsDefinition- opt : Type Directives-Con
type Object_ struct {
	Desc        string
	Name_             // inherits AssignName  (from Name_). Overidden
	Implements  NameS //TODO  = create type NameS []*Name_ and add method AppendField to NameS and then embedded this type in Object_ struct
	Directives_       // inherits AssignName  (from Name_) + others. Overidden
	//	Fields      FieldSet // TODO - embed anonymous this FieldSet in Object_
	FieldSet
}

func (o *Object_) TypeSystemNode() {}
func (o *Object_) TypeName() NameValue_ {
	return o.Name
}

func (f *Object_) CheckImplements(err *[]error) {
	for _, v := range f.Implements {
		// check name represents a interface type in repo
		if itf, ok, str := fetchInterface(v); !ok {
			*err = append(*err, errors.New(fmt.Sprintf(str)))
		} else {
			// check object implements the interface
			satisfied := make(map[NameValue_]bool)
			for _, v := range itf.FieldSet {
				satisfied[v.Name] = false
			}
			for _, ifn := range itf.FieldSet { // interface fields
				for _, fn := range f.FieldSet { // object fields
					if ifn.Name_.String() == fn.Name_.String() {
						if ifn.Type.Equals(fn.Type) {
							satisfied[fn.Name] = true
						}
					}
				}
				var s strings.Builder
				for k, v := range satisfied {
					if !v {
						s.WriteString(` "`)
						s.WriteString(k.String())
						s.WriteString(`"`)
					}
				}
				if len(s.String()) > 0 {
					*err = append(*err, fmt.Errorf(`Object type "%s" does not implement interface "%s", missing%s`, f.Name_, itf.Name_, s))
				}
			}
		}
	}
}

func (f *Object_) CheckUnresolvedTypes(unresolved *[]Name_) {
	f.FieldSet.CheckUnresolvedTypes(unresolved)
	f.Implements.CheckUnresolvedTypes(unresolved)
}

func (f *Object_) CheckIsOutputType(err *[]error) {
	for _, v := range f.FieldSet {
		if !IsOutputType(v.Type) {
			//loc := v.Name_.Loc
			*err = append(*err, fmt.Errorf(`Field "%s" type "%s", is not an output type %s`, v.Name_, v.Type.Name, v.Type.Name_.AtPosition()))
		}
	}

}

// use following method to disambiguate the promoted AssignName method from Name_ and Directives_ fields. Forces use of Name_ method.
func (f *Object_) AssignName(s string, loc *Loc_, unresolved *[]error) {
	f.Name_.AssignName(s, loc, unresolved) // assign Name_{Name, Loc} and addErr if error found
}

func (f *Object_) String() string {
	var s strings.Builder
	s.WriteString("\ntype " + f.Name_.String())
	for i, v := range f.Implements {
		if i == 0 {
			s.WriteString(" implements ")
		}
		if i > 0 {
			s.WriteString(" & ")
		}
		s.WriteString(v.String())
	}
	s.WriteString(" " + f.Directives_.String())
	s.WriteString(f.FieldSet.String())

	return s.String()
}

// ================ FieldSet =================

type FieldSet []*Field_

func (f *FieldSet) String() string {
	var s strings.Builder
	for i, v := range *f {
		if i == 0 {
			s.WriteString("{")
		}
		s.WriteString(v.String())
		if i == len(*f)-1 {
			s.WriteString("\n}")
		}
	}
	return s.String()
}

func (fs *FieldSet) CheckUnresolvedTypes(unresolved *[]Name_) {
	for _, v := range *fs {
		v.CheckUnresolvedTypes(unresolved)
		// if v.Type.TypeFlag == 0 { // ie. a user defined type not known to the Lexer
		// 	if typ, ok := Fetch(v.Type.Name); !ok {
		// 		*unresolved = append(*unresolved, v.Type.Name_)
		// 	} else {
		// 		// TODO - is this necessary
		// 		// if v.Type.TypeFlag = fetchTypeFlag(typ); v.Type.TypeFlag == NA {
		// 		// 	err := fmt.Errorf("Type not known [%c]", v.Type.Name_.String())
		// 		// 	*unresolved = append(*unresolved, err)
		// 		// }
		// 		v.Type.AST = typ
		// 	}
		// }
		// // check argument types
		// InputValueS.CheckUnresolvedTypes(unresolved *[]Name_)
	}
}

func (fs *FieldSet) AppendField(f_ *Field_) error {
	for _, v := range *fs {
		// check field (Name and Type) not already present
		if v.Equals(f_) {
			loc := f_.Name_.Loc
			return fmt.Errorf(`Duplicate Field name "%s" at line: %d, column: %d`, f_.Name_.String(), loc.Line, loc.Column)
		}
	}
	*fs = append(*fs, f_)
	return nil
}

// ===============================================================
type HasTypeI interface {
	AssignType(t *Type_)
}

// ==================== Field_ ================================
// instance of Object Field
// FieldDefinition
//		 Description-opt	Name	ArgumentsDefinition-opt	:	Type	Directives-opt
type Field_ struct {
	Desc string
	Name_
	InputValueS // []*InputValueDef
	// :
	Type *Type_
	Directives_
}

func (f *Field_) AssignType(t *Type_) {
	f.Type = t
}

func (a *Field_) Equals(b *Field_) bool {
	return a.Name_.Equals(b.Name_) && a.Type.Equals(b.Type)
}

func (f *Field_) CheckUnresolvedTypes(unresolved *[]Name_) {
	if f.Type == nil {
		log.Panic(fmt.Errorf("Severe Error - not expected: Field.Type is not assigned for [%s]", f.Name_.String()))
	}
	if f.Type.TypeFlag == 0 { // non-zero means its been defined in the Lexer as a Go Scalar type.
		if _, ok := Fetch(f.Type.Name); !ok {
			*unresolved = append(*unresolved, f.Type.Name_)
		} else {
			// TODO - is this necessary
			// if f.Type.TypeFlag = fetchTypeFlag(typ); f.Type.TypeFlag == NA {
			// 	err := fmt.Errorf("Type not known [%c]", f.Type.TypeFlag)
			// 	*unresolved = append(*unresolved, err)
			// }
			//f.Type.AST = nil //typ
		}
	}
	//
	f.InputValueS.CheckUnresolvedTypes(unresolved)
}

// use following method to override the promoted methods from Name_ and Directives_ fields. Forces use of Name_ method.
func (f *Field_) AssignName(s string, loc *Loc_, unresolved *[]error) {
	f.Name_.AssignName(s, loc, unresolved) // assign Name_{Name, Loc} and addErr if error found
}

func (f *Field_) String() string {
	var encl [2]token.TokenType = [2]token.TokenType{token.LPAREN, token.RPAREN}
	var s strings.Builder
	s.WriteString("\n" + f.Name_.String())
	s.WriteString(f.InputValueS.String(encl))
	s.WriteString(" : ")
	s.WriteString(f.Type.String())
	s.WriteString(f.Directives_.String())
	return s.String()
}

// ==================== InputValueS ================================
// Slice of *InputValueDef
type InputValueS []*InputValueDef

func (fa *InputValueS) AppendField(f *InputValueDef, unresolved *[]error) {
	for _, v := range *fa {
		if v.Name_.String() == f.Name_.String() && v.Type.Equals(f.Type) {
			loc := f.Name_.Loc
			*unresolved = append(*unresolved, fmt.Errorf(`Duplicate input value name "%s" at line: %d, column: %d`, f.Name_.String(), loc.Line, loc.Column))
		}
	}
	*fa = append(*fa, f)
}

func (fa *InputValueS) String(encl [2]token.TokenType) string {
	var s strings.Builder
	for i, v := range *fa {
		if i == 0 {
			//s.WriteString("\n")
			s.WriteString(string(encl[0]))
		}
		//s.WriteString("\n")
		s.WriteString(v.String())
		if i == len(*fa)-1 {
			//	s.WriteString("\n")
			s.WriteString(string(encl[1]))
		}
	}
	return s.String()
}

func (fa InputValueS) CheckUnresolvedTypes(unresolved *[]Name_) {

	for _, v := range fa {
		if v.Type.TypeFlag == 0 {
			if _, ok := Fetch(v.Type.Name); !ok {
				*unresolved = append(*unresolved, v.Type.Name_)
			} else {
				//TODO - do I need this
				// if v.Type.TypeFlag = fetchTypeFlag(typ); v.Type.TypeFlag == NA {
				// 	err := fmt.Errorf("Type not known [%c]", v.Type.TypeFlag)
				// 	*unresolved = append(*unresolved, err)
				// }
				//	v.Type.AST = nil //typ
			}
		}
	}
}

// ==================== . InputValueDef . ============================
// ArgumentsDefinition
//		(InputValueDefinitionlist)
// InputValueDefinition
//		Description-opt Name : Type DefaultValue-opt Directives-opt
type InputValueDef struct {
	Desc string
	Name_
	Type       *Type_
	DefaultVal *InputValue_
	Directives_
}

func (fa *InputValueDef) checkUnresolvedType(unresolved *[]Name_) {
	if fa.Type == nil {
		err := fmt.Errorf("Severe Error - not expected: InputValueDef.Type is not assigned for [%s]", fa.Name_.String())
		log.Panic(err)
	}
	if fa.Type.TypeFlag == 0 {
		if _, ok := Fetch(fa.Type.Name); !ok {
			*unresolved = append(*unresolved, fa.Type.Name_)
		} else {
			// TODO - do I need this?
			// if fa.Type.TypeFlag = fetchTypeFlag(typ); fa.Type.TypeFlag == NA {
			// 	err := fmt.Errorf("Type not known [%c]", fa.Type.Name_.String())
			// 	*unresolved = append(*unresolved, err)
			// }
			//fa.Type.AST = nil //typ
		}
	}
}

func (fa *InputValueDef) AssignName(input string, loc *Loc_, unresolved *[]error) {
	fa.Name_.AssignName(input, loc, unresolved)
}

func (fa *InputValueDef) AssignType(t *Type_) {
	fa.Type = t
}

func (fa *InputValueDef) String() string {
	var s strings.Builder
	s.WriteString(fa.Name_.String())
	s.WriteString(" : " + fa.Type.String() + " ")
	if fa.DefaultVal != nil {
		s.WriteString("=")
		s.WriteString(fa.DefaultVal.String())
	}
	s.WriteString(fa.Directives_.String())
	return s.String()
}

// ======================  Enum =========================

// Enum
//	Descriptio-nopt enum Name Directives-const-opt EnumValuesDefinition-opt
//		EnumValuesDefinition
//		{EnumValueDefinitionlist}
type Enum_ struct {
	Desc string
	Name_
	Directives_
	Values []*EnumValue_
}

func (e *Enum_) TypeSystemNode()                          {}
func (e *Enum_) CheckUnresolvedTypes(unresolved *[]Name_) {}

func (e *Enum_) TypeName() NameValue_ {
	return e.Name
}

func (e *Enum_) String() string {
	var s strings.Builder
	s.WriteString("enum " + e.Name_.String())
	s.WriteString(e.Directives_.String())
	for i, v := range e.Values {
		if i == 0 {
			s.WriteString("{\n")
		}
		s.WriteString(v.String() + "\n")
		if i == len(e.Values)-1 {
			s.WriteString("}\n")
		}
	}
	return s.String()
}

// ======================  EnumValue =========================

//	EnumValueDefinition
//		Description-opt EnumValue Directives-const-opt
type EnumValue_ struct {
	Desc string
	Name_
	Directives_
}

func (e *EnumValue_) ValueNode() {} // instane of InputValue_

func (e *EnumValue_) AssignName(s string, l *Loc_, unresolved *[]error) {
	e.Name_.AssignName(s, l, unresolved)
}

func (e *EnumValue_) String() string {
	var s strings.Builder
	s.WriteString(e.Name_.String())
	s.WriteString(" " + e.Directives_.String())
	return s.String()
}

// ======================  Schema =========================

type Schema struct {
	rootQuery        *Type_
	rootMutation     *Type_
	rootSubscription *Type_
}

// ======================  Interface =========================

// InterfaceTypeDefinition
//		Description-opt interface Name Directives-opt FieldsDefinition-opt
type Interface_ struct {
	Desc string
	Name_
	Directives_
	FieldSet
}

func (i *Interface_) TypeSystemNode() {}

func (i *Interface_) TypeName() NameValue_ {
	return i.Name
}

//func (i *Interface_) AssignUnresolvedTypes(repo TypeRepo) error {}
func (i *Interface_) AssignName(input string, loc *Loc_, unresolved *[]error) {
	i.Name_.AssignName(input, loc, unresolved)
}

func (i *Interface_) String() string {
	var s strings.Builder
	s.WriteString("interface ")
	s.WriteString(i.Name_.String())
	s.WriteString(" " + i.Directives_.String())
	s.WriteString(" " + i.FieldSet.String())
	return s.String()
}

// ======================  Union =========================

// InterfaceTypeDefinition
//		Description-opt interface Name Directives-opt FieldsDefinition-opt
type Union_ struct {
	Desc string
	Name_
	Directives_
	NameS // Members
}

func (u *Union_) TypeSystemNode() {}
func (u *Union_) TypeName() NameValue_ {
	return u.Name
}

// func (u *Union_) Equals(b *Union_) bool {
// 	if u.Name.Equals(b.Name) {
// 		return false
// 	}
// 	if !u.NameS.Equals(b) {
// 		return false
// 	}
// 	if !u.Directives_.Equasl(b) {
// 		return false
// 	}
// 	return true
// }

func (u *Union_) String() string {
	var s strings.Builder
	s.WriteString("\nunion ")
	s.WriteString(u.Name_.String())
	s.WriteString(" " + u.Directives_.String())
	for i, v := range u.NameS {
		if i == 0 {
			s.WriteString(token.ASSIGN)
		}
		s.WriteString(token.BAR)
		s.WriteString(" ")
		s.WriteString(v.String())
	}

	return s.String()
}

// ======================  Input_ =========================
// InputObjectTypeDefinition
//		Description-opt	input	Name	DirectivesConst-opt	InputFieldsDefinition-opt
type Input_ struct {
	Desc string
	Name_
	Directives_
	InputValueS // []*InputValueDef
}

func (e *Input_) TypeSystemNode() {}

func (i *Input_) TypeName() NameValue_ {
	return i.Name
}

func (u *Input_) String() string {
	var encl [2]token.TokenType = [2]token.TokenType{token.LBRACE, token.RBRACE}
	var s strings.Builder
	s.WriteString("\ninput ")
	s.WriteString(u.Name.String())
	s.WriteString(" " + u.Directives_.String())
	s.WriteString(u.InputValueS.String(encl))
	return s.String()
}
