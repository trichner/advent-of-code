package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	partOne()
}

type INode struct {
	IsDir    bool
	Parent   *INode
	Children map[string]*INode
	Size     int
	Name     string
}

func newInode(parent *INode, name string, size int, isDir bool) *INode {
	return &INode{
		IsDir:    isDir,
		Parent:   parent,
		Children: make(map[string]*INode),
		Size:     size,
		Name:     name,
	}
}

func partOne() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	mustScan(scanner)

	root := newInode(nil, "", 0, true)
	pwd := root

	text := mustScan(scanner)
	for {

		text = scanner.Text()
		if text == "" {
			break
		}

		if strings.HasPrefix(text, "$ ls") {
			for scanner.Scan() {
				text := scanner.Text()
				if strings.HasPrefix(text, "$") {
					break
				} else if strings.HasPrefix(text, "dir") {
					// don't care
				} else {
					var size int
					var fname string
					_, err := fmt.Sscanf(text, "%d %s", &size, &fname)
					if err != nil {
						log.Fatal(err)
					}
					pwd.Children[fname] = newInode(pwd, fname, size, false)
				}
			}
		} else if strings.HasPrefix(text, "$ cd") {
			var dir string
			_, err := fmt.Sscanf(text, "$ cd %s", &dir)
			if err != nil {
				log.Fatal(err)
			}
			if dir == ".." {
				pwd = pwd.Parent
			} else {
				n, ok := pwd.Children[dir]
				if ok {
					pwd = n
				} else {
					child := newInode(pwd, dir, 0, true)
					pwd.Children[dir] = child
					pwd = child
				}
			}
			if !scanner.Scan() {
				break
			}
		} else {
			log.Fatalf("unexpected text: %q", text)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	Du(root)

	// part one
	fmt.Printf("%d\n", SumFolders(root))

	// part two
	space := 70000000
	used := root.Size
	currentFree := space - used
	mustFree := 30000000 - currentFree
	d := ListDirs(root)
	sort.Sort(sort.IntSlice(d))

	for _, v := range d {
		if v >= mustFree {
			fmt.Printf("%d\n", v)
			break
		}
	}
}

func SumFolders(n *INode) int {
	max := 100000
	total := 0
	if n.IsDir && n.Size <= max {
		total += n.Size
	}
	for _, v := range n.Children {
		total += SumFolders(v)
	}
	return total
}

func ListDirs(n *INode) []int {
	var l []int
	if n.IsDir {
		l = append(l, n.Size)
	}
	for _, v := range n.Children {
		l = append(l, ListDirs(v)...)
	}
	return l
}

func Du(n *INode) int {
	total := 0
	if !n.IsDir {
		total = n.Size
	}
	for _, v := range n.Children {
		total += Du(v)
	}
	if n.IsDir {
		n.Size = total
	}
	return total
}

func Print(root *INode) {
	printRec(root, "")
}

func printRec(n *INode, ident string) {
	fmt.Printf("%s%s (size=%d)\n", ident, n.Name, n.Size)
	for _, v := range n.Children {
		printRec(v, ident+" ")
	}
}

func mustScan(scanner *bufio.Scanner) string {
	if !scanner.Scan() {
		panic(scanner.Err())
	}
	return scanner.Text()
}
