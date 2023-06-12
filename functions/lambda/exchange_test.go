package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestExchange(t *testing.T) {

	d := &dependencies{
		externalService: mockApiCall{},
	}

	t.Run("currency exchange", func(t *testing.T) {
		tests := map[string]struct {
			amount float64
			from   string
			to     string
			date   string
			want   float64
		}{
			"USD to GBP":                        {amount: 200, from: "USD", to: "GBP", date: "2023-06-01", want: 159.497818},
			"invalid but well-formed returns 0": {amount: 100, from: "invalid", to: "invalid", date: "invalid", want: 0},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				got, err := d.externalService.exchangeGetRequest(tc.amount, tc.from, tc.to, tc.date)

				assert.Equal(t, tc.want, got, "they should be equal")
				assert.Nil(t, err, "not expecting an error but got one")
			})
		}
	})

	t.Run("error with GET", func(t *testing.T) {
		_, err := d.externalService.exchangeGetRequest(1, "mal formed", "test", "test")

		assert.Equal(t, getRequestErrorMsg, err.Error())
		assert.NotNil(t, err, "expecting an error but didn't get one")
	})
}

type mockApiCall struct{}

func (m mockApiCall) exchangeGetRequest(amount float64, from, to, date string) (float64, error) {

	params := []string{from, to, date}
	invalidParam := false

	for _, param := range params {
		// returns error if malformed parameters
		if strings.Contains(param, " ") {
			return 0, errors.New(getRequestErrorMsg)
		}

		if param == "invalid" {
			invalidParam = true
		}
	}

	// returns 0 and no error if invalid params (invalid country code/date)
	if invalidParam {
		return 0, nil
	}

	// conversion of 200 USD to GBP at 1/6/23 rate
	return 159.497818, nil
}

func TestHandler(t *testing.T) {
	tests := map[string]struct {
		from, to, date, amount, wantBody string
		wantCode                         int
	}{
		"successful conversion": {"USD", "GBP", "2023-06-01", "200", "159.497818", 200},
		"missing params":        {"USD", "GBP", "", "200", missingParamsErrorMsg, 400},
		"malformed params":      {"US D", "G BP", "2023-06-08", "200", getRequestErrorMsg, 400},
		"non number for amount": {"USD", "GBP", "2023-01-01", "notanumber", floatConversionErrorMsg, 400},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			request := events.APIGatewayV2HTTPRequest{
				QueryStringParameters: map[string]string{
					"from":   tc.from,
					"to":     tc.to,
					"date":   tc.date,
					"amount": tc.amount,
				},
			}

			d := &dependencies{
				externalService: mockApiCall{},
			}

			got, err := d.Handler(request)

			assert.Equal(t, tc.wantBody, got.Body, "they should be equal")
			assert.Equal(t, tc.wantCode, got.StatusCode, "they should be equal")
			assert.Nil(t, err, "not expecting an error but got one")
		})
	}
}
