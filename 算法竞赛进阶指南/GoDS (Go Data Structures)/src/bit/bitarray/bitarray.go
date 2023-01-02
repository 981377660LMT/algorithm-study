package main

import "fmt"

func main() {
	bit := NewBITArray(10)
	bit.Build([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Println(bit.Query(10))
	bit.Add(10, 1)
	fmt.Println(bit.Query(9) - bit.Query(8))
	fmt.Println(bit.QueryRange(2, 9))
}

// !单点修改,区间查询,数组实现的树状数组
type BITArray struct {
	n    int
	tree []int
}

func NewBITArray(n int) *BITArray {
	return &BITArray{n: n, tree: make([]int, n+1)}
}

// 常数优化: dp O(n) 建树
// https://oi-wiki.org/ds/fenwick/#tricks
func (b *BITArray) Build(nums []int) {
	for i := 1; i < len(b.tree); i++ {
		b.tree[i] += nums[i-1]
		if j := i + (i & -i); j < len(b.tree) {
			b.tree[j] += b.tree[i]
		}
	}
}

// 位置 index 增加 delta
//  1<=i<=n
func (b *BITArray) Add(index int, delta int) {
	for ; index < len(b.tree); index += index & -index {
		b.tree[index] += delta
	}
}

// 求前缀和
//  1<=i<=n
func (b *BITArray) Query(index int) (res int) {
	if index > b.n {
		index = b.n
	}
	for ; index > 0; index -= index & -index {
		res += b.tree[index]
	}
	return
}

// 1<=left<=right<=n
func (b *BITArray) QueryRange(left, right int) int {
	return b.Query(right) - b.Query(left-1)
}

func (b *BITArray) Len() int {
	return b.n
}
