package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	n := 100000
	root := NewPersistentArraySqrt(make([]int, n))
	time1 := time.Now()
	for i := 0; i < n; i++ {
		root = root.Set(i, i)
		root.Get(i)
	}
	fmt.Println(time.Since(time1))
}

type SnapshotArray struct {
	pa  *PersistentArraySqrt
	git []*PersistentArraySqrt
}

func Constructor(length int) SnapshotArray {
	return SnapshotArray{pa: NewPersistentArraySqrt(make([]int, length))}
}

func (this *SnapshotArray) Set(index int, val int) {
	this.pa = this.pa.Set(index, val)
}

func (this *SnapshotArray) Snap() int {
	this.git = append(this.git, this.pa)
	return len(this.git) - 1
}

func (this *SnapshotArray) Get(index int, snap_id int) int {
	return this.git[snap_id].Get(index)
}

/**
 * Your SnapshotArray object will be instantiated and called as such:
 * obj := Constructor(length);
 * obj.Set(index,val);
 * param_2 := obj.Snap();
 * param_3 := obj.Get(index,snap_id);
 */
type T = int
type PersistentArraySqrt struct {
	arr     []T
	opIndex []int
	opValue []T
	opLen   int
}

func NewPersistentArraySqrt(arr []T) *PersistentArraySqrt {
	sqrt := 2 * (int(math.Sqrt(float64(len(arr)))) + 1)
	return &PersistentArraySqrt{arr: arr, opIndex: make([]int, 0, sqrt), opValue: make([]T, 0, sqrt)}
}

func (sa *PersistentArraySqrt) Get(i int) T {
	for j := sa.opLen - 1; j >= 0; j-- {
		if sa.opIndex[j] == i {
			return sa.opValue[j]
		}
	}
	return sa.arr[i]
}

func (sa *PersistentArraySqrt) Set(i int, v T) *PersistentArraySqrt {
	sa.opIndex = append(sa.opIndex, i)
	sa.opValue = append(sa.opValue, v)
	n := len(sa.arr)
	if tmp := len(sa.opIndex); tmp*tmp <= 4*n {
		return &PersistentArraySqrt{arr: sa.arr, opIndex: sa.opIndex, opValue: sa.opValue, opLen: tmp}
	}
	newArr := make([]T, n)
	copy(newArr, sa.arr)
	for i := range sa.opIndex {
		newArr[sa.opIndex[i]] = sa.opValue[i]
	}
	return NewPersistentArraySqrt(newArr)
}
