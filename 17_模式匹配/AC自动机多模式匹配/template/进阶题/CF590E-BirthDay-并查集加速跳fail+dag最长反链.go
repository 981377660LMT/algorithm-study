// CF590E-BirthDay-并查集加速跳fail+dag最长反链

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const INF int = 1e18

// Birthday
// https://www.luogu.com.cn/problem/CF590E
// 给定n个只包含a,b的字符串，求出最大的集合使得没有任意两个字符串是包含关系
// n<=750,sum(s)<=1e7
// !总时间复杂度O(n^3+O(sum(s)))
//
// 子串关系是一个偏序关系.需要利用传递闭包求出所有子串关系.。
// !1.如何快速求出所有子串的包含关系?
// !对字符的每个前缀，都记录下来它经过跳fail链能够跳到的第一个终止节点(非自身)是哪一个, 连一条有向边.
// !上跳fail过程中需要快速跳过不合法的节点，因此需要并查集加速.
// !求出每个串的转移后，利用传递闭包求出所有子串包含关系(这个过程可以bitset加速).
// !2.求dag最大反链方案.
// 拆点，转换为最大流问题.
func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)

	acm := NewAC()
	words := make([]string, n)
	posId := make(map[int]int)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		pos := acm.AddString(s)
		posId[pos] = i
		words[i] = s
	}
	acm.BuildSuffixLink(true)

	uf := NewUF(acm.Size())

	// !每个对每个字符的前缀，都记录下来它经过跳fail链能够跳到的第一个终止节点是哪一个.
	trans := NewTransitiveClosure(n)
	for i := 0; i < n; i++ {
		pos := 0
		for _, s := range words[i] {
			pos = int(acm.children[pos][int32(s)-OFFSET])
			// 这一段需要加速，不重复计算(即如果为空，则之后都要跳过这个节点)
			// 用并查集加速
			for cur := pos; cur != 0; cur = uf.Find(int(acm.suffixLink[cur])) {
				id, ok := posId[cur]
				if !ok {
					uf.Union(int(acm.suffixLink[cur]), cur, nil)
				}
				if ok && id != i {
					trans.AddDirectedEdge(i, id)
					break
				}
			}
		}
	}
	trans.Build()

	edges := make([][2]int, 0)
	trans.EnumerateEdges(func(from, to int) {
		edges = append(edges, [2]int{from, to})
	})
	S := NewMaxAntiChainSolver(n, edges)
	maxAntiChain := S.MaxAntiChain()
	fmt.Fprintln(out, len(maxAntiChain))
	for _, v := range maxAntiChain {
		fmt.Fprint(out, v+1, " ")
	}
}

const SIGMA int32 = 2
const OFFSET int32 = 'a'

// 不调用 BuildSuffixLink 就是Trie，调用 BuildSuffixLink 就是AC自动机.
// 每个状态对应Trie中的一个结点，也对应一个字符串.
type AC struct {
	children           [][SIGMA]int32 // children[v][c] 表示节点v通过字符c转移到的节点.
	suffixLink         []int32        // 又叫fail.指向当前trie节点(对应一个前缀)的最长真后缀对应结点，例如"bc"是"abc"的最长真后缀.
	bfsOrder           []int32        // 结点的拓扑序,0表示虚拟节点.
	needUpdateChildren bool           // 是否需要更新children数组.
}

func NewAC() *AC {
	res := &AC{}
	res.newNode()
	return res
}

// 添加一个字符串，返回最后一个字符对应的节点编号.
func (trie *AC) AddString(str string) int {
	if len(str) == 0 {
		return 0
	}
	pos := 0
	for _, s := range str {
		ord := int32(s) - OFFSET
		if trie.children[pos][ord] == -1 {
			trie.children[pos][ord] = trie.newNode()
		}
		pos = int(trie.children[pos][ord])
	}

	return pos
}

// pos: DFA的状态集, ord: DFA的字符集
func (trie *AC) Move(pos int, ord int) int {
	ord -= int(OFFSET)
	if trie.needUpdateChildren {
		return int(trie.children[pos][ord])
	}
	for {
		nexts := trie.children[pos]
		if nexts[ord] != -1 {
			return int(nexts[ord])
		}
		if pos == 0 {
			return 0
		}
		pos = int(trie.suffixLink[pos])
	}
}

// 自动机中的节点(状态)数量，包括虚拟节点0.
func (trie *AC) Size() int {
	return len(trie.children)
}

func (trie *AC) Empty() bool {
	return len(trie.children) == 1
}

// 构建后缀链接(失配指针).
// needUpdateChildren 表示是否需要更新children数组(连接trie图).
//
// !move调用较少时，设置为false更快.
func (trie *AC) BuildSuffixLink(needUpdateChildren bool) {
	trie.needUpdateChildren = needUpdateChildren
	trie.suffixLink = make([]int32, len(trie.children))
	for i := range trie.suffixLink {
		trie.suffixLink[i] = -1
	}
	trie.bfsOrder = make([]int32, len(trie.children))
	head, tail := 0, 0
	trie.bfsOrder[tail] = 0
	tail++
	for head < tail {
		v := trie.bfsOrder[head]
		head++
		for i, next := range trie.children[v] {
			if next == -1 {
				continue
			}
			trie.bfsOrder[tail] = next
			tail++
			f := trie.suffixLink[v]
			for f != -1 && trie.children[f][i] == -1 {
				f = trie.suffixLink[f]
			}
			trie.suffixLink[next] = f
			if f == -1 {
				trie.suffixLink[next] = 0
			} else {
				trie.suffixLink[next] = trie.children[f][i]
			}
		}
	}
	if !needUpdateChildren {
		return
	}
	for _, v := range trie.bfsOrder {
		for i, next := range trie.children[v] {
			if next == -1 {
				f := trie.suffixLink[v]
				if f == -1 {
					trie.children[v][i] = 0
				} else {
					trie.children[v][i] = trie.children[f][i]
				}
			}
		}
	}
}

// 按照拓扑序进行转移(EnumerateFail).
func (acm *AC) Dp(f func(from, to int)) {
	for _, v := range acm.bfsOrder {
		if v != 0 {
			f(int(acm.suffixLink[v]), int(v))
		}
	}
}

func (acm *AC) BuildFailTree() [][]int {
	res := make([][]int, acm.Size())
	acm.Dp(func(pre, cur int) {
		res[pre] = append(res[pre], cur)
	})
	return res
}

func (trie *AC) newNode() int32 {
	nexts := [SIGMA]int32{-1, -1}
	trie.children = append(trie.children, nexts)
	return int32(len(trie.children) - 1)
}

type MaxAntiChainSolver struct {
	n            int
	source, sink int
	maxFlow      *MaxFlowAtcoder
	calculated   bool
}

func NewMaxAntiChainSolver(n int, edges [][2]int) *MaxAntiChainSolver {
	source, sink := 2*n, 2*n+1
	mf := NewMaxFlowAtcoder(2*n + 2)
	for _, e := range edges {
		u, v := e[0], e[1]
		mf.AddEdge(u+n, v, INF)
	}
	for i := 0; i < n; i++ {
		mf.AddEdge(source, i+n, 1)
		mf.AddEdge(i, sink, 1)
		mf.AddEdge(i, i+n, INF)
	}
	return &MaxAntiChainSolver{n: n, source: source, sink: sink, maxFlow: mf}
}

// dag最长反链.
func (solver *MaxAntiChainSolver) MaxAntiChain() []int {
	solver.calFlow()
	isCut := solver.maxFlow.MinCut(solver.source)
	res := make([]int, 0)
	for i := 0; i < solver.n; i++ {
		if !isCut[i] && isCut[i+solver.n] {
			res = append(res, i)
		}
	}
	return res
}

// 每个点能否在最长反链上.
func (solver *MaxAntiChainSolver) OnAntiChain() []bool {
	solver.calFlow()
	adjList := make([][]int, solver.n*2+2)
	for _, e := range solver.maxFlow.GetEdges() {
		if e.cap != e.flow {
			adjList[e.from] = append(adjList[e.from], e.to)
		}
		if e.flow != 0 {
			adjList[e.to] = append(adjList[e.to], e.from)
		}
	}

	reach := make([][]bool, solver.n*2+2)
	for i := 0; i < solver.n*2+2; i++ {
		reach[i] = make([]bool, solver.n*2+2)
	}
	for i := 0; i < solver.n*2+2; i++ {
		queue := []int{i}
		reach[i][i] = true
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, to := range adjList[cur] {
				if !reach[i][to] {
					reach[i][to] = true
					queue = append(queue, to)
				}
			}
		}
	}

	res := make([]bool, solver.n)
	for i := 0; i < solver.n; i++ {
		res[i] = !reach[i+solver.n][solver.sink] && !reach[i+solver.n][i] && !reach[solver.source][i]
	}
	return res
}

func (solver *MaxAntiChainSolver) calFlow() {
	if solver.calculated {
		return
	}
	solver.calculated = true
	solver.maxFlow.Flow(solver.source, solver.sink)
}

type Edge struct{ from, to, cap, flow int }

type _edge struct {
	to, rev int32
	cap     int
}

type MaxFlowAtcoder struct {
	n   int32
	pos [][2]int32
	g   [][]_edge
}

// https://github.com/atcoder/ac-library/blob/master/atcoder/maxflow.hpp
func NewMaxFlowAtcoder(n int) *MaxFlowAtcoder {
	return &MaxFlowAtcoder{
		n: int32(n),
		g: make([][]_edge, n),
	}
}

// 添加一条从from到to的容量为cap的边，返回边的编号.
func (mf *MaxFlowAtcoder) AddEdge(from, to int, cap int) int {
	m := len(mf.pos)
	mf.pos = append(mf.pos, [2]int32{int32(from), int32(len(mf.g[from]))})
	fromId := int32(len(mf.g[from]))
	toId := int32(len(mf.g[to]))
	if from == to {
		toId++
	}
	mf.g[from] = append(mf.g[from], _edge{int32(to), toId, cap})
	mf.g[to] = append(mf.g[to], _edge{int32(from), fromId, 0})
	return m
}

func (mf *MaxFlowAtcoder) GetEdge(i int) Edge {
	first, second := mf.pos[i][0], mf.pos[i][1]
	e := mf.g[first][second]
	re := mf.g[e.to][e.rev]
	return Edge{int(first), int(e.to), e.cap + re.cap, re.cap}
}

func (mf *MaxFlowAtcoder) GetEdges() []Edge {
	m := len(mf.pos)
	res := make([]Edge, 0, m)
	for i := 0; i < m; i++ {
		res = append(res, mf.GetEdge(i))
	}
	return res
}

func (mf *MaxFlowAtcoder) ChangeEdge(i int, newCap, newFlow int) {
	e := &mf.g[mf.pos[i][0]][mf.pos[i][1]]
	re := &mf.g[e.to][e.rev]
	e.cap = newCap - newFlow
	re.cap = newFlow
}

func (mf *MaxFlowAtcoder) Flow(s, t int) int {
	return mf.FlowWithLimit(s, t, INF)
}

func (mf *MaxFlowAtcoder) FlowWithLimit(s, t int, flowLimit int) int {
	level := make([]int32, mf.n)
	iter := make([]int32, mf.n)
	queue := make([]int32, 0, mf.n)
	s32, t32 := int32(s), int32(t)

	bfs := func() {
		for i := range level {
			level[i] = -1
		}
		level[s] = 0
		queue = queue[:0]
		queue = append(queue, s32)
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for _, e := range mf.g[v] {
				if e.cap == 0 || level[e.to] >= 0 {
					continue
				}
				level[e.to] = level[v] + 1
				if e.to == t32 {
					return
				}
				queue = append(queue, e.to)
			}
		}
	}

	var dfs func(int32, int) int
	dfs = func(v int32, up int) int {
		if v == s32 {
			return up
		}
		res := 0
		levelV := level[v]
		for i := &iter[v]; *i < int32(len(mf.g[v])); *i++ {
			e := &mf.g[v][*i]
			if levelV <= level[e.to] || mf.g[e.to][e.rev].cap == 0 {
				continue
			}
			d := dfs(e.to, min(up-res, mf.g[e.to][e.rev].cap))
			if d <= 0 {
				continue
			}
			e.cap += d
			mf.g[e.to][e.rev].cap -= d
			res += d
			if res == up {
				return res
			}
		}
		level[v] = mf.n
		return res
	}

	flow := 0
	for flow < flowLimit {
		bfs()
		if level[t] == -1 {
			break
		}
		for i := range iter {
			iter[i] = 0
		}
		for {
			f := dfs(t32, flowLimit-flow)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

// 返回剩余网络中从源点s可以到达的所有节点的集合。
func (mf *MaxFlowAtcoder) MinCut(s int) (visited []bool) {
	visited = make([]bool, mf.n)
	q := []int32{int32(s)}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		visited[p] = true
		for _, e := range mf.g[p] {
			if e.cap != 0 && !visited[e.to] {
				visited[e.to] = true
				q = append(q, e.to)
			}
		}
	}
	return visited
}

// 有向图的传递闭包.
type TransitiveClosure struct {
	n        int
	canReach []_BitSet64
	hasBuilt bool
}

func NewTransitiveClosure(n int) *TransitiveClosure {
	canReach := make([]_BitSet64, n)
	for i := range canReach {
		canReach[i] = NewBitset(n)
	}
	return &TransitiveClosure{n: n, canReach: canReach}
}

func (tc *TransitiveClosure) AddDirectedEdge(from, to int) {
	tc.hasBuilt = false
	tc.canReach[from].Set(to)
}

func (tc *TransitiveClosure) Build() {
	if tc.hasBuilt {
		return
	}
	tc.hasBuilt = true
	n, canReach := tc.n, tc.canReach
	for k := 0; k < n; k++ {
		cacheK := canReach[k]
		for i := 0; i < n; i++ {
			cacheI := canReach[i]
			if cacheI.Has(k) {
				cacheI.IOr(cacheK)
			}
		}
	}
}

func (tc *TransitiveClosure) CanReach(from, to int) bool {
	if !tc.hasBuilt {
		tc.Build()
	}
	return tc.canReach[from].Has(to)
}

func (tc *TransitiveClosure) EnumerateEdges(f func(from, to int)) {
	if !tc.hasBuilt {
		tc.Build()
	}
	for from, bs := range tc.canReach {
		bs.ForEach(func(to int) {
			f(from, to)
		})
	}
}

type _BitSet64 []uint64

func NewBitset(n int) _BitSet64 { return make(_BitSet64, n>>6+1) }

func (b _BitSet64) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 }
func (b _BitSet64) Set(p int)      { b[p>>6] |= 1 << (p & 63) }

func (b _BitSet64) IOr(c _BitSet64) _BitSet64 {
	for i, v := range c {
		b[i] |= v
	}
	return b
}

// 遍历所有 1 的位置.
func (b _BitSet64) ForEach(f func(pos int)) {
	for i, v := range b {
		for ; v != 0; v &= v - 1 {
			j := (i << 6) | _lowbit(v)
			f(j)
		}
	}
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 0, 2)
func _lowbit(x uint64) int {
	if x == 0 {
		return -1
	}
	return bits.TrailingZeros64(x)
}

type UF struct {
	Part int
	n    int
	data []int32
}

func NewUF(n int) *UF {
	data := make([]int32, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UF{
		Part: n,
		n:    n,
		data: data,
	}
}

// 将child结点合并到parent结点上,返回是否合并成功.
func (ufa *UF) Union(parent, child int, f func(parentRoot, childRoot int)) bool {
	parent, child = ufa.Find(parent), ufa.Find(child)
	if parent == child {
		return false
	}
	ufa.data[parent] += ufa.data[child]
	ufa.data[child] = int32(parent)
	ufa.Part--
	if f != nil {
		f(parent, child)
	}
	return true
}

func (ufa *UF) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = int32(ufa.Find(int(ufa.data[key])))
	return int(ufa.data[key])
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
