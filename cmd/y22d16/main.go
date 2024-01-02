package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"math"
	"math/bits"
	"os"
	"slices"
	"strings"
	"time"
)

//go:embed input.txt
var input string

type Tunnel struct {
	Destination string
	Cost        int
}

type Valve struct {
	ID       string
	FlowRate int
	Tunnels  []*Tunnel
}

type TunnelInt struct {
	Destination uint8
	Cost        int
}

type ValveInt struct {
	ID       uint8
	FlowRate int
	Tunnels  []*TunnelInt
}

func (v *Valve) String() string {
	return fmt.Sprintf("{ ID: %q, Tunnels: %v}", v.ID, v.Tunnels)
}

type BitSet uint16

func (b *BitSet) Set(e uint8) BitSet {
	if e > 15 {
		panic(fmt.Errorf("too big: %d", e))
	}
	return (*b) | (1 << e)
}

func (b *BitSet) Has(e uint8) bool {
	return ((*b) & (1 << e)) != 0
}

func (b *BitSet) Key() uint16 {
	return uint16(*b)
}

func (b *BitSet) Len() int {
	return bits.OnesCount16(uint16(*b))
}

func (b *BitSet) String() string {
	var buf strings.Builder
	buf.WriteString("{")
	for i := uint8(0); i < 32; i++ {
		if b.Has(i) {
			fmt.Fprintf(&buf, " %d", i)
		}
	}
	buf.WriteString(" }")
	return buf.String()
}

type LinkedSet struct {
	size     int
	Value    string
	Previous *LinkedSet
}

func (l *LinkedSet) Key() string {
	if l == nil {
		return ""
	}
	node := l
	var s []string
	for {
		if node == nil {
			break
		}
		s = append(s, node.Value)
		node = node.Previous
	}
	slices.Sort(s)
	return strings.Join(s, "-")
}

func (l *LinkedSet) Put(e string) *LinkedSet {
	return &LinkedSet{
		Value:    e,
		Previous: l,
		size:     l.Len() + 1,
	}
}

func (l *LinkedSet) Len() int {
	if l == nil {
		return 0
	}
	return l.size
}

func (l *LinkedSet) Has(e string) bool {
	node := l
	for {
		if node == nil {
			return false
		}
		if node.Value == e {
			return true
		}
		node = node.Previous
	}
}

func main() {
	valves := parseInput(strings.NewReader(input))

	graph := prepareGraph(valves)

	flow := solve(graph)

	fmt.Printf("%d\n", flow)
}

func prepareGraph(valves map[string]*Valve) []ValveInt {
	dumpParsedGraph(valves, "full.dot")

	removeZeroFlows(valves)
	dumpParsedGraph(valves, "nozero.dot")

	mapped := mapGraph(valves)
	dumpMappedGraph(mapped, "mapped.dot")
	return mapped
}

func solve(graph []ValveInt) int {
	count = 0
	best = 0
	startTime = time.Now()

	solutions := make([][2]int, int(math.Pow(2, 16+4)))
	var opened BitSet
	return walkTheTunnels(solutions, graph, &graph[0], opened.Set(0), 0, 0, 30)
}

func dumpParsedGraph(valves map[string]*Valve, fname string) {
	var buf strings.Builder
	buf.WriteString("strict graph { \n")
	for _, v := range valves {
		buf.WriteString(fmt.Sprintf("  %s [label=\"%s (%d)\"]\n", v.ID, v.ID, v.FlowRate))
		for _, t := range v.Tunnels {
			buf.WriteString(fmt.Sprintf("  %s -- %s [label=%d]\n", v.ID, t.Destination, t.Cost))
		}
	}
	buf.WriteString("}\n")

	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(buf.String())
}

func dumpMappedGraph(valves []ValveInt, fname string) {
	var buf strings.Builder
	buf.WriteString("strict graph { \n")
	for _, v := range valves {
		buf.WriteString(fmt.Sprintf("  %d [label=\"%d (%d)\"]\n", v.ID, v.ID, v.FlowRate))
		for _, t := range v.Tunnels {
			buf.WriteString(fmt.Sprintf("  %d -- %d [label=%d]\n", v.ID, t.Destination, t.Cost))
		}
	}
	buf.WriteString("}\n")

	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(buf.String())
}

var (
	count     int
	best      int
	startTime time.Time
)

const updateFrequency = 4_194_304

// TODO
func keyIntWithElephant(opened BitSet, p1 *ValveInt, p2 *ValveInt) uint32 {
	id := uint32(p1.ID) | (uint32(p2.ID) << 4)
	if p1.ID > p2.ID {
		id = uint32(p2.ID) | (uint32(p1.ID) << 4)
	}
	k := uint32(opened.Key()) | uint32(id)<<16
	if (k & 0xFF000000) != 0 {
		panic(fmt.Errorf("invalid bits: %0x", k))
	}
	return k
}

func walkTheTunnels(solutions [][2]int, valves []ValveInt, current *ValveInt, opened BitSet, currentFlow int, currentTotal int, minutesRemaining int) int {
	if current == nil || minutesRemaining < 0 {
		panic("wtf????")
	}

	if minutesRemaining == 0 {
		count++
		best = max(best, currentTotal)
		if count%updateFrequency == 0 {
			elapsedSeconds := int(time.Now().Sub(startTime) / time.Second)
			elapsedSeconds = max(1, elapsedSeconds)
			log.Printf("c: %d, t: %d, r: %dk/s\n", count, best, count/1000/elapsedSeconds)
		}
		return currentTotal
	}

	k := keyInt(opened, current)
	v := solutions[k]
	if v[0] >= minutesRemaining && v[1] >= currentTotal {
		// we can't do better
		return 0
	}
	v[0] = minutesRemaining
	v[1] = currentTotal

	// idle around
	totalReleased := currentTotal + currentFlow*minutesRemaining

	// is there anything left to open?
	if opened.Len() == len(valves) {
		count++
		best = max(best, totalReleased)
		if count%updateFrequency == 0 {
			log.Printf("c: %d, t: %d\n", count, best)
		}
		return totalReleased
	}

	// guess whether we should skip the current node or not, makes us converge on the solution faster
	doSkip := false
	for _, t := range current.Tunnels {
		if opened.Has(t.Destination) || minutesRemaining <= t.Cost {
			continue
		}
		target := valves[t.Destination]

		// going there directly
		skipOpening := (target.FlowRate+currentFlow)*(minutesRemaining-t.Cost) + currentFlow*t.Cost

		// going there directly
		doOpen := (target.FlowRate+currentFlow+current.FlowRate)*(minutesRemaining-t.Cost-1) + currentFlow + (currentFlow+current.FlowRate)*(t.Cost+1)
		if skipOpening > doOpen {
			doSkip = true
		}
	}

	if !doSkip && !opened.Has(current.ID) {
		// open the valve if It's not already open
		newOpened := opened.Set(current.ID)
		released := walkTheTunnels(solutions, valves, current, newOpened, currentFlow+current.FlowRate, currentTotal+currentFlow, minutesRemaining-1)
		totalReleased = max(totalReleased, released)
	}

	// go down all connected tunnels
	for _, t := range current.Tunnels {
		if minutesRemaining-t.Cost <= 0 {
			continue
		}
		released := walkTheTunnels(solutions, valves, &valves[t.Destination], opened, currentFlow, currentTotal+t.Cost*currentFlow, minutesRemaining-t.Cost)
		totalReleased = max(totalReleased, released)
	}

	if doSkip && !opened.Has(current.ID) {
		// open the valve if It's not already open
		newOpened := opened.Set(current.ID)
		released := walkTheTunnels(solutions, valves, current, newOpened, currentFlow+current.FlowRate, currentTotal+currentFlow, minutesRemaining-1)
		totalReleased = max(totalReleased, released)
	}

	return totalReleased
}

func keyInt(opened BitSet, p *ValveInt) uint32 {
	k := uint32(opened.Key()) | uint32(p.ID)<<16
	if (k & 0xFFF00000) != 0 {
		panic(fmt.Errorf("invalid bits: %0x", k))
	}
	return k
}

func parseInput(r io.Reader) map[string]*Valve {
	scanner := bufio.NewScanner(r)

	var valves []*Valve

	for scanner.Scan() {
		text := scanner.Text()

		splits := strings.Split(text, ";")
		var flow int
		var id string
		_, err := fmt.Sscanf(splits[0], "Valve %s has flow rate=%d", &id, &flow)
		if err != nil {
			log.Fatal(err)
		}

		s := strings.TrimPrefix(splits[1], " tunnels lead to valves ")
		s = strings.TrimPrefix(s, " tunnel leads to valve ")
		next := strings.Split(s, ", ")
		var tunnels []*Tunnel
		for _, s := range next {
			tunnels = append(tunnels, &Tunnel{Destination: s, Cost: 1})
		}

		valves = append(valves, &Valve{
			ID:       id,
			FlowRate: flow,
			Tunnels:  tunnels,
		})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	valveMap := map[string]*Valve{}
	for _, v := range valves {
		valveMap[v.ID] = v
	}

	return valveMap
}

func mapGraph(valves map[string]*Valve) []ValveInt {
	idSeq := uint8(0)
	idToNum := map[string]uint8{"AA": idSeq}
	idSeq++

	mapped := make([]ValveInt, len(valves))

	for _, v := range valves {
		mappedId, ok := idToNum[v.ID]
		if !ok {
			mappedId = idSeq
			idSeq++
			idToNum[v.ID] = mappedId
		}

		var mappedTunnels []*TunnelInt
		for _, t := range v.Tunnels {
			tId, ok := idToNum[t.Destination]
			if !ok {
				tId = idSeq
				idSeq++
				idToNum[t.Destination] = tId
			}
			mappedTunnels = append(mappedTunnels, &TunnelInt{
				Destination: tId,
				Cost:        t.Cost,
			})
		}

		mapped[mappedId] = ValveInt{
			ID:       mappedId,
			FlowRate: v.FlowRate,
			Tunnels:  mappedTunnels,
		}
	}

	return mapped
}

func removeZeroFlows(valves map[string]*Valve) {
	var valveList []*Valve
	for _, v := range valves {
		valveList = append(valveList, v)
	}
	for _, v := range valveList {
		if v.FlowRate == 0 && v.ID != "AA" {
			removeValve(valves, v)
		}
	}
}

func removeValve(valves map[string]*Valve, v *Valve) {
	for i, t1 := range v.Tunnels {
		for j := i + 1; j < len(v.Tunnels); j++ {
			t2 := v.Tunnels[j]
			if t1 == nil || t2 == nil {
				panic("wtf?")
			}
			if t1.Destination == t2.Destination {
				continue
			}

			cost := t1.Cost + t2.Cost

			v1 := valves[t1.Destination]
			v2 := valves[t2.Destination]

			// do we already have a connection?
			idx1 := slices.IndexFunc(v1.Tunnels, matchDestination(t2.Destination))
			idx2 := slices.IndexFunc(v2.Tunnels, matchDestination(t1.Destination))

			if idx1 >= 0 {
				existingCost := v1.Tunnels[idx1].Cost
				newCost := min(cost, existingCost)
				v1.Tunnels[idx1].Cost = newCost
				v2.Tunnels[idx2].Cost = newCost
			} else {
				v1.Tunnels = append(v1.Tunnels, &Tunnel{Destination: v2.ID, Cost: cost})
				v2.Tunnels = append(v2.Tunnels, &Tunnel{Destination: v1.ID, Cost: cost})
			}
		}
	}
	for _, t := range v.Tunnels {
		other := valves[t.Destination]
		idx := slices.IndexFunc(other.Tunnels, func(t *Tunnel) bool { return t.Destination == v.ID })
		if idx >= 0 {
			newTunnels := slices.Delete(other.Tunnels, idx, idx+1)
			other.Tunnels = newTunnels
		}
	}
	delete(valves, v.ID)
}

func matchDestination(s string) func(t *Tunnel) bool {
	return func(t *Tunnel) bool { return t.Destination == s }
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
