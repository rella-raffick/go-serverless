package handlers

import (
	"GoLangServerless/package/user"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"net/http"
)

var errorMethodNotAllowed = "Method Not Allowed"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, tableName string, dynamoDbClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	email := req.QueryStringParameters["email"]
	if len(email) > 0 {
		result, err := user.FetchUser(email, tableName, dynamoDbClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}
	result, err := user.FetchUsers(tableName, dynamoDbClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynamoDbClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.CreateUser(req, tableName, dynamoDbClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func EditUser(req events.APIGatewayProxyRequest, tableName string, dynamoDbClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	result, err := user.UpdateUser(req, tableName, dynamoDbClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynamoDbClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	err := user.DeleteUser(req, tableName, dynamoDbClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, nil)
}

func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, errorMethodNotAllowed)
}
