package user

import (
	"github.com/aws/aws-lambda-go/events"
	"errors"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var ()

// you can use structs as the model
// which allows you to put models and controllers in the same file

type User struct {
	Email     string `json:"email`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email)
			}
		}
	}
}

func FetchUsers() {

}

func CreateUser() {

}

func UpdateUser() {

}

func DeleteUser() error {

}
