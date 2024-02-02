// https://www.luogu.com.cn/problem/CF786C
// CF786C-Till I Collapse-每段中不同数字的个数不超过k个
// !对于k=1~n, 分别求出最小的m, 使得存在一种将n个数划分成m段的方案，每段中不同数字的个数不超过k个。
// n<=1e5
//
// !1. 答案的范围不超过n//k(最坏情况是都不一样)
// !2. 当k单调递增时，答案单调不减，即答案具有二分性质

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	D := NewDictionary()
	for i := 0; i < n; i++ {
		nums[i] = D.Id(nums[i])
	}

	counter := make([]int, D.Size())
	count := 0

	// 每段中不同数字的个数不超过k个，求最小的m
	// 分组循环求解.
	f := func(k int) int {
		groupCount := 0
		ptr := 0
		for ptr < n {
			counter[nums[ptr]]++
			count = 1
			start := ptr
			ptr++
			for ptr < n {
				v := nums[ptr]
				if counter[v] == 0 {
					count++
				}
				if count > k {
					break
				}
				counter[v]++
				ptr++
			}
			groupCount++
			for i := start; i < ptr; i++ {
				counter[nums[i]]--
			}
		}
		return groupCount
	}

	res := ContinuousResultFunctionTrick(1, n, f)
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// 答案连续的函数分治求值.
// !给定一个函数f, 当i=left,left+1,...,right时，分别求出f(i)的值.
// !要求:f不同答案的个数不超过某个数量级(例如sqrt(n))，且答案相同的i一定连续.
// 可以对 [left,right] 值域分治，如果区间内的结果都相同，则不再分治
// 时间复杂度类似整除分块 O(f(n)√n)
// https://codeforces.com/problemset/problem/786/C
func ContinuousResultFunctionTrick(left, right int, f func(int) int) []int {
	res := make([]int, right-left+1)
	var solve func(int, int)
	solve = func(l, r int) {
		if l > r {
			return
		}
		resL, resR := f(l), f(r)
		if resL == resR {
			for i := l; i <= r; i++ {
				res[i-left] = resL
			}
			return
		}
		res[l-left] = resL
		res[r-left] = resR
		mid := (l + r) / 2
		solve(l+1, mid)
		solve(mid+1, r-1)
	}
	solve(left, right)
	return res
}

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
}
