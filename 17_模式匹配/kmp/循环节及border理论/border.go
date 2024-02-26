// 循环节及border理论
// 循环节:
// 一个字符串的循环节是指一个非空字符串, 使得原字符串是由循环节重复若干次(>=2)得到的.
// 例如, "ababab"的循环节是"ab".
// https://www.cnblogs.com/alex-wei/p/Common_String_Theory_Theory.html

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// P4391()
	P5829()

	// CF526D()

	// acwing143()
}

// 459. 重复的子字符串(是否存在循环节)
// https://leetcode.cn/problems/repeated-substring-pattern/description/
// 给定一个非空的字符串 s ，检查是否可以通过由它的一个子串重复多次构成。
func repeatedSubstringPattern(s string) bool {
	next := GetNext(len(s), func(i int) int { return int(s[i]) })
	n := len(s)
	return Period(next, n-1) > 0
}

// P3538 [POI2012] OKR-A Horrible Poem
// https://www.luogu.com.cn/problem/P3538

// P4391 [BOI2009] Radio Transmission 无线传输
// https://www.luogu.com.cn/problem/P4391
// !给定字符串s，需要找到一个尽可能短的前缀p满足s是pp..pp的子串.
// 答案为 n-next[n-1]，其中next为s的next数组.
func P4391() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)
	next := GetNext(n, func(i int) int { return int(s[i]) })
	fmt.Fprintln(out, n-next[n-1])
}

// P5829 【模板】失配树 (求解两个前缀的最长公共 border)
// https://www.luogu.com.cn/problem/P5829
// 在 KMP 算法中注意到 next[i]<i，因此若从 next[i] 向 i 连边，我们最终会得到一棵有根树。这就是失配树。
// 失配树有很好的性质：对于树上两个点a,b,若a是b的祖先，则s[:a+1]是s[:b+1]的border。这一点由 next 的性质可以得到.
// 因此，若需要查询u前缀和v前缀的最长公共 border, 只需要查询u和v在失配树上的 LCA 即可.
//
//   - aaaabbabbaa
//   - [[1 5 6 8 9] [2 7 10] [3] [4] [] [] [] [] [] [] []]
//   - 0:
//     1:a
//     2:aa
//     3:aaa
//     4:aaaa
//     5:aaaab
//     6:aaaabb
//     7:aaaabba
//     8:aaaabbab
//     9:aaaabbabb
//     10:aaaabbabba
//     11:aaaabbabbaa
//
// - 失配树(failTree)，每个点指向当前前缀的最长真后缀，深度表示当前字符串长度.
//
//	         0
//	     / / | \ \
//	    1  5 6 8  9
//	  / | \
//	 2  7  10
//	/ \
//	3  11
//	|
//	4
func P5829() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var q int
	fmt.Fscan(in, &q)

	n := len(s)
	next := GetNext(n, func(i int) int { return int(s[i]) })
	tree := make([][]int, n+1) // 结点i表示前缀s[:i]
	for i := 1; i <= n; i++ {
		p := next[i-1]
		tree[p] = append(tree[p], i)
	}

	L := NewLCA(tree, []int{0})
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)

		// !s[:a]和s[:b]的最长公共 border
		if a == 0 || b == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		lca := L.LCA(a, b)
		if lca == a || lca == b {
			lca = next[lca-1] // !Border 不能是自己，所以当 LCA 是两数之一时，要再往上跳一格
		}
		fmt.Fprintln(out, lca)
	}
}

// Om Nom and Necklace
// https://www.luogu.com.cn/problem/CF526D
// 给定一个长度为 n 的字符串和一个正整数k，判断其每个前缀是否形如 ABABA...BA
// A、B 可以为空，也可以是一个字符串，A 有 k+1 个，B 有 k 个.
//
// 令C=A+B,那么字符串可以变成 CC……CA，其中 C 有 k 个，且 A 是 C 的前缀（由定义而知）。
// 枚举C...C的长度，如果满足条件，则当前字符串一定具有周期 i-next[i-1](可能不完整).
func CF526D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)

	Z := ZAlgo(s)
	Z = append(Z, 0) // 哨兵
	next := GetNext(n, func(i int) int { return int(s[i]) })
	D := NewDiffArray(n + 1)

	for i := k; i <= n; i += k { // C...C可能的长度
		maybePeriod := i - next[i-1] // border性质
		if (i/k)%maybePeriod == 0 {
			start, end := i, i+min(i/k, Z[i])+1
			D.Add(start, end, 1)
		}
	}

	res := strings.Builder{}
	for i := 1; i <= n; i++ {
		if D.Get(i) > 0 {
			res.WriteByte('1')
		} else {
			res.WriteByte('0')
		}
	}
	fmt.Fprintln(out, res.String())
}

// 对每个前缀s[:i+1]，求具有循环节的前缀的长度和对应的循环次数
// 要求循环次数>=2.
// https://www.acwing.com/problem/content/discussion/index/143/1/
func acwing143() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(s string) [][2]int {
		n := len(s)
		res := make([][2]int, n)
		next := GetNext(n, func(i int) int { return int(s[i]) })
		for i := 0; i < n; i++ {
			period := Period(next, i)
			if period > 0 {
				res[i][0] = i + 1
				res[i][1] = (i + 1) / period
			}
		}
		return res
	}
	ptr := 1
	for {
		var n int
		fmt.Fscan(in, &n)
		if n == 0 {
			break
		}
		var s string
		fmt.Fscan(in, &s)
		res := solve(s)
		fmt.Fprintf(out, "Test case #%d\n", ptr)
		ptr++
		for _, v := range res {
			if v[0] == 0 {
				continue
			}
			fmt.Fprintf(out, "%d %d\n", v[0], v[1])
		}
		fmt.Fprintln(out)
	}
}

// `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度.
func GetNext(n int, f func(i int) int) []int {
	next := make([]int, n)
	j := 0
	for i := 1; i < n; i++ {
		for j > 0 && f(i) != f(j) {
			j = next[j-1]
		}
		if f(i) == f(j) {
			j++
		}
		next[i] = j
	}
	return next
}

// 求s的前缀[0:i+1)的最小周期.如果不存在,则返回0.
//
// !要求循环节次数>1 且 循环节完整.
//
//	0<=i<len(s).
func Period(next []int, i int) int {
	res := i + 1 - next[i]
	// !循环节次数>1 且 循环节完整
	if i+1 > res && (i+1)%res == 0 {
		return res
	}
	return 0
}

// z算法求字符串每个后缀与原串的最长公共前缀长度
//
// z[0]=0
// z[i]是后缀s[i:]与s的最长公共前缀(LCP)的长度 (i>=1)
func ZAlgo(s string) []int {
	n := len(s)
	z := make([]int, n)
	left, right := 0, 0
	for i := 1; i < n; i++ {
		z[i] = max(min(z[i-left], right-i+1), 0)
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			left, right = i, i+z[i]
			z[i]++
		}
	}
	return z
}

func ZAlgoNums(nums []int) []int {
	n := len(nums)
	z := make([]int, n)
	left, right := 0, 0
	for i := 1; i < n; i++ {
		z[i] = max(min(z[i-left], right-i+1), 0)
		for i+z[i] < n && nums[z[i]] == nums[i+z[i]] {
			left, right = i, i+z[i]
			z[i]++
		}
	}
	return z
}

type S = string

func GetHash(s S, base uint) uint {
	if len(s) == 0 {
		return 0
	}
	res := uint(0)
	for i := 0; i < len(s); i++ {
		res = res*base + uint(s[i])
	}
	return res
}

type RollingHash struct {
	base  uint
	mod   uint
	power []uint
}

// eg:
// NewRollingHash(37, 2102001800968)
// NewRollingHash(131, 999999751)
// mod: 999999937/999999929/999999893/999999797/999999761/999999757/999999751/999999739
func NewRollingHash(base uint, mod uint) *RollingHash {
	return &RollingHash{
		base:  base,
		mod:   mod,
		power: []uint{1},
	}
}

func (r *RollingHash) Build(s S) (hashTable []uint) {
	sz := len(s)
	hashTable = make([]uint, sz+1)
	for i := 0; i < sz; i++ {
		hashTable[i+1] = (hashTable[i]*r.base + uint(s[i])) % r.mod
	}
	return hashTable
}

func (r *RollingHash) Query(sTable []uint, start, end int) uint {
	r.expand(end - start)
	return (r.mod + sTable[end] - sTable[start]*r.power[end-start]%r.mod) % r.mod
}

func (r *RollingHash) Combine(h1, h2 uint, h2len int) uint {
	r.expand(h2len)
	return (h1*r.power[h2len] + h2) % r.mod
}

func (r *RollingHash) AddChar(hash uint, c byte) uint {
	return (hash*r.base + uint(c)) % r.mod
}

// 两个字符串的最长公共前缀长度.
func (r *RollingHash) LCP(sTable []uint, start1, end1 int, tTable []uint, start2, end2 int) int {
	len1 := end1 - start1
	len2 := end2 - start2
	len := min(len1, len2)
	low := 0
	high := len + 1
	for high-low > 1 {
		mid := (low + high) / 2
		if r.Query(sTable, start1, start1+mid) == r.Query(tTable, start2, start2+mid) {
			low = mid
		} else {
			high = mid
		}
	}
	return low
}

func (r *RollingHash) expand(sz int) {
	if len(r.power) < sz+1 {
		preSz := len(r.power)
		r.power = append(r.power, make([]uint, sz+1-preSz)...)
		for i := preSz - 1; i < sz; i++ {
			r.power[i+1] = (r.power[i] * r.base) % r.mod
		}
	}
}

type DiffArray struct {
	diff  []int
	dirty bool
}

func NewDiffArray(n int) *DiffArray {
	return &DiffArray{
		diff: make([]int, n+1),
	}
}

func (d *DiffArray) Add(start, end, delta int) {
	if start < 0 {
		start = 0
	}
	if end >= len(d.diff) {
		end = len(d.diff) - 1
	}
	if start >= end {
		return
	}
	d.dirty = true
	d.diff[start] += delta
	d.diff[end] -= delta
}

func (d *DiffArray) Build() {
	if d.dirty {
		preSum := make([]int, len(d.diff))
		for i := 1; i < len(d.diff); i++ {
			preSum[i] = preSum[i-1] + d.diff[i]
		}
		d.diff = preSum
		d.dirty = false
	}
}

func (d *DiffArray) Get(pos int) int {
	d.Build()
	return d.diff[pos]
}

func (d *DiffArray) GetAll() []int {
	d.Build()
	return d.diff[:len(d.diff)-1]
}

type LCAFast struct {
	Depth, Parent      []int32
	Tree               [][]int
	dfn, top, heavySon []int32
	idToNode           []int32
	dfnId              int32
}

func NewLCA(tree [][]int, roots []int) *LCAFast {
	n := len(tree)
	dfn := make([]int32, n)      // vertex => dfn
	top := make([]int32, n)      // 所处轻/重链的顶点（深度最小），轻链的顶点为自身
	depth := make([]int32, n)    // 深度
	parent := make([]int32, n)   // 父结点
	heavySon := make([]int32, n) // 重儿子
	idToNode := make([]int32, n)
	res := &LCAFast{
		Tree:     tree,
		dfn:      dfn,
		top:      top,
		Depth:    depth,
		Parent:   parent,
		heavySon: heavySon,
		idToNode: idToNode,
	}
	for _, root := range roots {
		root32 := int32(root)
		res._build(root32, -1, 0)
		res._markTop(root32, root32)
	}
	return res
}

// 树上多个点的 LCA，就是 DFS 序最小的和 DFS 序最大的这两个点的 LCA.
func (hld *LCAFast) LCAMultiPoint(nodes []int) int {
	if len(nodes) == 1 {
		return nodes[0]
	}
	if len(nodes) == 2 {
		return hld.LCA(nodes[0], nodes[1])
	}
	minDfn, maxDfn := int32(1<<31-1), int32(-1)
	for _, root := range nodes {
		root32 := int32(root)
		if hld.dfn[root32] < minDfn {
			minDfn = hld.dfn[root32]
		}
		if hld.dfn[root32] > maxDfn {
			maxDfn = hld.dfn[root32]
		}
	}
	u, v := hld.idToNode[minDfn], hld.idToNode[maxDfn]
	return hld.LCA(int(u), int(v))
}

func (hld *LCAFast) LCA(u, v int) int {
	u32, v32 := int32(u), int32(v)
	for {
		if hld.dfn[u32] > hld.dfn[v32] {
			u32, v32 = v32, u32
		}
		if hld.top[u32] == hld.top[v32] {
			return int(u32)
		}
		v32 = hld.Parent[hld.top[v32]]
	}
}

func (hld *LCAFast) Dist(u, v int) int {
	return int(hld.Depth[u] + hld.Depth[v] - 2*hld.Depth[hld.LCA(u, v)])
}

func (hld *LCAFast) EnumeratePathDecomposition(u, v int, vertex bool, f func(start, end int)) {
	u32, v32 := int32(u), int32(v)
	for {
		if hld.top[u32] == hld.top[v32] {
			break
		}
		if hld.dfn[u32] < hld.dfn[v32] {
			a, b := hld.dfn[hld.top[v32]], hld.dfn[v32]
			if a > b {
				a, b = b, a
			}
			f(int(a), int(b+1))
			v32 = hld.Parent[hld.top[v32]]
		} else {
			a, b := hld.dfn[u32], hld.dfn[hld.top[u32]]
			if a > b {
				a, b = b, a
			}
			f(int(a), int(b+1))
			u32 = hld.Parent[hld.top[u32]]
		}
	}

	edgeInt := int32(1)
	if vertex {
		edgeInt = 0
	}

	if hld.dfn[u32] < hld.dfn[v32] {
		a, b := hld.dfn[u32]+edgeInt, hld.dfn[v32]
		if a > b {
			a, b = b, a
		}
		f(int(a), int(b+1))
	} else if hld.dfn[v32]+edgeInt <= hld.dfn[u32] {
		a, b := hld.dfn[u32], hld.dfn[v32]+edgeInt
		if a > b {
			a, b = b, a
		}
		f(int(a), int(b+1))
	}
}

func (hld *LCAFast) _build(cur, pre, dep int32) int {
	subSize, heavySize, heavySon := 1, 0, int32(-1)
	for _, next := range hld.Tree[cur] {
		next32 := int32(next)
		if next32 != pre {
			nextSize := hld._build(next32, cur, dep+1)
			subSize += nextSize
			if nextSize > heavySize {
				heavySize, heavySon = nextSize, next32
			}
		}
	}
	hld.Depth[cur] = dep
	hld.heavySon[cur] = heavySon
	hld.Parent[cur] = pre
	return subSize
}

func (hld *LCAFast) _markTop(cur, top int32) {
	hld.top[cur] = top
	hld.dfn[cur] = hld.dfnId
	hld.idToNode[hld.dfnId] = cur
	hld.dfnId++
	if hld.heavySon[cur] != -1 {
		hld._markTop(hld.heavySon[cur], top)
		for _, next := range hld.Tree[cur] {
			next32 := int32(next)
			if next32 != hld.heavySon[cur] && next32 != hld.Parent[cur] {
				hld._markTop(next32, next32)
			}
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
