// 空间复杂度O(n*O(merge))
// 时间复杂度O((n+q)*O(merge))

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
// n<=4e4,q<=2e5,c<=200,scores[i]<=1e7
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	goods := make([]*[2]int32, n) // [重量,分数]
	for i := range goods {
		goods[i] = &[2]int32{}
	}
	for i := 0; i < n; i++ {
		var weight int32
		fmt.Fscan(in, &weight)
		goods[i][0] = weight
	}
	for i := 0; i < n; i++ {
		var score int32
		fmt.Fscan(in, &score)
		goods[i][1] = score
	}

	D := NewDivideAndConquerOffline(goods)
	queries := make([][3]int32, q)
	for i := 0; i < q; i++ {
		var left, right, capacity int32
		fmt.Fscan(in, &left, &right, &capacity)
		left--
		queries[i] = [3]int32{left, right, capacity}
		D.AddQuery(int(left), int(right))
	}

	res := D.Run(
		// init
		func() Merger {
			return &[C]int32{}
		},
		// clear
		func(e Merger) {
			for i := range e {
				e[i] = 0
			}
		},
		// add
		func(e Merger, value ArrayItem) {
			weight, score := value[0], value[1]
			for i := int32(len(e) - 1); i >= weight; i-- {
				e[i] = max32(e[i], e[i-weight]+score)
			}
		},
		// copy
		func(e Merger) Merger {
			copy_ := [C]int32{}
			copy(copy_[:], e[:])
			return &copy_
		},
		// queryMerge
		func(qid int, e1, e2 Merger) QueryRes {
			res := int32(0)
			capacity := queries[qid][2]
			for i := int32(0); i <= capacity; i++ {
				res = max32(res, int32(e1[i]+e2[capacity-i]))
			}
			return res
		},
		// queryLeaf
		func(qid int, value ArrayItem) QueryRes {
			capacity := queries[qid][2]
			if value[0] <= capacity {
				return value[1]
			} else {
				return 0
			}
		},
	)

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

const C int32 = 201

type ArrayItem = *[2]int32
type QueryRes = int32
type Merger = *[C]int32

// 猫树分治.
// !一共调用了n次Init、Clear, nlogn次Add、Copy, q次Merge、Query、QueryLeaf.
type DivideAndConquerOffline struct {
	arr               []ArrayItem
	merger            []Merger // 用于维护arr区间查询的结果
	qid, qStart, qEnd []int32
	tmpStart, tmpEnd  []int32
	res               []QueryRes

	init       func() Merger
	clear      func(e Merger)
	add        func(e Merger, value ArrayItem)
	copy       func(e Merger) Merger
	queryMerge func(qid int, e1, e2 Merger) QueryRes
	queryLeaf  func(qid int, value ArrayItem) QueryRes
}

func NewDivideAndConquerOffline(arr []ArrayItem) *DivideAndConquerOffline {
	return &DivideAndConquerOffline{arr: arr}
}

func (dc *DivideAndConquerOffline) AddQuery(start, end int) {
	dc.qid = append(dc.qid, int32(len(dc.qid)))
	dc.qStart = append(dc.qStart, int32(start))
	dc.qEnd = append(dc.qEnd, int32(end))
}

func (dc *DivideAndConquerOffline) Run(
	init func() Merger,
	clear func(e Merger),
	add func(e Merger, value ArrayItem),
	copy func(e Merger) Merger,
	queryMerge func(qid int, e1, e2 Merger) QueryRes,
	queryLeaf func(qid int, value ArrayItem) QueryRes,
) []QueryRes {
	dc.merger = make([]Merger, len(dc.arr))
	for i := range dc.arr {
		dc.merger[i] = init()
	}
	n := int32(len(dc.merger))
	q := int32(len(dc.qid))
	dc.tmpStart = make([]int32, q)
	dc.tmpEnd = make([]int32, q)
	dc.res = make([]QueryRes, q)
	dc.init = init
	dc.clear = clear
	dc.add = add
	dc.copy = copy
	dc.queryMerge = queryMerge
	dc.queryLeaf = queryLeaf
	dc.solve(0, n, 0, q)
	return dc.res
}

func (dc *DivideAndConquerOffline) solve(nStart, nEnd, qStart, qEnd int32) {
	if qStart >= qEnd {
		return
	}
	if nStart+1 == nEnd {
		for i := qStart; i < qEnd; i++ {
			id := int(dc.qid[i])
			dc.res[id] = dc.queryLeaf(id, dc.arr[nStart])
		}
		return
	}
	leftCount, rightCount := int32(0), int32(0)
	mid := (nStart + nEnd) >> 1
	dc.clear(dc.merger[mid])
	for i := mid - 1; i >= nStart; i-- {
		dc.merger[i] = dc.copy(dc.merger[i+1])
		dc.add(dc.merger[i], dc.arr[i])
	}
	dc.add(dc.merger[mid], dc.arr[mid])
	for i := mid + 1; i < nEnd; i++ {
		dc.merger[i] = dc.copy(dc.merger[i-1])
		dc.add(dc.merger[i], dc.arr[i])
	}
	for i := qStart; i < qEnd; i++ {
		id := dc.qid[i]
		if start := dc.qStart[id]; start < mid {
			if end := dc.qEnd[id]; end > mid {
				dc.res[id] = dc.queryMerge(int(id), dc.merger[start], dc.merger[end-1])
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
