// ParallelUnionFindOffline
// https://github.com/yosupo06/library-checker-problems/issues/934
// https://atcoder.jp/contests/yahoo-procon2018-final/submissions/8391439
// http://kmyk.github.io/competitive-programming-library/old/range-union-find-tree.inc.cpp
// https://yosupo.hatenablog.com/entry/2019/11/12/001535

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	// P3295()
	// atc2018()
	// abc349_g()
}

func demo() {
	uf := NewUnionFindRangeParallel(10)
	uf.UnionParallelly(0, 5, 5)
	uf.Build(func(big, small int32) {
		fmt.Printf("big=%d small=%d\n", big, small)
	})
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
	ufp := NewUnionFindRangeParallel(n)
	for i := int32(0); i < m; i++ {
		var start1, end1, start2, end2 int32
		fmt.Fscan(in, &start1, &end1, &start2, &end2)
		start1, start2 = start1-1, start2-1
		ufp.UnionParallelly(start1, start2, end1-start1)
	}

	uf := ufp.Build(func(big, small int32) {})

	part := int(uf.Part)

	fmt.Fprintln(out, 9*qpow(10, part-1)%MOD)

}

func atc2018() {
	// https://atcoder.jp/contests/yahoo-procon2018-final/tasks/yahoo_procon2018_final_d
	// !前缀和后缀的LCP为lens[i]的字符串
	// !LCP => 并查集
	// 给定长为n的数组lens, 问是否存在一个长度为s的字符串,满足:
	// !s[0:i+1] 和 s[n-(i+1):n] 的最长公共前缀为 lens[i] (0<=i<n)
	// n<=3e5 0<=lens[i]<=i+1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	lens := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &lens[i])
	}

	ufrp := NewUnionFindRangeParallel(n)
	for i := int32(0); i < n; i++ {
		ufrp.UnionParallelly(0, n-(i+1), lens[i]) // 各个位置的字符相同
	}

	uf := ufrp.Build(func(big, small int32) {})
	for i := int32(0); i < n; i++ {
		if lens[i] == i+1 {
			continue
		}
		if uf.IsConnected(lens[i], n-(i+1)+lens[i]) { // !s[len[i]]!=s[n-(i+1)+len[i]] (因为前后缀LCP只有len[i])
			fmt.Fprintln(out, "No")
			return
		}
	}
	fmt.Fprintln(out, "Yes")
}

// G - Palindrome Construction (并查集)
// https://atcoder.jp/contests/abc349/tasks/abc349_g
// 给定一个数组A[i], A[i]表示以位置i为中心的极长回文串的半径(奇数长度回文).
// 求出满足条件的字典序最小的正整数序列.如果不存在输出No.
// https://atcoder.jp/contests/abc349/submissions/52345482
//
// !类似后缀数组解决最长回文串问题，将字符串反串添加在后面.
// !回文关系转化为区间相等关系.
func abc349_g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	rangeUf := NewUnionFindRangeParallel(n * 2)
	for i := int32(0); i < n; i++ {
		l, r := i-nums[i], i+nums[i]
		rangeUf.UnionParallelly(l, n+(n-1-r), 2*nums[i]+1) // !必须相等
	}
	uf := rangeUf.Build(nil)

	groups := make([][]int32, n)
	for i := int32(0); i < n; i++ {
		root := uf.Find(i)
		groups[root] = append(groups[root], i)
	}
	notSame := make([][]int32, n)
	for i := int32(0); i < n; i++ {
		l, r := i-nums[i]-1, i+nums[i]+1 // !不能相等
		if 0 <= l && r < n {
			root1, root2 := uf.Find(l), uf.Find(r)
			if root1 == root2 {
				fmt.Println("No")
				return
			}
			notSame[root1] = append(notSame[root1], root2)
			notSame[root2] = append(notSame[root2], root1)
		}
	}

	res := make([]int32, n)
	for i := int32(0); i < n; i++ {
		if res[i] != 0 {
			continue
		}

		root := uf.Find(i)
		ng := make(map[int32]struct{})
		for _, v := range notSame[root] {
			if res[v] != 0 {
				ng[res[v]] = struct{}{}
			}
		}
		for v := int32(1); ; v++ {
			if _, ok := ng[v]; !ok {
				for _, u := range groups[root] {
					res[u] = v
				}
				break
			}
		}
	}

	fmt.Println("Yes")
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// 并行合并的并查集.
type UnionFindRangeParallel struct {
	n      int32
	queues [][][2]int32
}

func NewUnionFindRangeParallel(n int32) *UnionFindRangeParallel {
	return &UnionFindRangeParallel{n: n, queues: make([][][2]int32, n+1)}
}

// !并行合并[(a,b),(a+1,b+1),...,(a+len-1,b+len-1)].
func (uf *UnionFindRangeParallel) UnionParallelly(a, b, len int32) {
	if len == 0 {
		return
	}
	min_ := min32(len, uf.n)
	uf.queues[min_] = append(uf.queues[min_], [2]int32{a, b})
}

func (uf *UnionFindRangeParallel) Build(f func(big, small int32)) *_unionFindRange {
	res := _newUnionFindRange(uf.n)
	queue, nextQueue := [][2]int32{}, [][2]int32{}
	for di := uf.n; di >= 1; di-- {
		queue = append(queue, uf.queues[di]...)
		nextQueue = nextQueue[:0]
		for _, p := range queue {
			if res.IsConnected(p[0], p[1]) {
				continue
			}
			res.Union(p[0], p[1], f)
			nextQueue = append(nextQueue, [2]int32{p[0] + 1, p[1] + 1})
		}
		queue, nextQueue = nextQueue, queue
	}
	return res
}

type _unionFindRange struct {
	Part   int32
	n      int32
	parent []int32
	rank   []int32
}

func _newUnionFindRange(n int32) *_unionFindRange {
	uf := &_unionFindRange{
		Part:   n,
		n:      n,
		parent: make([]int32, n),
		rank:   make([]int32, n),
	}
	for i := int32(0); i < n; i++ {
		uf.parent[i] = i
		uf.rank[i] = 1
	}
	return uf
}

func (uf *_unionFindRange) Find(x int32) int32 {
	for x != uf.parent[x] {
		uf.parent[x] = uf.parent[uf.parent[x]]
		x = uf.parent[x]
	}
	return x
}

// Union 后, 大的编号的组会指向小的编号的组.
func (uf *_unionFindRange) Union(x, y int32, f func(big, small int32)) bool {
	if x < y {
		x, y = y, x
	}
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	if rootX == rootY {
		return false
	}
	uf.parent[rootX] = rootY
	uf.rank[rootY] += uf.rank[rootX]
	uf.Part--
	if f != nil {
		f(rootY, rootX)
	}
	return true
}

// UnionRange 合并区间 [left, right] 的所有元素, 返回合并次数.
func (uf *_unionFindRange) UnionRange(left, right int32, f func(big, small int32)) int {
	if left >= right {
		return 0
	}
	leftRoot := uf.Find(left)
	rightRoot := uf.Find(right)
	unionCount := 0
	for rightRoot != leftRoot {
		unionCount++
		uf.Union(rightRoot, rightRoot-1, f)
		rightRoot = uf.Find(rightRoot - 1)
	}
	return unionCount
}

func (uf *_unionFindRange) IsConnected(x, y int32) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *_unionFindRange) GetSize(x int32) int32 {
	return uf.rank[uf.Find(x)]
}

func (uf *_unionFindRange) GetGroups() map[int32][]int32 {
	group := make(map[int32][]int32)
	for i := int32(0); i < uf.n; i++ {
		group[uf.Find(i)] = append(group[uf.Find(i)], i)
	}
	return group
}

func min(a, b int) int {
	if a < b {
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
