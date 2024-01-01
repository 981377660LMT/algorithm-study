// 判断数据结构中每个数出现的次数是否均为k.
// 等价于:
//  1. 数据结构中每个数出现的次数均为k的倍数：异或哈希.
//  2. 数据结构中每个数出现的次数均不超过k：双指针.
//     在右指针扫到 i 的时候，不停将左指针向右移动并减去这个桶的出现次数，
//     直到 nums[i] 的出现次数小于等于 k 为止。此时再统计答案，两个限制都可以满足。

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
)

func main() {
	CF1418G()
}

// https://www.luogu.com.cn/problem/solution/CF1418G
func CF1418G() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	fmt.Fprintln(out, CountSubarrayWithFrequencyEqualToK(nums, 3))
}

// https://leetcode.cn/problems/count-complete-substrings
func countCompleteSubstrings(word string, k int) int {
	n := len(word)
	ords := make([]int, n)
	for i := 0; i < n; i++ {
		ords[i] = int(word[i] - 'a')
	}

	groups := [][]int{}
	ptr := 0
	for ptr < n {
		leader := ords[ptr]
		group := []int{leader}
		ptr++
		for ptr < n && abs(ords[ptr]-ords[ptr-1]) <= 2 {
			group = append(group, ords[ptr])
			ptr++
		}
		groups = append(groups, group)
	}

	res := 0
	for _, group := range groups {
		res += CountSubarrayWithFrequencyEqualToK(group, k)
	}
	return res
}

type Value = int

func RandomHash(min, max uint64) func(value Value) uint64 {
	pool := make(map[Value]uint64)
	f := func(value Value) uint64 {
		if hash, ok := pool[value]; ok {
			return hash
		}
		rand := rand.Uint64()%(max-min+1) + min
		pool[value] = rand
		return rand
	}
	return f
}

// 统计满足`每个元素出现的次数均为k`条件的子数组的个数.
func CountSubarrayWithFrequencyEqualToK(arr []Value, k int) int {
	n := len(arr)
	if n == 0 || k <= 0 || k > n {
		return 0
	}

	R := RandomHash(1, math.MaxUint64/uint64(n))

	pool := make(map[Value]int)
	getId := func(value Value) int {
		if id, ok := pool[value]; ok {
			return id
		}
		id := len(pool)
		pool[value] = id
		return id
	}
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = getId(arr[i])
	}
	counter := make([]int, len(pool))
	random := make([]uint64, n)
	for i := 0; i < n; i++ {
		random[i] = R(nums[i])
	}
	hashPreSum := make([]uint64, n+1)
	for i := 0; i < n; i++ {
		hashPreSum[i+1] = hashPreSum[i]
		hashPreSum[i+1] -= uint64(counter[nums[i]]) * random[i]
		counter[nums[i]] = (counter[nums[i]] + 1) % k
		hashPreSum[i+1] += uint64(counter[nums[i]]) * random[i]
	}

	countPreSum := make(map[uint64]int)
	countPreSum[0] = 1
	counter = make([]int, len(pool))
	res, left := 0, 0
	for right, num := range nums {
		counter[num]++
		for counter[num] > k {
			counter[nums[left]]--
			countPreSum[hashPreSum[left]]--
			left++
		}
		res += countPreSum[hashPreSum[right+1]]
		countPreSum[hashPreSum[right+1]]++
	}

	return res
}

type hashKey = struct {
	value    Value
	modCount int
}

// 判断数据结构中每个数出现的次数是否均恰好为k.
// 如果为空集合,则返回True.
type AllCountKChecker struct {
	hash       uint64
	counter    map[Value]int
	modCounter map[Value]int
	k          int
	countPq    *ErasableHeap
	pool       map[hashKey]uint64
}

func NewAllCountKChecker(k int) *AllCountKChecker {
	return &AllCountKChecker{
		hash:       0,
		counter:    make(map[Value]int),
		modCounter: make(map[Value]int),
		k:          k,
		countPq:    NewErasableHeap(func(a, b int) bool { return a > b }, nil),
		pool:       make(map[hashKey]uint64),
	}
}

func (c *AllCountKChecker) Add(v Value) {
	count := c.modCounter[v]
	random := c.randomHash(hashKey{v, count})
	c.hash ^= random
	count++
	if count == c.k {
		count = 0
	}
	c.hash ^= c.randomHash(hashKey{v, count})
	c.modCounter[v] = count

	preCount := c.counter[v]
	c.counter[v]++
	if preCount > 0 {
		c.countPq.Erase(preCount)
	}
	c.countPq.Push(preCount + 1)
}

// 删除前需要保证v在集合中存在.
func (c *AllCountKChecker) Remove(v Value) {
	count := c.modCounter[v]
	c.hash ^= c.randomHash(hashKey{v, count})
	count--
	if count == -1 {
		count = c.k - 1
	}
	c.hash ^= c.randomHash(hashKey{v, count})
	c.modCounter[v] = count

	preCount := c.counter[v]
	c.counter[v]--
	if preCount == 1 {
		delete(c.counter, v)
	}
	c.countPq.Erase(preCount)
	if preCount > 1 {
		c.countPq.Push(preCount - 1)
	}
}

// 询问数据结构中每个数出现的次数是否均k的倍数.
func (c *AllCountKChecker) Query() bool {
	if c.countPq.Len() == 0 {
		return true
	}
	return c.hash == 0 && c.countPq.Peek() == c.k
}

func (c *AllCountKChecker) GetHash() uint64 {
	return c.hash
}

func (c *AllCountKChecker) Clear() {
	c.hash = 0
	c.modCounter = make(map[Value]int)
	c.counter = make(map[Value]int)
	c.countPq.Clear()
}

func (c *AllCountKChecker) String() string {
	return fmt.Sprintf("hash=%d, counter=%v", c.hash, c.counter)
}

func (c *AllCountKChecker) randomHash(key hashKey) uint64 {
	if hash, ok := c.pool[key]; ok {
		return hash
	}
	rand := rand.Uint64()%(1<<61-1) + 1
	c.pool[key] = rand
	return rand
}

type H = int

type ErasableHeap struct {
	base   *Heap
	erased *Heap
}

func NewErasableHeap(less func(a, b H) bool, nums []H) *ErasableHeap {
	return &ErasableHeap{NewHeap(less, nums), NewHeap(less, nil)}
}

// 从堆中删除一个元素,要保证堆中存在该元素.
func (h *ErasableHeap) Erase(value H) {
	h.erased.Push(value)
	h.normalize()
}

func (h *ErasableHeap) Push(value H) {
	h.base.Push(value)
	h.normalize()
}

func (h *ErasableHeap) Pop() (value H) {
	value = h.base.Pop()
	h.normalize()
	return
}

func (h *ErasableHeap) Peek() (value H) {
	value = h.base.Top()
	return
}

func (h *ErasableHeap) Len() int {
	return h.base.Len()
}

func (h *ErasableHeap) Clear() {
	h.base.Clear()
	h.erased.Clear()
}

func (h *ErasableHeap) normalize() {
	for h.base.Len() > 0 && h.erased.Len() > 0 && h.base.Top() == h.erased.Top() {
		h.base.Pop()
		h.erased.Pop()
	}
}

func NewHeap(less func(a, b H) bool, nums []H) *Heap {
	nums = append(nums[:0:0], nums...)
	heap := &Heap{less: less, data: nums}
	if len(nums) > 1 {
		heap.heapify()
	}
	return heap
}

type Heap struct {
	data []H
	less func(a, b H) bool
}

func (h *Heap) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *Heap) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *Heap) Top() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *Heap) Len() int { return len(h.data) }

func (h *Heap) Clear() {
	h.data = h.data[:0]
}

func (h *Heap) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *Heap) pushDown(root int) {
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
