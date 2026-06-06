// https://www.luogu.com.cn/problem/CF551E
// https://codeforces.com/problemset/problem/551/E
//
// 输入 n(1≤n≤5e5) q(1≤q≤5e4) 和长为 n 的数组 a(1≤a[i]≤1e9)。
//
// 然后输入 q 个询问，格式如下：
// "1 L R x"：把下标在闭区间 [L,R] 中的 a[i] 都增加 x(0≤x≤1e9)。注：a 的下标从 1 开始。
// "2 y"：设 i 和 j 为元素 y(1≤y≤1e9) 在 a 中的最左下标和最右下标，输出 j-i。如果 a 中没有 y，输出 -1。
//
// 注：本题时间限制为 10s。
//
// 这种涉及精确元素位置+带修的题目，很难用线段树维护。由于本题给了 10s 时限，考虑分块。
// 把 a 分成若干段（块）。
//
// 每一块需要维护哪些信息？
// 有区间加，需要一个 lazy tag，表示这个块整体增加的值。（类似线段树）
// 需要知道区间中的每个数的最左位置和最右位置，用哈希表维护。
//
// 区间加：
// 块被 [L,R] 完全覆盖时，只把 lazy tag 增加 x。
// !块被 [L,R] 部分覆盖时（最多有两个这样的块），暴力更新区间中的每个数（注意应用 lazy tag），然后重置 lazy tag = 0，重新计算块中每个数的最左位置和最右位置。
//
// 查询：
// 正着遍历块，找最左位置。找到就退出循环。
// 倒着遍历块，找最右位置。找到就退出循环。
// 注：Go1.22 推荐块大小为 sqrt(n/5)，比 sqrt(n) 快一倍。
//
//
// !另一种方法：珂朵莉树维护区间 + 区间定期暴力重构(没有区间推平操作，只能隔一段时间暴力推平一次).

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	Solve2()
	// Solve1()
}

// 分块，每个块维护一个哈希表，记录每个元素的最左位置和最右位置.
// 修改时，整块直接修改 lazy tag，零散块暴力重构.
// 查询时，遍历每个块，找到最左位置和最右位置，然后返回 j-i.
//
// 5281 ms
func Solve1() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	arr := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	blockSize := int32(math.Ceil(math.Sqrt(float64(n) / 5)))

	blockCount := 1 + (n / blockSize)
	blockStart := make([]int32, blockCount)
	blockEnd := make([]int32, blockCount)
	belong := make([]int32, n)
	for i := int32(0); i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := int32(0); i < n; i++ {
		belong[i] = i / blockSize
	}

	blockLazy := make([]int, blockCount)
	blockFirst := make([]map[int]int32, blockCount)
	blockLast := make([]map[int]int32, blockCount)

	rebuild := func(bid int32) {
		start, end := blockStart[bid], blockEnd[bid]
		for i := start; i < end; i++ {
			arr[i] += blockLazy[bid]
		}
		blockLazy[bid] = 0
		fmp := make(map[int]int32)
		for i := end - 1; i >= start; i-- {
			fmp[arr[i]] = i
		}
		blockFirst[bid] = fmp
		lmp := make(map[int]int32)
		for i := start; i < end; i++ {
			lmp[arr[i]] = i
		}
		blockLast[bid] = lmp
	}

	update := func(l, r int32, x int) {
		bid1, bid2 := belong[l], belong[r-1]
		if bid1 == bid2 {
			for i := l; i < r; i++ {
				arr[i] += x
			}
			rebuild(bid1)
		} else {
			for i := l; i < blockEnd[bid1]; i++ {
				arr[i] += x
			}
			rebuild(bid1)
			for bid := bid1 + 1; bid < bid2; bid++ {
				blockLazy[bid] += x
			}
			for i := blockStart[bid2]; i < r; i++ {
				arr[i] += x
			}
			rebuild(bid2)
		}
	}

	query := func(y int) int32 {
		posL := n
		for bid := int32(0); bid < blockCount; bid++ {
			target := y - blockLazy[bid]
			if first, exists := blockFirst[bid][target]; exists {
				posL = first
				break
			}
		}
		if posL == n {
			return -1
		}
		for bid := blockCount - 1; bid >= 0; bid-- {
			target := y - blockLazy[bid]
			if last, exists := blockLast[bid][target]; exists {
				return last - posL
			}
		}
		return -1
	}

	for bid := int32(0); bid < blockCount; bid++ {
		rebuild(bid)
	}

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var l, r int32
			var x int
			fmt.Fscan(in, &l, &r, &x)
			l--
			update(l, r, x)
		} else {
			var y int
			fmt.Fscan(in, &y)
			fmt.Fprintln(out, query(y))
		}
	}
}

// 珂朵莉重构
// 把增量相同的段看成一个区间，区间用珂朵莉树维护.
// 刚开始增量都是0，只有一个区间.
// 修改一次，最多增加两个区间.
// 阈值B=800，区间个数到达800的时候就重构一次,最多重构q/B次.
// 查找时，只需要遍历每个段在邻接表二分查找即可.
//
// 3750ms.
func Solve2() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	arr := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	odt := NewODT32(n, -1)
	odt.Set(0, n, 0)
	valueToIndexes := make(map[int][]int32)

	b := int32(5000)

	rebuild := func() {
		valueToIndexes = make(map[int][]int32)
		odt.EnumerateAll(func(start, end int32, value *int) {
			add := *value
			for i := start; i < end; i++ {
				arr[i] += add
				valueToIndexes[arr[i]] = append(valueToIndexes[arr[i]], i)
			}
		})
		odt.Set(0, n, 0)
	}

	update := func(l, r int32, x int) {
		if odt.Len > b {
			rebuild()
		}

		odt.Split(l)
		odt.Split(r)
		odt.EnumerateRange(l, r, func(start, end int32, value *int) {
			*value += x
		}, false)
	}

	query := func(y int) int32 {
		minPos, maxPos := int32(n), int32(-1)

		odt.EnumerateAll(func(start, end int32, addPtr *int) {
			add := *addPtr
			target := y - add
			positions, exists := valueToIndexes[target]
			if !exists {
				return
			}

			// 二分查找第一个 >= start 的位置
			left := 0
			right := len(positions)
			for left < right {
				mid := (left + right) / 2
				if positions[mid] >= start {
					right = mid
				} else {
					left = mid + 1
				}
			}

			if left < len(positions) && positions[left] < end {
				if positions[left] < minPos {
					minPos = positions[left]
				}
			} else {
				return
			}

			// 二分查找最后一个 < end 的位置
			left2 := left
			right2 := len(positions)
			for left2 < right2 {
				mid := (left2 + right2) / 2
				if positions[mid] < end {
					left2 = mid + 1
				} else {
					right2 = mid
				}
			}
			if left2 > 0 && positions[left2-1] >= start {
				maxPos = positions[left2-1]
			}
		})

		if maxPos == -1 {
			return -1
		}
		return maxPos - minPos
	}

	rebuild()

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var l, r int32
			var x int
			fmt.Fscan(in, &l, &r, &x)
			l--
			update(l, r, x)
		} else {
			var y int
			fmt.Fscan(in, &y)
			fmt.Fprintln(out, query(y))
		}
	}
}

type ODT32[V comparable] struct {
	Len        int32 // 区间数
	Count      int32 // 区间元素个数之和
	llim, rlim int32
	noneValue  V
	data       []V
	ss         *FastSet32
}

// 指定区间长度 n 和哨兵 noneValue 建立一个 ODT.
//
//	区间为[0,n).
func NewODT32[V comparable](n int32, noneValue V) *ODT32[V] {
	res := &ODT32[V]{}
	dat := make([]V, n)
	for i := int32(0); i < n; i++ {
		dat[i] = noneValue
	}
	ss := NewFastSet32(n)
	ss.Insert(0)

	res.rlim = n
	res.noneValue = noneValue
	res.data = dat
	res.ss = ss
	return res
}

// 返回包含 x 的区间的信息.
func (odt *ODT32[V]) Get(x int32, erase bool) (start, end int32, value V) {
	start, end = odt.ss.Prev(x), odt.ss.Next(x+1)
	value = odt.data[start]
	if erase && value != odt.noneValue {
		odt.Len--
		odt.Count -= end - start
		odt.data[start] = odt.noneValue
		odt.Merge(start)
		odt.Merge(end)
	}
	return
}

func (odt *ODT32[V]) Set(start, end int32, value V) {
	odt.EnumerateRange(start, end, func(l, r int32, x *V) {}, true)
	odt.ss.Insert(start)
	odt.data[start] = value
	if value != odt.noneValue {
		odt.Len++
		odt.Count += end - start
	}
	odt.Merge(start)
	odt.Merge(end)
}

func (odt *ODT32[V]) EnumerateAll(f func(start, end int32, value *V)) {
	odt.EnumerateRange(0, odt.rlim, f, false)
}

// 遍历范围 [L, R) 内的所有数据.
func (odt *ODT32[V]) EnumerateRange(start, end int32, f func(start, end int32, value *V), erase bool) {
	if !(odt.llim <= start && start <= end && end <= odt.rlim) {
		panic(fmt.Sprintf("invalid range [%d, %d)", start, end))
	}

	NONE := odt.noneValue
	if !erase {
		l := odt.ss.Prev(start)
		for l < end {
			r := odt.ss.Next(l + 1)
			f(max32(l, start), min32(r, end), &odt.data[l])
			l = r
		}
		return
	}

	// 分割区间
	p := odt.ss.Prev(start)
	if p < start {
		odt.ss.Insert(start)
		v := odt.data[p]
		odt.data[start] = v
		if v != NONE {
			odt.Len++
		}
	}
	p = odt.ss.Next(end)
	if end < p {
		v := odt.data[odt.ss.Prev(end)]
		odt.data[end] = v
		odt.ss.Insert(end)
		if v != NONE {
			odt.Len++
		}
	}
	p = start
	for p < end {
		q := odt.ss.Next(p + 1)
		x := odt.data[p]
		f(p, q, &x)
		if x != NONE {
			odt.Len--
			odt.Count -= q - p
		}
		odt.ss.Erase(p)
		p = q
	}
	odt.ss.Insert(start)
	odt.data[start] = NONE
}

// 在位置 pos 处分割区间.
// 如果 pos 已经是区间的起始位置，则不进行分割.
func (odt *ODT32[V]) Split(pos int32) {
	if pos >= odt.rlim || pos <= odt.llim || odt.ss.Has(pos) {
		return
	}
	start := odt.ss.Prev(pos)
	odt.ss.Insert(pos)
	odt.data[pos] = odt.data[start]
	if odt.data[pos] != odt.noneValue {
		odt.Len++
	}
}

func (odt *ODT32[V]) Merge(p int32) {
	if p <= 0 || odt.rlim <= p {
		return
	}
	q := odt.ss.Prev(p - 1)
	if dataP, dataQ := odt.data[p], odt.data[q]; dataP == dataQ {
		if dataP != odt.noneValue {
			odt.Len--
		}
		odt.ss.Erase(p)
	}
}

func (odt *ODT32[V]) String() string {
	sb := []string{}
	odt.EnumerateAll(func(start, end int32, value *V) {
		var v interface{} = value
		if *value == odt.noneValue {
			v = "nil"
		}
		sb = append(sb, fmt.Sprintf("[%d,%d):%v", start, end, v))
	})
	return fmt.Sprintf("ODT{%v}", strings.Join(sb, ", "))
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

type FastSet32 struct {
	n, lg int32
	seg   [][]uint64
	size  int32
}

func NewFastSet32(n int32) *FastSet32 {
	res := &FastSet32{n: n}
	seg := [][]uint64{}
	n_ := n
	for {
		seg = append(seg, make([]uint64, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = int32(len(seg))
	return res
}

func NewFastSet32From(n int32, f func(i int32) bool) *FastSet32 {
	res := NewFastSet32(n)
	for i := int32(0); i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := int32(0); h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet32) Has(i int32) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet32) Insert(i int32) bool {
	if fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet32) Erase(i int32) bool {
	if !fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet32) Next(i int32) int32 {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == int32(len(cache)) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet32) Prev(i int32) int32 {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := int32(0); h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet32) Enumerate(start, end int32, f func(i int32)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet32) String() string {
	res := []string{}
	for i := int32(0); i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(int(i)))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet32) Size() int32 {
	return fs.size
}

func (*FastSet32) bsr(x uint64) int32 {
	return 63 - int32(bits.LeadingZeros64(x))
}

func (*FastSet32) bsf(x uint64) int32 {
	return int32(bits.TrailingZeros64(x))
}
