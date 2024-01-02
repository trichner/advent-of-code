package main

import (
	"cmp"
	"encoding/hex"
	"image/color"
	"image/png"
	"io"
	"slices"
	"strings"

	"aoc/pkg/iso3d"
	"aoc/pkg/util"
)

var colorPalette = []color.RGBA{
	hexDecode("#8931ef"),
	hexDecode("#f2ca19"),
	hexDecode("#ff00bd"),
	hexDecode("#0057e9"),
	hexDecode("#87e911"),
	hexDecode("#e11845"),
}

func hexDecode(s string) color.RGBA {
	s = strings.TrimPrefix(s, "#")
	if len(s) < 6 {
		panic("invalid hex color: " + s)
	}
	r := util.Must(hex.DecodeString(s[:2]))
	g := util.Must(hex.DecodeString(s[2:4]))
	b := util.Must(hex.DecodeString(s[4:6]))
	return color.RGBA{r[0], g[0], b[0], 0xff}
}

func RenderCuboids(w io.Writer, cubes []*iso3d.Cuboid) error {
	canvas := iso3d.NewCanvas(256, 1024*3)

	for i, cube := range cubes {
		canvas.AddCube(&iso3d.Cuboid{
			Position: cube.Position,
			Size:     cube.Size,
			Color:    colorPalette[i%len(colorPalette)],
		})
	}

	canvas.Draw()

	return png.Encode(w, canvas.AsImage())
}

func depthSort(cubes []*iso3d.Cuboid) {
	slices.SortFunc(cubes, depthCmp)
}

func depthCmp(a, b *iso3d.Cuboid) int {
	r := cmp.Compare(a.Position.X+a.Position.Y, b.Position.X+b.Position.Y)
	if r != 0 {
		return r
	}
	return cmp.Compare(a.Position.Z, b.Position.Z)
}
