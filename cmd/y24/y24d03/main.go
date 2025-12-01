package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"aoc/pkg/util"

	"aoc/pkg/in"
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

	sum := 0
	src := string(util.Must(io.ReadAll(file)))
	enabled := true
	for {

		if enabled {
			ok, l := matchDont(src)
			if ok {
				src = src[l:]
				enabled = false
				continue
			}

			ok, l, a, b := matchMul(src)
			if ok {
				sum += a * b
				src = src[l:]
				continue
			}
		} else {
			ok, l := matchDo(src)
			if ok {
				src = src[l:]
				enabled = true
				continue
			}
		}

		if len(src) < 1 {
			break
		}

		src = src[1:]
	}

	fmt.Printf("part two: %d\n", sum)
}

func matchDo(src string) (ok bool, l int) {
	return matchString(src, "do()")
}

func matchDont(src string) (ok bool, l int) {
	return matchString(src, "don't()")
}

func matchString(src string, s string) (ok bool, l int) {
	if len(src) < len(s) {
		return false, 0
	}
	return src[:len(s)] == s, len(s)
}

var mulPrefixPattern = regexp.MustCompile("^mul\\(([0-9]{1,3}),([0-9]{1,3})\\)")

func matchMul(src string) (ok bool, l int, a int, b int) {
	matches := mulPrefixPattern.FindStringSubmatch(src)
	if len(matches) == 0 {
		return false, 0, 0, 0
	}

	l = len(matches[0])
	a = util.Must(strconv.Atoi(matches[1]))
	b = util.Must(strconv.Atoi(matches[2]))
	ok = true
	return
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	pattern := regexp.MustCompile("mul\\(([0-9]+),([0-9]+)\\)")
	scanner := bufio.NewScanner(file)

	sum := 0
	for scanner.Scan() {
		text := scanner.Text()

		matches := pattern.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			a := util.Must(strconv.Atoi(match[1]))
			b := util.Must(strconv.Atoi(match[2]))
			sum += a * b
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part one: %d\n", sum)
}
