// P3203 [HNOI2010] 弹飞绵羊
// https://www.luogu.com.cn/problem/P3203

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// 每个弹力装置初始系数为nums[i].
// 当绵羊踩到第i个弹力装置时，会被弹到第i+nums[i]个弹力装置上。如果不存在第i+nums[i]个弹力装置，则绵羊会被弹飞。
//
// 1 index : 输出从inedx出发被弹几次后弹飞
// 2 index newValue : 将index位置的弹力装置系数改为newValue
// n<=2e5 q<=1e5

// 方法1：LCT
// !方法2：分块：
// 线性遍历太慢，需要`分块加速遍历`.
// !维护`跳出当前块后到达的位置`与`跳出当前块所用的步数`
// !需要从后往前维护。 修改时也是在块内从后往前。
func 弹飞绵羊(nums []int, operations [][3]int) []int {
	n := len(nums)
	B := UseBlock(n, int(math.Sqrt(float64(n))+1))
	belong, blockStart, blockEnd, blockCount := B.belong, B.blockStart, B.blockEnd, B.blockCount

	jumpTo := make([]int, n)
	jumpStep := make([]int, n)

	// 倒序重构块内信息
	rebuildBlock := func(bid int) {
		for i := blockEnd[bid] - 1; i >= blockStart[bid]; i-- {
			// 一次就跳出当前块，直接更新
			if i+nums[i] >= blockEnd[bid] {
				jumpTo[i] = i + nums[i]
				jumpStep[i] = 1
			} else {
				jumpTo[i] = jumpTo[i+nums[i]]
				jumpStep[i] = jumpStep[i+nums[i]] + 1
			}
		}
	}
	for bid := blockCount - 1; bid >= 0; bid-- {
		rebuildBlock(bid)
	}

	// 弹出序列需要多少步
	query := func(pos int) int {
		res := 0
		for pos < n {
			res += jumpStep[pos]
			pos = jumpTo[pos]
		}
		return res
	}

	res := []int{}
	for _, op := range operations {
		kind := op[0]
		if kind == 1 {
			cur := op[1]
			res = append(res, query(cur))
		} else {
			cur, newValue := op[1], op[2]
			nums[cur] = newValue
			rebuildBlock(belong[cur])
		}
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	var q int
	fmt.Fscan(in, &q)
	var operations [][3]int
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var index int
			fmt.Fscan(in, &index)
			operations = append(operations, [3]int{op, index, 0})
		} else {
			var index, newValue int
			fmt.Fscan(in, &index, &newValue)
			operations = append(operations, [3]int{op, index, newValue})
		}
	}

	res := 弹飞绵羊(nums, operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}
