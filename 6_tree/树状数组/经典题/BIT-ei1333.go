// https://ei1333.github.io/library/structure/others/binary-indexed-tree.hpp

// 使い方
// BinaryIndexedTree(sz): 長さ sz の 0で初期化された配列で構築する.
// BinaryIndexedTree(vs): 配列 vs で構築する.
// apply(i, x): 要素 i に値 x を加える.
// prod(r): [0, r) の要素の総和を求める.
// prod(l, r): [l, r) の要素の総和を求める.
// lower_bound(x): 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
// upper_bound(x): 区間[0,k]の総和がxを上回る最小のkを求める.数列が単調増加であることを要求する.

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

// test LowerBound and UpperBound
func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	bit := NewBitArrayFrom(nums)
	fmt.Println(bit.UpperBound(10), bit.Query(4))
}

type BitArray struct {
	n    int
	log  int
	data []int
}

func NewBitArray(n int) *BitArray {
	return &BitArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

func NewBitArrayFrom(arr []int) *BitArray {
	res := NewBitArray(len(arr))
	res.Build(arr)
	return res
}

func (b *BitArray) Build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}

func (b *BitArray) Add(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r)
func (b *BitArray) Query(r int) int {
	res := 0
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *BitArray) QueryRange(l, r int) int {
	return b.Query(r) - b.Query(l)
}

// 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) LowerBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 区間[0,k]の総和がxを上回る最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

func (b *BitArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}
