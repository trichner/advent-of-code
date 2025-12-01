package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Tile = byte

const (
	TileAir       Tile = '.'
	TileRoundRock Tile = 'O'
	TileCubeRock  Tile = '#'
)

func WeighTiles(t *Tiles) int {
	_, weightFactor := t.Size()
	sum := 0

	m, n := t.Size()
	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			v := t.Get(x, y)
			if v == TileRoundRock {
				sum += weightFactor
			}

		}
		weightFactor--
	}
	return sum
}

func TiltCycle(t *Tiles) *Tiles {
	for i := 0; i < 4; i++ {
		t = TiltNorth(t)
		t.Rotate(90)
	}
	return t
}

func TiltNorth(t *Tiles) *Tiles {
	m, _ := t.Size()
	for x := 0; x < m; x++ {
		t = tiltColumnNorth(t, x)
	}
	return t
}

func tiltColumnNorth(t *Tiles, x int) *Tiles {
	_, n := t.Size()

	for y := 0; y < n; y++ {
		prev := y
		if t.Get(x, y) != TileRoundRock {
			continue
		}

		for p := y - 1; p >= 0; p-- {
			if t.Get(x, p) == TileAir {
				// swap
				tmp := t.Get(x, prev)
				t.Set(x, prev, t.Get(x, p))
				t.Set(x, p, tmp)
			} else {
				break
			}
			prev = p
		}
	}
	return t
}

func NewTiles(t [][]Tile) *Tiles {
	m := len(t[0])
	n := len(t)
	backing := make([]Tile, m*n)
	for y, row := range t {
		for x, v := range row {
			backing[y*m+x] = v
		}
	}

	tiles := &Tiles{
		m:       m,
		n:       n,
		backing: backing,
	}
	tiles.mapPosition = tiles.getRot0
	return tiles
}

func (t *Tiles) getRot0(x, y int) (int, int) {
	return x, y
}

func (t *Tiles) getTransposed(x, y int) (int, int) {
	x, y = y, x
	return x, y
}

func (t *Tiles) getRot270(x, y int) (int, int) {
	x, y = y, x
	x = t.m - 1 - x
	return x, y
}

func (t *Tiles) getRot180(x, y int) (int, int) {
	x = t.m - 1 - x
	y = t.n - 1 - y
	return x, y
}

func (t *Tiles) getRot90(x, y int) (int, int) {
	x, y = y, x
	y = t.n - 1 - y
	return x, y
}

func (t *Tiles) Size() (int, int) {
	if t.rotation%180 == 0 {
		return t.m, t.n
	}
	return t.n, t.m
}

type Tiles struct {
	m, n    int
	backing []Tile

	rotation int

	mapPosition func(x, y int) (int, int)
}

func (t *Tiles) Hash() string {
	h := sha256.New()

	m, n := t.Size()
	h.Write([]byte{byte(m), byte(n), byte(t.rotation)})

	h.Write(t.backing)

	s := h.Sum(nil)
	return hex.EncodeToString(s)
}

func (t *Tiles) Rotate(deg int) {
	if deg%90 != 0 {
		panic("not supported")
	}
	t.rotation += deg + 360
	t.rotation = t.rotation % 360
	switch t.rotation {
	case 0:
		t.mapPosition = t.getRot0
		break
	case 90:
		t.mapPosition = t.getRot90
		break
	case 180:
		t.mapPosition = t.getRot180
		break
	case 270:
		t.mapPosition = t.getRot270
		break
	default:
		panic("wut?")
	}
}

func (t *Tiles) Get(x, y int) Tile {
	x, y = t.mapPosition(x, y)
	return t.backing[y*t.m+x]
}

func (t *Tiles) Set(x, y int, v Tile) {
	x, y = t.mapPosition(x, y)
	t.backing[y*t.m+x] = v
}

func (t *Tiles) String() string {
	var buf strings.Builder
	m, n := t.Size()
	for y := 0; y < n; y++ {
		buf.WriteString(fmt.Sprintf("%2d ", y))
		for x := 0; x < m; x++ {
			buf.WriteByte(byte(t.Get(x, y)))
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}
