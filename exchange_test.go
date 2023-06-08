package main

import "testing"

func TestExchange(t *testing.T) {
	
	got := Exchange("USD", "GBP")
	want := 0.802732

	if got != want {
		t.Errorf("got %f, want %f", got, want)
	}
}
