// FastHashRange-区间哈希值
// 支持单点仿射变换+查询区间哈希值
// !哈希函数为仿射变换
// TODO 支持仿射变换线段树获取区间哈希值???
// TODO 支持单点修改的线段树获取区间哈希值???

package main

func main() {

}

type E = [3]int

var MODS = [3]int{1000000007, 1000000009, 1000000021}
var BASES = [3]int{131, 13331, 13333331}

func add(hash1, hash2 E) E {
	return [3]int{(hash1[0] + hash2[0]) % MODS[0], (hash1[1] + hash2[1]) % MODS[1], (hash1[2] + hash2[2]) % MODS[2]}
}

func mul(hash E, k int) E {
	return [3]int{
		((hash[0]*k)%MODS[0] + MODS[0]) % MODS[0],
		((hash[1]*k)%MODS[1] + MODS[1]) % MODS[1],
		((hash[2]*k)%MODS[2] + MODS[2]) % MODS[2],
	}
}

func initHash(n int) []E {
	res := make([]E, n)
	res[0] = [3]int{1, 1, 1}
	for i := 1; i < n; i++ {
		res[i][0] = (res[i-1][0] * BASES[0]) % MODS[0]
		res[i][1] = (res[i-1][1] * BASES[1]) % MODS[1]
		res[i][2] = (res[i-1][2] * BASES[2]) % MODS[2]
	}
	return res
}

// 哈希函数需要支持仿射变换.
type FastHashRange struct {
	n   int
	seg *SegmentTree
}

// 初始时,所有位置都为0.
func NewFastHashRange(n int) *FastHashRange {
	leaves := initHash(n)
	return &FastHashRange{n: n, seg: NewSegmentTree(leaves)}
}

func (fhr *FastHashRange) Add(index int, value int) {
	fhr.seg.Update(index, value)
}

func (fhr *FastHashRange) Query(start, end int) int {}

func (*SegmentTree) e() E { return [3]int{0, 0, 0} }
func (*SegmentTree) op(a, b E) E { // 哈希函数需要满足仿射变换性质.
	return [3]int{(a[0] + b[0]) % MODS[0], (a[1] + b[1]) % MODS[1], (a[2] + b[2]) % MODS[2]}
}

type SegmentTree struct {
	n, size int
	seg     []E
	unit    E
}

func NewSegmentTree(leaves []E) *SegmentTree {
	res := &SegmentTree{}
	res.unit = res.e()
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}

func (st *SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.unit
	}
	return st.seg[index+st.size]
}

func (st *SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

func (st *SegmentTree) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.unit
	}
	leftRes, rightRes := st.unit, st.unit
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTree) QueryAll() E { return st.seg[1] }
