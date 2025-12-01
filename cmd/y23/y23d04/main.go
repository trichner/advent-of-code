package main

import (
	"bufio"
	"embed"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"aoc/pkg/maps"

	"aoc/pkg/in"
	"aoc/pkg/sets"
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

	scanner := bufio.NewScanner(file)

	cards := make(map[int]int)

	for scanner.Scan() {
		text := scanner.Text()

		cardId, winning, ours := mustParseLine(text)

		matches := countMatches(winning, ours)
		cards[cardId] = matches
	}

	sum := 0
	cardIds := sortedKeys(cards)
	for _, card := range cardIds {
		sum += playCard(cards, card)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part two: %d\n", sum)
}

func sortedKeys(m map[int]int) []int {
	keys := maps.Keys(m)
	slices.Sort(keys)
	return keys
}

func playCard(cards map[int]int, id int) int {
	matches, ok := cards[id]
	if !ok {
		return 0
	}

	sum := 1
	for i := 0; i < matches; i++ {
		sum += playCard(cards, id+1+i)
	}
	return sum
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0

	for scanner.Scan() {
		text := scanner.Text()

		_, winning, ours := mustParseLine(text)

		matches := countMatches(winning, ours)
		sum += score(matches)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part one: %d\n", sum)
	if sum != 23847 {
		panic("bad result")
	}
}

func score(winners int) int {
	if winners == 0 {
		return 0
	}
	return int(math.Pow(2, float64(winners-1)))
}

func countMatches(winning, ours []int) int {
	winningSet := sets.New[int]()
	winningSet.PutAll(winning)

	oursSet := sets.New[int]()
	oursSet.PutAll(ours)

	winners := sets.Intersect(winningSet, oursSet)
	return len(winners)
}

func mustParseLine(s string) (int, []int, []int) {
	splits := mustSplitN(s, ":", 2)

	var cardId int
	if _, err := fmt.Sscanf(splits[0], "Card %d", &cardId); err != nil {
		panic(err)
	}

	splits = mustSplitN(splits[1], "|", 2)

	l1 := mustParseNumbers(splits[0])
	l2 := mustParseNumbers(splits[1])
	return cardId, l1, l2
}

func mustParseNumbers(s string) []int {
	parts := strings.Fields(s)

	numbers := make([]int, len(parts))
	for i := range parts {
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			panic(fmt.Errorf("cant parse %q: %w", parts[i], err))
		}
		numbers[i] = n
	}
	return numbers
}

func mustSplitN(s string, sep string, N int) []string {
	splits := strings.Split(s, sep)
	if len(splits) != N {
		panic(fmt.Errorf("unexpected number of splits, expected %d, got %d : %v", N, len(splits), splits))
	}
	return splits
}
