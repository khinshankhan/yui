package slug

import "testing"

// Some slug collision cases are inspired by github-slugger:
// https://github.com/flet/github-slugger

func TestMake(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "basic", input: "Some Title", want: "some-title"},
		{name: "whitespace and punctuation", input: "  Hello, World!  ", want: "hello-world"},
		{name: "underscores and dashes", input: "foo_bar---baz", want: "foo-bar-baz"},
		{name: "camelcase splits into words", input: "MyAPIResponse", want: "my-api-response"},
		{name: "digit boundary splits", input: "RFC123Response", want: "rfc123-response"},
		{name: "unicode letters preserved", input: "Привет 你好", want: "привет-你好"},
		{name: "symbols stripped", input: "Rock & Roll / 1977", want: "rock-roll-1977"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Make(tt.input); got != tt.want {
				t.Fatalf("Make(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNextAvailable(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		reserved []string
		want     string
	}{
		{name: "unused", base: "some-title", want: "some-title"},
		{name: "taken", base: "some-title", reserved: []string{"some-title"}, want: "some-title-1"},
		{name: "multiple taken", base: "some-title", reserved: []string{"some-title", "some-title-1"}, want: "some-title-2"},
		{name: "gaps", base: "some-title", reserved: []string{"some-title", "some-title-2"}, want: "some-title-1"},
		{name: "empty base", base: "", reserved: []string{"x"}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NextAvailable(tt.base, tt.reserved); got != tt.want {
				t.Fatalf("NextAvailable(%q, %v) = %q, want %q", tt.base, tt.reserved, got, tt.want)
			}
		})
	}
}
