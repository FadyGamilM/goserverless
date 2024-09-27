package main

import (
	"github.com/FadyGamilM/goserverless/pkg/dynamodb"
	"github.com/FadyGamilM/goserverless/pkg/golambda"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	dynamoClient := dynamodb.CreateDynamodbSession()
	lambda.Start(golambda.LambdaHandler(dynamoClient))
}
