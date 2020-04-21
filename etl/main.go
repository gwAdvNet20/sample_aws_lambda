package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"

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
	Fname string `json:"filename"`
	File  string `json:"contents"`
}

//parseBrowser parses user agent
func parseBrowser(ua string) string {
	browsers := []string{"Firefox", "Chrome", "Opera", "Safari", "MSIE"}

	for _, b := range browsers {
		if strings.Contains(ua, b) {
			return b
		}
	}
	return "Other"
}

//parseFile will take a slice of strings and parse the fields.
func parseFile(lines []string) map[string]int {
	browsers := make(map[string]int)
	browsers["Firefox"] = 0
	browsers["Chrome"] = 0
	browsers["Opera"] = 0
	browsers["Safari"] = 0
	browsers["MSIE"] = 0
	browsers["other"] = 0

	for _, line := range lines {

		lineSplit := strings.Split(line, " ")
		userAgent := strings.Join(lineSplit[11:], " ")
		browsers[parseBrowser(userAgent)]++

	}
	return browsers
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
	log.Printf("Processing Lambda request %s\n", bodyEvent.Fname)

	data, err := base64.StdEncoding.DecodeString(bodyEvent.File)
	if err != nil {
		log.Printf("error:", err, "\n")
		return events.APIGatewayProxyResponse{}, ErrFileFormatIncorrect
	}

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	bc := parseFile(lines)

	jsonString, err := json.Marshal(bc)
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
