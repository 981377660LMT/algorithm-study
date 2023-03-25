// Update
// QueryRange
// QueryAll
// QueryPrefix
// MaxRight

package main

import "fmt"

// https://leetcode.cn/problems/longest-uploaded-prefix/
// 最长上传前缀
type LUPrefix struct {
	tree *FenwickTree
}

func Constructor(n int) LUPrefix {
	return LUPrefix{tree: NewFenwickTree(n + 10)}
}

func (this *LUPrefix) Upload(video int) {
	this.tree.Update(video-1, 1)
}

func (this *LUPrefix) Longest() int {
	return this.tree.MaxRight(func(preSum E, right int) bool { return preSum >= right })
}

func main() {
	lru := Constructor(4)
	lru.Upload(3)
	fmt.Println(lru.Longest())
	lru.Upload(1)
	fmt.Println(lru.Longest())
	lru.Upload(2)
	fmt.Println(lru.Longest())
}

// ["LUPrefix", "upload", "longest", "upload", "longest", "upload", "longest"]
// [[4], [3], [], [1], [], [2], []]
// 输出：
// [null, null, 0, null, 1, null, 3]

//
//
//
type E = int

func e() E          { return 0 }
func op(e1, e2 E) E { return e1 + e2 }
func inv(e E) E     { return -e } // 如果只查询前缀, 可以不需要是群

type FenwickTree struct {
	n     int
	data  []E
	total E
	unit  E
}

func NewFenwickTree(n int, nums ...E) *FenwickTree {
	fw := &FenwickTree{n: n, unit: e()}
	if len(nums) == 0 {
		nums = make([]E, n)
		for i := range nums {
			nums[i] = fw.unit
		}
	}
	fw.build(nums)
	return fw
}

func (fw *FenwickTree) QueryAll() E { return fw.total }

// [0, right)
func (fw *FenwickTree) QueryPrefix(right int) E {
	if right > fw.n {
		right = fw.n
	}
	res := fw.unit
	for right > 0 {
		res = op(res, fw.data[right-1])
		right &= right - 1
	}
	return res
}

// [left, right)
func (fw *FenwickTree) QueryRange(left, right int) E {
	if left < 0 {
		left = 0
	}
	if right > fw.n {
		right = fw.n
	}
	if left == 0 {
		return fw.QueryPrefix(right)
	}
	if left > right {
		return fw.unit
	}
	pos, neg := fw.unit, fw.unit
	for right > left {
		pos = op(pos, fw.data[right-1])
		right &= right - 1
	}
	for left > right {
		neg = op(neg, fw.data[left-1])
		left &= left - 1
	}
	return op(pos, inv(neg))
}

func (fw *FenwickTree) Update(i int, x E) {
	fw.total = op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = op(fw.data[i-1], x)
	}
}

// 最大的 right 使得 check(QueryPrefix(right)) == true.
//  check(value, right): value 对应的是 [0, right) 的和.
//
//  e.g.:
//  0/1 树状数组找到第 k(0-indexed) 个 1:
//  func (fw *FenwickTree) Kth(k E) int {
//  	return fw.MaxRight(func(preSum E, _ int) bool {
//  		return preSum <= k
//  	})
//  }
func (fw FenwickTree) MaxRight(check func(value E, right int) bool) int {
	i := 0
	cur := fw.unit
	k := 1
	for 2*k <= fw.n {
		k *= 2
	}
	for k > 0 {
		if i+k-1 < len(fw.data) {
			t := op(cur, fw.data[i+k-1])
			if check(t, i+k) {
				i += k
				cur = t
			}
		}
		k >>= 1
	}
	return i
}

func (fw *FenwickTree) String() string {
	res := []string{}
	for i := 0; i < fw.n; i++ {
		res = append(res, fmt.Sprintf("%d", fw.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("FenwickTree%v", res)
}

func (fw *FenwickTree) build(nums []E) {
	n := fw.n
	fw.data = append(fw.data, nums...)
	fw.total = fw.unit
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			fw.data[j-1] = op(fw.data[i-1], fw.data[j-1])
		}
	}
	fw.total = fw.QueryPrefix(n)
}
