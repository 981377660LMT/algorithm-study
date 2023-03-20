// 值域为1e9时需要先离散化

package main

import "fmt"

func main() {
	bit := NewFenwickTreePrefix(10)
	bit.Update(0, 1)
	bit.Update(1, 2)
	bit.Update(2, 3)
	fmt.Println(bit.Query(3))
	fmt.Println(bit.MaxRight(func(s int) bool { return s <= 3 }))
}

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
	res := &FenwickTreePrefix{n, make([]S, n)}
	for i := 0; i < n; i++ {
		res.data[i] = res.e()
	}
	return res
}

func NewFenwickTreePrefixWithSlice(nums []S) *FenwickTreePrefix {
	n := len(nums)
	res := &FenwickTreePrefix{n, make([]S, n)}
	for i := 0; i < n; i++ {
		res.data[i] = nums[i]
	}
	for i := 1; i <= n; i++ {
		if j := i + (i & -i); j <= n {
			res.data[j-1] = res.op(res.data[i-1], res.data[j-1])
		}
	}
	return res
}

// 单点更新index处的元素.
// 0 <= index < n
func (f *FenwickTreePrefix) Update(index int, value S) {
	for index++; index <= f.n; index += index & -index {
		f.data[index-1] = f.op(f.data[index-1], value)
	}
}

// 查询前缀区间 [0,right) 的值.
// 0 <= right <= n
func (f *FenwickTreePrefix) Query(right int) S {
	if right > f.n {
		right = f.n
	}
	res := f.e()
	for ; right > 0; right &= right - 1 {
		res = f.op(res, f.data[right-1])
	}
	return res
}

// 返回最大的right使得[0,right)中的元素之和满足f(s)。
func (ftp *FenwickTreePrefix) MaxRight(check func(s S) bool) int {
	i := 0
	s := ftp.e()
	k := 1
	for k<<1 <= ftp.n {
		k <<= 1
	}
	for k > 0 {
		if i+k-1 < ftp.n {
			t := ftp.op(s, ftp.data[i+k-1])
			if check(t) {
				i += k
				s = t
			}
		}
		k >>= 1
	}
	return i
}

// 0/1 树状数组查找第 k(0-based) 个1的位置.
func (ftp *FenwickTreePrefix) Kth(k int) int {
	return ftp.MaxRight(func(s S) bool { return s <= k })
}

// func (*FenwickTreePrefix) inv(a S) S   { return -a }
// func (ftp *FenwickTreePrefix) QueryRange(start, end int) S {
// 	if start < 0 {
// 		start = 0
// 	}
// 	if end > ftp.n {
// 		end = ftp.n
// 	}
// 	if start == 0 {
// 		return ftp.Query(end)
// 	}
// 	pos, neg := ftp.e(), ftp.e()
// 	for start < end {
// 		pos = ftp.op(pos, ftp.data[end-1])
// 		end &= end - 1
// 	}
// 	for start > end {
// 		neg = ftp.op(neg, ftp.data[start-1])
// 		start &= start - 1
// 	}
// 	return ftp.op(pos, ftp.inv(neg))
// }
