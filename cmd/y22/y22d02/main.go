package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	oRock    = "A"
	oPaper   = "B"
	oScissor = "C"
	pRock    = "X"
	pPaper   = "Y"
	pScissor = "Z"
	pLoss    = "X"
	pDraw    = "Y"
	pWin     = "Z"
)

type Shape int

const (
	ROCK Shape = iota
	PAPER
	SCISSORS
)

type GameResult int

const (
	LOSS GameResult = iota
	DRAW
	WIN
)

func main() {
	partOne()
	partTwo()
}

func partOne() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	total := 0

	for scanner.Scan() {
		text := scanner.Text()

		var opponentPlayStr string
		var myPlayStr string
		_, err := fmt.Sscanf(text, "%s %s", &opponentPlayStr, &myPlayStr)
		if err != nil {
			log.Fatal(err)
		}

		opponentPlay := ParseShape(opponentPlayStr)
		myPlay := ParseShape(myPlayStr)

		res := CalculateResult(myPlay, opponentPlay)

		score := ScoreShape(myPlay) + ScoreResult(res)
		total += score

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%d\n", total)
}

func partTwo() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	total := 0

	for scanner.Scan() {
		text := scanner.Text()

		var opponentPlayStr string
		var expectedResultStr string
		_, err := fmt.Sscanf(text, "%s %s", &opponentPlayStr, &expectedResultStr)
		if err != nil {
			log.Fatal(err)
		}

		opponentPlay := ParseShape(opponentPlayStr)
		expectedResult := ParseResult(expectedResultStr)
		myPlay := GetPlay(opponentPlay, expectedResult)

		res := CalculateResult(myPlay, opponentPlay)

		score := ScoreShape(myPlay) + ScoreResult(res)
		total += score

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("%d\n", total)
}

func GetPlay(opponent Shape, result GameResult) Shape {
	play := map[Shape]map[GameResult]Shape{
		ROCK: {
			WIN:  PAPER,
			LOSS: SCISSORS,
			DRAW: ROCK,
		},
		PAPER: {
			WIN:  SCISSORS,
			LOSS: ROCK,
			DRAW: PAPER,
		},
		SCISSORS: {
			WIN:  ROCK,
			LOSS: PAPER,
			DRAW: SCISSORS,
		},
	}
	return play[opponent][result]
}

func ParseResult(r string) GameResult {
	results := map[string]GameResult{
		pDraw: DRAW,
		pLoss: LOSS,
		pWin:  WIN,
	}
	return results[r]
}

func ScoreShape(s Shape) int {
	scores := map[Shape]int{
		ROCK:     1,
		PAPER:    2,
		SCISSORS: 3,
	}
	return scores[s]
}

func ScoreResult(r GameResult) int {
	scores := map[GameResult]int{
		LOSS: 0,
		DRAW: 3,
		WIN:  6,
	}
	return scores[r]
}

func CalculateResult(a, b Shape) GameResult {
	if a == b {
		return DRAW
	}
	switch a {
	case ROCK:
		if b == SCISSORS {
			return WIN
		} else {
			return LOSS
		}
	case PAPER:
		if b == ROCK {
			return WIN
		} else {
			return LOSS
		}
	case SCISSORS:
		if b == PAPER {
			return WIN
		} else {
			return LOSS
		}
	}
	panic(fmt.Errorf("unexpected shape: %d", a))
}

func ParseShape(s string) Shape {
	if s == oRock || s == pRock {
		return ROCK
	}

	if s == oPaper || s == pPaper {
		return PAPER
	}

	if s == oScissor || s == pScissor {
		return SCISSORS
	}

	panic(fmt.Errorf("unexpected shape: %q", s))
}
