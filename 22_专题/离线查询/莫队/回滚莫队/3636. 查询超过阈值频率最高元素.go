// 3636. 查询超过阈值频率最高元素
// https://leetcode.cn/problems/threshold-majority-queries/description/
// 给你一个长度为 n 的整数数组 nums 和一个查询数组 queries，其中 queries[i] = [li, ri, thresholdi]。
//
// 返回一个整数数组 ans，其中 ans[i] 等于子数组 nums[li...ri] 中出现 至少 thresholdi 次的元素，选择频率 最高 的元素（如果频率相同则选择 最小 的元素），如果不存在这样的元素则返回 -1。

package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	// nums = [1,1,2,2,1,1], queries = [[0,5,4],[0,3,3],[2,3,2]]
	fmt.Println(subarrayMajority([]int{1, 1, 2, 2, 1, 1}, [][]int{{0, 5, 4}, {0, 3, 3}, {2, 3, 2}}))
}

// 区间众数查询，存在多个众数时，返回值最小的.
func subarrayMajority(nums []int, queries [][]int) []int {
	n, q := len(nums), len(queries)
	mo := NewMoRollback(n, q)
	for _, query := range queries {
		left, right := query[0], query[1]
		mo.AddQuery(left, right+1)
	}

	newNums, origin := Discretize(nums)
	size := int32(len(newNums))

	res := make([]int, q)

	counter := make([]int, size)
	history := make([]int32, 0, n)

	maxCount, maxKey := 0, int32(0)
	snapState := 0
	snapCount, snapKey := 0, int32(0)

	add := func(index int) {
		x := newNums[index]
		history = append(history, x)
		counter[x]++
		// 存在多个众数时，返回值最小的
		if counter[x] > maxCount || (counter[x] == maxCount && x < maxKey) {
			maxCount = counter[x]
			maxKey = x
		}
	}

	reset := func() {
		for _, v := range history {
			counter[v] = 0
		}
		history = history[:0]
		maxCount, maxKey = 0, 0
	}

	snapshot := func() {
		snapState = len(history)
		snapCount, snapKey = maxCount, maxKey
	}

	rollback := func() {
		for len(history) > snapState {
			x := history[len(history)-1]
			history = history[:len(history)-1]
			counter[x]--
		}
		maxCount, maxKey = snapCount, snapKey
	}

	query := func(qi int) {
		threshold := queries[qi][2]
		if maxCount >= threshold {
			res[qi] = origin[maxKey] // 返回原始值
		} else {
			res[qi] = -1 // 不存在满足条件的元素
		}
	}

	mo.Run(add, add, reset, snapshot, rollback, query, -1)
	return res
}

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}

type MoRollback struct {
	left, right []int32
}

func NewMoRollback(n, q int) *MoRollback {
	return &MoRollback{left: make([]int32, 0, q), right: make([]int32, 0, q)}
}

func (mo *MoRollback) AddQuery(start, end int) {
	mo.left = append(mo.left, int32(start))
	mo.right = append(mo.right, int32(end))
}

// addLeft : 将index位置的元素加入到区间左端.
// addRight: 将index位置的元素加入到区间右端.
// reset: 重置区间.
// snapShot: 快照当前状态.
// rollback: 回滚到快照状态.
// query: 查询当前区间.
// blockSize: 分块大小.-1表示使用默认值.
func (mo *MoRollback) Run(
	addLeft func(index int),
	addRight func(index int),
	reset func(),
	snapShot func(),
	rollback func(),
	query func(qid int),
	blockSize int,
) {
	q := int32(len(mo.left))
	if q == 0 {
		return
	}
	n := int32(0)
	for _, v := range mo.right {
		n = max32(n, v)
	}
	blockSize32 := int32(blockSize)
	if blockSize32 == -1 {
		blockSize32 = int32(max32(1, n/max32(1, int32(math.Sqrt(float64(q*2/3))))))
	}
	queryId := make([][]int32, (n-1)/blockSize32+1)
	naive := func(qid int32) {
		snapShot()
		for i := mo.left[qid]; i < mo.right[qid]; i++ {
			addRight(int(i))
		}
		query(int(qid))
		rollback()
	}

	for qid := int32(0); qid < q; qid++ {
		l, r := mo.left[qid], mo.right[qid]
		iL, iR := l/blockSize32, r/blockSize32
		if iL == iR {
			naive(qid)
			continue
		}
		queryId[iL] = append(queryId[iL], qid)
	}

	for _, order := range queryId {
		if len(order) == 0 {
			continue
		}
		sort.Slice(order, func(i, j int) bool {
			return mo.right[order[i]] < mo.right[order[j]]
		})
		lMax := int32(0)
		for _, qid := range order {
			lMax = max32(lMax, mo.left[qid])
		}
		reset()
		l, r := lMax, lMax
		for _, qid := range order {
			L, R := mo.left[qid], mo.right[qid]
			for r < R {
				addRight(int(r))
				r++
			}
			snapShot()
			for L < l {
				l--
				addLeft(int(l))
			}
			query(int(qid))
			rollback()
			l = lMax
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
