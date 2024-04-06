// DivideIntervalBinaryLift
// 倍增拆分序列上的区间 `[start,end)`
// 一共拆分成[0,log]层，每层有n个元素.
// !jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n).

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type DivideIntervalByBinaryLift struct {
	n, log int32
	size   int32
}

func NewDivideIntervalByBinaryLift(n int32) *DivideIntervalByBinaryLift {
	log := int32(bits.Len(uint(n))) - 1
	size := n * (log + 1)
	return &DivideIntervalByBinaryLift{n: n, log: log, size: size}
}

// 遍历[start,end)区间内的所有jump.
func (d *DivideIntervalByBinaryLift) EnumerateRange(start, end int32, f func(level, index int32)) {
	if start >= end {
		return
	}
	cur := start
	log := d.log
	len := end - start
	for k := log; k >= 0; k-- {
		if len&(1<<k) != 0 {
			f(k, cur)
			cur += 1 << k
			if cur >= end {
				return
			}
		}
	}
	f(0, cur)
}

// 遍历[start1,end1)区间和[start2,end2)区间内的所有jump.要求区间长度相等.
func (d *DivideIntervalByBinaryLift) EnumerateRange2(start1, end1 int32, start2, end2 int32, f func(level, index1, index2 int32)) {
	if end1-start1 != end2-start2 {
		panic("not same length")
	}
	if start1 >= end1 {
		return
	}
	cur1, cur2 := start1, start2
	log := d.log
	len := end1 - start1
	for k := log; k >= 0; k-- {
		if len&(1<<k) != 0 {
			f(k, cur1, cur2)
			cur1 += 1 << k
			cur2 += 1 << k
			if cur1 >= end1 {
				return
			}
		}
	}
	f(0, cur1, cur2)
}

// 从高的jump开始下推信息，更新底部jump的答案.
// O(n*log(n)).
func (d *DivideIntervalByBinaryLift) PushDown(
	f func(parentLevel, parentIndex int32, childLevel, childIndex1, childIndex2 int32),
) {
	n, log := d.n, d.log
	for k := log - 1; k >= 0; k-- {
		for i := int32(0); i < n-(1<<k); i++ {
			f(k+1, i, k, i, i+(1<<k))
		}
	}
}

func (d *DivideIntervalByBinaryLift) Size() int32 { return d.size }
func (d *DivideIntervalByBinaryLift) Log() int32  { return d.log }

func main() {
	P3295()
	// demo()
}

func demo() {
	n := int32(10)
	D := NewDivideIntervalByBinaryLift(n)
	values := make([]int32, D.Size())
	D.EnumerateRange(3, 5, func(level, index int32) { values[level*n+index] = 10 })
	D.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
		p, c1, c2 := pLevel*n+pIndex, cLevel*n+cIndex1, cLevel*n+cIndex2
		values[c1] = max(values[c1], values[p])
		values[c2] = max(values[c2], values[p])
	})
	fmt.Println(values[:n])
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
	D := NewDivideIntervalByBinaryLift(n)
	ufs := make([]*UnionFindArraySimple32, D.Log()+1)
	for i := range ufs {
		ufs[i] = NewUnionFindArraySimple32(n)
	}

	for i := int32(0); i < m; i++ {
		var start1, end1, start2, end2 int32
		fmt.Fscan(in, &start1, &end1, &start2, &end2)
		start1, start2 = start1-1, start2-1
		D.EnumerateRange2(start1, end1, start2, end2, func(level, i1, i2 int32) {
			ufs[level].Union(i1, i2)
		})
	}

	D.PushDown(func(pLevel, pIndex, cLevel, cIndex1, _ int32) {
		root := ufs[pLevel].Find(pIndex)
		ufs[cLevel].Union(cIndex1, root)
		ufs[cLevel].Union(cIndex1+1<<cLevel, root+1<<cLevel)
	})

	uf := ufs[0]
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
