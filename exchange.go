package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

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

func HandleGetResponse(r *http.Response) (float64, error) {
	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, errors.New("error reading the body")
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Result, nil
}

func ExchangeGetRequest(amount float64, from, to, date string) (float64, error) {
	
	url := fmt.Sprintf("https://api.exchangerate.host/convert?from=%s&to=%s&amount=%f&date=%s", from, to, amount, date)
	
	response, err := http.Get(url)

	if err != nil {
		return 0, errors.New("error getting response")
	}

	return HandleGetResponse(response)
}
