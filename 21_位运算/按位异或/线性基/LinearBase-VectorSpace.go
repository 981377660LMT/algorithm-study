package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	P4151()
}

func demo() {
	v1, v2 := NewVectorSpace(nil), NewVectorSpace(nil)
	v1.Add(1)
	v1.Add(2)
	v1.Add(3)
	v2.Add(1)
	v2.Add(5)
	v2.Add(7)
	v2.Add(8)

	fmt.Println(v1.Or(v2))
	fmt.Println(v1)
	fmt.Println(v1.And(v2))
	fmt.Println(v1.And(NewVectorSpace(nil)))
	fmt.Println(v1)
}

// P4151 [WC2011] 最大XOR和路径
// https://www.luogu.com.cn/problem/P4151
// 考虑一个边权为非负整数的`无向连通图`，节点编号为 0 到 N-1.
// 试求出一条从 0 号节点到 N-1 号节点的路径，使得路径上经过的边的权值的 XOR 和最大。
// !路径可以重复经过某些点或边
//
// !将所有环的异或扔进线性基，答案就是0到n-1的路径的权值与线性基的最大异或和.
func P4151() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][3]int, m)
	for i := range edges {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}
	start, end := 0, n-1

	uf := NewUnionFindArrayWithDist(n)
	vs := NewVectorSpace(nil)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		root1, root2 := uf.Find(u), uf.Find(v)
		if root1 != root2 {
			uf.Union(u, v, w)
		} else {
			cycleXor := uf.Dist(u, v) ^ w
			vs.Add(cycleXor)
		}
	}

	dist := uf.Dist(start, end)
	fmt.Fprintln(out, vs.Max(dist))
}

// VectorSpace，线性基空间.支持线性基合并.
type VectorSpace struct {
	bases  []int
	maxBit int
}

func NewVectorSpace(nums []int) *VectorSpace {
	res := &VectorSpace{}
	for _, num := range nums {
		res.Add(num)
	}
	return res
}

// 插入一个向量,如果插入成功(不能被表出)返回True,否则返回False.
func (lb *VectorSpace) Add(num int) bool {
	for _, base := range lb.bases {
		if base == 0 || num == 0 {
			break
		}
		num = min(num, num^base)
	}
	if num != 0 {
		lb.bases = append(lb.bases, num)
		lb.maxBit = max(lb.maxBit, num)
		return true
	}
	return false
}

// 求xor与所有向量异或的最大值.
func (lb *VectorSpace) Max(xor int) int {
	res := xor
	for _, base := range lb.bases {
		res = max(res, res^base)
	}
	return res
}

// 求xor与所有向量异或的最小值.
func (lb *VectorSpace) Min(xorVal int) int {
	res := xorVal
	for _, base := range lb.bases {
		res = min(res, res^base)
	}
	return res
}

func (lb *VectorSpace) Copy() *VectorSpace {
	res := &VectorSpace{}
	res.bases = append(res.bases, lb.bases...)
	res.maxBit = lb.maxBit
	return res
}

func (lb *VectorSpace) Len() int {
	return len(lb.bases)
}

func (lb *VectorSpace) ForEach(f func(base int)) {
	for _, base := range lb.bases {
		f(base)
	}
}

func (lb *VectorSpace) Has(v int) bool {
	for _, w := range lb.bases {
		if v == 0 {
			break
		}
		v = min(v, v^w)
	}
	return v == 0
}

// Merge.
func (lb *VectorSpace) Or(other *VectorSpace) *VectorSpace {
	v1, v2 := lb, other
	if v1.Len() < v2.Len() {
		v1, v2 = v2, v1
	}
	res := v1.Copy()
	for _, base := range v2.bases {
		res.Add(base)
	}
	return res
}

func (lb *VectorSpace) And(other *VectorSpace) *VectorSpace {
	maxDim := max(lb.maxBit, other.maxBit)
	x := lb.orthogonalSpace(maxDim)
	y := other.orthogonalSpace(maxDim)
	if x.Len() < y.Len() {
		x, y = y, x
	}
	for _, base := range y.bases {
		x.Add(base)
	}
	return x.orthogonalSpace(maxDim)
}

func (lb *VectorSpace) String() string {
	return fmt.Sprintf("%v", lb.bases)
}

// 正交空间.
func (lb *VectorSpace) orthogonalSpace(maxDim int) *VectorSpace {
	lb.normalize(true)
	m := maxDim
	tmp := make([]int, m)
	for _, base := range lb.bases {
		tmp[bits.Len(uint(base))-1] = base
	}
	tmp = Transpose(m, m, tmp, true)
	res := &VectorSpace{}
	for j, v := range tmp {
		if v>>j&1 == 1 {
			continue
		}
		res.Add(v | 1<<j)
	}
	return res
}

func (lb *VectorSpace) normalize(reverse bool) {
	for j, v := range lb.bases {
		for i := 0; i < j; i++ {
			lb.bases[i] = min(lb.bases[i], lb.bases[i]^v)
		}
	}
	if !reverse {
		sort.Ints(lb.bases)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(lb.bases)))
	}
}

// 矩阵转置,O(n+m)log(n+m)
func Transpose(row, col int, grid []int, inPlace bool) []int {
	if len(grid) != row {
		panic("row not match")
	}
	if !inPlace {
		grid = append(grid[:0:0], grid...)
	}
	log := 0
	max_ := max(row, col)
	for 1<<log < max_ {
		log++
	}
	if len(grid) < 1<<log {
		*&grid = append(grid, make([]int, 1<<log-len(grid))...)
	}
	width := 1 << log
	mask := int(1)
	for i := 0; i < log; i++ {
		mask |= (mask << (1 << i))
	}
	for t := 0; t < log; t++ {
		width >>= 1
		mask ^= (mask >> width)
		for i := 0; i < 1<<t; i++ {
			for j := 0; j < width; j++ {
				x := &grid[width*(2*i)+j]
				y := &grid[width*(2*i+1)+j]
				*x = ((*y << width) & mask) ^ *x
				*y = ((*x & mask) >> width) ^ *y
				*x = ((*y << width) & mask) ^ *x
			}
		}
	}
	return grid[:col]
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

//
//

type T = int // 距离为异或

func e() T        { return 0 }
func op(x, y T) T { return x ^ y }
func inv(x T) T   { return x }

// 数组实现的带权并查集(维护到每个组根节点距离的并查集).
// 用于维护环的权值，树上的距离等.
type UnionFindArrayWithDist struct {
	Part      int
	data      []int
	potential []T
}

func NewUnionFindArrayWithDist(n int) *UnionFindArrayWithDist {
	uf := &UnionFindArrayWithDist{
		Part:      n,
		data:      make([]int, n),
		potential: make([]T, n),
	}
	for i := range uf.data {
		uf.data[i] = -1
		uf.potential[i] = e()
	}
	return uf
}

// p[x] = p[y] + dist.
//
//	如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
func (uf *UnionFindArrayWithDist) Union(x, y int, dist T) bool {
	dist = op(dist, op(uf.DistToRoot(y), inv(uf.DistToRoot(x))))
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return dist == e()
	}
	if uf.data[x] < uf.data[y] {
		x, y = y, x
		dist = inv(dist)
	}
	uf.data[y] += uf.data[x]
	uf.data[x] = y
	uf.potential[x] = dist
	uf.Part--
	return true
}

// p[x] = p[y] + dist.
//
//	如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
func (uf *UnionFindArrayWithDist) UnionWithCallback(x, y int, dist T, cb func(big, small int)) bool {
	dist = op(dist, op(uf.DistToRoot(y), inv(uf.DistToRoot(x))))
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return dist == e()
	}
	if uf.data[x] < uf.data[y] {
		x, y = y, x
		dist = inv(dist)
	}
	uf.data[y] += uf.data[x]
	uf.data[x] = y
	uf.potential[x] = dist
	uf.Part--
	if cb != nil {
		cb(y, x)
	}
	return true
}

func (uf *UnionFindArrayWithDist) Find(x int) int {
	if uf.data[x] < 0 {
		return x
	}
	root := uf.Find(uf.data[x])
	uf.potential[x] = op(uf.potential[x], uf.potential[uf.data[x]])
	uf.data[x] = root
	return root
}

// f[x]-f[find(x)].
//
//	点x到所在组根节点的距离.
func (uf *UnionFindArrayWithDist) DistToRoot(x int) T {
	uf.Find(x)
	return uf.potential[x]
}

// f[x] - f[y].
func (uf *UnionFindArrayWithDist) Dist(x, y int) T {
	return op(uf.DistToRoot(x), inv(uf.DistToRoot(y)))
}

func (uf *UnionFindArrayWithDist) GetSize(x int) int {
	return -uf.data[uf.Find(x)]
}

func (uf *UnionFindArrayWithDist) GetGroups() map[int][]int {
	res := make(map[int][]int)
	for i := range uf.data {
		root := uf.Find(i)
		res[root] = append(res[root], i)
	}
	return res
}

func (uf *UnionFindArrayWithDist) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}
