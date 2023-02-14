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
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	bit := NewBinaryIndexedTree(n)
	for i := 0; i < q; i++ {
		var op, x, y int
		fmt.Fscan(in, &op, &x, &y)
		if op == 0 {
			bit.Apply(x-1, y)
		} else {
			fmt.Fprintln(out, bit.ProdRange(x-1, y))
		}
	}
}

type BinaryIndexedTree struct {
	n    int
	log  int
	data []int
}

// 長さ n の 0で初期化された配列で構築する.
func NewBinaryIndexedTree(n int) *BinaryIndexedTree {
	return &BinaryIndexedTree{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

// 配列で構築する.
func NewBinaryIndexedTreeFrom(arr []int) *BinaryIndexedTree {
	res := NewBinaryIndexedTree(len(arr))
	res.build(arr)
	return res
}

// 要素 i に値 v を加える.
func (b *BinaryIndexedTree) Apply(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r) の要素の総和を求める.
func (b *BinaryIndexedTree) Prod(r int) int {
	res := int(0)
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *BinaryIndexedTree) ProdRange(l, r int) int {
	return b.Prod(r) - b.Prod(l)
}

// 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
func (b *BinaryIndexedTree) LowerBound(x int) int {
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
func (b *BinaryIndexedTree) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

func (b *BinaryIndexedTree) build(arr []int) {
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
