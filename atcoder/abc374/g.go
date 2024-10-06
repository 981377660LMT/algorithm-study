package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"sort"
	"strconv"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N := int32(io.NextInt())
	words := make([]string, N)
	for i := int32(0); i < N; i++ {
		words[i] = io.Text()
	}

	starts := make([][]int32, 26)
	for i := int32(0); i < N; i++ {
		starts[words[i][0]-'A'] = append(starts[words[i][0]-'A'], i)
	}

	adjList := make([][]int32, N)
	for i := int32(0); i < N; i++ {
		endV := words[i][len(words[i])-1] - 'A'
		adjList[i] = starts[endV]
	}

	groups, belong := FindScc(N, adjList)
	dag, _ := ToDag(adjList, groups, belong, nil)

	m := int32(len(dag))
	adjMatrix := make([][]bool, m)
	for i := range adjMatrix {
		adjMatrix[i] = make([]bool, m)
	}
	for i := int32(0); i < m; i++ {
		for _, to := range dag[i] {
			adjMatrix[i][to] = true
		}
	}

	for k := int32(0); k < m; k++ {
		for i := int32(0); i < m; i++ {
			for j := int32(0); j < m; j++ {
				adjMatrix[i][j] = adjMatrix[i][j] || (adjMatrix[i][k] && adjMatrix[k][j])
			}
		}
	}

	dag2 := make([][]int, m)
	for i := int32(0); i < m; i++ {
		for j := int32(0); j < m; j++ {
			if i != j && adjMatrix[i][j] {
				dag2[i] = append(dag2[i], int(j))
			}
		}
	}

	res := MaxAntiChain(len(dag2), dag2)
	io.Println(len(res))
}

const INF int = 1e18

func FindScc(n int32, graph [][]int32) (groups [][]int32, belong []int32) {
	dfsOrder := make([]int32, n)
	dfsId := int32(0)
	stack := []int32{}
	inStack := make([]bool, n)

	var dfs func(int32) int32
	dfs = func(cur int32) int32 {
		dfsId++
		dfsOrder[cur] = dfsId
		curLow := dfsId
		stack = append(stack, cur)
		inStack[cur] = true
		for _, next := range graph[cur] {
			if dfsOrder[next] == 0 {
				nextLow := dfs(next)
				if nextLow < curLow {
					curLow = nextLow
				}
			} else if inStack[next] && dfsOrder[next] < curLow {
				curLow = dfsOrder[next]
			}
		}
		if dfsOrder[cur] == curLow {
			group := []int32{}
			for {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inStack[top] = false
				group = append(group, top)
				if top == cur {
					break
				}
			}
			groups = append(groups, group)
		}
		return curLow
	}

	for i, order := range dfsOrder {
		if order == 0 {
			dfs(int32(i))
		}
	}

	for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
		groups[i], groups[j] = groups[j], groups[i]
	}
	belong = make([]int32, n)
	for i := int32(0); i < int32(len(groups)); i++ {
		for _, v := range groups[i] {
			belong[v] = i
		}
	}
	return
}

func ToDag(
	graph [][]int32, groups [][]int32, belong []int32,
	forEachEdge func(from, fromId, to, toId int32),
) (dag [][]int32, indeg []int32) {
	m := int32(len(groups))
	dag = make([][]int32, m)
	visitedEdge := map[int]struct{}{}
	indeg = make([]int32, m)
	for cur, nexts := range graph {
		curId := belong[cur]
		for _, next := range nexts {
			nextId := belong[next]
			if curId != nextId {
				hash := int(curId)*int(m) + int(nextId)
				if _, ok := visitedEdge[hash]; ok {
					continue
				}
				visitedEdge[hash] = struct{}{}
				dag[curId] = append(dag[curId], nextId)
				indeg[nextId]++
			}
			if forEachEdge != nil {
				forEachEdge(int32(cur), curId, next, nextId)
			}
		}
	}
	return
}

type BipartiteMatching struct {
	n       int32
	graph   [][]int
	color   []int8
	dist    []int32
	match   []int32
	visited []bool
}

func NewBipartiteMatching(graph [][]int, colors []int8) *BipartiteMatching {
	n := len(graph)
	bm := &BipartiteMatching{
		n:       int32(n),
		graph:   graph,
		dist:    make([]int32, n),
		match:   make([]int32, n),
		visited: make([]bool, n),
	}
	if colors != nil {
		bm.color = colors
	} else {
		bm.color = BipartiteVertexColoring(n, graph)
	}
	if n > 0 && len(bm.color) == 0 {
		panic("not bipartite graph")
	}
	for i := range bm.dist {
		bm.dist[i] = -1
	}
	for i := range bm.match {
		bm.match[i] = -1
	}
	for {
		bm.bfs()
		for i := range bm.visited {
			bm.visited[i] = false
		}
		flow := 0
		for v := int32(0); v < bm.n; v++ {
			if bm.color[v] == 0 && bm.match[v] == -1 && bm.dfs(v) {
				flow++
			}
		}
		if flow == 0 {
			break
		}
	}
	return bm
}

func (bm *BipartiteMatching) MaxMatching() (res [][2]int) {
	for v := int32(0); v < bm.n; v++ {
		if v < bm.match[v] {
			res = append(res, [2]int{int(v), int(bm.match[v])})
		}
	}
	return
}

func (bm *BipartiteMatching) MinVertexCover() (res []int) {
	for v := int32(0); v < bm.n; v++ {
		if (bm.color[v] != 0) != (bm.dist[v] == -1) {
			res = append(res, int(v))
		}
	}
	return
}

func (bm *BipartiteMatching) MaxIndependentSet() (res []int) {
	for v := int32(0); v < bm.n; v++ {
		if (bm.color[v] != 0) == (bm.dist[v] == -1) {
			res = append(res, int(v))
		}
	}
	return
}

func (bm *BipartiteMatching) MinEdgeCover(edges [][2]int) (res []int) {
	done := make([]bool, bm.n)
	for ei, e := range edges {
		u, v := e[0], e[1]
		if done[u] || done[v] {
			continue
		}
		if bm.match[u] == int32(v) {
			res = append(res, ei)
			done[u] = true
			done[v] = true
		}
	}
	for ei, e := range edges {
		u, v := e[0], e[1]
		if !done[u] {
			res = append(res, ei)
			done[u] = true
		}
		if !done[v] {
			res = append(res, ei)
			done[v] = true
		}
	}
	sort.Ints(res)
	return
}

func (bm *BipartiteMatching) Debug() {
	fmt.Println("match", bm.match)
	fmt.Println("MinVertexCoverr", bm.MinVertexCover())
	fmt.Println("MaxIndependentSet", bm.MaxIndependentSet())
}

func (bm *BipartiteMatching) DMDecomposition(edges [][2]int) (compCount int, belong []int) {
	belong = make([]int, bm.n)
	for i := range belong {
		belong[i] = -1
	}
	queue := []int{}
	add := func(v, x int) {
		if belong[v] == -1 {
			belong[v] = x
			queue = append(queue, v)
		}
	}
	for v := 0; v < int(bm.n); v++ {
		if bm.match[v] == -1 && bm.color[v] == 0 {
			add(v, 0)
		}
	}
	for v := 0; v < int(bm.n); v++ {
		if bm.match[v] == -1 && bm.color[v] == 1 {
			add(v, INF)
		}
	}
	for len(queue) > 0 {
		v := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if bm.match[v] != -1 {
			add(int(bm.match[v]), belong[v])
		}
		if bm.color[v] == 0 && belong[v] == 0 {
			for _, to := range bm.graph[v] {
				add(to, belong[v])
			}
		}
		if bm.color[v] == 1 && belong[v] == INF {
			for _, to := range bm.graph[v] {
				add(to, belong[v])
			}
		}
	}

	vs := []int{}
	for v := 0; v < int(bm.n); v++ {
		if belong[v] == -1 {
			vs = append(vs, v)
		}
	}
	m := len(vs)
	dg := make([][]int, m)
	for i := range dg {
		v := vs[i]
		if bm.match[v] != -1 {
			j := sort.SearchInts(vs, int(bm.match[v]))
			dg[i] = append(dg[i], j)
		}
		if bm.color[v] == 0 {
			for _, to := range bm.graph[v] {
				if belong[to] != -1 || to == int(bm.match[v]) {
					continue
				}
				j := sort.SearchInts(vs, to)
				dg[i] = append(dg[i], j)
			}
		}
	}

	compCount, comp := StronglyConnectedComponent(dg)
	compCount++

	for i := 0; i < m; i++ {
		belong[vs[i]] = 1 + comp[i]
	}
	for v := 0; v < int(bm.n); v++ {
		if belong[v] == INF {
			belong[v] = compCount
		}
	}
	return
}

func (bm *BipartiteMatching) bfs() {
	for i := range bm.dist {
		bm.dist[i] = -1
	}
	queue := []int32{}
	for v := int32(0); v < bm.n; v++ {
		if bm.color[v] == 0 && bm.match[v] == -1 {
			queue = append(queue, v)
			bm.dist[v] = 0
		}
	}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, to := range bm.graph[v] {
			bm.dist[to] = 0
			w := bm.match[to]
			if w != -1 && bm.dist[w] == -1 {
				bm.dist[w] = bm.dist[v] + 1
				queue = append(queue, w)
			}
		}
	}
}

func (bm *BipartiteMatching) dfs(v int32) bool {
	bm.visited[v] = true
	for _, to := range bm.graph[v] {
		w := bm.match[to]
		if w == -1 || (!bm.visited[w] && bm.dist[w] == bm.dist[v]+1 && bm.dfs(w)) {
			bm.match[to] = v
			bm.match[v] = int32(to)
			return true
		}
	}
	return false
}

func BipartiteVertexColoring(n int, graph [][]int) (colors []int8) {
	uf := NewUf(2 * n)
	for cur, nexts := range graph {
		for _, next := range nexts {
			if cur < next {
				uf.Union(cur+n, next)
				uf.Union(cur, next+n)
			}
		}
	}
	colors = make([]int8, 2*n)
	for i := range colors {
		colors[i] = -1
	}
	for v := 0; v < n; v++ {
		if root := uf.Find(v); root == v && colors[root] < 0 {
			colors[root] = 0
			colors[uf.Find(v+n)] = 1
		}
	}
	for v := 0; v < n; v++ {
		colors[v] = colors[uf.Find(v)]
	}
	colors = colors[:n]
	for v := 0; v < n; v++ {
		if uf.Find(v) == uf.Find(v+n) {
			return nil
		}
	}
	return
}

type Uf struct {
	data []int
}

func NewUf(n int) *Uf {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &Uf{data: data}
}

func (ufa *Uf) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *Uf) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func StronglyConnectedComponent(graph [][]int) (compCount int, belong []int) {
	n32 := int32(len(graph))
	compId := int32(0)
	comp := make([]int32, n32)
	low := make([]int32, n32)
	ord := make([]int32, n32)
	for i := range ord {
		ord[i] = -1
	}
	path := []int32{}
	now := int32(0)

	var dfs func(int32)
	dfs = func(v int32) {
		low[v] = now
		ord[v] = now
		now++
		path = append(path, v)
		for _, to := range graph[v] {
			if ord[to] == -1 {
				dfs(int32(to))
				if low[v] > low[to] {
					low[v] = low[to]
				}
			} else if low[v] > ord[to] {
				low[v] = ord[to]
			}
		}
		if low[v] == ord[v] {
			for {
				u := path[len(path)-1]
				path = path[:len(path)-1]
				ord[u] = n32
				comp[u] = compId
				if u == v {
					break
				}
			}
			compId++
		}
	}

	for v := int32(0); v < n32; v++ {
		if ord[v] == -1 {
			dfs(v)
		}
	}

	compCount = int(compId)
	belong = make([]int, n32)
	for v := int32(0); v < n32; v++ {
		belong[v] = compCount - 1 - int(comp[v])
	}
	return
}

func SccDag(graph [][]int, compCount int, belong []int) (dag [][]int) {
	unique := func(nums []int32) []int32 {
		set := make(map[int32]struct{})
		for _, v := range nums {
			set[v] = struct{}{}
		}
		res := make([]int32, 0, len(set))
		for k := range set {
			res = append(res, k)
		}
		return res
	}

	edges := make([][]int32, compCount)
	for cur, nexts := range graph {
		curComp := belong[cur]
		for _, next := range nexts {
			nextComp := belong[next]
			if curComp != nextComp {
				edges[curComp] = append(edges[curComp], int32(nextComp))
			}
		}
	}

	dag = make([][]int, compCount)
	for cur := 0; cur < compCount; cur++ {
		edges[cur] = unique(edges[cur])
		for _, next := range edges[cur] {
			dag[cur] = append(dag[cur], int(next))
		}
	}

	return
}

func MaxAntiChain(n int, dag [][]int) []int {
	newGraph := make([][]int, n+n)
	for i := 0; i < n; i++ {
		for _, to := range dag[i] {
			newGraph[i] = append(newGraph[i], to+n)
		}
	}
	bm := NewBipartiteMatching(newGraph, nil)
	cover := bm.MinVertexCover()
	ok := make([]bool, n)
	for i := range ok {
		ok[i] = true
	}
	for _, v := range cover {
		ok[v%n] = false
	}
	antichain := []int{}
	for v := 0; v < n; v++ {
		if ok[v] {
			antichain = append(antichain, v)
		}
	}
	return antichain
}

type Edge = [2]int

func Dijkstra(n int, adjList [][]Edge, start int) (dist, preV []int) {
	type pqItem struct{ node, dist int }
	dist = make([]int, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	preV = make([]int, n)
	for i := range preV {
		preV[i] = -1
	}

	pq := nhp(func(a, b H) int {
		return a.(pqItem).dist - b.(pqItem).dist
	}, nil)
	pq.Push(pqItem{start, 0})

	for pq.Len() > 0 {
		curNode := pq.Pop().(pqItem)
		cur, curDist := curNode.node, curNode.dist
		if curDist > dist[cur] {
			continue
		}

		for _, edge := range adjList[cur] {
			next, weight := edge[0], edge[1]
			if cand := curDist + weight; cand < dist[next] {
				dist[next] = cand
				preV[next] = cur
				pq.Push(pqItem{next, cand})
			}
		}
	}

	return
}

type H = interface{}

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
