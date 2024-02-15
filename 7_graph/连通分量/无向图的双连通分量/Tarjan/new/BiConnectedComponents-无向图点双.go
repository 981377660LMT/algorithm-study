package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	// ArticulationPoints()
	// yosupo()
	abc334g()
}

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=GRL_3_A
func ArticulationPoints() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	_, _, isCut := BiConnectedComponent(graph)
	for i := 0; i < n; i++ {
		if isCut[i] {
			fmt.Fprintln(out, i)
		}
	}
}

// https://judge.yosupo.jp/problem/biconnected_components
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	group, _, _ := BiConnectedComponent(graph)
	fmt.Fprintln(out, len(group))
	for _, g := range group {
		fmt.Fprint(out, len(g))
		for _, v := range g {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}

// G - Christmas Color Grid 2
// https://atcoder.jp/contests/abc334/tasks/abc334_g
// 给定一个网格图，由红/绿两种点组成.
// !询问：将每个绿色点变为红色(删除这个点)后，剩余的绿色点连通块的大小.
// ROW,COL<=1000.
//
// 解法1：线段树分治删点(ok)
// 解法2：点双连通分量(good)，按照割点、孤立点和其他三种情况讨论.
func abc334g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const GREEN byte = '#'
	const RED byte = '.'
	const MOD int = 998244353
	pow := func(base, exp, mod int) int {
		base %= mod
		res := 1 % mod
		for ; exp > 0; exp >>= 1 {
			if exp&1 == 1 {
				res = res * base % mod
			}
			base = base * base % mod
		}
		return res
	}

	var ROW, COL int
	fmt.Fscan(in, &ROW, &COL)
	grid := make([][]byte, ROW)
	for i := 0; i < ROW; i++ {
		grid[i] = make([]byte, COL)
		fmt.Fscan(in, &grid[i])
	}

	vid := make([][]int, ROW)
	for i := 0; i < ROW; i++ {
		vid[i] = make([]int, COL)
		for j := 0; j < COL; j++ {
			vid[i][j] = -1
		}
	}
	curId := 0
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			if grid[i][j] == GREEN {
				vid[i][j] = curId
				curId++
			}
		}
	}

	graph := make([][]int, curId)
	uf := NewUnionFindArray(curId)
	rawDeg := make([]int, curId)
	for i := 0; i < ROW; i++ {
		for j := 0; j < COL; j++ {
			if grid[i][j] == GREEN {
				cur := vid[i][j]
				if i+1 < ROW && grid[i+1][j] == GREEN {
					next := vid[i+1][j]
					graph[cur] = append(graph[cur], next)
					graph[next] = append(graph[next], cur)
					rawDeg[cur]++
					rawDeg[next]++
					uf.Union(cur, next)
				}
				if j+1 < COL && grid[i][j+1] == GREEN {
					next := vid[i][j+1]
					graph[cur] = append(graph[cur], next)
					graph[next] = append(graph[next], cur)
					rawDeg[cur]++
					rawDeg[next]++
					uf.Union(cur, next)
				}
			}
		}
	}
	basePart := uf.Part

	_, belong, isCut := BiConnectedComponent(graph)
	res := []int{}
	for i := 0; i < curId; i++ {
		if isCut[i] { // !割点
			res = append(res, basePart+len(belong[i])-1)
		} else if rawDeg[i] == 0 { // !孤立点
			res = append(res, basePart-1)
		} else { // !其他
			res = append(res, basePart)
		}
	}

	sum := 0
	for _, v := range res {
		sum = (sum + v) % MOD
	}
	fmt.Fprintln(out, sum*pow(len(res), MOD-2, MOD)%MOD)
}

// 求无向图的点双连通分量.
// 返回值为 (每个双连通分量包含哪些点, 每个点所属的双连通分量编号(割点有多个), 每个点是否为割点)
// - 原图的割点`至少`在两个不同的 v-BCC 中
// - 原图不是割点的点都`只存在`于一个 v-BCC 中
// - v-BCC 形成的子图内没有割点
func BiConnectedComponent(graph [][]int) (groups [][]int, belong [][]int, isCut []bool) {
	n := len(graph)
	dfsId := int32(0)
	dfsOrder := make([]int32, n)
	vbccId := make([]int32, n)
	idCount := int32(0)
	isCut = make([]bool, n)
	stack := [][2]int32{} // (u,v)

	var dfs func(cur, pre int32) int32
	dfs = func(cur, pre int32) int32 {
		dfsId++
		dfsOrder[cur] = dfsId
		curLow := dfsId
		childCount := 0
		for _, next := range graph[cur] {
			if dfsOrder[next] == 0 {
				stack = append(stack, [2]int32{cur, int32(next)})
				childCount++
				nextLow := dfs(int32(next), cur)
				if nextLow >= dfsOrder[cur] {
					isCut[cur] = true
					idCount++
					group := []int{}
					for {
						topEdge := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						v1, v2 := topEdge[0], topEdge[1]
						if vbccId[v1] != idCount {
							vbccId[v1] = idCount
							group = append(group, int(v1))
						}
						if vbccId[v2] != idCount {
							vbccId[v2] = idCount
							group = append(group, int(v2))
						}
						if v1 == cur && v2 == int32(next) {
							break
						}
					}
					groups = append(groups, group)
				}
				if nextLow < curLow {
					curLow = nextLow
				}
			} else if next != int(pre) && dfsOrder[next] < dfsOrder[cur] {
				stack = append(stack, [2]int32{cur, int32(next)})
				if dfsOrder[next] < curLow {
					curLow = dfsOrder[next]
				}
			}
		}
		if pre == -1 && childCount == 1 {
			isCut[cur] = false
		}
		return curLow
	}

	for i, v := range dfsOrder {
		if v == 0 {
			if len(graph[i]) == 0 { // 零度，即孤立点（isolated vertex）
				idCount++
				vbccId[i] = idCount
				groups = append(groups, []int{i})
				continue
			}
			dfs(int32(i), -1)
		}
	}

	belong = make([][]int, n)
	for i, group := range groups {
		for _, v := range group {
			belong[v] = append(belong[v], i)
		}
	}
	return
}

// 边双缩点成树.
// !BCC 和割点作为新图中的节点，并在每个割点与包含它的所有 BCC 之间连边
// !bcc1 - 割点1 - bcc2 - 割点2 - ...
func ToTree(graph [][]int, groups [][]int, isCut []bool) (tree [][]int) {
	n := len(graph)
	idCount := len(groups)
	cutId := make([]int, n)
	for i, v := range isCut {
		if v {
			cutId[i] = idCount
			idCount++
		}
	}
	tree = make([][]int, idCount)
	for cur, group := range groups {
		for _, g := range group {
			if isCut[g] {
				next := cutId[g]
				tree[cur] = append(tree[cur], next)
				tree[next] = append(tree[next], cur)
			}
		}
	}
	return
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int
	n    int
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{
		Part: n,
		n:    n,
		data: data,
	}
}

// 按秩合并.
func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	if cb != nil {
		cb(root1, root2)
	}
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetSize(key int) int {
	return -ufa.data[ufa.Find(key)]
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
