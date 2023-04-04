// https://yukicoder.me/problems/no/114
// !一个带权的无向图上有k个关键点，求联通k个关键点最小的代价(边权之和)。
// n<=35,k<=n,m<=n*(n-1)/2

// k<=15时最小斯坦纳树
// k>15时枚举+并查集

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u, v = u-1, v-1
		edges = append(edges, []int{u, v, w})
	}
	criticals := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &criticals[i])
		criticals[i]--
	}

	fmt.Fprintln(out, solve(n, edges, criticals))
}

func solve(n int, edges [][]int, criticals []int) int {
	if len(criticals) <= 15 {
		res, _, _ := MinimumSteinerTree(n, edges, criticals, nil)
		return res
	}

	sort.Slice(edges, func(i, j int) bool {
		return edges[i][2] < edges[j][2]
	})
	notCriticals := []int{}
	set := map[int]struct{}{}
	for _, c := range criticals {
		set[c] = struct{}{}
	}
	for i := 0; i < n; i++ {
		if _, ok := set[i]; !ok {
			notCriticals = append(notCriticals, i)
		}
	}

	// 选择点集state时,连通criticals的最小边权之和
	cal := func(state int) int {
		uf := NewUnionFindArray(n)
		ok := 0
		cost := 0
		for _, e := range edges {
			u, v, w := e[0], e[1], e[2]
			if state&(1<<u) != 0 && state&(1<<v) != 0 && uf.Union(u, v) {
				cost += w
				ok++
			}
		}
		if ok == bits.OnesCount(uint(state))-1 {
			return cost
		}
		return INF
	}

	res := INF
	state := 0
	for _, c := range criticals {
		state |= 1 << c
	}
	for i := 0; i < (1 << len(notCriticals)); i++ {
		curState := state
		for j := 0; j < len(notCriticals); j++ {
			if i&(1<<j) != 0 {
				curState |= 1 << notCriticals[j]
			}
		}
		res = min(res, cal(curState))
	}

	return res
}

const INF int = 1e18

// 一个联通的无向带权图上有k个关键点 criticals，求联通所有点最小的代价(边权之和)。
//  vWeights: 每个顶点的附加权重(一般为make([]int, n))。
func MinimumSteinerTree(n int, edges [][]int, criticals, vWeights []int) (cost int, es, vs []int) {
	graph := make([][][3]int, n) // (to,w,ei)
	for i, e := range edges {
		u, v, w := e[0], e[1], e[2]
		graph[u] = append(graph[u], [3]int{v, w, i})
		graph[v] = append(graph[v], [3]int{u, w, i})
	}
	if vWeights == nil {
		vWeights = make([]int, n)
	}

	m := len(edges)
	k := len(criticals)

	dp := make([][]int, 1<<k)
	for i := range dp {
		dp[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	for i := 0; i < n; i++ {
		dp[0][i] = vWeights[i]
	}

	par := make([][]int, 1<<k)
	for i := range par {
		par[i] = make([]int, n)
		for j := range par[i] {
			par[i][j] = -1
		}
	}

	for s := 1; s < 1<<k; s++ {
		curDp := dp[s]
		curS := s
		for curS != 0 {
			i := bits.TrailingZeros(uint(curS))
			v := criticals[i]
			curDp[v] = min(curDp[v], dp[curS^(1<<i)][v])
			curS ^= 1 << i
		}
		for t := s; t >= 0; {
			if t == 0 || t == s {
				t--
				continue
			}
			for v := 0; v < n; v++ {
				cand := dp[t][v] + dp[s^t][v] - vWeights[v]
				if cand < curDp[v] {
					curDp[v] = cand
					par[s][v] = 2 * t
				}
			}
			t = (t - 1) & s
		}

		inits := make([]H, n)
		for i := 0; i < n; i++ {
			inits[i] = H{curDp[i], i}
		}
		pq := NewHeap(func(a, b H) bool {
			return a[0] < b[0]
		}, inits)
		for pq.Len() > 0 {
			item := pq.Pop()
			dv, v := item[0], item[1]
			if dv != curDp[v] {
				continue
			}
			for _, e := range graph[v] {
				to, cost, id := e[0], e[1], e[2]
				cand := dv + cost + vWeights[to]
				if cand < curDp[to] {
					curDp[to] = cand
					par[s][to] = 2*id + 1
					pq.Push(H{cand, to})
				}
			}
		}

	}

	// 复元
	usedV, usedE := make([]bool, n), make([]bool, m)
	vToK := make([]int, n)
	for i := range vToK {
		vToK[i] = -1
	}
	for i := 0; i < k; i++ {
		vToK[criticals[i]] = i
	}

	root, min_ := 0, INF
	for i, v := range dp[len(dp)-1] {
		if v < min_ {
			root, min_ = i, v
		}
	}
	queue := [][2]int{{(1 << k) - 1, root}}
	usedV[root] = true

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]
		s, v := item[0], item[1]
		if s == 0 {
			continue
		}
		if par[s][v] == -1 {
			k := vToK[v]
			queue = append(queue, [2]int{s ^ (1 << k), v})
			continue
		} else if par[s][v]&1 == 1 {
			eid := par[s][v] / 2
			e := edges[eid]
			w := v ^ e[0] ^ e[1]
			usedV[w] = true
			usedE[eid] = true
			queue = append(queue, [2]int{s, w})
			continue
		} else {
			t := par[s][v] / 2
			queue = append(queue, [2]int{t, v}, [2]int{s ^ t, v})
		}

	}

	for i := 0; i < n; i++ {
		if usedV[i] {
			vs = append(vs, i)
		}
	}
	for i := 0; i < m; i++ {
		if usedE[i] {
			es = append(es, i)
		}
	}
	for _, v := range vs {
		cost += vWeights[v]
	}
	for _, e := range es {
		cost += edges[e][2]
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

type H = [2]int

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Peek() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
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
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}

		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}

// NewUnionFindWithCallback ...
func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	cb(root2, root1)
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
