package main

import (
	"bufio"
	"fmt"
	"os"
)

// P6240 好吃的题目(分治+背包,区间01背包)
// https://www.luogu.com.cn/problem/P6240
// 有n个物品，每个物品有一个重量和一个分数，
// 现在有q个询问，每个询问给出一个区间[l,r]和一个容量c，
// 要求在这个区间内选出若干个物品，使得选出的物品的重量和不超过c，且选出的物品的分数之和最大。
// n<=4e4,q<=2e5,c<=200
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &weights[i])
	}
	scores := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &scores[i])
	}

	queries := make([][3]int, q) // [start, end, capacity]

	for i := 0; i < q; i++ {
		var left, right, capacity int
		fmt.Fscan(in, &left, &right, &capacity)
		left--
		queries[i] = [3]int{left, right, capacity}
	}
}

type ArrayItem = int32
type QueryRes = int32
type DataStrcuture = []int32

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
