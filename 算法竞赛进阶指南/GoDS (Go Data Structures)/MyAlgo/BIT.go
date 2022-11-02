// 区间修改、区间查询的树状数组

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://atcoder.jp/contests/practice2/tasks/practice2_b
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

	bit := CreateBIT(n)
	for i := 0; i < n; i++ {
		bit.Add(i+1, i+1, nums[i])
	}

	for i := 0; i < q; i++ {
		var opt uint8
		var left, deltaOrRight int // 0<=left<=right<=n
		fmt.Fscan(in, &opt, &left, &deltaOrRight)
		left++
		if opt == 0 {
			delta := deltaOrRight
			bit.Add(left, left, delta)
		} else {
			right := deltaOrRight
			fmt.Fprintln(out, bit.Query(left, right))
		}
	}
}

type BIT interface {
	// 区间 [left, right] 里的每个数增加 delta
	Add(left, right, delta int)

	// 查询区间 [left, right] 的和
	Query(left, right int) int
}

func CreateBIT(n int) BIT {
	if n <= int(1e6) {
		return NewBITSlice(n)
	}

	return NewBITMap(n)
}

type BITMap struct {
	n     int
	tree1 map[int]int
	tree2 map[int]int
}

type BITSlice struct {
	n     int
	tree1 []int
	tree2 []int
}

func NewBITMap(n int) *BITMap {
	return &BITMap{
		n:     n,
		tree1: make(map[int]int, 1<<10),
		tree2: make(map[int]int, 1<<10),
	}
}

func NewBITSlice(n int) *BITSlice {
	return &BITSlice{
		n:     n,
		tree1: make([]int, n+10),
		tree2: make([]int, n+10),
	}
}

func (bit *BITMap) Add(left, right, delta int) {
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *BITSlice) Add(left, right, delta int) {
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *BITMap) Query(left, right int) int {
	return bit.query(right) - bit.query(left-1)
}

func (bit *BITSlice) Query(left, right int) int {
	return bit.query(right) - bit.query(left-1)
}

func (bit *BITMap) add(index, delta int) {
	if index <= 0 {
		errorInfo := fmt.Sprintf("index must be greater than 0, but got %d", index)
		panic(errorInfo)
	}

	rawIndex := index
	for index <= bit.n {
		bit.tree1[index] += delta
		bit.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (bit *BITSlice) add(index, delta int) {
	if index <= 0 {
		errorInfo := fmt.Sprintf("index must be greater than 0, but got %d", index)
		panic(errorInfo)
	}

	rawIndex := index
	for index <= bit.n {
		bit.tree1[index] += delta
		bit.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (bit *BITMap) query(index int) int {
	if index > bit.n {
		index = bit.n
	}

	rawIndex := index
	res := 0
	for index > 0 {
		res += rawIndex*bit.tree1[index] - bit.tree2[index]
		index -= index & -index
	}
	return res
}

func (bit *BITSlice) query(index int) int {
	if index > bit.n {
		index = bit.n
	}

	rawIndex := index
	res := 0
	for index > 0 {
		res += rawIndex*bit.tree1[index] - bit.tree2[index]
		index -= index & -index
	}
	return res
}

func (bit *BITMap) String() string {
	return "not implemented"
}

func (bit *BITSlice) String() string {
	nums := make([]int, bit.n+1)
	for i := 0; i < bit.n; i++ {
		nums[i+1] = bit.Query(i+1, i+1)
	}
	return fmt.Sprint(nums)
}
