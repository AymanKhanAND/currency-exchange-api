package main

import "testing"

func TestExchange(t *testing.T) {
	
	got := Exchange(530, "USD", "GBP","2023-06-08")
	want := 425.448008

	if got != want {
		t.Errorf("got %f, want %f", got, want)
	}
}
