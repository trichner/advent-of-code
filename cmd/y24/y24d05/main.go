package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"

	"aoc/pkg/in"
	"aoc/pkg/sio"
	"aoc/pkg/util"
)

//go:embed *.txt
var inputs embed.FS

func main() {
	partOne()
	partTwo()
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	updates, rules := readIn(file)

	updates = filterUnsortedUpdate(updates, rules)

	sum := 0
	ruleGraph := buildRuleGraph(rules)
	for _, update := range updates {
		sorted := toposort(update, ruleGraph)
		sum += sorted[len(sorted)/2]
	}

	fmt.Printf("part two: %d\n", sum)
}

func buildRuleGraph(rules [][]int) map[int][]int {
	graph := make(map[int][]int)
	for _, r := range rules {
		graph[r[0]] = append(graph[r[0]], r[1])
	}
	return graph
}

func toposort(update []int, graph map[int][]int) []int {
	inDegrees := make(map[int]int)
	for _, page := range update {
		edges := graph[page]
		for _, edge := range edges {
			inDegrees[edge] = inDegrees[edge] + 1
		}
	}

	newInDegrees := make(map[int]int)
	for _, page := range update {
		count := inDegrees[page]
		newInDegrees[page] = count
	}

	sorted := make([]int, len(update))
	for node, inCount := range newInDegrees {
		sorted[inCount] = node
	}

	return sorted
}

func filterUnsortedUpdate(updates [][]int, rules [][]int) [][]int {
	var filtered [][]int
	for _, update := range updates {
		if !checkUpdateForRules(update, rules) {
			filtered = append(filtered, update)
		}
	}
	return filtered
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	updates, rules := readIn(file)

	sum := 0
	for _, update := range updates {
		if checkUpdateForRules(update, rules) {
			sum += update[len(update)/2]
		}
	}

	fmt.Printf("part one: %d\n", sum)
}

func checkUpdateForRules(update []int, rule [][]int) bool {
	for _, r := range rule {
		if !checkUpdateForRule(update, r) {
			return false
		}
	}

	return true
}

func checkUpdateForRule(update []int, rule []int) bool {
	for i := 0; i < len(update); i++ {
		head := update[i]
		if rule[0] == head {
			// do we violate anything before?
			for j := 0; j < i; j++ {
				after := update[j]
				if rule[1] == after {
					return false
				}
			}
		}

	}
	for i := 0; i < len(update); i++ {
		head := update[i]
		if rule[1] == head {
			// do we violate anything after?
			for j := i; j < len(update); j++ {
				after := update[j]
				if after == rule[0] {
					return false
				}
			}
		}

	}
	return true
}

func readIn(r io.Reader) ([][]int, [][]int) {
	var updates, rules [][]int

	scanner := bufio.NewScanner(r)

	// rules
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		var a, b int
		util.Must(fmt.Sscanf(text, "%d|%d", &a, &b))
		rules = append(rules, []int{a, b})
	}

	// updates
	for scanner.Scan() {
		text := scanner.Text()

		update := sio.IntFieldsByChar(text, ',')
		updates = append(updates, update)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return updates, rules
}
