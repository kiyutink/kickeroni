package db

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var MatchesTable = "kickeroni_matches"
var PlayersTable = "kickeroni_players"

func Connect() *dynamodb.DynamoDB {
	AWSSession := session.Must(session.NewSessionWithOptions(session.Options{}))
	dynamoDBClient := dynamodb.New(AWSSession)
	return dynamoDBClient
}

var dbConnection = Connect()

func CreateRecord(tableName string, record interface{}) error {
	attributeValue, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return err
	}
	_, err = dbConnection.PutItem(&dynamodb.PutItemInput{
		Item:      attributeValue,
		TableName: aws.String(tableName),
	})

	return err
}

func ReplaceRecord(tableName string, id string, record interface{}) error {

	err := DeleteRecord(tableName, id)
	if err != nil {
		return err
	}

	err = CreateRecord(tableName, record)
	return err
}

func DeleteRecord(tableName string, id string) error {
	attributeValue, err := dynamodbattribute.MarshalMap(struct {
		Id string `json:"id"`
	}{Id: id})

	if err != nil {
		return err
	}

	_, err = dbConnection.DeleteItem(
		&dynamodb.DeleteItemInput{
			Key:       attributeValue,
			TableName: aws.String(tableName),
		},
	)
	return err
}

func GetAllRecords(tableName string) ([]map[string]*dynamodb.AttributeValue, error) {
	result, err := dbConnection.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		return []map[string]*dynamodb.AttributeValue{}, err
	}

	return result.Items, nil
}

func GetFirstRecordByFieldValue(tableName string, fieldName string, fieldValue interface{}, bindTo interface{}) (bool, error) {
	filter := expression.Name(fieldName).Equal(expression.Value(fieldValue))
	expression, err := expression.NewBuilder().WithFilter(filter).Build()

	if err != nil {
		return false, err
	}

	result, err := dbConnection.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		FilterExpression:          expression.Filter(),
		TableName:                 aws.String(tableName),
	})

	if err != nil {
		return false, err
	}

	if len(result.Items) == 0 {
		return false, nil
	}

	if err = dynamodbattribute.UnmarshalMap(result.Items[0], bindTo); err != nil {
		return false, err
	}

	return true, nil
}
