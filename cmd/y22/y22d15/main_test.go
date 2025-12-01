package main

import "testing"

func Test_area(t *testing.T) {
	type args struct {
		r int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			area(tt.args.r)
		})
	}
}
