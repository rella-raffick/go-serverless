package user

import (
	"GoLangServerless/package/validators"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorFailedToFetchRecord     = "Error! Failed to Fetch Record"
	ErrorFailedToMarshalItem     = "Error! Failed To Marshal Item"
	ErrorFailedToUnmarshalRecord = "Error! Failed To Unmarshal Record"
	ErrorInvalidUserData         = "Error! Invalid User Data"
	ErrorInvalidEmailData        = "Error! Invalid Email"
	ErrorFailedToDelete          = "Error! Failed to Delete Record"
	ErrorFailedToPut             = "Error! Failed to Edit Record"
	ErrorUserExists              = "Error! User Already Exists"
	ErrorUserDoesNotExists       = "Error! User Does Not Exists"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email, tableName string, dynamoDbClient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := dynamoDbClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil
}

func FetchUsers(tableName string, dynamoDbClient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynamoDbClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	item := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item)
	return item, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynamoDbClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User

	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}
	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmailData)
	}
	currentUser, _ := FetchUser(u.Email, tableName, dynamoDbClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserExists)
	}

	result, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorFailedToMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      result,
		TableName: aws.String(tableName),
	}

	_, err = dynamoDbClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToPut)
	}
	return &u, nil
}
func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynamoDbClient dynamodbiface.DynamoDBAPI) (*User, error) {
	var u User

	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}
	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmailData)
	}
	currentUser, _ := FetchUser(u.Email, tableName, dynamoDbClient)
	if currentUser != nil && len(currentUser.Email) == 0 {
		return nil, errors.New(ErrorUserExists)
	}

	result, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorFailedToMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      result,
		TableName: aws.String(tableName),
	}

	_, err = dynamoDbClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToPut)
	}
	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynamoDbClient dynamodbiface.DynamoDBAPI) error {

	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynamoDbClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorFailedToDelete)
	}
	return nil
}
