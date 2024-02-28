package main

import (
	"fmt"
	"index/suffixarray"
	"reflect"
	"unsafe"
)

func main() {
	words := []string{"banana", "apple", "orange"}
	sc := NewSequnceConcater[byte](0)
	for _, word := range words {
		sc.Concat(int32(len(word)), func(i int32) byte { return word[i] })
	}
	sequence := sc.GetSequence()
	_, rank, _ := SuffixArray32(int32(len(sequence)), func(i int32) int32 { return int32(sequence[i]) })
	belong, start := sc.Build(int32(len(sequence)), func(i int32) int32 { return rank[i] })
	fmt.Println(belong, start)
}

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type SequnceConcater[S Int] struct {
	Sequence []S
	lid, rid []int32
	splitter S
}

func NewSequnceConcater[S Int](splitter S) *SequnceConcater[S] {
	return &SequnceConcater[S]{splitter: splitter}
}

func (sc *SequnceConcater[S]) Concat(n int32, f func(int32) S) {
	sc.Sequence = append(sc.Sequence, sc.splitter)
	sc.lid = append(sc.lid, int32(len(sc.Sequence)))
	for i := int32(0); i < n; i++ {
		sc.Sequence = append(sc.Sequence, f(i))
	}
	sc.rid = append(sc.rid, int32(len(sc.Sequence)))
}

func (sc *SequnceConcater[S]) GetSequence() []S {
	return sc.Sequence
}

// 根据后缀数组的排名数组 rank 构造 belong 和 start 数组.
// 返回值:
// belong: 每个后缀所属的字符串编号, -1 表示该后缀不属于任何字符串.
// start: 每个后缀在所属字符串中的起始位置, -1 表示该后缀不属于任何字符串.
func (sc *SequnceConcater[S]) Build(n int32, getRank func(int32) int32) (belong, start []int32) {
	belong = make([]int32, n)
	start = make([]int32, n)
	for i := range belong {
		belong[i] = -1
		start[i] = -1
	}
	for sid := int32(0); sid < int32(len(sc.lid)); sid++ {
		for j := sc.lid[sid]; j < sc.rid[sid]; j++ {
			p := getRank(j)
			belong[p] = sid
			start[p] = j - sc.lid[sid]
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
