// 小偷问题
// 维护两个序列 a 和 b，支持以下操作：
// 1. 插入新元素 a[i] 和 b[i]（a[i] 应当是 a 中的最大元素）
// 2. 给定 x，对所有 i，将 a[i] 设为 min(a[i], x)
// 3. 查询 max(a[i]+b[i])
// 每个操作的均摊时间复杂度为 O(1).

package main

func main() {
	C := NewMinPairContainer(10)
	C.Add(1, 2)
	C.Add(2, 3)
	C.Add(3, 4)
	C.Update(2)
}

const INF int = 4e18

type MinPairContainer struct {
	as, bs []int
	size   int32
}

func NewMinPairContainer(n int) *MinPairContainer {
	return &MinPairContainer{
		as: make([]int, n),
		bs: make([]int, n),
	}
}

func (mpc *MinPairContainer) Add(a, b int) {
	for mpc.size > 1 {
		mpc._pushDown()
		e1 := mpc._event(mpc.as[mpc.size-2], mpc.bs[mpc.size-2], mpc.bs[mpc.size-1])
		e2 := mpc._event(mpc.as[mpc.size-1], mpc.bs[mpc.size-1], b)
		if e1 >= e2 {
			mpc._pop()
		} else {
			break
		}
	}
	if mpc.size == 0 || mpc.as[mpc.size-1]+mpc.bs[mpc.size-1] < a+b {
		mpc._insert(a, b)
	}
}

func (mpc *MinPairContainer) Update(a int) {
	if mpc.size == 0 {
		return
	}
	mpc.as[mpc.size-1] = min(mpc.as[mpc.size-1], a)
	for mpc.size > 1 {
		mpc._pushDown()
		e := mpc._event(mpc.as[mpc.size-2], mpc.bs[mpc.size-2], mpc.bs[mpc.size-1])
		if e >= mpc.as[mpc.size-1] {
			mpc._pop()
		} else {
			break
		}
	}
}

func (mpc *MinPairContainer) Query() int {
	return mpc.as[mpc.size-1] + mpc.bs[mpc.size-1]
}

func (mpc *MinPairContainer) Empty() bool {
	return mpc.size == 0
}

func (mpc *MinPairContainer) Clear() {
	mpc.size = 0
}

func (mpc *MinPairContainer) _event(a0, b0, b1 int) int {
	if b0 <= b1 {
		return -INF
	}
	return a0 + b0 - b1
}

func (mpc *MinPairContainer) _insert(a, b int) {
	mpc.as[mpc.size] = a
	mpc.bs[mpc.size] = b
	mpc.size++
}

func (mpc *MinPairContainer) _pop() {
	mpc.size--
}

func (mpc *MinPairContainer) _pushDown() {
	mpc.as[mpc.size-2] = min(mpc.as[mpc.size-2], mpc.as[mpc.size-1])
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
