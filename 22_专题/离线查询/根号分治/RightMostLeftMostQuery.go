// !左侧右侧包含当前位置.

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
)

func main() {
	// checkWithBruteForce()
	abc382c()
}

// C - Kaiten Sushi
// https://atcoder.jp/contests/abc382/tasks/abc382_c
// 回转寿司
//
// 在某个回转寿司店，有 N 人从 1 到 N 的编号到访。人 i 的美食度是 A i ​ 。
// 传送带上流动着 M 个寿司。第 j 个流动的寿司的美味程度是 B j ​ 。
// 每个寿司依次流过每个人 1,2,…,N 的面前。
// 每个人在美味程度>=自己美食度的寿司流到自己面前时会取走并食用该寿司，其他情况下则不做任何事情。
// 请确定每个寿司由谁食用，或者是否没有人食用。
//
// 等价于：对于每个寿司，找到第一个美食度大于等于它的人。
func abc382c() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M int32
	fmt.Fscan(in, &N, &M)
	A, B := make([]int, N), make([]int, M)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &A[i])
	}
	for i := int32(0); i < M; i++ {
		fmt.Fscan(in, &B[i])
	}

	res := make([]int32, M)
	Q := NewRightMostLeftMostQuery(A)
	for i := int32(0); i < M; i++ {
		res[i] = int32(Q.RightNearestFloor(0, B[i]))
	}
	for _, v := range res {
		if v == -1 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, v+1)
		}
	}
}

type RightMostLeftMostQuery struct {
	_nums       []int
	_belong     []int32
	_blockStart []int32
	_blockEnd   []int32
	_blockCount int32
	_blockMin   []int
	_blockMax   []int
	_blockLazy  []int
}

// 对每个下标，`O(sqrt)`查询 最右侧/最左侧 lower/floor/ceiling/higher 的元素.
func NewRightMostLeftMostQuery(arr []int) *RightMostLeftMostQuery {
	arr = append(arr[:0:0], arr...)
	n := int32(len(arr))
	blockSize := (int32(math.Sqrt(float64(n))) + 1)
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int32, blockCount)
	blockEnd := make([]int32, blockCount)
	belong := make([]int32, n)
	for i := int32(0); i < blockCount; i++ {
		blockStart[i] = i * blockSize
		blockEnd[i] = min32((i+1)*blockSize, n)
	}
	for i := int32(0); i < n; i++ {
		belong[i] = (i / blockSize)
	}
	res := &RightMostLeftMostQuery{
		_nums:       arr,
		_belong:     belong,
		_blockStart: blockStart,
		_blockEnd:   blockEnd,
		_blockCount: blockCount,
		_blockMin:   make([]int, blockCount),
		_blockMax:   make([]int, blockCount),
		_blockLazy:  make([]int, blockCount),
	}
	for bid := int32(0); bid < blockCount; bid++ {
		res._rebuildBlock(bid)
	}
	return res
}

func (rm *RightMostLeftMostQuery) Get(index int32) int {
	if index < 0 || index >= int32(len(rm._nums)) {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	return rm._nums[index] + rm._blockLazy[rm._belong[index]]
}

func (rm *RightMostLeftMostQuery) Set(index int32, value int) {
	if index < 0 || index >= int32(len(rm._nums)) {
		return
	}
	bid := rm._belong[index]
	lazy := rm._blockLazy[bid]
	pre := rm._nums[index] + lazy
	if pre == value {
		return
	}
	rm._nums[index] = value - lazy
	rm._rebuildBlock(bid)
}

func (rm *RightMostLeftMostQuery) AddRange(start, end int32, delta int) {
	if start < 0 {
		start = 0
	}
	if end > int32(len(rm._nums)) {
		end = int32(len(rm._nums))
	}
	if start >= end {
		return
	}
	bid1 := rm._belong[start]
	bid2 := rm._belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			rm._nums[i] += delta
		}
		rm._rebuildBlock(bid1)
	} else {
		for i := start; i < rm._blockEnd[bid1]; i++ {
			rm._nums[i] += delta
		}
		rm._rebuildBlock(bid1)
		for bid := bid1 + 1; bid < bid2; bid++ {
			rm._blockLazy[bid] += delta
		}
		for i := rm._blockStart[bid2]; i < end; i++ {
			rm._nums[i] += delta
		}
		rm._rebuildBlock(bid2)
	}
}

// 查询`index`右侧最远的下标`j`，使得 `nums[j] < nums[index]`.
// 如果不存在，返回`-1`.
func (rm *RightMostLeftMostQuery) RightMostLower(index int32, target int) int32 {
	return rm._queryRightMost(
		index,
		func(bid int32) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] < target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] < target
		},
	)
}

func (rm *RightMostLeftMostQuery) RightMostFloor(index int32, target int) int32 {
	return rm._queryRightMost(
		index,
		func(bid int32) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] <= target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] <= target
		},
	)
}

func (rm *RightMostLeftMostQuery) RightMostCeiling(index int32, target int) int32 {
	return rm._queryRightMost(
		index,
		func(bid int32) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] >= target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] >= target
		},
	)
}

func (rm *RightMostLeftMostQuery) RightMostHigher(index int32, target int) int32 {
	return rm._queryRightMost(
		index,
		func(bid int32) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] > target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] > target
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftMostLower(index int32, target int) int32 {
	return rm._queryLeftMost(
		index,
		func(bid int32) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] < target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] < target
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftMostFloor(index int32, target int) int32 {
	return rm._queryLeftMost(
		index,
		func(bid int32) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] <= target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] <= target
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftMostCeiling(index int32, target int) int32 {
	return rm._queryLeftMost(
		index,
		func(bid int32) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] >= target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] >= target
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftMostHigher(index int32, target int) int32 {
	return rm._queryLeftMost(
		index,
		func(bid int32) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] > target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] > target
		},
	)
}

func (rm *RightMostLeftMostQuery) RightNearestLower(index int32, target int) int32 {
	return rm._queryRightNearest(
		index,
		func(bid int32) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] < target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] < target
		},
	)
}

func (rm *RightMostLeftMostQuery) RightNearestFloor(index int32, target int) int32 {
	return rm._queryRightNearest(
		index,
		func(bid int32) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] <= target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] <= target
		},
	)
}

func (rm *RightMostLeftMostQuery) RightNearestCeiling(index int32, target int) int32 {
	return rm._queryRightNearest(
		index,
		func(bid int32) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] >= target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] >= target
		},
	)
}

func (rm *RightMostLeftMostQuery) RightNearestHigher(index int32, target int) int32 {
	return rm._queryRightNearest(
		index,
		func(bid int32) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] > target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] > target
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftNearestLower(index int32, target int) int32 {
	return rm._queryLeftNearest(
		index,
		func(bid int32) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] < target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] < target
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftNearestFloor(index int32, target int) int32 {
	return rm._queryLeftNearest(
		index,
		func(bid int32) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] <= target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] <= target
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftNearestCeiling(index int32, target int) int32 {
	return rm._queryLeftNearest(
		index,
		func(bid int32) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] >= target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] >= target
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftNearestHigher(index int32, target int) int32 {
	return rm._queryLeftNearest(
		index,
		func(bid int32) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] > target
		},
		func(eid, bid int32) bool {
			return rm._nums[eid]+rm._blockLazy[bid] > target
		},
	)
}

func (rm *RightMostLeftMostQuery) _queryRightMost(
	pos int32,
	predicateBlock func(bid int32) bool,
	predicateElement func(eid, bid int32) bool,
) int32 {
	bid := rm._belong[pos]
	for i := rm._blockCount - 1; i > bid; i-- {
		if !predicateBlock(i) {
			continue
		}
		for j := rm._blockEnd[i] - 1; j >= rm._blockStart[i]; j-- {
			if predicateElement(j, i) {
				return j
			}
		}
	}
	for i := rm._blockEnd[bid] - 1; i >= pos; i-- {
		if predicateElement(i, bid) {
			return i
		}
	}
	return -1
}

func (rm *RightMostLeftMostQuery) _queryLeftMost(
	pos int32,
	predicateBlock func(bid int32) bool,
	predicateElement func(eid, bid int32) bool,
) int32 {
	bid := rm._belong[pos]
	for i := int32(0); i < bid; i++ {
		if !predicateBlock(i) {
			continue
		}
		for j := rm._blockStart[i]; j < rm._blockEnd[i]; j++ {
			if predicateElement(j, i) {
				return j
			}
		}
	}
	for i := rm._blockStart[bid]; i <= pos; i++ {
		if predicateElement(i, bid) {
			return i
		}
	}
	return -1
}

func (rm *RightMostLeftMostQuery) _queryRightNearest(
	pos int32,
	predicateBlock func(bid int32) bool,
	predicateElement func(eid, bid int32) bool,
) int32 {
	bid := rm._belong[pos]
	for i := pos; i < rm._blockEnd[bid]; i++ {
		if predicateElement(i, bid) {
			return i
		}
	}
	for i := bid + 1; i < rm._blockCount; i++ {
		if !predicateBlock(i) {
			continue
		}
		for j := rm._blockStart[i]; j < rm._blockEnd[i]; j++ {
			if predicateElement(j, i) {
				return j
			}
		}
	}
	return -1
}

func (rm *RightMostLeftMostQuery) _queryLeftNearest(
	pos int32,
	predicateBlock func(bid int32) bool,
	predicateElement func(eid, bid int32) bool,
) int32 {
	bid := rm._belong[pos]
	for i := pos; i >= rm._blockStart[bid]; i-- {
		if predicateElement(i, bid) {
			return i
		}
	}
	for i := bid - 1; i >= 0; i-- {
		if !predicateBlock(i) {
			continue
		}
		for j := rm._blockEnd[i] - 1; j >= rm._blockStart[i]; j-- {
			if predicateElement(j, i) {
				return j
			}
		}
	}
	return -1
}

func (rm *RightMostLeftMostQuery) _rebuildBlock(bid int32) {
	rm._blockMin[bid] = math.MaxInt64
	rm._blockMax[bid] = math.MinInt64
	lazy := rm._blockLazy[bid]
	rm._blockLazy[bid] = 0
	for i := rm._blockStart[bid]; i < rm._blockEnd[bid]; i++ {
		rm._nums[i] += lazy
		rm._blockMin[bid] = min(rm._blockMin[bid], rm._nums[i])
		rm._blockMax[bid] = max(rm._blockMax[bid], rm._nums[i])
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

// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
type Mocker struct {
	nums []int
}

func NewMocker(nums []int) *Mocker {
	return &Mocker{nums: nums}
}

func (m *Mocker) Set(index int32, value int) {
	m.nums[index] = value
}
func (m *Mocker) AddRange(start, end int32, delta int) {
	for i := start; i < end; i++ {
		m.nums[i] += delta
	}
}

func (m *Mocker) RightMostLower(index int32) int32 {
	cur := m.nums[index]
	for i := int32(len(m.nums) - 1); i >= index; i-- {
		if m.nums[i] < cur {
			return i
		}
	}
	return -1
}
func (m *Mocker) RightMostFloor(index int32) int32 {
	cur := m.nums[index]
	for i := int32(len(m.nums) - 1); i >= index; i-- {
		if m.nums[i] <= cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) RightMostCeiling(index int32) int32 {
	cur := m.nums[index]
	for i := int32(len(m.nums) - 1); i >= index; i-- {
		if m.nums[i] >= cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) RightMostHigher(index int32) int32 {
	cur := m.nums[index]
	for i := int32(len(m.nums) - 1); i >= index; i-- {
		if m.nums[i] > cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) LeftMostLower(index int32) int32 {
	cur := m.nums[index]
	for i := int32(0); i <= index; i++ {
		if m.nums[i] < cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) LeftMostFloor(index int32) int32 {
	cur := m.nums[index]
	for i := int32(0); i <= index; i++ {
		if m.nums[i] <= cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) LeftMostCeiling(index int32) int32 {
	cur := m.nums[index]
	for i := int32(0); i <= index; i++ {
		if m.nums[i] >= cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) LeftMostHigher(index int32) int32 {
	cur := m.nums[index]
	for i := int32(0); i <= index; i++ {
		if m.nums[i] > cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) RightNearestLower(index int32) int32 {
	cur := m.nums[index]
	for i := index; i < int32(len(m.nums)); i++ {
		if m.nums[i] < cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) RightNearestFloor(index int32) int32 {
	cur := m.nums[index]
	for i := index; i < int32(len(m.nums)); i++ {
		if m.nums[i] <= cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) RightNearestCeiling(index int32) int32 {
	cur := m.nums[index]
	for i := index; i < int32(len(m.nums)); i++ {
		if m.nums[i] >= cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) RightNearestHigher(index int32) int32 {
	cur := m.nums[index]
	for i := index; i < int32(len(m.nums)); i++ {
		if m.nums[i] > cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) LeftNearestLower(index int32) int32 {
	cur := m.nums[index]
	for i := index; i >= 0; i-- {
		if m.nums[i] < cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) LeftNearestFloor(index int32) int32 {
	cur := m.nums[index]
	for i := index; i >= 0; i-- {
		if m.nums[i] <= cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) LeftNearestCeiling(index int32) int32 {
	cur := m.nums[index]
	for i := index; i >= 0; i-- {
		if m.nums[i] >= cur {
			return i
		}
	}
	return -1
}

func (m *Mocker) LeftNearestHigher(index int32) int32 {
	cur := m.nums[index]
	for i := index; i >= 0; i-- {
		if m.nums[i] > cur {
			return i
		}
	}
	return -1
}

// checkWithBruteForce
func checkWithBruteForce() {
	N := int(1e5)
	MAX := int(1e9)
	randNums := make([]int, N)
	for i := 0; i < N; i++ {
		randNums[i] = rand.Intn(MAX)
	}
	mocker := NewMocker(append([]int{}, randNums...))
	real := NewRightMostLeftMostQuery(append([]int{}, randNums...))

	adapter := func(f func(index int32, target int) int32, dep *RightMostLeftMostQuery) func(index int32) int32 {
		return func(index int32) int32 {
			return f(index, dep.Get(index))
		}
	}
	debug := func(f1 func(int32) int32, f2 func(int32, int) int32) {
		index := int32(rand.Intn(N))
		res1 := f1(index)
		res2 := adapter(f2, real)(index)
		if res1 != res2 {
			panic("error")
		}
	}

	for i := 0; i < 100000; i++ {
		op := rand.Intn(18)
		switch op {
		case 0:
			index := int32(rand.Intn(N))
			value := rand.Intn(MAX)
			mocker.Set(index, value)
			real.Set(index, value)

		case 1:
			start := int32(rand.Intn(N))
			end := int32(start + int32(rand.Intn(N-int(start))))
			delta := rand.Intn(MAX)
			mocker.AddRange(start, end, delta)
			real.AddRange(start, end, delta)
		case 2:
			debug(mocker.RightMostLower, real.RightMostLower)
		case 3:
			debug(mocker.RightMostFloor, real.RightMostFloor)
		case 4:
			debug(mocker.RightMostCeiling, real.RightMostCeiling)
		case 5:
			debug(mocker.RightMostHigher, real.RightMostHigher)
		case 6:
			debug(mocker.LeftMostLower, real.LeftMostLower)
		case 7:
			debug(mocker.LeftMostFloor, real.LeftMostFloor)
		case 8:
			debug(mocker.LeftMostCeiling, real.LeftMostCeiling)
		case 9:
			debug(mocker.LeftMostHigher, real.LeftMostHigher)
		case 10:
			debug(mocker.RightNearestLower, real.RightNearestLower)
		case 11:
			debug(mocker.RightNearestFloor, real.RightNearestFloor)
		case 12:
			debug(mocker.RightNearestCeiling, real.RightNearestCeiling)
		case 13:
			debug(mocker.RightNearestHigher, real.RightNearestHigher)
		case 14:
			debug(mocker.LeftNearestLower, real.LeftNearestLower)
		case 15:
			debug(mocker.LeftNearestFloor, real.LeftNearestFloor)
		case 16:
			debug(mocker.LeftNearestCeiling, real.LeftNearestCeiling)
		case 17:
			debug(mocker.LeftNearestHigher, real.LeftNearestHigher)
		}
	}

	fmt.Println("done")
}
