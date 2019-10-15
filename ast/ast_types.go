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

func (tf TypeFlag_) String() string {
	switch tf {
	case INT:
		return token.INT
	case FLOAT:
		return token.FLOAT
	case BOOLEAN:
		return token.BOOLEAN
	case STRING:
		return token.STRING
	case RAWSTRING:
		return token.STRING
	case SCALAR:
		return token.SCALAR
	case ENUM:
		return token.ENUM
	case OBJECT:
		return token.OBJECT
	case INPUT: // aka INPUTOBJ
		return token.INPUT
	// case INPUTOBJ:
	// 	return "INPUTOBJ
	case LIST:
		return token.LIST
	case NULL:
		return token.NULL
	}
	return "NoTypeFound"
}

const (
	//  input value types
	SCALAR TypeFlag_ = 1 << iota
	INT
	FLOAT
	BOOLEAN
	STRING
	RAWSTRING
	//
	NULL
	OBJECT
	ENUM
	INPUT
	LIST
	INTERFACE
	UNION
	// error - not available
	NA
)

// func isScalar(n TypeSystemDef) bool {
// 	switch n.isType {
// 	case SCALAR, INT, FLOAT, BOOLEAN, STRING, RAWSTRING: // TODO: add ID
// 		return true
// 	default:
// 		return false
// 	}
// }

// func fetchTypeFlag(n TypeI_) TypeFlag_ {
// 	// output non-Scalar Go types
// 	switch n.(type) {
// 	case *Object_:
// 		return OBJECT
// 	case *Enum_:
// 		return ENUM
// 	case *Interface_:
// 		return INTERFACE
// 	case *Union_:
// 		return UNION
// 	default:
// 		return NA
// 	}
// }

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
	// determine inputType from t.Name
	if t.isScalar() {
		return true
	}
	switch t.isType() {
	case ENUM, INPUT:
		return true
	default:
		return false
	}
}

// IsOutputType(type)
//	If type is a List type or Non‐Null type:
//		 Let unwrappedType be the unwrapped type of type.
//			Return IsOutputType(unwrappedType)
//	If type is a Scalar, Object, Interface, Union, or Enum type:
//		Return true
//	Return false

func IsOutputType(t *Type_) bool {
	if t.isScalar() {
		return true
	}
	switch t.isType() {
	case ENUM, OBJECT, INTERFACE, UNION:
		return true
	default:
		return false
	}
}

// ============================ Type_ ======================

type Type_ struct {
	Constraint byte          // each on bit from right represents not-null constraint applied e.g. in nested list type [type]! is 00000010, [type!]! is 00000011, type! 00000001
	AST        TypeSystemDef // AST instance of type. WHen would this be used??. Used for non-Scalar types. AST in cache(typeName), then in Type_(typeName). If not in Type_, check cache, then DB.
	Depth      int           // depth of nested List e.g. depth 2 is [[type]]. Depth 0 implies non-list type, depth > 0 is a list type
	Name_                    // type name. inherit AssignName()
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

func (t Type_) TypeName() string {
	return t.Name.String()
}

func (a *Type_) Equals(b *Type_) bool {
	return a.Name_.String() == b.Name_.String() && a.Constraint == b.Constraint && a.Depth == b.Depth
}

// ==================== interfaces ======================

type FieldI interface {
	AppendField(f_ *Field_) error
}

type FieldArgI interface {
	AppendField(f_ *InputValueDef, unresolved *[]error)
}

// ========= Argument ==========

type ArgumentI interface {
	String() string
	AppendArgument(s *ArgumentT)
}

type ArgumentT struct {
	//( name:value )
	Name_
	Value *InputValue_
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

// ================ ObjectVal ====================
// used as input values in Input Value Type

type ObjectVals []*ArgumentT

func (o ObjectVals) TypeSystemNode() {}
func (o ObjectVals) ValueNode()      {}
func (o ObjectVals) String() string {
	var s strings.Builder
	s.WriteString("{")
	for _, v := range o {
		s.WriteString(v.Name_.String() + ":" + v.Value.String() + " ")
	}
	s.WriteString("} ")
	return s.String()
}
func (o ObjectVals) Exists() bool {
	if len(o) > 0 {
		return true
	}
	return false
}

// type Type_ struct {
// 	Constraint byte          // each on bit from right represents not-null constraint applied e.g. in nested list type [type]! is 00000010, [type!]! is 00000011, type! 00000001
// 	AST        TypeSystemDef // AST instance of type. WHen would this be used??. Used for non-Scalar types. AST in cache(typeName), then in Type_(typeName). If not in Type_, check cache, then DB.
// 	Depth      int           // depth of nested List e.g. depth 2 is [[type]]. Depth 0 implies non-list type, depth > 0 is a list type
// 	Name_                    // type name. inherit AssignName()
// }
// type Input_ struct {
// 	Desc string
// 	Name_
// 	Directives_
// 	InputValueDefs // []*InputValueDef
// }

// type InputValueDef struct {  								<== an ArgumentDef
// 	Desc string
// 	Name_
// 	Type       *Type_   	// ** argument type specification   	<==== required type for argument
// 	DefaultVal *InputValue_ // ** input value(s) type(s)        	<==== instance data to check against required type
// 	Directives_
// }

// 	type ArgumentT struct {
// Name_
// Value *InputValue_

func (o ObjectVals) ValidateInputObjectValues(ref *Type_, err *[]error) {
	//
	//  ref{ name:value name:value ... } -- ref is the input object type specifed for the argument and { } is the default data
	//
	refFields := make(map[NameValue_]*Type_)
	// check if default input fields has fields not in field Type, PET, MEASURE
	if ref.isType() != INPUT { // redundant check
		*err = append(*err, fmt.Errorf(`Argument "%s", type %s is not an INPUT object  %s %s`, "XXX", ref.Name, ref.AtPosition()))
		return
	}
	// get reference type AST object.
	refOV := ref.AST.(*Input_)
	// build a map of fields in reference type - which defines the types of each item in {}
	for _, v := range refOV.InputValueDefs { //
		refFields[v.Name] = v.Type
	}
	//
	// loop thru name:value pairs using the ref type spec to match against name and its associated type for each pair.
	//
	for _, v := range o { // []*ArgumentT    type ArgumentT struct { Name_, Value *InputValue_}  type InputValue_ struct {Value ValueI,Loc *Loc_}
		//    ValueI populated by parser.parseInputValue_(): ast.Int_, ast.Flaat_, ast.List_, ast.ObjectVals, ast.EnumValue_ etc
		if reftype, ok := refFields[v.Name]; !ok {
			*err = append(*err, fmt.Errorf(`field %s does not exist in type %s  %s`, v.Name, ref.TypeName(), v.AtPosition()))

		} else {
			if v.Value.isType() != reftype.isType() {
				if (reftype.Depth > 0) && v.Value.isType() != LIST {
					*err = append(*err, fmt.Errorf(`Argument "%s", has type %s should be  %s %s`, ref.Name, v.Value.isType(), reftype.isType(), v.Value.AtPosition()))
					return
				}
			}
			// look at argument value type as it may be a list or another input object type
			switch inobj := v.Value.Value.(type) {

			case List_:
				if reftype.Depth == 0 {
					*err = append(*err, fmt.Errorf(`Field, %s, is not a LIST type but input data is  %s`, v.Name, v.AtPosition()))
				}
				// maxd records maximum depth of list(d=1) [] list of lists [[]](d=2) = [[][][][]] list of lists of lists (d=3) [[[]]] = [[[][][]],[[][][][]],[[]]]
				d := 0
				maxd := 0
				inobj.ValidateListValues(reftype, &d, &maxd, err)
				d--
				if maxd != reftype.Depth {
					*err = append(*err, fmt.Errorf(`Argument "%s", nested List type depth different reqired %d, got %d %s`, v.Name, reftype.Depth, maxd, v.AtPosition()))
				}

			case ObjectVals:
				inobj.ValidateInputObjectValues(reftype, err)
			}
		}

	}

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
func (o *Object_) ValueNode()      {}
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
			*err = append(*err, fmt.Errorf(`Field "%s" type "%s", is not an output type %s`, v.Name_, v.Type.Name, v.Type.Name_.AtPosition()))
		}
	}

}

func (f *Object_) CheckIsInputType(err *[]error) {
	for _, v := range f.FieldSet {
		for _, p := range v.ArgumentDefs {
			if !IsInputType(p.Type) {
				*err = append(*err, fmt.Errorf(`Argument "%s" type "%s", is not an input type %s`, p.Name_, p.Type.Name, p.Type.Name_.AtPosition()))
			}
			//	_ := p.DefaultVal.isType() // e.g. scalar, int | List
		}
	}
}

// type InputValue_ struct {
// 	Value ValueI //  IV:type|value = assert type to determine InputValue_'s type
// 	Loc   *Loc_
// // }
// type Type_ struct {
// 	Constraint byte          // each on bit from right represents not-null constraint applied e.g. in nested list type [type]! is 00000010, [type!]! is 00000011, type! 00000001
// 	AST        TypeSystemDef // AST instance of type. WHen would this be used??. Used for non-Scalar types. AST in cache(typeName), then in Type_(typeName). If not in Type_, check cache, then DB.
// 	Depth      int           // depth of nested List e.g. depth 2 is [[type]]. Depth 0 implies non-list type, depth > 0 is a list type
// 	Name_                    // type name. inherit AssignName()
// }
// Type       *Type_
// DefaultVal *InputValue_
func (f *Object_) CheckInputValueType(err *[]error) {
	for _, v := range f.FieldSet {

		// for each field in the object check if it has any default values to check
		// type Input_ struct {                                       <== Input Object
		// 	Desc string
		// 	Name_
		// 	Directives_
		// 	InputValueDefs // []*InputValueDef                          <== fields of input object
		// }
		// type Field_ struct {
		// 	Desc string
		// 	Name_
		// 	ArgumentDefs InputValueDefs //[]*InputValueDef      		<== arguments in field in object
		// 	// :
		// 	Type *Type_
		// 	Directives_
		// }
		// type InputValueDef struct {  								<== an ArgumentDef
		// 	Desc string
		// 	Name_
		// 	Type       *Type_   	// ** argument type specification   	<==== required type for argument
		// 	DefaultVal *InputValue_ // ** input value(s) type(s)        	<==== instance data to check against required type
		// 	Directives_
		// }
		// type InputValue_ struct {
		// 	Value ValueI //  IV:type|value = assert type to determine InputValue_'s type via dTString
		// 	Loc   *Loc_
		// }
		// a.Type is required type -  check against a.DefaultVal.Value.isType()
		for _, a := range v.ArgumentDefs { // go thru each of the argument field objects [] {} scalar

			if a.DefaultVal != nil {

				switch defvalObj := a.DefaultVal.Value.(type) {

				case List_: // [ "ads", "wer" ]
					if a.Type.Depth == 0 { // required type is not a LIST
						*err = append(*err, fmt.Errorf(`Argument "%s", type is not a list but default value is a list %s`, a.Name_, a.DefaultVal.AtPosition()))
						return
					}
					var d int = 0
					var maxd int
					defvalObj.ValidateListValues(a.Type, &d, &maxd, err) // a.Type is the list data type.
					//
					if maxd != a.Type.Depth {
						*err = append(*err, fmt.Errorf(`Argument "%s", nested List type depth different reqired %d, got %d %s`, a.Name_, a.Type.Depth, maxd, a.DefaultVal.AtPosition()))
					}
					//
				case ObjectVals:
					// { x: "ads", y: 234 }
					defvalObj.ValidateInputObjectValues(a.Type, err)

				default:
					if a.DefaultVal.isType() == NULL {
						if a.Type.Constraint>>uint(a.Type.Depth)&1 == 1 { // not null constraint set
							*err = append(*err, fmt.Errorf(`Value cannot be NULL %s`, a.DefaultVal.AtPosition()))
						}
					} else {
						// check types match
						if a.DefaultVal.isType() != a.Type.isType() {
							*err = append(*err, fmt.Errorf(`Required type "%s", got "%s" %s`, a.Type.isType(), a.DefaultVal.isType(), a.DefaultVal.AtPosition()))
						}
					}
				}
			}
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
	ArgumentDefs InputValueDefs //[]*InputValueDef
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
	if !f.Type.isScalar() && f.Type.AST == nil {
		if obj, ok := CacheFetch(f.Type.Name); !ok {
			*unresolved = append(*unresolved, f.Type.Name_)
		} else {
			f.Type.AST = obj
		}
	}
	//
	f.ArgumentDefs.CheckUnresolvedTypes(unresolved)
}

// use following method to override the promoted methods from Name_ and Directives_ fields. Forces use of Name_ method.
func (f *Field_) AssignName(s string, loc *Loc_, unresolved *[]error) {
	f.Name_.AssignName(s, loc, unresolved) // assign Name_{Name, Loc} and addErr if error found
}

func (f *Field_) String() string {
	var encl [2]token.TokenType = [2]token.TokenType{token.LPAREN, token.RPAREN}
	var s strings.Builder
	s.WriteString("\n" + f.Name_.String())
	s.WriteString(f.ArgumentDefs.String(encl))
	s.WriteString(" : ")
	s.WriteString(f.Type.String())
	s.WriteString(f.Directives_.String())
	return s.String()
}

func (f *Field_) AppendField(f_ *InputValueDef, unresolved *[]error) {
	f.ArgumentDefs.AppendField(f_, unresolved)
}

// ==================== ArgumentDefs ================================
// Slice of *InputValueDef
type InputValueDefs []*InputValueDef

func (fa *InputValueDefs) AppendField(f *InputValueDef, unresolved *[]error) {
	for _, v := range *fa {
		if v.Name_.String() == f.Name_.String() { //&& v.Type.Equals(f.Type) {
			loc := f.Name_.Loc
			*unresolved = append(*unresolved, fmt.Errorf(`Duplicate input value name "%s" at line: %d, column: %d`, f.Name_, loc.Line, loc.Column))
		}
	}
	*fa = append(*fa, f)
}

func (fa *InputValueDefs) String(encl [2]token.TokenType) string {
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

func (fa InputValueDefs) CheckUnresolvedTypes(unresolved *[]Name_) {

	for _, v := range fa {
		if !v.Type.isScalar() && v.Type.AST == nil {
			if ast, ok := CacheFetch(v.Type.Name); !ok {
				*unresolved = append(*unresolved, v.Type.Name_)
			} else {
				v.Type.AST = ast
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
	if !fa.Type.isScalar() && fa.Type.AST == nil {
		if ast, ok := CacheFetch(fa.Type.Name); !ok {
			*unresolved = append(*unresolved, fa.Type.Name_)
		} else {
			fa.Type.AST = ast
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
	InputValueDefs // []*InputValueDef
}

func (e *Input_) TypeSystemNode() {}
func (e *Input_) ValueNode()      {}

func (i *Input_) TypeName() NameValue_ {
	return i.Name
}

func (u *Input_) String() string {
	var encl [2]token.TokenType = [2]token.TokenType{token.LBRACE, token.RBRACE}
	var s strings.Builder
	s.WriteString("\ninput ")
	s.WriteString(u.Name.String())
	s.WriteString(" " + u.Directives_.String())
	s.WriteString(u.InputValueDefs.String(encl))
	return s.String()
}

// ======================  Scalar_ =========================
// ScalarTypeDefinition:
//		Description-opt	input	Name	DirectivesConst-opt	InputFieldsDefinition-opt
type Scalar_ struct {
	Desc string
	Name_
	Directives_
}

func (e *Scalar_) TypeSystemNode() {}
func (e *Scalar_) ValueNode()      {}

func (i *Scalar_) TypeName() NameValue_ {
	return i.Name
}
func (i *Scalar_) CheckUnresolvedTypes(unresolved *[]Name_) {}

func (u *Scalar_) String() string {
	var s strings.Builder
	s.WriteString("\nscalar ")
	s.WriteString(u.Name.String())
	s.WriteString(" " + u.Directives_.String())
	return s.String()
}
