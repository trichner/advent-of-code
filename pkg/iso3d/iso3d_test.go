package iso3d

import (
	"encoding/hex"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"strings"
	"testing"

	"aoc/pkg/util"
	"aoc/pkg/vec"
)

func TestRender(t *testing.T) {
	size := 1024
	// center := vec.Vec2i{size / 2, size / 2}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	drawCoordinate(img, vec.Vec3i{})

	//drawUnitCube(img, vec.Vec3i{1, 1, 3})
	//drawSquare(img, projected, Red)
	//
	cube := &Cuboid{
		Position: vec.Vec3i{0, 0, 0},
		Size:     vec.Vec3i{1, 1, 1},
	}
	drawIsoCube(img, cube, Red)

	// img.Set(center.X, center.Y, Green)

	f := util.Must(os.Create("rendered.png"))
	defer f.Close()

	err := png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

func TestCanvas_AsImage(t *testing.T) {
	canvas := NewCanvas(1024, 1024)

	var cube Cuboid

	cube = Cuboid{
		Position: vec.Vec3i{2, 0, 0},
		Size:     vec.Vec3i{2, 4, 8},
		Color:    color.RGBA{127, 127, 0, 255},
	}
	canvas.AddCube(&cube)

	cube = Cuboid{
		Position: vec.Vec3i{0, 1, 0},
		Size:     vec.Vec3i{2, 4, 8},
		Color:    Green,
	}
	canvas.AddCube(&cube)

	cube = Cuboid{
		Position: vec.Vec3i{},
		Size:     vec.Vec3i{1, 1, 1},
		Color:    Red,
	}
	canvas.AddCube(&cube)

	canvas.Draw()

	f := util.Must(os.Create("canvas.png"))
	defer f.Close()

	err := png.Encode(f, canvas.AsImage())
	if err != nil {
		panic(err)
	}
}

func TestCanvas_ThreeBoxProblem(t *testing.T) {
	canvas := NewCanvas(300, 300)

	cubes := []*Cuboid{
		{
			Position: vec.Vec3i{0, 1, 0},
			Size:     vec.Vec3i{1, 1, 2},
			Color:    colorPalette[0],
		},
		{
			Position: vec.Vec3i{0, 0, 0},
			Size:     vec.Vec3i{2, 1, 1},
			Color:    colorPalette[1],
		},
		{
			Position: vec.Vec3i{1, 0, 1},
			Size:     vec.Vec3i{1, 2, 1},
			Color:    colorPalette[2],
		},
	}

	for _, cuboid := range cubes {
		canvas.AddCube(cuboid)
	}

	canvas.Draw()

	f := util.Must(os.Create("three_blocks.png"))
	defer f.Close()

	err := png.Encode(f, canvas.AsImage())
	if err != nil {
		panic(err)
	}
}

func TestCanvas_DrawLots(t *testing.T) {
	canvas := NewCanvas(300, 300)

	rnd := rand.New(rand.NewSource(1337))

	var cubes []*Cuboid
	for i := 0; i < 100; i++ {
		cubes = append(cubes, &Cuboid{
			Position: vec.Vec3i{rnd.Intn(13), rnd.Intn(13), rnd.Intn(13)},
			Size:     vec.Vec3i{rnd.Intn(3) + 1, rnd.Intn(3) + 1, rnd.Intn(3) + 1},
			Color:    colorPalette[i%len(colorPalette)],
		})
	}

	for _, cuboid := range cubes {
		canvas.AddCube(cuboid)
	}

	canvas.Draw()

	f := util.Must(os.Create("lots.png"))
	defer f.Close()

	err := png.Encode(f, canvas.AsImage())
	if err != nil {
		panic(err)
	}
}

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
