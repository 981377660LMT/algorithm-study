// LexMaxSuffixOfAllPrefix/MaxLexSuffixOfAllPrefix
// 所有前缀的字典序最大后缀(的长度).
// https://www.codechef.com/START137A/problems/CABABAA

package main

import (
	"fmt"
	"index/suffixarray"
	"reflect"
	"unsafe"
)

func main() {
	s := "abacaba"
	res := LexMaxSuffixOfAllPrefix(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	fmt.Println(res)
}

// res[i] 表示前缀 s[0:i) 的字典序最大后缀的长度.
// 按照sa的顺序消除字符.
func LexMaxSuffixOfAllPrefix(n int32, f func(i int32) int32) []int32 {
	sa, _, height := SuffixArray32(n, f)
	toRemove := make([][]int32, n)
	var stack [][2]int32
	for i := int32(0); i < n; i++ {
		j := sa[i]
		var k int32
		if i == 0 {
			k = n
		} else {
			k = height[i]
		}
		for len(stack) > 0 && stack[len(stack)-1][0] > j {
			if tmp := stack[len(stack)-1][1]; tmp < k {
				k = tmp
			}
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			toRemove[j+k] = append(toRemove[j+k], j)
		}
		stack = append(stack, [2]int32{j, k})
	}

	p := n - 1
	res := make([]int32, n+1)
	bad := make([]bool, n)
	for i := n; i > 0; i-- {
		for bad[sa[p]] || i <= sa[p] {
			p--
		}
		res[i] = i - sa[p]
		for _, j := range toRemove[i-1] {
			bad[j] = true
		}
	}
	return res
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
