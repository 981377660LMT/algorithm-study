// 链表实现的珂朵莉树 随机数据时O(nlogn)
// 不要指望什么题目都套珂朵莉树（虽然它能水过很多数据结构题），特别是在数据非随机的情况下不要使用。
// 如果题目让你求区间x次方和而在题目条件下你想不出巧算，那写一颗珂朵莉树还是很OK的。

// 什么时候用珂朵莉树
//  !珂朵莉树实质上是一种可以维护区间上的分裂与合并的数据结构，但要求数据是随机的，
//   !或者有大量的随机合并操作，这样才能保证维护的区间个数是一个很小的值。
//   !在若干次`随机合并`后，区间个数会骤降至一个稳定的范围

// https://www.luogu.com.cn/problem/CF896C
// https://www.luogu.com.cn/blog/endlesscheng/solution-cf896c
// n个数 m次操作 n,m<=1e5 数据随机
// 区间加
// 区间赋值
// 区间第k小
// 求区间幂次和

// API:
//  !Set(left, right, value) 区间合并+赋值
//  !EnumerateRange(left, right, func(block *ODTBlock)) 区间分裂+遍历区间内的块

package main

import (
	"fmt"
	"runtime/debug"
	"sort"
)

func init() {
	debug.SetGCPercent(-1)
}

// !下标从0开始,闭区间[0,n-1]
func demo() {
	n := 10
	odt := NewODT([]Value{12, 3, 5, 2, 7, 8, 9, 1, 4, 6})
	fmt.Println(odt)
	odt.add(2, 3, 1)
	fmt.Println(odt)
	fmt.Println(odt.kth(0, 3, 4))
	fmt.Println(odt.powSum(2, 3, 2, 1e9+7))
	odt.add(0, 0, 1)
	fmt.Println(odt)
	odt.Set(1, 5, 2)
	fmt.Println(odt)
	odt.split(2)
	odt.Set(1, 1, 1)
	fmt.Println(odt)
	odt.EnumerateRange(0, n-1, func(block *ODTBlock) {
		fmt.Println(block)
	})
}

func main() {
	demo()
}

type Value = int

// 链表/数组能更简洁地维护分裂与合并操作
type ODTBlock struct {
	left, right int
	value       Value
}

// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/odt.go
type ODT []ODTBlock

// 当区间随机时,可以用珂朵莉树来维护区间信息.
func NewODT(nums []Value) ODT {
	n := len(nums)
	res := make(ODT, n)
	for i := range res {
		res[i] = ODTBlock{i, i, nums[i]}
	}
	return res
}

// [left, right] 区间赋值.
func (t *ODT) Set(left, right int, value Value) {
	begin, end := t.prepare(left, right)
	t.assign(begin, end, right, value)
}

// 获取区间 [left,right] 中的块的信息.
func (t ODT) EnumerateRange(left, right int, f func(block *ODTBlock)) {
	begin, end := t.prepare(left, right)
	for i := begin; i < end; i++ {
		f(&t[i])
	}
}

// !所有操作前需要预分裂 保证后续操作在 [left, right] 内部
func (t *ODT) prepare(left, right int) (begin, end int) {
	begin = t.split(left)
	end = t.split(right + 1)
	return
}

// 区间加
func (t ODT) add(left, right int, val Value) {
	t.EnumerateRange(left, right, func(block *ODTBlock) {
		block.value += val
	})
}

// 区间第k小(1-based)
func (t ODT) kth(left, right, k int) Value {
	blocks := []ODTBlock{}
	t.EnumerateRange(left, right, func(block *ODTBlock) {
		blocks = append(blocks, *block)
	})
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
func (t ODT) powSum(left, right int, exp int, mod int) (res int) {
	t.EnumerateRange(left, right, func(block *ODTBlock) {
		res += (block.right - block.left + 1) * pow(block.value, exp, mod)
		res %= mod
	})
	return
}

// [l, r] => [l, mid-1] [mid, r], 返回指向后者的迭代器.
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

// !区间[left,right]赋值操作
//  这里传入right是为了将begin块的右端点更新为right
func (t *ODT) assign(begin, end, right int, value Value) {
	ot := *t
	ot[begin].right = right
	ot[begin].value = value
	if begin+1 < end {
		*t = append(ot[:begin+1], ot[end:]...)
	}
}

func pow(base int, exp int, mod int) int {
	base %= mod
	res := 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}
