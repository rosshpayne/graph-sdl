package ast

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/graph-sdl/token"
)

// =================  InputValueProvider =================================

//  InputValueProvider represents the Graph QL Input Value types (see parseInputValue:) &Int_, &Float_,...,&Enum_, &List_, &ObjectVals
type InputValueProvider interface {
	ValueNode()
	IsType() TypeFlag_
	String() string
	//	Exists() bool
}

// =================  InputValue =================================

// input values used for "default values" in arguments in type and field arguments and input objecs.
type InputValue_ struct {
	InputValueProvider // Important: this is an Interface (embedded value|type), so the type of the input value is defined in the interface value.
	Loc                *Loc_
}

//func (iv *InputValue_) InputValueNode() {}

func (iv *InputValue_) String() string {

	switch x := iv.InputValueProvider.(type) {
	case RawString_:
		return token.RAWSTRINGDEL + iv.InputValueProvider.String() + token.RAWSTRINGDEL //+ "-" + iv.dTString()
	case String_:
		return token.STRINGDEL + iv.InputValueProvider.String() + token.STRINGDEL //+ "-" + iv.dTString() + iv.Loc.String()
	case *Scalar_:
		switch x.Name {
		case "Time":
			return fmt.Sprintf("%q", x.TimeV.String())
		}
	}
	if iv.InputValueProvider == nil { // interface is not populated with concrete value
		return ""
	}
	return iv.InputValueProvider.String() //+ "-" + iv.dTString()
}

func (iv *InputValue_) AtPosition() string {
	return iv.Loc.String()
}

func (iv *InputValue_) isType() TypeFlag_ {
	// Union are not a valid input value
	switch iv.InputValueProvider.(type) {
	case ID_:
		return ID
	case Int_:
		return INT
	case Float_:
		return FLOAT
	case Bool_:
		return BOOLEAN
	case String_:
		return STRING
	case RawString_:
		return STRING
	case *Scalar_:
		return SCALAR
	case *EnumValue_:
		return ENUM
		//	case *Union_: // Union is not a valid input value
	case ObjectVals:
		return OBJECT
	// case *Input_:
	// 	return INPUT
	case Null_:
		return NULL
	case List_:
		return LIST
	}
	return NA
}

func (iv *InputValue_) IsScalar() bool {
	// Union are not a valid input value
	switch iv.InputValueProvider.(type) {
	case Int_, Bool_, Float_, RawString_:
		return true
	case *Scalar_:
		return true
	}
	return false
}

// CheckInputValueType__ called from graphql package to validate
// * default values of variables
// * arguments in fields and directives .
// * output from field resolvers
// refType is the type definition the value of the InputValue_ must match
// nm is the name of the associated argument or input - used for its Loc value
// err contains all errors caught during validation
func (a *InputValue_) CheckInputValueType(refType *Type_, nm Name_, err *[]error) {

	fmt.Println("=========== CheckInputValueType ==============")
	if a == nil {
		return
	}
	// what type is the default value
	var atPosition string
	if nm.Loc == nil {
		atPosition = ""
	} else {
		atPosition = a.AtPosition()
	}
	switch valueType := a.InputValueProvider.(type) {

	case List_:
		// [ "ads", "wer" ]
		// single instance data
		fmt.Printf("name: %s\n", refType.Name_)
		fmt.Printf("constrint: %08b\n", refType.Constraint)
		fmt.Printf("depth: %d\n", refType.Depth)
		fmt.Println("defType ", a.isType(), a.IsScalar())
		fmt.Println("refType ", refType.isType())
		fmt.Println("=========== CheckInputValueType  List_ ==============")
		if refType.Depth == 0 { // required type is not a LIST
			*err = append(*err, fmt.Errorf(`Input value  "%s" named "%s"  is not a list but required type is a list %s`, valueType.String(), nm, atPosition))
			return
		}
		var d int = 0
		var maxd int
		valueType.ValidateListValues(refType, &d, &maxd, err) // m.Type is the data type of the list items
		//
		if maxd != refType.Depth {
			*err = append(*err, fmt.Errorf(`Input value "%s", nested List type depth different reqired %d, got %d %s`, nm, refType.Depth, maxd, atPosition))
		}

	case ObjectVals:
		//  { name:value name:value ... } - match the name,value pairs against the refType (object type fields or input type fields)
		fmt.Println("=========== CheckInputValueType  ObjVals ==============")

		valueType.ValidateObjectValues(refType, err)

	case *EnumValue_:
		if refType.Depth > 0 { // required type is not a LIST
			*err = append(*err, fmt.Errorf(`List type expected, got an enum value "%s" instead for "%s" %s`, valueType.String(), nm, atPosition))
			return
		}
		// EAST WEST NORHT SOUTH
		fmt.Println("=========== CheckInputValueType  EnumValue ==============")
		if refType.isType() != ENUM {
			*err = append(*err, fmt.Errorf(`"%s" is an enum like value but the argument type "%s" is not an Enum type %s`, valueType.Name, refType.Name_, atPosition))
		} else {
			valueType.CheckEnumValue(refType, err)
		}

	default:
		// single instance data
		fmt.Printf("name: %s\n", refType.Name_)
		fmt.Printf("constrint: %08b\n", refType.Constraint)
		fmt.Printf("depth: %d\n", refType.Depth)
		fmt.Println("defType ", a.isType(), a.IsScalar())
		fmt.Println("refType ", refType.isType())

		// save default type before potential coercing
		defType := a.isType()

		if a.isType() == NULL {
			// test case FieldArgListInt3_6 [int]!  null  - value cannot be null
			if refType.Constraint>>uint(refType.Depth)&1 == 1 {
				*err = append(*err, fmt.Errorf(`Value cannot be NULL %s`, atPosition))
			}

		} else if refType.isType() == SCALAR { //a.IsScalar() {
			// can the input value be coerced e.g. from string to Time
			// try coercing default value to the appropriate scalar e.g. string to Time
			if s, ok := refType.AST.(ScalarProvider); ok { // assert interface supported - normal assert type (*Scalar_) would also work just as well because there is only 1 scalar type really
				if civ, cerr := s.Coerce(a.InputValueProvider); cerr != nil {
					*err = append(*err, cerr)
					return
				} else {
					a.InputValueProvider = civ
					defType = a.isType()
				}
			}
			// coerce to a list of appropriate depth. Current value is not a list as this is switch case default - see other cases.
			if refType.Depth > 0 {
				var coerce2list func(i *InputValue_, depth int) *InputValue_
				// type List_ []*InputValue_

				coerce2list = func(i *InputValue_, depth int) *InputValue_ {
					if depth == 0 {
						return i
					}
					vallist := make(List_, 1, 1)
					vallist[0] = i
					vi := &InputValue_{InputValueProvider: vallist, Loc: i.Loc}
					depth--
					return coerce2list(vi, depth)
				}
				a = coerce2list(a, refType.Depth)
			}

		} else {
			// coerce to a list of appropriate depth. Current value is not a list as this is case default - see other cases.
			if refType.Depth > 0 {
				var coerce2list func(i *InputValue_, depth int) *InputValue_
				// type List_ []*InputValue_

				coerce2list = func(i *InputValue_, depth int) *InputValue_ {
					if depth == 0 {
						return i
					}
					vallist := make(List_, 1, 1)
					vallist[0] = i
					vi := &InputValue_{InputValueProvider: vallist, Loc: i.Loc}
					depth--
					return coerce2list(vi, depth)
				}
				a = coerce2list(a, refType.Depth)
			}
		}

		if defType != NULL && defType != refType.isType() {
			*err = append(*err, fmt.Errorf(`Required type "%s", got "%s" %s`, refType.isType(), defType, atPosition))
		}
	}
}

func BaseType(t GQLTypeProvider) string {
	return IsGLType(t)
}

func IsGLType(t GQLTypeProvider) string {
	//
	//
	// non-standard defined types
	//
	switch t.(type) {
	case *Object_:
		return "O"
	case *Interface_:
		return "I"
	case *Enum_:
		return "E"
	case *Input_:
		return "In"
	case *Union_:
		return "U"
	case *Scalar_:
		return "S"
	}
	//
	return "X"
}

// ============================ Type_ ======================

type Type_ struct {
	Constraint byte            // each on bit from right represents not-null constraint applied e.g. in nested list type [type]! is 00000010, [type!]! is 00000011, type! 00000001
	AST        GQLTypeProvider // AST instance of type. WHen would this be used??. Used for non-Scalar types. AST in cache(typeName), then in Type_(typeName). If not in Type_, check cache, then DB.
	Depth      int             // depth of nested List e.g. depth 2 is [[type]]. Depth 0 implies non-list type, depth > 0 is a list type
	Name_                      // type name. inherit AssignName(). Use Name_ to access AST via cache lookup. ALternatively, use AST above or TypeFlag_ instead of string.
	Base       string          // base type e.g. Name_ = "Episode" has Base = E(num)
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

// dataTypeString - prints the datatype of the type specification

func (t *Type_) isType() TypeFlag_ {
	//
	// Object types have nested types i.e. each field has a *Type attribute
	//  the *Type.AST can itself be another object or a scalar (system or user defined)
	// Scalars do not have a *Type attribute as they represent the leaf node in a tree of types.
	//

	switch t.Name.String() {
	//
	// system scalars
	//
	case token.ID:
		return ID
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
	// case token.Time: // could implement new scalar time at this level or embedded in Scalar_ like user defined scales would be.
	// return TIME
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
			case *EnumValue_:
				return ENUMVALUE
			case *Input_:
				return INPUT
			case *Union_:
				return UNION
			case *Scalar_:
				return SCALAR
			}
			//
			return NA
		}
	}
	return NA
}

func (t *Type_) isType2() TypeFlag_ {
	//
	// Object types have nested types i.e. each field has a *Type attribute
	//  the *Type.AST can itself be another object or a scalar (system or user defined)
	// Scalars do not have a *Type attribute as they represent the leaf node in a tree of types.
	//
	if t.Depth > 0 {
		return LIST
	} else {
		return t.isType()
	}
}

func (t *Type_) IsScalar() bool {
	switch t.isType() {
	case INT, FLOAT, STRING, BOOLEAN, SCALAR, ID, ENUM, ENUMVALUE:
		return true
	default:
		return false
	}

}

func (t *Type_) IsType() TypeFlag_ {
	return t.isType()
}

func (t *Type_) IsType2() TypeFlag_ {
	if t.Depth > 0 {
		return LIST
	}
	return t.isType()
}

// ================= Input Value scalar datatypes ===================

type Null_ bool // moved from Scalar to it's own type. No obvious reason why - no obvious advantage at this stage

func (n Null_) ValueNode() {}
func (n Null_) IsType() TypeFlag_ {
	return NULL
}
func (n Null_) String() string {
	if n == false {
		return ""
	}
	return "null"
}

type Int_ string //int

func (i Int_) ValueNode() {}
func (n Int_) IsType() TypeFlag_ {
	return INT
}
func (i Int_) String() string {
	//return strconv.FormatInt(int64(i), 10)
	return string(i)
}

type ID_ string //float64

func (f ID_) ValueNode() {}
func (f ID_) IsType() TypeFlag_ {
	return ID
}
func (f ID_) String() string {
	return string(f)
	//return strconv.FormatFloat(float64(f), 'G', -1, 64)
}

type Float_ string //float64

func (f Float_) ValueNode() {}
func (f Float_) IsType() TypeFlag_ {
	return FLOAT
}
func (f Float_) String() string {
	return string(f)
	//return strconv.FormatFloat(float64(f), 'G', -1, 64)
}

type String_ string

func (s String_) ValueNode() {}
func (s String_) IsType() TypeFlag_ {
	return STRING
}
func (s String_) String() string {
	return string(s)
}

type RawString_ string

func (s RawString_) IsType() TypeFlag_ {
	return STRING
}
func (s RawString_) ValueNode() {}

func (s RawString_) String() string {
	return string(s)
}

type Bool_ bool //bool

func (b Bool_) ValueNode() {}
func (b Bool_) IsType() TypeFlag_ {
	return BOOLEAN
}
func (b Bool_) String() string {
	//return strconv.FormatBool(bool(i))
	if b {
		return "true"
	}
	return "false"
}

// List for input values only - just a bunch of types (any) can be the same or different. The base type is defined elsewhere in the TYPE field of a argument for example.

type List_ []*InputValue_

func (l List_) ValueNode() {}
func (l List_) IsType() TypeFlag_ {
	return LIST
}
func (l List_) TypeName() string {
	return "List"
}
func (l List_) String() string {
	var s strings.Builder
	s.WriteString("[")
	for i, v := range l {
		//fmt.Printf("string() len(l)  %d %T  %T %d \n", len(l), v, v.InputValueProvider, i)
		s.WriteString(v.String())
		if i < len(l)-1 {
			s.WriteString(" ")
		}
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
// 	Value InputValueProvider //  IV:type|value = assert type to determine InputValue_'s type
// 	Loc   *Loc_
// // }
// type Type_ struct {
// 	Constraint byte          // each on bit from right represents not-null constraint applied e.g. in nested list type [type]! is 00000010, [type!]! is 00000011, type! 00000001
// 	AST        GQLTypeProvider // AST instance of type. WHen would this be used??. Used for non-Scalar types. AST in cache(typeName), then in Type_(typeName). If not in Type_, check cache, then DB.
// 	Depth      int           // depth of nested List e.g. depth 2 is [[type]]. Depth 0 implies non-list type, depth > 0 is a list type
// 	Name_                    // type name. inherit()
// }

func (l List_) ValidateListValues(iv *Type_, d *int, maxd *int, err *[]error) {
	reqType := iv.isType() // INT, FLOAT, OBJECT, PET, MEASURE etc            note: OBJECT is for specification of a type, OBJECTVAL is an object literal for input purposes
	reqDepth := iv.Depth
	//
	// for each element in the LIST
	///
	fmt.Println("++++++++ ValidateListValues ++++++++")
	*d++ // current depth of [ in [[[]]]
	if *d > *maxd {
		*maxd = *d
	}
	for _, v := range l { // []*InputValue_ // Measure items {Name_, InputValue_}
		// what is the type of the list element. Scalar, another LIST, a OBJECT
		switch in := v.InputValueProvider.(type) {

		case List_:
			fmt.Println("++++++++ ValidateListValues List_ ++++++++")
			// maxd records maximum depth of list(d=1) [] list of lists [[]](d=2) = [[][][][]] list of lists of lists (d=3) [[[]]] = [[[][][]],[[][][][]],[[]]]
			in.ValidateListValues(iv, d, maxd, err)
			*d--

		case ObjectVals:
			fmt.Println("++++++++ ValidateListValues ObjectVals ++++++++")
			// default values in input object form { name:value name:value ... }: []*ArgumentT type ArgumentT: struct {Name_, InputValueProvider}
			// reqType is the type of the input object  - which defines the name and associated type for each item in the { }
			if *d != reqDepth {
				if reqDepth == 0 {
					*err = append(*err, fmt.Errorf(`Value %s should not be contained in a List %s`, v, v.AtPosition()))
				} else {
					*err = append(*err, fmt.Errorf(`Value %s is not at required nesting of %d %s`, v, reqDepth, v.AtPosition()))
				}
			}
			in.ValidateObjectValues(iv, err)

		default:
			fmt.Println("++++++++ ValidateListValues Default ++++++++")
			// check the item - this is matched against the type specification for the list ie. [type]
			if *d != reqDepth && v.isType() != NULL {
				if reqDepth == 0 {
					*err = append(*err, fmt.Errorf(`Value %s should not be contained in a List %s`, v, v.AtPosition()))
				} else {
					*err = append(*err, fmt.Errorf(`Value %s is not at required nesting of %d %s`, v, reqDepth, v.AtPosition()))
				}
			}
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
	//	return "@" + d.Name_.String() + d.Arguments_.String()
	return d.Name_.String() + d.Arguments_.String()
}

func (d *DirectiveT) CoerceDirectiveName() {
	d.Name_.Name = NameValue_("@" + d.Name_.String())
}

type Directives_ struct {
	Directives []*DirectiveT
}

//func (d *Directives_) CheckDirectiveLocation(location string, err *[]error) {}

func (d *Directives_) AppendDirective(s *DirectiveT) error {
	s.CoerceDirectiveName()
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

func (d *Directives_) Len() int {
	return len(d.Directives)
}

func (d *Directives_) SolicitNonScalarTypes(unresolved UnresolvedMap) {
	// for _, v := range d.Directives {
	// 	unresolved[v.Name_] = nil
	// }
}

func (d *Directives_) CheckDirectiveRef(dir NameValue_, err *[]error) {
	for _, v := range d.Directives {
		if v.Name_.String() == dir.String() {
			*err = append(*err, fmt.Errorf(`Directive "%s" references itself, is not permitted %s`, dir, v.Name_.AtPosition()))
		}
	}
}

func (d *Directives_) checkDirectiveLocation_(input DirectiveLoc, err *[]error) {
	var found bool
	for _, v := range d.Directives {
		//	get the use named directive's AST
		if e, ok := TyCache[v.Name.String()]; ok {
			found = false
			if x, ok := e.(*Directive_); ok {
				for _, loc := range x.Location {
					if loc == input {
						found = true
					}
				}
				if !found {
					if dloc, ok := DirectiveLocationMap[input]; ok {
						*err = append(*err, fmt.Errorf(`Directive "%s" is not registered for %s usage %s`, v.Name, dloc, v.Name_.AtPosition()))
					} else {
						*err = append(*err, fmt.Errorf(`System Error: Directive %s not found in map `, v.Name, dloc, v.Name_.AtPosition()))
					}
				}
			} else {
				*err = append(*err, fmt.Errorf(`AST for type %s is not a Directive_ type %s`, v.Name, v.Name_.AtPosition()))
			}
		}
	}
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

func (a NameValue_) EqualString(b string) bool {
	return string(a) == b
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

func (a Name_) EqualString(b string) bool {
	return a.Name.EqualString(b)
}

func (n Name_) AtPosition() string {
	// if n.Loc == nil {
	// 	fmt.
	// 	panic(fmt.Errorf("Error in AtPosition(), Loc not set"))
	// }
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
	ValidateName(s, errS, loc)
	n.Name = NameValue_(s)
}

// ======== Document ===================================

type Document struct {
	Statements    []GQLTypeProvider
	StatementsMap map[NameValue_]GQLTypeProvider
	ErrorMap      map[NameValue_][]error
}

func (d Document) String() string {
	var (
		s    strings.Builder
		name []string
	)
	tc = 2
	for k, _ := range d.StatementsMap {
		name = append(name, k.String())
	}
	sort.Strings(name) // conversion method to acquire sort methods and perform inplace sort

	for _, v := range name {
		stmt := d.StatementsMap[NameValue_(v)]
		s.WriteString(stmt.String())
		s.WriteString("\n")
	}
	return s.String()
}

// ======== type statements ==========

// GQLTypeProvider reperesents all the GraphQL types, SCALAR (user defined), OBJECTS, UNIONS, INTERFACES, ENUMS, INPUTOBJECTS, LISTS
type GQLTypeProvider interface {
	TypeSystemNode()
	TypeName() NameValue_
	SolicitNonScalarTypes(UnresolvedMap) // while not all Types contain nested types that need to be resolved e.g scalar must still include this method
	CheckDirectiveRef(dir NameValue_, err *[]error)
	CheckDirectiveLocation(err *[]error)
	String() string
	Type() string // equiv to IsType however IsType has been used for Input Types (subset of GQLTypeProvider
}

// ======================================================

var tc = 2

type opType byte

const (
	QUERY_OP opType = 1 << iota
	MUTATION_OP
	SUBSCRIPTION_OP
)

type Schema_ struct {
	Directives_
	Query        Name_ // named type to use as root type of query into graph of types e.g. "Query" -> type Query { allPersons(last : Int ) : [Person!]! }
	Mutation     Name_
	Subscription Name_
	Op           opType //  current operation used during parsing of statement
}

func (sc *Schema_) TypeSystemNode() {}

func (sc *Schema_) Type() string {
	return "Schema"
}

func (sc *Schema_) AssignName(s string, loc *Loc_, errS *[]error) {
	switch sc.Op {
	case QUERY_OP:
		sc.Query.AssignName(s, loc, errS)
	case MUTATION_OP:
		sc.Mutation.AssignName(s, loc, errS)
	case SUBSCRIPTION_OP:
		sc.Subscription.AssignName(s, loc, errS)
	}
}

func (sc *Schema_) CheckDirectiveLocation(err *[]error) {
	sc.checkDirectiveLocation_(SCHEMA_DL, err)
}

func (sc *Schema_) String() string {
	var s strings.Builder
	sc.Directives_.String()
	s.WriteString("schema {")
	if sc.Query.Exists() {
		s.WriteString("\nquery : ")
		s.WriteString(sc.Query.String())
	}
	if sc.Mutation.Exists() {
		s.WriteString(" \nmutation : ")
		s.WriteString(sc.Mutation.String())
	}
	if sc.Subscription.Exists() {
		s.WriteString("\nsubscription : ")
		s.WriteString(sc.Subscription.String())
	}
	s.WriteString("\n}")
	return s.String()
}

func (sc *Schema_) TypeName() NameValue_ {
	return NameValue_("schema")
}

// ======================================================

var blank string = ""
var errNameChar string = "Invalid character in identifer at line: %d, column: %d"
var errNameBegin string = "identifer [%s] cannot start with two underscores at line: %d, column: %d"

func ValidateName(name string, errS *[]error, loc *Loc_) {
	// /[_A-Za-z][_0-9A-Za-z]*/
	var err error
	if len(name) == 0 {
		err = fmt.Errorf("Error: zero length name passed to ValidateName")
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
