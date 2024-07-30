package main

import "testing"

func TestPing(t *testing.T) {
	args := []Value{{
		dataType: "bulk",
		bulk:     "baz",
	}}

	pingResponse := ping(args)
	if pingResponse.bulk != "baz" {
		t.Fatalf("ping(%v) = %s, want %s", pingResponse, pingResponse.bulk, "baz")
	}
}
