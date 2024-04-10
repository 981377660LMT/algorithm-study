// DivideIntervalBinaryLift/反向st表
// 倍增拆分序列上的区间 `[start,end)`
// 一共拆分成[0,log]层，每层有n个元素.
// !jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n).
//
// api:
//
//	NewDivideIntervalByBinaryLift(n int32) *DivideIntervalByBinaryLift
//	EnumerateRange(start, end int32, f func(level, index int32))
//	EnumerateRangeDangerously(start, end int32, f func(level, index int32))
//	EnumerateRange2(start1, end1 int32, start2, end2 int32, f func(level, index1, index2 int32))
//	EnumerateRange2Dangerously(start1, end1 int32, start2, end2 int32, f func(level, index1, index2 int32))
//	PushDown(f func(parentLevel, parentIndex int32, childLevel, childIndex1, childIndex2 int32))
//	Size() int32
//	Log() int32

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
)

type DivideIntervalByBinaryLift struct {
	n, log int32
	size   int32
}

func NewDivideIntervalByBinaryLift(n int32) *DivideIntervalByBinaryLift {
	log := int32(bits.Len32(uint32(n))) - 1
	size := n * (log + 1)
	return &DivideIntervalByBinaryLift{n: n, log: log, size: size}
}

// O(logn)遍历[start,end)区间内的所有jump.
func (d *DivideIntervalByBinaryLift) EnumerateRange(start, end int32, f func(level, index int32)) {
	if start >= end {
		return
	}
	cur := start
	log := d.log
	len := end - start
	for k := log; k >= 0; k-- {
		if len&(1<<k) != 0 {
			f(k, cur)
			cur += 1 << k
			if cur >= end {
				return
			}
		}
	}
	f(0, cur)
}

// O(1)遍历[start,end)区间内的所有jump.
// !要求运算幂等(idempotent).
func (d *DivideIntervalByBinaryLift) EnumerateRangeDangerously(start, end int32, f func(level, index int32)) {
	if start >= end {
		return
	}
	k := int32(bits.Len32(uint32(end-start))) - 1
	f(k, start)
	f(k, end-(1<<k))
}

// O(logn)遍历[start1,end1)区间和[start2,end2)区间内的所有jump.要求区间长度相等.
func (d *DivideIntervalByBinaryLift) EnumerateRange2(start1, end1 int32, start2, end2 int32, f func(level, index1, index2 int32)) {
	if end1-start1 != end2-start2 {
		panic("not same length")
	}
	if start1 >= end1 {
		return
	}
	cur1, cur2 := start1, start2
	log := d.log
	len := end1 - start1
	for k := log; k >= 0; k-- {
		if len&(1<<k) != 0 {
			f(k, cur1, cur2)
			cur1 += 1 << k
			cur2 += 1 << k
			if cur1 >= end1 {
				return
			}
		}
	}
	f(0, cur1, cur2)
}

// O(1)遍历[start1,end1)区间和[start2,end2)区间内的所有jump.要求区间长度相等.
// !要求运算幂等(idempotent).
func (d *DivideIntervalByBinaryLift) EnumerateRange2Dangerously(start1, end1 int32, start2, end2 int32, f func(level, index1, index2 int32)) {
	if end1-start1 != end2-start2 {
		panic("not same length")
	}
	if start1 >= end1 {
		return
	}
	k := int32(bits.Len32(uint32(end1-start1))) - 1
	f(k, start1, start2)
	f(k, end1-(1<<k), end2-(1<<k))
}

// 从高的jump开始下推信息，更新底部jump的答案.
// O(n*log(n)).
func (d *DivideIntervalByBinaryLift) PushDown(
	f func(parentLevel, parentIndex int32, childLevel, childIndex1, childIndex2 int32),
) {
	n, log := d.n, d.log
	for k := log - 1; k >= 0; k-- {
		for i := int32(0); i < n-(1<<k); i++ {
			f(k+1, i, k, i, i+(1<<k))
		}
	}
}

func (d *DivideIntervalByBinaryLift) Size() int32 { return d.size }
func (d *DivideIntervalByBinaryLift) Log() int32  { return d.log }

func main() {
	// P3295()
	// demo()

	solve := func(n int32, rangeToRangeInfo [][5]int32, start int32) []int {
		D := NewDivideIntervalByBinaryLift(n)
		size := D.Size()
		newGraph := make([][]Neighbor, size*2+n) // 入点：[0,size)，出点：[size,2*size), 实际的底层结点：[2*size,2*size+n)
		addEdge := func(from, to, w int32) {
			newGraph[from] = append(newGraph[from], Neighbor{to, w})
		}
		D.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
			p, c1, c2 := pLevel*n+pIndex, cLevel*n+cIndex1, cLevel*n+cIndex2
			addEdge(c1, p, 0)
			addEdge(c2, p, 0)
			addEdge(p+size, c1+size, 0)
			addEdge(p+size, c2+size, 0)
		})

		for i := int32(0); i < n; i++ {
			// to[u][0]=fa;in[u][0]=out[u][0]=u;
			offset := 2 * size
			p, c := i, i+offset
			addEdge(c, p, 0)
			addEdge(p+size, c, 0)
		}

		// !2.区间入点和区间出点之间相互连边.
		addRangeToRange := func(u1, v1, u2, v2, w int32) {
			from, to := make([]int32, 0, 2), make([]int32, 0, 2)
			D.EnumerateRangeDangerously(u1, v1, func(level, index int32) {
				id := level*n + index
				from = append(from, id)
			})
			D.EnumerateRangeDangerously(u2, v2, func(level, index int32) {
				id := (level*n + index) + size
				to = append(to, id)
			})
			for _, a := range from {
				for _, b := range to {
					addEdge(a, b, w)
				}
			}
		}

		for _, info := range rangeToRangeInfo {
			start1, end1, start2, end2, w := info[0], info[1], info[2], info[3], info[4]
			addRangeToRange(start1, end1, start2, end2, w)
		}

		dist := Dijkstra(int32(len(newGraph)), newGraph, start+2*size)
		res := make([]int, n)
		for i := int32(0); i < n; i++ {
			res[i] = dist[i+2*size]
		}
		return res
	}

	bruteForce := func(n int32, rangeToRangeInfo [][5]int32, start int32) []int {
		adjList := make([][]Neighbor, n)
		for _, info := range rangeToRangeInfo {
			start1, end1, start2, end2, w := info[0], info[1], info[2], info[3], info[4]
			for i := start1; i < end1; i++ {
				for j := start2; j < end2; j++ {
					adjList[i] = append(adjList[i], Neighbor{j, w})
				}
			}
		}

		dist := Dijkstra(n, adjList, start)
		return dist
	}

	n := int32(rand.Intn(10000)) + 1
	m := int32(rand.Intn(5)) + 1
	rangeToRangeInfo := make([][5]int32, 0, m)
	for i := int32(0); i < m; i++ {
		start1, end1 := int32(rand.Intn(int(n))), int32(rand.Intn(int(n)))+1
		start2, end2 := int32(rand.Intn(int(n))), int32(rand.Intn(int(n)))+1
		if start1 > end1 {
			start1, end1 = end1, start1
		}
		if start2 > end2 {
			start2, end2 = end2, start2
		}
		w := int32(rand.Intn(100) + 1)
		rangeToRangeInfo = append(rangeToRangeInfo, [5]int32{start1, end1, start2, end2, w})
	}

	start := int32(rand.Intn(int(n)))
	res1 := solve(n, rangeToRangeInfo, start)
	res2 := bruteForce(n, rangeToRangeInfo, start)
	for i := range res1 {
		if res1[i] != res2[i] {
			fmt.Println("not equal")
			return
		}
	}

	fmt.Println("pass", n, len(rangeToRangeInfo), start)
}

func demo() {
	n := int32(10)
	D := NewDivideIntervalByBinaryLift(n)
	values := make([]int32, D.Size())
	D.EnumerateRange(1, 9, func(level, index int32) { values[level*n+index] += 5 })
	D.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
		p, c1, c2 := pLevel*n+pIndex, cLevel*n+cIndex1, cLevel*n+cIndex2
		values[c1] = max(values[c1], values[p])
		values[c2] = max(values[c2], values[p])
	})
	fmt.Println(values[:n])
}

// 萌萌哒
// https://www.luogu.com.cn/problem/P3295
// 给定一个长度为n的大数，每个大数元素为0到9之间的整数(注意不能有前导零)。
// 再给定一些约束条件，形如[start1,end1,start2,end2]，表示[start1,end1)区间内的数和[start2,end2)区间内的数相等。
// 问满足以上所有条件的数有多少个，对1e9+7取模。
func P3295() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7
	qpow := func(a, b int) int {
		res := 1
		for b > 0 {
			if b&1 == 1 {
				res = res * a % MOD
			}
			a = a * a % MOD
			b >>= 1
		}
		return res
	}

	var n, m int32
	fmt.Fscan(in, &n, &m)
	D := NewDivideIntervalByBinaryLift(n)
	ufs := make([]*UnionFindArraySimple32, D.Log()+1)
	for i := range ufs {
		ufs[i] = NewUnionFindArraySimple32(n)
	}

	for i := int32(0); i < m; i++ {
		var start1, end1, start2, end2 int32
		fmt.Fscan(in, &start1, &end1, &start2, &end2)
		start1, start2 = start1-1, start2-1
		D.EnumerateRange2Dangerously(start1, end1, start2, end2, func(level, i1, i2 int32) {
			ufs[level].Union(i1, i2)
		})
	}

	D.PushDown(func(pLevel, pIndex, cLevel, cIndex1, _ int32) {
		root := ufs[pLevel].Find(pIndex)
		ufs[cLevel].Union(cIndex1, root)
		ufs[cLevel].Union(cIndex1+1<<cLevel, root+1<<cLevel)
	})

	uf := ufs[0]
	part := int(uf.Part)
	fmt.Fprintln(out, 9*qpow(10, part-1)%MOD)
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = u.Find(u.data[key])
	return u.data[key]
}

func (u *UnionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}

// 370. 区间加法
// https://leetcode.cn/problems/range-addition/description/
// 用于验证两种拆分方式对幂等性的要求.
func getModifiedArray(length int, updates [][]int) []int {
	n := int32(length)
	D := NewDivideIntervalByBinaryLift(n)
	values := make([]int, D.Size())
	for _, u := range updates {
		start, end, inc := u[0], u[1]+1, u[2]
		// 这里不能用EnumerateRangeDangerously，因为加法不是幂等操作.
		D.EnumerateRange(int32(start), int32(end), func(level, index int32) {
			values[level*n+index] += inc
		})
	}
	D.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
		p, c1, c2 := pLevel*n+pIndex, cLevel*n+cIndex1, cLevel*n+cIndex2
		values[c1] += values[p]
		values[c2] += values[p]
	})
	return values[:n]
}

const INF int = 1e18

type Neighbor struct {
	to, weight int32
}

func Dijkstra(n int32, adjList [][]Neighbor, start int32) (dist []int) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0

	pq := nhp(func(a, b H) int {
		return a.dist - b.dist
	}, []H{{start, 0}})

	for pq.Len() > 0 {
		curNode := pq.Pop()
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge.to, edge.weight
			if cand := curDist + int(weight); cand < dist[next] {
				dist[next] = cand
				pq.Push(H{next, cand})
			}
		}
	}

	return
}

type H = struct {
	node int32
	dist int
}

// Should return a number:
//
//	negative , if a < b
//	zero     , if a == b
//	positive , if a > b
type Comparator func(a, b H) int

func nhp(comparator Comparator, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{comparator: comparator, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data       []H
	comparator Comparator
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		return
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Peek() (value H) {
	if h.Len() == 0 {
		return
	}
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.comparator(h.data[root], h.data[parent]) < 0; parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.comparator(h.data[left], h.data[minIndex]) < 0 {
			minIndex = left
		}

		if right < n && h.comparator(h.data[right], h.data[minIndex]) < 0 {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}
