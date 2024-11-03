// 比较多个字符串的字典序.
//
// 1. 比较两个字符串的字典序大小.
// 2. 求两个字符串的最长公共前缀.
// 3. 比较两个字符串拼接后的字典序大小.用于求拼接最小数/拼接最大数.

package main

import (
	"fmt"
	"index/suffixarray"
	"math/bits"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

// LCR 164. 破解闯关密码
// https://leetcode.cn/problems/ba-shu-zu-pai-cheng-zui-xiao-de-shu-lcof/
// 拼接最小数/拼接字典序最小的字符串.
func crackPassword(password []int) string {
	strs := make([]string, len(password))
	for i := range password {
		strs[i] = strconv.Itoa(password[i])
	}
	C := NewManyStringCompare(int32(len(password)), func(nth int32) (int32, func(int32) int32) {
		return int32(len(strs[nth])), func(i int32) int32 {
			return int32(strs[nth][i])
		}
	})

	order := make([]int, len(strs))
	for i := range order {
		order[i] = i
	}
	less := func(i, j int) bool {
		i, j = order[i], order[j]
		return C.CompareSTTS(int32(i), 0, C.Len(int32(i)), int32(j), 0, C.Len(int32(j))) < 0 // 拼接最小数
	}
	swap := func(i, j int) {
		order[i], order[j] = order[j], order[i]
		strs[i], strs[j] = strs[j], strs[i]
	}
	Sort(len(strs), less, swap)

	res := strings.Builder{}
	for i := range strs {
		res.WriteString(strs[i])
	}
	return res.String()
}

func demo() {
	words := []string{"ab", "abc", "abcd", "abcde", "abcdef"}
	C := NewManyStringCompare(int32(len(words)), func(nth int32) (int32, func(int32) int32) {
		return int32(len(words[nth])), func(i int32) int32 {
			return int32(words[nth][i])
		}
	})

	// 比较words[0]和words[1]的字典序大小.
	fmt.Println(C.Compare(0, 0, 2, 1, 0, 2)) // -1

	// words[0]和words[1]的最长公共前缀.
	fmt.Println(C.Lcp(0, 0, 2, 1, 0, 3)) // 2

	// words[0]+words[1] 与 words[1]+words[0]的字典序大小.
	fmt.Println(C.CompareSTTS(0, 0, 2, 1, 0, 3)) // -1

	{

		// 求能拼接出的字典序最小的字符串(拼接最小数/拼接最大数).
		nums := [][]int32{
			{1, 2, 3},
			{3, 4, 5},
			{1, 2, 4},
		}

		C := NewManyStringCompare(int32(len(nums)), func(nth int32) (int32, func(int32) int32) {
			return int32(len(nums[nth])), func(i int32) int32 {
				return nums[nth][i]
			}
		})

		order := make([]int, len(nums))
		for i := range order {
			order[i] = i
		}
		less := func(i, j int) bool {
			i, j = order[i], order[j]
			return C.CompareSTTS(int32(i), 0, C.Len(int32(i)), int32(j), 0, C.Len(int32(j))) > 0 // 拼接最大数
		}
		swap := func(i, j int) {
			order[i], order[j] = order[j], order[i]
			nums[i], nums[j] = nums[j], nums[i]
		}
		Sort(len(nums), less, swap)

		fmt.Println(nums) // [[1 2 3] [1 2 4] [3 4 5]]
	}
}

type ManyStringCompare struct {
	n   int32
	all []int32
	pos []int32
	s   *S
}

func NewManyStringCompare(n int32, f func(nth int32) (len int32, g func(int32) int32)) *ManyStringCompare {
	all := []int32{}
	pos := make([]int32, 0, n+1)
	pos = append(pos, 0)
	for i := int32(0); i < n; i++ {
		v, g := f(i)
		for j := int32(0); j < v; j++ {
			all = append(all, g(j))
		}
		pos = append(pos, int32(len(all)))
	}
	s := NewS(int32(len(all)), func(i int32) int32 { return all[i] })
	return &ManyStringCompare{n: n, all: all, pos: pos, s: s}
}

// S[i][la:ra), S[j][lb:rb)
func (msc *ManyStringCompare) Lcp(i, la, ra, j, lb, rb int32) int32 {
	pos := msc.pos
	n := msc.s.Lcp(pos[i]+la, pos[i]+ra, pos[j]+lb, pos[j]+rb)
	return min32(n, min32(ra-la, rb-lb))
}

// 比较S[i][la:ra)和S[j][lb:rb)字典序大小.
// 返回-1/0/1.
func (msc *ManyStringCompare) Compare(i, la, ra, j, lb, rb int32) int32 {
	na := ra - la
	nb := rb - lb
	if na > nb {
		return -msc.Compare(j, lb, rb, i, la, ra)
	}
	n := msc.Lcp(i, la, ra, j, lb, rb)
	if n == na {
		if na == nb {
			return 0
		}
		return -1
	}
	all, pos := msc.all, msc.pos
	if all[pos[i]+la+n] < all[pos[j]+lb+n] {
		return -1
	}
	return 1
}

// 另s=S[i][la:ra), t=S[j][lb:rb).
// 比较s+t和t+s的字典序大小.
// 返回-1/0/1.
//
// 使用场景：求能拼接出的字典序最小的字符串(拼接最小序).
func (msc *ManyStringCompare) CompareSTTS(i, la, ra, j, lb, rb int32) int32 {
	na := ra - la
	nb := rb - lb
	if na > nb {
		return -msc.CompareSTTS(j, lb, rb, i, la, ra)
	}
	k := msc.Compare(i, la, la+na, j, lb, lb+na)
	if k != 0 {
		return k
	}
	k = msc.Compare(j, lb, lb+nb-na, j, lb+na, rb)
	if k != 0 {
		return k
	}
	return msc.Compare(j, lb+nb-na, rb, i, la, ra)
}

func (msc *ManyStringCompare) Len(i int32) int32 {
	return msc.pos[i+1] - msc.pos[i]
}

type S struct {
	Sa     []int32 // 排名第i的后缀是谁.
	Rank   []int32 // 后缀s[i:]的排名是多少.
	Height []int32 // 排名相邻的两个后缀的最长公共前缀.Height[0] = 0,Height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	Ords   []int32
	n      int32
	minSt  *LinearRMQ32 // 维护lcp的最小值
}

func NewS(n int32, f func(i int32) int32) *S {
	res := &S{n: n}
	sa, rank, lcp := SuffixArray32Simple(n, f)
	res.Sa, res.Rank, res.Height = sa, rank, lcp
	return res
}

// 求任意两个子串s[a,b)和s[c,d)的最长公共前缀(lcp).
func (suf *S) Lcp(a, b int32, c, d int32) int32 {
	if a >= b || c >= d {
		return 0
	}
	cand := suf._lcp(a, c)
	return min32(cand, min32(b-a, d-c))
}

// 求任意两个后缀s[i:]和s[j:]的最长公共前缀(lcp).
func (suf *S) _lcp(i, j int32) int32 {
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

// 用于记录排序过程中的交换.
func Sort(n int, less func(i, j int) bool, swap func(i, j int)) {
	sort.Sort(&sorter{n: n, less: less, swap: swap})
}

type sorter struct {
	n    int
	less func(i, j int) bool
	swap func(i, j int)
}

func (s *sorter) Len() int           { return s.n }
func (s *sorter) Less(i, j int) bool { return s.less(i, j) }
func (s *sorter) Swap(i, j int) {
	if i == j {
		return
	}
	s.swap(i, j)
}
