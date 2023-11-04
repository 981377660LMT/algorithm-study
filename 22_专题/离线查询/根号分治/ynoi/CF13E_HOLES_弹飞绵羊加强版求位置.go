package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// 最后一次落在哪个洞,球被弹出前共被弹了多少次.
func 弹飞绵羊2(nums []int, operations [][3]int) [][2]int {
	n := len(nums)
	B := UseBlock(nums, int(math.Sqrt(float64(n))+1))
	belong, blockStart, blockEnd, blockCount := B.belong, B.blockStart, B.blockEnd, B.blockCount

	jumpTo := make([]int, n)   // 跳出当前块后到达的位置
	jumpStep := make([]int, n) // 跳出当前块所用的步数

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
	// !由于每一个点记录的只是跳出当前块后到达的位置，所以求得的最终位置并不是实际的位置，还要再往后找.
	query := func(pos int) [2]int {
		res := 0
		for nextPos := jumpTo[pos]; nextPos < n; nextPos = jumpTo[pos] {
			res += jumpStep[pos]
			pos = nextPos
		}
		for nextPos := pos + nums[pos]; nextPos < n; nextPos = pos + nums[pos] {
			res++
			pos = nextPos
		}
		res++
		return [2]int{pos, res}
	}

	res := [][2]int{}
	for _, op := range operations {
		kind := op[0]
		if kind == 0 {
			pos, newValue := op[1], op[2]
			nums[pos] = newValue
			rebuildBlock(belong[pos])
		} else {
			pos := op[1]
			res = append(res, query(pos))
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
		if op == 0 {
			var index, newValue int
			fmt.Fscan(in, &index, &newValue)
			index--
			operations = append(operations, [3]int{op, index, newValue})
		} else {
			var index int
			fmt.Fscan(in, &index)
			index--
			operations = append(operations, [3]int{op, index, 0})
		}
	}

	res := 弹飞绵羊2(nums, operations)
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
