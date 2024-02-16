// 可以插入字符的在线z函数
// https://ei1333.github.io/library/string/z-algorithm.hpp
// 每个后缀s[i:]与s的最长公共前缀长度
// aaabaaaab
// 921034210

// ZAlgorithm(s)  // s: string
// Append(c)      // c: string
// Get(i)         // i: index
// GetAll()       // return []int

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

func demo() {
	s := "aaabaaaab"
	Z := NewZAlgoOnline(0, func(i int) int { return int(s[i]) })
	fmt.Println(Z.GetAll())
	Z.Append(97) // 97 is the ascii code of 'a'
	fmt.Println(Z.GetAll())
	Z.Append(97) // 97 is the ascii code of 'a'
	fmt.Println(Z.GetAll())
	fmt.Println(Z.ords)
}

func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	Z := NewZAlgoOnline(0, func(i int) int { return int(s[i]) })
	for i := 0; i < n; i++ {
		Z.Append(int(s[i]))
	}
	res := Z.GetAll()
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

type ZAlgoOnline struct {
	ords        []int
	deleted     []bool
	deletedHist [][]int32
	z           []int32
	cur         []int32
}

func NewZAlgoOnline(n int, f func(i int) int) *ZAlgoOnline {
	res := &ZAlgoOnline{}
	for i := 0; i < n; i++ {
		res.Append(f(i))
	}
	return res
}

func (za *ZAlgoOnline) Append(c int) {
	za.ords = append(za.ords, c)
	za.deletedHist = append(za.deletedHist, []int32{})
	za.deleted = append(za.deleted, false)
	za.z = append(za.z, 0)
	za.z[0]++

	len_ := int32(len(za.ords))
	if len_ == 1 {
		return
	}
	if za.ords[0] == c {
		za.cur = append(za.cur, len_-1)
	} else {
		za.deleted[len_-1] = true
	}

	set := func(t, l int32) {
		za.deleted[t] = true
		for int32(len(za.deletedHist)) <= l {
			za.deletedHist = append(za.deletedHist, []int32{})
		}
		za.deletedHist[l] = append(za.deletedHist[l], t)
		za.z[t] = l - t - 1
	}

	for len(za.cur) > 0 {
		t := za.cur[0]
		if za.deleted[t] {
			za.cur = za.cur[1:]
		} else if za.ords[len_-t-1] == za.ords[len(za.ords)-1] {
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

func (za *ZAlgoOnline) Get(i int) int {
	if za.deleted[i] {
		return int(za.z[i])
	}
	return len(za.ords) - i
}

func (za *ZAlgoOnline) GetAll() []int {
	res := make([]int, len(za.ords))
	for i := 0; i < len(za.ords); i++ {
		res[i] = za.Get(i)
	}
	return res
}
