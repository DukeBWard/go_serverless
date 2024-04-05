package main

import (
	"os"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/DukeBWard/go_serverless/pkg/handlers"
)

/**
* attached to a type (*string) indicates a pointer to the type.

* attached to a variable in an assignment (*v = ...) indicates an indirect assignment.
 That is, change the value pointed at by the variable.

* attached to a variable or expression (*v) indicates a pointer dereference.
 That is, take the value the variable is pointing at.

& attached to a variable or expression (&v) indicates a reference.
 That is, create a pointer to the value of the variable or to the field.
**/

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

// * and & is standard C mem reference
// & gets the mem address
// * gets info at mem address
// standard dereferencing, but no pointer arithmetic
func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)
	},)

	if err != nil { return }

	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

const tableNames = "LambdaUser"

// first parens is what the function accepts, second is what it returns
func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, tableName, dynaClient)
	case "POST":
		return handlers.createUser(req, tableName, dynaClient)
	case "PUT":
		return handlers.UpdateUser(req, tableName, dynaClient)
	case "DELETE":
		return handlers.DeleteUser(req, tableName, dynaClient)
	}
	default:
		return handlers.UnhandledMethod()

}
