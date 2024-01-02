package iso3d

import (
	"image"
	"image/color"

	"aoc/pkg/vec"
)

func drawLine(rgba *image.RGBA, start, end vec.Vec2i, c color.Color) {
	// https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm#Algorithm_for_integer_arithmetic
	if abs(end.Y-start.Y) < abs(end.X-start.X) {
		if start.X > end.X {
			start, end = end, start
		}
		drawLineLow(rgba, start, end, c)
	} else {
		if start.Y > end.Y {
			start, end = end, start
		}
		drawLineHigh(rgba, start, end, c)
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func drawLineLow(rgba *image.RGBA, start, end vec.Vec2i, c color.Color) {
	diff := end.Sub(start)

	yi := 1
	if diff.Y < 0 {
		yi = -1
		diff.Y = -diff.Y
	}

	D := 2*diff.Y - diff.X
	y := start.Y

	for x := start.X; x <= end.X; x++ {
		rgba.Set(x, y, c)
		if D > 0 {
			y += yi
			D += 2 * (diff.Y - diff.X)
		} else {
			D += 2 * diff.Y
		}
	}
}

func drawLineHigh(rgba *image.RGBA, start, end vec.Vec2i, c color.Color) {
	diff := end.Sub(start)

	xi := 1
	if diff.X < 0 {
		xi = -1
		diff.X = -diff.X
	}

	D := 2*diff.X - diff.Y
	x := start.X

	for y := start.Y; y <= end.Y; y++ {
		rgba.Set(x, y, c)
		if D > 0 {
			x += xi
			D += 2 * (diff.X - diff.Y)
		} else {
			D += 2 * diff.X
		}
	}
}
