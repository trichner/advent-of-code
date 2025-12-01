package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
	"os"
)

func main() {
	findMarker(4)
	findMarker(14)
}

func findMarker(l int) {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanBytes)

	pos := 0
	buf := make([]byte, l)

	for i := 0; i < l && scanner.Scan(); i++ {
		b := scanner.Bytes()
		buf[pos%l] = b[0]
		pos++
	}

	for scanner.Scan() {

		b := scanner.Bytes()
		buf[pos%l] = b[0]
		pos++
		if isMarker(buf) {
			break
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d\n", pos)
}

func isMarker(s []byte) bool {
	var set uint32
	for _, b := range s {
		p := b - 'a'
		set |= 1 << p
	}
	return bits.OnesCount32(set) == len(s)
}
