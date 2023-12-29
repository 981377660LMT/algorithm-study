// 空间复杂度O(n*O(merge))
// 时间复杂度O((n+q)*O(merge))

package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ivan and Burgers
// https://codeforces.com/contest/1100/problem/F
// 给定一个数组和q个查询,每个查询给定一个区间[lefti, righti],
// 求在原数组区间[lefti, righti]中选取任意个数,能凑出的最大异或和
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int32, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	D := NewDivideAndConquerOffline(nums)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var left, right int
		fmt.Fscan(in, &left, &right)
		left--
		D.AddQuery(left, right)
	}

	res := D.QueryAll()
	for _, r := range res {
		fmt.Fprintln(out, r)
	}
}

type ArrayItem = int32
type QueryRes = int32
type DataStrcuture = *VS

func (dc *DivideAndConquerOffline) Init() DataStrcuture {
	return &VS{}
}
func (dc *DivideAndConquerOffline) Clear(e DataStrcuture) {
	e.Clear()
}
func (dc *DivideAndConquerOffline) Add(e DataStrcuture, value ArrayItem) {
	e.Add(value)
}
func (dc *DivideAndConquerOffline) Copy(e DataStrcuture) DataStrcuture {
	return e.Copy()
}
func (dc *DivideAndConquerOffline) Merge(e1, e2 DataStrcuture) DataStrcuture {
	return e1.Or(e2)
}
func (dc *DivideAndConquerOffline) Query(e DataStrcuture) QueryRes {
	return e.Max(0)
}
func (dc *DivideAndConquerOffline) QueryLeaf(value ArrayItem) QueryRes {
	return value
}

// 猫树分治.
// !一共调用了n次Init、Clear, nlogn次Add、Copy, q次Merge、Query、QueryLeaf.
type DivideAndConquerOffline struct {
	arr               []ArrayItem
	data              []DataStrcuture // 用于维护arr区间查询的结果
	qid, qStart, qEnd []int32
	queries           [][3]int32
	tmpStart, tmpEnd  []int32
	res               []QueryRes
}

func NewDivideAndConquerOffline(arr []ArrayItem) *DivideAndConquerOffline {
	res := &DivideAndConquerOffline{arr: arr}
	data := make([]DataStrcuture, len(arr))
	for i := range arr {
		data[i] = res.Init()
	}
	res.data = data
	return res
}

func (dc *DivideAndConquerOffline) AddQuery(start, end int) {
	dc.qid = append(dc.qid, int32(len(dc.qid)))
	dc.qStart = append(dc.qStart, int32(start))
	dc.qEnd = append(dc.qEnd, int32(end))
}

func (dc *DivideAndConquerOffline) QueryAll() []QueryRes {
	n := int32(len(dc.data))
	q := int32(len(dc.qid))
	dc.tmpStart = make([]int32, q)
	dc.tmpEnd = make([]int32, q)
	dc.res = make([]QueryRes, q)
	dc.solve(0, n, 0, q)
	return dc.res
}

func (dc *DivideAndConquerOffline) solve(nStart, nEnd, qStart, qEnd int32) {
	if qStart >= qEnd {
		return
	}
	if nStart+1 == nEnd {
		for i := qStart; i < qEnd; i++ {
			dc.res[dc.qid[i]] = dc.QueryLeaf(dc.arr[nStart])
		}
		return
	}
	leftCount, rightCount := int32(0), int32(0)
	mid := (nStart + nEnd) >> 1
	dc.Clear(dc.data[mid])
	dc.Add(dc.data[mid], dc.arr[mid])
	for i := mid - 1; i >= nStart; i-- {
		dc.data[i] = dc.Copy(dc.data[i+1])
		dc.Add(dc.data[i], dc.arr[i])
	}
	for i := mid + 1; i < nEnd; i++ {
		dc.data[i] = dc.Copy(dc.data[i-1])
		dc.Add(dc.data[i], dc.arr[i])
	}
	for i := qStart; i < qEnd; i++ {
		id := dc.qid[i]
		if start := dc.qStart[id]; start < mid {
			if end := dc.qEnd[id]; end > mid {
				dc.res[id] = dc.Query(dc.Merge(dc.data[start], dc.data[end-1]))
			} else {
				dc.tmpStart[leftCount] = id
				leftCount++
			}
		} else {
			dc.tmpEnd[rightCount] = id
			rightCount++
		}
	}

	for i := int32(0); i < leftCount; i++ {
		dc.qid[qStart+i] = dc.tmpStart[i]
	}
	for i := int32(0); i < rightCount; i++ {
		dc.qid[qStart+leftCount+i] = dc.tmpEnd[i]
	}
	dc.solve(nStart, mid, qStart, qStart+leftCount)
	dc.solve(mid, nEnd, qStart+leftCount, qStart+leftCount+rightCount)
}

type VS [20]int32

func NewVectorSpaceArray(nums []int32) *VS {
	res := VS{}
	for _, num := range nums {
		res.Add(num)
	}
	return &res
}

func (lb *VS) Add(num int32) {
	if num != 0 {
		for i := len(lb) - 1; i >= 0; i-- {
			if num>>i&1 == 1 {
				if lb[i] == 0 {
					lb[i] = num
					break
				}
				num ^= lb[i]
			}
		}
	}
}

// 求xor与所有向量异或的最大值.
func (lb *VS) Max(xor int32) int32 {
	res := xor
	for i := len(lb) - 1; i >= 0; i-- {
		if tmp := res ^ lb[i]; tmp > res {
			res = tmp
		}
	}
	return res
}

// 求xor与所有向量异或的最小值.
func (lb *VS) Min(xorVal int32) int32 {
	res := xorVal
	for i := len(lb) - 1; i >= 0; i-- {
		if tmp := res ^ lb[i]; tmp < res {
			res = tmp
		}
	}
	return res
}

func (lb *VS) Copy() *VS {
	res := &VS{}
	copy(res[:], lb[:])
	return res
}

func (lb *VS) Clear() {
	for i := range lb {
		lb[i] = 0
	}
}

func (lb *VS) Len() int {
	return len(lb)
}

func (lb *VS) ForEach(f func(base int32)) {
	for _, base := range lb {
		f(base)
	}
}

func (lb *VS) Has(v int32) bool {
	for i := len(lb) - 1; i >= 0; i-- {
		if v == 0 {
			break
		}
		v = min32(v, v^lb[i])
	}
	return v == 0
}

// Merge.
func (lb *VS) Or(other *VS) *VS {
	v1, v2 := lb, other
	if v1.Len() < v2.Len() {
		v1, v2 = v2, v1
	}
	res := v1.Copy()
	for _, base := range v2 {
		res.Add(base)
	}
	return res
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}
