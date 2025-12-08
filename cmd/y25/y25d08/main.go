package main

import (
	"bufio"
	"cmp"
	"embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"aoc/pkg/vec"

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

func partTwo() {
	positions := readJunctions()

	var junctions []*junction
	for _, p := range positions {
		j := &junction{pos: p}
		junctions = append(junctions, j)
	}

	var edges []*edge
	for i, a := range junctions {
		for j := i + 1; j < len(junctions); j++ {
			n := junctions[j]
			e := &edge{
				distance: a.pos.Sub(n.pos).Abs(),
				from:     a,
				to:       n,
			}
			edges = append(edges, e)
		}
	}

	slices.SortFunc(edges, func(a, b *edge) int {
		return cmp.Compare(a.distance, b.distance)
	})

	// wire them up
	circuitId := 1
	sum := 0
	for _, e := range edges {

		if e.to.circuit != nil && e.from.circuit != nil {
			_, ok := e.to.circuit.nodes[e.from.pos]
			if ok {
				continue
			}
		}

		if e.to.circuit == nil && e.from.circuit == nil {
			// none is in a circuit yet
			nc := make(map[vec.Vec3i]*junction)
			nc[e.from.pos] = e.from
			nc[e.to.pos] = e.to

			c := &circuit{id: circuitId, nodes: nc}
			circuitId++

			e.from.circuit = c
			e.to.circuit = c
		} else if e.to.circuit == nil {
			// 'from' is in a circuit
			e.from.circuit.nodes[e.to.pos] = e.to
			e.to.circuit = e.from.circuit
			if len(e.from.circuit.nodes) == len(junctions) {
				sum = e.from.pos.X * e.to.pos.X
				break
			}
		} else if e.from.circuit == nil {
			// 'to' is in a circuit
			e.to.circuit.nodes[e.from.pos] = e.from
			e.from.circuit = e.to.circuit
			if len(e.to.circuit.nodes) == len(junctions) {
				sum = e.from.pos.X * e.to.pos.X
				break
			}
		} else {
			// merge two circuits
			first := e.from.circuit
			for k, v := range e.to.circuit.nodes {
				first.nodes[k] = v
				v.circuit = first
			}

			if len(first.nodes) == len(junctions) {
				sum = e.from.pos.X * e.to.pos.X
				break
			}
		}
	}

	// 3767453340
	fmt.Printf("part two: %d\n", sum)
}

type junction struct {
	pos     vec.Vec3i
	circuit *circuit
}

type edge struct {
	distance float32
	from     *junction
	to       *junction
}

type circuit struct {
	id    int
	nodes map[vec.Vec3i]*junction
}

func partOne() {
	positions := readJunctions()

	var junctions []*junction
	for _, p := range positions {
		j := &junction{pos: p}
		junctions = append(junctions, j)
	}

	var edges []*edge
	for i, a := range junctions {
		for j := i + 1; j < len(junctions); j++ {
			n := junctions[j]
			e := &edge{
				distance: a.pos.Sub(n.pos).Abs(),
				from:     a,
				to:       n,
			}
			edges = append(edges, e)
		}
	}

	slices.SortFunc(edges, func(a, b *edge) int {
		return cmp.Compare(a.distance, b.distance)
	})

	// wire them up
	cables := 1000
	circuitId := 1
	for _, e := range edges {
		if cables <= 0 {
			break
		}
		cables--

		if e.to.circuit != nil && e.from.circuit != nil {
			_, ok := e.to.circuit.nodes[e.from.pos]
			if ok {
				continue
			}
		}

		if e.to.circuit == nil && e.from.circuit == nil {
			// none is in a circuit yet
			nc := make(map[vec.Vec3i]*junction)
			nc[e.from.pos] = e.from
			nc[e.to.pos] = e.to

			c := &circuit{id: circuitId, nodes: nc}
			circuitId++

			e.from.circuit = c
			e.to.circuit = c
		} else if e.to.circuit == nil {
			// 'from' is in a circuit
			e.from.circuit.nodes[e.to.pos] = e.to
			e.to.circuit = e.from.circuit
		} else if e.from.circuit == nil {
			// 'to' is in a circuit
			e.to.circuit.nodes[e.from.pos] = e.from
			e.from.circuit = e.to.circuit
		} else {
			// merge two circuits
			first := e.from.circuit
			for k, v := range e.to.circuit.nodes {
				first.nodes[k] = v
				v.circuit = first
			}
		}
	}

	circuitMap := map[int]int{}
	for _, j := range junctions {
		if j.circuit == nil {
			continue
		}
		circuitMap[j.circuit.id] = len(j.circuit.nodes)
	}

	var lengths []int
	for k, v := range circuitMap {
		lengths = append(lengths, v)
		fmt.Printf("circuit %d: %d nodes\n", k, v)
	}

	slices.Sort(lengths)
	slices.Reverse(lengths)

	sum := 1
	for i := 0; i < 3; i++ {
		sum *= lengths[i]
	}
	// 67488
	fmt.Printf("part one: %d\n", sum)
}

func readJunctions() []vec.Vec3i {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var junctions []vec.Vec3i
	for scanner.Scan() {
		text := scanner.Text()
		parts := strings.Split(text, ",")

		x, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		y, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		z, _ := strconv.Atoi(strings.TrimSpace(parts[2]))

		junctions = append(junctions, vec.Vec3i{X: x, Y: y, Z: z})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	slices.SortFunc(junctions, func(a, b vec.Vec3i) int {
		return cmp.Compare(a.Abs(), b.Abs())
	})

	return junctions
}
