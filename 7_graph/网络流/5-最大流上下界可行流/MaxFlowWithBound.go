// api:
// 	NewMaxFlowWithBound(n, start, target int32) *MaxFlowWithBound
// 	(m *MaxFlowWithBound) Add(from, to int32, lo, hi int)
// 	(m *MaxFlowWithBound) Flow() int
// 	(m *MaxFlowWithBound) GetFlowResult() []int // 返回每条边的流量
// 	(m *MaxFlowWithBound) Debug()  					 // 打印边的信息

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, START, END int32
	fmt.Fscan(in, &n, &m, &START, &END)
	START--
	END--
	M := NewMaxFlowWithBound(n, START, END)
	for i := int32(0); i < m; i++ {
		var u, v int32
		var lower, upper int
		fmt.Fscan(in, &u, &v, &lower, &upper)
		u--
		v--
		M.Add(u, v, lower, upper)
	}

	res := M.Flow()
	if res == -1 {
		fmt.Fprintln(out, "No Solution")
	} else {
		fmt.Fprintln(out, res)
	}
}

const INF int = 1e18
const INF32 int32 = 1e9

type rawEdge struct {
	from, to int32
	lo, hi   int
}

type edge struct {
	rev, to   int32
	cap, flow int
}

type MaxFlowWithBound struct {
	n, start, target, source, sink int32
	flowRes                        int
	prepared                       bool
	data                           []rawEdge
	graph                          []edge
	indptr                         []int32
	idx                            []int32
	level, queue, prog             []int32
}

func NewMaxFlowWithBound(n, start, target int32) *MaxFlowWithBound {
	return &MaxFlowWithBound{
		n:      n,
		start:  start,
		target: target,
		source: n,
		sink:   n + 1,
	}
}

func (m *MaxFlowWithBound) Add(from, to int32, lo, hi int) {
	if m.prepared {
		panic("already prepared")
	}
	if !(0 <= from && from < m.n) {
		panic("invalid from")
	}
	if !(0 <= to && to < m.n) {
		panic("invalid to")
	}
	if !(0 <= lo && lo <= hi) {
		panic("invalid bounds")
	}
	m.data = append(m.data, rawEdge{from: from, to: to, lo: lo, hi: hi})
}

func (m *MaxFlowWithBound) Flow() int {
	m.build()
	a := m.flowSt(m.source, m.sink)
	b := m.flowSt(m.source, m.target)
	c := m.flowSt(m.start, m.sink)
	d := m.flowSt(m.start, m.target)
	valid := true
	data, graph, idx := m.data, m.graph, m.idx
	for i := int32(0); i < int32(len(m.data)); i++ {
		lo := data[i].lo
		if lo > 0 && graph[idx[6*i+2]].cap > 0 {
			valid = false
		}
		if lo > 0 && graph[idx[6*i+4]].cap > 0 {
			valid = false
		}
	}
	if !valid {
		m.flowRes = -1
		return -1
	}
	if a+b != a+c || c+d != b+d {
		panic("invalid")
	}
	m.flowRes = c + d
	return c + d
}

func (m *MaxFlowWithBound) GetFlowResult() []int {
	if m.flowRes == -1 {
		panic("no flow")
	}
	res := make([]int, len(m.data))
	data, graph, idx := m.data, m.graph, m.idx
	for i := int32(0); i < int32(len(m.data)); i++ {
		lo, hi := data[i].lo, data[i].hi
		var flow int
		if lo < hi {
			flow = graph[idx[6*i+1]].cap + lo
		} else {
			flow = lo
		}
		res[i] = flow
	}
	return res
}

func (m *MaxFlowWithBound) Debug() {
	for _, e := range m.data {
		fmt.Printf("from: %d, to: %d, lo: %d, hi: %d\n", e.from, e.to, e.lo, e.hi)
	}
}

func (m *MaxFlowWithBound) build() {
	if m.prepared {
		panic("already prepared")
	}
	m.prepared = true
	d := int32(len(m.data))
	m.idx = make([]int32, 6*d)
	for i := int32(0); i < 6*d; i++ {
		m.idx[i] = -1
	}
	cnt := make([]int32, m.n+2)
	data := m.data
	for i := int32(0); i < d; i++ {
		from, to, lo, hi := data[i].from, data[i].to, data[i].lo, data[i].hi
		if from == to {
			continue
		}
		if lo < hi {
			cnt[from]++
			cnt[to]++
		}
		if lo > 0 {
			cnt[m.source]++
			cnt[to]++
			cnt[from]++
			cnt[m.sink]++
		}
	}
	m.indptr = make([]int32, len(cnt)+1)
	for i := int32(0); i < int32(len(cnt)); i++ {
		m.indptr[i+1] = m.indptr[i] + cnt[i]
	}
	size := m.indptr[len(m.indptr)-1]
	m.graph = make([]edge, size)
	m.prog = append(m.indptr[:0:0], m.indptr...)
	add := func(i, j, a, b int32, c int) {
		p, q := m.prog[a], m.prog[b]
		m.prog[a]++
		m.prog[b]++
		m.idx[i] = p
		m.idx[j] = q
		m.graph[p] = edge{rev: q, to: b, cap: c, flow: 0}
		m.graph[q] = edge{rev: p, to: a, cap: 0, flow: 0}
	}
	for i := int32(0); i < d; i++ {
		from, to, lo, hi := data[i].from, data[i].to, data[i].lo, data[i].hi
		if from == to {
			continue
		}
		if lo < hi {
			add(6*i+0, 6*i+1, from, to, hi-lo)
		}
		if lo > 0 {
			add(6*i+2, 6*i+3, m.source, to, lo)
			add(6*i+4, 6*i+5, from, m.sink, lo)
			cnt[m.source]++
			cnt[to]++
			cnt[from]++
			cnt[m.sink]++
		}
	}
}

func (m *MaxFlowWithBound) setLevel(start int32) {
	m.level = make([]int32, m.n+2)
	m.queue = make([]int32, m.n+2)
	for i := int32(0); i < m.n+2; i++ {
		m.level[i] = INF32
	}
	ql, qr := int32(0), int32(0)
	update := func(v, d int32) {
		if m.level[v] > d {
			m.level[v] = d
			m.queue[qr] = v
			qr++
		}
	}
	update(start, 0)
	for ql < qr {
		v := m.queue[ql]
		ql++
		for i := m.indptr[v]; i < m.indptr[v+1]; i++ {
			e := m.graph[i]
			if e.cap > 0 {
				update(e.to, m.level[v]+1)
			}
		}
	}
}

func (m *MaxFlowWithBound) flowDfs(start, target int32) int {
	m.prog = append(m.indptr[:0:0], m.indptr...)
	var dfs func(v int32, lim int) int
	dfs = func(v int32, lim int) int {
		if v == target {
			return lim
		}
		res := 0
		for i := &m.prog[v]; *i < m.indptr[v+1]; *i++ {
			e := &m.graph[*i]
			if e.cap > 0 && m.level[e.to] == m.level[v]+1 {
				a := dfs(e.to, min(lim, e.cap))
				if a == 0 {
					continue
				}
				e.cap -= a
				e.flow += a
				m.graph[e.rev].cap += a
				m.graph[e.rev].flow -= a
				res += a
				lim -= a
				if lim == 0 {
					break
				}
			}
		}
		return res
	}
	return dfs(start, INF)
}

func (m *MaxFlowWithBound) flowSt(start, target int32) int {
	res := 0
	for {
		m.setLevel(start)
		if m.level[target] == INF32 {
			break
		}
		res += m.flowDfs(start, target)
	}
	return res
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
