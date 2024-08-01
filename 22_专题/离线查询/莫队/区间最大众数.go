package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

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

	const MAX int32 = 3e5 + 5

	n, q := int32(NextInt()), int32(NextInt())
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		nums[i] = int32(NextInt())
	}
	lefts, rights, ceils := make([]int32, q), make([]int32, q), make([]int32, q)
	mo := NewMoRollback32(n, q)
	for i := int32(0); i < q; i++ {
		lefts[i], rights[i], ceils[i] = int32(NextInt()), int32(NextInt()), int32(NextInt())
		lefts[i]--
		mo.AddQuery(lefts[i], rights[i])
	}

	res := make([]int32, q)
	maxCount, maxKey := int32(0), int32(0)
	maxCountSnap, maxKeySnap := int32(0), int32(0)
	counter := [MAX + 1]int32{} // 可撤销counter
	counterHistory := make([]int32, 0, n)
	counterTime := int32(0)
	seg := NewLazySegTree32Rollbackable(MAX, func(i int32) E { return 0 }) // 可撤销值域线段树
	segTime0, segTime1 := int32(0), int32(0)

	add := func(index int32) {
		x := nums[index]
		counterHistory = append(counterHistory, x)
		counter[x]++
		if counter[x] > maxCount || counter[x] == maxCount && x > maxKey {
			maxCount, maxKey = counter[x], x
		}
		seg.Update(x, MAX+1, int(counter[x])<<20|int(x))
	}

	reset := func() {
		maxCount, maxKey = 0, 0
		for _, v := range counterHistory {
			counter[v] = 0
		}
		counterHistory = counterHistory[:0]
		seg.Clear()
	}

	snapshot := func() {
		maxCountSnap, maxKeySnap = maxCount, maxKey
		counterTime = int32(len(counterHistory))
		segTime0, segTime1 = seg.GetTime()
	}

	rollback := func() {
		maxCount, maxKey = maxCountSnap, maxKeySnap
		for int32(len(counterHistory)) > counterTime {
			x := counterHistory[len(counterHistory)-1]
			counterHistory = counterHistory[:len(counterHistory)-1]
			counter[x]--
		}
		seg.Rollback(segTime0, segTime1)
	}

	query := func(qi int32) {
		res[qi] = int32(seg.Query(0, ceils[qi]+1) & mask20)
	}

	mo.Run(add, add, reset, snapshot, rollback, query, -1)
	for _, v := range res {
		fmt.Fprintln(out, v)
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

// RangeChmaxRangeMax

type E = int  // (count, key)
type Id = int // (count, key)
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
func (*LazySegTree32Rollbackable) e() E   { return 0 }
func (*LazySegTree32Rollbackable) id() Id { return 0 }
func (*LazySegTree32Rollbackable) op(left, right E) E {
	return maxMask(left, right)
}
func (*LazySegTree32Rollbackable) mapping(f Id, g E) E {
	return maxMask(f, g)
}
func (*LazySegTree32Rollbackable) composition(f, g Id) Id {
	return maxMask(f, g)
}

// !template
type LazySegTree32Rollbackable struct {
	n    int32
	size int32
	log  int32
	data *rollbackArraySpecified
	lazy *rollbackArraySpecified
}

func NewLazySegTree32Rollbackable(n int32, f func(int32) E) *LazySegTree32Rollbackable {
	tree := &LazySegTree32Rollbackable{}
	tree.n = n
	tree.log = 1
	for 1<<tree.log < n {
		tree.log++
	}
	tree.size = 1 << tree.log
	data := make([]E, tree.size<<1)
	for i := range data {
		data[i] = tree.e()
	}
	for i := int32(0); i < n; i++ {
		data[tree.size+i] = f(i)
	}
	// pushUp
	for i := tree.size - 1; i >= 1; i-- {
		data[i] = tree.op(data[i<<1], data[i<<1|1])
	}
	tree.data = newRollbackArraySpecifiedFrom(data)
	tree.lazy = newRollbackArraySpecified(tree.size, func(int32) Id { return tree.id() })
	return tree
}

func NewLazySegTreeRollbackableFrom(leaves []E) *LazySegTree32Rollbackable {
	return NewLazySegTree32Rollbackable(int32(len(leaves)), func(i int32) E { return leaves[i] })
}

func (tree *LazySegTree32Rollbackable) GetTime() (dataTime, lazyTime int32) {
	return tree.data.GetTime(), tree.lazy.GetTime()
}

func (tree *LazySegTree32Rollbackable) Rollback(dataTime, lazyTime int32) {
	tree.data.Rollback(dataTime)
	tree.lazy.Rollback(lazyTime)
}

func (tree *LazySegTree32Rollbackable) Clear() {
	tree.Rollback(0, 0)
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32Rollbackable) Query(left, right int32) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data.Get(left))
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data.Get(right), smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}
func (tree *LazySegTree32Rollbackable) QueryAll() E {
	return tree.data.Get(1)
}
func (tree *LazySegTree32Rollbackable) GetAll() []E {
	tmp := tree.data.GetAll()
	for k := int32(1); k < tree.size; k++ {
		tree.pushDown(k)
	}
	return tmp[tree.size : tree.size+tree.n]
	// return append(res[:0:0], res...)
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32Rollbackable) Update(left, right int32, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := int32(1); i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

func (tree *LazySegTree32Rollbackable) pushUp(root int32) {
	tree.data.Set(root, tree.op(tree.data.Get(root<<1), tree.data.Get(root<<1|1)))
}

func (tree *LazySegTree32Rollbackable) pushDown(root int32) {
	if tmp := tree.lazy.Get(root); tmp != tree.id() {
		tree.propagate(root<<1, tmp)
		tree.propagate(root<<1|1, tmp)
		tree.lazy.Set(root, tree.id())
	}
}
func (tree *LazySegTree32Rollbackable) propagate(root int32, f Id) {
	tree.data.Set(root, tree.mapping(f, tree.data.Get(root)))
	if root < tree.size {
		tree.lazy.Set(root, tree.composition(f, tree.lazy.Get(root)))
	}
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
	return &rollbackArraySpecified{
		n:    n,
		data: data,
	}
}

func newRollbackArraySpecifiedFrom(data []int) *rollbackArraySpecified {
	return &rollbackArraySpecified{n: int32(len(data)), data: data}
}

func (r *rollbackArraySpecified) GetTime() int32 {
	return int32(len(r.history))
}

func (r *rollbackArraySpecified) Rollback(time int32) {
	for int32(len(r.history)) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair&mask20] = pair >> 20
	}
}

func (r *rollbackArraySpecified) Get(index int32) int {
	return r.data[index]
}

func (r *rollbackArraySpecified) Set(index int32, value int) {
	r.history = append(r.history, r.data[index]<<20|int(index))
	r.data[index] = value
}

func (r *rollbackArraySpecified) GetAll() []int {
	return append(r.data[:0:0], r.data...)
}

func (r *rollbackArraySpecified) GetAllMut() []int {
	return r.data
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

func Discretize(nums []int32) (newNums []int32, origin []int32) {
	newNums = make([]int32, len(nums))
	origin = make([]int32, 0, len(newNums))
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
