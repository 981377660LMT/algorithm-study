// 并行合并的并查集(ParallelRangeUnionFindOffline).
// 倍增结构与线段树的区别
// 对于任意两段长度相等的区间，倍增的子树结构是相同的，而线段树的子树结构是不同的.
// !倍增的子区间经过平移之后，仍然对应一个合法的子区间，而线段树的子区间平移之后不一定.
// 基于倍增结构的这一特点，可以实现"区间并行操作".

// 倍增技术的强大是基于一个很简单的倍增结构。
// jump(u,i)表示u沿着出边移动2^i步所在的位置，
// link(u,i)表示从u出发，通过少于i次的转移能抵达的所有结点的集合。
// !将jump(u,i)视为一个结点，它覆盖了所有link(u,2^i)上的结点。

// 实现思路：
//  0. 并查集的合并是满足结合律的，且幂等;
//  1. 区间link(start1,len)和区间link(start2,len)的合并，可以分解为两次合并，即合并
//     link(start1,1<<k)和link(start2,1<<k)的元素，以及合并
//     link(start1+len-(1<<k),1<<k)和link(start2+len-(1<<k),1<<k)的元素;
//     !这两个区间可以用两对jump点代替，即jump(start1,k)、jump(start2,k)以及jump(start1+len-(1<<k),k)、jump(start2+len-(1<<k),k);
//     !进而将一对区间的合并转化为两对jump点的合并;
//  2. 向下传递：
//     为了让区间合并关系下推到最底层的并查集，可以从最高层开始，逐层向下合并;
//     对于两个父结点jump(start1,k)与jump(start2,k)的合并，
//     可以变成它们拆分后得到的四个子区间祖先相同。那么合并四个子区间就可以了.
//     父结点代表的区间长度是2^k，所以子区间长度就应该是2^(k-1);
//     对子区间每个0<=i<n，
//     子区间需要合并jump(i,k-1)与jump(root,k-1)以及jump(i+2^(k-1),k-1)与jump(root+2^(k-1),k-1)的元素;

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type IUnionFind interface {
	Union(a, b int32) bool
	Find(a int32) int32
}

type RangeUnionFindTreeOffline[U IUnionFind] struct {
	n    int32
	data []U
}

func NewRangeUnionFindTreeOffline[U IUnionFind](n int32, createUnionFind func(n int32) U) *RangeUnionFindTreeOffline[U] {
	log := bits.Len32(uint32(n)) - 1
	data := make([]U, log+1)
	for i := range data {
		data[i] = createUnionFind(n)
	}
	return &RangeUnionFindTreeOffline[U]{n: n, data: data}
}

// 对0<=k<len，合并start1+k和start2+k的元素.
// O(1).
func (u *RangeUnionFindTreeOffline[U]) UnionRange(start1, start2 int32, len int32) {
	if len <= 0 {
		return
	}
	// !由于并查集合并满足幂等性，所以可以直接合并两个子jump.
	k := bits.Len32(uint32(len)) - 1                      // log2
	u.data[k].Union(start1, start2)                       // union jump(start1,k) and jump(start2,k)
	u.data[k].Union(start1+len-(1<<k), start2+len-(1<<k)) // union jump(end1-(1<<k),k) and jump(end2-(1<<k),k)
}

// 下推合并关系.
// O(n*log(n)).
func (u *RangeUnionFindTreeOffline[U]) Build() U {
	n, log := u.n, int32(len(u.data))-1
	for k := log - 1; k >= 0; k-- {
		for i := int32(0); i < n-(1<<k); i++ {
			root := u.data[k+1].Find(i)
			u.data[k].Union(i, root)               // union jump(i,k) and jump(root,k)
			u.data[k].Union(i+(1<<k), root+(1<<k)) // union jump(i+2^k,k) and jump(root+2^k,k)
		}
	}
	return u.data[0]
}

func main() {
	P3295()
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
	rangeUf := NewRangeUnionFindTreeOffline[*UnionFindArraySimple32](
		n, func(n int32) *UnionFindArraySimple32 { return NewUnionFindArraySimple32(n) },
	)
	for i := int32(0); i < m; i++ {
		var start1, end1, start2, end2 int32
		fmt.Fscan(in, &start1, &end1, &start2, &end2)
		start1, start2 = start1-1, start2-1
		rangeUf.UnionRange(start1, start2, end1-start1)
	}

	uf := rangeUf.Build()

	part := int(uf.Part)

	fmt.Fprintln(out, 9*qpow(10, part-1)%MOD)

}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = u.Find(u.data[key])
	return u.data[key]
}

func (u *UnionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}
