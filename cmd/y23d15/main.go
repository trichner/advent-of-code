package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

func main() {
	partOne()
	partTwo()
}

const (
	OpReplace = '='
	OpRemove  = '-'
)

type Step struct {
	Lense Lense
	Op    int
	Box   int
}

type Lense struct {
	Label       string
	FocalLength int
}

func (l *Lense) String() string {
	return fmt.Sprintf("[%s %d]", l.Label, l.FocalLength)
}

type Box struct {
	ID     int
	Lenses []Lense
}

func (b *Box) FocussingPower() int {
	power := 0
	for i, l := range b.Lenses {
		power += (1 + b.ID) * (i + 1) * l.FocalLength
	}
	return power
}

func (b *Box) String() string {
	return fmt.Sprintf("Box %d: %v", b.ID, b.Lenses)
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	raw := parse(file)
	seq := parseSequence(raw)

	boxes := make([]Box, 256)
	for i := range boxes {
		boxes[i] = Box{ID: i}
	}

	runSequence(boxes, seq)
	sum := focussingPower(boxes)

	fmt.Printf("part two: %d\n", sum)
	if sum != 269410 {
		panic("bad result")
	}
}

func runSequence(boxes []Box, seq []Step) {
	for _, step := range seq {
		box := boxes[step.Box]
		if step.Op == OpRemove {
			box.Lenses = slices.DeleteFunc(box.Lenses, func(l Lense) bool {
				return l.Label == step.Lense.Label
			})
		} else if step.Op == OpReplace {
			i := slices.IndexFunc(box.Lenses, func(l Lense) bool {
				return l.Label == step.Lense.Label
			})

			if i >= 0 {
				box.Lenses[i] = step.Lense
			} else {
				box.Lenses = append(box.Lenses, step.Lense)
			}

		}
		boxes[step.Box] = box
	}
}

func focussingPower(boxes []Box) int {
	sum := 0
	for _, b := range boxes {
		sum += b.FocussingPower()
	}
	return sum
}

func printBoxes(boxes []Box) {
	for _, b := range boxes {
		if len(b.Lenses) > 0 {
			fmt.Println(b.String())
		}
	}
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	steps := parse(file)

	sum := 0
	for _, s := range steps {
		sum += hash(s)
	}

	fmt.Printf("part one: %d\n", sum)

	if sum != 510792 {
		panic("bad result")
	}
}

func hash(s string) int {
	state := 0

	for i := 0; i < len(s); i++ {
		v := s[i]
		state += int(v)
		state *= 17
		state = state % 256
	}
	return state
}

func parseSequence(s []string) []Step {
	steps := make([]Step, len(s))
	for i, raw := range s {
		steps[i] = parseStep(raw)
	}

	return steps
}

var pattern = regexp.MustCompile("([a-z]+)([-=])([0-9]*)")

func parseStep(s string) Step {
	matches := pattern.FindStringSubmatch(s)
	if len(matches) == 0 {
		panic("no macth: " + s)
	}
	seq := matches[1]
	op := matches[2]
	focalLength := 0
	if len(matches[3]) > 0 {
		focalLength, _ = strconv.Atoi(matches[3])
	}
	return Step{
		Lense: Lense{
			Label:       seq,
			FocalLength: focalLength,
		},
		Op:  int(op[0]),
		Box: hash(seq),
	}
}

func parse(f fs.File) []string {
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := scanner.Text()

		return strings.Split(text, ",")

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return nil
}
