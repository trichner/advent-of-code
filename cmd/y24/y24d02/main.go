package main

import (
	"bufio"
	"embed"
	"fmt"
	"strconv"
	"strings"

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

	scanner := bufio.NewScanner(file)

	lvl := -1
	safeCount := 0
	for scanner.Scan() {
		lvl++
		text := scanner.Text()

		nums := readLine(text)

		if checkValid(nums) {
			safeCount++
			continue
		}

		// let's just bruteforce our options
		for i := 0; i < len(nums); i++ {
			mut := cutAt(nums, i)
			if checkValid(mut) {
				safeCount++
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part two: %d\n", safeCount)
}

func cutAt(nums []int, i int) []int {
	mut := make([]int, len(nums)-1)
	p := 0
	for d := 0; d < len(nums); d++ {
		if d == i {
			continue
		}
		mut[p] = nums[d]
		p++
	}
	return mut
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	safeCount := 0
	for scanner.Scan() {
		text := scanner.Text()

		nums := readLine(text)

		if checkValid(nums) {
			safeCount++
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part one: %d\n", safeCount)
}

func checkValid(nums []int) bool {
	last := nums[0]

	ordering := sign(nums[1] - nums[0])

	for i := 1; i < len(nums); i++ {
		current := nums[i]
		diff := current - last
		s := sign(diff)
		if s != ordering {
			return false
		}
		adiff := abs(current - last)
		if adiff < 1 || adiff > 3 {
			return false
		}

		last = current
	}

	return true
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func sign(a int) int {
	if a >= 0 {
		return 1
	}
	return -1
}

func readLine(s string) []int {
	splits := strings.Split(s, " ")

	nums := make([]int, len(splits))
	for i, split := range splits {
		nums[i] = util.Must(strconv.Atoi(split))
	}

	return nums
}
