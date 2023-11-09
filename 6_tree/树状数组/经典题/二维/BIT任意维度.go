package main

import "fmt"

func main() {
	bit := NewBITMultiDimension(10, 10, 10, 10)
	bit.Add([]int{1, 1, 1, 1}, 1)
	bit.Add([]int{1, 1, 1, 1}, 2)
	fmt.Println(bit.QueryPrefix([]int{2, 3, 2, 2}))
}

type BITMultiDimension struct {
	_dim  []int
	_data []int
}

func NewBITMultiDimension(dimension ...int) *BITMultiDimension {
	n := 1
	for _, v := range dimension {
		n *= v
	}
	return &BITMultiDimension{_dim: dimension, _data: make([]int, n)}
}

// 0<=indices[i]<dimension[i]
func (bm *BITMultiDimension) Add(indices []int, x int) {
	bm._addRec(indices, 0, 0, x)
}

// 0<=indices[i]<dimension[i]
func (bm *BITMultiDimension) QueryPrefix(indices []int) int {
	return bm._queryRec(indices, 0, 0)
}

// 0<=a[i]<=b[i]<=dimension[i]
func (bm *BITMultiDimension) QueryRange(a, b []int) int {
	t := make([]int, len(a))
	return bm._queryRangeRec(0, a, b, t)
}

func (bm *BITMultiDimension) _addRec(indices []int, k, t, x int) {
	d := bm._dim[k]
	t *= d
	if k+1 == len(bm._dim) {
		for i := indices[k]; i < d; i |= i + 1 {
			bm._data[t+i] += x
		}
	} else {
		for i := indices[k]; i < d; i |= i + 1 {
			bm._addRec(indices, k+1, t+i, x)
		}
	}
}

func (bm *BITMultiDimension) _queryRec(indices []int, k, t int) int {
	d := bm._dim[k]
	t *= d
	res := 0
	if k+1 == len(bm._dim) {
		for i := indices[k]; i > 0; i -= i & -i {
			res += bm._data[t+i-1]
		}
	} else {
		for i := indices[k]; i > 0; i -= i & -i {
			res += bm._queryRec(indices, k+1, t+i-1)
		}
	}
	return res
}

func (bm *BITMultiDimension) _queryRangeRec(d int, a, b, t []int) int {
	if d == len(bm._dim) {
		return bm._queryRec(t, 0, 0)
	}
	res := 0
	t[d] = b[d]
	res += bm._queryRangeRec(d+1, a, b, t)
	t[d] = a[d]
	res -= bm._queryRangeRec(d+1, a, b, t)
	return res
}
