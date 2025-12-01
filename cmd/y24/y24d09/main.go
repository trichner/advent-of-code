package main

import (
	"aoc/pkg/lists"
	"aoc/pkg/util"
	"embed"
	"fmt"
	"io"
	"slices"
	"strings"
	"time"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

func main() {
	start := time.Now()
	partOne()
	partTwo()
	elapsed := time.Since(start)
	fmt.Printf("executed in: %s\n", elapsed)
}
func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	blocks := readDisk(file)

	totalLength := 0
	for _, b := range blocks {
		totalLength += b.size
	}

	rawDisk := make([]int16, totalLength)

	start := 0
	for _, b := range blocks {
		for j := 0; j < b.size; j++ {
			rawDisk[start+j] = int16(b.fid)
		}
		start += b.size
	}

	reversed := make([]int16, len(rawDisk))
	copy(reversed, rawDisk)
	slices.Reverse(reversed)
	compactDisk(rawDisk)

	fmt.Printf("part one: %d\n", checksumDisk(rawDisk))
}

func checksumDisk(disk []int16) int {
	sum := 0
	for i, el := range disk {
		if el < 0 {
			continue
		}
		sum += int(el) * i
	}
	return sum
}
func compactDisk(disk []int16) {

	last := len(disk) - 1
	for i := 0; i < len(disk); i++ {
		c := disk[i]
		if c >= 0 {
			continue
		}
		f := int16(-1)
		for j := 0; j < len(disk); j++ {
			if last-j < 0 || (last-j <= i) {
				break
			}
			fc := disk[last-j]
			if fc >= 0 {
				f = fc
				disk[last-j] = -1
				break
			}
		}
		if f < 0 {
			break
		}
		disk[i] = f
	}
}

func printDisk(disk []int16) {
	for _, u := range disk {
		if u < 0 {
			fmt.Print(". ")
		} else {
			fmt.Printf("%d ", u)
		}
	}
	fmt.Println()
}

type block struct {
	fid  int
	size int
}

func readDisk(r io.Reader) []*block {

	raw := util.Must(io.ReadAll(r))
	str := string(raw)
	str = strings.TrimSpace(str)

	id := 0
	blockRanges := lists.New[*block]()
	for i, rn := range str {

		num := int(byte(rn) - '0')

		if i%2 == 0 {
			//file
			blockRanges.Append(&block{
				fid:  id,
				size: num,
			})

			id++
		} else {
			blockRanges.Append(&block{
				fid:  -1,
				size: num,
			})
		}
	}

	return blockRanges.ToSlice()
}
