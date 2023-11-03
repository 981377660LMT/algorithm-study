package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func 弹飞绵羊(nums []int, operations [][3]int) [][2]int {
	n := len(nums)
	B := UseBlock(nums, int(math.Sqrt(float64(n))+1))
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
				// 否则继承同一个块中下一个跳到的位置
				jumpTo[i] = jumpTo[i+nums[i]]
				jumpStep[i] = jumpStep[i+nums[i]] + 1
			}
		}
	}
	for bid := blockCount - 1; bid >= 0; bid-- {
		rebuildBlock(bid)
	}

	// 弹出序列需要多少步
	query := func(pos int) [2]int {
		res := 0
		last := pos
		for pos < n {
			res += jumpStep[pos]
			pos = jumpTo[pos]
			if pos < n {
				last = pos
			}
		}
		return [2]int{last, res}
	}

	res := [][2]int{}
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

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	var operations [][3]int
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var index int
			fmt.Fscan(in, &index)
			index--
			operations = append(operations, [3]int{op, index, 0})
		} else {
			var index, newValue int
			fmt.Fscan(in, &index, &newValue)
			index--
			operations = append(operations, [3]int{op, index, newValue})
		}
	}

	res := 弹飞绵羊(nums, operations)
	for _, v := range res {
		fmt.Fprintln(out, v[0]+1, v[1])
	}
}

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(nums []int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	n := len(nums)

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
