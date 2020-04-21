package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided     = errors.New("no name was provided in the HTTP body")
	ErrFileFormatIncorrect = errors.New("file was not encoded properly")
	ErrJsonEncode          = errors.New("failure in encoding json")
)

type MyEvent struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	bodyEvent := MyEvent{}

	// Unmarshal the json, return 404 if error
	err := json.Unmarshal([]byte(request.Body), &bodyEvent)
	if err != nil {
		log.Println("Failed run: ", err)
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", bodyEvent.User)

	var result string
	for _, v := range bodyEvent.Message {
		result = string(v) + result
	}

	response := map[string]string{"message": "Hello " + bodyEvent.User + ", " + result}

	//marshall output
	jsonString, err := json.Marshal(response)
	if err != nil {
		log.Printf("error: %s \n", err)
		return events.APIGatewayProxyResponse{}, ErrJsonEncode
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonString),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
