package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	StaticRangeModeQuery()
}

// 区间众数查询，存在多个众数时，返回任意一个.
// https://judge.yosupo.jp/problem/static_range_mode_query
func StaticRangeModeQuery() {
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
	_ = NextInt

	n, q := int32(NextInt()), int32(NextInt())

	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		nums[i] = int32(NextInt())
	}
	D := NewDictionary[int32]()
	for i := int32(0); i < n; i++ {
		nums[i] = D.Id(nums[i])
	}

	M := NewMoRollback32(n, q)
	for i := int32(0); i < q; i++ {
		l, r := int32(NextInt()), int32(NextInt())
		M.AddQuery(l, r)
	}
	res := make([][2]int32, q) // (众数, 众数出现次数)

	counter := make([]int32, D.Size())
	history := make([]int32, 0, n)

	maxCount, maxKey := int32(0), int32(0)
	snapState := int32(0)
	snapCount, snapKey := int32(0), int32(0)

	add := func(index int32) {
		x := nums[index]
		history = append(history, x)
		counter[x]++
		if counter[x] > maxCount {
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
		snapState = int32(len(history))
		snapCount, snapKey = maxCount, maxKey
	}

	rollback := func() {
		for int32(len(history)) > snapState {
			x := history[len(history)-1]
			history = history[:len(history)-1]
			counter[x]--
		}
		maxCount, maxKey = snapCount, snapKey
	}

	query := func(qi int32) {
		res[qi] = [2]int32{D.Value(maxKey), maxCount}
	}

	M.Run(add, add, reset, snapshot, rollback, query, -1)
	for _, v := range res {
		fmt.Fprintln(out, v[0], v[1])
	}
}

type MoRollback32 struct {
	left, right []int32
}

func NewMoRollback32(n, q int32) *MoRollback32 {
	return &MoRollback32{left: make([]int32, 0, q), right: make([]int32, 0, q)}
}

func (mo *MoRollback32) AddQuery(start, end int32) {
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
func (mo *MoRollback32) Run(
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
		sort.Slice(order, func(i, j int) bool {
			return mo.right[order[i]] < mo.right[order[j]]
		})
		lMax := int32(0)
		for _, qid := range order {
			lMax = max32(lMax, mo.left[qid])
		}
		reset()
		l, r := lMax, lMax
		for _, qi := range order {
			L, R := mo.left[qi], mo.right[qi]
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

type Dictionary[V comparable] struct {
	_idToValue []V
	_valueToId map[V]int32
}

// A dictionary that maps values to unique ids.
func NewDictionary[V comparable]() *Dictionary[V] {
	return &Dictionary[V]{
		_valueToId: map[V]int32{},
	}
}
func (d *Dictionary[V]) Id(value V) int32 {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := int32(len(d._idToValue))
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary[V]) Value(id int32) V {
	return d._idToValue[id]
}
func (d *Dictionary[V]) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary[V]) Size() int32 {
	return int32(len(d._idToValue))
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
