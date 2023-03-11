// 值域为1e9时需要先离散化
package main

type S = int

const INF int = 1e18

func (*FenwickTreePrefix) e() S        { return -INF }
func (*FenwickTreePrefix) op(a, b S) S { return max(a, b) }
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type FenwickTreePrefix struct {
	n    int
	data []S
}

func NewFenwickTreePrefix(n int) *FenwickTreePrefix {
	res := &FenwickTreePrefix{n, make([]S, n+1)}
	for i := 0; i < n+1; i++ {
		res.data[i] = res.e()
	}
	return res
}

func NewFenwickTreePrefixWithSlice(nums []S) *FenwickTreePrefix {
	n := len(nums)
	res := &FenwickTreePrefix{n, make([]S, n+1)}
	for i := 1; i < n+1; i++ {
		res.data[i] = nums[i-1]
	}
	for i := 1; i < n+1; i++ {
		if j := i + (i & -i); j <= n {
			res.data[j] = res.op(res.data[j], res.data[i])
		}
	}
	return res
}

// 单点更新index处的元素.
// 0 <= index < n
func (f *FenwickTreePrefix) Update(index int, value S) {
	for index++; index <= f.n; index += index & -index {
		f.data[index] = f.op(f.data[index], value)
	}
}

// 查询前缀区间 [0,right) 的值.
// 0 <= right <= n
func (f *FenwickTreePrefix) Query(right int) S {
	res := f.e()
	if right > f.n {
		right = f.n
	}
	for ; right > 0; right &= right - 1 {
		res = f.op(res, f.data[right])
	}
	return res
}
