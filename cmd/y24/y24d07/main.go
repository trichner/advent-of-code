package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"aoc/pkg/sio"
	"aoc/pkg/util"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

func main() {
	partOne()
	partTwo()
}

type equation struct {
	result int
	values []int
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()
	equations := readEquations(file)

	sum := 0
	for _, eq := range equations {
		if checkValid2(eq, nil) {
			sum += eq.result
		}
	}

	fmt.Printf("part two: %d\n", sum)
}

func checkValid2(e equation, operators []byte) bool {
	current := eval(e.values, operators)
	if len(operators) == len(e.values)-1 {
		return e.result == current
	}

	if current > e.result {
		return false
	}

	// must. go. deeper.
	nextOperators := append(operators, '*')
	if checkValid2(e, nextOperators) {
		return true
	}

	nextOperators[len(nextOperators)-1] = '|'
	if checkValid2(e, nextOperators) {
		return true
	}

	nextOperators[len(nextOperators)-1] = '+'
	return checkValid2(e, nextOperators)
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	equations := readEquations(file)

	sum := 0
	for _, eq := range equations {
		if checkValid1(eq, nil) {
			sum += eq.result
		}
	}

	fmt.Printf("part one: %d\n", sum)
}

func checkValid1(e equation, operators []byte) bool {
	current := eval(e.values, operators)
	if len(operators) == len(e.values)-1 {
		return e.result == current
	}

	if current > e.result {
		return false
	}

	// must. go. deeper.
	nextOperators := append(operators, '*')
	if checkValid1(e, nextOperators) {
		return true
	}

	nextOperators[len(nextOperators)-1] = '+'
	return checkValid1(e, nextOperators)
}

func eval(values []int, operators []byte) int {
	result := values[0]
	for i := 0; i < len(operators); i++ {
		if operators[i] == '+' {
			result += values[i+1]
		} else if operators[i] == '*' {
			result *= values[i+1]
		} else if operators[i] == '|' {
			// concat
			c := values[i+1]
			digits := countDigits(c)
			result = result*pow10(digits) + c
		} else {
			log.Fatalf("unexpected op: %c", operators[i])
		}
	}

	return result
}

func pow10(exp int) int {
	r := 1
	for i := 0; i < exp; i++ {
		r *= 10
	}
	return r
}

func countDigits(i int) int {
	if i == 0 {
		return 1
	}
	n := 0
	for i != 0 {
		i /= 10
		n++
	}
	return n
}

func readEquations(r io.Reader) []equation {
	scanner := bufio.NewScanner(r)

	var equations []equation
	for scanner.Scan() {
		text := scanner.Text()

		splits := strings.Split(text, ":")
		res := util.Must(strconv.Atoi(splits[0]))
		values := sio.IntFieldsByWhitespace(splits[1])

		equations = append(equations, equation{res, values})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return equations
}
