package main

import (
	"bytes"
	"fmt"
)

func main() {
	slice02()
}

func slice01() {
	// Slices Overlapped
	path := []byte("AAAA/BBBBBBBBB")
	sepIndex := bytes.IndexByte(path, '/')

	// dir1 := path[:sepIndex:sepIndex]
	dir1 := path[:sepIndex]
	dir2 := path[sepIndex+1:]
	fmt.Println("dir1 =>", string(dir1)) //prints: dir1 => AAAA
	fmt.Println("dir2 =>", string(dir2)) //prints: dir2 => BBBBBBBBB

	dir1 = append(dir1, "suffix"...)

	dir1 = append(dir1, "suffix"...)
	fmt.Println("dir1 =>", string(dir1)) //prints: dir1 => AAAAsuffix
	fmt.Println("dir2 =>", string(dir2)) //prints: dir2 => uffixBBBB
	// Slices Overlapped
}

func slice02() {
	path := []byte("AAAA/BBBBBBBBB")
	seqIndex := bytes.IndexByte(path, '/')

	dir1 := path[:seqIndex:seqIndex] // limited Capacity  The appending cause a new buffer allocation
	dir2 := path[seqIndex+1:]

	fmt.Println("dir1 =>", string(dir1)) //prints: dir1 => AAAA fmt.Println("dir2 =>",string(dir2)) //prints: dir2 => BBBBBBBBB
	dir1 = append(dir1, "suffix"...)
	fmt.Println("dir1 =>", string(dir1)) //prints: dir1 => AAAAsuffix
	fmt.Println("dir2 =>", string(dir2)) //prints: dir2 => BBBBBBBBB
}
