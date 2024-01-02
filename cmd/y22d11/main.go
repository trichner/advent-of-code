package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Monkey struct {
	id             int
	items          []int
	operation      func(w int) int // weight -> adjusted weight
	next           func(w int) int // weight -> monkey
	itemsInspected int
	divisibleBy    int
}

func main() {
	playMonkeyGame(20, true)
	playMonkeyGame(10000, false)
}

func playMonkeyGame(N int, isLessWorry bool) {
	monkeys := parseMonkeys("input.txt")

	modulus := 1
	for _, m := range monkeys {
		modulus *= m.divisibleBy
	}

	for i := 0; i < N; i++ {
		for _, m := range monkeys {
			for _, item := range m.items {
				m.itemsInspected++
				worry := m.operation(item)
				if isLessWorry {
					worry /= 3
				}
				worry = worry % modulus
				next := m.next(worry)
				monkeys[next].items = append(monkeys[next].items, worry)
			}
			m.items = []int{}
		}
	}

	score := calculateScore(monkeys)
	fmt.Printf("%d\n", score)
}

func makeNext(divisible int, monkeyIfTrue int, monkeyIfFalse int) func(w int) int {
	return func(w int) int {
		if w%divisible == 0 {
			return monkeyIfTrue
		}
		return monkeyIfFalse
	}
}

func mustScanf(str string, format string, a ...any) {
	_, err := fmt.Sscanf(str, format, a...)
	if err != nil {
		panic(err)
	}
}

func parseOperation(s string) func(w int) int {
	if s == "new = old * old" {
		return func(w int) int {
			return w * w
		}
	}
	var op string
	var arg int
	_, err := fmt.Sscanf(s, "new = old %s %d", &op, &arg)
	if err != nil {
		panic(err)
	}
	if op == "*" {
		return func(w int) int {
			return w * arg
		}
	}
	if op == "+" {
		return func(w int) int {
			return w + arg
		}
	}
	panic(fmt.Errorf("unknown op: %q", op))
}

func mustScanLine(scanner *bufio.Scanner, format string, a ...any) {
	ok := scanner.Scan()
	if !ok {
		panic(scanner.Err())
	}
	text := scanner.Text()
	mustScanf(text, format, a...)
}

func mustReadLine(scanner *bufio.Scanner) string {
	ok := scanner.Scan()
	if !ok {
		panic(scanner.Err())
	}
	return scanner.Text()
}

func mustScanItems(scanner *bufio.Scanner) []int {
	ok := scanner.Scan()
	if !ok {
		panic(scanner.Err())
	}
	text := scanner.Text()
	text = strings.TrimPrefix(text, "  Starting items:")
	splits := strings.Split(text, ",")

	var items []int
	for _, s := range splits {
		d, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			panic(err)
		}
		items = append(items, d)
	}
	return items
}

func parseMonkeys(fname string) []*Monkey {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var monkeys []*Monkey
	for {
		var id int
		mustScanLine(scanner, "Monkey %d:", &id)

		items := mustScanItems(scanner)

		op := mustReadLine(scanner)
		op = strings.TrimPrefix(op, "  Operation: ")
		operation := parseOperation(op)

		var divisibleBy int
		mustScanLine(scanner, "  Test: divisible by %d", &divisibleBy)
		var monkeyA, monkeyB int
		mustScanLine(scanner, "    If true: throw to monkey %d", &monkeyA)
		mustScanLine(scanner, "    If false: throw to monkey %d", &monkeyB)

		m := &Monkey{
			id:          id,
			items:       items,
			operation:   operation,
			next:        makeNext(divisibleBy, monkeyA, monkeyB),
			divisibleBy: int(divisibleBy),
		}
		monkeys = append(monkeys, m)

		if !scanner.Scan() {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return monkeys
}

func calculateScore(monkeys []*Monkey) int {
	var inspections []int
	for _, m := range monkeys {
		inspections = append(inspections, m.itemsInspected)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(inspections)))
	return inspections[0] * inspections[1]
}
