package main

import (
	"cmp"
	"fmt"
	"slices"

	"aoc/pkg/vec"
)

func partTwo(loop []vec.Vec2i) {
	segments := segmentLoop(loop)
	printSegments(segments)
	area := calculateArea(segments)

	fmt.Printf("part two: %d\n", area)
	if area != 417 {
		panic("bad result")
	}
}

func boundingBox(segments []*Segment) vec.AABB {
	var points []vec.Vec2i
	for _, s := range segments {
		points = append(points, s.To)
		points = append(points, s.From)
	}
	return vec.BoundingBox2i(points)
}

func calculateArea(segments []*Segment) int {
	bb := boundingBox(segments)

	area := 0
	for y := 0; y < bb.To.Y; y++ {
		area += calculateAreaPerRow(bb, segments, y)
	}
	return area
}

func calculateAreaPerRow(bb vec.AABB, segments []*Segment, y int) int {
	area := 0
	opened := 0
	for x := 0; x < bb.To.X; x++ {
		p := vec.Vec2i{x, y}
		isect := stabsAnySegment(segments, p)
		if isect == IntersectsInside {
			opened++
		} else if isect == IntersectsEnd || isect == IntersectsStart {
			i2 := IntersectsNone
			for i2 != IntersectsEnd && i2 != IntersectsStart {
				x++
				p2 := vec.Vec2i{x, y}
				i2 = stabsAnySegment(segments, p2)
			}
			if i2 != isect {
				opened++
			}
		} else if isect == IntersectsNone {
			if opened%2 == 1 {
				area++
			}
		} else {
			panic("Wut?")
		}

	}
	return area
}

func newSegment(from, to vec.Vec2i) *Segment {
	c := vec.Compare2i(from, to)
	if c > 0 {
		from, to = to, from
	}
	return &Segment{
		From: from,
		To:   to,
	}
}

type Segment struct {
	From, To vec.Vec2i
}

type Intersection int

const (
	IntersectsNone Intersection = iota
	IntersectsStart
	IntersectsEnd
	IntersectsInside
)

func stabsAnySegment(segments []*Segment, p vec.Vec2i) Intersection {
	idx := slices.IndexFunc(segments, func(segment *Segment) bool {
		return stabSegment(segment, p) != IntersectsNone
	})
	if idx < 0 {
		return IntersectsNone
	}
	seg := segments[idx]
	return stabSegment(seg, p)
}

func stabSegment(seg *Segment, p vec.Vec2i) Intersection {
	if p.X != seg.From.X {
		return IntersectsNone
	}
	if p == seg.From {
		return IntersectsStart
	}
	if p == seg.To {
		return IntersectsEnd
	}
	if seg.From.Y < p.Y && seg.To.Y > p.Y {
		return IntersectsInside
	}
	return IntersectsNone
}

func (s *Segment) String() string {
	return fmt.Sprintf("(%v, %v)", s.From, s.To)
}

func segmentLoop(loop []vec.Vec2i) []*Segment {
	var segments []*Segment

	tail := loop[len(loop)-1]
	prev := tail
	for i := 0; i < len(loop); i++ {
		head := loop[i]
		v := head.Sub(prev)
		if v.X == 0 {
			// we proceed
		} else if tail != prev {
			segments = append(segments, newSegment(tail, prev))
			prev = head
			tail = head
		} else {
			prev = head
			tail = head
		}
		prev = head
	}

	slices.SortFunc(segments, func(a, b *Segment) int {
		d := cmp.Compare(a.From.X, b.From.X)
		if d != 0 {
			return d
		}
		return cmp.Compare(a.From.Y, b.From.Y)
	})

	return segments
}
