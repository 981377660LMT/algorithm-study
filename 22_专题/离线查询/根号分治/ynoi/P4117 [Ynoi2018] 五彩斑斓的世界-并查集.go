// P4117 [Ynoi2018] 五彩斑斓的世界(第二分块)
// https://www.luogu.com.cn/problem/P4117

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// 0 start end x：把区间 [start, end) 内大于 x 的数减去 x.
// 1 start end x: 查询区间 [start, end) 内等于 x 的数的个数.
//
// TODO:
//

func RangeClampRangeFreq(nums []int, operations [][4]int) []int {
	block := UseBlock(nums, int(math.Sqrt(float64(len(nums)))+1))
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount

	res := []int{}
	for _, op := range operations {
		kind := op[0]
		if kind == 0 {
			start, end, x := op[1], op[2], op[3]
			bid1, bid2 := belong[start], belong[end-1]
			if bid1 == bid2 {
			} else {
			}
		} else {
			start, end, x := op[1], op[2], op[3]
			bid1, bid2 := belong[start], belong[end-1]
			if bid1 == bid2 {
			} else {

			}
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
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	operations := make([][4]int, q)
	for i := range operations {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var start, end, x int
			fmt.Fscan(in, &start, &end, &x)
			start--
			operations[i] = [4]int{0, start, end, x}
		} else {
			var start, end, x int
			fmt.Fscan(in, &start, &end, &x)
			start--
			operations[i] = [4]int{1, start, end, x}
		}
	}

	res := RangeClampRangeFreq(nums, operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
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

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
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

func bisectLeft(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}
func bisectRight(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] <= target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}
