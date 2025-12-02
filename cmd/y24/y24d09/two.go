package main

import (
	"fmt"
	"log"

	"aoc/pkg/in"
	"aoc/pkg/lists"
)

// some kinda doubly linked list
type Node[T any] struct {
	data T
	next *Node[T]
	prev *Node[T]
}

func (n *Node[T]) Next() *Node[T] {
	return n.next
}

func (n *Node[T]) Previous() *Node[T] {
	return n.prev
}

func (n *Node[T]) Set(v T) {
	n.data = v
}

func (n *Node[T]) Get() T {
	return n.data
}

func (n *Node[T]) AddAfter(v T) {
	oldNext := n.next
	newNode := &Node[T]{
		data: v,
		next: oldNext,
		prev: n,
	}
	n.next = newNode
}

func partTwo() {
	file := in.MustOpenInputTxt(inputs)
	defer file.Close()

	blocks := readDisk(file)
	head, tail := asLinkedList(blocks)

	// printList(head)
	compact(head, tail)
	// printList(head)
	// printExample(head)

	// too high: 6401555874641
	sum := checksumBlocks(head)
	upperBound := 6401528557129
	lowerBound := 6138836566452
	if sum >= upperBound {
		log.Fatalf("too high %d >= %d", sum, upperBound)
	}

	if sum <= lowerBound {
		log.Fatalf("too low %d <= %d", sum, lowerBound)
	}

	fmt.Printf("part two: %d\n", checksumBlocks(head))
}

func compact(head, tail *Node[*block]) {
	start := head

	for start != nil {
		current := start.Get()
		if current.fid >= 0 {
			start = start.Next()
			continue
		}

		candidate := tail
		for candidate != nil && candidate != start {
			candidateValue := candidate.Get()
			if candidateValue.fid < 0 {
				// whitespace, try next
				candidate = candidate.Previous()
				continue
			}

			//if candidateValue.size == current.size {
			//	//swap positions
			//	start.Set(&block{
			//		fid:  candidateValue.fid,
			//		size: candidateValue.size,
			//	})
			//	candidate.Set(&block{
			//		fid:  -1,
			//		size: current.size,
			//	})
			//	break
			//}

			if candidateValue.size < current.size {
				// swap positions
				start.Set(&block{
					fid:  candidateValue.fid,
					size: candidateValue.size,
				})
				remaining := current.size - candidateValue.size
				if remaining > 0 {
					start.AddAfter(&block{
						fid:  -1,
						size: remaining,
					})
				}

				// old position can be all whitespace

				// TODO merge this! what if we need to adjust the tail???

				prev := candidate.Previous()
				next := candidate.Next().Next()

				prev.next = next
				next.prev = prev
				prev.Set(&block{
					fid:  -1,
					size: prev.data.size + candidateValue.size + next.data.size,
				})
				//candidate.Set(&block{
				//	fid:  -1,
				//	size: candidateValue.size,
				//})

				break
			}

			candidate = candidate.Previous()
		}

		start = start.Next()
	}
}

func printExample(head *Node[*block]) {
	next := head
	for next != nil {
		c := "."
		el := next.Get()
		if el.fid >= 0 {
			c = fmt.Sprintf("%d", el.fid)
		}

		for i := 0; i < el.size; i++ {
			fmt.Print(c)
		}

		next = next.Next()
	}
	fmt.Println()
}

func printList(head *Node[*block]) {
	next := head
	for next != nil {
		c := ". "
		el := next.Get()
		if el.fid >= 0 {
			c = fmt.Sprintf("%d(%d) ", el.fid, el.size)
		}

		for i := 0; i < el.size; i++ {
			fmt.Print(c)
		}

		next = next.Next()
	}
	fmt.Println()
}

func asLinkedList(blocks []*block) (*Node[*block], *Node[*block]) {
	head := &Node[*block]{
		data: blocks[0],
	}

	next := head
	for i := 1; i < len(blocks); i++ {
		next.AddAfter(blocks[i])
		next = next.Next()
	}
	tail := next
	return head, tail
}

func checksumBlocks(head *Node[*block]) int {
	sum := 0
	start := 0

	next := head
	for next != nil {
		v := next.Get()
		if v.fid < 0 {
			start += v.size
			next = next.Next()
			continue
		}

		for j := 0; j < v.size; j++ {
			sum += start * v.fid
			start++
		}
		next = next.Next()
	}
	return sum
}

func printBlocks(blocks lists.LinkedList[*block]) {
	blocks.ForEach(func(_ int, aBlock *block) {
		c := '.'
		if aBlock.fid >= 0 {
			c = rune(aBlock.fid + '0')
		}

		for i := 0; i < aBlock.size; i++ {
			fmt.Printf("%c", c)
		}
	})
	fmt.Println()
}
