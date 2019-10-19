package ast

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// cache returns the AST type for a given TypeName
type typeCache map[NameValue_]TypeDefiner

type TypeRow struct {
	PKey string
	Def  string
}

type PkRow struct {
	PKey string
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
// func Fetch(input NameValue_) (TypeDefiner, bool) {
// 	return CacheFetch(input)
// }

func CacheFetch(input NameValue_) (TypeDefiner, bool) { // TODO: use TypeDefiner instead of TypeDefiner??
	if ast, ok := typeCache_[input]; !ok {
		return nil, false
	} else {
		return ast, true
	}
}

func Persist(input NameValue_, ast TypeDefiner) {
	// save GraphQL statement to Dynamodb
	dbPersist(input, ast)
}

func Add2Cache(input NameValue_, obj TypeDefiner) {
	fmt.Println("** Add2Cache ", input)
	typeCache_[input] = obj
}

func fetchInterface(input Name_) (*Interface_, bool, string) {
	if ast, ok := typeCache_[input.Name]; ok {
		if ast_, ok := ast.(*Interface_); !ok {
			return nil, false, fmt.Sprintf(`Implements type "%s" is not an Interface %s`, input, input.AtPosition())
		} else {
			return ast_, true, ""
		}
	} else {
		return nil, true, ""
	}

}

func dbPersist(input NameValue_, ast TypeDefiner) error {
	//
	typeDef := TypeRow{PKey: input.String(), Def: ast.String()}
	av, err := dynamodbattribute.MarshalMap(typeDef)
	if err != nil {
		return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
	}
	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("GraphTypes"),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("%s: %s", "Error: failed to PutItem ", err.Error())
	}
	return nil
}

func DeleteType(input string) error {

	//
	typeDef := PkRow{PKey: input}
	av, err := dynamodbattribute.MarshalMap(typeDef)
	if err != nil {
		return fmt.Errorf("%s: %s", "Error: failed to marshal type definition ", err.Error())
	}
	_, err = db.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("GraphTypes"),
		Key:       av,
	})
	if err != nil {
		return fmt.Errorf(`Error: failed to DeleteItem: "%s"  %s`, input, err.Error())
	}
	return nil
}

func ListCache() []TypeDefiner {
	l := make([]TypeDefiner, len(typeCache_), len(typeCache_))
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
	type pKey struct {
		PKey string
	}
	fmt.Printf("DB Fetch name: [%s]\n", name.String())

	errmsg := "Error in marshall of pKey "
	pkey := pKey{PKey: name.String()}
	av, err := dynamodbattribute.MarshalMap(&pkey)
	if err != nil {
		return "", fmt.Errorf("%s. MarshalMap: %s", errmsg, err.Error())
	}
	input := &dynamodb.GetItemInput{
		Key:       av,
		TableName: aws.String("GraphTypes"),
	}
	input = input.SetReturnConsumedCapacity("TOTAL").SetConsistentRead(false)
	//
	result, err := db.GetItem(input)
	if err != nil {
		fmt.Println("ERROROR")
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
		return "", fmt.Errorf("%s %s: %s", errmsg, "GetItem", errmsg, err.Error())
	}
	fmt.Println("dbFetch: GetItem: Query ConsumedCapacity: \n", result.ConsumedCapacity)
	if len(result.Item) == 0 {
		return "", nil
	}
	rec := &TypeRow{}
	err = dynamodbattribute.UnmarshalMap(result.Item, rec)
	if err != nil {
		errmsg := "error in unmarshal "
		return "", fmt.Errorf("%s. UnmarshalMaps:  %s", errmsg, err.Error())
	}
	return rec.Def, nil
}
