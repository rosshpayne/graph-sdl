package ast

import (
	_ "errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/graph-sdl/token"
)

type TypeFlag_ uint8

// TypeFlag constants shared by GQLType & InputValue types - not all the below are applicable to each
const (
	//  input value types
	_ TypeFlag_ = iota
	ID
	INT
	FLOAT
	BOOLEAN
	STRING
	RAWSTRING
	SCALAR
	//
	NULL
	OBJECT
	ENUM
	ENUMVALUE
	INPUT
	LIST
	INTERFACE
	UNION
	//
	ILLEGAL
)

func (tf TypeFlag_) String() string {
	switch tf {
	case ID:
		return token.ID
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
	case ENUMVALUE:
		return token.ENUM
	case OBJECT:
		return token.OBJECT
	case INPUT:
		return token.INPUT
	case LIST:
		return token.LIST
	case NULL:
		return token.NULL
	case INTERFACE:
		return token.INTERFACE
	case UNION:
		return token.UNION

	}
	return token.ILLEGAL
}

// Type cache - each type and its associated AST is held in the cache.
var TyCache map[string]GQLTypeProvider

func InitCache(size int) {
	TyCache = make(map[string]GQLTypeProvider, size)
}

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

type UnresolvedMap map[Name_]*GQLtype

// Directive Locations
type DirectiveLoc uint8

// Directive Locations
const (
	_ DirectiveLoc = iota
	SCHEMA_DL
	SCALAR_DL
	OBJECT_DL
	FIELD_DEFINITION_DL
	ARGUMENT_DEFINITION_DL
	INTERFACE_DL
	UNION_DL
	ENUM_DL
	ENUM_VALUE_DL
	INPUT_OBJECT_DL
	INPUT_FIELD_DEFINITION_DL
	//
	QUERY_DL
	MUTATION_DL
	SUBSCRIPTION_DL
	FIELD_DL
	FRAGMENT_DEFINITION_DL
	FRAGMENT_SPREAD_DL
	INLINE_FRAGMENT_DL
)

var DirectiveLocationMap map[DirectiveLoc]string

func init() {
	DirectiveLocationMap = make(map[DirectiveLoc]string)

	DirectiveLocationMap = map[DirectiveLoc]string{
		SCHEMA_DL:                 "SCHEMA",
		SCALAR_DL:                 "SCALAR",
		OBJECT_DL:                 "OBJECT",
		FIELD_DEFINITION_DL:       "FIELD_DEFINITION",
		ARGUMENT_DEFINITION_DL:    "ARGUMENT_DEFINITION",
		INTERFACE_DL:              "INTERFACE",
		UNION_DL:                  "UNION",
		ENUM_DL:                   "ENUM",
		ENUM_VALUE_DL:             "ENUM_VALUE",
		INPUT_OBJECT_DL:           "INPUT_OBJECT",
		INPUT_FIELD_DEFINITION_DL: "INPUT_FIELD_DEFINITION",
		//
		QUERY_DL:               "QUERY",
		MUTATION_DL:            "MUTATION",
		SUBSCRIPTION_DL:        "SUBSCRIPTION",
		FIELD_DL:               "FIELD",
		FRAGMENT_DEFINITION_DL: "FRAGMENT_DEFINITION",
		FRAGMENT_SPREAD_DL:     "FRAGMENT_SPREAD",
		INLINE_FRAGMENT_DL:     "INLINE_FRAGMENT",
	}

}

// ============== maps =============================

// IsInputType(type)
//	 If type is a List type or Non‐Null type:
//		Let unwrappedType be the unwrapped type of type.
//			Return IsInputType(unwrappedType)
//	 If type is a Scalar, Enum, or Input Object type:
//				Return true
//	 Return false

func IsInputType(t *GQLtype) bool {
	// determine inputType from t.Name
	fmt.Println("***** IsInputType ***** ", t.Name)
	if t.IsScalar() {
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

func IsOutputType(t *GQLtype) bool {
	if t.IsScalar() {
		return true
	}
	switch t.isType() {
	case ENUM, OBJECT, INTERFACE, UNION, SCALAR:
		return true
	default:
		return false
	}
}

// ==================== interfaces ======================

type FieldAppender interface {
	AppendField(f_ *Field_) error
}

type FieldArgAppender interface {
	AppendField(f_ *InputValueDef, unresolved *[]error)
}

// ========= Argument ==========

type ArgumentAppender interface {
	String() string
	AppendArgument(s *ArgumentT)
}

// ========= ArgumentT ==========

type ArgumentT struct {
	//( name : value ) e.g. picture(size: 300): Url    where Name_ is size and Value is 300
	Name_
	Value *InputValue_
}

func (a *ArgumentT) StmtType() string {
	return ""
} // to support ql.NameI

func (a *ArgumentT) String(last bool) string {
	if last {
		return a.Name_.String() + ":" + a.Value.String()
	}
	return a.Name_.String() + ":" + a.Value.String() + " "
}

// ======================================

type ArgumentS []*ArgumentT // same as type ObjectVals []*ArgumentT

func (a ArgumentS) String() string {
	var s strings.Builder
	if len(a) > 0 {
		s.WriteString("(")
		for i, v := range a {
			s.WriteString(v.String(i == len(a)-1))
		}
		s.WriteString(")")
		return s.String()
	}
	return ""
}

// ========================================

type Arguments_ struct {
	Arguments []*ArgumentT
}

// func (a *Arguments_) CheckInputValueType(err *) {
// }

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

func (a *Arguments_) SolicitAbstractTypes(unresolved UnresolvedMap) {
	// TODO - need to have the type of Name_
	// for _, v := range a.Arguments {
	// 	if !v.Value.IsScalar() {
	// 		unresolved[v.Name_] = nil
	// 	}
	// }
}

// ================ ObjectVal ====================
// used as input values in Input Value Type

type ObjectVals []*ArgumentT

func (o ObjectVals) TypeSystemNode() {}
func (o ObjectVals) ValueNode()      {}

func (o ObjectVals) IsType() TypeFlag_ {
	return OBJECT
}

func (o ObjectVals) Type() string {
	return "ObjectVals"
}

func (a ObjectVals) AppendArgument(ss *ArgumentT) {
	a = append(a, ss)
}

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

// ValidateObjectValues compares the value of the objectVal ie. value in { name:value name:value ... } against the TYPE (ref) from the type definition Object/Input
// . e.g. {name:"Ross", age:33} ==> root type "Person" {name: String, age: Int}
func (o ObjectVals) ValidateObjectValues(ref *GQLtype, err *[]error) {
	//
	var errObj string
	fmt.Println(" ----------- ValidateObjectValues ------------------")
	fmt.Printf("ref:  %s %T  Name:  %s   List: %v\n", ref.IsType(), ref, ref.Name_, ref.IsList())
	refFields := make(map[NameValue_]*GQLtype)
	//
	// What is the reference type the objectValue value should match
	//
	switch x := ref.AST.(type) {
	case *Input_:
		for _, v := range x.InputValueDefs { //
			refFields[v.Name] = v.Type
		}
		errObj = "Argument"
	case *Object_:
		for _, v := range x.FieldSet { //
			refFields[v.Name] = v.Type
		}
		errObj = "Object"
	default:
		*err = append(*err, fmt.Errorf(`Mismatched types. The input data (object values in this case) does not match a Object or Input type. The reference type is a %s`, ref.TypeName())) //TODO location required
		return
	}
	fmt.Println("****** ", refFields)
	//
	// loop thru name:value pairs using the ref type (object or Input types) to match against name and its associated type for each pair.
	//
	for _, v := range o {

		// reference type of the object/input field
		if reftype, ok := refFields[v.Name]; !ok {
			*err = append(*err, fmt.Errorf(`field "%s" does not exist in type %s  %s`, v.Name, ref.TypeName(), v.AtPosition()))

		} else {
			// compare reference type against field  data

			//	fmt.Printf("Field, , v.Value.isType(), refType.isType2(): %s, %T %T, %s, %s, %s\n", v.Name, v.Value, reftype, v.Value.isType(), reftype.isType2(), reftype.isType()) // InputValue.isType, *GQLtype.isType()
			fmt.Println("++++++++++++")
			fmt.Printf("Field v.Name:  [%s]\n", v.Name)
			fmt.Printf("v.Value:       [%T]\n", v.Value)
			fmt.Printf("v.value.isType(): %s\n", v.Value.isType())
			fmt.Printf("reftype.isType2():  %s\n", reftype.isType2())
			fmt.Printf("reftype.isType():   %s\n", reftype.isType())
			fmt.Printf("reftype:        %T\n", reftype)
			fmt.Println("++++++++++++")
			// value == LIST isType2 == LIST isType == INT	    // LIST appropriate but no check for internal types made in ValidateListValues.
			// value == LIST isTYpe2 == INT  isType == INT		// should not be in list
			// value == INT  isType2 == LIST isType = INT       // must be in list
			// note: for reftype (ie. *Type) isType2() display LIST when it is a LIST, whereas isType() displays the base type of the LIST ie. its member type.
			//       for value (ie *InputValue) isType() displays OBJECTVALS, LIST, <SCALARS>, ENUMs of the embedded InputValueProvider
			if v.Value.isType() != reftype.isType2() { // ie. both are not LISTs or different scalar types
				// only LIST differences consider here
				// value is a LIST but ref type is not
				if v.Value.isType() == LIST && !reftype.isList() {
					*err = append(*err, fmt.Errorf(`%s "%s" from type "%s" should not be a List type %s`, errObj, v.Name, ref.Name, v.Value.AtPosition()))
					// abort any further validation on this item
					return
				} else if v.Value.isType() != LIST && reftype.isType2() == LIST {
					*err = append(*err, fmt.Errorf(`%s "%s" from type "%s" expected %s %s`, errObj, v.Name, ref.Name, reftype.isType2(), v.Value.AtPosition()))
				}
			}
			// when value not LIST check types

			if v.Value.isType() != reftype.isType() && v.Value.isType() != LIST { // List is validated ValidateListValues
				// for the purpose of this validation OBJECT and INPUT are the same
				if !(v.Value.isType() == OBJECT && reftype.isType() == INPUT) {
					//	if v.Value.isType() != LIST {
					*err = append(*err, fmt.Errorf(`%s "%s" from type "%s" expected %s got %s %s`, errObj, v.Name, ref.Name, reftype.isType(), v.Value.isType(), v.Value.AtPosition()))
				}
			}
			//
			// look at  value type as it may be a list or another object/input type
			//
			switch iv := v.Value.InputValueProvider.(type) { // y inob:Float_

			case List_:
				fmt.Println(" ----------- ValidateObjectValues --LIST----------------")
				// maxd records maximum depth of list(d=1) [] list of lists [[]](d=2) = [[][][][]] list of lists of lists (d=3) [[[]]] = [[[][][]],[[][][][]],[[]]]
				var d, maxd uint8
				iv.ValidateListValues(reftype, &d, &maxd, err)
				d--
				if maxd != reftype.Depth && reftype.Depth != 0 { // reftype.Depth == 0 check performed above
					*err = append(*err, fmt.Errorf(`Argument "%s", nested List type depth different reqired %d, got %d %s`, v.Name, reftype.Depth, maxd, v.AtPosition()))
				}

			case ObjectVals:
				fmt.Println(" ----------- ValidateObjectValues ------OBJVAL------------")
				iv.ValidateObjectValues(reftype, err)

			}
		}
	}
	//
	// check mandatory fields present - for Input types only.
	//
	var at string
	if _, ok := ref.AST.(*Input_); ok {
		for k, v := range refFields { // k Name, v *Type
			if (v.Constraint>>uint(v.Depth))&1 == 1 { // mandatory field. Check present.
				found := false
				for _, v := range o {
					if v.Name == k {
						found = true
					}
					at = v.AtPosition()
				}
				if !found {
					*err = append(*err, fmt.Errorf(`Mandatory field "%s" missing in type "%s" %s `, k, ref.TypeName(), at))
				}
			}
		}
	}
}

// =================================================================
// Slice of Name_
type NameS []Name_

// SolicitAbstractTypes is typically promoted to type that embedds the NameS type.
func (f NameS) SolicitAbstractTypes(unresolved UnresolvedMap) { //TODO rename to checkUnresolvedTypes
	//  handled by type in which NameS is nested
}

type SDLObjectInterfacer interface {
	GetSelectionSet() FieldSet
	TypeName() NameValue_
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
	Name_             // instane name of Object_ type e.g. Person, Pet. inherits fields and method, AssignName from Name_. Overidden
	Implements  NameS //TODO  = create type NameS []*Name_ and add method AppendField to NameS and then embedded this type in Object_ struct
	Directives_       // inherits AssignName  (from Name_) + others. Overidden
	//	Fields      FieldSet // TODO - embed anonymous this FieldSet in Object_
	FieldSet
}

func (o *Object_) TypeSystemNode() {}

func (o *Object_) Type() string {
	return "Object"
}

//func (o *Object_) ValueNode()      {}
func (o *Object_) TypeName() NameValue_ {
	return o.Name
}
func (o *Object_) GetSelectionSet() FieldSet {
	return o.FieldSet
}
func (o *Object_) CheckDirectiveRef(dirName NameValue_, err *[]error) {

	o.Directives_.CheckDirectiveRef(dirName, err)

	for _, v := range o.FieldSet {
		v.CheckDirectiveRef(dirName, err)
	}

}

func (o *Object_) CheckDirectiveLocation(err *[]error) {
	o.checkDirectiveLocation_(OBJECT_DL, err)
}

func (f *Object_) CheckImplements(err *[]error) {
	for _, v := range f.Implements {
		var (
			ok   bool
			itf_ GQLTypeProvider
		)
		// check name represents a interface type in ast
		// TODO - requires fetchInterface to use cache - rethink - MAYBE SHOULD NOT BE A OBJECT METHOD but a check in the parser itself as it has access to the cache.
		fmt.Printf("CACHE LOOPKUP: Look for interface %s in cache\n", v.Name)
		if itf_, ok = TyCache[v.Name.String()]; !ok {
			return
		}
		// check object implements the interface
		if itf_.Type() != "Interface" {
			*err = append(*err, fmt.Errorf(`"%s" is not an interface type, %s`, v.Name, v.AtPosition()))
			return
		}
		itf := itf_.(*Interface_)
		satisfied := make(map[NameValue_]bool)
		for _, v := range itf.FieldSet {
			fmt.Println()
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
		}
		//
		// publish in repeatable order because maps cannot
		//
		var s strings.Builder
		for _, ifn := range itf.FieldSet { // interface fields
			if v, ok := satisfied[ifn.Name]; ok {
				if !v {
					s.WriteString(` "`)
					s.WriteString(ifn.Name.String())
					s.WriteString(`"`)
				}
			}
		}
		if len(s.String()) > 0 {
			*err = append(*err, fmt.Errorf(`Type "%s" does not implement interface "%s", missing %s`, f.Name_, itf.Name_, s.String()))
		}
	}
}

func (o *Object_) SolicitAbstractTypes(unresolved UnresolvedMap) {

	o.FieldSet.SolicitAbstractTypes(unresolved)
	// check existence of implement type(s)
	o.Implements.SolicitAbstractTypes(unresolved)
	for _, v := range o.Implements {
		unresolved[v] = nil
	}
	// check existence of each directive type used.
	o.Directives_.SolicitAbstractTypes(unresolved)
}

func (f *Object_) CheckIsOutputType(err *[]error) {
	for _, v := range f.FieldSet {
		if !IsOutputType(v.Type) {
			*err = append(*err, fmt.Errorf(`Field "%s" type "%s", is not an output type %s`, v.Name_, v.Type.Name, v.Type.Name_.AtPosition()))
		}
	}

}

func (f *Object_) CheckIsInputType(err *[]error) {
	//
	for _, v := range f.FieldSet {
		for _, p := range v.ArgumentDefs {
			if !IsInputType(p.Type) {
				*err = append(*err, fmt.Errorf(`Argument "%s" type "%s", is not an input type %s`, p.Name_, p.Type.Name, p.Type.Name_.AtPosition()))
			}
		}
	}
}

func (f *Object_) CheckInputValueType(err *[]error) {
	f.Directives_.CheckInputValueType(err)
	f.FieldSet.CheckInputValueType(err)
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

func (fs *FieldSet) CheckInputValueType(err *[]error) {
	for _, v := range *fs {
		v.CheckInputValueType(err)
	}
}

func (fs *FieldSet) SolicitAbstractTypes(unresolved UnresolvedMap) {
	for _, v := range *fs {
		v.SolicitAbstractTypes(unresolved)
	}
}

func (fs *FieldSet) AppendField(f_ *Field_) error {
	for _, v := range *fs {
		// check field (Name and Type) not already present
		//if v.Equals(f_) { // TODO - where is it necessary to compare Name & Type
		if v.Name.String() == f_.Name.String() {
			loc := f_.Name_.Loc
			if loc != nil {
				return fmt.Errorf(`Duplicate Field name "%s" at line: %d, column: %d`, f_.Name_, loc.Line, loc.Column)
			} else {
				return fmt.Errorf(`Duplicate Field name "%s" at line: %d, column: %d`, f_.Name_)
			}
		}
	}
	*fs = append(*fs, f_)
	return nil
}

// ===============================================================
type AssignTyper interface {
	AssignType(t *GQLtype)
}

// ==================== Field_ ================================
// instance of Object Field
// FieldDefinition
//		 Description-opt	Name	ArgumentsDefinition-opt	:	Type	Directives-opt
type Field_ struct {
	Desc string
	Name_
	ArgumentDefs InputValueDefs //[]*InputValueDef []*ObjectVal
	// :
	Type *GQLtype
	Directives_
}

//TODO  - check argumentsDefs

func (f *Field_) AssignType(t *GQLtype) {
	f.Type = t
}

func (f *Field_) CheckDirectiveLocation(err *[]error) {
	f.checkDirectiveLocation_(FIELD_DEFINITION_DL, err)
	f.ArgumentDefs.CheckDirectiveLocation(err)
}

// func (a *Field_) Equals(b *Field_) bool {
// 	return a.Name_.Equals(b.Name_) && a.Type.Equals(b.Type)
// }

func (f *Field_) SolicitAbstractTypes(unresolved UnresolvedMap) {
	if f.Type == nil {
		log.Panic(fmt.Errorf("Severe Error - not expected: Field.Type is not assigned for [%s]", f.Name_.String()))
	}
	if !f.Type.IsScalar() && f.Type.AST == nil {
		unresolved[f.Type.Name_] = f.Type
	}
	//
	f.ArgumentDefs.SolicitAbstractTypes(unresolved)
	f.Directives_.SolicitAbstractTypes(unresolved)
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
	s.WriteString(" ")
	s.WriteString(f.Directives_.String())
	return s.String()
}

func (f *Field_) AppendField(f_ *InputValueDef, unresolved *[]error) {
	f.ArgumentDefs.AppendField(f_, unresolved)
}

func (f *Field_) CheckDirectiveRef(dir NameValue_, err *[]error) {

	refCheck := func(dirName NameValue_, x GQLTypeProvider) {
		x.CheckDirectiveRef(dirName, err)
	}

	f.ArgumentDefs.CheckDirectiveRef(dir, err)

	f.Directives_.CheckDirectiveRef(dir, err)

	if !f.Type.IsScalar() && f.Type.AST != nil {
		refCheck(dir, f.Type.AST)
	}

}

func (f *Field_) CheckInputValueType(err *[]error) {
	f.ArgumentDefs.CheckInputValueType(err)
	f.Directives_.CheckInputValueType(err)
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

func (fa *InputValueDefs) CheckDirectiveLocation(err *[]error) {
	for _, v := range *fa {
		v.CheckDirectiveLocation(err)
	}
}

func (fa *InputValueDefs) String(encl [2]token.TokenType) string {
	var s strings.Builder
	for i, v := range *fa {
		s.WriteString(" ")
		if i == 0 {
			s.WriteString(" ")
			s.WriteString(string(encl[0]))
		}
		s.WriteString(" ")
		s.WriteString(v.String())
		if i != len(*fa)-1 {
			s.WriteString(" ")
		}
		if i == len(*fa)-1 {
			s.WriteString(" ")
			s.WriteString(string(encl[1]))
		}
	}
	return s.String()
}

func (fa InputValueDefs) SolicitAbstractTypes(unresolved UnresolvedMap) {

	for _, v := range fa {
		v.SolicitAbstractTypes(unresolved)
	}
}

func (fa InputValueDefs) CheckDirectiveRef(dirName NameValue_, err *[]error) {
	for _, v := range fa {
		v.CheckDirectiveRef(dirName, err)
	}
}

func (fa InputValueDefs) CheckIsInputType(err *[]error) {
	for _, p := range fa {
		if !IsInputType(p.Type) {
			*err = append(*err, fmt.Errorf(`Field "%s" of input type "%s", must be an input type %s`, p.Name_, p.Type.Name, p.Type.Name_.AtPosition()))
		}
		//	_ := p.DefaultVal.isType() // e.g. scalar, int | List
	}
}

func (fa InputValueDefs) CheckInputValueType(err *[]error) {
	for _, a := range fa { // go thru each of the argument field objects [] {} scalar
		a.CheckInputValueType(err)
	}
}

// ==================== . InputValueDef . ============================
// ArgumentsDefinition
//		(InputValueDefinitionlist)
// InputValueDefinition
//		Description-opt  Name : Type  =  DefaultValue-opt   Directives-opt
type InputValueDef struct {
	Desc string
	Name_
	Type       *GQLtype
	DefaultVal *InputValue_
	Directives_
}

func (fa *InputValueDef) SolicitAbstractTypes(unresolved UnresolvedMap) { //TODO - check this..should it use unresolvedMap?
	if fa.Type == nil {
		err := fmt.Errorf("Severe Error - not expected: InputValueDef.Type is not assigned for [%s]", fa.Name_.String())
		log.Panic(err)
	}
	if !fa.Type.IsScalar() && fa.Type.AST == nil {
		unresolved[fa.Type.Name_] = fa.Type
	}

	fa.Directives_.SolicitAbstractTypes(unresolved)
}

func (fa *InputValueDef) CheckDirectiveLocation(err *[]error) {
	fa.checkDirectiveLocation_(ARGUMENT_DEFINITION_DL, err)
}

func (fa *InputValueDef) CheckDirectiveRef(dir NameValue_, err *[]error) {

	refCheck := func(dirName NameValue_, x GQLTypeProvider) {
		fmt.Println("inputValueDef refcheck for ", dirName, x.TypeName())
		x.CheckDirectiveRef(dirName, err)
	}

	fa.Directives_.CheckDirectiveRef(dir, err)

	if !fa.Type.IsScalar() && fa.Type.AST != nil {
		refCheck(dir, fa.Type.AST)
	}

}

func (fa *InputValueDef) AssignName(input string, loc *Loc_, unresolved *[]error) {
	fa.Name_.AssignName(input, loc, unresolved)
}

func (fa *InputValueDef) AssignType(t *GQLtype) {
	fa.Type = t
}

func (fa *InputValueDef) String() string {
	var s strings.Builder
	s.WriteString(" ")
	s.WriteString(fa.Name_.String())
	s.WriteString(" : " + fa.Type.String() + " ")
	if fa.DefaultVal != nil {
		s.WriteString(" = ")
		s.WriteString(fa.DefaultVal.String())
	}
	s.WriteString(" ")
	s.WriteString(fa.Directives_.String())
	return s.String()
}

func (a *InputValueDef) CheckInputValueType(err *[]error) {
	//
	a.DefaultVal.CheckInputValueType(a.Type, a.Name_, err)
	a.Directives_.CheckInputValueType(err)
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

func (e *Enum_) TypeSystemNode() {}
func (e *Enum_) SolicitAbstractTypes(unresolved UnresolvedMap) {
	e.Directives_.SolicitAbstractTypes(unresolved)
	for _, v := range e.Values {
		v.SolicitAbstractTypes(unresolved)
	}
}

func (e *Enum_) Type() string {
	return "Enum"
}

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

func (e *Enum_) CheckDirectiveLocation(err *[]error) {
	e.checkDirectiveLocation_(ENUM_DL, err)
	for _, v := range e.Values {
		v.CheckDirectiveLocation(err)
	}
}

func (e *Enum_) CheckInputValueType(err *[]error) {
	for _, v := range e.Values {
		v.CheckInputValueType(err)
	}
}

// ======================  EnumValue =========================

//	EnumValueDefinition
//		Description-opt EnumValue Directives-const-opt
type EnumValue_ struct {
	Desc string
	Name_
	Directives_
	hostValue InputValueProvider
	//
	//parent *Enum_ // is a member to
}

func (e *EnumValue_) ValueNode() {} // instane of InputValue_
func (e *EnumValue_) IsType() TypeFlag_ {
	return ENUMVALUE
}
func (e *EnumValue_) TypeSystemNode() {}
func (e *EnumValue_) SolicitAbstractTypes(unresolved UnresolvedMap) {
	for _, v := range e.Directives {
		unresolved[v.Name_] = nil
	}
}

func (e *EnumValue_) Type() string {
	return "EnumValue"
}
func (e *EnumValue_) CheckDirectiveLocation(err *[]error) {
	e.checkDirectiveLocation_(ENUM_VALUE_DL, err)
}

func (e *EnumValue_) AssignName(s string, l *Loc_, unresolved *[]error) {
	e.Name_.AssignName(s, l, unresolved)
}
func (e *EnumValue_) TypeName() NameValue_ {
	return e.Name
}
func (e *EnumValue_) String() string {
	var s strings.Builder
	s.WriteString(e.Name_.String())
	if e.Directives != nil {
		s.WriteString(" " + e.Directives_.String())
	}
	return s.String()
}

// CheckEnumValue checks the ENUM value (as Argument in Field object) is a member of the ENUM Type.
func (e *EnumValue_) CheckEnumValue(a *GQLtype, err *[]error) {
	// get Enum type and compare it against the instance value
	// TODO - rethink this solution - should not use CacheFEtch in type mthod
	if ast_, ok := TyCache[a.Name.String()]; ok {
		switch enum_ := ast_.(type) {
		case *Enum_:
			found := false
			for _, v := range enum_.Values {
				//	if v.Name_.String() == e.Name_.String() {
				if v.Name_.Equals(e.Name_) {
					found = true
					break
				}
			}
			if !found {
				*err = append(*err, fmt.Errorf(` "%s" is not a member of Enum type %s %s`, e.Name_, a.Name, e.Name_.AtPosition()))
			}
		default:
			*err = append(*err, fmt.Errorf(`Type "%s" is not an ENUM but argument value "%s" is an ENUM value `, a.Name, e.Name_, e.Name_.AtPosition()))
		}

	} else {
		*err = append(*err, fmt.Errorf(`Enum type "%s" is not found in cache %s`, a.Name, e.Name_.AtPosition()))
	}
}

func (e *EnumValue_) CheckInputValueType(err *[]error) {
	e.Directives_.CheckInputValueType(err)
}

// ======================  Schema =========================

type Schema struct {
	rootQuery        *GQLtype
	rootMutation     *GQLtype
	rootSubscription *GQLtype
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
func (i *Interface_) SolicitAbstractTypes(unresolved UnresolvedMap) {
	i.Directives_.SolicitAbstractTypes(unresolved)
	i.FieldSet.SolicitAbstractTypes(unresolved)
}

func (i *Interface_) Type() string {
	return "Interface"
}
func (i *Interface_) TypeName() NameValue_ {
	return i.Name
}
func (i *Interface_) GetSelectionSet() FieldSet {
	return i.FieldSet
}

func (i *Interface_) CheckDirectiveLocation(err *[]error) {
	i.checkDirectiveLocation_(INTERFACE_DL, err)
}

//func (i *Interface_) AssignUnresolvedTypes(ast TypeRepo) error {}
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

func (i *Interface_) Conform(obj GQLTypeProvider) bool {
	obj_, ok := obj.(*Object_)
	if !ok {
		return false
	}
	if len(i.FieldSet) > len(obj_.FieldSet) {
		return false
	}
	for _, fld := range i.FieldSet {
		var found bool
		for _, objFld := range obj_.FieldSet {
			if fld.Name.Equals(objFld.Name) {
				found = true
				break
			}
			if !found {
				return false
			}

		}
	}
	return true
}
func (i *Interface_) CheckFieldMembers(err *[]error) {
	// Fields on a GraphQL interface have the same rules as fields on a GraphQL object;
	// their type can be Scalar, Object, Enum, Interface, or Union, or any wrapping type whose base type is one of those five.
	for _, v := range i.FieldSet {
		switch v.Type.isType() {
		case OBJECT, ENUM, INTERFACE, UNION, FLOAT, INT, BOOLEAN, ID, SCALAR, STRING:
		default:
			*err = append(*err, fmt.Errorf(`Member %q of interface %q is not an appropriate type. Must be an object, enum, interface or union, %s`, v.Name, i.TypeName(), v.Name_.AtPosition()))
		}
	}
}

func (i *Interface_) CheckInputValueType(err *[]error) {
	i.Directives_.CheckInputValueType(err)
	i.FieldSet.CheckInputValueType(err)
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
func (u *Union_) SolicitAbstractTypes(unresolved UnresolvedMap) { // TODO check this is being executed
	u.Directives_.SolicitAbstractTypes(unresolved)
	for _, v := range u.NameS {
		unresolved[v] = nil
	}
}

func (u *Union_) Type() string {
	return "Union"
}
func (u *Union_) TypeName() NameValue_ {
	return u.Name
}

func (u *Union_) CheckDirectiveLocation(err *[]error) {
	u.checkDirectiveLocation_(UNION_DL, err)
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

// Union inherits this from Directive_
// func (u *Union_) CheckInputValueType(err *[]error) {
// 	u.CheckInputValueType(err)
// }

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

//func (e *Input_) ValueNode()      {}// commented out 19/3/2020
func (e *Input_) SolicitAbstractTypes(unresolved UnresolvedMap) { // TODO check this is being executed
	e.Directives_.SolicitAbstractTypes(unresolved)
	e.InputValueDefs.SolicitAbstractTypes(unresolved)
}
func (e *Input_) IsType() TypeFlag_ {
	return INPUT
}

func (e *Input_) Type() string {
	return "Input"
}

//func (e *Input_) ValueNode() {}

func (i *Input_) TypeName() NameValue_ {
	return i.Name
}

func (i *Input_) CheckDirectiveLocation(err *[]error) {
	i.Directives_.checkDirectiveLocation_(INPUT_OBJECT_DL, err)
	for _, v := range i.InputValueDefs {
		v.checkDirectiveLocation_(INPUT_FIELD_DEFINITION_DL, err)
	}
}

func (i *Input_) CheckDirectiveRef(dir NameValue_, err *[]error) {
	i.Directives_.CheckDirectiveRef(dir, err)
	i.InputValueDefs.CheckDirectiveRef(dir, err)
}

func (u *Input_) String() string {
	var encl [2]token.TokenType = [2]token.TokenType{token.LBRACE, token.RBRACE}
	var s strings.Builder
	s.WriteString("\ninput ")
	s.WriteString(" ")
	s.WriteString(u.Name.String())
	s.WriteString(" ")
	s.WriteString(" " + u.Directives_.String())
	s.WriteString(" ")
	s.WriteString(u.InputValueDefs.String(encl))
	return s.String()
}

func (i *Input_) CheckInputValueType(err *[]error) {
	i.Directives_.CheckInputValueType(err)
	i.InputValueDefs.CheckInputValueType(err)
}

// ======================  ScalarProvider =========================

type ScalarProvider interface {
	GQLTypeProvider
	Coerce(i InputValueProvider) (InputValueProvider, error)
}

// ======================  Scalar_ =========================
// ScalarTypeDefinition:
//		Description-opt	input	Name	DirectivesConst-opt	InputFieldsDefinition-opt
type Scalar_ struct {
	Desc string
	Name string // no need to hold Location as its stored in InputValue, parent of this object
	Loc  *Loc_
	Directives_
	Data string
	//
	// scalar data validator types (may not be necessary)
	//
	TimeV   time.Time // any date-time
	NumberV float64   // any number
	IntV    int64     // any int
}

func (e *Scalar_) TypeSystemNode() {}
func (e *Scalar_) ValueNode()      {}
func (e *Scalar_) IsType() TypeFlag_ {
	return SCALAR
}

func (e *Scalar_) Type() string {
	return "scalar"
}
func (e *Scalar_) SolicitAbstractTypes(unresolved UnresolvedMap) { // TODO check this is being executed
	e.Directives_.SolicitAbstractTypes(unresolved)
}
func (i *Scalar_) TypeName() NameValue_ {
	return NameValue_(i.Name)
}

func (e *Scalar_) CheckDirectiveLocation(err *[]error) {
	e.checkDirectiveLocation_(SCALAR_DL, err)
}

func (e *Scalar_) AssignName(s string, loc *Loc_, errS *[]error) {
	ValidateName(s, errS, loc)
	e.Name = s
	e.Loc = loc
}

func (u *Scalar_) String() string {
	var s strings.Builder
	s.WriteString("\nscalar ")
	s.WriteString(u.Name)
	s.WriteString(" " + u.Directives_.String())
	return s.String()
}

// func (s *Scalar_) Print() string {
// 	switch s.Name {
// 	case "Time":
// 		// print time
// 		// print ?
// 		return fmt.Errorf("%s must be a string to be coerced to Time", s.Name)
// 	default:
// 		return fmt.Errorf("%s is not a Scalar ", s.Name)
// 	}
// }

func (s *Scalar_) Coerce(input InputValueProvider) (InputValueProvider, error) {

	fmt.Printf(" %#v, s.Name  %s = \n", s, s.Name)
	switch s.Name {

	case "Time":
		// coerce input to scalar Time
		switch in := input.(type) {
		case String_:
			// try to convert string input value to a scalar (time) iv
			const longForm = "Jan 2, 2006 at 3:04pm (MST)"
			if t, err := time.Parse(longForm, in.String()); err != nil {
				const shortForm = "2006-Jan-02"
				if t, err = time.Parse(shortForm, in.String()); err != nil {
					if t, err = time.Parse(time.RFC3339, in.String()); err != nil {
						return nil, fmt.Errorf("Error in parsing of Time value ")
					}
				}
			} else {
				// convert input value from string to time

				b := &Scalar_{Name: "Time", Data: in.String(), TimeV: t}
				fmt.Printf("*** Cource input from String_ b %#v, s.Name  %s = \n", b, b.Name)
				return b, nil
			}
		}
		// if t_, ok := s.(RawString_); ok {
		// 	// convert string to time
		// 	return tc_, nil
		// }
		return nil, fmt.Errorf("%s must be a string to be coerced to Time", s.Name)

	default:
		return nil, fmt.Errorf("%s is not a Scalar ", s.Name)
	}
}

// ======================  Directive_ =========================
// ScalarTypeDefinition:
//		Description-opt	input	Name	DirectivesConst-opt	InputFieldsDefinition-opt

type Directive_ struct {
	Desc         string
	Name_                       // no need to hold Location as its stored in InputValue, parent of this object
	ArgumentDefs InputValueDefs //TODO consider making InputValueDefs an embedded type ie. an anonymous field
	Location     []DirectiveLoc
}

func (d *Directive_) TypeSystemNode() {}

//
func (d *Directive_) Type() string {
	return "Directive"
}

//func (d *Directive_) ValueNode()      {}
func (d *Directive_) SolicitAbstractTypes(unresolved UnresolvedMap) {
	d.ArgumentDefs.SolicitAbstractTypes(unresolved)
}
func (d *Directive_) CheckDirectiveRef(dir NameValue_, err *[]error) {
	for _, v := range d.ArgumentDefs {
		v.CheckDirectiveRef(dir, err)
	}
}
func (d *Directive_) CheckDirectiveLocation(err *[]error) {
	d.ArgumentDefs.CheckDirectiveLocation(err)
}

func (d *Directive_) CoerceDirectiveName() {
	d.Name_.Name = NameValue_("@" + d.Name.String())
}
func (d *Directive_) TypeName() NameValue_ {
	return d.Name
}

func (d *Directive_) AssignName(input string, loc *Loc_, err *[]error) {
	d.Name_.AssignName(input, loc, err)
}

func (d *Directive_) String() string {
	var (
		s    strings.Builder
		encl [2]token.TokenType = [2]token.TokenType{token.LPAREN, token.RPAREN}
	)
	s.WriteString("\ndirective ")
	s.WriteString(d.Name.String())
	s.WriteString(d.ArgumentDefs.String(encl))
	if len(d.Location) > 0 {

		s.WriteString(" on ")
		for _, v := range d.Location {
			s.WriteString("| ")
			if dloc, ok := DirectiveLocationMap[v]; ok {
				s.WriteString(dloc)
			} else {
				s.WriteString(" not-found ")
			}
		}
	}
	return s.String()
}

func (d *Directive_) CheckIsInputType(err *[]error) {
	for _, p := range d.ArgumentDefs {
		if !IsInputType(p.Type) {
			*err = append(*err, fmt.Errorf(`Argument "%s" type "%s", is not an input type %s`, p.Name_, p.Type.Name, p.Type.Name_.AtPosition()))
		}
		//	_ := p.DefaultVal.isType() // e.g. scalar, int | List
	}
}

func (d *Directive_) AppendField(f_ *InputValueDef, err *[]error) {
	d.ArgumentDefs.AppendField(f_, err)
}

func (d *Directive_) CheckInputValueType(err *[]error) {
	d.ArgumentDefs.CheckInputValueType(err)
}
