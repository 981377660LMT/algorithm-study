// 区间最小众数(RangeMode)

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// 给出一个长为 n 的数列，以及 n 个操作，操作涉及询问区间的最小众数。
// 第一行输入一个数字 n。
// 第二行输入 n 个数字，第 i 个数字为 a_i，以空格隔开。
// 接下来输入 n 行询问，每行输入两个数字 l、r，以空格隔开。
// 表示查询位于 [l,r] 的数字的众数。如果有多个众数，输出最小的那个。
//
// !1.O(nsqrt(n))预处理从第i块到第j块的区间最小众数(0<=i<=j<blockCount).
// !2.![start,end)的区间众数，一定是在[start,end)中一段连续的整块的众数和两边非完整块的数的并集内。

func main() {
	// https://loj.ac/p/6277
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	block := UseBlock(nums)
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount

	D := NewDictionary()
	newNums := make([]int, n)
	mp := make([][]int, n+1) // n 表示不存在的数字.
	for i := range nums {
		newNums[i] = D.Id(nums[i])
		mp[newNums[i]] = append(mp[newNums[i]], i)
	}

	// [start,end)区间内等于target的元素个数.
	rangeFreq := func(start, end int, target int) int {
		pos := mp[target]
		return sort.SearchInts(pos, end) - sort.SearchInts(pos, start)
	}

	// !O(nsqrt(n))预处理从第i块到第j块的区间最小众数(0<=i<=j<blockCount).
	blockMode := make([][]int, blockCount)
	for i := 0; i < blockCount; i++ {
		blockMode[i] = make([]int, blockCount)
		for j := 0; j < blockCount; j++ {
			blockMode[i][j] = n
		}
	}
	for bid1 := 0; bid1 < blockCount; bid1++ {
		counter := make([]int, D.Size())
		res, max_ := 0, 0
		for j := blockStart[bid1]; j < n; j++ {
			bid2 := belong[j]
			counter[newNums[j]]++
			if counter[newNums[j]] > max_ || (counter[newNums[j]] == max_ && D.Value(newNums[j]) < D.Value(res)) {
				max_ = counter[newNums[j]]
				res = newNums[j]
			}
			blockMode[bid1][bid2] = res
		}
	}

	// 查询[start,end)区间的最小众数.
	// ![start,end)的区间众数，一定是在[start,end)中一段连续的整块的众数和两边非完整块的数的并集内。
	queryMode := func(start, end int) int {
		bid1, bid2 := belong[start], belong[end-1]
		res, max_ := 0, 0
		if bid1 == bid2 {
			for i := start; i < end; i++ {
				freq := rangeFreq(start, end, newNums[i])
				if freq > max_ || (freq == max_ && D.Value(newNums[i]) < D.Value(res)) {
					max_ = freq
					res = newNums[i]
				}
			}
		} else {
			res = blockMode[bid1+1][bid2-1]
			max_ = rangeFreq(start, end, res)
			for i := start; i < blockEnd[bid1]; i++ {
				freq := rangeFreq(start, end, newNums[i])
				if freq > max_ || (freq == max_ && D.Value(newNums[i]) < D.Value(res)) {
					res = newNums[i]
					max_ = freq
				}
			}
			for i := blockStart[bid2]; i < end; i++ {
				freq := rangeFreq(start, end, newNums[i])
				if freq > max_ || (freq == max_ && D.Value(newNums[i]) < D.Value(res)) {
					res = newNums[i]
					max_ = freq
				}
			}
		}

		return D.Value(res)
	}

	for i := 0; i < n; i++ {
		var left, right int
		fmt.Fscan(in, &left, &right)
		if left > right {
			left, right = right, left
		}
		left--
		fmt.Println(queryMode(left, right))
	}
}

func UseBlock(nums []int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	n := len(nums)
	// blockSize := int(math.Sqrt(float64(n)) + 1)
	blockSize := int((math.Sqrt(float64(n)) / 4) + 1)
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

type Dictionary struct {
	_idToValue []int
	_valueToId map[int]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[int]int{},
	}
}

func (d *Dictionary) Id(value int) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}

func (d *Dictionary) Value(id int) (value int) {
	return d._idToValue[id]
}

func (d *Dictionary) Size() int {
	return len(d._idToValue)
}
