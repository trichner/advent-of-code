package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"runtime"
	"strings"

	"aoc/pkg/sets"
	"aoc/pkg/vec"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

type Part byte

type Contraption struct {
	Parts       [][]Part
	Light       [][]sets.Set[vec.Vec2i]
	BoundingBox vec.AABB
}

func NewContraption(parts [][]Part) *Contraption {
	size := vec.Vec2i{len(parts[0]), len(parts)}

	lights := makeLights(size)

	bb := vec.AABB{From: vec.Vec2i{}, To: size.Sub(vec.Vec2i{X: 1, Y: 1})}
	return &Contraption{
		Parts:       parts,
		Light:       lights,
		BoundingBox: bb,
	}
}

func CopyContraption(c *Contraption) *Contraption {
	size := c.Size()
	lights := makeLights(size)

	return &Contraption{
		Parts:       c.Parts,
		Light:       lights,
		BoundingBox: c.BoundingBox,
	}
}

func makeLights(size vec.Vec2i) [][]sets.Set[vec.Vec2i] {
	lights := make([][]sets.Set[vec.Vec2i], size.Y)
	for i := range lights {
		row := make([]sets.Set[vec.Vec2i], size.X)
		for j := range row {
			row[j] = make(sets.Set[vec.Vec2i])
		}
		lights[i] = row
	}
	return lights
}

func (c *Contraption) Size() vec.Vec2i {
	return vec.Vec2i{X: len(c.Parts[0]), Y: len(c.Parts)}
}

func (c *Contraption) GetPart(p vec.Vec2i) Part {
	if c.BoundingBox.Contains(p) {
		return c.Parts[p.Y][p.X]
	}
	return PartVoid
}

func (c *Contraption) AddLight(p vec.Vec2i, heading vec.Vec2i) bool {
	if c.BoundingBox.Contains(p) {
		set := c.Light[p.Y][p.X]
		if set.Has(heading) {
			return false
		}
		set.Put(heading)
		return true
	}
	return false
}

func (c *Contraption) LightCount() int {
	sum := 0
	for _, row := range c.Light {
		for _, l := range row {
			if len(l) > 0 {
				sum++
			}
		}
	}
	return sum
}

const (
	PartVoid Part = iota
	PartAir
	PartMirrorUp
	PartMirrorDown
	PartHSplitter
	PartVSplitter
)

func main() {
	partOne()
	partTwo()
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	parts := parse(file)

	contraption := NewContraption(parts)

	size := contraption.Size()

	tasks := make(chan []vec.Vec2i, size.X*2+size.Y*2+1)

	results := make(chan int, 256)

	running := 0
	for i := 0; i < runtime.NumCPU()*2; i++ {
		running++
		go func() {
			for {
				t := <-tasks
				if t == nil {
					tasks <- nil
					results <- -1
					return
				}

				c := CopyContraption(contraption)
				traceBeam(c, t[0], t[1])
				results <- c.LightCount()
			}
		}()
	}

	for x := 0; x < size.X; x++ {
		tasks <- []vec.Vec2i{{X: x}, HeadDown}
		tasks <- []vec.Vec2i{{x, size.Y - 1}, HeadUp}
	}
	for y := 0; y < size.Y; y++ {
		tasks <- []vec.Vec2i{{0, y}, HeadRight}
		tasks <- []vec.Vec2i{{size.X - 1, y}, HeadLeft}
	}

	tasks <- nil

	candidate := 0

	for running > 0 {
		r := <-results
		if r < 0 {
			running--
			continue
		}
		candidate = max(candidate, r)
	}
	close(tasks)
	close(results)

	fmt.Printf("part two: %d\n", candidate)
	if candidate != 8314 {
		panic("bad")
	}
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	parts := parse(file)

	contraption := NewContraption(parts)

	traceBeam(contraption, vec.Vec2i{0, 0}, HeadRight)

	sum := contraption.LightCount()
	fmt.Printf("part one: %d\n", sum)

	if sum != 8112 {
		panic("bad")
	}
}

func traceBeam(c *Contraption, position, heading vec.Vec2i) {
	ok := c.AddLight(position, heading)
	if !ok {
		return
	}

	part := c.GetPart(position)
	headings := getHeadings(part, heading)
	for _, h := range headings {
		next := position.Add(h)
		traceBeam(c, next, h)
	}
}

type key struct {
	P Part
	H vec.Vec2i
}

var (
	HeadUp    = vec.Vec2i{0, -1}
	HeadDown  = vec.Vec2i{0, 1}
	HeadRight = vec.Vec2i{1, 0}
	HeadLeft  = vec.Vec2i{-1, 0}
)

var headingLut = map[key][]vec.Vec2i{
	{PartAir, HeadUp}:    {HeadUp},
	{PartAir, HeadDown}:  {HeadDown},
	{PartAir, HeadLeft}:  {HeadLeft},
	{PartAir, HeadRight}: {HeadRight},

	{PartHSplitter, HeadUp}:    {HeadLeft, HeadRight},
	{PartHSplitter, HeadDown}:  {HeadLeft, HeadRight},
	{PartHSplitter, HeadLeft}:  {HeadLeft},
	{PartHSplitter, HeadRight}: {HeadRight},

	{PartVSplitter, HeadUp}:    {HeadUp},
	{PartVSplitter, HeadDown}:  {HeadDown},
	{PartVSplitter, HeadLeft}:  {HeadUp, HeadDown},
	{PartVSplitter, HeadRight}: {HeadUp, HeadDown},

	{PartMirrorUp, HeadUp}:    {HeadRight},
	{PartMirrorUp, HeadDown}:  {HeadLeft},
	{PartMirrorUp, HeadLeft}:  {HeadDown},
	{PartMirrorUp, HeadRight}: {HeadUp},

	{PartMirrorDown, HeadUp}:    {HeadLeft},
	{PartMirrorDown, HeadDown}:  {HeadRight},
	{PartMirrorDown, HeadLeft}:  {HeadUp},
	{PartMirrorDown, HeadRight}: {HeadDown},

	{PartVoid, HeadUp}:    {},
	{PartVoid, HeadDown}:  {},
	{PartVoid, HeadLeft}:  {},
	{PartVoid, HeadRight}: {},
}

func getHeadings(part Part, heading vec.Vec2i) []vec.Vec2i {
	headings, ok := headingLut[key{P: part, H: heading}]
	if !ok {
		panic("wut?")
	}
	return headings
}

func parse(r io.Reader) [][]Part {
	scanner := bufio.NewScanner(r)

	var parts [][]Part

	for scanner.Scan() {
		text := scanner.Text()
		parts = append(parts, parseLine(text))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return parts
}

func parseLine(s string) []Part {
	lut := map[byte]Part{'.': PartAir, '/': PartMirrorUp, '\\': PartMirrorDown, '|': PartVSplitter, '-': PartHSplitter}
	parts := make([]Part, len(s))
	s = strings.TrimSpace(s)
	for i, b := range []byte(s) {
		parts[i] = lut[b]
	}
	return parts
}
