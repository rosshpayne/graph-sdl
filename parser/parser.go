package parser

import (
	"errors"
	"fmt"
	_ "os"
	"strings"

	"github.com/graph-sdl/ast"
	"github.com/graph-sdl/lexer"
	"github.com/graph-sdl/token"
)

type (
	parseFn func(op string) ast.TypeSystemDef
)

const (
	cErrLimit = 8 // how many parse errors are permitted before processing stops
)

//var typeRepo ast.TypeRepo_
var enumRepo ast.EnumRepo_

//var enumValue ast.EnumRepo_

type Parser struct {
	l *lexer.Lexer

	extend bool

	abort bool

	curToken  token.Token
	peekToken token.Token

	parseFns map[token.TokenType]parseFn
	perror   []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.parseFns = make(map[token.TokenType]parseFn)
	p.registerFn(token.TYPE, p.ParseObjectType)
	p.registerFn(token.ENUM, p.ParseEnumType)
	p.registerFn(token.INTERFACE, p.ParseInterfaceType)
	p.registerFn(token.UNION, p.ParseUnionType)
	p.registerFn(token.INPUT, p.ParseInputValueType)
	p.registerFn(token.SCALAR, p.ParseScalarType)
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	//fmt.Println("New 1 ", p.curToken.Literal)
	p.nextToken()
	//fmt.Println("New 1 ", p.curToken.Literal)

	return p

}

// repository of all types defined in the graph

func init() {
	//typeRepo = make(ast.TypeRepo_)
	//	typeRepo[token.STRING] = ast.String_("")
	enumRepo = make(ast.EnumRepo_)
}

func (p *Parser) Loc() *ast.Loc_ {
	loc := p.curToken.Loc
	return &ast.Loc_{loc.Line, loc.Col}
}

func (p *Parser) printToken(s ...string) {
	if len(s) > 0 {
		fmt.Printf("** Current Token: [%s] %s %s %s %v %s %s [%s]\n", s[0], p.curToken.Type, p.curToken.Literal, p.curToken.Cat, p.curToken.IsScalarType, "Next Token:  ", p.peekToken.Type, p.peekToken.Literal)
	} else {
		fmt.Println("** Current Token: ", p.curToken.Type, p.curToken.Literal, p.curToken.Cat, "Next Token:  ", p.peekToken.Type, p.peekToken.Literal)
	}
}
func (p *Parser) hasError() bool {
	if len(p.perror) > 15 || p.abort {
		return true
	}
	return false
}
func (p *Parser) addErr(s string) error {
	if strings.Index(s, " at line: ") == -1 {
		s += fmt.Sprintf(" at line: %d, column: %d", p.curToken.Loc.Line, p.curToken.Loc.Col)
	}
	e := errors.New(s)
	p.perror = append(p.perror, e)
	return e
}

func (p *Parser) registerFn(tokenType token.TokenType, fn parseFn) {
	p.parseFns[tokenType] = fn
}

func (p *Parser) nextToken(s ...string) {
	p.curToken = p.peekToken

	p.peekToken = p.l.NextToken() // get another token from lexer:    [,+,(,99,Identifier,keyword etc.
	if len(s) > 0 {
		fmt.Printf("** Current Token: [%s] %s %s %s %s %s %s\n", s[0], p.curToken.Type, p.curToken.Literal, p.curToken.Cat, "Next Token:  ", p.peekToken.Type, p.peekToken.Literal)
	}
	if p.curToken.Illegal {
		p.addErr(fmt.Sprintf("Illegal %s token, [%s]", p.curToken.Type, p.curToken.Literal))
	}
	// if $variable present then mark the identier as a VALUE
	if p.curToken.Literal == token.DOLLAR {
		p.peekToken.Cat = token.VALUE
	}
}

// ==================== Start =========================

func (p *Parser) ParseDocument() (*ast.Document, []error) {
	program := &ast.Document{}
	program.Statements = []ast.TypeSystemDef{} // slice is initialised  with no elements - each element represents an interface value of type ast.TypeSystemDef
	program.StatementsMap = make(map[ast.NameValue_]ast.TypeSystemDef)
	// Build AST from GraphQL SDL script

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if p.hasError() {
			break
		}
		if stmt != nil {
			// if no errors add to cache / DB
			if len(p.perror) == 0 {
				ast.Add2Cache(stmt.TypeName(), stmt)
			}
			//
			program.Statements = append(program.Statements, stmt)
			name := stmt.TypeName()
			program.StatementsMap[name] = stmt
		}
		if p.extend {
			p.extend = false
		}

	}
	if p.hasError() {
		return program, p.perror
	}
	// Validation - via multiple parses of AST

	var validationFail bool
	var stillUnresolved []ast.Name_
	//
	for _, v := range program.Statements {
		var unresolved []ast.Name_
		// returns slice of unresolved types from each statement in the document
		v.CheckUnresolvedTypes(&unresolved)

		// resolve unresolved types by checking in DB
		for _, v_ := range unresolved {

			if typeDef, err := ast.DBFetch(v_.Name); err != nil {
				p.addErr(err.Error())
				continue
			} else {
				if len(typeDef) == 0 {
					// not resolved - no Type exists in repo (cache + DB)
					stillUnresolved = append(stillUnresolved, v_)
					validationFail = true
				} else {
					// resolved - generate AST and add type to cache
					l := lexer.New(typeDef)
					p := New(l)
					p.ParseDocument() // TODO - is parseStatement workable
					//p.parseStatement()
				}
			}
		}
		for _, v := range stillUnresolved {
			p.addErr(fmt.Sprintf(`Unresolved type "%s" %s`, v.Name, v.AtPosition()))
		}
		if !validationFail {
			switch x := v.(type) {
			case *ast.Object_:
				x.CheckIsOutputType(&p.perror)
				x.CheckIsInputType(&p.perror)
				x.CheckInputValueType(&p.perror)
				x.CheckImplements(&p.perror) // check implements are interfaces
			case *ast.Enum_:
			case *ast.Interface_:
			}

			// pass validation - add type to repo
		}
	}

	if len(p.perror) == 0 {
		for _, v := range program.Statements {
			ast.Add(v.TypeName(), v)
		}
	}
	return program, p.perror
}

var opt bool = true // is optional

func (p *Parser) parseStatement() ast.TypeSystemDef {
	p.skipComment()
	if p.curToken.Type == token.EXTEND {
		p.extend = true
		p.nextToken() // read over extend
	}
	stmtType := p.curToken.Literal
	if f, ok := p.parseFns[p.curToken.Type]; ok {
		return f(stmtType)
	} else {
		p.abort = true
		p.addErr(fmt.Sprintf(`Parse aborted. "%s" is not a statement keyword`, p.curToken.Literal))
	}
	return nil
}

// ====================  fetchAST  =============================

func (p *Parser) fetchAST(name string) ast.TypeSystemDef {
	// search in cache, then database and parse database statement
	//  to generate the type's AST and return it
	var (
		ast_ ast.TypeSystemDef
		ok   bool
	)
	name_ := ast.NameValue_(name)
	if ast_, ok = ast.CacheFetch(name_); !ok {
		if typeDef, err := ast.DBFetch(name_); err != nil {
			p.addErr(err.Error())
			p.abort = true
			return nil
		} else {
			if len(typeDef) == 0 { // no type found in DB
				p.addErr(fmt.Sprintf(`Type "%s" does not exist %s`, name, p.Loc()))
				p.abort = true
				return nil
			} else {
				// generate the AST
				l := lexer.New(typeDef)
				p2 := New(l)
				ast_ := p2.parseStatement()
				ast.Add2Cache(name_, ast_)
			}
		}
	}
	return ast_
}

// ==================== Object Type  ============================

// Description-opt TYPE Name  ImplementsInterfaces-opt Directives-Const-opt  FieldsDefinition-opt
//		{FieldDefinition-list}
//         FieldDefinition:
//			Description-opt Name ArgumentsDefinition- opt : Type Directives-Con
func (p *Parser) ParseObjectType(op string) ast.TypeSystemDef {
	// Types: query, mutation, subscription
	p.nextToken() // read over type
	if !p.extend {
		obj := &ast.Object_{}

		p.parseName(obj).parseImplements(obj, opt).parseDirectives(obj, opt).parseFields(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj := p.parseExtendName()
		if obj != nil {
			if inp, ok := obj.(*ast.Object_); !ok {
				p.addErr(fmt.Sprintf(`specified extend type "%s" is not an Input Value Type`, obj.TypeName()))
				p.abort = true
			} else {
				icnt, dcnt, fcnt := len(inp.Implements), len(inp.Directives), len(inp.FieldSet)

				p.parseImplements(inp, opt).parseDirectives(inp, opt).parseFields(inp, opt)

				if icnt == len(inp.Implements) && dcnt == len(inp.Directives) && fcnt == len(inp.FieldSet) {
					p.addErr(fmt.Sprintf(`extend for type "%s" contains no changes`, inp.TypeName()))
				}
				return inp
			}
		}
	}
	return nil
}

// ====================  Enum Type  ============================

// EnumTypeDefinition
//		 Description-opt enum Name Directives-opt EnumValueDefinition-opt
// EnumValuesDefinition
//		{EnumValueDefinition-list}
// EnumValueDefinition
//		Description-opt EnumValue Directives-opt
func (p *Parser) ParseEnumType(op string) ast.TypeSystemDef {
	p.nextToken() // read type
	if !p.extend {
		obj := &ast.Enum_{}

		p.parseName(obj).parseDirectives(obj, opt).parseEnumValues(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj := p.parseExtendName()
		if obj != nil {
			if inp, ok := obj.(*ast.Enum_); !ok {
				p.addErr(fmt.Sprintf(`specified extend type "%s" is not an Input Value Type`, obj.TypeName()))
				p.abort = true
			} else {
				dcnt, fcnt := len(inp.Directives), len(inp.Values)

				p.parseDirectives(inp, opt).parseEnumValues(inp)

				if dcnt == len(inp.Directives) && fcnt == len(inp.Values) {
					p.addErr(fmt.Sprintf(`extend for type "%s" contains no changes`, inp.TypeName()))
				}
				return inp
			}
		}
	}
	return nil
}

// ====================== Interface ===========================
// InterfaceTypeDefinition
//		Description-opt	interface	Name	Directives-opt	FieldsDefinition-opt
func (p *Parser) ParseInterfaceType(op string) ast.TypeSystemDef {
	p.nextToken() // read over interfcae keyword
	if !p.extend {
		obj := &ast.Interface_{}

		p.parseName(obj).parseDirectives(obj, opt).parseFields(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj := p.parseExtendName()
		if obj != nil {
			if inp, ok := obj.(*ast.Interface_); !ok {
				p.addErr(fmt.Sprintf(`specified extend type "%s" is not an Input Value Type`, obj.TypeName()))
				p.abort = true
			} else {
				dcnt, fcnt := len(inp.Directives), len(inp.FieldSet)

				p.parseDirectives(inp, opt).parseFields(inp, opt)

				if dcnt == len(inp.Directives) && fcnt == len(inp.FieldSet) {
					p.addErr(fmt.Sprintf(`extend for type "%s" contains no changes`, inp.TypeName()))
				}
				return inp
			}
		}
	}
	return nil
}

// ====================== Union ===============================
// UnionTypeDefinition
//		Descriptionopt union Name  Directives-opt  UnionMemberTypes-opt
// UnionMemberTypes
//		=|-opt	NamedType
//		UnionMemberTypes | NamedType
func (p *Parser) ParseUnionType(op string) ast.TypeSystemDef {
	p.nextToken() // read over interfcae keyword
	if !p.extend {
		obj := &ast.Union_{}

		p.parseName(obj).parseDirectives(obj, opt).parseUnionMembers(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj := p.parseExtendName()
		if obj != nil {
			if inp, ok := obj.(*ast.Union_); !ok {
				p.addErr(fmt.Sprintf(`specified extend type "%s" is not an Input Value Type`, obj.TypeName()))
				p.abort = true
			} else {
				dcnt, fcnt := len(inp.Directives), len(inp.NameS)

				p.parseDirectives(inp, opt).parseUnionMembers(inp)

				if dcnt == len(inp.Directives) && fcnt == len(inp.NameS) {
					p.addErr(fmt.Sprintf(`extend for type "%s" contains no changes`, inp.TypeName()))
				}
				return inp
			}
		}
	}
	return nil
}

// ====================== Input ===============================
// InputObjectTypeDefinition
//		Description-opt	input	Name	DirectivesConst-opt	InputFieldsDefinition-opt
func (p *Parser) ParseInputValueType(op string) ast.TypeSystemDef {

	p.nextToken() // read over input keyword
	if !p.extend {
		inp := &ast.Input_{}

		p.parseName(inp).parseDirectives(inp, opt).parseInputFieldDefs(inp)

		return inp
	} else {
		// return original AST associated with the extend Name.
		obj := p.parseExtendName()
		if obj != nil {
			if inp, ok := obj.(*ast.Input_); !ok {
				p.addErr(fmt.Sprintf(`specified extend type "%s" is not an Input Value Type`, obj.TypeName()))
				p.abort = true
			} else {
				dcnt, fcnt := len(inp.Directives), len(inp.InputValueDefs)

				p.parseDirectives(inp, opt).parseInputFieldDefs(inp)

				if dcnt == len(inp.Directives) && fcnt == len(inp.InputValueDefs) {
					p.addErr(fmt.Sprintf(`extend for type "%s" contains no changes`, inp.TypeName()))
				}
				return inp
			}
		}
	}
	return nil

}

// ====================== Scalar_ ===============================
// InputObjectTypeDefinition
//		Description-opt	input	Name	DirectivesConst-opt	InputFieldsDefinition-opt
func (p *Parser) ParseScalarType(op string) ast.TypeSystemDef {

	p.nextToken() // read over input keyword
	if !p.extend {
		inp := &ast.Scalar_{}

		p.parseName(inp).parseDirectives(inp, opt)

		return inp
	} else {
		// return original AST associated with the extend Name.
		obj := p.parseExtendName()
		if obj != nil {
			if inp, ok := obj.(*ast.Input_); !ok {
				p.addErr(fmt.Sprintf(`specified extend type "%s" is not an Input Value Type`, obj.TypeName()))
				p.abort = true
			} else {
				dcnt := len(inp.Directives)

				p.parseDirectives(inp, opt)

				if dcnt == len(inp.Directives) {
					p.addErr(fmt.Sprintf(`extend for type "%s" contains no changes`, inp.TypeName()))
				}
				return inp
			}
		}
	}
	return nil

}

// =============================================================

func (p *Parser) skipComment() {
	if p.curToken.Type == token.STRING {
		p.nextToken() // read over comment string
	}
}

func (p *Parser) readDescription() string {
	var s string
	if p.curToken.Type == token.STRING {
		s = p.curToken.Literal
		p.nextToken() // read over comment string
	}
	return s
}

func (p *Parser) parseDescription() *Parser {
	if p.curToken.Type == token.STRING {
		p.nextToken() // read over comment string
	}
	return p
}

// ==================== parseName ===============================
// parseName is always mandatory
func (p *Parser) parseName(f ast.NameI) *Parser {
	// check if appropriate thing to do
	if p.hasError() {
		return p
	}
	if p.curToken.Type == token.IDENT {
		f.AssignName(p.curToken.Literal, p.Loc(), &p.perror)
	} else {
		p.addErr(fmt.Sprintf(`Expected name identifer got %s of "%s"`, p.curToken.Type, p.curToken.Literal))
	}
	p.nextToken() // read over name

	return p
}

// ==================== parseExtendName ===============================

func (p *Parser) parseExtendName() ast.TypeSystemDef {
	if p.hasError() {
		return nil
	}
	var extName string
	// does name entity exist
	if p.curToken.Type == token.IDENT {
		extName = p.curToken.Literal
	} else {
		p.addErr(fmt.Sprintf(`Expected name identifer got %s of "%s"`, p.curToken.Type, p.curToken.Literal))
	}
	if obj, ok := ast.Fetch(ast.NameValue_(extName)); !ok {
		if typeDef, err := ast.DBFetch(ast.NameValue_(extName)); err != nil {
			p.addErr(err.Error())
		} else {
			if len(typeDef) == 0 { // no type found in DB
				p.addErr(fmt.Sprintf(`Cannot extend, type "%s" does not exist %s`, extName, p.Loc()))
				p.abort = true
			} else {
				// generate the AST
				l := lexer.New(typeDef)
				p2 := New(l)
				obj = p2.parseStatement()
			}
		}
		p.nextToken() // read over name
		return obj
	} else {
		p.nextToken() // read over name
		return obj
	}

}

func (p *Parser) parseEnumValues(enum *ast.Enum_, optional ...bool) *Parser {

	if p.hasError() || p.curToken.Type != token.LBRACE {
		return p
	}

	for p.nextToken(); p.curToken.Type != token.RBRACE; {

		ev := &ast.EnumValue_{}

		_ = p.parseDescription().parseName(ev).parseDirectives(ev, opt)

		if p.hasError() {
			break
		}
		// check for duplicate values
		for _, v := range enum.Values {
			if ev.Loc == nil {
				continue
			}
			if v.Name_.String() == ev.Name_.String() {
				loc := ev.Name_.Loc
				p.addErr(fmt.Sprintf("Duplicate Enum Value [%s] at line: %d column: %d", ev.Name_.String(), loc.Line, loc.Column))
				//return p
			}
		}
		enum.Values = append(enum.Values, ev)
		enumRepo[string(ev.Name)+"|"+string(enum.Name)] = struct{}{}

	}
	p.nextToken() // read over }
	return p
}

//==============================================================================
// UnionMemberTypes
//		=|optNamedType
//		UnionMemberTypes | NamedType
func (p *Parser) parseUnionMembers(u *ast.Union_, optional ...bool) *Parser {

	if p.hasError() || p.curToken.Type != token.ASSIGN {
		return p
	}
	for p.nextToken(); p.curToken.Type == token.BAR || p.curToken.Type == token.IDENT; p.nextToken() {
		if p.curToken.Type == token.BAR && p.peekToken.Type != token.IDENT {
			p.addErr(fmt.Sprintf("expected Union identifer, got %s, %s", p.curToken.Type, p.curToken.Literal))
		} else if p.curToken.Type == token.IDENT {
			var memberName ast.Name_
			memberName.AssignName(p.curToken.Literal, p.Loc(), &p.perror) // appends to perror if invalid name
			for _, v := range u.NameS {
				if v.String() == memberName.String() {
					loc := memberName.Loc
					p.addErr(fmt.Sprintf("Duplicate member name at line: %d column: %d", loc.Line, loc.Column))
				}
			}
			u.NameS = append(u.NameS, memberName) // save string component of Name_
		}
	}
	return p
}

//==============================================================================

// parseImplements
// ImplementsInterfaces
//		implements &-opt NamedType
//		ImplementsInterfaces & NamedType
//func (p *Parser) parseImplements(f ast.ImplementI, optional ...bool) *Parser {
func (p *Parser) parseImplements(o *ast.Object_, optional ...bool) *Parser {

	if p.hasError() || p.curToken.Type != token.IMPLEMENTS {
		return p
	}
	for p.nextToken(); p.curToken.Type == token.AND || p.curToken.Type == token.IDENT; p.nextToken() {
		if p.curToken.Type == token.AND && p.peekToken.Type != token.IDENT {
			p.addErr(fmt.Sprintf("expected interface identifer, got %s, %s", p.curToken.Type, p.curToken.Literal))
		} else if p.curToken.Type == token.IDENT {
			var impName ast.Name_
			impName.AssignName(p.curToken.Literal, p.Loc(), &p.perror) // appends to perror if invalid name
			for _, v := range o.Implements {
				if v.String() == impName.String() {
					loc := impName.Loc
					p.addErr(fmt.Sprintf("Duplicate interface name at line: %d column: %d", loc.Line, loc.Column))
				}
			}
			o.Implements = append(o.Implements, impName) // save string component of Name_
		}
	}
	return p
}

// Directives[Const]
// 		Directive[?Const]list
// Directive[Const] :
// 		@ Name Arguments[?Const]opt ...
// hero(episode: $episode) {
//     name
//     friends @include(if: $withFriends) @ Size (aa:1 bb:2) @ Pack (filter: true) {
//       name
//     }
func (p *Parser) parseDirectives(f ast.DirectiveI, optional ...bool) *Parser { // f is a iv initialised from concrete types *ast.Field,*OperationStmt,*FragementStmt. It will panic if they don't satisfy DirectiveI

	if p.hasError() {
		return p
	}
	if p.curToken.Type != token.ATSIGN {
		if len(optional) == 0 {
			p.addErr("Variable is mandatory")
		}
		return p
	}
	for p.curToken.Type == token.ATSIGN {
		p.nextToken() // read over @
		a := []*ast.ArgumentT{}
		d := &ast.DirectiveT{Arguments_: ast.Arguments_{Arguments: a}}

		p.parseName(d).parseArguments(d, opt)
		if err := f.AppendDirective(d); err != nil {
			p.addErr(err.Error())
		}
	}
	return p
}

//
// type Argument struct {
// 	//( name:value )
// 	Name  Name_
// 	Value []InputValue_ // could use string as this value is mapped directly to get function - at this stage we don't care about its type, maybe?
// }
//
// Arguments[Const] :
//		( Argument[?Const]list )
// Argument[Const] :
//		Name : Value [?Const]
// only fields have arguments so not interface argument is necessary to support multiple types
func (p *Parser) parseArguments(f ast.ArgumentI, optional ...bool) *Parser {

	if p.hasError() {
		return p
	}
	if p.curToken.Type != token.LPAREN {
		if p.curToken.Type == token.LBRACKET {
			p.addErr(fmt.Sprintf("Expected a ( or } or { instead got %s", p.curToken.Type))
		} else {
			if len(optional) == 0 {
				p.addErr(fmt.Sprintf("Expected a ( instead got %s", p.curToken.Type))
			}
			return p
		}
	}
	//
	for p.nextToken(); p.curToken.Type != token.RPAREN; p.nextToken() {
		v := new(ast.ArgumentT)

		p.parseName(v).parseColon().parseInputValue(v)

		f.AppendArgument(v)
	}
	p.nextToken() // read over )
	return p
}

func (p *Parser) parseColon() *Parser {

	if !(p.curToken.Type == token.COLON) {
		p.addErr(fmt.Sprintf(`Expected a colon got an "%s"`, p.curToken.Literal))
	}
	p.nextToken() // read over :
	return p
}

func (p *Parser) parseInputValue(v *ast.ArgumentT) *Parser {
	if p.hasError() {
		return p
	}
	if !((p.curToken.Cat == token.VALUE && (p.curToken.Type == token.DOLLAR && p.peekToken.Cat == token.VALUE)) ||
		(p.curToken.Cat == token.VALUE && (p.peekToken.Cat == token.NONVALUE || p.peekToken.Type == token.RPAREN)) ||
		(p.curToken.Type == token.LBRACKET || p.curToken.Type == token.LBRACE)) { // [  or {
		p.addErr(fmt.Sprintf(`Expected an argument Value followed by IDENT or RPAREN got an %s:%s:%s %s:%s:%s`, p.curToken.Cat, p.curToken.Type, p.curToken.Literal, p.peekToken.Cat, p.peekToken.Type, p.peekToken.Literal))
	}
	v.Value = p.parseInputValue_()

	return p
}

//==============================================================================
// {FieldDefinition ...} :
// .  Description-opt Name ArgumentsDefinition-opt : Type Directives-opt
func (p *Parser) parseFields(f ast.FieldI, optional ...bool) *Parser {

	if p.hasError() || p.curToken.Type != token.LBRACE {
		if len(optional) == 0 {
			p.addErr("Field definitions is required")
		}
		return p
	}
	for p.nextToken(); p.curToken.Type != token.RBRACE; { // p.nextToken("next token in parseFields..") {

		field := &ast.Field_{}

		_ = p.parseDecription().parseName(field).parseFieldArgumentDefs(field).parseType(field).parseDirectives(field, opt)

		if p.hasError() {
			return p
		}
		if err := f.AppendField(field); err != nil {
			p.addErr(err.Error())
		}
	}
	p.nextToken() // read over }
	return p
}

func (p *Parser) parseDecription() *Parser {
	p.skipComment()
	return p
}

func (p *Parser) parseType(f ast.HasTypeI) *Parser {
	if p.hasError() {
		return p
	}
	if p.curToken.Type == token.COLON {
		p.nextToken() // read over :
	}
	if !p.curToken.IsScalarType { // ie not a Int, Float, String, Boolean, ID, <namedType>
		if !(p.curToken.Type == token.IDENT || p.curToken.Type == token.LBRACKET) {
			p.addErr(fmt.Sprintf("Expected a Type, got %s, %s", p.curToken.Type, p.curToken.Literal))
		}
	}
	var (
		bit  byte
		name string
		ast_ ast.TypeSystemDef
		//typedef ast.TypeFlag_ // token defines SCALAR types only. All other types will be populated in repoType map.
		depth   int
		nameLoc *ast.Loc_
	)
	nameLoc = p.Loc()
	switch p.curToken.Type {

	case token.LBRACKET:
		// [ typeName ]
		var (
			depthClose uint
		)
		p.nextToken() // read over [
		for depth = 1; p.curToken.Type == token.LBRACKET; p.nextToken() {
			depth++
		}
		if depth > 7 {
			p.addErr("Nested list type cannot be greater than 8 deep ")
			break
		}
		if !(p.curToken.Type == token.IDENT || p.curToken.IsScalarType) {
			p.addErr(fmt.Sprintf("Expected type identifer got %s, %s", p.curToken.Type, p.curToken.Literal))
			break
		}
		nameLoc = p.Loc()
		name = p.curToken.Literal // actual type name, Int, Float, Pet ...
		// System ScalarTypes are defined by the Type_.Name_, Non-system Scalar and non-scalar are defined by the AST.
		if !p.curToken.IsScalarType {
			ast_ = p.fetchAST(string(name))
		}
		p.nextToken() // read over IDENT
		for bangs := 0; p.curToken.Type == token.RBRACKET || p.curToken.Type == token.BANG; {
			if p.curToken.Type == token.BANG {
				bangs++
				if bangs > depth+1 {
					p.addErr("redundant !")
					p.nextToken() // read over !
					//return p
				} else {
					bit |= (1 << depthClose)
					p.nextToken() // read over !
				}
			} else {
				depthClose++
				p.nextToken() // read over ]
			}
		}
		if depth != int(depthClose) {
			p.addErr("close ] does not match opening [ in type specification")
			return p
		}
	default:
		// typeName
		//if p.curToken.IsScalarType {
		if p.curToken.Type == token.IDENT || p.curToken.IsScalarType {
			name = p.curToken.Literal
			if !p.curToken.IsScalarType {
				ast_ = p.fetchAST(string(name))
			}
			if p.peekToken.Type == token.BANG {
				bit = 1 << 0
				p.nextToken() // read over IDENT
			}
			p.nextToken() // read over ! or IDENT
		} else {
			p.addErr(fmt.Sprintf("Expected type identifer got %s, %s %v", p.curToken.Type, p.curToken.Literal, p.curToken.IsScalarType))
		}
	}

	if p.hasError() {
		return p
	}
	// name is the type name Int, Person, [name], ...
	t := &ast.Type_{Constraint: bit, Depth: depth, AST: ast_}
	t.AssignName(name, nameLoc, &p.perror)
	f.AssignType(t) // assign the name of the named type. Later pass of AST will confirm if the named type has been defined.
	return p

}

// ArgumentsDefinition
//		(InputValueDefinitionlist)
// InputValueDefinition
//		Description-opt Name : Type DefaultValue-opt Directives-opt
// DefaultValue
//		= Value
// type FieldArgument_ struct {
// 	Desc       string
// 	Name       Name_
// 	Type       Type_
// 	DefValue   *InputValueDef_
// 	Directives []*Directive_
// }
// InputFieldsDefinition
//		{InputValueDefinition-list}
// InputValueDefinition
//		Description-opt	Name	:	Type	DefaultValue-opt	Directives-opt
func (p *Parser) parseFieldArgumentDefs(f ast.FieldArgI) *Parser { // st is an iv initialised from passed in argument which is a *OperationStmt

	if p.hasError() {
		return p
	}
	var encl [2]token.TokenType = [2]token.TokenType{token.LPAREN, token.RPAREN}
	return p.parseInputValueDefs(f, encl)
}

func (p *Parser) parseInputFieldDefs(f ast.FieldArgI) *Parser {

	if p.hasError() {
		return p
	}
	var encl [2]token.TokenType = [2]token.TokenType{token.LBRACE, token.RBRACE}
	return p.parseInputValueDefs(f, encl)
}

func (p *Parser) parseInputValueDefs(f ast.FieldArgI, encl [2]token.TokenType) *Parser {

	if p.curToken.Type == encl[0] {
		p.nextToken() // read over ( or {
		for p.curToken.Type != encl[1] {

			v := &ast.InputValueDef{}
			v.Loc = p.Loc()

			p.parseDecription().parseName(v).parseType(v).parseDefaultVal(v, opt).parseDirectives(v, opt)

			if p.hasError() {
				return p
			}
			f.AppendField(v, &p.perror)
		}
		p.nextToken() //read over )
	}

	return p
}

func (p *Parser) parseDefaultVal(v *ast.InputValueDef, optional ...bool) *Parser {

	if p.hasError() {
		return p
	}

	if p.curToken.Type != token.ASSIGN {
		if len(optional) == 0 {
			p.addErr("Default Value is mandatory")
		}
		return p
	}

	if p.curToken.Type == token.ASSIGN {
		//	p.nextToken() // read over Datatype
		p.nextToken() // read over ASSIGN
		v.DefaultVal = p.parseInputValue_(v)
		p.nextToken() // read over input value
	}
	return p
}

// parseObjectArguments - used for input object values
func (p *Parser) parseObjectArguments(argS []*ast.ArgumentT) []*ast.ArgumentT {

	for p.nextToken(); p.curToken.Type == token.IDENT; {
		//for p.nextToken(); p.curToken.Type != token.RBRACE ; p.nextToken() { // TODO: use this
		v := new(ast.ArgumentT)

		p.parseName(v).parseColon().parseInputValue(v)

		argS = append(argS, v)

		if p.curToken.Type == token.RBRACE {
			p.nextToken() // read over }
			break
		}
		p.nextToken()
	}
	return argS
}

// parseInputValue_ used to interpret "default value" in argument and field values.
//  parseInputValue_ expects an InputValue_ literal (true,false, 234, 23.22, "abc" or $variable in the next token.  The value is a type bool,int,flaot,string..
//  if it is a variable then the variable value (which is an InputValue_ type) will be sourced
//  TODO: currently called from parseArgument only. If this continues to be the case then add this func as anonymous func to it.
func (p *Parser) parseInputValue_(iv ...*ast.InputValueDef) *ast.InputValue_ {
	if p.curToken.Cat != token.VALUE {
		// maybe its a enum value in ENUM type iv.Type_.Name
		if len(iv) != 0 {
			iv := iv[0]
			p.printToken(p.curToken.Literal + "|" + string(iv.Type.Name))
			if _, ok := enumRepo[p.curToken.Literal+"|"+string(iv.Type.Name)]; !ok {
				p.addErr(fmt.Sprintf("Value expected got %s of %s", p.curToken.Type, p.curToken.Literal))
				return &ast.InputValue_{}
			}
		} else {
			p.addErr(fmt.Sprintf("Value expected got %s of %s", p.curToken.Type, p.curToken.Literal))
			return &ast.InputValue_{}
		}
	}
	switch p.curToken.Type {

	// case token.DOLLAR:
	// 	// variable supplied - need to fetch value
	// 	p.nextToken() // IDENT variable name
	// 	// change category of token to VALUE as previous token was $ - otherwise this step would not be executed.
	// 	p.curToken.Cat = token.VALUE
	// 	if p.curToken.Type == token.IDENT {
	// 		// get variable value....
	// 		if val, ok := p.getVarValue(p.curToken.Literal); !ok {
	// 			return ast.InputValue_{}, p.addErr(fmt.Sprintf("Variable, %s not defined ", p.curToken.Literal))
	// 		} else {
	// 			return val, nil
	// 		}
	// 	} else {
	// 		return ast.InputValue_{}, p.addErr(fmt.Sprintf("Expected Variable Name Identifer got %s", p.curToken.Type))
	// 	}

	case token.LBRACKET:
		// [ value value value .. ]
		p.nextToken() // read over [
		if p.curToken.Cat != token.VALUE {
			p.addErr(fmt.Sprintf("Expect an Input Value followed by another Input Value or a ], got %s %s ", p.curToken.Literal, p.peekToken.Literal))
			return &ast.InputValue_{}
		}
		// edge case: empty, []
		if p.curToken.Type == token.RBRACKET {
			p.nextToken() // ]
			var null ast.Null_ = true
			iv := ast.InputValue_{Value: null, Loc: p.Loc()}
			return &iv
		}
		// process list of values - all value types should be the same
		var vallist ast.List_
		for ; p.curToken.Type != token.RBRACKET; p.nextToken() {
			v := p.parseInputValue_()
			vallist = append(vallist, v)
		}
		// completed processing values, return List type
		iv := ast.InputValue_{Value: vallist, Loc: p.Loc()}
		return &iv

	case token.LBRACE:
		//  { name:value name:value ... }
		var ObjList ast.ObjectVals // []*ArgumentT {Name_,Value *InputValue_}
		for p.curToken.Type != token.RBRACE {
			ObjList = p.parseObjectArguments(ObjList)
			if p.hasError() {
				return &ast.InputValue_{}
			}
			if p.curToken.Type == token.RBRACE {
				break
			}
			p.nextToken()
		}
		iv := ast.InputValue_{Value: ObjList, Loc: p.Loc()}
		return &iv

	case token.NULL:
		var null ast.Null_ = true
		iv := ast.InputValue_{Value: null, Loc: p.Loc()}
		return &iv
	case token.INT:
		i := ast.Int_(p.curToken.Literal)
		iv := ast.InputValue_{Value: i, Loc: p.Loc()}
		return &iv
	case token.FLOAT:
		f := ast.Float_(p.curToken.Literal)
		iv := ast.InputValue_{Value: f, Loc: p.Loc()}
		return &iv
	case token.STRING:
		f := ast.String_(p.curToken.Literal)
		iv := ast.InputValue_{Value: &f, Loc: p.Loc()}
		return &iv
	case token.RAWSTRING:
		f := ast.RawString_(p.curToken.Literal)
		iv := ast.InputValue_{Value: &f, Loc: p.Loc()}
		return &iv
	case token.TRUE, token.FALSE: //token.BOOLEAN:
		b := ast.Bool_(p.curToken.Literal)
		iv := ast.InputValue_{Value: b, Loc: p.Loc()}
		return &iv
	default:
		// possible ENUM value
		b := &ast.EnumValue_{}
		b.AssignName(string(p.curToken.Literal), p.Loc(), &p.perror)
		iv := ast.InputValue_{Value: b, Loc: p.Loc()}
		return &iv
	}
	return nil

}
