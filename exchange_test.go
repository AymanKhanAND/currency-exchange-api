package main

import "testing"

func TestExchange(t *testing.T) {

	tests := map[string]struct {
		amount float64
		from string
		to string
		date string
		want float64		
	}{
		"USD to GBP": {amount: 530, from: "USD", to: "GBP", date: "2023-06-08", want: 425.448008},
		"GBP to USD": {amount: 342, from: "GBP", to: "USD", date: "2023-06-01", want: 428.845991},
		"faulty params": {amount: 100, from: "invalid", to: "invalid", date: "invalid", want: 0},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := Exchange(tc.amount, tc.from, tc.to, tc.date)

			if got != tc.want {
				t.Errorf("got %f, want %f", got, tc.want)
			}
		})
	}
}
