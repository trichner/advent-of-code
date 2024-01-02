package iso3d

// https://mazebert.com/forum/news/isometric-depth-sorting--id775/
// http://shaunlebron.github.io/IsometricBlocks/
// https://stackoverflow.com/questions/892811/drawing-isometric-game-worlds

import (
	"cmp"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"slices"

	"aoc/pkg/vec"
)

var (
	Green = color.RGBA{0, 255, 0, 255}
	Red   = color.RGBA{255, 0, 0, 255}
	Blue  = color.RGBA{0, 0, 255, 255}
)

const (
	tileWidth      = 16
	tileWidthHalf  = tileWidth / 2
	tileHeight     = 8
	tileHeightHalf = tileHeight / 2
)

var unitFaceLeft = []vec.Vec3i{
	{},
	{0, 0, 1},
	{0, 1, 1},
	{0, 1, 0},
}

var unitFaceTop = []vec.Vec3i{
	{0, 0, 0},
	{1, 0, 0},
	{1, 1, 0},
	{0, 1, 0},
}

var unitFaceRight = []vec.Vec3i{
	{},
	{1, 0, 0},
	{1, 0, 1},
	{0, 0, 1},
}

type Face struct {
	Vertices []vec.Vec3i
	Color    color.Color
}

const (
	// order matters for depth sort, top must be first
	tileTop = iota
	tileLeft
	tileRight
)

type cuboidTile struct {
	Position vec.Vec3i
	Side     uint8
	Color    color.RGBA
}

func (c *cuboidTile) vertices() []vec.Vec3i {
	origin := c.Position
	switch c.Side {
	case tileLeft:
		return []vec.Vec3i{
			origin.Add(unitFaceLeft[0]),
			origin.Add(unitFaceLeft[1]),
			origin.Add(unitFaceLeft[2]),
			origin.Add(unitFaceLeft[3]),
		}
	case tileTop:
		return []vec.Vec3i{
			origin.Add(unitFaceTop[0]),
			origin.Add(unitFaceTop[1]),
			origin.Add(unitFaceTop[2]),
			origin.Add(unitFaceTop[3]),
		}
	case tileRight:
		return []vec.Vec3i{
			origin.Add(unitFaceRight[0]),
			origin.Add(unitFaceRight[1]),
			origin.Add(unitFaceRight[2]),
			origin.Add(unitFaceRight[3]),
		}
	}
	panic(fmt.Errorf("unknown tile side: %d", c.Side))
}

// Cuboid represents a rectangular cuboid in 3D space
//
// In isometric coordinates:
// .       z
// .    _.-`-._
// .   |`-_._-'|
// .   |   |   |
// . y `--_|_--' x
// .       ^ Position
type Cuboid struct {
	// Position defines the bottom center corner of the cuboid
	Position vec.Vec3i

	// Size defines the length of the edges in all three dimension
	Size vec.Vec3i

	// Color is the color of the cuboid
	Color color.RGBA
}

func (c *Cuboid) String() string {
	return fmt.Sprintf("{ Position: %s, Size: %s }", c.Position, c.Size)
}

func (c *Cuboid) tileFaces() []*cuboidTile {
	size := c.Size.Y*c.Size.Y + c.Size.X*c.Size.Z + c.Size.X*c.Size.Y
	faces := make([]*cuboidTile, 0, size)

	for y := 0; y < c.Size.Y; y++ {
		for z := 0; z < c.Size.Z; z++ {
			origin := c.Position.Add(vec.Vec3i{0, y, z})
			faces = append(faces, &cuboidTile{
				Position: origin,
				Side:     tileLeft,
				Color:    c.Color,
			})
		}
	}

	rightColor := mulColor(c.Color, 0.8)
	for x := 0; x < c.Size.X; x++ {
		for z := 0; z < c.Size.Z; z++ {
			origin := c.Position.Add(vec.Vec3i{x, 0, z})
			faces = append(faces, &cuboidTile{
				Position: origin,
				Side:     tileRight,
				Color:    rightColor,
			})
		}
	}

	topColor := mulColor(c.Color, 1.2)
	for x := 0; x < c.Size.X; x++ {
		for y := 0; y < c.Size.Y; y++ {
			origin := c.Position.Add(vec.Vec3i{x, y, c.Size.Z})
			faces = append(faces, &cuboidTile{
				Position: origin,
				Side:     tileTop,
				Color:    topColor,
			})
		}
	}

	return faces
}

func mulColor(c color.RGBA, f float32) color.RGBA {
	return color.RGBA{
		R: mulChannel(c.R, f),
		G: mulChannel(c.G, f),
		B: mulChannel(c.B, f),
		A: 255,
	}
}

func mulChannel(v uint8, f float32) uint8 {
	return clamp(float32(v) * f)
}

func clamp(v float32) uint8 {
	return uint8(max(0, min(v, 255)))
}

type Canvas struct {
	img           *image.RGBA
	width, height int

	faces []*cuboidTile
}

func NewCanvas(width, height int) *Canvas {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	c := &Canvas{
		img:    img,
		width:  width,
		height: height,
	}

	c.drawFloor()

	return c
}

func (c *Canvas) AsImage() image.Image {
	return c.img
}

func (c *Canvas) AddCube(cube *Cuboid) {
	faces := cube.tileFaces()

	c.faces = append(c.faces, faces...)
}

func (c *Canvas) Draw() {
	c.depthSortFaces()

	projectedFace := make([]vec.Vec2i, 4)
	for _, f := range c.faces {
		for i, vertex := range f.vertices() {
			projectedFace[i] = c.isoProject(vertex)
		}
		drawSquare(c.img, projectedFace, f.Color)
	}
}

func (c *Canvas) depthSortFaces() {
	slices.SortFunc(c.faces, depthCmp)
}

func depthCmp(a, b *cuboidTile) int {
	pa := a.Position
	pb := b.Position

	// draw bottom first
	r := cmp.Compare(pa.Z, pb.Z)
	if r != 0 {
		return r
	}

	// last diagonal first
	r = cmp.Compare(pb.X+pb.Y, pa.X+pa.Y)
	if r != 0 {
		return r
	}

	r = cmp.Compare(pa.X, pb.X)
	if r != 0 {
		return r
	}

	// same position
	if pa != pb {
		panic(fmt.Errorf("expected %s == %s", a.Position, b.Position))
	}

	if a.Side == b.Side {
		panic(fmt.Errorf("overlapping face: %s / %s", a, b))
	}

	return cmp.Compare(a.Side, b.Side)
}

func (c *Canvas) isoProject(p vec.Vec3i) vec.Vec2i {
	projected := project(p)
	return vec.Vec2i{projected.X + c.width/2, c.height - projected.Y}
}

func (c *Canvas) drawFloor() {
	length := 16

	for i := 0; i <= length; i++ {
		start := c.isoProject(vec.Vec3i{0, i, 0})
		end := c.isoProject(vec.Vec3i{length, i, 0})
		c.drawLine(start, end)

		start = c.isoProject(vec.Vec3i{i, 0, 0})
		end = c.isoProject(vec.Vec3i{i, length, 0})
		c.drawLine(start, end)
	}
}

func (c *Canvas) drawLine(start, end vec.Vec2i) {
	drawLine(c.img, start, end, color.Black)
}

func drawIsoCube(rgba *image.RGBA, cube *Cuboid, c color.Color) {
	faces := cube.tileFaces()

	projectedFace := make([]vec.Vec2i, 4)
	for _, f := range faces {
		for i, vertex := range f.vertices() {
			projectedFace[i] = project(vertex)
		}
		drawSquare(rgba, projectedFace, c)
	}
}

func project(p vec.Vec3i) vec.Vec2i {
	x := (p.X-p.Y)*tileWidthHalf - tileWidthHalf
	y := (p.X + p.Y) * tileHeightHalf
	y += p.Z * tileHeight
	return vec.Vec2i{x, y}
}

func drawCoordinate(rgba *image.RGBA, origin vec.Vec3i) {
	xUnity := project(vec.Vec3i{1, 0, 0})
	yUnity := project(vec.Vec3i{0, 1, 0})
	zUnity := project(vec.Vec3i{0, 0, 1})

	p := project(origin)
	rgba.Set(p.X, p.Y, color.Black)

	p = project(origin.Add(vec.Vec3i{1, 0, 0}))
	rgba.Set(xUnity.X, xUnity.Y, Red)

	p = project(origin.Add(vec.Vec3i{0, 1, 0}))
	rgba.Set(yUnity.X, yUnity.Y, Blue)

	p = project(origin.Add(vec.Vec3i{0, 0, 1}))
	rgba.Set(zUnity.X, zUnity.Y, Green)
}

func drawUnitCube(rgba *image.RGBA, pos vec.Vec3i) {
	cube := &Cuboid{
		Position: pos,
		Size:     vec.Vec3i{1, 1, 1},
	}
	faces := cube.tileFaces()

	projectedFace := make([]vec.Vec2i, 4)
	colors := []color.Color{Red, Blue, Green}
	for n, f := range faces {
		for i, vertex := range f.vertices() {
			projectedFace[i] = project(vertex)
		}
		drawSquare(rgba, projectedFace, colors[n])
	}
}

func drawSquare(img *image.RGBA, vertices []vec.Vec2i, c color.Color) {
	bb := vec.BoundingBox2i(vertices)

	for x := bb.From.X; x <= bb.To.X; x++ {
		for y := bb.From.Y; y <= bb.To.Y; y++ {
			p := vec.Vec2i{x, y}

			c1 := cross(vertices[0], vertices[1], p)
			if c1 < 0 {
				continue
			}
			c2 := cross(vertices[1], vertices[2], p)
			if c2 < 0 {
				continue
			}
			c3 := cross(vertices[2], vertices[3], p)
			if c3 < 0 {
				continue
			}
			c4 := cross(vertices[3], vertices[0], p)
			if c4 < 0 {
				continue
			}

			//if c1 == 0 || c2 == 0 || c3 == 0 || c4 == 0 {
			//	img.Set(x, y, color.Black)
			//} else {
			img.Set(x, y, c)
			//}
		}
	}
}

func isNeedleInside(a, b, needle vec.Vec2i) bool {
	return cross(a, b, needle) >= 0
}

func cross(a, b, needle vec.Vec2i) int {
	return (needle.X-a.X)*(b.Y-a.Y) - (needle.Y-a.Y)*(b.X-a.X)
}
