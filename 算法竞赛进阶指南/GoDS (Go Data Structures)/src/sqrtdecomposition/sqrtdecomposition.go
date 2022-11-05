// !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/sqrt_decomposition.go#L1-L110

package copypasta

import (
	"math"
	"sort"
)

/* 根号分治 Sqrt Decomposition
一种技巧：组合两种算法从而降低复杂度 O(n^2) -> O(n√n)

有 n 个对象，每个对象有一个「关于其他对象的统计量」ci（一个数、一个集合的元素个数，等等）
为方便起见，假设 ∑ci 的数量级和 n 一样，下面用 n 表示 ∑ci
当 ci > √n 时，这样的对象不超过 √n 个，暴力枚举这些对象之间的关系（或者，该对象与其他所有对象的关系），时间复杂度为 O(n) 或 O(n√n)。此乃算法一
当 ci ≤ √n 时，这样的对象有 O(n) 个，由于统计量 ci 很小，暴力枚举当前对象的统计量，时间复杂度为 O(n√n)。此乃算法二
这样，以 √n 为界，我们将所有对象划分成了两组，并用两个不同的算法处理
这两种算法是看待同一个问题的两种不同方式，通过恰当地组合（平衡）这两个算法，复杂度由 O(n^2) 降至 O(n√n)
注意：**枚举时要做到不重不漏**
*/

/*
分块数据结构
https://oi-wiki.org/ds/decompose/
https://oi-wiki.org/ds/block-array/
【推荐】https://www.luogu.com.cn/blog/220037/Sqrt1

教主的魔法 https://www.luogu.com.cn/problem/solution/P2801 => 区间修改、区间查询有多少个数大于k (线段树做不到)
对于整块，我们可以打上一个add标记，这样二分查找就要查 >= k-add 的值。
对于不完整的块，我们暴力修改，再直接排序整个块。
*/
func SqrtDecompotision() {
	type block struct {
		left, right    int
		origin, sorted []int
		lazyAdd        int
	}

	var blocks []block
	sqrtInit := func(nums []int) {
		n := len(nums)
		blockSize := int(math.Sqrt(float64(n))) //blockSize := int(math.Sqrt(float64(n) * math.Log2(float64(n+1))))
		blockNum := (n-1)/blockSize + 1
		blocks = make([]block, blockNum)
		for i, v := range nums {
			pos := i / blockSize
			if i%blockSize == 0 {
				blocks[pos] = block{left: i, origin: make([]int, 0, blockSize)}
			}
			blocks[pos].origin = append(blocks[pos].origin, v)
		}

		for i := range blocks {
			block := &blocks[i]
			block.right = block.left + len(block.origin) - 1
			// 对每个块的元素进行初始化
			block.sorted = append([]int(nil), block.origin...)
			sort.Ints(block.sorted)
		}
	}

	// 区间更新或者区间查询
	sqrtOp := func(left, right int, value int) { // [l,r], starts at 0
		for i := range blocks {
			block := &blocks[i]
			if block.right < left {
				continue
			}

			if block.left > right {
				break
			}

			if left <= block.left && block.right <= right {
				// do op on full block
				// 区间更新，类似线段树的懒标记
			} else {
				// do op on part block
				// 暴力更新
				bl := max(block.left, left)
				br := min(block.right, right)
				for j := bl - block.left; j <= br-block.left; j++ {
					// do b.origin[j]...
				}
			}
		}
	}

	_ = []interface{}{sqrtInit, sqrtOp}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
