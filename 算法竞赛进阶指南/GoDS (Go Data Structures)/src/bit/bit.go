package bit

// 区间修改、区间查询的树状数组

import (
	"bufio"
	"fmt"
	"os"
)

// https://atcoder.jp/contests/practice2/tasks/practice2_b
func demo() {
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

func CreateBIT(n int) BIT {
	if n <= int(1e6) {
		return newSliceBIT(n)
	}

	return newMapBIT(n)
}

type BIT interface {
	// 区间 [left, right] 里的每个数增加 delta
	//  1 <= left <= right <= n
	Add(left, right, delta int)

	// 查询区间 [left, right] 的和
	Query(left, right int) int
}

type mapBIT struct {
	n     int
	tree1 map[int]int
	tree2 map[int]int
}

type sliceBIT struct {
	n     int
	tree1 []int
	tree2 []int
}

func newMapBIT(n int) *mapBIT {
	return &mapBIT{
		n:     n,
		tree1: make(map[int]int, 1<<10),
		tree2: make(map[int]int, 1<<10),
	}
}

func newSliceBIT(n int) *sliceBIT {
	return &sliceBIT{
		n:     n,
		tree1: make([]int, n+1),
		tree2: make([]int, n+1),
	}
}

// left>=1
func (bit *mapBIT) Add(left, right, delta int) {
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

// left>=1
func (bit *sliceBIT) Add(left, right, delta int) {
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

// left>=1
func (bit *mapBIT) Query(left, right int) int {
	return bit.query(right) - bit.query(left-1)
}

// left>=1
func (bit *sliceBIT) Query(left, right int) int {
	return bit.query(right) - bit.query(left-1)
}

// index >= 1
func (bit *mapBIT) add(index, delta int) {
	rawIndex := index
	for index <= bit.n {
		bit.tree1[index] += delta
		bit.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

// index >= 1
func (bit *sliceBIT) add(index, delta int) {
	rawIndex := index
	for index <= bit.n {
		bit.tree1[index] += delta
		bit.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (bit *mapBIT) query(index int) int {
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

func (bit *sliceBIT) query(index int) int {
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

func (bit *mapBIT) String() string {
	return "not implemented"
}

func (bit *sliceBIT) String() string {
	nums := make([]int, bit.n+1)
	for i := 0; i < bit.n; i++ {
		nums[i+1] = bit.Query(i+1, i+1)
	}
	return fmt.Sprint(nums)
}
