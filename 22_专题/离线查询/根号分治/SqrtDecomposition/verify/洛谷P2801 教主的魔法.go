// SqrtDecomposition
// https://nyaannyaan.github.io/library/data-structure/square-root-decomposition.hpp

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	// https://www.luogu.com.cn/problem/P2801
	// n<=1e6 q<=3000 k<=1e9
	// !区间更新:加法 区间查询:大于等于k的元素个数
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	nums := make([]E, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	sqrt := NewSqrtDecomposition(nums, int(math.Sqrt(float64(n))))
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "M" {
			var left, right, add int
			fmt.Fscan(in, &left, &right, &add)
			left--
			right--
			sqrt.Update(left, right, add)
		} else {
			var left, right int
			fmt.Fscan(in, &left, &right, &k)
			left--
			right--
			res := 0
			sqrt.Query(left, right, func(cur E) { res += cur })
			fmt.Fprintln(out, res)
		}
	}
}

var k int

type E = int
type Id = int

type Block struct {
	// dont modify
	id, left, right int
	nums            []E // block内的原序列

	sorted  []E // 排序后的序列
	lazyAdd Id
}

// 初始化块内数据(只会调用一次)
func (b *Block) Init() {
	b.Build()
}

// 重构
func (b *Block) Build() {
	b.sorted = append(b.sorted[:0:0], b.nums...)
	sort.Ints(b.sorted)
}
func (b *Block) UpdateAll(lazy Id) { b.lazyAdd += lazy }
func (b *Block) UpdatePart(start, end int, lazy Id) {
	for i := start; i <= end; i++ {
		b.nums[i] += lazy
	}
	b.Build() // !注意重构
}
func (b *Block) QueryAll() E {
	lower := sort.SearchInts(b.sorted, k-b.lazyAdd)
	return len(b.sorted) - lower
}
func (b *Block) QueryPart(start, end int) E {
	res := 0
	for i := start; i <= end; i++ {
		if b.nums[i]+b.lazyAdd >= k {
			res++
		}
	}
	return res
}

//
//
//
// dont modify the template below
//
//
//

type SqrtDecomposition struct {
	n   int
	bs  int
	bls []Block
}

// 指定维护的序列和分块大小初始化.
//  blockSize:分块大小,一般取根号n(300)
func NewSqrtDecomposition(nums []E, blockSize int) *SqrtDecomposition {
	n := len(nums)
	res := &SqrtDecomposition{n: len(nums), bs: blockSize, bls: make([]Block, (n-1)/blockSize+1)}
	for i, v := range nums {
		pos := i / blockSize
		if i%blockSize == 0 {
			res.bls[pos] = Block{left: i, id: pos, nums: make([]E, 0, blockSize)}
		}
		res.bls[pos].nums = append(res.bls[pos].nums, v)
	}
	for i := range res.bls {
		block := &res.bls[i]
		block.right = block.left + len(block.nums) - 1
		block.Init()
	}
	return res
}

// 更新闭区间[left,right]的值.
//  0<=left<=right<n
func (s *SqrtDecomposition) Update(left, right int, lazy Id) {
	for i := range s.bls {
		block := &s.bls[i]
		if block.right < left {
			continue
		}
		if block.left > right {
			break
		}

		if left <= block.left && block.right <= right {
			// !区间更新完整的块:类似线段树，只需要打上懒标记
			block.UpdateAll(lazy)
		} else {
			bl := max(block.left, left)
			br := min(block.right, right)
			// !区间修改不完整的块：暴力更新实际值
			block.UpdatePart(bl-block.left, br-block.left, lazy)
		}
	}
}

// 查询闭区间[left,right]的值.
//  0<=left<=right<n
func (s *SqrtDecomposition) Query(left, right int, forEach func(blockRes E)) {
	for i := range s.bls {
		block := &s.bls[i]
		if block.right < left {
			continue
		}
		if block.left > right {
			break
		}

		if left <= block.left && block.right <= right {
			// !区间查询完整的块:实际值+懒标记里的值
			forEach(block.QueryAll())
		} else {
			bl := max(block.left, left)
			br := min(block.right, right)
			// !区间查询不完整的块：暴力计算 实际值+懒标记里的值
			forEach(block.QueryPart(bl-block.left, br-block.left))
		}
	}
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
