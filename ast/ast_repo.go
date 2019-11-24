package ast

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TableName string = "GraphQL2"
)

var (
	document   string
	defaultDoc string
)

type TypeRow struct {
	PKey  string
	SortK string
	Stmt  string
	Type  string //TODO: is this necessary?  Reason: is saves having to parse stmt in order to determine its "type"
}

// cache returns the AST type for a given TypeName
type typeCache map[string]GQLTypeProvider

type PkRow struct {
	PKey  string
	SortK string
}

var typeCache_ typeCache
var db *dynamodb.DynamoDB

func init() {
	typeCache_ = make(typeCache)

	dynamodbService := func() *dynamodb.DynamoDB {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		})
		if err != nil {
			log.Panic(err)
		}
		return dynamodb.New(sess, aws.NewConfig())
	}

	db = dynamodbService()
}

// Fetch - when type is in cache it is said to be "resolved".
//  unresolved types are therefore not in the typeCaches
// func Fetch(input NameValue_) (GQLTypeProvider, bool) {
// 	return CacheFetch(input)
// }

func CacheClear() {
	fmt.Println("******************************************")
	fmt.Println("************ CLEAR CACHE *****************")
	fmt.Println("******************************************")
	typeCache_ = map[string]GQLTypeProvider{} // map literal to zero cache
}

func buildKey(input NameValue_) string {
	var s strings.Builder
	if len(document) == 0 {
		document = defaultDoc
	}
	s.WriteString(input.String())
	s.WriteString("/")
	s.WriteString(document)
	return s.String()
}
func CacheFetch(input NameValue_) (GQLTypeProvider, bool) { // TODO: use GQLTypeProvider instead of GQLTypeProvider?
	fmt.Printf("** CacheFetch [%s]\n", input)

	if ast, ok := typeCache_[buildKey(input)]; !ok {
		return nil, false
	} else {
		return ast, true
	}
}

func Persist(input NameValue_, ast GQLTypeProvider) {
	// save GraphQL statement to Dynamodb
	dbPersist(input, ast)
}

func Add2Cache(input NameValue_, obj GQLTypeProvider) {
	//	fmt.Printf("** Add2Cache  %s [%s]\n", input, obj.String())
	typeCache_[buildKey(input)] = obj
}

func fetchInterface(input Name_) (*Interface_, bool, string) {
	if ast, ok := typeCache_[buildKey(input.Name)]; ok {
		if ast_, ok := ast.(*Interface_); !ok {
			return nil, false, fmt.Sprintf(`Implements type "%s" is not an Interface %s`, input, input.AtPosition())
		} else {
			return ast_, true, ""
		}
	} else {
		return nil, true, ""
	}

}

func dbPersist(pkey NameValue_, ast GQLTypeProvider) error {
	//
	// TODO: check to see if item already exists, and if type is different error otherwise give a warning.
	//		 table design ensures uniqueness of type with a given name, however currently it will overrite existing item
	//
	switch x := ast.(type) {
	case *Directive_:
		type DirRow struct {
			PKey  string
			SortK string
			Stmt  string
			Dir   string // Part of Secondary index - identifies Directives only
			Type  string // Type of stmt - saves having to parse stmt to determine type
			PKey_ string // Object belonging to interface
		}
		//	typeDef := DirRow{PKey: pkey.String(), SortK: "D", Stmt: ast.String(), Dir: "D", Type: "D"}
		typeDef := DirRow{PKey: pkey.String(), SortK: document, Stmt: ast.String(), Dir: "D", Type: "D", PKey_: ast.String()}
		av, err := dynamodbattribute.MarshalMap(typeDef)
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
		}
		_, err = db.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(TableName),
			Item:      av,
		})
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to PutItem ", err.Error())
		}
	case *Object_:
		//	typeDef := TypeRow{PKey: pkey.String(), SortK: "__", Stmt: ast.String(), Type: "O"}
		typeDef := TypeRow{PKey: pkey.String(), SortK: document, Stmt: ast.String(), Type: "O"}
		av, err := dynamodbattribute.MarshalMap(typeDef)
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
		}
		_, err = db.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(TableName),
			Item:      av,
		})
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to PutItem ", err.Error())
		}
		for _, imp := range x.Implements {
			PersistImplements(imp.Name, x.TypeName())
		}
	default:
		//typeDef := TypeRow{PKey: pkey.String(), SortK: "__", Stmt: ast.String(), Type: isType(ast)}
		typeDef := TypeRow{PKey: pkey.String(), SortK: document, Stmt: ast.String(), Type: isType(ast)}
		av, err := dynamodbattribute.MarshalMap(typeDef)
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
		}
		_, err = db.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(TableName),
			Item:      av,
		})
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to PutItem ", err.Error())
		}
	}
	return nil
}
func SetDocument(doc string) {
	document = doc
}

func SetDefaultDoc(doc string) {
	defaultDoc = doc
}

func PersistImplements(pkey NameValue_, sortk NameValue_) error {
	//
	type ImplementRow struct {
		PKey  string
		SortK string
		In    string
	}
	fmt.Println("PersistImplements: ", pkey.String(), sortk.String(), sortk.String())
	//typeDef := ImplementRow{PKey: pkey.String(), SortK: sortk.String(), In: sortk.String()}
	typeDef := ImplementRow{PKey: pkey.String(), SortK: document, In: sortk.String()}
	av, err := dynamodbattribute.MarshalMap(typeDef)
	if err != nil {
		return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
	}
	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("%s: %s", "Error: failed to PutItem ", err.Error())
	}
	return nil
}

func DeleteType(input string) error {

	//
	if len(document) == 0 {
		if len(defaultDoc) == 0 {
			defaultDoc = "DefaultDoc"
		}
		document = defaultDoc
	}
	fmt.Println("delete from document: ", document)
	typeDef := PkRow{PKey: input, SortK: document}
	av, err := dynamodbattribute.MarshalMap(typeDef)
	if err != nil {
		return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
	}
	_, err = db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key:       av,
	})
	if err != nil {
		return fmt.Errorf(`Error: failed to DeleteItem: "%s"  %s`, input, err.Error())
	}
	//TODO - delete any implement items
	return nil
}

func ListCache() []GQLTypeProvider {
	l := make([]GQLTypeProvider, len(typeCache_), len(typeCache_))
	i := 0
	for _, v := range typeCache_ {
		l[i] = v
		i++
	}
	return l
}

func DBFetch(name NameValue_) (string, error) {
	//
	// query on recipe name to get RecipeId and  book name
	//
	///var sortK string
	fmt.Printf("DB Fetch name: [%s]\n", name.String())
	if len(document) == 0 {
		document = defaultDoc
	}
	if len(name) == 0 {
		return "", fmt.Errorf("No DB search value provided")
	}
	errmsg := "Error in marshall of pKey "
	// if name[0] == '@' {
	// 	sortK = "D"
	// } else {
	// 	sortK = "__"
	// }
	//pkey := PkRow{PKey: name.String(), SortK: sortK}
	pkey := PkRow{PKey: name.String(), SortK: document}
	av, err := dynamodbattribute.MarshalMap(&pkey)
	if err != nil {
		return "", fmt.Errorf("%s. MarshalMap: %s", errmsg, err.Error())
	}
	input := &dynamodb.GetItemInput{
		Key:       av,
		TableName: aws.String(TableName),
	}
	input = input.SetReturnConsumedCapacity("TOTAL").SetConsistentRead(true)
	//
	result, err := db.GetItem(input)
	if err != nil {
		fmt.Println("ERROR")
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			//case dynamodb.ErrCodeRequestLimitExceeded:
			//	fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", fmt.Errorf("%s %s: %s", errmsg, "GetItemX", err.Error())
	}
	fmt.Println("dbFetch: GetItem: Query ConsumedCapacity: \n", result.ConsumedCapacity)
	if len(result.Item) == 0 {
		return "", nil //fmt.Errorf(`No type data found in database for "%s"`, name)
	}
	rec := &TypeRow{}
	err = dynamodbattribute.UnmarshalMap(result.Item, rec)
	if err != nil {
		errmsg := "error in unmarshal "
		return "", fmt.Errorf("%s. UnmarshalMaps:  %s", errmsg, err.Error())
	}
	fmt.Printf("DBfetch result: [%s] \n", rec.Stmt)
	return rec.Stmt, nil
}
