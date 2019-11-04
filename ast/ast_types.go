package ast

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

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

type TypeFlag_ uint16

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
	case LIST:
		return token.LIST
	case NULL:
		return token.NULL
	}
	return "NoTypeFound "
}

type UnresolvedMap map[Name_]*Type_

const (
	//  input value types
	_ TypeFlag_ = 1 << iota
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

	// error - not available
	NA
)

// Directive Locations
type DirectiveLoc uint32

// Directive Locations
const (
	_ DirectiveLoc = 1 << iota
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

func IsInputType(t *Type_) bool {
	// determine inputType from t.Name
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

func IsOutputType(t *Type_) bool {
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

// ============================ Type_ ======================

type Type_ struct {
	Constraint byte            // each on bit from right represents not-null constraint applied e.g. in nested list type [type]! is 00000010, [type!]! is 00000011, type! 00000001
	AST        GQLTypeProvider // AST instance of type. WHen would this be used??. Used for non-Scalar types. AST in cache(typeName), then in Type_(typeName). If not in Type_, check cache, then DB.
	Depth      int             // depth of nested List e.g. depth 2 is [[type]]. Depth 0 implies non-list type, depth > 0 is a list type
	Name_                      // type name. inherit AssignName(). Use Name_ to access AST via cache lookup. ALternatively, use AST above.
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

func (o ObjectVals) ValidateInputObjectValues(ref *Type_, err *[]error) {
	//
	//  ref{ name:value name:value ... } -- ref is the input object type specifed for the argument and { } is the argument data
	//
	refFields := make(map[NameValue_]*Type_)
	// check if default input fields has fields not in field Type, PET, MEASURE
	if ref.isType() != INPUT { // required: donot remove
		return
	}
	// get reference type AST object.
	refOV := ref.AST.(*Input_)
	// build a map of fields in reference type - which defines the types of each item in {}
	for _, v := range refOV.InputValueDefs { //
		refFields[v.Name] = v.Type
	}
	fmt.Println("****** ", refFields)
	//
	// loop thru name:value pairs using the ref type spec to match against name and its associated type for each pair.
	//
	for _, v := range o { // []*ArgumentT    type ArgumentT struct { Name_, Value *InputValue_}  type InputValue_ struct {Value ValueI,Loc *Loc_}
		//    ValueI populated by parser.parseInputValue_(): ast.Int_, ast.Flaat_, ast.List_, ast.ObjectVals, ast.EnumValue_ etc
		if reftype, ok := refFields[v.Name]; !ok {
			*err = append(*err, fmt.Errorf(`field "%s" does not exist in type %s  %s`, v.Name, ref.TypeName(), v.AtPosition()))

		} else {

			//fmt.Println(v.Value.isType(), reftype.isType())

			if v.Value.isType() != reftype.isType() && v.Value.isType() != LIST { // && v.Value.isType() != LIST {
				if reftype.Depth > 0 && reftype.Constraint == 1 && v.Value.isType() != LIST {
					*err = append(*err, fmt.Errorf(`Argument type "%s", value should be a List type %s %s`, ref.Name, reftype.isType(), v.Value.AtPosition()))
					*err = append(*err, fmt.Errorf(`Argument type "%s", value has type %s should be %s %s`, ref.Name, v.Value.isType(), reftype.isType(), v.Value.AtPosition()))
				} else {
					*err = append(*err, fmt.Errorf(`Argument type "%s", value has type %s should be %s %s`, ref.Name, v.Value.isType(), reftype.isType(), v.Value.AtPosition()))
				}
			} else {
				if reftype.Depth > 0 && reftype.Constraint == 1 && v.Value.isType() != LIST {
					*err = append(*err, fmt.Errorf(`Argument type "%s", value should be a List type %s %s`, ref.Name, reftype.isType(), v.Value.AtPosition()))
				}
			}

			// look at argument value type as it may be a list or another input object type
			switch inobj := v.Value.Value.(type) { // y inob:Float_

			case List_:
				var errSet bool
				if reftype.Depth == 0 {
					*err = append(*err, fmt.Errorf(`Field, %s, is not a LIST type but input data is a LIST type, %s`, v.Name, v.AtPosition()))
					errSet = true
				}
				// maxd records maximum depth of list(d=1) [] list of lists [[]](d=2) = [[][][][]] list of lists of lists (d=3) [[[]]] = [[[][][]],[[][][][]],[[]]]
				d := 0
				maxd := 0
				inobj.ValidateListValues(reftype, &d, &maxd, err)
				d--
				if maxd != reftype.Depth && !errSet {
					*err = append(*err, fmt.Errorf(`Argument "%s", nested List type depth different reqired %d, got %d %s`, v.Name, reftype.Depth, maxd, v.AtPosition()))
				}

			case ObjectVals:
				inobj.ValidateInputObjectValues(reftype, err)

			}
		}
	}
	//
	// check mandatory fields present
	//
	var at string
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

// =================================================================
// Slice of Name_
type NameS []Name_

// CheckUnresolvedTypes is typically promoted to type that embedds the NameS type.
func (f NameS) CheckUnresolvedTypes(unresolved UnresolvedMap) { //TODO rename to checkUnresolvedTypes
	for _, v := range f {
		// check if the implement type is cached.
		if _, ok := CacheFetch(v.Name); !ok {
			unresolved[v] = nil
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
}

func (o *Object_) CheckUnresolvedTypes(unresolved UnresolvedMap) {
	o.FieldSet.CheckUnresolvedTypes(unresolved)
	o.Implements.CheckUnresolvedTypes(unresolved)
	o.Directives_.CheckUnresolvedTypes(unresolved)
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
// 	Value InputValueProvider //  IV:type|value = assert type to determine InputValue_'s type
// 	Loc   *Loc_
// // }
// type Type_ struct {
// 	Constraint byte          // each on bit from right represents not-null constraint applied e.g. in nested list type [type]! is 00000010, [type!]! is 00000011, type! 00000001
// 	AST        GQLTypeProvider // AST instance of type. WHen would this be used??. Used for non-Scalar types. AST in cache(typeName), then in Type_(typeName). If not in Type_, check cache, then DB.
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
		// 	Type       *Type_   	// ** argument type specification   	<==== required type of argument
		// 	DefaultVal *InputValue_ // ** input value(s) type(s)        	<==== instance data to check against required type
		// 	Directives_
		// }
		// type InputValue_ struct {
		// 	Value InputValueProvider
		// 	Loc   *Loc_
		// }
		// a.Type is argument type -  check it against a.DefaultVal.Value.isType()
		for _, a := range v.ArgumentDefs { // go thru each of the argument field objects [] {} scalar

			if a.DefaultVal != nil {

				//a.DefaultVal.CheckInputValueType(a.Type, err)

				// what type is the default value
				switch defval := a.DefaultVal.Value.(type) {

				case List_: // [ "ads", "wer" ]
					if a.Type.Depth == 0 { // required type is not a LIST
						*err = append(*err, fmt.Errorf(`Argument "%s", type is not a list but default value is a list %s`, a.Name_, a.DefaultVal.AtPosition()))
						return
					}
					var d int = 0
					var maxd int
					defval.ValidateListValues(a.Type, &d, &maxd, err) // a.Type is the data type of the list items
					//
					if maxd != a.Type.Depth {
						*err = append(*err, fmt.Errorf(`Argument "%s", nested List type depth different reqired %d, got %d %s`, a.Name_, a.Type.Depth, maxd, a.DefaultVal.AtPosition()))
					}

				case ObjectVals:
					// { x: "ads", y: 234 }
					defval.ValidateInputObjectValues(a.Type, err)

				case *EnumValue_:
					// EAST WEST NORHT SOUTH
					if a.Type.isType() != ENUM {
						*err = append(*err, fmt.Errorf(`"%s" is an enum like value but the argument type "%s" is not an Enum type %s`, defval.Name, a.Type.Name_, a.DefaultVal.AtPosition()))
					} else {
						defval.CheckEnumValue(a.Type, err)
					}

				default:
					// single instance data
					fmt.Printf("name: %s\n", a.Type.Name_)
					fmt.Printf("constrint: %08b\n", a.Type.Constraint)
					fmt.Printf("depth: %d\n", a.Type.Depth)
					fmt.Println("defType ", a.DefaultVal.isType(), a.DefaultVal.IsScalar())
					fmt.Println("refType ", a.Type.isType())

					// save default type before potential coercing
					defType := a.DefaultVal.isType()

					if a.DefaultVal.isType() == NULL {
						// test case FieldArgListInt3_6 [int]!  null  - value cannot be null
						if a.Type.Constraint>>uint(a.Type.Depth)&1 == 1 {
							*err = append(*err, fmt.Errorf(`Value cannot be NULL %s`, a.DefaultVal.AtPosition()))
						}

					} else if a.Type.isType() == SCALAR { //a.DefaultVal.IsScalar() {
						// can the input value be coerced e.g. from string to Time
						// try coercing default value to the appropriate scalar e.g. string to Time
						if s, ok := a.Type.AST.(ScalarProvider); ok { // assert interface supported - normal assert type (*Scalar_) would also work just as well because there is only 1 scalar type really
							if civ, cerr := s.Coerce(a.DefaultVal.Value); cerr != nil {
								*err = append(*err, cerr)
								return
							} else {
								a.DefaultVal.Value = civ
								defType = a.DefaultVal.isType()
							}
						}
						// coerce to a list of appropriate depth. Current value is not a list as this is switch case default - see other cases.
						if a.Type.Depth > 0 {
							var coerce2list func(i *InputValue_, depth int) *InputValue_
							// type List_ []*InputValue_

							coerce2list = func(i *InputValue_, depth int) *InputValue_ {
								if depth == 0 {
									return i
								}
								vallist := make(List_, 1, 1)
								vallist[0] = i
								vi := &InputValue_{Value: vallist, Loc: i.Loc}
								depth--
								return coerce2list(vi, depth)
							}
							a.DefaultVal = coerce2list(a.DefaultVal, a.Type.Depth)
						}

					} else {
						// coerce to a list of appropriate depth. Current value is not a list as this is case default - see other cases.
						if a.Type.Depth > 0 {
							var coerce2list func(i *InputValue_, depth int) *InputValue_
							// type List_ []*InputValue_

							coerce2list = func(i *InputValue_, depth int) *InputValue_ {
								if depth == 0 {
									return i
								}
								vallist := make(List_, 1, 1)
								vallist[0] = i
								vi := &InputValue_{Value: vallist, Loc: i.Loc}
								depth--
								return coerce2list(vi, depth)
							}
							a.DefaultVal = coerce2list(a.DefaultVal, a.Type.Depth)
						}
					}

					if defType != NULL && defType != a.Type.isType() {
						*err = append(*err, fmt.Errorf(`Required type "%s", got "%s" %s`, a.Type.isType(), defType, a.DefaultVal.AtPosition()))
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

func (fs *FieldSet) CheckUnresolvedTypes(unresolved UnresolvedMap) {
	for _, v := range *fs {
		v.CheckUnresolvedTypes(unresolved)
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

func (f *Field_) CheckDirectiveLocation(err *[]error) {
	f.checkDirectiveLocation_(FIELD_DEFINITION_DL, err)
	f.ArgumentDefs.CheckDirectiveLocation(err)
}

// func (a *Field_) Equals(b *Field_) bool {
// 	return a.Name_.Equals(b.Name_) && a.Type.Equals(b.Type)
// }

func (f *Field_) CheckUnresolvedTypes(unresolved UnresolvedMap) {
	if f.Type == nil {
		log.Panic(fmt.Errorf("Severe Error - not expected: Field.Type is not assigned for [%s]", f.Name_.String()))
	}
	if !f.Type.IsScalar() {
		if f.Type.AST == nil {
			// check in cache only at this stage.
			// When control passes back to parser we resolved the unresolved using the DB and parse stmt if found.
			if ast, ok := CacheFetch(f.Type.Name); !ok {
				unresolved[f.Type.Name_] = f.Type
			} else {
				f.Type.AST = ast
			}
		}
	}
	//
	f.ArgumentDefs.CheckUnresolvedTypes(unresolved)
	f.Directives_.CheckUnresolvedTypes(unresolved)
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

func (fa InputValueDefs) CheckUnresolvedTypes(unresolved UnresolvedMap) {

	for _, v := range fa {
		v.CheckUnresolvedTypes(unresolved)
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
			*err = append(*err, fmt.Errorf(`Argument "%s" type "%s", is not an input type %s`, p.Name_, p.Type.Name, p.Type.Name_.AtPosition()))
		}
		//	_ := p.DefaultVal.isType() // e.g. scalar, int | List
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

func (fa *InputValueDef) CheckUnresolvedTypes(unresolved UnresolvedMap) { //TODO - check this..should it use unresolvedMap?
	if fa.Type == nil {
		err := fmt.Errorf("Severe Error - not expected: InputValueDef.Type is not assigned for [%s]", fa.Name_.String())
		log.Panic(err)
	}
	if !fa.Type.IsScalar() && fa.Type.AST == nil {
		unresolved[fa.Type.Name_] = fa.Type
	}
	fa.Directives_.CheckUnresolvedTypes(unresolved)
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

func (e *Enum_) TypeSystemNode() {}
func (e *Enum_) CheckUnresolvedTypes(unresolved UnresolvedMap) {
	e.Directives_.CheckUnresolvedTypes(unresolved)
	for _, v := range e.Values {
		v.CheckUnresolvedTypes(unresolved)
	}
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

// ======================  EnumValue =========================

//	EnumValueDefinition
//		Description-opt EnumValue Directives-const-opt
type EnumValue_ struct {
	Desc string
	Name_
	Directives_
}

func (e *EnumValue_) ValueNode()      {} // instane of InputValue_
func (e *EnumValue_) TypeSystemNode() {}
func (e *EnumValue_) CheckUnresolvedTypes(unresolved UnresolvedMap) {
	e.Directives_.CheckUnresolvedTypes(unresolved)
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

// CheckUnresolvedTypes checks the ENUM value (as Argument in Field object) is a member of the ENUM Type.
func (e *EnumValue_) CheckEnumValue(a *Type_, err *[]error) {
	// get Enum type and compare it against the instance value
	if ast_, ok := CacheFetch(a.Name); ok {
		switch enum_ := ast_.(type) {
		case *Enum_:
			found := false
			for _, v := range enum_.Values {
				if v.Name_.String() == e.Name_.String() {
					found = true
					break
				}
			}
			if !found {
				*err = append(*err, fmt.Errorf(` "%s" is not a member of type Enum %s %s`, e.Name_, a.Name, e.Name_.AtPosition()))
			}
		default:
			*err = append(*err, fmt.Errorf(`Type "%s" is not an ENUM but argument value "%s" is an ENUM value `, a.Name, e.Name_, e.Name_.AtPosition()))
		}

	} else {
		*err = append(*err, fmt.Errorf(`Enum type "%s" is not found in cache`, a.Name, e.Name_.AtPosition()))
	}
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
func (i *Interface_) CheckUnresolvedTypes(unresolved UnresolvedMap) {
	i.Directives_.CheckUnresolvedTypes(unresolved)
	i.FieldSet.CheckUnresolvedTypes(unresolved)
}
func (i *Interface_) TypeName() NameValue_ {
	return i.Name
}

func (i *Interface_) CheckDirectiveLocation(err *[]error) {
	i.checkDirectiveLocation_(INTERFACE_DL, err)
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
func (u *Union_) CheckUnresolvedTypes(unresolved UnresolvedMap) { // TODO check this is being executed
	u.Directives_.CheckUnresolvedTypes(unresolved)
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
func (e *Input_) CheckUnresolvedTypes(unresolved UnresolvedMap) { // TODO check this is being executed
	e.Directives_.CheckUnresolvedTypes(unresolved)
	e.InputValueDefs.CheckUnresolvedTypes(unresolved)
}
func (e *Input_) ValueNode() {}

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
	s.WriteString(u.Name.String())
	s.WriteString(" " + u.Directives_.String())
	s.WriteString(u.InputValueDefs.String(encl))
	return s.String()
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
func (e *Scalar_) CheckUnresolvedTypes(unresolved UnresolvedMap) { // TODO check this is being executed
	e.Directives_.CheckUnresolvedTypes(unresolved)
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
	Name_        // no need to hold Location as its stored in InputValue, parent of this object
	ArgumentDefs InputValueDefs
	Location     []DirectiveLoc
}

func (d *Directive_) TypeSystemNode() {}
func (d *Directive_) ValueNode()      {}
func (d *Directive_) CheckUnresolvedTypes(unresolved UnresolvedMap) {
	d.ArgumentDefs.CheckUnresolvedTypes(unresolved)
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

func (d *Directive_) CheckInputValueType(err *[]error) { // TODO try merging wih *Object_ version
	for _, a := range d.ArgumentDefs { // go thru each of the argument field objects [] {} scalar

		if a.DefaultVal != nil {

			// what type is the default value
			switch defval := a.DefaultVal.Value.(type) {

			case List_: // [ "ads", "wer" ]
				if a.Type.Depth == 0 { // required type is not a LIST
					*err = append(*err, fmt.Errorf(`Argument "%s", type is not a list but default value is a list %s`, a.Name_, a.DefaultVal.AtPosition()))
					return
				}
				var d int = 0
				var maxd int
				defval.ValidateListValues(a.Type, &d, &maxd, err) // a.Type is the data type of the list items
				//
				if maxd != a.Type.Depth {
					*err = append(*err, fmt.Errorf(`Argument "%s", nested List type depth different reqired %d, got %d %s`, a.Name_, a.Type.Depth, maxd, a.DefaultVal.AtPosition()))
				}

			case ObjectVals:
				// { x: "ads", y: 234 }
				defval.ValidateInputObjectValues(a.Type, err)

			case *EnumValue_:
				// EAST WEST NORHT SOUTH
				if a.Type.isType() != ENUM {
					*err = append(*err, fmt.Errorf(`"%s" is an enum like value but the argument type "%s" is not an Enum type %s`, defval.Name, a.Type.Name_, a.DefaultVal.AtPosition()))
				} else {
					defval.CheckEnumValue(a.Type, err)
				}

			default:
				// single instance data
				fmt.Printf("name: %s\n", a.Type.Name_)
				fmt.Printf("constrint: %08b\n", a.Type.Constraint)
				fmt.Printf("depth: %d\n", a.Type.Depth)
				fmt.Println("defType ", a.DefaultVal.isType(), a.DefaultVal.IsScalar())
				fmt.Println("refType ", a.Type.isType())

				// save default type before potential coercing
				defType := a.DefaultVal.isType()

				if a.DefaultVal.isType() == NULL {
					// test case FieldArgListInt3_6 [int]!  null  - value cannot be null
					if a.Type.Constraint>>uint(a.Type.Depth)&1 == 1 {
						*err = append(*err, fmt.Errorf(`Value cannot be NULL %s`, a.DefaultVal.AtPosition()))
					}

				} else if a.Type.isType() == SCALAR { //a.DefaultVal.IsScalar() {
					// can the input value be coerced e.g. from string to Time
					// try coercing default value to the appropriate scalar e.g. string to Time
					if s, ok := a.Type.AST.(ScalarProvider); ok { // assert interface supported - normal assert type (*Scalar_) would also work just as well because there is only 1 scalar type really
						if civ, cerr := s.Coerce(a.DefaultVal.Value); cerr != nil {
							*err = append(*err, cerr)
							return
						} else {
							a.DefaultVal.Value = civ
							defType = a.DefaultVal.isType()
						}
					}
					// coerce to a list of appropriate depth. Current value is not a list as this is switch case default - see other cases.
					if a.Type.Depth > 0 {
						var coerce2list func(i *InputValue_, depth int) *InputValue_
						// type List_ []*InputValue_

						coerce2list = func(i *InputValue_, depth int) *InputValue_ {
							if depth == 0 {
								return i
							}
							vallist := make(List_, 1, 1)
							vallist[0] = i
							vi := &InputValue_{Value: vallist, Loc: i.Loc}
							depth--
							return coerce2list(vi, depth)
						}
						a.DefaultVal = coerce2list(a.DefaultVal, a.Type.Depth)
					}

				} else {
					// coerce to a list of appropriate depth. Current value is not a list as this is case default - see other cases.
					if a.Type.Depth > 0 {
						var coerce2list func(i *InputValue_, depth int) *InputValue_
						// type List_ []*InputValue_

						coerce2list = func(i *InputValue_, depth int) *InputValue_ {
							if depth == 0 {
								return i
							}
							vallist := make(List_, 1, 1)
							vallist[0] = i
							vi := &InputValue_{Value: vallist, Loc: i.Loc}
							depth--
							return coerce2list(vi, depth)
						}
						a.DefaultVal = coerce2list(a.DefaultVal, a.Type.Depth)
					}
				}

				if defType != NULL && defType != a.Type.isType() {
					*err = append(*err, fmt.Errorf(`Required type "%s", got "%s" %s`, a.Type.isType(), defType, a.DefaultVal.AtPosition()))
				}
			}
		}
	}
}

func (a *InputValue_) CheckInputValueType__(m *Type_, nm Name_, err *[]error) {

	if a == nil {
		return
	}
	// what type is the default value
	switch defval := a.Value.(type) {

	case List_: // [ "ads", "wer" ]
		if m.Depth == 0 { // required type is not a LIST
			*err = append(*err, fmt.Errorf(`Argument "%s", type is not a list but default value is a list %s`, nm, a.AtPosition()))
			return
		}
		var d int = 0
		var maxd int
		defval.ValidateListValues(m, &d, &maxd, err) // m.Type is the data type of the list items
		//
		if maxd != m.Depth {
			*err = append(*err, fmt.Errorf(`Argument "%s", nested List type depth different reqired %d, got %d %s`, nm, m.Depth, maxd, a.AtPosition()))
		}

	case ObjectVals:
		// { x: "ads", y: 234 }
		defval.ValidateInputObjectValues(m, err)

	case *EnumValue_:
		// EAST WEST NORHT SOUTH
		if m.isType() != ENUM {
			*err = append(*err, fmt.Errorf(`"%s" is an enum like value but the argument type "%s" is not an Enum type %s`, defval.Name, m.Name_, a.AtPosition()))
		} else {
			defval.CheckEnumValue(m, err)
		}

	default:
		// single instance data
		fmt.Printf("name: %s\n", m.Name_)
		fmt.Printf("constrint: %08b\n", m.Constraint)
		fmt.Printf("depth: %d\n", m.Depth)
		fmt.Println("defType ", a.isType(), a.IsScalar())
		fmt.Println("refType ", m.isType())

		// save default type before potential coercing
		defType := a.isType()

		if a.isType() == NULL {
			// test case FieldArgListInt3_6 [int]!  null  - value cannot be null
			if m.Constraint>>uint(m.Depth)&1 == 1 {
				*err = append(*err, fmt.Errorf(`Value cannot be NULL %s`, a.AtPosition()))
			}

		} else if m.isType() == SCALAR { //a.IsScalar() {
			// can the input value be coerced e.g. from string to Time
			// try coercing default value to the appropriate scalar e.g. string to Time
			if s, ok := m.AST.(ScalarProvider); ok { // assert interface supported - normal assert type (*Scalar_) would also work just as well because there is only 1 scalar type really
				if civ, cerr := s.Coerce(a.Value); cerr != nil {
					*err = append(*err, cerr)
					return
				} else {
					a.Value = civ
					defType = a.isType()
				}
			}
			// coerce to a list of appropriate depth. Current value is not a list as this is switch case default - see other cases.
			if m.Depth > 0 {
				var coerce2list func(i *InputValue_, depth int) *InputValue_
				// type List_ []*InputValue_

				coerce2list = func(i *InputValue_, depth int) *InputValue_ {
					if depth == 0 {
						return i
					}
					vallist := make(List_, 1, 1)
					vallist[0] = i
					vi := &InputValue_{Value: vallist, Loc: i.Loc}
					depth--
					return coerce2list(vi, depth)
				}
				a = coerce2list(a, m.Depth)
			}

		} else {
			// coerce to a list of appropriate depth. Current value is not a list as this is case default - see other cases.
			if m.Depth > 0 {
				var coerce2list func(i *InputValue_, depth int) *InputValue_
				// type List_ []*InputValue_

				coerce2list = func(i *InputValue_, depth int) *InputValue_ {
					if depth == 0 {
						return i
					}
					vallist := make(List_, 1, 1)
					vallist[0] = i
					vi := &InputValue_{Value: vallist, Loc: i.Loc}
					depth--
					return coerce2list(vi, depth)
				}
				a = coerce2list(a, m.Depth)
			}
		}

		if defType != NULL && defType != m.isType() {
			*err = append(*err, fmt.Errorf(`Required type "%s", got "%s" %s`, m.isType(), defType, a.AtPosition()))
		}
	}
}
