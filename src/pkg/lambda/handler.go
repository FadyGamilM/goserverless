package lambda

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type LambdaHandlerDef func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func LambdaHandler(dynamodbClient *dynamodb.DynamoDB) LambdaHandlerDef {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.HTTPMethod {
		case "POST":
			return createItem(dynamodbClient, request)
		case "GET":
			return getItem(dynamodbClient, request)
		case "PUT":
			return updateItem(dynamodbClient, request)
		case "DELETE":
			return deleteItem(dynamodbClient, request)
		default:
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid method"}, nil
		}
	}
}

func createItem(dynamoClient *dynamodb.DynamoDB, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var note Note
	if err := json.Unmarshal([]byte(request.Body), &note); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request body"}, nil
	}

	av, err := dynamodbattribute.MarshalMap(note)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error marshalling item"}, nil
	}
	log.Println("av is : ", av)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME")),
	}

	_, err = dynamoClient.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error putting item in DynamoDB"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 201, Body: "Note created successfully"}, nil
}

func getItem(dynamoClient *dynamodb.DynamoDB, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	query := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME")),
	}

	result, err := dynamoClient.GetItem(query)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error getting item from DynamoDB"}, nil
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Note not found"}, nil
	}

	var note Note
	err = dynamodbattribute.UnmarshalMap(result.Item, &note)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error unmarshalling item"}, nil
	}

	body, _ := json.Marshal(note)
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(body)}, nil
}

// TODO : we need to implement the optimistic concurrency control via the concept of version that dynamodb provides
func updateItem(dynamoClient *dynamodb.DynamoDB, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var note Note
	if err := json.Unmarshal([]byte(request.Body), &note); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request body"}, nil
	}

	query := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {S: aws.String(note.Content)},
			":p": {S: aws.String(note.Password)},
			":u": {S: aws.String(note.URL)},
		},
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(note.Id)},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set content = :c, password = :p, url = :u"),
	}

	_, err := dynamoClient.UpdateItem(query)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error updating item in DynamoDB"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Note updated successfully"}, nil
}

func deleteItem(dynamoClient *dynamodb.DynamoDB, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	query := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		TableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME")),
	}

	_, err := dynamoClient.DeleteItem(query)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error deleting item from DynamoDB"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "Note deleted successfully"}, nil
}
