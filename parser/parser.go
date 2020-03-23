package parser

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/graph-sdl/ast"
	"github.com/graph-sdl/db"
	"github.com/graph-sdl/lexer"
	"github.com/graph-sdl/token"
)

const (
	cErrLimit  = 8 // how many parse errors are permitted before processing stops
	Executable = 'E'
	TypeSystem = 'T'
	defaultDoc = "DefaultDoc"
	logrFlags  = log.LstdFlags | log.Lshortfile
)

// Error exit codes
const (
	FATAL int = 0
)

type stateT uint8

// Parse State
const (
	_ stateT = iota
	parseOperationTypes_
	parseArguments_
	parseFields_
	parseArgumentDefs_
	parseObjectArguments_
	parseFieldArgumentDefs_
	parseInputFieldDefs_
	//
	parseObjectType
	parseEnumType
	parseInterfaceType
	parseUnionType
	parseInputValueType
	parseScalar
	parseDirective
	parseSchema
	//
	parseEnumValues_
	parseDirectiveLocations_
	parseUnionMembers_
	parseImplements_
	parseDirectives_
	parseInputValue_
	parseType_
	parseDefaultVal_
	parseInputValue__
)

type (
	parseFn func(op string) ast.GQLTypeProvider

	Parser struct {
		l *lexer.Lexer

		extend bool

		cache *Cache_

		logr *log.Logger
		logf *os.File

		abort     bool
		stmtType  string
		state     stateT
		curToken  *token.Token
		peekToken *token.Token

		parseFns map[token.TokenType]parseFn
		perror   []error
	}
)

func (p *Parser) setState(o stateT) func() {
	var oldState = o
	return func() { p.state = oldState }
}

var (
	//	enumRepo      ast.EnumRepo

	directiveLocation map[string]ast.DirectiveLoc = map[string]ast.DirectiveLoc{
		//	Executable DirectiveLocation
		"QUERY":               ast.QUERY_DL,
		"MUTATION":            ast.MUTATION_DL,
		"SUBSCRIPTION":        ast.SUBSCRIPTION_DL,
		"FIELD":               ast.FIELD_DL,
		"FRAGMENT_DEFINITION": ast.FRAGMENT_DEFINITION_DL,
		"FRAGMENT_SPREAD":     ast.FRAGMENT_SPREAD_DL,
		"INLINE_FRAGMENT":     ast.INLINE_FRAGMENT_DL,
		//	TypeSystem DirectiveLocation
		"SCHEMA":                 ast.SCHEMA_DL,
		"SCALAR":                 ast.SCALAR_DL,
		"OBJECT":                 ast.OBJECT_DL,
		"FIELD_DEFINITION":       ast.FIELD_DEFINITION_DL,
		"ARGUMENT_DEFINITION":    ast.ARGUMENT_DEFINITION_DL,
		"INTERFACE":              ast.INTERFACE_DL,
		"UNION":                  ast.UNION_DL,
		"ENUM":                   ast.ENUM_DL,
		"ENUM_VALUE":             ast.ENUM_VALUE_DL,
		"INPUT_OBJECT":           ast.INPUT_OBJECT_DL,
		"INPUT_FIELD_DEFINITION": ast.INPUT_FIELD_DEFINITION_DL,
	}
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	// assigns "shared" cache across multiple parsers or parser uses goroutines for processing. None of which is currently employeed.
	//  sharing cache was an exercise in making the cache concurrency safe rather than an actual design necessarity.
	p.cache = NewCache()

	p.parseFns = make(map[token.TokenType]parseFn)
	p.registerFn(token.TYPE, p.ParseObjectType)
	p.registerFn(token.ENUM, p.ParseEnumType)
	p.registerFn(token.INTERFACE, p.ParseInterfaceType)
	p.registerFn(token.UNION, p.ParseUnionType)
	p.registerFn(token.INPUT, p.ParseInputValueType)
	p.registerFn(token.SCALAR, p.ParseScalar)
	p.registerFn(token.DIRECTIVE, p.ParseDirective)
	p.registerFn(token.SCHEMA, p.ParseSchema)
	// Read two tokens, to initialise curToken and peekToken
	p.nextToken()
	p.nextToken()
	//
	// remove cacheClar before releasing..
	//
	//ast.CacheClear()
	return p
}

// astsitory of all types defined in the graph

func (p *Parser) Getperror() []error {
	return p.perror
}

func (p *Parser) Loc() *ast.Loc_ {
	loc := p.curToken.Loc
	return &ast.Loc_{loc.Line, loc.Col}
}

func (p *Parser) ClearCache() {
	p.cache.CacheClear()
}
func (p *Parser) printToken(s ...string) {
	if len(s) > 0 {
		fmt.Printf("** Current Token: [%s] %s %s %s %v %s %s [%s]\n", s[0], p.curToken.Type, p.curToken.Literal, p.curToken.Cat, p.curToken.IsScalarType, "Next Token:  ", p.peekToken.Type, p.peekToken.Literal)
	} else {
		fmt.Println("** Current Token: ", p.curToken.Type, p.curToken.Literal, p.curToken.Cat, "Next Token:  ", p.peekToken.Type, p.peekToken.Literal)
	}
}

// containersErr accepts a "validation" method value and the number of errors that can be generated in its call to abort the process,
// preventing next validation task from running.
func (p *Parser) executeWithErrLimit(mv func(*[]error), num ...int) (b bool) {
	var delta int
	if len(num) == 0 {
		delta = 1
	} else {
		delta = num[0]
	}
	defer func() func() {
		var errBefore = len(p.perror)
		return func() { b = len(p.perror) > errBefore+delta }
	}()()
	// execute method value, mv
	mv(&p.perror)
	return
}

func (p *Parser) hasError() bool {

	if len(p.perror) > 10 || p.abort {
		return true
	}
	return false
}

// addErr appends to error slice held in parser.
func (p *Parser) addErr(s string, xCode ...int) error {

	if strings.Index(s, " at line: ") == -1 {
		s += fmt.Sprintf(" at line: %d, column: %d", p.curToken.Loc.Line, p.curToken.Loc.Col)
	}
	e := errors.New(s)
	p.perror = append(p.perror, e)
	if len(xCode) > 0 {
		p.abort = true
	}
	return e
}

// addErr2 appends to error slice held in parser.
func (p *Parser) addErr2(e error) error {

	fmt.Println("addErr2: before ", len(p.perror))
	p.perror = append(p.perror, e)
	fmt.Println("addErr2: after ", len(p.perror))
	return e
}

func (p *Parser) registerFn(tokenType token.TokenType, fn parseFn) {
	p.parseFns[tokenType] = fn
}

func (p *Parser) nextToken(s ...string) {
	p.curToken = p.peekToken

	p.peekToken = p.l.NextToken() // get another token from lexer:    [,+,(,99,Identifier,keyword,EOF etc.
	//fmt.Println("nextToken: ", p.peekToken.Type, p.peekToken.Literal)
	if len(s) > 0 {
		fmt.Printf("** Current Token: [%s] %s %s %s %s %s %s\n", s[0], p.curToken.Type, p.curToken.Literal, p.curToken.Cat, "Next Token:  ", p.peekToken.Type, p.peekToken.Literal)
	}
	if p.curToken != nil {
		if p.curToken.Illegal {
			p.addErr(fmt.Sprintf("Illegal %s token, [%s]", p.curToken.Type, p.curToken.Literal))
		}
		// if $variable present then mark the identier as a VALUE
		if p.curToken.Literal == token.DOLLAR {
			p.peekToken.Cat = token.VALUE
		}
	}
}

//
func openLogFile() *os.File {
	logf, err := os.OpenFile("sdlserver.sys.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return logf
}

func (p *Parser) closeLogFile() {
	if err := p.logf.Close(); err != nil {
		log.Fatal(err)
	}
}

// ==================== Start =========================

func (p *Parser) ParseDocument(doc ...string) (api *ast.Document, errs []error) {
	var holderr []error
	api = &ast.Document{}
	api.Statements = []ast.GQLTypeProvider{} // slice is initialised  with no elements - each element represents an interface value of type ast.GQLTypeProvider
	api.StatementsMap = make(map[ast.NameValue_]ast.GQLTypeProvider)
	api.ErrorMap = make(map[ast.NameValue_][]error)
	//
	// create cache
	//
	fmt.Println("^^^^^^^ typeNotExists : ", len(typeNotExists))
	fmt.Println("^^^^^^^ Cache_ :", len(p.cache.Cache))
	defer p.closeLogFile()
	defer func() {
		//
		//p.perror = nil
		p.perror = append(p.perror, holderr...)
		for _, v := range api.StatementsMap { //range api.Statements {
			p.perror = append(p.perror, api.ErrorMap[v.TypeName()]...)
		}
		// persist error free statements to db
		for _, v := range api.StatementsMap { //api.Statements {
			if len(api.ErrorMap[v.TypeName()]) == 0 {
				// TODO - what if another type by that name exists
				//  auto overrite or raise an error
				if err := db.Persist(v.TypeName().String(), v); err != nil {
					p.addErr(err.Error())
				}
			}
		}
		//	ast.CacheClear()
		errs = p.perror
	}()
	//
	//  open log file and set logger
	//
	p.logf = openLogFile()
	p.logr = log.New(p.logf, "SDL:", logrFlags)
	p.cache.SetLogger(p.logr)
	//
	// set document
	//
	db.SetDefaultDoc(defaultDoc)
	if len(doc) == 0 {
		db.SetDocument(defaultDoc)
	} else {
		db.SetDocument(doc[0])
	}
	//
	// parse phase - build AST from GraphQL document
	//
	p.logr.Println("Start server...")
	//
	for p.curToken.Type != token.EOF {
		stmtAST := p.ParseStatement()

		// handle any abort error
		if p.hasError() {
			return api, p.perror
		}
		if stmtAST != nil {
			fmt.Printf("Parsed statement: %s %s   (errors: %d) ", stmtAST.Type(), stmtAST.TypeName(), len(p.perror))
			api.Statements = append(api.Statements, stmtAST)
			name := stmtAST.TypeName()
			api.StatementsMap[name] = stmtAST
			api.ErrorMap[name] = p.perror
			// add all stmts to cache (even errored ones). This prevents db searches for errored stmts.
			p.cache.AddEntry(stmtAST.TypeName(), stmtAST)
			//	if len(p.perror) == 0 {
			//		p.cache.AddEntry(stmtAST.TypeName(), stmtAST) // stmts define GL types
			//	}
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
	// validate phase 1 - resolve ALL types ie. check on the existence of all a types nested abstract type, and those nested types
	//                    until all abstract types (ie. non-scalar types) have been validated.
	//                    This process will also assign the AST to *GQLtype.AST  where applicable.
	//					  if cache returns no value then don't generate error as this was done at cache populate time for that item.
	//
	for _, v := range api.Statements {
		fmt.Println("A out to resolve types for ", v.TypeName())
		p.ResolveNestedTypes(v, p.cache)
		if len(p.perror) > 0 {
			api.ErrorMap[v.TypeName()] = append(api.ErrorMap[v.TypeName()], p.perror...)
			p.perror = nil
		}
	}
	//
	//  This is a fix because I incorrectly moved the type cache from the ast package to the parser, towards the end of development.
	//  Fundamentally ast needs this data as does the parser, but because of cyclic dependency ast cannot access the cache when its in the parser.
	//  TODO: put the cache back in package ast. The parser can always access the cache when its in ast.
	//
	//	initialise ast Cache
	ast.InitCache(len(p.cache.Cache))
	for k, v := range p.cache.Cache {
		ast.TyCache[k] = v.data
	}
	fmt.Println("*** entries transfered to ast cache - ", len(ast.TyCache))
	//
	// Build perror from statement errors to use in hasError() counting
	//
	p.perror = holderr
	for _, v := range api.StatementsMap { //api.Statements {
		p.perror = append(p.perror, api.ErrorMap[v.TypeName()]...)
	}
	if p.hasError() {
		return api, p.perror
	}
	//
	// validate phase 2 - generic validations
	//
	p.perror = nil
	for _, v := range api.StatementsMap {
		v.CheckDirectiveLocation(&p.perror)
		if len(p.perror) > 0 {
			api.ErrorMap[v.TypeName()] = append(api.ErrorMap[v.TypeName()], p.perror...)
			p.perror = nil
		}
	}
	//
	// Build perror from statement errors to use in hasError() counting
	//
	p.perror = holderr
	for _, v := range api.StatementsMap { //api.Statements {
		p.perror = append(p.perror, api.ErrorMap[v.TypeName()]...)
	}
	if p.hasError() {
		p.perror = nil
		return api, p.perror
	}
	//
	// validate phase 3 - type specific validations
	//
	errCollect := func(typeName ast.NameValue_) {
		api.ErrorMap[typeName] = append(api.ErrorMap[typeName], p.perror...)
		p.perror = nil
	}

	p.perror = nil
	fmt.Println("stmt map: ", api.StatementsMap)
	for _, v := range api.StatementsMap {
		if p.hasError() {
			break
		}
		//
		// *** proceed for stmt if there are no type resolve errors.
		//
		if len(api.ErrorMap[v.TypeName()]) != 0 {
			var abortMoreValidation bool
			for _, v2 := range api.ErrorMap[v.TypeName()] {
				if errors.Is(v2, TypeResolveErr) {
					abortMoreValidation = true
					continue
				}
			}
			if abortMoreValidation {
				errCollect(v.TypeName())
				continue
			}
		}
		if !p.checkFieldASTAssigned(v) {
			continue
		}
		if p.executeWithErrLimit(v.CheckInputValueType, 5) {
			errCollect(v.TypeName())
			continue
		}
		switch x := v.(type) {
		case *ast.Input_:
			x.CheckIsInputType(&p.perror)
		case *ast.Object_:
			if p.executeWithErrLimit(x.CheckIsOutputType, 5) {
				errCollect(v.TypeName())
				continue
			}
			if p.executeWithErrLimit(x.CheckIsInputType, 5) {
				errCollect(v.TypeName())

				continue
			}
			x.CheckImplements(&p.perror) // check implements are interfaces
		case *ast.Enum_:
		case *ast.Interface_:
		case *ast.Union_:
			p.CheckUnionMembers(x)
		case *ast.Directive_:
			if p.executeWithErrLimit(x.CheckIsInputType, 5) {
				errCollect(v.TypeName())
				continue
			}
			p.CheckSelfReference(v.TypeName(), x)
		}
		//
		errCollect(v.TypeName())
	}
	return api, p.perror
}

// for _, v := range api.ErrorMap {
// 	for _, x := range v {
// 		fmt.Println("xErr: ", x.Error())
// 	}
// }

// ================ parseStatement ==========================

// parseStatement takes predefined parser routine and applies it to a valid statement
func (p *Parser) ParseStatement() ast.GQLTypeProvider {
	p.skipComment()
	if p.curToken.Type == token.EXTEND {
		p.extend = true
		p.nextToken() // read over extend
	}
	p.stmtType = strings.ToLower(p.curToken.Literal)
	if f, ok := p.parseFns[p.curToken.Type]; ok {
		return f(p.stmtType)
	} else {
		p.abort = true
		p.addErr(fmt.Sprintf(`Parse aborted. "%s" is not a statement keyword`, p.stmtType))
	}

	return nil
}

// TypeResolveErr used only to categorise the error not to provided extra information.
var TypeResolveErr = errors.New("")

// ResolveNestedTypes is a validation check performed after parsing completes.
// all nested abstract types for the passed in AST are confirmed to exist either
// in cache or database.  Resolving continutes in FetchAST, until all nested types are resolved
// within each type.
func (p *Parser) ResolveNestedTypes(v ast.GQLTypeProvider, t *Cache_) {
	//
	// find all Abstract Types (ie. non-scalar) nested within the current type v. These need to be resolved.
	// Note: SolicitAbstractTypes does not recursively evaluate all nested types, only the
	// types immediately associated with v, the type under investigation.
	// ResolveNestedType is called via FetchAST, to walk the graph of nested types beyond v.
	//
	fmt.Println("************** ResolveType: ", v.TypeName())
	nestedAbstractTypes := make(ast.UnresolvedMap)
	v.SolicitAbstractTypes(nestedAbstractTypes)
	//
	// load all resolved types from the cache into a map
	//
	resolved := make(ast.UnresolvedMap)
	t.Lock()
	//
	// purge current type from map of all types to be resolved.
	//
	for tyName := range nestedAbstractTypes {
		// if all ready cached then add to resolved map
		if _, ok := t.Cache[tyName.String()]; ok {
			resolved[tyName] = nil
			if tyName.Name == v.TypeName() {
				// remove type that is under consideration from list of types to be resolved.
				delete(nestedAbstractTypes, tyName)
			}
		}
	}
	t.Unlock()
	fmt.Println("AbstractType & resolved: ", nestedAbstractTypes, resolved)
	//
	//  nestedAbstractTypes should now contain abstract types except current type under investigation.
	//  As a side effect of this proecssing we populate the AST attribute in the GQLtype when the AST exists.
	//
	// typeName, *GQLType
	for tyName, gqltype := range nestedAbstractTypes {
		//
		// resolve type - note FetchAST will recursively call ResolveNestedTypes to evalute tyName.
		//
		ast_, err := t.FetchAST(tyName.Name)
		if err != nil {
			switch {
			case errors.Is(err, ErrNotCached):
				p.addErr2(fmt.Errorf(`Item %q %s in document %q %s %w`, tyName, err, db.GetDocument(), tyName.AtPosition(), TypeResolveErr))
			case errors.Is(err, db.NoItemFoundErr):
				p.addErr2(fmt.Errorf(`%s %s %w`, err, tyName.AtPosition(), TypeResolveErr))
			default:
				p.addErr2(fmt.Errorf(`%s %s %w`, err, tyName.AtPosition(), TypeResolveErr))
			}
		} else {
			//
			// we have reached the leaf nodes when gqltype
			//
			if gqltype != nil {
				gqltype.AST = ast_
			}
		}
	}
}

//  ===================== CheckDirectives ================

func (p *Parser) CheckSelfReference(directive ast.NameValue_, x *ast.Directive_) {
	// search for directives in the current stmt making sure it doesn't find itself

	refCheck := func(dirName ast.NameValue_, x ast.GQLTypeProvider) {
		x.CheckDirectiveRef(dirName, &p.perror)
	}

	for _, v := range x.ArgumentDefs {
		for _, dir := range v.Directives {
			if directive.String() == dir.Name_.String() {
				p.addErr(fmt.Sprintf(`Directive "%s" that references itself, is not permitted %s`, directive, dir.Name_.AtPosition()))
			}
		}
		if v.Type.AST != nil {
			refCheck(directive, v.Type.AST)
		}
	}
}

// ===================== CheckUnionMembers ================

func (p *Parser) CheckUnionMembers(x *ast.Union_) {
	//
	for _, m := range x.NameS {
		ast_, err := p.cache.FetchAST(m.Name)
		if err != nil { //ast_ == nil || err != nil {
			if errors.Is(err, ErrNotCached) {
				p.addErr(fmt.Sprintf(`%s. Union member "%s" does not exist %s`, err, m, m.AtPosition()))
			} else {
				p.addErr(fmt.Sprintf(`Union member %s %s %s`, m, err, m.AtPosition()))
			}
		} else {
			switch ast_.(type) {
			case *ast.Object_, *ast.Union_, *ast.Interface_, *ast.Scalar_: //, *ast.Int_, *ast.Float_, *ast.String_, *ast.Boolean_, *ast.ID_:
			default:
				if x, ok := ast_.(ast.InputValueProvider); ok {
					switch x.(type) {
					case *ast.Int_, *ast.Float_, *ast.String_, *ast.Bool_, *ast.ID_:
					default:
						p.addErr(fmt.Sprintf(`Union member "%s" must be an object based type %s`, m, m.AtPosition()))
					}
				} else {
					p.addErr(fmt.Sprintf(`Union member "%s" must be an object based type %s`, m, m.AtPosition()))
				}
			}

		}
	}
}

var opt bool = true // is optional

// ==================== Schema  ============================

//SchemaDefinition
//		schemaDirectivesConstopt{RootOperationTypeDefinitionlist}
// RootOperationTypeDefinition
//  	 OperationType : NamedType
func (p *Parser) ParseSchema(op string) ast.GQLTypeProvider {
	defer p.setState(p.state)()
	// Types: query, mutation, subscription
	p.state = parseSchema
	p.nextToken() // read over type
	if !p.extend {
		inp := &ast.Schema_{}

		p.parseDirectives(inp, opt).parseOperationTypes(inp)

		return inp
	} else {
		// return original AST associated with the extend Name.
		obj, name, err := p.parseExtendName()
		if obj != nil {
			if inp, ok := obj.(*ast.Schema_); !ok {
				p.addErr(fmt.Sprintf(`specified extend type "%s" is not an Input Value Type`, obj.TypeName()))
				p.abort = true
			} else {
				d := inp.Directives_.Len()

				p.parseDirectives(inp, opt)

				if inp.Directives_.Len() > d {
					p.parseOperationTypes(inp, opt)
				} else {
					p.parseOperationTypes(inp)
				}
				return inp
			}
		} else {
			p.addErr(strings.Replace(err.Error(), "Item", "Schema", 1) + p.Loc().String())
			return &ast.Object_{Name_: name}
		}
	}
	return nil
}

func (p *Parser) parseOperationTypes(v *ast.Schema_, optional ...bool) *Parser {
	defer p.setState(p.state)()

	p.state = parseOperationTypes_
	if p.hasError() {
		return p
	}
	if p.curToken.Type != token.LBRACE {
		if len(optional) != 0 {
			p.addErr(fmt.Sprintf("Expected a ( instead got %s", p.curToken.Type))
		}
		return p
	}
	//
	for p.nextToken(); p.curToken.Type != token.RBRACE; {

		p.parseOperation(v).parseColon().parseName(v)

		if p.hasError() {
			return p
		}
	}

	p.nextToken() // read over }
	return p
}

func (p *Parser) parseOperation(inp *ast.Schema_) *Parser {

	switch p.curToken.Type {
	case token.QUERY:
		inp.Op = ast.QUERY_OP
	case token.MUTATION:
		inp.Op = ast.MUTATION_OP
	case token.SUBSCRIPTION:
		inp.Op = ast.SUBSCRIPTION_OP
	default:
		p.addErr(fmt.Sprintf("%s is not a valid operation. Must be query, mutation or subscription ", p.curToken.Type))
	}

	p.nextToken() // read over op type

	return p
}

// checkFieldASTAssigned return false if AST is not assigned. Further validations should not be carried out if AST is not assigned
func (p *Parser) checkFieldASTAssigned(stmt ast.GQLTypeProvider) bool {

	if x, ok := stmt.(ast.SelectionGetter); ok {

		for _, fld := range x.GetSelectionSet() {
			//
			// Confirm argument value type against type definition
			//
			if !fld.Type.IsScalar() && fld.Type.AST == nil {
				//var err error
				fld.Type.AST, _ = p.cache.FetchAST(fld.Type.Name)
				if fld.Type.AST == nil {
					return false
				}
			}
		}
	}
	return true
}

// ==================== Object Type  ============================

// Description-opt TYPE Name  ImplementsInterfaces-opt Directives-Const-opt  FieldsDefinition-opt
//		{FieldDefinition-list}
//         FieldDefinition:
//			Description-opt Name ArgumentsDefinition- opt : Type Directives-Con
func (p *Parser) ParseObjectType(op string) ast.GQLTypeProvider {
	// Types: query, mutation, subscription
	defer p.setState(p.state)()

	p.state = parseObjectType
	p.nextToken() // read over type
	fmt.Println("parseObjectType..........")
	if !p.extend {
		obj := &ast.Object_{}

		p.parseName(obj).parseImplements(obj, opt).parseDirectives(obj, opt).parseFields(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj, name, err := p.parseExtendName()
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
			p.addErr(strings.Replace(err.Error(), "Item", "Type", 1) + p.Loc().String())
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
	defer p.setState(p.state)()

	p.state = parseEnumType
	p.nextToken() // read type
	if !p.extend {
		obj := &ast.Enum_{}

		p.parseName(obj).parseDirectives(obj, opt).parseEnumValues(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj, name, err := p.parseExtendName()
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
			p.addErr(strings.Replace(err.Error(), "Item", "Enum", 1) + p.Loc().String())
			return &ast.Enum_{Name_: name}
		}
	}
	return nil
}

// ====================== Interface ===========================
// InterfaceTypeDefinition
//		Description-opt	interface	Name	Directives-opt	FieldsDefinition-opt
func (p *Parser) ParseInterfaceType(op string) ast.GQLTypeProvider {
	defer p.setState(p.state)()

	p.state = parseInterfaceType
	p.nextToken() // read over interfcae keyword
	if !p.extend {
		obj := &ast.Interface_{}

		p.parseName(obj).parseDirectives(obj, opt).parseFields(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj, name, err := p.parseExtendName()
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
			p.addErr(strings.Replace(err.Error(), "Item", "Interface", 1) + p.Loc().String())
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
	defer p.setState(p.state)()

	p.state = parseUnionType
	p.nextToken() // read over interfcae keyword
	if !p.extend {
		obj := &ast.Union_{}

		p.parseName(obj).parseDirectives(obj, opt).parseUnionMembers(obj, opt)

		return obj
	} else {
		// return original AST associated with the extend Name.
		obj, name, err := p.parseExtendName()
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
			p.addErr(strings.Replace(err.Error(), "Item", "Union", 1) + p.Loc().String())
			return &ast.Union_{Name_: name}
		}
	}
	return nil
}

//====================== Input ===============================
// InputObjectTypeDefinition
//		Description-opt	input	Name	DirectivesConst-opt	InputFieldsDefinition-opt
func (p *Parser) ParseInputValueType(op string) ast.GQLTypeProvider {
	defer p.setState(p.state)()

	p.state = parseInputValueType
	p.nextToken() // read over input keyword
	if !p.extend {
		inp := &ast.Input_{}

		p.parseName(inp).parseDirectives(inp, opt).parseInputFieldDefs(inp)

		return inp
	} else {
		// return original AST associated with the extend Name.
		obj, name, err := p.parseExtendName()
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
			p.addErr(strings.Replace(err.Error(), "Item", "Input type", 1) + p.Loc().String())
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
	defer p.setState(p.state)()

	p.state = parseScalar
	p.nextToken() // read over input keyword
	if !p.extend {
		inp := &ast.Scalar_{}

		p.parseName(inp).parseDirectives(inp, opt)

		return inp
	} else {
		// return original AST associated with the extend Name.
		obj, name, err := p.parseExtendName()
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
			p.addErr(strings.Replace(err.Error(), "Item", "Scalar", 1) + p.Loc().String())
			return &ast.Scalar_{Name: name.String()}
		}
	}
	return nil

}

// ====================== Directive_ ===============================
// DirectiveDefinition
//	Descriptiono-pt directive @ Name ArgumentsDefinition-opt  on  DirectiveLocations
// DirectiveLocations
//   | optDirectiveLocation
// DirectiveLocations | DirectiveLocation
// DirectiveLocation
//      ExecutableDirectiveLocation
//      TypeSystemDirectiveLocation
func (p *Parser) ParseDirective(op string) ast.GQLTypeProvider {
	defer p.setState(p.state)()

	p.state = parseDirective
	p.nextToken() // read over input keyword

	inp := &ast.Directive_{}

	p.parseAtSign().parseName(inp).parseFieldArgumentDefs(inp).parseOn().parseDirectiveLocations(inp)

	if p.hasError() {
		return nil
	}

	inp.CoerceDirectiveName() // prepend @ to name

	return inp
}

// =============================================================

func (p *Parser) skipComment() {
	if p.curToken.Type == token.STRING {
		p.nextToken() // read over comment string
	}
}

func (p *Parser) parseAtSign() *Parser {
	if p.curToken.Type == token.ATSIGN {
		p.nextToken() // read over @
	} else {
		p.addErr(fmt.Sprintf("Expected @ got %s of %s ", p.curToken.Type, p.curToken.Literal))
		p.abort = true
	}
	return p
}

func (p *Parser) parseOn() *Parser {
	if p.curToken.Type == token.ON {
		p.nextToken() // read over ON
	} else {
		p.addErr(fmt.Sprintf("Expected @ got %s of %s ", p.curToken.Type, p.curToken.Literal))
		p.abort = true
	}
	return p
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
		if p.curToken.Type == token.ILLEGAL {
			p.abort = true
		}
	}
	p.nextToken() // read over name
	return p
}

// ==================== parseExtendName ===============================
// parseExtendName will consume the type name to be extended. Returns the type's AST.
func (p *Parser) parseExtendName() (ast.GQLTypeProvider, ast.Name_, error) {
	// if p.hasError() {
	// 	return nil,
	// }
	var extName string
	// does name entity exist
	if p.curToken.Type == token.IDENT {
		extName = p.curToken.Literal
	} else {
		p.addErr(fmt.Sprintf(`Expected name identifer got %s of "%s"`, p.curToken.Type, p.curToken.Literal), FATAL)
		if p.curToken.Type == token.ILLEGAL {
			p.abort = true
			return nil, ast.Name_{}, nil
		}
	}
	name_ := ast.Name_{Name: ast.NameValue_(extName), Loc: p.Loc()}
	ast, err := p.cache.FetchAST(name_.Name) // ignore error as as ast value of nil means no data found
	// handle err to calling routine, which can add extra value
	if ast != nil {
		p.nextToken() // read over name
	}
	return ast, name_, err
}

func (p *Parser) parseEnumValues(enum *ast.Enum_, optional ...bool) *Parser {
	defer p.setState(p.state)()

	p.state = parseEnumValues_

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

//========================= parseDirectiveLocations ====================================

func (p *Parser) parseDirectiveLocations(d *ast.Directive_) *Parser {
	defer p.setState(p.state)()

	p.state = parseDirectiveLocations_
	if p.hasError() {
		return p
	}

	for ; p.curToken.Type == token.BAR || p.curToken.Type == token.IDENT; p.nextToken() {
		if p.curToken.Type == token.BAR && p.peekToken.Type != token.IDENT {
			p.addErr(fmt.Sprintf("expected directive location identifer, got %s, %s", p.curToken.Type, p.curToken.Literal))
		} else if p.curToken.Type == token.IDENT {
			if dl, ok := directiveLocation[p.curToken.Literal]; !ok {
				p.addErr(fmt.Sprintf(`Invalid directive location "%s" %s`, p.curToken.Literal, p.Loc()))
			} else {
				d.Location = append(d.Location, dl)
			}
		}
	}
	return p
}

//========================= parseUnionMembers ====================================
// UnionMemberTypes
//		=|optNamedType
//		UnionMemberTypes | NamedType
func (p *Parser) parseUnionMembers(u *ast.Union_, optional ...bool) *Parser {
	defer p.setState(p.state)()

	p.state = parseUnionMembers_
	if p.curToken.Type != token.ASSIGN {
		p.abort = true
	}
	if p.hasError() {
		return p
	}
	for p.nextToken(); p.curToken.Type == token.BAR || p.curToken.Type == token.IDENT; p.nextToken() {
		if p.curToken.Type == token.BAR && p.peekToken.Type != token.IDENT {
			p.addErr(fmt.Sprintf("expected Union  member  identifer, got %s, %s", p.curToken.Type, p.curToken.Literal))
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
	defer p.setState(p.state)()

	p.state = parseImplements_

	if p.curToken.Type != token.IMPLEMENTS {
		if len(optional) == 0 {
			p.abort = true
		}
		return p
	}
	if p.hasError() {
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
	defer p.setState(p.state)()

	p.state = parseDirectives_
	fmt.Println("ParseDirective.........", p.abort)
	if p.hasError() {
		return p
	}
	if p.curToken.Type != token.ATSIGN {
		if len(optional) == 0 {
			p.addErr("Variable is mandatory")
		}
		fmt.Println("parseDirective return...", p.curToken.Type)
		return p
	}

	for p.curToken.Type == token.ATSIGN {
		p.nextToken() // read over @
		a := []*ast.ArgumentT{}
		d := &ast.DirectiveT{Arguments_: ast.Arguments_{Arguments: a}}

		p.parseName(d).parseArguments(d, opt)

		if p.hasError() {
			return p
		}

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
	defer p.setState(p.state)()

	p.state = parseArguments_
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
	for p.nextToken(); p.curToken.Type != token.RPAREN; { // p.nextToken() is now redundant as parseInputValue handles nexttoken()
		v := new(ast.ArgumentT)

		p.parseName(v).parseColon().parseInputValue(v)

		if p.hasError() {
			return p
		}

		f.AppendArgument(v)
	}
	p.nextToken() // read over )
	return p
}

// ParseReponse - API for testing purposes only. Not part of normal processing.
func (p *Parser) ParseResponse() ast.InputValueProvider {

	value := p.parseInputValue_()

	return value.InputValueProvider

}

func (p *Parser) parseColon() *Parser {

	if p.hasError() {
		return p
	}
	if p.curToken.Type != token.COLON {
		var next string
		switch p.state {
		case parseOperationTypes_:
			next = "an operation name (Query,Mutation) "
		case parseFields_, parseInputFieldDefs_, parseFieldArgumentDefs_:
			next = "a GQL-Type"
		case parseObjectArguments_:
			next = "an argument value"
		case parseDirectives_:
			next = "a directive argument value"
		case parseArguments_:
			next = "an argument value"
		default:
			next = strconv.Itoa(int(p.state))
		}
		p.addErr(fmt.Sprintf(`Expected a colon followed by %s, got "%s" `, next, p.curToken.Literal), FATAL)
		return p
	}
	p.nextToken() // read over :
	return p
}

func (p *Parser) parseInputValue(v *ast.ArgumentT) *Parser {
	defer p.setState(p.state)()

	p.state = parseInputValue_
	if p.hasError() {
		return p
	}
	if !((p.curToken.Cat == token.VALUE && (p.curToken.Type == token.DOLLAR && p.peekToken.Cat == token.VALUE)) ||
		(p.curToken.Cat == token.VALUE && (p.peekToken.Cat == token.NONVALUE || p.peekToken.Type == token.RPAREN)) ||
		(p.curToken.Type == token.LBRACKET || p.curToken.Type == token.LBRACE)) { // [  or {
		p.addErr(fmt.Sprintf(`Expected an argument value followed by an identifer or close parenthesis got "%s"`, p.curToken.Literal))
	}
	v.Value = p.parseInputValue_()

	return p
}

//======================== parseFields ==========================================
// {FieldDefinition ...} :
// .  Description-opt Name ArgumentsDefinition-opt : Type Directives-opt
func (p *Parser) parseFields(f ast.FieldAppender, optional ...bool) *Parser {
	defer p.setState(p.state)()

	p.state = parseFields_
	if p.hasError() || p.curToken.Type != token.LBRACE {
		fmt.Println("here..", len(optional))
		if len(optional) == 0 {
			p.addErr("Field definitions is required")
		}
		return p
	}

	for p.nextToken(); p.curToken.Type != token.RBRACE; { // p.nextToken("next token in parseFields..") {

		field := &ast.Field_{}

		_ = p.parseDecription().parseName(field).parseFieldArgumentDefs(field).parseColon().parseType(field).parseDirectives(field, opt)

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

// func  (p *Parser) ParseType(p ParserI, f ast.AssignTyper) {
// 	// need to create methods for all parser variables
// 	return p.parseType(f)
// }

func (p *Parser) parseType(f ast.AssignTyper) *Parser {
	defer p.setState(p.state)()

	p.state = parseType_
	if p.hasError() {
		return p
	}
	if p.curToken.Type == token.COLON {
		p.addErr("A second colon detected")
		p.nextToken() // read over :
	}
	// else {
	// 	p.addErr(fmt.Sprintf("Colon expected got %s of %s", p.curToken.Type, p.curToken.Literal))
	// }

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
		//typedef ast.TypeFlag_ // token defines SCALAR types only. All other types will be populated in astType map.
		depth   uint8
		nameLoc *ast.Loc_
	)
	nameLoc = p.Loc()
	switch p.curToken.Type {

	case token.LBRACKET:
		// [ typeName ]
		var (
			depthClose uint8
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
		// //System ScalarTypes are defined by the GQLtype.Name_, Non-system Scalar and non-scalar are defined by the AST.
		// if !p.curToken.IsScalarType {
		// 	ast_ = p.fetchAST(name_)
		// }
		p.nextToken() // read over IDENT
		for bangs := uint8(0); p.curToken.Type == token.RBRACKET || p.curToken.Type == token.BANG; {
			if p.curToken.Type == token.BANG {
				bangs++
				fmt.Println("parseType: bangs:", bangs)
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
		if depth != depthClose {
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
	t := &ast.GQLtype{Constraint: bit, Depth: depth} //, AST: ast_}
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
// 	Type       GQLtype
// 	DefValue   *InputValueDef_
// 	Directives []*Directive_
// }
// InputFieldsDefinition
//		{InputValueDefinition-list}
// InputValueDefinition
//		Description-opt	Name	:	Type	DefaultValue-opt	Directives-opt
func (p *Parser) parseFieldArgumentDefs(f ast.FieldArgAppender) *Parser { // st is an iv initialised from passed in argument which is a *OperationStmt
	defer p.setState(p.state)()

	encl := [2]token.TokenType{token.LPAREN, token.RPAREN} // ()
	p.state = parseFieldArgumentDefs_
	fmt.Println("parseFieldArgumentDefs......")
	if p.hasError() {
		return p
	}

	return p.parseArgumentDefs(f, encl)
}

func (p *Parser) parseInputFieldDefs(f ast.FieldArgAppender) *Parser {
	defer p.setState(p.state)()
	encl := [2]token.TokenType{token.LBRACE, token.RBRACE} // {}
	p.state = parseInputFieldDefs_
	if p.hasError() {
		return p
	}
	return p.parseArgumentDefs(f, encl)
}

func (p *Parser) parseArgumentDefs(f ast.FieldArgAppender, encl [2]token.TokenType) *Parser {

	if p.curToken.Type == encl[0] {
		p.nextToken() // read over ( or {
		for p.curToken.Type != encl[1] {
			//for p.curToken.Type != ":" { //TODP fix should be encl[1]
			v := &ast.InputValueDef{}
			v.Loc = p.Loc()
			//	p.parseDecription().parseName(v).parseType(v).parseDefaultVal(v, opt).parseDirectives(v, opt)
			p.parseDecription().parseName(v).parseColon().parseType(v).parseDefaultVal(v, opt).parseDirectives(v, opt)
			if p.hasError() {
				return p
			}
			f.AppendField(v, &p.perror)
		}
		p.nextToken() //read over )
	}

	return p
}

// type InputValueDef struct {
// 	Desc string
// 	Name_
// 	Type       *GQLtype          // ENUM
// 	DefaultVal *InputValue_    // ENUMVALUE
// 	Directives_
// }

func (p *Parser) parseDefaultVal(v *ast.InputValueDef, optional ...bool) *Parser {
	defer p.setState(p.state)()

	p.state = parseDefaultVal_
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
		//	p.nextToken() // redundant now that parseInputValue_ performs it
		// if ev, ok := v.DefaultVal.InputValueProvider.(*ast.EnumValue_); ok {
		// 	enum := p.Cache.FetchAST(v.Type.Name_)
		// 	// check def value is member of ENUM
		// 	var found bool
		// 	for _, e := range enum.Values {
		// 		if ev.Name_.Equals(e.Name_) {
		// 			found = true
		// 			break
		// 		}
		// 	}
		// 	if !found {
		// 		p.addErr(`Default value "%s", not a member of ENUM %s %s`, ev.Name, enum.Name_, ev.Name.AtPosition())
		// 	}
		// }
	}

	fmt.Printf("Default Val: %#v\n", v.DefaultVal)
	return p
}

// parseObjectArguments - used for input object values
func (p *Parser) parseObjectArguments(argS []*ast.ArgumentT) []*ast.ArgumentT {
	defer p.setState(p.state)()

	p.state = parseObjectArguments_
	for p.curToken.Type == token.IDENT {

		v := new(ast.ArgumentT)

		p.parseName(v).parseColon().parseInputValue(v)
		if p.hasError() {
			return nil
		}

		argS = append(argS, v)

	}
	return argS
}

// parseInputValue_ used to interpret "default value" in argument and field values.
//  parseInputValue_ expects an InputValue_ literal (true,false, 234, 23.22, "abc" or $variable in the next token.  The value is a type bool,int,flaot,string..
//  if it is a variable then the variable value (which is an InputValue_ type) will be sourced
//  TODO: currently called from parseArgument only. If this continues to be the case then add this func as anonymous func to it.
//func (p *Parser) parseInputValue_(iv ...*ast.InputValueDef) *ast.InputValue_ { //TODO remove iv argeument now redundant
func (p *Parser) parseInputValue_() *ast.InputValue_ {
	defer p.nextToken() // this func will finish paused on next token - always
	defer p.setState(p.state)()

	p.state = parseInputValue__
	if p.hasError() {
		return nil
	}
	fmt.Println("parseInputValue_............................", p.curToken.Type, p.curToken.Literal)
	if p.curToken.Type == "ILLEGAL" {
		p.addErr(fmt.Sprintf("Value expected got %s of %s", p.curToken.Type, p.curToken.Literal))
		p.abort = true
		return nil
	}
	switch p.curToken.Type {
	//
	//  List type
	//
	case token.LBRACKET:
		// [ value value value .. ]
		p.nextToken() // read over [
		// if p.curToken.Cat != token.VALUE { // Comment out - as Token can be an ENUMVALUE or an IDENT both of which have a CAT of NONVALUE - so lets ignore this check as its uncessary anyway
		// 	fmt.Printf("Expect an Input Value followed by another Input Value or a ], got %s %s ", p.curToken.Literal, p.peekToken.Literal)
		// 	p.addErr(fmt.Sprintf("Expect an Input Value followed by another Input Value or a ], got %s %s ", p.curToken.Literal, p.peekToken.Literal))
		// 	return &ast.InputValue_{}
		// }
		// edge case: empty, []
		if p.curToken.Type == token.RBRACKET {
			p.nextToken() // ]
			var null ast.Null_ = true
			iv := ast.InputValue_{InputValueProvider: null, Loc: p.Loc()}
			return &iv
		}
		// process list of values - all value types should be the same
		var vallist ast.List_
		for p.curToken.Type != token.RBRACKET {
			v := p.parseInputValue_()
			vallist = append(vallist, v)
		}
		// completed processing values, return List type
		iv := ast.InputValue_{InputValueProvider: vallist, Loc: p.Loc()}
		return &iv
	//
	//  Object type
	//
	case token.LBRACE:
		//  { name:value name:value ... }
		p.nextToken()              // read over {
		var ObjList ast.ObjectVals // []*ArgumentT {Name_,Value *InputValue_}
		for p.curToken.Type != token.RBRACE {

			ObjList = p.parseObjectArguments(ObjList)
			if p.hasError() {
				return &ast.InputValue_{}
			}
		}
		iv := ast.InputValue_{InputValueProvider: ObjList, Loc: p.Loc()}
		return &iv
	//
	//  Standard Scalar types
	//
	case token.NULL:
		var null ast.Null_ = true
		iv := ast.InputValue_{InputValueProvider: null, Loc: p.Loc()}
		return &iv
	case token.INT:
		i := ast.Int_(p.curToken.Literal)
		iv := ast.InputValue_{InputValueProvider: i, Loc: p.Loc()}
		return &iv
	case token.FLOAT:
		f := ast.Float_(p.curToken.Literal)
		iv := ast.InputValue_{InputValueProvider: f, Loc: p.Loc()}
		return &iv
	case token.STRING:
		f := ast.String_(p.curToken.Literal)
		iv := ast.InputValue_{InputValueProvider: f, Loc: p.Loc()}
		return &iv
	case token.ID:
		id := ast.ID_(p.curToken.Literal)
		iv := ast.InputValue_{InputValueProvider: id, Loc: p.Loc()}
		return &iv
	case token.RAWSTRING:
		f := ast.RawString_(p.curToken.Literal)
		iv := ast.InputValue_{InputValueProvider: f, Loc: p.Loc()}
		return &iv
	case token.TRUE, token.FALSE: //token.BOOLEAN:
		var b ast.Bool_
		if p.curToken.Literal == "true" {
			b = ast.Bool_(true)
		} else {
			b = ast.Bool_(false)
		}
		iv := ast.InputValue_{InputValueProvider: b, Loc: p.Loc()}
		return &iv
	// case token.Time:
	// 	b := ast.Time_(p.curToken.Literal)
	// 	iv := ast.InputValue_{Value: b, Loc: p.Loc()}
	// 	return &iv
	default:
		// possible ENUM value
		b := &ast.EnumValue_{}
		b.AssignName(string(p.curToken.Literal), p.Loc(), &p.perror)
		iv := ast.InputValue_{InputValueProvider: b, Loc: p.Loc()}
		return &iv
	}
	return nil

}
