package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/graph-sdl/ast"
)

const (
	TableName string = "GraphQL2"
)

var (
	document   string
	defaultDoc string

	db *dynamodb.DynamoDB
)

type TypeRow struct {
	PKey  string
	SortK string
	Stmt  string
	Type  string //this maps to ast.Type.Base - reqired for ENUM types but maybe useful for others
}

type PkRow struct {
	PKey  string
	SortK string
}

func init() {
	fmt.Println("***************************************************************** init ast_repo ***********************************************")

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

func buildKey(input string) string {
	var s strings.Builder
	if len(document) == 0 {
		document = defaultDoc
	}
	s.WriteString(input)
	s.WriteString("/")
	s.WriteString(document)
	return s.String()
}

func Persist(input string, ast_ ast.GQLTypeProvider) error {
	// save GraphQL statement to Dynamodb
	if err := dbPersist(input, ast_); err != nil {
		return err
	}
	return nil
}

// func FetchInterface(input string) (*Interface_, bool, string) {
// 	if ast, ok := typeCache_[buildKey(input.Name)]; ok {
// 		if ast_, ok := ast.(*Interface_); !ok {
// 			return nil, false, fmt.Sprintf(`Implements type "%s" is not an Interface %s`, input, input.AtPosition())
// 		} else {
// 			return ast_, true, ""
// 		}
// 	} else {
// 		return nil, true, ""
// 	}

// }

func dbPersist(pkey string, ast_ ast.GQLTypeProvider) error {
	//
	// TODO: check to see if item already exists, and if type is different error otherwise give a warning.
	//		 table design ensures uniqueness of type with a given name, however currently it will overrite existing item
	//
	switch ast_.(type) {

	case *ast.Directive_:
		type DirRow struct {
			PKey  string
			SortK string
			Stmt  string
			Dir   string // Part of Secondary index - identifies Directives only
			Type  string // Type of stmt - saves having to parse stmt to determine type
			PKey_ string // Object belonging to interface
		}
		//	typeDef := DirRow{PKey: pkey.String(), SortK: "D", Stmt: ast.String(), Dir: "D", Type: "D"}
		typeDef := DirRow{PKey: pkey, SortK: document, Stmt: ast_.String(), Dir: "D", Type: "D", PKey_: ast_.String()}
		av, err := dynamodbattribute.MarshalMap(typeDef)
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
		}
		_, err = db.PutItem(&dynamodb.PutItemInput{
			TableName:           aws.String(TableName),
			ConditionExpression: aws.String("attribute_not_exists(Pkey)"),
			Item:                av,
		})
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to PutItem ", err.Error())
		}

	case *ast.Object_:
		//	typeDef := TypeRow{PKey: pkey.String(), SortK: "__", Stmt: ast.String(), Type: "O"}
		typeDef := TypeRow{PKey: pkey, SortK: document, Stmt: ast_.String(), Type: "O"}
		av, err := dynamodbattribute.MarshalMap(typeDef)
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
		}
		// Note: attribute_not_exists(Pkey) - means check for the existence of a tuple with the supplied PKey + SortK
		//  and then check for the existence of the attribute_not_exists attribute ie. PKey - if it exists (meaning an item was found) then return false and prevent insert.
		// In the case of a non-key field e.g. attribute_not_exists(email), the process is to use the suppled pkey + sortk
		// to find a tuple and then check to see if the attribute email exists. If it does exit return false and prevent insert.
		//  so the emphasis is on "find tuple then check to see if attribute exists".
		// Without the condition expression PutItem will simply overwrite any data. You can prevent the default insert operation using condition express.
		//
		_, err = db.PutItem(&dynamodb.PutItemInput{
			TableName:           aws.String(TableName),
			ConditionExpression: aws.String("attribute_not_exists(Pkey)"),
			Item:                av,
		})
		if err != nil {
			return fmt.Errorf("%s: %s", "Error: failed to PutItem ", err.Error())
		}
		// for _, imp := range x.Implements {
		// 	if err := persistImplements(imp.Name, x.TypeName()); err != nil {
		// 		return err
		// 	}
		// }

	default:
		//typeDef := TypeRow{PKey: pkey.String(), SortK: "__", Stmt: ast.String(), Type: isType(ast)}
		typeDef := TypeRow{PKey: pkey, SortK: document, Stmt: ast_.String(), Type: ast.IsGLType(ast_)}
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

// func persistImplements(interface_ NameValue_, object_ NameValue_) error {
// 	//
// 	type ImplementRow struct {
// 		PKey  string
// 		SortK string
// 	}
// 	// Key design permits searching for all objects that reference an interface
// 	// pkey=? and sortk = startWith("?/")
// 	typeDef := ImplementRow{PKey: interface_.String(), SortK: document + "/" + object_.String()}
// 	av, err := dynamodbattribute.MarshalMap(typeDef)
// 	if err != nil {
// 		return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
// 	}
// 	// Note: attribute_not_exists(Pkey) - means check for the existence of a tuple with the supplied PKey + SortK
// 	//  and then check for the existence of the PKey - if PKey exists (meaning an item was found) then return false and prevent insert.
// 	// In the case of a non-key field e.g. attribute_not_exists(email), the process is to use the suppled pkey + sortk
// 	// to find a tuple and then check to see if the attribute email exists. If it does exit return false and prevent insert.
// 	//  so the emphasis is on "find tuple then check to see if attribute exists".
// 	// Without the condition expression PutItem will simply overwrite any data. You can prevent the default insert operation using condition express.
// 	//
// 	_, err = db.PutItem(&dynamodb.PutItemInput{
// 		TableName:           aws.String(TableName),
// 		ConditionExpression: aws.String("attribute_not_exists(Pkey)"),
// 		Item:                av,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("%s: %s", "Error: failed to PutItem ", err.Error())
// 	}
// 	return nil
// }

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

func DBFetch(name string) (string, error) {
	//
	// query on recipe name to get RecipeId and  book name
	//
	///var sortK string
	fmt.Printf("XX DB Fetch name: [%s]\n", name)
	if len(document) == 0 {
		document = defaultDoc
	}
	fmt.Println("DBFetch document : ", document)
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
	pkey := PkRow{PKey: name, SortK: document}
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
		return "", fmt.Errorf(`Type "%s" not found`, name)
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
