// https://github.com/EndlessCheng/codeforces-go/blob/646deb927bbe089f60fc0e9f43d1729a97399e5f/copypasta/strings.go#L556
// https://visualgo.net/zh/suffixarray
// !常用分隔符 #(35) $(36) _(95) |(124)
// SA-IS 与 DC3 的效率对比 https://riteme.site/blog/2016-6-19/sais.html#5
// 注：Go1.13 开始使用 SA-IS 算法
//
// - 支持sa/rank/lcp
// - 比较任意两个子串的字典序
// - 求出任意两个子串的最长公共前缀(lcp)

//  sa : 排第几的后缀是谁.
//  rank : 每个后缀排第几.
//  lcp : 排名相邻的两个后缀的最长公共前缀.
// 	lcp[0] = 0
// 	lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
//
//  "banana" -> sa: [5 3 1 0 4 2], rank: [3 2 5 1 4 0], lcp: [0 1 3 0 0 2]
//
//  !lcp(sa[i],sa[j]) = min(height[i+1..j])
//
// !api:
//  func NewSuffixArray(ords []int) *SuffixArray
//  func NewSuffixArrayWithString(s string) *SuffixArray
//  func (suf *SuffixArray) Lcp(a, b int, c, d int) int
//  func (suf *SuffixArray) CompareSubstr(a, b int, c, d int) int
//  func (suf *SuffixArray) LcpRange(left int, k int) (start, end int)
//	func (suf *SuffixArray) Count(start, end int) int
//  func GetSA(ords []int) (sa []int)
//  func UseSA(ords []int) (sa, rank, lcp []int)

package main

import (
	"fmt"
	"index/suffixarray"
	"reflect"
	"strings"
	"unsafe"

	"math/bits"
)

func main() {
	test()
}

func test() {
	n := int32(1000)
	ords := make([]int32, n)
	for i := int32(1); i < n; i++ {
		ords[i] = int32(i * i)
		ords[i] ^= ords[i-1]
	}

	S := NewSuffixArray(int32(len(ords)), func(i int32) int32 { return ords[i] })
	S2 := NewSuffixArray(int32(len(ords)), func(i int32) int32 { return ords[i] })
	LcpRange2 := func(left, k int32) (start, end int32) {
		curRank := S2.Rank[left]
		for i := curRank; i >= 0; i-- {
			sa := S2.Sa[i]
			if S2.Lcp(sa, n, left, n) >= k {
				start = i
			} else {
				break
			}
		}
		for i := curRank; i < n; i++ {
			sa := S2.Sa[i]
			if S2.Lcp(sa, n, left, n) >= k {
				end = i + 1
			} else {
				break
			}
		}
		if start == 0 && end == 0 {
			return -1, -1
		}
		return
	}

	count2 := func(ords []int32, start, end int32) int32 {

		target := ords[start:end]
		len_ := end - start
		if len_ == 0 {
			return 0
		}
		res := int32(0)

		for i := int32(0); i+len_ <= n; i++ {
			curOrds := ords[i : i+len_]
			allOk := true
			for j := int32(0); j < len_; j++ {
				if curOrds[j] != target[j] {
					allOk = false
					break
				}
			}
			if allOk {
				res++
			}
		}
		return res
	}

	// lcpRange
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < n; j++ {
			start1, end1 := S.LcpRange(i, j)
			start2, end2 := LcpRange2(i, j)
			if start1 != start2 || end1 != end2 {
				fmt.Println(i, j, start1, end1, start2, end2)
				panic("")
			}
		}
	}

	// count
	for i := int32(0); i < n; i++ {
		for j := i; j < n; j++ {
			c1 := S.Count(i, j)
			c2 := count2(ords, i, j)
			if c1 != c2 {
				fmt.Println(i, j, c1, c2)
				panic("")
			}
		}
	}

	fmt.Println("pass ")
}

type SuffixArray32 struct {
	Sa     []int32 // 排名第i的后缀是谁.
	Rank   []int32 // 后缀s[i:]的排名是多少.
	Height []int32 // 排名相邻的两个后缀的最长公共前缀.Height[0] = 0,Height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	Ords   []int32
	n      int32
	minSt  *LinearRMQ32 // 维护lcp的最小值
}

func NewSuffixArray(n int32, f func(i int32) int32) *SuffixArray32 {
	res := &SuffixArray32{n: n}
	sa, rank, lcp := SuffixArray32Simple(n, f)
	res.Sa, res.Rank, res.Height = sa, rank, lcp
	return res
}

// 求任意两个子串s[a,b)和s[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray32) Lcp(a, b int32, c, d int32) int32 {
	cand := suf._lcp(a, c)
	return min32(cand, min32(b-a, d-c))
}

// 比较任意两个子串s[a,b)和s[c,d)的字典序.
//
//	s[a,b) < s[c,d) 返回-1.
//	s[a,b) = s[c,d) 返回0.
//	s[a,b) > s[c,d) 返回1.
func (suf *SuffixArray32) CompareSubstr(a, b int32, c, d int32) int32 {
	len1, len2 := b-a, d-c
	lcp := suf.Lcp(a, b, c, d)
	if len1 == len2 && lcp >= len1 {
		return 0
	}
	if lcp >= len1 || lcp >= len2 { // 一个是另一个的前缀
		if len1 < len2 {
			return -1
		}
		return 1
	}
	if suf.Rank[a] < suf.Rank[c] {
		return -1
	}
	return 1
}

// 与 s[left:] 的 lcp 大于等于 k 的后缀数组(sa)上的区间.
// 如果不存在,返回(-1,-1).
func (suf *SuffixArray32) LcpRange(left int32, k int32) (start, end int32) {
	if k > suf.n-left {
		return -1, -1
	}
	if k == 0 {
		return 0, suf.n
	}
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ32(suf.Height)
	}
	i := suf.Rank[left] + 1
	start = suf.minSt.MinLeft(i, func(e int32) bool { return e >= k }) - 1 // 向左找
	end = suf.minSt.MaxRight(i, func(e int32) bool { return e >= k })      // 向右找
	return
}

// 查询s[start:end)在s中的出现次数.
func (suf *SuffixArray32) Count(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > suf.n {
		end = suf.n
	}
	if start >= end {
		return 0
	}
	a, b := suf.LcpRange(start, end-start)
	return b - a
}

// 求任意两个后缀s[i:]和s[j:]的最长公共前缀(lcp).
func (suf *SuffixArray32) _lcp(i, j int32) int32 {
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ32(suf.Height)
	}
	if i == j {
		return suf.n - i
	}
	r1, r2 := suf.Rank[i], suf.Rank[j]
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	return suf.minSt.Query(r1+1, r2+1)
}

type LinearRMQ32 struct {
	n     int32
	nums  []int32
	small []int
	large [][]int32
}

func NewLinearRMQ32(nums []int32) *LinearRMQ32 {
	n := int32(len(nums))
	res := &LinearRMQ32{n: n, nums: nums}
	stack := make([]int32, 0, 64)
	small := make([]int, 0, n)
	var large [][]int32
	large = append(large, make([]int32, 0, n>>6))
	for i := int32(0); i < n; i++ {
		for len(stack) > 0 && nums[stack[len(stack)-1]] > nums[i] {
			stack = stack[:len(stack)-1]
		}
		tmp := 0
		if len(stack) > 0 {
			tmp = small[stack[len(stack)-1]]
		}
		small = append(small, tmp|(1<<(i&63)))
		stack = append(stack, i)
		if (i+1)&63 == 0 {
			large[0] = append(large[0], stack[0])
			stack = stack[:0]
		}
	}

	for i := int32(1); (i << 1) <= n>>6; i <<= 1 {
		csz := int32(n>>6 + 1 - (i << 1))
		v := make([]int32, csz)
		for k := int32(0); k < csz; k++ {
			back := large[len(large)-1]
			v[k] = res._getMin(back[k], back[k+i])
		}
		large = append(large, v)
	}

	res.small = small
	res.large = large
	return res
}

// 查询区间`[start, end)`中的最小值.
func (rmq *LinearRMQ32) Query(start, end int32) int32 {
	end--
	left := start>>6 + 1
	right := end >> 6
	if left < right {
		msb := bits.Len64(uint64(right-left)) - 1
		cache := rmq.large[msb]
		i := (left-1)<<6 + int32(bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63)))))
		cand1 := rmq._getMin(i, cache[left])
		j := right<<6 + int32(bits.TrailingZeros64(uint64(rmq.small[end])))
		cand2 := rmq._getMin(cache[right-(1<<msb)], j)
		return rmq.nums[rmq._getMin(cand1, cand2)]
	}
	if left == right {
		i := (left-1)<<6 + int32(bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63)))))
		j := left<<6 + int32(bits.TrailingZeros64(uint64(rmq.small[end])))
		return rmq.nums[rmq._getMin(i, j)]
	}
	return rmq.nums[right<<6+int32(bits.TrailingZeros64(uint64(rmq.small[end]&(^0<<(start&63)))))]
}

func (rmq *LinearRMQ32) _getMin(i, j int32) int32 {
	if rmq.nums[i] < rmq.nums[j] {
		return i
	}
	return j
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (st *LinearRMQ32) MaxRight(left int32, check func(e int32) bool) int32 {
	if left == st.n {
		return st.n
	}
	ok, ng := left, st.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(st.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (st *LinearRMQ32) MinLeft(right int32, check func(e int32) bool) int32 {
	if right == 0 {
		return 0
	}
	ok, ng := right, int32(-1)
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(st.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 用于求解`两个字符串s和t`相关性质的后缀数组.
type SuffixArray2 struct {
	SA     *SuffixArray32
	offset int32
}

// !ord值很大时,需要先离散化.
// !ords[i]>=0.
func NewSuffixArray2(ords1, ords2 []int32) *SuffixArray2 {
	newNums := append(ords1, ords2...)
	sa := NewSuffixArray(int32(len(newNums)), func(i int32) int32 { return newNums[i] })
	return &SuffixArray2{SA: sa, offset: int32(len(ords1))}
}

func NewSuffixArray2FromString(s, t string) *SuffixArray2 {
	ords1 := make([]int32, len(s))
	for i, c := range s {
		ords1[i] = c
	}
	ords2 := make([]int32, len(t))
	for i, c := range t {
		ords2[i] = c
	}
	return NewSuffixArray2(ords1, ords2)
}

// 求任意两个子串s[a,b)和t[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray2) Lcp(a, b int32, c, d int32) int32 {
	return suf.SA.Lcp(a, b, c+suf.offset, d+suf.offset)
}

// 比较任意两个子串s[a,b)和t[c,d)的字典序.
//
//	s[a,b) < t[c,d) 返回-1.
//	s[a,b) = t[c,d) 返回0.
//	s[a,b) > t[c,d) 返回1.
func (suf *SuffixArray2) CompareSubstr(a, b int32, c, d int32) int32 {
	return suf.SA.CompareSubstr(a, b, c+suf.offset, d+suf.offset)
}

func SuffixArray32Simple(n int32, f func(i int32) int32) (sa, rank, height []int32) {
	s := make([]byte, 0, n*4)
	for i := int32(0); i < n; i++ {
		v := f(i)
		s = append(s, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
	_sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))
	sa = make([]int32, 0, n)
	for _, v := range _sa {
		if v&3 == 0 {
			sa = append(sa, v>>2)
		}
	}
	rank = make([]int32, n)
	for i := int32(0); i < n; i++ {
		rank[sa[i]] = i
	}
	height = make([]int32, n)
	h := int32(0)
	for i := int32(0); i < n; i++ {
		rk := rank[i]
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := sa[rk-1]; i+h < n && j+h < n && f(i+h) == f(j+h); h++ {
			}
		}
		height[rk] = h
	}
	return
}

type ClampableStackItem = struct {
	value int
	count int32
}

type ClampableStack struct {
	clampMin bool
	total    int
	count    int
	stack    []ClampableStackItem
}

// clampMin：
//
//	为true时，调用AddAndClamp(x)后，容器内所有数最小值被截断(小于x的数变成x)；
//	为false时，调用AddAndClamp(x)后，容器内所有数最大值被截断(大于x的数变成x).
func NewClampableStack(clampMin bool) *ClampableStack {
	return &ClampableStack{clampMin: clampMin}
}

func (h *ClampableStack) AddAndClamp(x int) {
	newCount := 1
	if h.clampMin {
		for len(h.stack) > 0 {
			top := h.stack[len(h.stack)-1]
			if top.value > x {
				break
			}
			h.stack = h.stack[:len(h.stack)-1]
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	} else {
		for len(h.stack) > 0 {
			top := h.stack[len(h.stack)-1]
			if top.value < x {
				break
			}
			h.stack = h.stack[:len(h.stack)-1]
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	}
	h.total += x * newCount
	h.count++
	h.stack = append(h.stack, ClampableStackItem{value: x, count: int32(newCount)})
}

func (h *ClampableStack) Sum() int {
	return h.total
}

func (h *ClampableStack) Len() int {
	return h.count
}

func (h *ClampableStack) Clear() {
	h.stack = h.stack[:0]
	h.total = 0
	h.count = 0
}

// 求每个元素作为最值的影响范围(闭区间).
func GetRange(nums []int, isMax, isLeftStrict, isRightStrict bool) (leftMost, rightMost []int) {
	compareLeft := func(stackValue, curValue int) bool {
		if isLeftStrict && isMax {
			return stackValue <= curValue
		} else if isLeftStrict && !isMax {
			return stackValue >= curValue
		} else if !isLeftStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	compareRight := func(stackValue, curValue int) bool {
		if isRightStrict && isMax {
			return stackValue <= curValue
		} else if isRightStrict && !isMax {
			return stackValue >= curValue
		} else if !isRightStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	n := len(nums)
	leftMost, rightMost = make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		rightMost[i] = n - 1
	}

	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && compareRight(nums[stack[len(stack)-1]], nums[i]) {
			rightMost[stack[len(stack)-1]] = i - 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	stack = []int{}
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && compareLeft(nums[stack[len(stack)-1]], nums[i]) {
			leftMost[stack[len(stack)-1]] = i + 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	return
}

type MonoQueueValue = int32
type MonoQueue struct {
	MinQueue       []MonoQueueValue
	_minQueueCount []int32
	_less          func(a, b MonoQueueValue) bool
	_len           int32
}

func NewMonoQueue(less func(a, b MonoQueueValue) bool) *MonoQueue {
	return &MonoQueue{
		_less: less,
	}
}

func (q *MonoQueue) Append(value MonoQueueValue) *MonoQueue {
	count := int32(1)
	for len(q.MinQueue) > 0 && q._less(value, q.MinQueue[len(q.MinQueue)-1]) {
		q.MinQueue = q.MinQueue[:len(q.MinQueue)-1]
		count += q._minQueueCount[len(q._minQueueCount)-1]
		q._minQueueCount = q._minQueueCount[:len(q._minQueueCount)-1]
	}
	q.MinQueue = append(q.MinQueue, value)
	q._minQueueCount = append(q._minQueueCount, count)
	q._len++
	return q
}

func (q *MonoQueue) Popleft() {
	q._minQueueCount[0]--
	if q._minQueueCount[0] == 0 {
		q.MinQueue = q.MinQueue[1:]
		q._minQueueCount = q._minQueueCount[1:]
	}
	q._len--
}

func (q *MonoQueue) Head() MonoQueueValue {
	return q.MinQueue[0]
}

func (q *MonoQueue) Min() MonoQueueValue {
	return q.MinQueue[0]
}

func (q *MonoQueue) Len() int32 {
	return q._len
}

func (q *MonoQueue) String() string {
	sb := []string{}
	for i := 0; i < len(q.MinQueue); i++ {
		sb = append(sb, fmt.Sprintf("%v", pair{q.MinQueue[i], q._minQueueCount[i]}))
	}
	return fmt.Sprintf("MonoQueue{%v}", strings.Join(sb, ", "))
}

type pair struct {
	value MonoQueueValue
	count int32
}

func (p pair) String() string {
	return fmt.Sprintf("(value: %v, count: %v)", p.value, p.count)
}

type UnionFindArray struct {
	data []int32
}

func NewUnionFindArray(n int32) *UnionFindArray {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{data: data}
}

// f: 合并两个分组前的钩子函数.
func (ufa *UnionFindArray) Union(key1, key2 int32, f func(big, small int32)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	if f != nil {
		f(root1, root2)
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *UnionFindArray) Find(key int32) int32 {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindArray) Size(key int32) int32 {
	return -ufa.data[ufa.Find(key)]
}

func NewHeap[H any](less func(a, b H) bool, nums []H) *Heap[H] {
	nums = append(nums[:0:0], nums...)
	heap := &Heap[H]{less: less, data: nums}
	heap.heapify()
	return heap
}

type Heap[H any] struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap[H]) Top() (value H) {
	value = h.data[0]
	return
}

func (h *Heap[H]) Len() int { return len(h.data) }

func (h *Heap[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root
		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}
		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}
		if minIndex == root {
			return
		}
		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
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
func (d *Dictionary[V]) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return int(res)
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = int32(id)
	return id
}
func (d *Dictionary[V]) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary[V]) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary[V]) Size() int {
	return len(d._idToValue)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b

}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}

func mins32(a []int32) int32 {
	mn := a[0]
	for _, x := range a {
		if x < mn {
			mn = x
		}
	}
	return mn
}

func maxs32(a []int32) int32 {
	mx := a[0]
	for _, x := range a {
		if x > mx {
			mx = x
		}
	}
	return mx
}
