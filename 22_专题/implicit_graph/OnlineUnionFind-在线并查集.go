// 不预先给出图,而是通过函数来寻找下一个可以合并(到达)的点.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	CF1242B()
}

// 0-1 MST
// https://www.luogu.com.cn/problem/CF1242B
// 给定一张完全图, 其中有 m 条边权值为 1,其余边权为0， 求这张图的最小生成树
// n,m<=1e5
// !求出所有可以用0边连成的联通块，那么将这些联通块连起来，就可以构造出一颗最小生成树
// 答案为0组成的联通分量数-1.
func CF1242B() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	ban := make([]map[int]struct{}, n)
	for i := 0; i < n; i++ {
		ban[i] = make(map[int]struct{})
	}
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		ban[u][v] = struct{}{}
		ban[v][u] = struct{}{}
	}

	unUsed := NewFastSetFrom(n, func(i int) bool { return true })
	nexts := make([]int, n)

	setUsed := func(u int) { unUsed.Erase(u) }
	findUnused := func(u int) int {
		tmp := &nexts[u]
		for {
			*tmp = unUsed.Next(*tmp)
			if *tmp == n {
				return -1
			}
			if _, has := ban[u][*tmp]; !has {
				return *tmp
			}
			*tmp++
		}
	}

	uf := OnlineUnionFind(n, setUsed, findUnused)
	fmt.Println(uf.Part - 1)
}

// 在线并查集, 通过函数来寻找下一个可以合并(到达)的点.
//
//	setUsed(u) : 删除数据结构中的u, 用于标记u已经被使用了.
//	findUnused(u) : 返回u可以到达的下一个点, 如果没有返回-1.
func OnlineUnionFind(
	n int, setUsed func(u int), findUnused func(u int) (v int),
) (res *_UnionFindArray) {
	res = NewUnionFindArray(n)
	visited := make([]bool, n)
	queue := []int{}
	for v := 0; v < n; v++ {
		if visited[v] {
			continue
		}
		queue = append(queue, v)
		visited[v] = true
		setUsed(v)
		for len(queue) > 0 {
			x := queue[0]
			queue = queue[1:]
			setUsed(x)
			visited[x] = true
			for {
				to := findUnused(x)
				if to == -1 {
					break
				}
				res.Union(v, to)
				queue = append(queue, to)
				visited[to] = true
				setUsed(to)
			}
		}
	}
	return
}

func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type _UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
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

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *_UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}

type FastSet struct {
	n, lg int
	seg   [][]int
	size  int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.size
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}
