package sio

import (
	"strconv"
	"strings"
)

func IntFields(s string) []int {
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
