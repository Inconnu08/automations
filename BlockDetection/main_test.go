package main

import (
	"reflect"
	"testing"
)

func TestIsBlock(t *testing.T) {
	tests := []struct {
		testName string
		input string
		want  bool
	}{
		{testName: "Test with proper conf block", input: "location ~ ^/(images|javascript|js|css|flash|media|static)/{\n  root/var/www/virtual/big.server.com/htdocs;\nexpires 30d;\n}", want: true},
		{testName: "Test with separate double blocks", input: "lgbjhfbg{nreigni}nrgn{jrngj}", want: true},
		{testName: "Test with not blocks provided", input: "gregergger", want: false},
	}

	for _, tc := range tests {
		got := IsBlock(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("test %s: expected: %v, got: %v", tc.testName, tc.want, got)
		}
	}
}

func TestIsLine(t *testing.T) {
	tests := []struct {
		testName string
		input string
		want  bool
	}{
		{testName: "Test with an actual line", input: "location ~ ^/(images|javascript|js|css|flash|media|static)/{\n", want: true},
		{testName: "Test with not an actual line", input: "lgbjhfbg{nreigni}nrgn{jrngj}", want: false},
	}

	for _, tc := range tests {
		got := IsLine(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("test %s: expected: %v, got: %v", tc.testName, tc.want, got)
		}
	}
}

func TestHasComment(t *testing.T) {
	tests := []struct {
		testName string
		input string
		want  bool
	}{
		{testName: "Test with an actual comment", input: "#this is a comment", want: true},
		{testName: "Test with double # comment", input: "##this si also a comment", want: true},
		{testName: "Test with space between # and comment", input: "# this si also a comment", want: true},
		{testName: "Test with not a comment", input: "this si not a comment", want: false},
	}

	for _, tc := range tests {
		got := HasComment(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("test %s: expected: %v, got: %v", tc.testName, tc.want, got)
		}
	}
}


