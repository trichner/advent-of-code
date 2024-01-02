package main

import (
	"bufio"
	"embed"
	"fmt"
	"strings"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

type BagSample struct {
	RedCount, GreenCount, BlueCount int
}

func (b BagSample) String() string {
	return fmt.Sprintf("(%d red, %d green, %d blue)", b.RedCount, b.GreenCount, b.BlueCount)
}

func main() {
	partOne()
	partTwo()
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0

	for scanner.Scan() {
		text := scanner.Text()

		id, bags := parseGame(text)

		var minRed, minGreen, minBlue int
		for _, b := range bags {
			minRed = max(minRed, b.RedCount)
			minGreen = max(minGreen, b.GreenCount)
			minBlue = max(minBlue, b.BlueCount)
		}
		power := minGreen * minRed * minBlue

		fmt.Printf("%d: %+v = %d\n", id, bags, power)

		sum += power

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("%d\n", sum)
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	maxRed := 12
	maxGreen := 13
	maxBlue := 14

	sum := 0

	for scanner.Scan() {
		text := scanner.Text()

		id, bags := parseGame(text)

		isPossible := true
		for _, b := range bags {
			if b.RedCount > maxRed || b.GreenCount > maxGreen || b.BlueCount > maxBlue {
				isPossible = false
				break
			}
		}
		if isPossible {
			sum += id
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("%d\n", sum)
}

func parseGame(ln string) (int, []BagSample) {
	var gameId int
	_, err := fmt.Sscanf(ln, "Game %d:", &gameId)
	if err != nil {
		panic(err)
	}

	splits := strings.SplitN(ln, ":", 2)
	ln = splits[1]

	splits = strings.Split(ln, ";")
	bags := make([]BagSample, len(splits))
	for i, bagsplit := range splits {
		bags[i] = parseBag(bagsplit)
	}

	return gameId, bags
}

func parseBag(s string) BagSample {
	var bag BagSample
	splits := strings.Split(s, ",")
	for _, cubes := range splits {
		c := strings.TrimSpace(cubes)

		var count int
		var color string
		_, err := fmt.Sscanf(c, "%d %s", &count, &color)
		if err != nil {
			panic("failed to parse bag: " + s)
		}

		switch color {
		case "red":
			bag.RedCount = count
			break
		case "green":
			bag.GreenCount = count
			break
		case "blue":
			bag.BlueCount = count
			break
		}
	}
	return bag
}
