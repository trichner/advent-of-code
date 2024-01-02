package y23d07p2

import (
	"bufio"
	"cmp"
	"embed"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

// A, K, Q, J, T, 9, 8, 7, 6, 5, 4, 3, or 2
type Card int

const (
	Ace Card = 13 - iota
	King
	Queen
	Ten
	Nine
	Eight
	Seven
	Six
	Five
	Four
	Three
	Two
	Jack
)

type Hand struct {
	Cards     []Card
	CardsRuns [][]Card
	Bid       int
}

func PartTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var hands []Hand
	for scanner.Scan() {
		text := scanner.Text()

		cards, bid := parseLine(text)
		runs := getRuns(cards)

		hands = append(hands, Hand{CardsRuns: runs, Bid: bid, Cards: cards})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	slices.SortFunc(hands, cmpHand)

	sum := 0
	for i, h := range hands {
		sum += (i + 1) * h.Bid
	}

	// printHands(hands)

	fmt.Printf("part two: %d\n", sum)

	if sum != 249817836 {
		panic("bad result")
	}
}

func printHands(hands []Hand) {
	for _, h := range hands {
		fmt.Printf("%+v %+v %d\n", h.Cards, h.CardsRuns, h.Bid)
	}
}

func parseLine(l string) ([]Card, int) {
	fields := strings.Fields(l)

	var cards []Card
	for i := 0; i < len(fields[0]); i++ {
		c := fields[0][i]
		card := parseCard(c)
		cards = append(cards, card)
	}

	bid, err := strconv.Atoi(fields[1])
	if err != nil {
		panic(err)
	}
	return cards, bid
}

func parseCard(c byte) Card {
	cardStrings := map[byte]Card{
		'A': Ace,
		'K': King,
		'Q': Queen,
		'J': Jack,
		'T': Ten,
		'9': Nine,
		'8': Eight,
		'7': Seven,
		'6': Six,
		'5': Five,
		'4': Four,
		'3': Three,
		'2': Two,
	}

	card, ok := cardStrings[c]
	if !ok {
		panic(fmt.Errorf("unknown card: %s", string([]byte{c})))
	}
	return card
}

func cmpHand(h1 Hand, h2 Hand) int {
	a := h1.CardsRuns
	b := h2.CardsRuns

	for i := 0; i < min(len(a), len(b)); i++ {
		ra := a[i]
		rb := b[i]
		order := cmp.Compare(len(ra), len(rb))
		if order != 0 {
			return order
		}
	}

	return cmpCards(h1.Cards, h2.Cards)
}

func cmpCards(c1, c2 []Card) int {
	for i := range c1 {
		order := cmp.Compare(c1[i], c2[i])
		if order != 0 {
			return order
		}
	}
	panic("wtf?")
}

func getRuns(c []Card) [][]Card {
	cards := make([]Card, len(c))
	copy(cards, c)

	slices.Sort(cards)

	var runs [][]Card
	var run []Card
	var jokers []Card
	for _, c := range cards {
		if c == Jack {
			jokers = append(jokers, c)
			continue
		}
		if len(run) == 0 {
			run = append(run, c)
			continue
		}
		if run[0] == c {
			run = append(run, c)
			continue
		}
		runs = append(runs, run)
		run = []Card{c}
	}
	if len(run) > 0 {
		runs = append(runs, run)
	}

	slices.SortFunc(runs, func(a, b []Card) int {
		if len(a) == len(b) {
			return int(b[0] - a[0])
		}
		return len(b) - len(a)
	})

	if len(runs) == 0 {
		return append(runs, jokers)
	}

	runs[0] = append(runs[0], jokers...)

	return runs
}
