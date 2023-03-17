// 链表实现的珂朵莉树 随机数据时O(nlogn)
// 不要指望什么题目都套珂朵莉树（虽然它能水过很多数据结构题），特别是在数据非随机的情况下不要使用。
// 如果题目让你求区间x次方和而在题目条件下你想不出巧算，那写一颗珂朵莉树还是很OK的。

// 什么时候用珂朵莉树
//  !珂朵莉树的核心操作在于推平(assign)一个区间。保证数据随机。

// https://www.luogu.com.cn/problem/CF896C
// https://www.luogu.com.cn/blog/endlesscheng/solution-cf896c
// n个数 m次操作 n,m<=1e5 数据随机
// 区间加
// 区间赋值
// 区间第k小
// 求区间幂次和

package main

import (
	"fmt"
	"sort"
)

// !下标从1开始,区间[left,right]为闭区间
func main() {
	odt := NewODT([]Value{12, 3, 5, 2, 7, 8, 9, 1, 4, 6})

	begin, end := odt.Prepare(2, 3)
	fmt.Println(odt.powSum(begin, end, 2, 1e9+7))

	begin, end = odt.Prepare(2, 7)
	odt.add(begin, end, 1)

	begin, end = odt.Prepare(2, 7)
	fmt.Println(odt.kth(begin, end, 1))
}

type Value = int

// 链表/数组能更简洁地维护分裂与合并操作
type ODTBlock struct {
	left, right int
	value       Value
}

type ODT []ODTBlock

func NewODT(nums []Value) ODT {
	n := len(nums)
	res := make(ODT, n)
	for i := range res {
		res[i] = ODTBlock{i, i, nums[i]}
	}
	return res
}

// !所有操作前需要预分裂 保证后续操作在 [left, right] 内部
func (t *ODT) Prepare(left, right int) (begin, end int) {
	begin = t.split(left)
	end = t.split(right + 1)
	return
}

// !区间[left,right]赋值操作
//  这里传入right是为了将begin块的右端点更新为right
func (t *ODT) Assign(begin, end, right int, value Value) {
	ot := *t
	ot[begin].right = right
	ot[begin].value = value
	if begin+1 < end {
		*t = append(ot[:begin+1], ot[end:]...)
	}
}

// 区间加
func (t ODT) add(begin, end int, val Value) {
	for i := begin; i < end; i++ {
		t[i].value += val
	}
}

// 区间第k小(1-based)
func (t ODT) kth(begin, end, k int) Value {
	blocks := append(t[:0:0], t[begin:end]...)
	sort.Slice(blocks, func(i, j int) bool { return blocks[i].value < blocks[j].value })
	k--
	for _, b := range blocks {
		if cnt := b.right - b.left + 1; k >= cnt {
			k -= cnt
		} else {
			return b.value
		}
	}
	panic(k)
}

// 区间幂和取模
func (t ODT) powSum(begin, end int, exp int, mod int) (res int) {
	for _, b := range t[begin:end] {
		res += (b.right - b.left + 1) * pow(b.value, exp, mod)
		res %= mod
	}
	return res % mod
}

// [l, r] => [l, mid-1] [mid, r]
//  return index of [mid, r]
//  return len(t) if not found
func (t *ODT) split(mid int) int {
	odt := *t
	for i, block := range odt {
		if block.left == mid {
			return i
		}
		if block.left < mid && mid <= block.right { // b.l <= mid-1
			*t = append(odt[:i+1], append(ODT{{mid, block.right, block.value}}, odt[i+1:]...)...)
			odt[i].right = mid - 1
			return i + 1
		}
	}
	return len(odt)
}

func pow(base int, exp int, mod int) int {
	base %= mod
	res := int(1) % mod
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}
