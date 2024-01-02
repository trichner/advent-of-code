package main

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"runtime"
	"strings"

	"aoc/pkg/in"
	"aoc/pkg/sio"
)

//go:embed *.txt
var inputs embed.FS

type MapRange struct {
	DestStart, SrcStart, Len int
}

type Mapping struct {
	From, To  string
	MapRanges []MapRange
}

func (m *Mapping) Map(from int) int {
	for _, e := range m.MapRanges {
		if e.SrcStart <= from && e.SrcStart+e.Len > from {
			to := e.DestStart + (from - e.SrcStart)
			return to
		}
	}
	return from
}

func main() {
	partOne()
	partTwoWorker()
}

type task struct {
	startId, endId int
}

func partTwoWorker() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	seeds, maps := parse(file)

	seedCount := countSeeds(seeds)
	fmt.Printf("seed count: %d\n", seedCount)

	TASK_SENTINEL := task{}
	RESULT_SENTINEL := -1
	tasks := make(chan task, 512)
	results := make(chan int, 512)

	runningCount := 0
	for i := 0; i < runtime.NumCPU()*2; i++ {
		runningCount++
		go func() {
			minLocation := math.MaxInt
			for {
				t := <-tasks
				if t == TASK_SENTINEL {
					tasks <- t
					break
				}
				for j := t.startId; j < t.endId; j++ {
					l := walk(maps, "seed", j)
					minLocation = min(minLocation, l)
				}
				results <- minLocation
			}
			results <- RESULT_SENTINEL
		}()
	}

	for i := 0; i < len(seeds); i += 2 {
		startId := seeds[i]
		endId := seeds[i] + seeds[i+1]
		step := 10_000_000
		for taskStart := startId; taskStart < endId; taskStart += step {
			tasks <- task{
				startId: taskStart,
				endId:   min(endId, taskStart+step),
			}
		}
	}
	tasks <- task{}

	totalCount := runningCount
	minLocation := math.MaxInt
	for runningCount > 0 {
		r := <-results
		if r == RESULT_SENTINEL {
			runningCount--
			fmt.Printf("%d/%d done\n", totalCount-runningCount, totalCount)
			continue
		}
		if r < minLocation {
			minLocation = r
			fmt.Printf("min: %d\n", r)
		}
	}

	fmt.Printf("part two: %d\n", minLocation)
	if minLocation != 1493866 {
		panic("bad result!")
	}
}

func countSeeds(seeds []int) int {
	sum := 0
	for i := 0; i < len(seeds); i += 2 {
		sum += seeds[i+1]
	}
	return sum
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	seeds, maps := parse(file)

	minLocation := math.MaxInt
	for _, s := range seeds {
		l := walk(maps, "seed", s)
		minLocation = min(l, minLocation)
	}

	fmt.Printf("part one: %d\n", minLocation)
	if minLocation != 174137457 {
		panic("bad result!")
	}
}

func walk(mappings map[string]*Mapping, to string, id int) int {
	mapping, ok := mappings[to]
	if !ok {
		return id
	}
	nextId := mapping.Map(id)

	nextTo := mapping.To
	return walk(mappings, nextTo, nextId)
}

func parse(file fs.File) ([]int, map[string]*Mapping) {
	scanner := bufio.NewScanner(file)

	text := mustScan(scanner)
	seeds := parseSeeds(text)
	skipEmptyLine(scanner)

	maps := make(map[string]*Mapping)

	for {

		from, to, mappings, err := parseMap(scanner)
		if errors.Is(err, io.EOF) {
			break
		}
		maps[from] = &Mapping{
			From:      from,
			To:        to,
			MapRanges: mappings,
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return seeds, maps
}

func skipEmptyLine(scanner *bufio.Scanner) {
	l := mustScan(scanner)
	if l != "" {
		panic("line was not empty but: " + l)
	}
}

func parseMap(scanner *bufio.Scanner) (string, string, []MapRange, error) {
	if !scanner.Scan() {
		return "", "", nil, io.EOF
	}
	txt := scanner.Text()
	from, to := parseMappingFromAndTo(txt)

	var mappings []MapRange
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		numbers := sio.IntFields(line)
		mappings = append(mappings, MapRange{
			DestStart: numbers[0],
			SrcStart:  numbers[1],
			Len:       numbers[2],
		})
	}

	return from, to, mappings, nil
}

func parseMappingFromAndTo(line string) (string, string) {
	var name string
	if _, err := fmt.Sscanf(line, "%s map:", &name); err != nil {
		panic(err)
	}
	splits := strings.Split(name, "-to-")
	return splits[0], splits[1]
}

func mustScan(scanner *bufio.Scanner) string {
	if !scanner.Scan() {
		panic("failed to scan")
	}
	return scanner.Text()
}

func parseSeeds(text string) []int {
	text = strings.TrimPrefix(text, "seeds: ")
	return sio.IntFields(text)
}
