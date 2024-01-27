// Inversion 逆序对

package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	inv := AllRangeInversion([]int{5, 4, 3, 2, 1})
	fmt.Println(inv[0][5])
	fmt.Println(CountInversionRotate([]int{5, 4, 3, 2, 1}))
}

// https://leetcode.cn/problems/shu-zu-zhong-de-ni-xu-dui-lcof/
func reversePairs(record []int) int {
	return CountInversion(record)
}

func CountInversion(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	res := 0
	var mergeSort func(nums []int, l int, r int, temp []int)
	mergeSort = func(nums []int, l int, r int, temp []int) {
		if l >= r {
			return
		}
		mid := (l + r) / 2
		mergeSort(nums, l, mid, temp)
		mergeSort(nums, mid+1, r, temp)
		i := l
		j := mid + 1
		k := l
		for i <= mid && j <= r {
			if nums[i] > nums[j] {
				res += (mid - i + 1)
				temp[k] = nums[j]
				j++
				k++
			} else {
				temp[k] = nums[i]
				i++
				k++
			}
		}
		for i <= mid {
			temp[k] = nums[i]
			i++
			k++
		}
		for j <= r {
			temp[k] = nums[j]
			j++
			k++
		}
		for t := l; t <= r; t++ {
			nums[t] = temp[t]
		}
	}
	mergeSort(nums, 0, len(nums)-1, make([]int, len(nums)))
	return res
}

// 轮转逆序对.
// 返回一个数组，第 i 个元素表示将nums[i]作为首元素时的逆序对数.
func CountInversionRotate(nums []int) []int {
	if len(nums) == 0 {
		return nil
	}
	n := len(nums)
	nums, size := Discretize(nums)
	base := 0
	bit := NewBitArray(size)
	for _, v := range nums {
		base += bit.QueryRange(v+1, size)
		bit.Add(v, 1)
	}
	res := make([]int, n)
	for i, v := range nums {
		res[i] = base
		base += bit.QueryRange(v+1, size) - bit.QueryPrefix(v)
	}
	return res
}

// 区间逆序对.
// 返回一个(n+1*n+1)的二维数组，inv[i][j] 表示 nums[i:j] 的逆序对数.
func AllRangeInversion(nums []int) (inv [][]int) {
	n := len(nums)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for left := n; left >= 0; left-- {
		for right := left + 2; right <= n; right++ {
			dp[left][right] = dp[left][right-1] + dp[left+1][right] - dp[left+1][right-1]
			if nums[left] > nums[right-1] {
				dp[left][right]++
			}
		}
	}
	return dp
}

func Discretize(nums []int) (newNums []int, size int) {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	allNums := make([]int, 0, len(set))
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	mp := make(map[int]int, len(allNums))
	for i, v := range allNums {
		mp[v] = i
	}
	newNums = make([]int, len(nums))
	for i, v := range nums {
		newNums[i] = mp[v]
	}
	return newNums, len(allNums)
}

func maxs(nums []int) int {
	res := nums[0]
	for _, v := range nums {
		if v > res {
			res = v
		}
	}
	return res
}

// !Point Add Range Sum, 0-based.
type BITArray struct {
	n     int
	total int
	data  []int
}

func NewBitArray(n int) *BITArray {
	res := &BITArray{n: n, data: make([]int, n)}
	return res
}

func NewBitArrayFrom(n int, f func(i int) int) *BITArray {
	total := 0
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &BITArray{n: n, total: total, data: data}
}

func (b *BITArray) Add(index int, v int) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *BITArray) QueryAll() int {
	return b.total
}

func (b *BITArray) MaxRight(check func(index, preSum int) bool) int {
	i := 0
	s := 0
	k := 1
	for 2*k <= b.n {
		k *= 2
	}
	for k > 0 {
		if i+k-1 < b.n {
			t := s + b.data[i+k-1]
			if check(i+k, t) {
				i += k
				s = t
			}
		}
		k >>= 1
	}
	return i
}

// 0/1 树状数组查找第 k(0-based) 个1的位置.
// UpperBound.
func (b *BITArray) Kth(k int) int {
	return b.MaxRight(func(index, preSum int) bool { return preSum <= k })
}

func (b *BITArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}
