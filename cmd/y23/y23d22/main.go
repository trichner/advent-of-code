package main

import (
	"bufio"
	"cmp"
	"embed"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"aoc/pkg/iso3d"
	"aoc/pkg/queue"
	"aoc/pkg/sets"
	"aoc/pkg/util"
	"aoc/pkg/vec"

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

	cuboids := parse(file)
	cuboids = stackCuboids(cuboids)
	sum := maxFallingBlocks(cuboids)

	fmt.Printf("part two: %d\n", sum)

	if sum != 71002 {
		panic("bad")
	}
}

func maxFallingBlocks(cuboids []*iso3d.Cuboid) int {
	integrities := calculateIntegrities(cuboids)

	sum := 0
	for _, cube := range cuboids {
		sum += countFallingPerBlock(integrities, cube)
	}
	return sum
}

func countFallingPerBlock(integrities map[*iso3d.Cuboid]*Integrity, cube *iso3d.Cuboid) int {
	integrity := integrities[cube]
	fallingBlocks := sets.New[*iso3d.Cuboid]()

	disintegrated := sets.New[*iso3d.Cuboid]()
	disintegrated.Put(cube)

	var toCheck queue.Queue[*iso3d.Cuboid]
	for _, other := range integrity.Supports.Keys() {
		toCheck.Push(other)
	}

	for {
		other, ok := toCheck.Pop()
		if !ok {
			break
		}
		otherIntegrity := integrities[other]
		isSupported := false
		for _, supportingCube := range otherIntegrity.SupportedBy.Keys() {
			if disintegrated.Has(supportingCube) {
				continue
			}
			isSupported = true
			break
		}

		if !isSupported {
			fallingBlocks.Put(other)
			for _, s := range otherIntegrity.Supports.Keys() {
				disintegrated.Put(other)
				toCheck.Push(s)
			}
		}
	}
	return fallingBlocks.Size()
}

func partOne() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	cuboids := parse(file)

	visualized := util.Must(os.Create("viz_before.png"))
	defer visualized.Close()
	RenderCuboids(visualized, cuboids)

	stacked := stackCuboids(cuboids)

	visualizedAfter := util.Must(os.Create("viz_after.png"))
	defer visualizedAfter.Close()
	RenderCuboids(visualizedAfter, cuboids)

	sum := countNonStructuralBlocks(stacked)

	fmt.Printf("part one: %d\n", sum)

	if sum != 505 {
		panic("bad")
	}
}

func countNonStructuralBlocks(cuboids []*iso3d.Cuboid) int {
	integrities := calculateIntegrities(cuboids)

	sum := 0
	for _, integrity := range integrities {
		isStructural := false
		for _, other := range integrity.Supports.Keys() {
			otherIntegrity := integrities[other]
			if otherIntegrity.SupportedBy.Size() == 1 {
				isStructural = true
				break
			}
		}
		if !isStructural {
			sum++
		}
	}
	return sum
}

type Integrity struct {
	SupportedBy sets.Set[*iso3d.Cuboid]
	Supports    sets.Set[*iso3d.Cuboid]
}

func calculateIntegrities(cuboids []*iso3d.Cuboid) map[*iso3d.Cuboid]*Integrity {
	// better be safe than sorry
	sortZ(cuboids)

	integrities := map[*iso3d.Cuboid]*Integrity{}

	heightMap := map[vec.Vec2i]*iso3d.Cuboid{}

	for _, cube := range cuboids {

		integrity := &Integrity{
			SupportedBy: sets.New[*iso3d.Cuboid](),
			Supports:    sets.New[*iso3d.Cuboid](),
		}
		for x := 0; x < cube.Size.X; x++ {
			for y := 0; y < cube.Size.Y; y++ {
				supportedBy, ok := heightMap[vec.Vec2i{cube.Position.X + x, cube.Position.Y + y}]
				if !ok {
					// on ground if cube.Position.Z == 0
					continue
				}
				if supportedBy.Position.Z+supportedBy.Size.Z < cube.Position.Z {
					// in air
					continue
				}
				if supportedBy.Position.Z+supportedBy.Size.Z > cube.Position.Z {
					panic("wtf?")
				}
				integrity.SupportedBy.Put(supportedBy)
				other := integrities[supportedBy]
				other.Supports.Put(cube)
			}
		}
		integrities[cube] = integrity

		// update height map
		for x := 0; x < cube.Size.X; x++ {
			for y := 0; y < cube.Size.Y; y++ {
				heightMap[vec.Vec2i{cube.Position.X + x, cube.Position.Y + y}] = cube
			}
		}
	}

	return integrities
}

func stackCuboids(cuboids []*iso3d.Cuboid) []*iso3d.Cuboid {
	sortZ(cuboids)
	heightMap := map[vec.Vec2i]int{}

	var stacked []*iso3d.Cuboid
	for _, cube := range cuboids {

		minHeight := 0
		for x := 0; x < cube.Size.X; x++ {
			for y := 0; y < cube.Size.Y; y++ {
				h := heightMap[vec.Vec2i{cube.Position.X + x, cube.Position.Y + y}]
				minHeight = max(h, minHeight)
			}
		}
		stacked = append(stacked, &iso3d.Cuboid{
			Position: vec.Vec3i{cube.Position.X, cube.Position.Y, minHeight},
			Size:     cube.Size,
		})

		for x := 0; x < cube.Size.X; x++ {
			for y := 0; y < cube.Size.Y; y++ {
				heightMap[vec.Vec2i{cube.Position.X + x, cube.Position.Y + y}] = minHeight + cube.Size.Z
			}
		}
	}
	return stacked
}

func sortZ(cuboids []*iso3d.Cuboid) {
	slices.SortFunc(cuboids, func(a, b *iso3d.Cuboid) int {
		return cmp.Compare(a.Position.Z, b.Position.Z)
	})
}

func parse(r io.Reader) []*iso3d.Cuboid {
	scanner := bufio.NewScanner(r)

	var blocks []*iso3d.Cuboid
	for scanner.Scan() {
		text := scanner.Text()

		blocks = append(blocks, parseBlock(text))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return blocks
}

func parseBlock(s string) *iso3d.Cuboid {
	splits := strings.Split(s, "~")

	start := parseVec3i(splits[0])
	end := parseVec3i(splits[1])
	size := end.Sub(start).Add(vec.Vec3i{1, 1, 1})
	return &iso3d.Cuboid{
		Position: start,
		Size:     size,
	}
}

func parseVec3i(s string) vec.Vec3i {
	splits := strings.Split(s, ",")

	x := util.Must(strconv.Atoi(splits[0]))
	y := util.Must(strconv.Atoi(splits[1]))
	z := util.Must(strconv.Atoi(splits[2]))
	return vec.Vec3i{x, y, z}
}
