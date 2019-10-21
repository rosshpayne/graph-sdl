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
	parseFn func(op string) ast.GQLTypeProvider
)

const (
	cErrLimit = 8 // how many parse errors are permitted before processing stops
)

var (
	//	enumRepo      ast.EnumRepo_
	typeNotExists map[ast.NameValue_]bool
)

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
	p.registerFn(token.SCALAR, p.ParseScalar)
	// Read two tokens, to initialise curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// repository of all types defined in the graph

func init() {
	//	enumRepo = make(ast.EnumRepo_)
	typeNotExists = make(map[ast.NameValue_]bool)
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
	if len(p.perror) > 7 || p.abort {
		return true
	}
	return false
}

// addErr appends to error slice held in parser.
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

func (p *Parser) ParseDocument() (program *ast.Document, errs []error) {
	var holderr []error
	program = &ast.Document{}
	program.Statements = []ast.GQLTypeProvider{} // slice is initialised  with no elements - each element represents an interface value of type ast.GQLTypeProvider
	program.StatementsMap = make(map[ast.NameValue_]ast.GQLTypeProvider)
	program.ErrorMap = make(map[ast.NameValue_][]error)

	defer func() {
		//
		//p.perror = nil
		p.perror = append(p.perror, holderr...)
		for _, v := range program.Statements {
			p.perror = append(p.perror, program.ErrorMap[v.TypeName()]...)
		}
		// persist error free statements to db
		for _, v := range program.Statements {
			if len(program.ErrorMap[v.TypeName()]) == 0 {
				// TODO - what if another type by that name exists
				//  auto overrite or raise an error
				ast.Persist(v.TypeName(), v)
			}
		}
		errs = p.perror
	}()

	//
	// parse phase - 	build AST from GraphQL SDL script
	//
	for p.curToken.Type != token.EOF {
		StmtAST := p.parseStatement()

		// handle any abort error
		if p.hasError() {
			return program, p.perror
		}
		if StmtAST != nil {
			program.Statements = append(program.Statements, StmtAST)

			name := StmtAST.TypeName()
			program.StatementsMap[name] = StmtAST
			program.ErrorMap[name] = p.perror
			if len(p.perror) == 0 {
				ast.Add2Cache(StmtAST.TypeName(), StmtAST)
			}
			p.perror = nil

		} else {
			// for no statements hold errors
			holderr = p.perror
			p.perror = nil
		}
		if p.extend {
			p.extend = false
		}
	}
	//
	// validate phase 1 - resolve types
	//
	for _, v := range program.Statements {
		p.checkUnresolvedTypes_(v)
		if len(p.perror) > 0 {
			program.ErrorMap[v.TypeName()] = append(program.ErrorMap[v.TypeName()], p.perror...)
			p.perror = nil
		}
	}
	//
	// Build perror from statement errors to use in hasError() counting
	//
	p.perror = holderr
	for _, v := range program.Statements {
		p.perror = append(p.perror, program.ErrorMap[v.TypeName()]...)
	}
	if p.hasError() {
		p.perror = nil
		return program, p.perror
	}
	//
	// validate phase 2
	//
	p.perror = nil
	for _, v := range program.Statements {
		// only proceed if zero errors for stmt
		if len(program.ErrorMap[v.TypeName()]) == 0 {
			switch x := v.(type) {
			case *ast.Object_:
				x.CheckIsOutputType(&p.perror)
				x.CheckIsInputType(&p.perror)
				x.CheckInputValueType(&p.perror)
				x.CheckImplements(&p.perror) // check implements are interfaces
			case *ast.Enum_:
			case *ast.Interface_:
			}
		}
		program.ErrorMap[v.TypeName()] = append(program.ErrorMap[v.TypeName()], p.perror...)
		p.perror = nil
	}

	return program, p.perror
}

// for _, v := range program.ErrorMap {
// 	for _, x := range v {
// 		fmt.Println("xErr: ", x.Error())
// 	}
// }

// ================ parseStatement ==========================

// parseStatement takes predefined parser routine and applies it to a valid statement
func (p *Parser) parseStatement() ast.GQLTypeProvider {
	p.skipComment()
	if p.curToken.Type == token.EXTEND {
		p.extend = true
		p.nextToken("read ver...") // read over extend
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
// fetchAST should only be used after all statements have been passed
//  As each statement is parsed its types are added to the cache
//  During validation phase each type is checked for existence using this func.
//  if not in cache then looks at DB for types that have been predefined.
func (p *Parser) fetchAST(name ast.Name_) ast.GQLTypeProvider {
	var (
		ast_ ast.GQLTypeProvider
		ok   bool
	)
	name_ := name.Name
	if ast_, ok = ast.CacheFetch(name_); !ok {
		if !typeNotExists[name_] {
			if typeDef, err := ast.DBFetch(name_); err != nil {
				p.addErr(err.Error())
				p.abort = true
				return nil
			} else {
				if len(typeDef) == 0 { // no type found in DB
					// mark type as being nonexistent
					typeNotExists[name_] = true
					return nil
				} else {
					// generate the AST
					l := lexer.New(typeDef)
					p2 := New(l)
					ast_ = p2.parseStatement()
					ast.Add2Cache(name_, ast_)
				}
			}
		} else {
			return nil
		}
	}
	return ast_
}

// ===================  checkUnresolvedTypes_  ==========================
// checkUnresolvedTypes_ is a validation check performed after parsing completed
//  unresolved Types from parsed types are then checked in DB.
//  check performed across nested types until all leaf finsihed or unresolved found
func (p *Parser) checkUnresolvedTypes_(v ast.GQLTypeProvider) {
	//returns slice of unresolved types from the statement passed in
	unresolved := make(ast.UnresolvedMap)
	v.CheckUnresolvedTypes(unresolved)

	//  unresolved should only contain non-scalar types known upto that point.
	for tyName, ty := range unresolved { // unresolvedMap: [name]*Type
		ast_ := p.fetchAST(tyName)
		// type ENUM values will have nil *Type
		if ast_ != nil {
			if ty != nil {
				ty.AST = ast_
				// if not scalar then check for unresolved types in nested type
				if !ty.IsScalar() {
					p.checkUnresolvedTypes_(ast_)
				}
			}

		} else {
			// nil ast_ means not found in db
			if ty != nil {
				p.addErr(fmt.Sprintf(`Type "%s" does not exist %s`, ty.Name, ty.AtPosition()))
			} else {
				p.addErr(fmt.Sprintf(`Type "%s" does not exist %s`, tyName, tyName.AtPosition()))
			}
		}
	}
}

var opt bool = true // is optional

// ==================== Object Type  ============================

// Description-opt TYPE Name  ImplementsInterfaces-opt Directives-Const-opt  FieldsDefinition-opt
//		{FieldDefinition-list}
//         FieldDefinition:
//			Description-opt Name ArgumentsDefinition- opt : Type Directives-Con
func (p *Parser) ParseObjectType(op string) ast.GQLTypeProvider {
	// Types: query, mutation, subscription
	p.nextToken() // read over type
	if !p.extend {
		obj := &ast.Object_{}

		p.parseName(obj).parseImplements(obj, opt).parseDirectives(obj, opt).parseFields(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj, name := p.parseExtendName()
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
		} else {
			p.addErr(fmt.Sprintf(`Type "%s" does not exist %s`, p.curToken.Literal, p.Loc()))
			return &ast.Object_{Name_: name}
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
func (p *Parser) ParseEnumType(op string) ast.GQLTypeProvider {
	p.nextToken() // read type
	if !p.extend {
		obj := &ast.Enum_{}

		p.parseName(obj).parseDirectives(obj, opt).parseEnumValues(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj, name := p.parseExtendName()
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
		} else {
			p.addErr(fmt.Sprintf(`Type "%s" does not exist %s`, p.curToken.Literal, p.Loc()))
			return &ast.Enum_{Name_: name}
		}
	}
	return nil
}

// ====================== Interface ===========================
// InterfaceTypeDefinition
//		Description-opt	interface	Name	Directives-opt	FieldsDefinition-opt
func (p *Parser) ParseInterfaceType(op string) ast.GQLTypeProvider {
	p.nextToken() // read over interfcae keyword
	if !p.extend {
		obj := &ast.Interface_{}

		p.parseName(obj).parseDirectives(obj, opt).parseFields(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj, name := p.parseExtendName()
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
		} else {
			p.addErr(fmt.Sprintf(`Type "%s" does not exist %s`, p.curToken.Literal, p.Loc()))
			return &ast.Interface_{Name_: name}
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
func (p *Parser) ParseUnionType(op string) ast.GQLTypeProvider {
	p.nextToken() // read over interfcae keyword
	if !p.extend {
		obj := &ast.Union_{}

		p.parseName(obj).parseDirectives(obj, opt).parseUnionMembers(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj, name := p.parseExtendName()
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
		} else {
			p.addErr(fmt.Sprintf(`Type "%s" does not exist %s`, p.curToken.Literal, p.Loc()))
			return &ast.Union_{Name_: name}
		}
	}
	return nil
}

//====================== Input ===============================
// InputObjectTypeDefinition
//		Description-opt	input	Name	DirectivesConst-opt	InputFieldsDefinition-opt
func (p *Parser) ParseInputValueType(op string) ast.GQLTypeProvider {

	p.nextToken() // read over input keyword
	if !p.extend {
		inp := &ast.Input_{}

		p.parseName(inp).parseDirectives(inp, opt).parseInputFieldDefs(inp)

		return inp
	} else {
		// return original AST associated with the extend Name.
		obj, name := p.parseExtendName()
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
		} else {
			p.addErr(fmt.Sprintf(`Type "%s" does not exist %s`, p.curToken.Literal, p.Loc()))
			p.abort = true
			return &ast.Input_{Name_: name}
		}
	}
	return nil

}

// ====================== Scalar_ ===============================
// InputObjectTypeDefinition
//		Description-opt	scalar	Name	DirectivesConst-opt
func (p *Parser) ParseScalar(op string) ast.GQLTypeProvider {

	p.nextToken() // read over input keyword
	if !p.extend {
		inp := &ast.Scalar_{}

		p.parseName(inp).parseDirectives(inp, opt)

		return inp
	} else {
		// return original AST associated with the extend Name.
		obj, name := p.parseExtendName()
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
		} else {
			p.addErr(fmt.Sprintf(`Type "%s" does not exist %s`, p.curToken.Literal, p.Loc()))
			return &ast.Scalar_{Name: name.String()}
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
func (p *Parser) parseName(f ast.NameAssigner) *Parser {
	// check if appropriate thing to do
	if p.hasError() {
		return p
	}
	if p.curToken.Type == token.IDENT {
		f.AssignName(p.curToken.Literal, p.Loc(), &p.perror)
	} else {
		p.addErr(fmt.Sprintf(`Expected name identifer got %s of "%s"`, p.curToken.Type, p.curToken.Literal))
		if p.curToken.Type == "ILLEGAL" {
			p.abort = true
		}
	}
	p.nextToken() // read over name

	return p
}

// ==================== parseExtendName ===============================
// parseExtendName will consume the type name to be extended. Returns the type's AST.
func (p *Parser) parseExtendName() (ast.GQLTypeProvider, ast.Name_) {
	// if p.hasError() {
	// 	return nil,
	// }
	var extName string
	// does name entity exist
	if p.curToken.Type == token.IDENT {
		extName = p.curToken.Literal
	} else {
		p.addErr(fmt.Sprintf(`Expected name identifer got %s of "%s"`, p.curToken.Type, p.curToken.Literal))
		if p.curToken.Type == "ILLEGAL" {
			p.abort = true
		}
	}
	name_ := ast.Name_{Name: ast.NameValue_(extName), Loc: p.Loc()}
	ast := p.fetchAST(name_)
	if ast != nil {
		p.nextToken() // read over name
	}
	return ast, name_
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
		//enumRepo[string(ev.Name)+"|"+string(enum.Name)] = struct{}{}

	}
	p.nextToken() // read over }
	return p
}

//========================= parseUnionMembers ====================================
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

//=========================== parseImplements =====================================

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

//============================ parseDirectives ========================================

// Directives[Const]
// 		Directive[?Const]list
// Directive[Const] :
// 		@ Name Arguments[?Const]opt ...
// hero(episode: $episode) {
//     name
//     friends @include(if: $withFriends) @ Size (aa:1 bb:2) @ Pack (filter: true) {
//       name
//     }
func (p *Parser) parseDirectives(f ast.DirectiveAppender, optional ...bool) *Parser { // f is a iv initialised from concrete types *ast.Field,*OperationStmt,*FragementStmt. It will panic if they don't satisfy DirectiveAppender

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
func (p *Parser) parseArguments(f ast.ArgumentAppender, optional ...bool) *Parser {

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

//======================== parseFields ==========================================
// {FieldDefinition ...} :
// .  Description-opt Name ArgumentsDefinition-opt : Type Directives-opt
func (p *Parser) parseFields(f ast.FieldAppender, optional ...bool) *Parser {

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

// ===================== parseType ===========================

func (p *Parser) parseType(f ast.AssignTyper) *Parser {
	if p.hasError() {
		return p
	}
	if p.curToken.Type == token.COLON {
		p.nextToken() // read over :
	} else {
		p.addErr(fmt.Sprintf("Colon expected got %s of %s", p.curToken.Type, p.curToken.Literal))
	}
	if !p.curToken.IsScalarType { // ie not a Int, Float, String, Boolean, ID, <namedType>
		if !(p.curToken.Type == token.IDENT || p.curToken.Type == token.LBRACKET) {
			p.addErr(fmt.Sprintf("Expected a Type, got %s, %s", p.curToken.Type, p.curToken.Literal))
		} else if p.curToken.Type == "ILLEGAL" {
			p.abort = true
			return p
		}
	}
	var (
		bit  byte
		name string
		//	ast_ ast.GQLTypeProvider
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
		// name_ := ast.Name_{Name: ast.NameValue_(name), Loc: nameLoc}
		// //System ScalarTypes are defined by the Type_.Name_, Non-system Scalar and non-scalar are defined by the AST.
		// if !p.curToken.IsScalarType {
		// 	ast_ = p.fetchAST(name_)
		// }
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
		if p.curToken.Type == token.IDENT || p.curToken.IsScalarType {
			name = p.curToken.Literal
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
	t := &ast.Type_{Constraint: bit, Depth: depth} //, AST: ast_}
	t.AssignName(name, nameLoc, &p.perror)
	f.AssignType(t) // assign the name of the named type. Later type validation pass of AST will confirm if the named type exists.
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
func (p *Parser) parseFieldArgumentDefs(f ast.FieldArgAppender) *Parser { // st is an iv initialised from passed in argument which is a *OperationStmt

	if p.hasError() {
		return p
	}
	var encl [2]token.TokenType = [2]token.TokenType{token.LPAREN, token.RPAREN} // ()
	return p.parseInputValueDefs(f, encl)
}

func (p *Parser) parseInputFieldDefs(f ast.FieldArgAppender) *Parser {

	if p.hasError() {
		return p
	}
	var encl [2]token.TokenType = [2]token.TokenType{token.LBRACE, token.RBRACE} // {}
	return p.parseInputValueDefs(f, encl)
}

func (p *Parser) parseInputValueDefs(f ast.FieldArgAppender, encl [2]token.TokenType) *Parser {

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
		//v.DefaultVal = p.parseInputValue_(v)
		v.DefaultVal = p.parseInputValue_()
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
//func (p *Parser) parseInputValue_(iv ...*ast.InputValueDef) *ast.InputValue_ { //TODO remove iv argeument now redundant
func (p *Parser) parseInputValue_() *ast.InputValue_ {
	if p.curToken.Type == "ILLEGAL" {
		p.addErr(fmt.Sprintf("Value expected got %s of %s", p.curToken.Type, p.curToken.Literal))
		p.abort = true
		return nil
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
	//
	//  List type
	//
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
	//
	//  Object type
	//
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
	//
	//  Standard Scalar types
	//
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
		iv := ast.InputValue_{Value: f, Loc: p.Loc()}
		return &iv
	case token.RAWSTRING:
		f := ast.RawString_(p.curToken.Literal)
		iv := ast.InputValue_{Value: f, Loc: p.Loc()}
		return &iv
	case token.TRUE, token.FALSE: //token.BOOLEAN:
		b := ast.Bool_(p.curToken.Literal)
		iv := ast.InputValue_{Value: b, Loc: p.Loc()}
		return &iv
	// case token.Time:
	// 	b := ast.Time_(p.curToken.Literal)
	// 	iv := ast.InputValue_{Value: b, Loc: p.Loc()}
	// 	return &iv
	default:
		// possible ENUM value
		b := &ast.EnumValue_{}
		b.AssignName(string(p.curToken.Literal), p.Loc(), &p.perror)
		iv := ast.InputValue_{Value: b, Loc: p.Loc()}
		return &iv
	}
	return nil

}
