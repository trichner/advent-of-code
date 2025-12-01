package main

import (
	"testing"

	"aoc/pkg/in"
)

func BenchmarkWalk(b *testing.B) {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	island := parse(file)

	b.Run("walk", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			walk(island, 0, 3)
		}
	})
	b.Run("walk_prio", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// walkPrio(island, 0, 3)
		}
	})
}
