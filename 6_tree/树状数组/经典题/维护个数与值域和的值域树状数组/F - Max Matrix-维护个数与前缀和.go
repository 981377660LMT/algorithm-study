package main

import (
	"bufio"
	"fmt"
	"os"
)

// # TODO
// # Run time Error
func maxMatrix(n, m int, queries [][3]int) []int {
	nums1, nums2 := make([]int, n), make([]int, m)
	sum1, sum2 := CreateBIT(int(1e8+10)), CreateBIT(int(1e8+10))
	count1, count2 := CreateBIT(int(1e8+10)), CreateBIT(int(1e8+10))

	count1.Add(0, 0, n)
	count2.Add(0, 0, m)
	cur := 0
	res := make([]int, 0, len(queries))
	for _, q := range queries {
		op, qi, qv := q[0], q[1], q[2]
		var curNums []int
		var curCount, curSum BIT
		var otherCount, otherSum BIT
		if op == 2 {
			curNums, curCount, curSum = nums1, count1, sum1
			otherCount, otherSum = count2, sum2
		} else {
			curNums, curCount, curSum = nums2, count2, sum2
			otherCount, otherSum = count1, sum1
		}

		// 移除旧值
		preNum := curNums[qi]
		// (as min/max)
		cur -= otherCount.Query(0, preNum)*preNum + otherSum.Query(preNum+1, int(1e8+10))
		cur += otherCount.Query(0, qv)*qv + otherSum.Query(qv+1, int(1e8+10))
		curCount.Add(preNum, preNum, -1)
		curSum.Add(preNum, preNum, -preNum)
		// 加入新值
		curNums[qi] = qv
		curCount.Add(qv, qv, 1)
		curSum.Add(qv, qv, qv)

		res = append(res, cur)
	}

	return res

}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)
	Q := make([][3]int, q)
	for i := range Q {
		fmt.Fscan(in, &Q[i][0], &Q[i][1], &Q[i][2])
		Q[i][1]--
	}

	res := maxMatrix(n, m, Q)
	for _, v := range res {
		fmt.Fprintln(out, v)
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

func (bit *mapBIT) Add(left, right, delta int) {
	left, right = left+1, right+1
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *sliceBIT) Add(left, right, delta int) {
	left, right = left+1, right+1
	bit.add(left, delta)
	bit.add(right+1, -delta)
}

func (bit *mapBIT) Query(left, right int) int {
	left, right = left+1, right+1
	return bit.query(right) - bit.query(left-1)
}

func (bit *sliceBIT) Query(left, right int) int {
	left, right = left+1, right+1
	return bit.query(right) - bit.query(left-1)
}

func (bit *mapBIT) add(index, delta int) {
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

func (bit *sliceBIT) add(index, delta int) {
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
