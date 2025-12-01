package sio

import (
	"strconv"
	"strings"
)

func IntFieldsByWhitespace(s string) []int {
	splits := strings.Fields(s)

	numbers := make([]int, len(splits))

	for i := range splits {
		n, err := strconv.Atoi(splits[i])
		if err != nil {
			panic(err)
		}
		numbers[i] = n
	}
	return numbers
}

func IntFieldsByChar(s string, splitter rune) []int {
	splits := strings.FieldsFunc(s, func(r rune) bool {
		return r == splitter
	})

	numbers := make([]int, len(splits))

	for i := range splits {
		n, err := strconv.Atoi(splits[i])
		if err != nil {
			panic(err)
		}
		numbers[i] = n
	}
	return numbers
}
