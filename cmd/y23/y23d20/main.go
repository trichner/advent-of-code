package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strings"

	"aoc/pkg/in"
)

//go:embed *.txt
var inputs embed.FS

type Metered interface {
	Metrics() *Metrics
}
type ModuleRunner interface {
	Metered
	ID() string

	Pulse(sender ModuleRunner, level bool)

	AddInput(runner ModuleRunner)
	AddOutput(runner ModuleRunner)

	Inputs() []ModuleRunner
	Outputs() []ModuleRunner
}

type Module struct {
	Id      string
	outputs []ModuleRunner
	inputs  []ModuleRunner
}

func (m *Module) ID() string {
	return m.Id
}

func (m *Module) Metrics() *Metrics {
	panic("implement me")
}

func (m *Module) Pulse(sender ModuleRunner, level bool) {
	// nop
}

func (m *Module) pulseAll(sender ModuleRunner, level bool) {
	// nop
	for _, n := range m.outputs {
		n.Pulse(sender, level)
	}
}

func (m *Module) AddInput(runner ModuleRunner) {
	m.inputs = append(m.inputs, runner)
}

func (m *Module) AddOutput(runner ModuleRunner) {
	m.outputs = append(m.outputs, runner)
}

func (m *Module) Inputs() []ModuleRunner {
	return m.inputs
}

func (m *Module) Outputs() []ModuleRunner {
	return m.outputs
}

type Button struct {
	Module
}

func (b *Button) Pulse(sender ModuleRunner, level bool) {
	b.pulseAll(b, false)
}

type Output struct {
	Module
}

type FlipFlop struct {
	Module
	state bool
}

func (f *FlipFlop) Pulse(sender ModuleRunner, level bool) {
	if level {
		return
	}
	f.state = !f.state
	f.pulseAll(f, f.state)
}

type Broadcaster struct {
	Module
}

func (b *Broadcaster) Pulse(sender ModuleRunner, level bool) {
	b.pulseAll(b, level)
}

type Conjunction struct {
	Module
	inputstate map[string]bool
}

func (c *Conjunction) Pulse(sender ModuleRunner, level bool) {
	c.inputstate[sender.ID()] = level

	for _, input := range c.inputs {
		v := c.inputstate[input.ID()]
		if !v {
			c.pulseAll(c, true)
			return
		}
	}
	c.pulseAll(c, false)
}

type Metrics struct {
	HighPulses int
	LowPulses  int
}

type TracedModule struct {
	delegate ModuleRunner
	metrics  *Metrics
}

func (t *TracedModule) ID() string {
	return t.delegate.ID()
}

func (t *TracedModule) Pulse(sender ModuleRunner, level bool) {
	if level {
		t.metrics.HighPulses++
	} else {
		t.metrics.LowPulses++
	}
	t.delegate.Pulse(sender, level)
}

func (t *TracedModule) AddInput(runner ModuleRunner) {
	t.delegate.AddInput(runner)
}

func (t *TracedModule) AddOutput(runner ModuleRunner) {
	t.delegate.AddOutput(runner)
}

func (t *TracedModule) Inputs() []ModuleRunner {
	return t.delegate.Inputs()
}

func (t *TracedModule) Outputs() []ModuleRunner {
	return t.delegate.Outputs()
}

func (t *TracedModule) Metrics() *Metrics {
	return t.metrics
}

func main() {
	partOne()
	partTwo()
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	definitions := parse(file)

	runners := buildCircuit2(definitions)

	// input dependent loops :/
	loops := [][]string{
		{"vl", "kb"},
		{"fb", "jt"},
		{"kn", "ks"},
		{"ln", "sx"},
	}

	steps := 1
	for _, l := range loops {
		steps *= stepsToLowPulse(runners, l[0], l[1])
	}
	fmt.Printf("part two: %d\n", steps)

	if steps != 243081086866483 {
		panic("bad")
	}
}

func stepsToLowPulse(runners map[string]ModuleRunner, start, end string) int {
	startNode := runners[start]
	endNode := runners[end]
	metrics := endNode.Metrics()

	count := 0
	for {
		count++
		startNode.Pulse(nil, false)
		if metrics.LowPulses > 0 {
			return count
		}
	}
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	definitions := parse(file)

	runners, metrics := buildCircuit(definitions)

	button := &Button{Module: Module{Id: "button"}}
	start := runners["broadcaster"]
	for i := 0; i < 1000; i++ {
		start.Pulse(button, false)
	}

	sum := metrics.HighPulses * metrics.LowPulses
	fmt.Printf("part one: %d\n", sum)

	if 777666211 != sum {
		panic("bad!")
	}
}

func buildCircuit2(definitions []ModuleDefinition) map[string]ModuleRunner {
	runners := map[string]ModuleRunner{}

	// instantiate
	for _, definition := range definitions {
		var runner ModuleRunner
		switch definition.Type {
		case ModuleFlipFlop:
			runner = &FlipFlop{Module: Module{Id: definition.ID}}
		case ModuleBroadcaster:
			runner = &Broadcaster{Module: Module{Id: definition.ID}}
		case ModuleButton:
			runner = &Button{Module: Module{Id: definition.ID}}
		case ModuleConjunction:
			runner = &Conjunction{Module: Module{Id: definition.ID}, inputstate: make(map[string]bool)}
		}
		runners[runner.ID()] = &TracedModule{delegate: runner, metrics: &Metrics{}}
	}

	runners["rx"] = &TracedModule{delegate: &Output{Module: Module{Id: "rx"}}, metrics: &Metrics{}}

	// wire up
	for _, definition := range definitions {
		runner := runners[definition.ID]
		for _, n := range definition.Next {
			if len(n) == 0 {
				panic("wut?")
			}
			out, ok := runners[n]
			if !ok {
				panic("wut?")
				// out = &Output{id: n}
				// runners[n] = out
			}
			runner.AddOutput(out)
			out.AddInput(runner)
		}
	}

	return runners
}

func buildCircuit(definitions []ModuleDefinition) (map[string]ModuleRunner, *Metrics) {
	runners := map[string]ModuleRunner{}

	metrics := &Metrics{}

	// instantiate
	for _, definition := range definitions {
		var runner ModuleRunner
		switch definition.Type {
		case ModuleFlipFlop:
			runner = &FlipFlop{Module: Module{Id: definition.ID}}
		case ModuleBroadcaster:
			runner = &Broadcaster{Module: Module{Id: definition.ID}}
		case ModuleButton:
			runner = &Button{Module: Module{Id: definition.ID}}
		case ModuleConjunction:
			runner = &Conjunction{Module: Module{Id: definition.ID}, inputstate: make(map[string]bool)}
		}
		runner = &TracedModule{delegate: runner, metrics: metrics}
		runners[runner.ID()] = runner
	}

	// wire up
	for _, definition := range definitions {
		runner := runners[definition.ID]
		for _, n := range definition.Next {
			if len(n) == 0 {
				panic("wut?")
			}
			out, ok := runners[n]
			if !ok {
				out = &TracedModule{delegate: &Output{Module: Module{Id: n}}, metrics: metrics}
				runners[n] = out
			}
			runner.AddOutput(out)
			out.AddInput(runner)
		}
	}

	return runners, metrics
}

type ModuleType byte

const (
	ModuleFlipFlop    ModuleType = '%'
	ModuleBroadcaster            = 'X'
	ModuleButton                 = '#'
	ModuleConjunction            = '&'
)

type ModuleDefinition struct {
	ID   string
	Type ModuleType
	Next []string
}

func parse(r io.Reader) []ModuleDefinition {
	scanner := bufio.NewScanner(r)

	var definitions []ModuleDefinition
	for scanner.Scan() {
		text := scanner.Text()

		definitions = append(definitions, parseModule(text))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return definitions
}

func parseModule(s string) ModuleDefinition {
	splits := strings.Split(s, " -> ")

	nexts := strings.Split(splits[1], ", ")
	if len(nexts) == 0 {
		panic("bad: " + s)
	}

	module := splits[0]
	if len(module) == 0 {
		panic("bad: " + s)
	}
	if module == "broadcaster" {
		return ModuleDefinition{
			ID:   module,
			Type: ModuleBroadcaster,
			Next: nexts,
		}
	}

	t := module[0]
	switch t {
	case '%':
		return ModuleDefinition{
			ID:   module[1:],
			Type: ModuleFlipFlop,
			Next: nexts,
		}
	case '&':
		return ModuleDefinition{
			ID:   module[1:],
			Type: ModuleConjunction,
			Next: nexts,
		}
	}
	panic("unknown: " + s)
}
