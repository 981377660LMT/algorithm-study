// 静态区间逆序对
// https://judge.yosupo.jp/submission/119304
// n*sqrt(n)*log(n)
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	nums := make([]int, n)
	D := NewDiscretizer() // 离散化
	for i := 0; i < n; i++ {
		var num int
		fmt.Fscan(in, &num)
		nums[i] = num
		D.Add(num)
	}
	D.Build()

	inv := 0
	bit := NewBITArray(D.Len())
	mo := NewMoAlgo(n, q, op{
		add: func(i, delta int) {
			num := D.Get(nums[i])
			if delta == 1 { // append
				inv += bit.QueryRange(num+1, bit.Len())
				bit.Add(num, 1)
			} else { // appendleft
				inv += bit.Query(num - 1)
				bit.Add(num, 1)
			}
		},
		remove: func(i, delta int) {
			num := D.Get(nums[i])
			if delta == 1 { // popleft
				inv -= bit.Query(num - 1)
				bit.Add(num, -1)
			} else { // pop
				inv -= bit.QueryRange(num+1, bit.Len())
				bit.Add(num, -1)
			}
		},
		query: func(qLeft, qRight int) int {
			return inv
		},
	},
	)

	for i := 0; i < q; i++ {
		var left, right int
		fmt.Fscan(in, &left, &right)
		right--
		mo.AddQuery(left, right)
	}

	res := mo.Work()
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type V = int
type R = int

type MoAlgo struct {
	queryOrder int
	chunkSize  int
	buckets    [][]query
	op         op
}

type query struct{ qi, left, right int }

type op struct {
	// 将数据添加到窗口
	add func(index, delta int)
	// 将数据从窗口中移除
	remove func(index, delta int)
	// 更新当前窗口的查询结果
	query func(qLeft, qRight int) R
}

func NewMoAlgo(n, q int, op op) *MoAlgo {
	chunkSize := max(1, n/int(math.Sqrt(float64(q))))
	buckets := make([][]query, n/chunkSize+1)
	return &MoAlgo{chunkSize: chunkSize, buckets: buckets, op: op}
}

// 0 <= left <= right < n
func (mo *MoAlgo) AddQuery(left, right int) {
	index := left / mo.chunkSize
	mo.buckets[index] = append(mo.buckets[index], query{mo.queryOrder, left, right + 1})
	mo.queryOrder++
}

// 返回每个查询的结果
func (mo *MoAlgo) Work() []R {
	buckets := mo.buckets
	res := make([]R, mo.queryOrder)
	left, right := 0, 0
	for i, bucket := range buckets {
		if i&1 == 1 {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right > bucket[j].right })
		} else {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right < bucket[j].right })
		}

		for _, q := range bucket {
			// !窗口扩张
			for left > q.left {
				left--
				mo.op.add(left, -1)
			}
			for right < q.right {
				mo.op.add(right, 1)
				right++
			}

			// !窗口收缩
			for left < q.left {
				mo.op.remove(left, 1)
				left++
			}
			for right > q.right {
				right--
				mo.op.remove(right, -1)
			}

			res[q.qi] = mo.op.query(q.left, q.right-1)
		}
	}

	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// !单点修改,区间查询,数组实现的树状数组
type BITArray struct {
	n    int
	tree []int
}

func NewBITArray(n int) *BITArray {
	return &BITArray{n: n, tree: make([]int, n+1)}
}

// 常数优化: dp O(n) 建树
// https://oi-wiki.org/ds/fenwick/#tricks
func (b *BITArray) Build(nums []int) {
	for i := 1; i < len(b.tree); i++ {
		b.tree[i] += nums[i-1]
		if j := i + (i & -i); j < len(b.tree) {
			b.tree[j] += b.tree[i]
		}
	}
}

// 位置 index 增加 delta
//  1<=i<=n
func (b *BITArray) Add(index int, delta int) {

	for ; index < len(b.tree); index += index & -index {
		b.tree[index] += delta
	}
}

// 求前缀和
//  1<=i<=n
func (b *BITArray) Query(index int) (res int) {
	if index > b.n {
		index = b.n
	}
	for ; index > 0; index -= index & -index {
		res += b.tree[index]
	}
	return
}

// 1<=left<=right<=n
func (b *BITArray) QueryRange(left, right int) int {
	return b.Query(right) - b.Query(left-1)
}

func (b *BITArray) Len() int {
	return b.n
}

// 离散化
type Discretizer struct {
	allNums []int
	set     map[int]struct{}
	mp      map[int]int // num -> index (1-based)
}

func NewDiscretizer() *Discretizer {
	return &Discretizer{
		set: make(map[int]struct{}, 16),
		mp:  make(map[int]int, 16),
	}
}

func (d *Discretizer) Add(num int) {
	d.set[num] = struct{}{}
}

func (d *Discretizer) Build() {
	allNums := make([]int, 0, len(d.set))
	for num := range d.set {
		allNums = append(allNums, num)
	}
	sort.Ints(allNums)
	d.allNums = allNums
	for i, num := range allNums {
		d.mp[num] = i + 1
	}
}

func (d *Discretizer) Get(num int) int {
	res, ok := d.mp[num]
	if ok {
		return res
	}
	// bisect right
	return sort.Search(len(d.allNums), func(i int) bool { return d.allNums[i] > num })
}

func (d *Discretizer) Len() int {
	return len(d.allNums)
}
