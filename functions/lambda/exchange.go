package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const getRequestErrorMsg = "error getting response"
const responseReadErrorMsg = "error reading the body"
const missingParamsErrorMsg = "error getting query parameters"
const floatConversionErrorMsg = "error converting given amount to float"


type Response struct {
	Success bool 				`json:"success"`
	Query 	Query 				`json:"query"`
	Info 	map[string]float64 	`json:"info"`
	Date 	string 				`json:"date"`
	Result 	float64 			`json:"result"`
}

type Query struct {
	From 	string 	`json:"from"`
	To 		string 	`json:"to"`
	Amount 	float64 `json:"amount"`
}

func ExchangeGetRequest(amount float64, from, to, date string) (float64, error) {
	
	url := fmt.Sprintf("https://api.exchangerate.host/convert?from=%s&to=%s&amount=%f&date=%s", from, to, amount, date)
	
	response, err := http.Get(url)
	if err != nil {
		return 0, errors.New(getRequestErrorMsg)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, errors.New(responseReadErrorMsg)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Result, nil
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	amount := request.QueryStringParameters["amount"]
	from := request.QueryStringParameters["from"]
	to := request.QueryStringParameters["to"]
	date := request.QueryStringParameters["date"]

	if amount == "" || from == "" || to == "" || date == "" {
		return events.APIGatewayProxyResponse{
			Body:       missingParamsErrorMsg,
			StatusCode: 400,
		}, nil
	}

	numericAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       floatConversionErrorMsg,
			StatusCode: 400,
		}, nil
	}

	converted, err := ExchangeGetRequest(numericAmount, from, to, date)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	textConverted := fmt.Sprintf("%f", converted)
	return events.APIGatewayProxyResponse{
		Body:       textConverted,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
