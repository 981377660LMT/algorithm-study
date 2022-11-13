package main

import "testing"

func Test_demo(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"Test_demo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			demo()
		})
	}
}
