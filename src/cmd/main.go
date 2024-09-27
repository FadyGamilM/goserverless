package main

import (
	goLambda "github.com/FadyGamilM/goserverless/pkg/lambda"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(goLambda.LambdaHandler)
}
