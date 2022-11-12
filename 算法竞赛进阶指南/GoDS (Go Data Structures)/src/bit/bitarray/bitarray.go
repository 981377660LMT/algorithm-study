package main

import "fmt"

func main() {
	bit := NewBIT([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Println(bit.PreSum(10))
	bit.Add(10, 1)
	fmt.Println(bit.PreSum(10) - bit.PreSum(9))
	fmt.Println(bit.RangeSum(2, 10))
}

// !单点修改,区间查询,数组实现的树状数组
type BIT struct {
	tree []int
}

// 常数优化: dp O(n) 建树
// https://oi-wiki.org/ds/fenwick/#tricks
func NewBIT(nums []int) *BIT {
	tree := make([]int, len(nums)+1)
	for i := 1; i < len(tree); i++ {
		tree[i] += nums[i-1]
		if j := i + (i & -i); j < len(tree) {
			tree[j] += tree[i]
		}
	}
	return &BIT{tree}
}

// 位置 index 增加 delta
//  1<=i<=n
func (b *BIT) Add(index int, delta int) {
	for ; index < len(b.tree); index += index & -index {
		b.tree[index] += delta
	}
}

// 求前缀和preSum[index]
//  0<=i<=n
func (b *BIT) PreSum(index int) (res int) {
	for ; index > 0; index -= index & -index {
		res += b.tree[index]
	}
	return
}

func (b *BIT) RangeSum(left, right int) int {
	return b.PreSum(right) - b.PreSum(left-1)
}
