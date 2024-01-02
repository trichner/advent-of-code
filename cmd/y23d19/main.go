package main

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"regexp"
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

	workflowsList, _ := parse(file)

	workflows := buildWorkflowMap(workflowsList)

	r := ValueRange{
		X: []int{1, 4000},
		M: []int{1, 4000},
		A: []int{1, 4000},
		S: []int{1, 4000},
	}
	combinations := treeEval(workflows, "in", r)

	fmt.Printf("part two: %d\n", combinations)

	if combinations != 124615747767410 {
		panic("bad result")
	}
}

type ValueRange struct {
	X []int
	M []int
	A []int
	S []int
}

func treeEval(workflows map[string]*Workflow, id string, valueRange ValueRange) int {
	if id == "A" {
		return scoreRange(valueRange)
	}
	if id == "R" {
		return 0
	}

	w := workflows[id]

	var branch ValueRange

	combinations := 0
	for _, r := range w.Rules {

		if r.Type == RuleConditionalNext {
			branch, valueRange = nextRange(r, valueRange)
			combinations += treeEval(workflows, r.Destination, branch)
			continue
		}

		if r.Type == RuleNext {
			combinations += treeEval(workflows, r.Destination, valueRange)
			break
		}

		if r.Type == RuleReject {
			break
		}
		if r.Type == RuleAccept {
			combinations += scoreRange(valueRange)
			break
		}
	}

	return combinations
}

func nextRange(r Rule, v ValueRange) (ValueRange, ValueRange) {
	v1 := v
	v2 := v

	switch r.Property {
	case 'x':
		a, b := splitRange(v.X, r.Literal, r.Operator)
		v1.X = a
		v2.X = b
	case 'm':
		a, b := splitRange(v.M, r.Literal, r.Operator)
		v1.M = a
		v2.M = b
	case 'a':
		a, b := splitRange(v.A, r.Literal, r.Operator)
		v1.A = a
		v2.A = b
	case 's':
		a, b := splitRange(v.S, r.Literal, r.Operator)
		v1.S = a
		v2.S = b
	}
	return v1, v2
}

func splitRange(a []int, split int, op byte) ([]int, []int) {
	if op == '<' {
		return splitRangeLt(a, split)
	}
	if op == '>' {
		return splitRangeGt(a, split)
	}
	panic("wut?")
}

func splitRangeLt(a []int, split int) ([]int, []int) {
	var v1, v2 []int
	if a[0] <= split-1 {
		v1 = []int{a[0], split - 1}
	}

	if split <= a[1] {
		v2 = []int{split, a[1]}
	}

	return v1, v2
}

func splitRangeGt(a []int, split int) ([]int, []int) {
	var v1, v2 []int
	if split+1 <= a[1] {
		v1 = []int{split + 1, a[1]}
	}

	if a[0] <= split {
		v2 = []int{a[0], split}
	}
	return v1, v2
}

func scoreRange(v ValueRange) int {
	return (v.X[1] - v.X[0] + 1) * (v.M[1] - v.M[0] + 1) * (v.A[1] - v.A[0] + 1) * (v.S[1] - v.S[0] + 1)
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	workflowsList, parts := parse(file)

	workflows := buildWorkflowMap(workflowsList)

	var accepted []Part
	for _, p := range parts {

		w := workflows["in"]

	workflowchain:
		for {
		rulechain:
			for _, rule := range w.Rules {
				result, next := rule.Eval(p)
				switch result {
				case ResultNoMatch:
					continue
				case ResultNext:
					w = workflows[next]
					break rulechain
				case ResultAccepted:
					accepted = append(accepted, p)
					break workflowchain
				case ResultRejected:
					break workflowchain
				}
			}
		}
	}

	sum := score(accepted)
	fmt.Printf("part one: %d\n", sum)

	if sum != 353553 {
		panic("bad result")
	}
}

func score(parts []Part) int {
	sum := 0
	for _, p := range parts {
		sum += p.X + p.M + p.A + p.S
	}
	return sum
}

func buildWorkflowMap(w []*Workflow) map[string]*Workflow {
	indexed := make(map[string]*Workflow, len(w))

	for _, workflow := range w {
		indexed[workflow.ID] = workflow
	}
	return indexed
}

type Part struct {
	X, M, A, S int
}

func (p *Part) String() string {
	return fmt.Sprintf("{x=%d,m=%d,a=%d,s=%d}", p.X, p.M, p.A, p.S)
}

type RuleType byte

const (
	RuleAccept RuleType = iota
	RuleReject
	RuleNext
	RuleConditionalNext
)

type Result byte

const (
	ResultNoMatch Result = iota
	ResultNext
	ResultAccepted
	ResultRejected
)

type Rule struct {
	Eval        RuleEval
	Destination string
	Type        RuleType

	Property byte
	Operator byte
	Literal  int
}

type RuleEval func(p Part) (r Result, next string)

type Workflow struct {
	ID    string
	Rules []Rule
}

func parse(r io.Reader) ([]*Workflow, []Part) {
	scanner := bufio.NewScanner(r)

	var workflows []*Workflow
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}
		workflows = append(workflows, parseWorkflow(text))
	}

	var parts []Part
	for scanner.Scan() {
		text := scanner.Text()

		parts = append(parts, parsePart(text))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return workflows, parts
}

var workflowPattern = regexp.MustCompile("([a-z]+)\\{([^}]*)}")

func parseWorkflow(s string) *Workflow {
	// px{a<2006:qkq,m>2090:A,rfg}
	groups := workflowPattern.FindStringSubmatch(s)
	if len(groups) == 0 {
		panic(":( : " + s)
	}

	id := groups[1]
	rawRules := groups[2]
	rules := parseRules(rawRules)

	return &Workflow{
		ID:    id,
		Rules: rules,
	}
}

func newRuleExpression(s string) Rule {
	splits := strings.Split(s, ":")
	property := splits[0][0]
	operator := splits[0][1]
	literal := util.Must(strconv.Atoi(splits[0][2:]))

	dst := splits[1]

	return Rule{
		Destination: dst,
		Type:        RuleConditionalNext,
		Eval:        buildEval(property, operator, dst, literal),

		Operator: operator,
		Literal:  literal,
		Property: property,
	}
}

func buildEval(property, operator byte, destination string, literal int) RuleEval {
	g := getter(property)
	c := comparator(operator)

	switch destination {
	case "A":
		return func(p Part) (r Result, next string) {
			if c(g(p), literal) {
				return ResultAccepted, ""
			}
			return ResultNoMatch, ""
		}
	case "R":
		return func(p Part) (r Result, next string) {
			if c(g(p), literal) {
				return ResultRejected, ""
			}
			return ResultNoMatch, ""
		}
	}

	return func(p Part) (r Result, next string) {
		if c(g(p), literal) {
			return ResultNext, destination
		}
		return ResultNoMatch, ""
	}
}

func parseRules(s string) []Rule {
	ruleParts := strings.Split(s, ",")
	var rules []Rule
	for _, part := range ruleParts {
		if part == "A" {
			rules = append(rules, Rule{
				Eval:        Accept,
				Destination: "A",
				Type:        RuleAccept,
			})
		} else if part == "R" {
			rules = append(rules, Rule{
				Eval:        Reject,
				Destination: "R",
				Type:        RuleReject,
			})
		} else if strings.ContainsRune(part, ':') {
			rules = append(rules, newRuleExpression(part))
		} else {
			dst := part
			rules = append(rules, Rule{
				Eval: func(p Part) (r Result, next string) {
					return ResultNext, dst
				},
				Destination: dst,
				Type:        RuleNext,
			})
		}
	}
	return rules
}

var Accept RuleEval = func(p Part) (r Result, next string) {
	return ResultAccepted, ""
}

var Reject RuleEval = func(p Part) (r Result, next string) {
	return ResultRejected, ""
}

func getter(p byte) func(p Part) int {
	switch p {
	case 'x':
		return func(p Part) int {
			return p.X
		}
	case 'm':
		return func(p Part) int {
			return p.M
		}
	case 'a':
		return func(p Part) int {
			return p.A
		}
	case 's':
		return func(p Part) int {
			return p.S
		}
	}
	panic("wut?")
}

func comparator(c byte) func(a, b int) bool {
	switch c {
	case '<':
		return func(a, b int) bool {
			return a < b
		}
	case '>':
		return func(a, b int) bool {
			return a > b
		}
	}
	panic("wut?")
}

var partPattern = regexp.MustCompile("\\{x=([0-9]+),m=([0-9]+),a=([0-9]+),s=([0-9]+)}")

func parsePart(s string) Part {
	groups := partPattern.FindStringSubmatch(s)
	if len(groups) == 0 {
		panic("bad: " + s)
	}

	return Part{
		X: util.Must(strconv.Atoi(groups[1])),
		M: util.Must(strconv.Atoi(groups[2])),
		A: util.Must(strconv.Atoi(groups[3])),
		S: util.Must(strconv.Atoi(groups[4])),
	}
}
