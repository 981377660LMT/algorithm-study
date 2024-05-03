package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// P4151()
	yosupo()
	// demo()

}

// https://judge.yosupo.jp/problem/intersection_of_f2_vector_spaces
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(arr1, arr2 []uint) []uint {
		return F2Intersection(arr1, arr2)
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var n int32
		fmt.Fscan(in, &n)
		arr1 := make([]uint, n)
		for i := int32(0); i < n; i++ {
			fmt.Fscan(in, &arr1[i])
		}
		var m int32
		fmt.Fscan(in, &m)
		arr2 := make([]uint, m)
		for i := int32(0); i < m; i++ {
			fmt.Fscan(in, &arr2[i])
		}
		res := solve(arr1, arr2)
		fmt.Fprint(out, len(res), " ")
		for _, base := range res {
			fmt.Fprint(out, base, " ")
		}
		fmt.Fprintln(out)
	}
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
			vs.Add(uint(cycleXor))
		}
	}

	dist := uf.Dist(start, end)
	fmt.Fprintln(out, vs.Max(uint(dist)))
}

// VectorSpace，线性基空间.支持线性基合并.
type VectorSpace struct {
	Bases []uint
}

func NewVectorSpace(nums []uint) *VectorSpace {
	res := &VectorSpace{}
	for _, num := range nums {
		res.Add(num)
	}
	return res
}

// 插入一个向量,如果插入成功(不能被表出)返回True,否则返回False.
func (lb *VectorSpace) Add(num uint) bool {
	for _, base := range lb.Bases {
		if base == 0 || num == 0 {
			break
		}
		num = minu(num, num^base)
	}
	if num != 0 {
		lb.Bases = append(lb.Bases, num)
		return true
	}
	return false
}

// 插入一个向量,如果插入成功(不能被表出)返回新的线性基,否则返回0.
func (lb *VectorSpace) Add2(num uint) uint {
	for _, base := range lb.Bases {
		if base == 0 || num == 0 {
			break
		}
		num = minu(num, num^base)
	}
	if num != 0 {
		lb.Bases = append(lb.Bases, num)
		return num
	}
	return 0
}

// 求xor与所有向量异或的最大值.
func (lb *VectorSpace) Max(xor uint) uint {
	res := xor
	for _, base := range lb.Bases {
		res = maxu(res, res^base)
	}
	return res
}

// 求xor与所有向量异或的最小值.
func (lb *VectorSpace) Min(xorVal uint) uint {
	res := xorVal
	for _, base := range lb.Bases {
		res = minu(res, res^base)
	}
	return res
}

func (lb *VectorSpace) Copy() *VectorSpace {
	return &VectorSpace{
		Bases: append(lb.Bases[:0:0], lb.Bases...),
	}
}

func (lb *VectorSpace) Clear() {
	lb.Bases = lb.Bases[:0]
}

func (lb *VectorSpace) Len() int32 {
	return int32(len(lb.Bases))
}

func (lb *VectorSpace) ForEach(f func(base uint)) {
	for _, base := range lb.Bases {
		f(base)
	}
}

func (lb *VectorSpace) Has(v uint) bool {
	for _, w := range lb.Bases {
		if v == 0 {
			break
		}
		v = minu(v, v^w)
	}
	return v == 0
}

func (lb *VectorSpace) String() string {
	return fmt.Sprintf("%v", lb.Bases)
}

// Or.
// 线性基合并.
func Merge(v1, v2 *VectorSpace) *VectorSpace {
	if v1.Len() < v2.Len() {
		v1, v2 = v2, v1
	}
	res := v1.Copy()
	for _, base := range v2.Bases {
		res.Add(base)
	}
	return res
}

// Or.
// 线性基合并.
func MergeDestructively(v1, v2 *VectorSpace) *VectorSpace {
	if v1.Len() < v2.Len() {
		v1, v2 = v2, v1
	}
	for _, base := range v2.Bases {
		v1.Add(base)
	}
	return v1
}

// Intersection.(And)
// 注意此时最大值不超过2**30.
func F2Intersection(A, B []uint) []uint {
	tmp := &VectorSpace{}
	for _, a := range A {
		tmp.Add(a<<32 + a)
	}
	upper := uint(1 << 32)
	var res []uint
	for _, b := range B {
		v := b << 32
		u := tmp.Add2(v)
		if u < upper {
			res = append(res, u)
		}
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minu(a, b uint) uint {
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

func maxu(a, b uint) uint {
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
