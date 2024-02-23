package main

import (
	"bufio"
	"fmt"
	"index/suffixarray"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

func main() {
	// P2852()
	abc141e()
}

// https://www.luogu.com.cn/problem/P2852
func P2852() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int32
	fmt.Fscan(in, &n, &k)
	nums := make([]int32, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	start, end := LongestRepeatSubstring(n, func(i int32) int32 { return nums[i] }, true, k)
	fmt.Fprintln(out, end-start)
}

// https://atcoder.jp/contests/abc141/tasks/abc141_e
func abc141e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	start, end := LongestRepeatSubstring(n, func(i int32) int32 { return int32(s[i]) }, false, 2)
	fmt.Fprintln(out, end-start)
}

// 返回任意一个符合条件的最长重复子串的起始位置和结束位置.
func LongestRepeatSubstring(
	n int32, f func(i int32) int32,
	allowDuplicate bool, // 是否允许重复
	minRepeatCount int32, // 至少出现的次数
) (start, end int32) {
	if allowDuplicate {
		if minRepeatCount == 2 {
			return solve1(n, f)
		} else {
			return solve2(n, f, minRepeatCount)
		}
	} else {
		if minRepeatCount != 2 {
			panic("暂不支持")
		}
		return solve3(n, f)
	}
}

// 可重叠最长重复子串 (可重叠的至少出现 2 次的最长重复子串)
// !高度数组中的最大值对应的就是可重叠最长重复子串
func solve1(n int32, f func(i int32) int32) (start, end int32) {
	sa, _, height := SuffixArray32(n, f)
	saIndex, maxHeight := int32(0), int32(0)
	for i := int32(0); i < n; i++ {
		h := height[i]
		if h > maxHeight {
			saIndex = i
			maxHeight = h
		}
	}
	return sa[saIndex], sa[saIndex] + maxHeight
}

// 可重叠的至少出现 k 次的最长重复子串(k>2)
// https://www.luogu.com.cn/problem/P2852
// 出现至少 k 次意味着后缀排序后有至少连续 k 个后缀以这个子串作为公共前缀。
// 所以，单调队列求出每相邻 k-1 个 height 的最小值，再求这些最小值的最大值就是答案。
func solve2(n int32, f func(i int32) int32, k int32) (start, end int32) {
	sa, _, height := SuffixArray32(n, f)
	minQueue := NewMonoQueue(func(i, j MonoQueueItem) bool {
		return height[i] < height[j]
	})
	maxHeight := int32(0)
	for i := int32(0); i < n; i++ {
		minQueue.Append(i)
		if minQueue.Len() == k-1 {
			minIndex := minQueue.Min()
			curHeight := height[minIndex]
			if curHeight > maxHeight {
				maxHeight = curHeight
				start, end = sa[minIndex], sa[minIndex]+curHeight
			}
			minQueue.Popleft()
		}
	}
	return
}

// 不可重叠最长重复子串 (不可重叠的至少出现 2 次的最长重复子串)
// 二分目标串的长度|s| ，将height数组划分成若干个连续 LCP 大于等于|s|的段，维护每个段中出现的数中最大和最小的sa下标，
// 若这两个下标的距离满足条件，则一定有长度为|s|的字符串不重叠地出现了两次。
// https://atcoder.jp/contests/abc141/tasks/abc141_e
// https://www.cnblogs.com/xiaoyh/p/10328219.html
func solve3(n int32, f func(i int32) int32) (start, end int32) {
	sa, _, height := SuffixArray32(n, f)

	check := func(mid int32) (int32, bool) {
		minSa, maxSa := sa[0], sa[0]
		for i := int32(1); i < n; i++ {
			if height[i] >= mid {
				minSa, maxSa = min32(sa[i], minSa), max32(sa[i], maxSa)
				if maxSa-minSa >= mid {
					return minSa, true
				}
			} else {
				minSa, maxSa = sa[i], sa[i]
			}
		}
		return minSa, maxSa-minSa >= mid
	}

	left, right := int32(1), n
	for left <= right {
		mid := (left + right) >> 1
		tmpStart, ok := check(mid)
		if ok {
			left = mid + 1
			start, end = tmpStart, tmpStart+mid
		} else {
			right = mid - 1
		}
	}
	return
}

func SuffixArray32(n int32, f func(i int32) int32) (sa, rank, height []int32) {
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

type MonoQueueItem = int32

// 单调队列维护滑动窗口最小值.
// 单调队列队头元素为当前窗口最小值.
type MonoQueue struct {
	MinQueue       []MonoQueueItem
	_minQueueCount []int32
	_less          func(a, b MonoQueueItem) bool
	_len           int32
}

func NewMonoQueue(less func(a, b MonoQueueItem) bool) *MonoQueue {
	return &MonoQueue{
		_less: less,
	}
}

func (q *MonoQueue) Append(value MonoQueueItem) *MonoQueue {
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

func (q *MonoQueue) Head() MonoQueueItem {
	return q.MinQueue[0]
}

func (q *MonoQueue) Min() MonoQueueItem {
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
	value MonoQueueItem
	count int32
}

func (p pair) String() string {
	return fmt.Sprintf("(value: %v, count: %v)", p.value, p.count)
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

func max32(a, b int32) int32 {
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
