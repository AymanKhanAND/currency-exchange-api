package main

import (
	"testing"
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
			"USD to GBP": {amount: 530, from: "USD", to: "GBP", date: "2023-06-08", want: 425.448008},
			"GBP to USD": {amount: 342, from: "GBP", to: "USD", date: "2023-06-01", want: 428.845991},
			"invalid but well-formed returns 0": {amount: 100, from: "invalid", to: "invalid", date: "invalid", want: 0},
		}
	
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				got, err := ExchangeGetRequest(tc.amount, tc.from, tc.to, tc.date)
	
				if err != nil {
					t.Fatal("got error where there shouldn't be")
				}
	
				if got != tc.want {
					t.Errorf("got %f, want %f", got, tc.want)
				}
			})
		}
	})

	t.Run("error with GET", func(t *testing.T) {
		_, err := ExchangeGetRequest(400, "bl  ah", "blah", "blah")
		want := "error getting response"

		if err == nil {
			t.Fatal("was expecting error and didn't get one")
		}
		
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}
