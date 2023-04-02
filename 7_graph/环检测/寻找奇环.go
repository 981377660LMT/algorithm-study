// https://maspypy.github.io/library/graph/find_odd_cycle.
// (TODO:不一定对)

package main

import "fmt"

const INF int = 1e18

func main() {
	// https://yukicoder.me/problems/no/1436
	n := 3
	edges := [][3]int{{0, 1, 0}, {2, 0, 1}}
	vs, es := FindOddCycle(n, edges, false)
	fmt.Println(vs, es)
}

// 寻找奇环，返回环上的点和边.
// !edges: (u, v, id)
func FindOddCycle(n int, edges [][3]int, directed bool) (vs, es []int) {
	comp := make([]int, n)
	if directed {
		scc := NewStronglyConnectedComponents(n)
		for _, e := range edges {
			scc.AddEdge(e[0], e[1], 1)
		}
		scc.Build()
		comp = scc.CompId
	}

	dist, parent := make([]int, 2*n), make([]int, 2*n)
	for i := range dist {
		dist[i] = INF
		parent[i] = -1
	}
	queue := NewDeque(2 * n)
	add := func(v, d, p int) {
		if dist[v] > d {
			dist[v] = d
			queue.Append(v)
			parent[v] = p
		}
	}

	for root := 0; root < n; root++ {
		if dist[2*root] < INF {
			continue
		}
		if dist[2*root+1] < INF {
			continue
		}
		add(2*root, 0, -1)

		for queue.Size() > 0 {
			v := queue.PopLeft()
			b := v % 2
			for _, e := range edges {
				if comp[e[0]] != comp[e[1]] {
					continue
				}
				w := 2*e[1] + (b ^ 1)
				add(w, dist[v]+1, e[2])
			}
		}

		if dist[2*root+1] == INF {
			continue
		}

		vs = append(vs, root)
		v := 2*root + 1
		for parent[v] != -1 {
			i := parent[v]
			es = append(es, i)
			e := edges[i]
			v = 2*(e[0]+e[1]) + 1 - v
			vs = append(vs, v/2)
		}

		for i, j := 0, len(vs)-1; i < j; i, j = i+1, j-1 {
			vs[i], vs[j] = vs[j], vs[i]
		}
		for i, j := 0, len(es)-1; i < j; i, j = i+1, j-1 {
			es[i], es[j] = es[j], es[i]
		}

		used := make([]int, n)
		for i := range used {
			used[i] = -1
		}
		l, r := -1, -1
		for i, v := range vs {
			if used[v] == -1 {
				used[v] = i
				continue
			}
			l, r = used[v], i
			break
		}

		if l == -1 {
			panic("invalid")
		}
		vs = vs[l:r]
		es = es[l:r]
		return
	}
	return
}

type WeightedEdge struct{ from, to, cost, index int }
type StronglyConnectedComponents struct {
	G      [][]WeightedEdge // 原图
	Dag    [][]int          // 强连通分量缩点后的DAG(有向图邻接表)
	CompId []int            // 每个顶点所属的强连通分量的编号
	Group  [][]int          // 每个强连通分量所包含的顶点
	rg     [][]WeightedEdge
	order  []int
	used   []bool
	eid    int
}

func NewStronglyConnectedComponents(n int) *StronglyConnectedComponents {
	return &StronglyConnectedComponents{G: make([][]WeightedEdge, n)}
}

func (scc *StronglyConnectedComponents) AddEdge(from, to, cost int) {
	scc.G[from] = append(scc.G[from], WeightedEdge{from, to, cost, scc.eid})
	scc.eid++
}

func (scc *StronglyConnectedComponents) Build() {
	scc.rg = make([][]WeightedEdge, len(scc.G))
	for i := range scc.G {
		for _, e := range scc.G[i] {
			scc.rg[e.to] = append(scc.rg[e.to], WeightedEdge{e.to, e.from, e.cost, e.index})
		}
	}

	scc.CompId = make([]int, len(scc.G))
	for i := range scc.CompId {
		scc.CompId[i] = -1
	}
	scc.used = make([]bool, len(scc.G))
	for i := range scc.G {
		scc.dfs(i)
	}
	for i, j := 0, len(scc.order)-1; i < j; i, j = i+1, j-1 {
		scc.order[i], scc.order[j] = scc.order[j], scc.order[i]
	}

	ptr := 0
	for _, v := range scc.order {
		if scc.CompId[v] == -1 {
			scc.rdfs(v, ptr)
			ptr++
		}
	}

	dag := make([][]int, ptr)
	visited := make(map[int]struct{}) // 边去重
	for i := range scc.G {
		for _, e := range scc.G[i] {
			x, y := scc.CompId[e.from], scc.CompId[e.to]
			if x == y {
				continue // 原来的边 x->y 的顶点在同一个强连通分量内,可以汇合同一个 SCC 的权值
			}
			hash := x*len(scc.G) + y
			if _, ok := visited[hash]; !ok {
				dag[x] = append(dag[x], y)
				visited[hash] = struct{}{}
			}
		}
	}
	scc.Dag = dag

	scc.Group = make([][]int, ptr)
	for i := range scc.G {
		scc.Group[scc.CompId[i]] = append(scc.Group[scc.CompId[i]], i)
	}
}

// 获取顶点k所属的强连通分量的编号
func (scc *StronglyConnectedComponents) Get(k int) int {
	return scc.CompId[k]
}

func (scc *StronglyConnectedComponents) dfs(idx int) {
	tmp := scc.used[idx]
	scc.used[idx] = true
	if tmp {
		return
	}
	for _, e := range scc.G[idx] {
		scc.dfs(e.to)
	}
	scc.order = append(scc.order, idx)
}

func (scc *StronglyConnectedComponents) rdfs(idx int, cnt int) {
	if scc.CompId[idx] != -1 {
		return
	}
	scc.CompId[idx] = cnt
	for _, e := range scc.rg[idx] {
		scc.rdfs(e.to, cnt)
	}
}

//
//
type D = int
type Deque struct{ l, r []D }

func NewDeque(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
