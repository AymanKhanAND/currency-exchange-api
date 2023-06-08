package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
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

func ExchangeHandler(w http.ResponseWriter, r *http.Request) {

	amount := r.URL.Query().Get("amount")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	date := r.URL.Query().Get("date")

	if amount == "" || from == "" || to == "" || date == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	numericAmount, _ := strconv.ParseFloat(amount, 64)

	converted, _ := ExchangeGetRequest(numericAmount, from, to, date)

	fmt.Fprint(w, converted)
}

func main() {
	handler := http.HandlerFunc(ExchangeHandler)
	log.Fatal(http.ListenAndServe(":5000", handler))
}
