// Update
// QueryRange
// QueryAll
// QueryPrefix
// MaxRight

package main

import "fmt"

func main() {
	bit2 := NewBITGroupRangeAddFrom(10, func(index int) E { return index })
	fmt.Println(bit2)
	bit2.UpdateRange(0, 10, 1)
	fmt.Println(bit2)
	fmt.Println(bit2.QueryRange(2, 5))
}

type E = int

func e() E           { return 0 }
func op(e1, e2 E) E  { return e1 + e2 }
func inv(e E) E      { return -e }    // 如果只查询前缀, 可以不需要是群
func mul(e E, k E) E { return e * k } // 如果不需要区间修改, 可以不需要是环

type BITGroup struct {
	n     int
	data  []E
	total E
}

func NewBITGroup(n int) *BITGroup {
	data := make([]E, n)
	for i := range data {
		data[i] = e()
	}
	return &BITGroup{n: n, data: data, total: e()}
}

func NewBITGroupFrom(n int, f func(index int) E) *BITGroup {
	total := e()
	data := make([]E, n)
	for i := range data {
		data[i] = f(i)
		total = op(total, data[i])
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = op(data[j-1], data[i-1])
		}
	}
	return &BITGroup{n: n, data: data, total: total}
}

func (fw *BITGroup) Update(i int, x E) {
	fw.total = op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = op(fw.data[i-1], x)
	}
}

func (fw *BITGroup) QueryAll() E { return fw.total }

// [0, end)
func (fw *BITGroup) QueryPrefix(end int) E {
	if end > fw.n {
		end = fw.n
	}
	res := e()
	for end > 0 {
		res = op(res, fw.data[end-1])
		end &= end - 1
	}
	return res
}

// [start, end)
func (fw *BITGroup) QueryRange(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > fw.n {
		end = fw.n
	}
	if start == 0 {
		return fw.QueryPrefix(end)
	}
	if start > end {
		return e()
	}
	pos, neg := e(), e()
	for end > start {
		pos = op(pos, fw.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = op(neg, fw.data[start-1])
		start &= start - 1
	}
	return op(pos, inv(neg))
}

// 最大的 right 使得 check(QueryPrefix(right)) == true.
//  check(value, right): value 对应的是 [0, right) 的和.
//
//  e.g.:
//  0/1 树状数组找到第 k(0-indexed) 个 1:
//  func (fw *BITGroup) Kth(k E) int {
//  	return fw.MaxRight(func(preSum E, _ int) bool {
//  		return preSum <= k
//  	})
//  }
func (fw BITGroup) MaxRight(check func(value E, right int) bool) int {
	i := 0
	cur := e()
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

func (fw *BITGroup) String() string {
	res := []string{}
	for i := 0; i < fw.n; i++ {
		res = append(res, fmt.Sprintf("%d", fw.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BITGroup%v", res)
}

type BITGroupRangeAdd struct {
	n    int
	bit0 *BITGroup
	bit1 *BITGroup
}

func NewBITGroupRangeAdd(n int) *BITGroupRangeAdd {
	return &BITGroupRangeAdd{n: n, bit0: NewBITGroup(n), bit1: NewBITGroup(n)}
}

func NewBITGroupRangeAddFrom(n int, f func(index int) E) *BITGroupRangeAdd {
	return &BITGroupRangeAdd{n: n, bit0: NewBITGroupFrom(n, f), bit1: NewBITGroup(n)}
}

func (fw *BITGroupRangeAdd) Update(index int, x E) {
	fw.bit0.Update(index, x)
}

func (fw *BITGroupRangeAdd) UpdateRange(start, end int, x E) {
	fw.bit0.Update(start, mul(x, -start))
	fw.bit0.Update(end, mul(x, end))
	fw.bit1.Update(start, x)
	fw.bit1.Update(end, inv(x))
}

func (fw *BITGroupRangeAdd) QueryRange(start, end int) E {
	rightRes := op(mul(fw.bit1.QueryPrefix(end), end), fw.bit0.QueryPrefix(end))
	leftRes := op(mul(fw.bit1.QueryPrefix(start), start), fw.bit0.QueryPrefix(start))
	return op(inv(leftRes), rightRes)
}

func (fw *BITGroupRangeAdd) String() string {
	res := []string{}
	for i := 0; i < fw.n; i++ {
		res = append(res, fmt.Sprintf("%d", fw.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BITGroupRangeAdd%v", res)
}
