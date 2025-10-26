package main

import (
	"index/suffixarray"
	"reflect"
	"sort"
	"unsafe"
)

// 和可被K整除的本质不同子数组个数
// https://leetcode.cn/problems/count-distinct-subarrays-divisible-by-k-in-sorted-array/solutions/3815839/golang-hou-zhui-shu-zu-by-981377660lmt-sghj/
func numGoodSubarrays(nums []int, k int) int64 {
	n := int32(len(nums))
	if n == 0 {
		return 0
	}

	presum := make([]int, n+1)
	for i := int32(0); i < n; i++ {
		x := (presum[i] + nums[i]) % k
		if x < 0 {
			x += k
		}
		presum[i+1] = x
	}

	mp := make(map[int][]int32, n+1)
	for e := int32(0); e <= n; e++ {
		mp[presum[e]] = append(mp[presum[e]], e)
	}

	newNums, _ := Discretize(nums)
	sa, _, height := SuffixArray32(n, func(i int32) int32 { return newNums[i] })

	res := 0
	for i := int32(0); i < n; i++ {
		s := sa[i]
		h := height[i]
		threshold := s + h + 1 // SA+LCP 去重：每个后缀只新增长度 > height 的子数组
		ends := mp[presum[s]]
		idx := sort.Search(len(ends), func(j int) bool { return ends[j] >= threshold })
		res += len(ends) - idx
	}

	return int64(res)
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

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
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
