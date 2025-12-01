package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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

	d := &Display{next: 20, signal: 1}

	for scanner.Scan() {
		text := scanner.Text()
		d.Write(text)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d\n", d.sum)
}

type Display struct {
	cycle  int
	next   int
	signal int

	sum int
}

func (d *Display) Write(text string) {
	if text == "noop" {
		d.tick()
	} else {
		var summand int
		_, err := fmt.Sscanf(text, "addx %d", &summand)
		if err != nil {
			log.Fatal(err)
		}
		d.tick()
		d.tick()
		d.signal += summand
	}
}

func (d *Display) tick() {
	d.cycle++
	if d.next == d.cycle {
		d.next += 40
		d.sum += d.cycle * d.signal
	}
}

func partTwo() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	d := &Drawer{display: os.Stdout}

	for scanner.Scan() {
		text := scanner.Text()
		d.Write(text)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type Drawer struct {
	cycle     int
	spritePos int
	display   io.Writer
}

func (d *Drawer) Write(text string) {
	if text == "noop" {
		d.tick()
	} else {
		var summand int
		_, err := fmt.Sscanf(text, "addx %d", &summand)
		if err != nil {
			log.Fatal(err)
		}
		d.tick()
		d.tick()
		d.spritePos += summand
	}
}

func (d *Drawer) tick() {
	pos := d.cycle % 40
	d.cycle++
	if pos >= d.spritePos && pos < d.spritePos+3 {
		fmt.Fprintf(d.display, "#")
	} else {
		fmt.Fprintf(d.display, ".")
	}
	if d.cycle%40 == 0 {
		fmt.Fprintf(d.display, "\n")
	}
}
