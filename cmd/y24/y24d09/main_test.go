package main

import (
	"fmt"
	"testing"
)

type a struct {
	i int
}

func TestA(t *testing.T) {
	s := []a{{1}, {1}, {1}}

	s[1].i = 3

	fmt.Printf("%+v\n", s)
}
