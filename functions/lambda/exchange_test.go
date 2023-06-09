package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestExchange(t *testing.T) {

	t.Run("currency exchange", func(t *testing.T) {
		tests := map[string]struct {
			amount float64
			from string
			to string
			date string
			want float64		
		}{
			"USD to GBP": {amount: 530, from: "USD", to: "GBP", date: "2023-06-01", want: 422.669218},
			"GBP to USD": {amount: 342, from: "GBP", to: "USD", date: "2023-06-01", want: 428.845991},
			"invalid but well-formed returns 0": {amount: 100, from: "invalid", to: "invalid", date: "invalid", want: 0},
		}
	
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				got, err := ExchangeGetRequest(tc.amount, tc.from, tc.to, tc.date)
	
				assert.Equal(t, got, tc.want, "they should be equal")
				assert.Nil(t, err, "not expecting an error but got one")
			})
		}
	})

	t.Run("error with GET", func(t *testing.T) {
		_, err := ExchangeGetRequest(1, "mal formed", "test", "test")

		assert.NotNil(t, err, "expecting an error but didn't get one")
		assert.Equal(t, err.Error(), getRequestErrorMsg)
	})

	t.Run("successful conversion through api gateway", func(t *testing.T) {
		request := events.APIGatewayV2HTTPRequest{
			QueryStringParameters: map[string]string{
				"from": "USD",
				"to": "GBP",
				"date": "2023-06-01",
				"amount": "200",
			},
		}

		want := "159.497818"

		got, err := Handler(request)

		if got.Body != want {
			t.Errorf("got %s want %s", got.Body, want)
		}

		assert.Equal(t, got.StatusCode, 200, "they should be equal")
		assert.Nil(t, err, "not expecting an error but got one")
	})


	t.Run("errors with query params", func(t *testing.T) {
		tests := map[string]struct{
			from, to, date, amount string
		}{
			"missing params": {"USD", "GBP", "", "200"},
			"malformed params": {"US D", "G BP", "2023-06-08", "200"},
			"non number for amount": {"USD", "GBP", "2023-01-01", "notanumber"},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
	
				request := events.APIGatewayV2HTTPRequest{
					QueryStringParameters: map[string]string{
						"from": tc.from,
						"to": tc.to,
						"date": tc.date,
						"amount": tc.amount,
					},
				}

				got, err := Handler(request)
	
				assert.Equal(t, got.StatusCode, 400, "they should be equal")
				assert.Nil(t, err, "not expecting an error but got one")
			})
		}
	})
}
