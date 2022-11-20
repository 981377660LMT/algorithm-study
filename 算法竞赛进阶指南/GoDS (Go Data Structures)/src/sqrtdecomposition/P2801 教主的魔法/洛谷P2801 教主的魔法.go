package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// https://www.luogu.com.cn/problem/P2801
// n<=1e6 q<=3000 k<=1e9
// 区间更新:加法 区间查询:大于等于k的元素个数
// https://zhuanlan.zhihu.com/p/114268236
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	type block struct {
		left, right int
		elements    []int // !raw data

		sorted  []int // data 排序后的元素
		lazyAdd int   // lazy
	}

	// sqrt分块模板
	// !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/sqrt_decomposition.go#L1-L110
	var blocks []block
	sqrtInit := func(nums []int) {
		n := len(nums)
		blockSize := int(math.Sqrt(float64(n))) //blockSize := int(math.Sqrt(float64(n) * math.Log2(float64(n+1))))
		blockNum := (n-1)/blockSize + 1
		blocks = make([]block, blockNum)
		for i, v := range nums {
			pos := i / blockSize
			if i%blockSize == 0 {
				blocks[pos] = block{left: i, elements: make([]int, 0, blockSize)}
			}
			blocks[pos].elements = append(blocks[pos].elements, v)
			// !对每个块的元素进行初始化1
			blocks[pos].lazyAdd = 0
		}

		for i := range blocks {
			block := &blocks[i]
			block.right = block.left + len(block.elements) - 1

			// !对每个块的元素进行初始化2
			block.sorted = append([]int(nil), block.elements...)
			sort.Ints(block.sorted)
		}
	}

	// 区间更新
	update := func(left, right int, delta int) { // [l,r], starts at 0
		for i := range blocks {
			block := &blocks[i]
			if block.right < left {
				continue
			}
			if block.left > right {
				break
			}

			if left <= block.left && block.right <= right {
				// !区间更新完整的块:类似线段树，只需要打上懒标记
				block.lazyAdd += delta
			} else {
				bl := max(block.left, left)
				br := min(block.right, right)
				for j := bl - block.left; j <= br-block.left; j++ {
					// !区间修改不完整的块：暴力更新实际值之后再排序(保证总是查询前sorted是有序的)
					block.elements[j] += delta
				}
				// !重新排序
				block.sorted = append([]int(nil), block.elements...)
				sort.Ints(block.sorted)
			}
		}
	}

	// 区间查询
	query := func(left, right int, k int) (res int) { // [l,r], starts at 0
		for i := range blocks {
			block := &blocks[i]
			if block.right < left {
				continue
			}
			if block.left > right {
				break
			}

			if left <= block.left && block.right <= right {
				// !区间查询完整的块:k减去标记的add之后，二分查找 (black.size - bisectLeft(sorted,k-lazyAdd))
				pos := sort.SearchInts(block.sorted, k-block.lazyAdd)
				res += block.right - block.left + 1 - pos
			} else {
				bl := max(block.left, left)
				br := min(block.right, right)
				for j := bl - block.left; j <= br-block.left; j++ {
					// !区间查询不完整的块：暴力计算 实际值+懒标记里的值
					cur := block.elements[j] + block.lazyAdd
					if cur >= k {
						res++
					}
				}
			}
		}

		return
	}

	// !更新与查询
	sqrtInit(nums)
	for i := 0; i < q; i++ {
		var op string
		var left, right, x int
		fmt.Fscan(in, &op, &left, &right, &x)
		left, right = left-1, right-1
		if op == "M" {
			update(left, right, x) // !区间所有元素加x
		} else {
			fmt.Fprintln(out, query(left, right, x)) // !区间内有多少个元素大于等于x
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
	if a < b {
		return a
	}
	return b
}
