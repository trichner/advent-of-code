package main

import (
	"bufio"
	"embed"
	"encoding/binary"
	"fmt"
	"runtime"
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

const (
	UNKNOWN SpringState = iota
	DAMAGED
)

type SpringState int

const patternMask = uint64(0x00FFFFFF_FFFFFFFF)

type Pattern struct {
	lo, hi uint64
}

func (p *Pattern) Len() int {
	return int(p.hi >> 56)
}

func (p *Pattern) First() SpringState {
	return SpringState(p.lo & 0x01)
}

func (p *Pattern) IsDamaged() bool {
	return (p.lo | (p.hi & patternMask)) != 0
}

func (p *Pattern) Shift(n int) Pattern {
	l := p.hi >> 56
	if l <= uint64(n) {
		return Pattern{}
	}

	lo := p.lo >> n
	hi := p.hi & patternMask
	if n <= 64 {
		lo |= hi << (64 - n)
	} else {
		lo = hi >> (n - 64)
	}
	hi = hi >> n

	l = l - uint64(n)
	hi |= l << 56
	return Pattern{lo, hi}
}

func (p *Pattern) Get(n int) SpringState {
	if n < 64 {
		return SpringState((p.lo >> n) & 0x01)
	}
	return SpringState((p.hi >> (n % 64)) & 0x01)
}

func newPattern(p []SpringState) Pattern {
	l := len(p)
	if l > 128-8 {
		panic(fmt.Errorf("too long: %d", l))
	}
	var lo, hi uint64
	for i := range p {
		b := uint64(p[i]) << (i % 64)
		if i < 64 {
			lo |= b
		} else {
			hi |= b
		}
	}
	hi |= uint64(l) << 56
	return Pattern{lo: lo, hi: hi}
}

type task struct {
	states []Pattern
	groups []int
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sentinelTask := task{}
	tasks := make(chan task, 1024)

	sentinelResult := -1
	results := make(chan int, 256)

	running := 0
	for i := 0; i < runtime.NumCPU()*2; i++ {
		running++
		go func() {
			for {
				t := <-tasks
				if t.states == nil {
					tasks <- sentinelTask
					results <- sentinelResult
					return
				}
				cache := make(map[string]int, 1024)
				cnt := arrangements(cache, t.states, t.groups)
				results <- cnt
			}
		}()
	}

	nTasks := 0
	for scanner.Scan() {
		text := scanner.Text()

		states, groups := parseUnfoldedLine(text)
		patterns := toPatterns(states)

		nTasks++
		tasks <- task{
			states: patterns,
			groups: groups,
		}
	}
	tasks <- sentinelTask

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	nDoneTasks := 0
	sum := 0
	for running > 0 {
		r := <-results
		if r == sentinelResult {
			running--
		} else {
			nDoneTasks++
			sum += r
			// fmt.Printf("%2.1f%% (%d/%d)\n", float64(nDoneTasks)/float64(nTasks)*100, nDoneTasks, nTasks)
			// fmt.Printf("sum: %d\n", sum)
		}
	}

	fmt.Printf("part two: %d\n", sum)

	if sum != 157383940585037 {
		panic("bad result")
	}
}

func parseUnfoldedLine(l string) ([][]SpringState, []int) {
	splits := strings.Split(l, " ")

	patternRaw := strings.Join(multiply(splits[0], 5), "?")

	pattern := parsePattern(patternRaw)
	groups := repeat(parseGroups(splits[1]), 5)
	return pattern, groups
}

func repeat(arr []int, n int) []int {
	repeated := make([]int, n*len(arr))
	for i := range repeated {
		repeated[i] = arr[i%len(arr)]
	}
	return repeated
}

func multiply[E any](s E, n int) []E {
	multiplied := make([]E, n)
	for i := range multiplied {
		multiplied[i] = s
	}
	return multiplied
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0
	for scanner.Scan() {
		text := scanner.Text()

		_ = text
		states, groups := parseLine(text)
		patterns := toPatterns(states)

		cache := make(map[string]int, 1024)
		cnt := arrangements(cache, patterns, groups)
		sum += cnt
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("part one: %d\n", sum)

	if sum != 7670 {
		panic("bad result")
	}
}

func toPatterns(states [][]SpringState) []Pattern {
	patterns := make([]Pattern, len(states))
	for i, s := range states {
		patterns[i] = newPattern(s)
	}
	return patterns
}

func toKey(patterns []Pattern, groups []int) string {
	var buf strings.Builder

	patternKey := make([]byte, len(patterns)*16)
	for i, pattern := range patterns {
		binary.LittleEndian.PutUint64(patternKey[16*i:], pattern.lo)
		binary.LittleEndian.PutUint64(patternKey[16*i+8:], pattern.hi)
	}
	buf.Write([]byte{byte(len(patterns))})
	buf.Write(patternKey)

	groupKey := make([]byte, len(groups)*2)
	for i, group := range groups {
		binary.LittleEndian.PutUint16(groupKey[2*i:], uint16(group))
	}
	buf.Write([]byte{byte(len(groups))})
	buf.Write(groupKey)
	return buf.String()
}

func arrangements(cache map[string]int, patterns []Pattern, groups []int) int {
	// get rid of the straight forward options

	// nothing to match anymore, a solution
	if len(patterns) == 0 && len(groups) == 0 {
		return 1
	}

	// groups left but no patterns, no solution
	if len(patterns) == 0 && len(groups) > 0 {
		return 0
	}

	// patterns left but no groups
	if len(patterns) > 0 && len(groups) == 0 {
		for _, s := range patterns {
			if s.IsDamaged() {
				// no solution, needs to match a group
				return 0
			}
		}
		// only wildcards left, we're good
		return 1
	}

	k := toKey(patterns, groups)
	r, ok := cache[k]
	if ok {
		return r
	}

	group := groups[0]
	pattern := patterns[0]

	// pattern cannot match next group
	if pattern.Len() < group {
		// damaged group but size too large
		if pattern.IsDamaged() {
			return 0
		} else {
			// skip a pattern
			return arrangements(cache, patterns[1:], groups)
		}
	}

	// if they are equal length
	if pattern.Len() == group {
		sum := 0

		// is skipping an option?
		if !pattern.IsDamaged() {
			sum += arrangements(cache, patterns[1:], groups)
		}

		// consume pattern
		sum += arrangements(cache, patterns[1:], groups[1:])

		cache[k] = sum

		return sum
	}

	// from here on, pattern MUST  be longer than the group
	sum := 0

	// can we skip the first one?
	if pattern.First() == UNKNOWN {
		copied := shallowCopy(patterns)
		copied[0] = copied[0].Shift(1)
		sum += arrangements(cache, copied, groups)
	}

	// can we skip the last one?
	if pattern.Get(group) != DAMAGED {
		nextPatterns := patterns
		first := nextPatterns[0]
		first = first.Shift(group + 1)
		if first.Len() == 0 {
			nextPatterns = nextPatterns[1:]
		} else {
			nextPatterns = shallowCopy(nextPatterns)
			nextPatterns[0] = first
		}
		sum += arrangements(cache, nextPatterns, groups[1:])
	}

	cache[k] = sum

	return sum
}

func shallowCopy(states []Pattern) []Pattern {
	copied := make([]Pattern, len(states))
	copy(copied, states)
	return copied
}

func parseLine(l string) ([][]SpringState, []int) {
	splits := strings.Split(l, " ")

	pattern := parsePattern(splits[0])
	groups := parseGroups(splits[1])
	return pattern, groups
}

func parseGroups(l string) []int {
	splits := strings.Split(l, ",")
	var groups []int
	for _, s := range splits {
		groups = append(groups, util.Must(strconv.Atoi(s)))
	}
	return groups
}

func parsePattern(l string) [][]SpringState {
	var states [][]SpringState

	var group []SpringState
	for _, c := range l {
		if c == '#' {
			group = append(group, DAMAGED)
		} else if c == '?' {
			group = append(group, UNKNOWN)
		} else if c == '.' {
			if group != nil {
				states = append(states, group)
			}
			group = nil
		} else {
			panic("wut?")
		}
	}
	if group != nil {
		states = append(states, group)
	}
	return states
}
