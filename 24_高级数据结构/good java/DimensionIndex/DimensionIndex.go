package main

import "fmt"

func main() {
	di := NewDimensionIndex(2, 3, 4)
	fmt.Println(di.Index(1, 2, 3))
	fmt.Println(di.Size())
}

// dimension: 是一个数组，例如 [2, 3, 4]
// index: 是一个数字，例如 8
type DimensionIndex struct {
	dimensions []int
	size       int
}

func NewDimensionIndex(dimension ...int) *DimensionIndex {
	size := 1
	for i := 0; i < len(dimension); i++ {
		size *= dimension[i]
	}
	return &DimensionIndex{
		dimensions: dimension,
		size:       size,
	}
}

func (di *DimensionIndex) Index(dimension ...int) int {
	res := dimension[0]
	for i := 1; i < len(dimension); i++ {
		res = res*di.dimensions[i] + dimension[i]
	}
	return res
}

func (di *DimensionIndex) Dimension(index int) []int {
	res := make([]int, len(di.dimensions))
	for i := len(di.dimensions) - 1; i >= 0; i-- {
		res[i] = index % di.dimensions[i]
		index = index / di.dimensions[i]
	}
	return res
}

func (di *DimensionIndex) IsValidDimension(dimension ...int) bool {
	for i := 0; i < len(dimension); i++ {
		if dimension[i] < 0 || dimension[i] >= di.dimensions[i] {
			return false
		}
	}
	return true
}

func (di *DimensionIndex) IndexOfSpecifiedDimension(index, d int) int {
	for i := len(di.dimensions) - 1; i >= 0; i-- {
		if i == d {
			return index % di.dimensions[i]
		}
		index = index / di.dimensions[i]
	}
	panic("Invalid dimension")
}

func (di *DimensionIndex) Size() int {
	return di.size
}
