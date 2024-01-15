// https://atcoder.jp/contests/abc324/editorial/7399
package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

// 给定一个数组，数组元素为1-n的排列
// 有两种操作：
// 1.把 A[version]中下标大于等于 x 的元素分裂成一个新的数组 Ai(A[version]中保留x个)。
// 2.把 A[version]中值大于 x 的元素分裂成一个新的数组 Ai。
// 这两种操作都不会改变元素相对顺序。
// 输出每次分裂出的数组大小。
//
// SortedList + Deque 维护.
// 启发式分裂：每次分裂出较小的那一半
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
	// 	int main() {
	//     using namespace std;
	//     unsigned N;
	//     cin >> N;
	//     vector<sequence> sequences(1);
	//     for (unsigned a; const auto i : views::iota(0U, N)) {
	//         cin >> a;
	//         sequences.front().emplace(i, a); // i 番目の要素 A[i] を追加する
	//     }
	//     unsigned Q;
	//     cin >> Q;
	//     sequences.resize(1 + Q);
	//     for (unsigned t, s, x; const auto i : views::iota(1U, 1 + Q)) {
	//         cin >> t >> s >> x;
	//         if (t == 1) { // x 個目より後を移動
	//             const auto M{size(sequences[s])};
	//             if (2 * x < M) { // 移動しない要素のほうが少ないとき
	//                 swap(sequences[s], sequences[i]); // 入れ替えて
	//                 for (const auto _ : views::iota(0U, x)) // x 個目以前を
	//                     sequences[s].emplace(sequences[i].pop_front()); // 前から順に移動
	//             } else if (x < M)
	//                 for (const auto _ : views::iota(x, M)) // x 個目より後を
	//                     sequences[i].emplace(sequences[s].pop_back()); // 後ろから順に移動
	//         } else { // x より大きな要素を移動
	//             if (x < sequences[s].mid()) { // x が中央値より小さいとき
	//                 swap(sequences[s], sequences[i]); // 入れ替えて
	//                 while (!empty(sequences[i]) && sequences[i].min() <= x) // x 以下の要素が残っている限り
	//                     sequences[s].emplace(sequences[i].pop_min()); // 小さいほうから順に移動
	//             } else
	//                 while (!empty(sequences[s]) && sequences[s].max() > x) // x より大きい要素が残っている限り
	//                     sequences[i].emplace(sequences[s].pop_max()); // 大きいほうから順に移動
	//         }
	//         cout << size(sequences[i]) << endl; // i 番の数列の大きさを出力
	//     }
	//     return 0;
	// }

	var q int
	fmt.Fscan(in, &q)

	gits := make([]*SortedList, q+1)
	gits[0] = NewSortedList(func(a, b S) bool { return a < b }, nums...)
	for i := 1; i < len(gits); i++ {
		gits[i] = NewSortedList(func(a, b S) bool { return a < b })
	}

	for i := 0; i < q; i++ {
		var kind, version, x int
		fmt.Fscan(in, &kind, &version, &x)
		if kind == 1 { // 将 A[version] 中下标大于等于 x 的元素分裂成一个新的数组 Ai
			len1, len2 := x, gits[version].Len()-x
			if len1 < len2 { // 前面少，拆前面
			} else { // 后面少，拆后面
			}
		} else { // 将 A[version] 中值大于 x 的元素分裂成一个新的数组 Ai
		}
	}
}

// 可删除元素、获取第k小值的双端队列.
type SortedDeque struct {
	sl *SortedList
	dq *Deque
}

func NewSortedDeque(less func(a, b S) bool, elements ...S) *SortedDeque {
	elements = append(elements[:0:0], elements...)
	res := &SortedDeque{sl: NewSortedList(less, elements...), dq: NewDeque(len(elements))}
	for _, v := range elements {
		res.dq.Append(v)
	}
	return res
}

func (sd *SortedDeque) Append(value S) {
	sd.sl.Add(value)
	sd.dq.Append(value)
}

func (sd *SortedDeque) AppendLeft(value S) {
	sd.sl.Add(value)
	sd.dq.AppendLeft(value)
}

func (sd *SortedDeque) Pop() S {
	value := sd.dq.Pop()
	sd.sl.Discard(value)
	return value
}

func (sd *SortedDeque) PopLeft() S {
	value := sd.dq.PopLeft()
	sd.sl.Discard(value)
	return value
}

func (sd *SortedDeque) Remove(value S) {}

func (sd *SortedDeque) Min() S {
	return sd.sl.Min()
}

func (sd *SortedDeque) Max() S {
	return sd.sl.Max()
}

func (sd *SortedDeque) Kth(k int) S {
	return sd.sl.At(k)
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 200

type S = int

// 使用分块+树状数组维护的有序序列.
type SortedList struct {
	less              func(a, b S) bool
	size              int
	blocks            [][]S
	mins              []S
	tree              []int
	shouldRebuildTree bool
}

func NewSortedList(less func(a, b S) bool, elements ...S) *SortedList {
	elements = append(elements[:0:0], elements...)
	res := &SortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]S{}
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
		blocks = append(blocks, elements[start:end:end]) // !各个块互不影响, max参数也需要指定为end
	}
	mins := make([]S, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	return res
}

func (sl *SortedList) Add(value S) *SortedList {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return sl
	}

	pos, index := sl._locRight(value)

	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]S{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]

	// n -> load + (n - load)
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]S{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]S{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *SortedList) Has(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < len(sl.blocks[pos]) && sl.blocks[pos][index] == value
}

func (sl *SortedList) Discard(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locRight(value)
	if index > 0 && sl.blocks[pos][index-1] == value {
		sl._delete(pos, index-1)
		return true
	}
	return false
}

func (sl *SortedList) Pop(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}

func (sl *SortedList) At(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SortedList) Erase(start, end int) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedList) Lower(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList) Higher(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *SortedList) Floor(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList) Ceiling(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *SortedList) BisectLeft(value S) int {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *SortedList) BisectRight(value S) int {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *SortedList) Count(value S) int {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *SortedList) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SortedList) ForEach(f func(value S, index int), reverse bool) {
	if !reverse {
		count := 0
		for i := 0; i < len(sl.blocks); i++ {
			block := sl.blocks[i]
			for j := 0; j < len(block); j++ {
				f(block[j], count)
				count++
			}
		}
		return
	}

	count := 0
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		for j := len(block) - 1; j >= 0; j-- {
			f(block[j], count)
			count++
		}
	}
}

func (sl *SortedList) Enumerate(start, end int, f func(value S), erase bool) {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return
	}

	pos, startIndex := sl._findKth(start)
	count := end - start
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endIndex := min(len(block), startIndex+count)
		if f != nil {
			for j := startIndex; j < endIndex; j++ {
				f(block[j])
			}
		}
		deleted := endIndex - startIndex

		if erase {
			if deleted == len(block) {
				// !delete block
				sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
				sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				for i := startIndex; i < endIndex; i++ {
					sl._updateTree(pos, -1)
				}
				block = append(block[:startIndex], block[endIndex:]...)
				sl.mins[pos] = block[0]
			}
			sl.size -= deleted
		}

		count -= deleted
		startIndex = 0
	}
}

func (sl *SortedList) Slice(start, end int) []S {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return nil
	}
	count := end - start
	res := make([]S, 0, count)
	pos, index := sl._findKth(start)
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endPos := min(len(block), index+count)
		curCount := endPos - index
		res = append(res, block[index:endPos]...)
		count -= curCount
		index = 0
	}
	return res
}

func (sl *SortedList) Range(min, max S) []S {
	if sl.less(max, min) {
		return nil
	}
	res := []S{}
	pos := sl._locBlock(min)
	for i := pos; i < len(sl.blocks); i++ {
		block := sl.blocks[i]
		for j := 0; j < len(block); j++ {
			x := block[j]
			if sl.less(max, x) {
				return res
			}
			if !sl.less(x, min) {
				res = append(res, x)
			}
		}
	}
	return res
}

func (sl *SortedList) Min() S {
	if sl.size == 0 {
		panic("Min() called on empty SortedList")
	}
	return sl.mins[0]
}

func (sl *SortedList) Max() S {
	if sl.size == 0 {
		panic("Max() called on empty SortedList")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *SortedList) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedList{")
	sl.ForEach(func(value S, index int) {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%v", value))
	}, false)
	sb.WriteByte('}')
	return sb.String()
}

func (sl *SortedList) Len() int {
	return sl.size
}

func (sl *SortedList) _delete(pos, index int) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], sl.blocks[pos][index+1:]...)
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
	sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
	sl.shouldRebuildTree = true
}

func (sl *SortedList) _locLeft(value S) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := -1
	right := len(sl.blocks) - 1
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	pos = right

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = len(cur)
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(cur[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SortedList) _locRight(value S) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := 0
	right := len(sl.blocks)
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, sl.mins[mid]) {
			right = mid
		} else {
			left = mid
		}
	}
	pos = left

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = len(cur)
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, cur[mid]) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SortedList) _locBlock(value S) int {
	left, right := -1, len(sl.blocks)-1
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	return right
}

func (sl *SortedList) _buildTree() {
	sl.tree = make([]int, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = len(sl.blocks[i])
	}
	tree := sl.tree
	for i := 0; i < len(tree); i++ {
		j := i | (i + 1)
		if j < len(tree) {
			tree[j] += tree[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *SortedList) _updateTree(index, delta int) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < len(tree); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SortedList) _queryTree(end int) int {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := 0
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *SortedList) _findKth(k int) (pos, index int) {
	if k < len(sl.blocks[0]) {
		return 0, k
	}
	last := len(sl.blocks) - 1
	lastLen := len(sl.blocks[last])
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	pos = -1
	bitLength := bits.Len32(uint32(len(tree)))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < len(tree) && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

type Deque struct{ l, r []S }

func NewDeque(cap int) *Deque { return &Deque{make([]S, 0, 1+cap/2), make([]S, 0, 1+cap/2)} }

func (q *Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q *Deque) Len() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v S) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v S) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v S) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v S) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q *Deque) Front() S {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q *Deque) Back() S {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q *Deque) At(i int) S {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
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
