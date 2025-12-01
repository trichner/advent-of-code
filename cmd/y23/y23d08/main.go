package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"regexp"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

type RawNode struct {
	ID          string
	Left, Right string
}

func main() {
	partOne()
	partTwo()
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	sequence, rawNodes := parse(file)

	var starts []string
	for k := range rawNodes {
		if k[2] == 'A' {
			starts = append(starts, k)
		}
	}

	l := 1
	for _, s := range starts {
		steps := walk(rawNodes, s, sequence, 0, func(s string) bool {
			return s[2] == 'Z'
		})
		l = lcm(l, steps)
	}

	fmt.Printf("part two: %d\n", l)
	if l != 12030780859469 {
		panic("bad result!")
	}
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()
	sequence, rawNodes := parse(file)

	steps := walk(rawNodes, "AAA", sequence, 0, func(s string) bool {
		return s == "ZZZ"
	})

	fmt.Printf("part one: %d\n", steps)
	if steps != 12169 {
		panic("bad result")
	}
}

func walk(nodes map[string]*RawNode, current string, sequence string, position int, dst func(s string) bool) int {
	if dst(current) {
		return 0
	}

	node := nodes[current]
	direction := sequence[position]
	position = (position + 1) % len(sequence)

	if direction == 'R' {
		return walk(nodes, node.Right, sequence, position, dst) + 1
	} else if direction == 'L' {
		return walk(nodes, node.Left, sequence, position, dst) + 1
	}
	panic("wtf?")
}

func parse(f fs.File) (string, map[string]*RawNode) {
	scanner := bufio.NewScanner(f)

	var sequence string

	rawNodes := make(map[string]*RawNode)

	for scanner.Scan() {
		text := scanner.Text()

		if sequence == "" {
			sequence = text
			scanner.Scan()
			continue
		}

		n, l, r := parseLine(text)

		rawNodes[n] = &RawNode{
			ID:    n,
			Left:  l,
			Right: r,
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return sequence, rawNodes
}

var pattern = regexp.MustCompile("([A-Z0-9]{3}) = \\(([A-Z0-9]{3}), ([A-Z0-9]{3})\\)")

func parseLine(line string) (string, string, string) {
	strings := pattern.FindStringSubmatch(line)
	if strings == nil {
		panic(fmt.Errorf("no match: %q", line))
	}

	return strings[1], strings[2], strings[3]
}

func lcm(a, b int) int {
	return (a * b) / gcd(a, b)
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}
