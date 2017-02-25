package main

import (
	"testing"
)

func TestURL__Defaults(t *testing.T) {
	cases := []struct {
		address, path, expected string
	}{
		{"localhost:8080", "submit", "https://localhost:8080/submit"},
		{"localhost:8080", "/submit", "https://localhost:8080/submit"},
		{"http://localhost:8080", "submit", "http://localhost:8080/submit"},
		{"localhost:8080/proxy", "submit", "https://localhost:8080/proxy/submit"},
	}
	for _,tt := range cases {
		u, err := createURL(tt.address, tt.path)
		if err != nil {
			t.Fatalf("'%s' / '%s' - err - %v", tt.address, tt.path, err)
		}
		if u != tt.expected {
			t.Fatalf("'%s' doesn't match expected default address", u)
		}
	}
}
