package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
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
			"USD to GBP": {amount: 530, from: "USD", to: "GBP", date: "2023-06-01", want: 422.669218},
			"GBP to USD": {amount: 342, from: "GBP", to: "USD", date: "2023-06-01", want: 428.845991},
			"invalid but well-formed returns 0": {amount: 100, from: "invalid", to: "invalid", date: "invalid", want: 0},
		}
	
		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				got, err := ExchangeGetRequest(tc.amount, tc.from, tc.to, tc.date)
	
				assertErrorNotPresent(t, err)
				assertConversion(t, got, tc.want)
			})
		}
	})

	t.Run("error with GET", func(t *testing.T) {
		_, err := ExchangeGetRequest(1, "mal formed", "test", "test")

		assertErrorPresent(t, err)
		assertErrorMessage(t, err.Error(), getRequestErrorMsg)
	})

	t.Run("successful conversion through http handler", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/convert?from=USD&to=GBP&amount=200&date=2023-06-01", nil)
		response := httptest.NewRecorder()

		ExchangeHandler(response, request)

		got, _ := strconv.ParseFloat(response.Body.String(), 64)

		assertConversion(t, got, 159.497818)
		assertStatus(t, response.Code, 200)
	})


	t.Run("errors with query params", func(t *testing.T) {
		tests := map[string]struct{
			url string
		}{
			"missing params": {"/convert?from=USD&to=GBP&amount=200"},
			"malformed params": {"/convert?from=US D&to=G BP&amount=200&date=2023-06-08"},
			"non number for amount": {"/convert?from=USD&to=GBP&amount=test&date=2023-06-08"},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				request, _ := http.NewRequest(http.MethodGet, tc.url, nil)
				response := httptest.NewRecorder()
	
				ExchangeHandler(response, request)
	
				assertStatus(t, response.Code, 400)
			})
		}
	})
}

func assertConversion(t testing.TB, got, want float64) {
	t.Helper()
	if got != want {
		t.Errorf("got %f, want %f", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertErrorMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func assertErrorPresent(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("was expecting error and didn't get one")
	}
}

func assertErrorNotPresent(t testing.TB, err error) {
	if err != nil {
		t.Fatal("got error where there shouldn't be")
	}
}
