// 可以插入字符的在线z函数
// https://ei1333.github.io/library/string/z-algorithm.hpp
// 每个后缀s[i:]与s的最长公共前缀长度
// aaabaaaab
// 921034210

// ZAlgorithm(s)  // s: string
// Append(c)      // c: string
// Get(i)         // i: index
// GetAll()          // return []int

package main

import (
	"fmt"
)

func main() {
	Z := NewZAlgorithm("aaabaaaab")
	fmt.Println(Z.GetAll())
	Z.Append("a")
	fmt.Println(Z.GetAll())
	Z.Append("a")
	fmt.Println(Z.GetAll())
	fmt.Println(Z.S)
}

type ZAlgorithm struct {
	S           []string
	deleted     []bool
	deletedHist [][]int
	z           []int
	cur         []int
}

func NewZAlgorithm(s string) *ZAlgorithm {
	res := &ZAlgorithm{}
	for _, c := range s {
		res.Append(string(c))
	}
	return res
}

func (za *ZAlgorithm) Append(char string) {
	za.S = append(za.S, char)
	za.deletedHist = append(za.deletedHist, []int{})
	za.deleted = append(za.deleted, false)
	za.z = append(za.z, 0)
	za.z[0]++

	len_ := len(za.S)
	if len_ == 1 {
		return
	}
	if za.S[0] == char {
		za.cur = append(za.cur, len_-1)
	} else {
		za.deleted[len_-1] = true
	}

	set := func(t, l int) {
		za.deleted[t] = true
		for len(za.deletedHist) <= l {
			za.deletedHist = append(za.deletedHist, []int{})
		}
		za.deletedHist[l] = append(za.deletedHist[l], t)
		za.z[t] = l - t - 1
	}

	for len(za.cur) > 0 {
		t := za.cur[0]
		if za.deleted[t] {
			za.cur = za.cur[1:]
		} else if za.S[len_-t-1] == za.S[len(za.S)-1] {
			break
		} else {
			set(t, len_)
			za.cur = za.cur[1:]
		}
	}

	if len(za.cur) > 0 {
		t := za.cur[0]
		for _, p := range za.deletedHist[len_-t] {
			set(p+t, len_)
		}
	}

}

func (za *ZAlgorithm) Get(i int) int {
	if za.deleted[i] {
		return za.z[i]
	}
	return len(za.S) - i
}

func (za *ZAlgorithm) GetAll() []int {
	res := make([]int, len(za.S))
	for i := 0; i < len(za.S); i++ {
		res[i] = za.Get(i)
	}
	return res
}
