// 有向图找到所有极小的环(无弦)
// 图中不存在多重边。
// O(V+E)*len(cycle+1)

package main

func main() {
	C := NewFindAllCycles(5)
	C.AddDirectedEdge(0, 1)
	C.AddDirectedEdge(1, 2)
	C.AddDirectedEdge(2, 3)
	C.AddDirectedEdge(3, 0)
	C.AddDirectedEdge(1, 3)

	C.Run(func(cycle []int) {
		for _, v := range cycle {
			print(v, " ")
		}
		println()
	})
}

type FindAllCycles struct {
	v                  int
	g                  [][]Node
	blockStack         [][]int
	fix, blocked, used []bool
	cb                 func([]int)
	selfCycles         [][]int
}

func NewFindAllCycles(n int) *FindAllCycles {
	res := &FindAllCycles{}
	res.v = n
	res.g = make([][]Node, n)
	for i := 0; i < n; i++ {
		res.g[i] = []Node{{0, 0, -1}}
	}
	res.blockStack = make([][]int, n)
	res.fix = make([]bool, n)
	res.blocked = make([]bool, n)
	res.used = make([]bool, n)
	return res
}

func (fac *FindAllCycles) AddDirectedEdge(u, v int) {
	if u == v {
		fac.selfCycles = append(fac.selfCycles, []int{u})
		return
	}
	fac.g[u][len(fac.g[u])-1].next = len(fac.g[u])
	fac.g[u][0].prev = len(fac.g[u])
	fac.g[u] = append(fac.g[u], Node{len(fac.g[u]) - 1, 0, v})
}

func (fac *FindAllCycles) Run(cb func([]int)) {
	fac.cb = cb
	fac.scc()
	for i := 0; i < fac.v; i++ {
		path := []int{}
		verList := []int{}
		fac.dfs(i, i, &path, &verList)
		fac.fix[i] = true
		for _, j := range verList {
			fac.used[j] = false
			fac.blocked[j] = false
			fac.blockStack[j] = []int{}
		}
	}

	for _, cycle := range fac.selfCycles {
		fac.cb(cycle)
	}
}

func (fac *FindAllCycles) eraseEdge(u, id int) {
	fac.g[u][fac.g[u][id].next].prev = fac.g[u][id].prev
	fac.g[u][fac.g[u][id].prev].next = fac.g[u][id].next
}

func (fac *FindAllCycles) sccDfs(u int, tm, cnt *int, ord, low, cmp []int, st *[]int) {
	ord[u] = *tm
	low[u] = *tm
	*tm++
	*st = append(*st, u)
	for id := fac.g[u][0].next; id != 0; id = fac.g[u][id].next {
		v := fac.g[u][id].to
		if ord[v] < 0 {
			fac.sccDfs(v, tm, cnt, ord, low, cmp, st)
			low[u] = min(low[u], low[v])
		} else if cmp[v] < 0 {
			low[u] = min(low[u], ord[v])
		}
		if cmp[v] >= 0 {
			fac.eraseEdge(u, id)
		}
	}
	if ord[u] == low[u] {
		for {
			v := (*st)[len(*st)-1]
			*st = (*st)[:len(*st)-1]
			cmp[v] = *cnt
			if v == u {
				break
			}
		}
		*cnt++
	}
}

func (fac *FindAllCycles) scc() {
	ord := make([]int, fac.v)
	low := make([]int, fac.v)
	cmp := make([]int, fac.v)
	for i := range ord {
		ord[i] = -1
		cmp[i] = -1
	}
	st := []int{}
	tm := 0
	cnt := 0
	for i := 0; i < fac.v; i++ {
		if ord[i] < 0 {
			fac.sccDfs(i, &tm, &cnt, ord, low, cmp, &st)
		}
	}
}

func (fac *FindAllCycles) unBlock(u int) {
	fac.blocked[u] = false
	for len(fac.blockStack[u]) > 0 {
		v := fac.blockStack[u][len(fac.blockStack[u])-1]
		fac.blockStack[u] = fac.blockStack[u][:len(fac.blockStack[u])-1]
		if fac.blocked[v] {
			fac.unBlock(v)
		}
	}
}

func (fac *FindAllCycles) dfs(u, s int, path, verList *[]int) bool {
	flag := false
	*path = append(*path, u)
	fac.blocked[u] = true
	if !fac.used[u] {
		fac.used[u] = true
		*verList = append(*verList, u)
	}
	for id := fac.g[u][0].next; id != 0; id = fac.g[u][id].next {
		v := fac.g[u][id].to
		if fac.fix[v] {
			fac.eraseEdge(u, id)
			continue
		}
		if v == s {
			fac.cb(*path)
			flag = true
		} else if !fac.blocked[v] {
			if fac.dfs(v, s, path, verList) {
				flag = true
			}
		}
	}

	if flag {
		fac.unBlock(u)
	} else {
		for id := fac.g[u][0].next; id != 0; id = fac.g[u][id].next {
			v := fac.g[u][id].to
			fac.blockStack[v] = append(fac.blockStack[v], u)
		}
	}
	*path = (*path)[:len(*path)-1]
	return flag
}

type Node struct {
	prev, next, to int
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
