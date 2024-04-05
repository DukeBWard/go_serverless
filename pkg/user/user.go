package user

import (
	"encoding/json"
	"errors"

	"github.com/DukeBWard/go_serverless/pkg/validators"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorInvalidUserData         = "invalid user data"
	ErrorInvalidEmail            = "invalid email"
	ErrorUserAlreadyExists       = "user already exists"
	ErrorCouldNotMarhsalItem     = "cannot marshal item"
	ErrorCouldNotPutItem         = "cannot dynamo put item"
	ErrorUserDoesNotExist        = "user does not exist"
	ErrorCouldNotDeleteItem      = "could not delete item"
)

// you can use structs as the model
// which allows you to put models and controllers in the same file

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		// pretty easy way to declare a map in go
		// the keys are string values, while the objects are time dynamodb.AttributeValue
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)
	// getting the result item from dynamodb (json) and unmarshalling it into the user struct (item)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

// don't forget that go uses "slices" not lists
func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	// input now has the mem address of the new dynamodb.ScanInput struct
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// scan is like find for dynamo
	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*User, error) {

	var user User

	// the []byte() turns whatever is inside into a byte slice
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(user.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currUser, _ := FetchUser(user.Email, tableName, dynaClient)
	if currUser != nil && len(currUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	// put the go struct into something dynamo readable
	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarhsalItem)
	}

	// setting up the input struct pointer for putting
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	// user the input struct pointer to actually put into the db
	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotPutItem)
	}

	return &user, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*User, error) {

	var user User
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return nil, errors.New(ErrorInvalidEmail)
	}

	// second param is error usually, where the _ is
	currUser, _ := FetchUser(user.Email, tableName, dynaClient)
	if currUser != nil && len(currUser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExist)
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarhsalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotPutItem)
	}

	return &user, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {

	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil
}
