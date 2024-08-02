// 区间最大众数/区间不超过阈值的最大众数.
// https://www.luogu.com.cn/problem/T482369?contestId=183908
//
// 回滚莫队，然后加一个可撤销的"树状数组"维护值域前缀的众数，这个"树状数组"用值域分块实现

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"sort"
)

func init() {
	debug.SetGCPercent(-1)
}

// 区间下不大于限制数的最大众数
func main() {
	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}

	n, q := int32(NextInt()), int32(NextInt())
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		nums[i] = int32(NextInt())
	}
	lefts, rights, ceils := make([]int32, q), make([]int32, q), make([]int32, q)
	mo := newMoRollback32(n, q)
	for i := int32(0); i < q; i++ {
		lefts[i], rights[i], ceils[i] = int32(NextInt()), int32(NextInt()), int32(NextInt())
		lefts[i]--
		mo.AddQuery(lefts[i], rights[i])
	}

	const MAX int32 = 3e5

	res := make([]int32, q)
	counter := [MAX + 1]int32{} // 可撤销counter
	counterHistory := make([]int32, 0, n)
	counterTime := int32(0)
	bit := newBitLikeRollbackable(func() int { return 0 }, maxMask)
	bit.Build(MAX+1, func(i int32) int { return 0 }) // 可撤销值域树状数组
	bitTime0, bitTime1 := int32(0), int32(0)

	add := func(index int32) {
		x := nums[index]
		counterHistory = append(counterHistory, x)
		counter[x]++
		bit.Update(x, int(counter[x])<<20|int(x))
	}

	reset := func() {
		for _, v := range counterHistory {
			counter[v] = 0
		}
		counterHistory = counterHistory[:0]
		bit.Reset()
	}

	snapshot := func() {
		counterTime = int32(len(counterHistory))
		bitTime0, bitTime1 = bit.GetTime()
	}

	rollback := func() {
		for i := int32(len(counterHistory)) - 1; i >= counterTime; i-- {
			counter[counterHistory[i]]--
		}
		counterHistory = counterHistory[:counterTime]
		bit.Rollback(bitTime0, bitTime1)
	}

	query := func(qi int32) {
		state := bit.QueryRange(0, ceils[qi]+1)
		if state>>20 == 0 {
			res[qi] = 0
		} else {
			res[qi] = int32(state & mask20)
		}
	}

	mo.Run(add, add, reset, snapshot, rollback, query, -1)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type moRollback32 struct {
	left, right []int32
}

func newMoRollback32(n, q int32) *moRollback32 {
	return &moRollback32{left: make([]int32, 0, q), right: make([]int32, 0, q)}
}

func (mo *moRollback32) AddQuery(start, end int32) {
	mo.left = append(mo.left, start)
	mo.right = append(mo.right, end)
}

// addLeft : 将index位置的元素加入到区间左端.
// addRight: 将index位置的元素加入到区间右端.
// reset: 重置区间.
// snapShot: 快照当前状态.
// rollback: 回滚到快照状态.
// query: 查询当前区间.
// blockSize: 分块大小.-1表示使用默认值.
func (mo *moRollback32) Run(
	addLeft func(i int32),
	addRight func(i int32),
	reset func(),
	snapShot func(),
	rollback func(),
	query func(qi int32),
	blockSize int32,
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
	naive := func(qi int32) {
		snapShot()
		for i := mo.left[qi]; i < mo.right[qi]; i++ {
			addRight(i)
		}
		query(qi)
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
		left, right := mo.left, mo.right
		sort.Slice(order, func(i, j int) bool {
			return right[order[i]] < right[order[j]]
		})
		lMax := int32(0)
		for _, qid := range order {
			lMax = max32(lMax, left[qid])
		}
		reset()
		l, r := lMax, lMax
		for _, qi := range order {
			L, R := left[qi], right[qi]
			for r < R {
				addRight(r)
				r++
			}
			snapShot()
			for L < l {
				l--
				addLeft(l)
			}
			query(qi)
			rollback()
			l = lMax
		}
	}
}

// a,b: (count, key)
func maxMask(a, b int) int {
	c1, c2 := a>>20, b>>20
	if c1 > c2 {
		return a
	}
	if c1 < c2 {
		return b
	}
	k1, k2 := a&mask20, b&mask20
	if k1 > k2 {
		return a
	}
	return b
}

type bitLikeRollbackable struct {
	_n          int32
	_belong     []int32
	_blockStart []int32
	_blockEnd   []int32
	_nums       *rollbackArraySpecified
	_blockSum   *rollbackArraySpecified
	e           func() int
	op          func(a, b int) int
}

func newBitLikeRollbackable(e func() int, op func(a, b int) int) *bitLikeRollbackable {
	return &bitLikeRollbackable{e: e, op: op}
}

func (b *bitLikeRollbackable) GetTime() (time0, time1 int32) {
	return b._nums.GetTime(), b._blockSum.GetTime()
}

func (b *bitLikeRollbackable) Rollback(time0, time1 int32) {
	b._nums.Rollback(time0)
	b._blockSum.Rollback(time1)
}

func (b *bitLikeRollbackable) Reset() {
	b._nums.Rollback(0)
	b._blockSum.Rollback(0)
}

func (b *bitLikeRollbackable) Build(n int32, f func(i int32) int) {
	blockSize := int32(math.Sqrt(float64(n)) + 1)
	blockCount := 1 + (n / blockSize)
	belong := make([]int32, n)
	for i := int32(0); i < n; i++ {
		belong[i] = i / blockSize
	}
	blockStart := make([]int32, blockCount)
	blockEnd := make([]int32, blockCount)
	for i := int32(0); i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := blockStart[i] + blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	nums := make([]int, n)
	blockSum := make([]int, blockCount)
	for i := int32(0); i < n; i++ {
		nums[i] = f(i)
		bid := belong[i]
		blockSum[bid] = b.op(blockSum[bid], f(i))
	}
	b._n = n
	b._belong = belong
	b._blockStart = blockStart
	b._blockEnd = blockEnd
	b._nums = newRollbackArraySpecifiedFrom(nums)
	b._blockSum = newRollbackArraySpecifiedFrom(blockSum)
}

func (b *bitLikeRollbackable) Update(index int32, delta int) {
	b._nums.Set(index, b.op(b._nums.Get(index), delta))
	bid := b._belong[index]
	b._blockSum.Set(bid, b.op(b._blockSum.Get(bid), delta))
}

func (b *bitLikeRollbackable) QueryRange(start, end int32) int {
	if start < 0 {
		start = 0
	}
	if end > b._n {
		end = b._n
	}
	if start >= end {
		return b.e()
	}
	res := b.e()
	bid1 := b._belong[start]
	bid2 := b._belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			res = b.op(res, b._nums.Get(i))
		}
		return res
	}
	for i := start; i < b._blockEnd[bid1]; i++ {
		res = b.op(res, b._nums.Get(i))
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res = b.op(res, b._blockSum.Get(bid))
	}
	for i := b._blockStart[bid2]; i < end; i++ {
		res = b.op(res, b._nums.Get(i))
	}
	return res
}

const mask20 int = 1<<20 - 1

type rollbackArraySpecified struct {
	n       int32
	data    []int
	history []int // (value, index)
}

func newRollbackArraySpecified(n int32, f func(index int32) int) *rollbackArraySpecified {
	data := make([]int, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &rollbackArraySpecified{n: n, data: data, history: make([]int, 0, n)}
}

func newRollbackArraySpecifiedFrom(data []int) *rollbackArraySpecified {
	return &rollbackArraySpecified{n: int32(len(data)), data: data, history: make([]int, 0, len(data))}
}

func (r *rollbackArraySpecified) GetTime() int32 {
	return int32(len(r.history))
}

func (r *rollbackArraySpecified) Rollback(time int32) {
	for i := int32(len(r.history)) - 1; i >= time; i-- { // 注意反向
		pair := r.history[i]
		r.data[pair&mask20] = pair >> 20
	}
	r.history = r.history[:time]
}

func (r *rollbackArraySpecified) Get(index int32) int {
	return r.data[index]
}

func (r *rollbackArraySpecified) Set(index int32, value int) bool {
	if r.data[index] == value {
		return false
	}
	r.history = append(r.history, r.data[index]<<20|int(index))
	r.data[index] = value
	return true
}

func (r *rollbackArraySpecified) Len() int32 {
	return r.n
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
