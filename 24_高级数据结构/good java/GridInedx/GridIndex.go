package main

import (
	"fmt"
	"math/bits"
)

func main() {
	G := NewGridIndex(3, 4)
	fmt.Println(G.Encode(1, 2))
	fmt.Println(G.DecodeRow(6))
	fmt.Println(G.DecodeCol(6))
}

// 二维矩阵的坐标编码解码.
type GridIndex struct {
	mask int32
	bit  int32
}

func NewGridIndex(row, col int32) *GridIndex {
	bit := CeilLog32(uint32(col))
	mask := (1 << bit) - 1
	return &GridIndex{mask: int32(mask), bit: int32(bit)}
}

func (g *GridIndex) Encode(x, y int32) int32 {
	return (x << g.bit) | y
}

func (g *GridIndex) DecodeRow(id int32) int32 {
	return id >> g.bit
}

func (g *GridIndex) DecodeCol(id int32) int32 {
	return id & g.mask
}

func CeilLog32(x uint32) int {
	if x <= 0 {
		return 0
	}
	return 32 - bits.LeadingZeros32(x-1)
}
