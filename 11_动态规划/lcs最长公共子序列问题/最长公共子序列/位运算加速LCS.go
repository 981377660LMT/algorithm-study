// lcsbit/bitlcs
// https://atcoder.jp/contests/dp/submissions/34604402
// https://loj.ac/s/1633431

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"time"
)

func main() {
	test()
}

func test() {
	// 10000*10000
	nums1 := make([]int, int(1e5))
	nums2 := make([]int, int(1e5))
	for i := range nums1 {
		nums1[i] = i
		nums2[i] = i
	}
	time1 := time.Now()
	fmt.Println(LCSBit(nums1, nums2))
	fmt.Println(time.Since(time1))
}

func demo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s1, s2 string
	fmt.Fscan(in, &s1, &s2)
	fmt.Fprintln(out, LCSBit2(s1, s2))
}

func LCSBit2(s1, s2 string) int {
	ords1, ords2 := make([]int, len(s1)), make([]int, len(s2))
	for i, v := range s1 {
		ords1[i] = int(v)
	}
	for i, v := range s2 {
		ords2[i] = int(v)
	}
	return LCSBit(ords1, ords2)
}

// 位运算加速LCS(最长公共子序列).
func LCSBit(arr1 []int, arr2 []int) int {
	arr1, arr2 = append(arr1[:0:0], arr1...), append(arr2[:0:0], arr2...)
	id := make(map[int]int)
	for i, v := range arr1 {
		if _, ok := id[v]; !ok {
			id[v] = len(id)
		}
		arr1[i] = id[v]
	}
	for i, v := range arr2 {
		if _, ok := id[v]; !ok {
			id[v] = len(id)
		}
		arr2[i] = id[v]
	}

	n := len(arr1)
	size := len(id)
	f := make([]*BS, size)
	for i := range f {
		f[i] = NewBS(n)
	}
	for i, v := range arr1 {
		f[v].Set(i)
	}

	dp := NewBS(n)
	for _, v := range arr2 {
		dp.Run(f[v])
	}
	return dp.Count()
}

type BS struct {
	data []uint64
}

func NewBS(n int) *BS {
	return &BS{
		data: make([]uint64, 1+(n/63)),
	}
}

func (bs *BS) Set(i int) {
	bs.data[i/63] |= 1 << (i % 63)
}

func (bs *BS) Count() int {
	res := 0
	for _, v := range bs.data {
		res += bits.OnesCount64(v)
	}
	return res
}

func (bs *BS) Run(o *BS) {
	c := uint64(1)
	for i, v := range bs.data {
		x, y := v, v|o.data[i]
		x += x + c + (^y & (1<<63 - 1))
		bs.data[i] = x & y
		c = x >> 63
	}
}
