// 动态z函数/在线z函数

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	Z := NewZAlgoOnline[byte]()
	for i := 0; i < n; i++ {
		Z.Append(s[i])
	}

	for i := 0; i < n; i++ {
		fmt.Fprint(out, Z.Query(int32(i)), " ")
	}
}

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type ZAlgoOnline[S Int] struct {
	sb   []S
	z    []int32
	memo [][]int32
	pos  int32
}

func NewZAlgoOnline[S Int]() *ZAlgoOnline[S] {
	return &ZAlgoOnline[S]{pos: 1}
}

// 末尾添加s[i]，返回被更新的位置j，满足s[j:i]是lcp.
func (za *ZAlgoOnline[S]) Append(c S) []int32 {
	last := int32(len(za.sb))
	za.sb = append(za.sb, c)
	len_ := int32(len(za.sb))
	za.z = append(za.z, -1)
	za.memo = append(za.memo, []int32{})
	end := []int32{}
	if len_ == 1 {
		return end
	}
	if za.sb[0] != c {
		za.z[last] = 0
		end = append(end, last)
	}
	delete := func(j int32) {
		za.z[j] = last - j
		za.memo[last] = append(za.memo[last], j)
		end = append(end, j)
	}
	for za.pos <= last {
		if za.z[za.pos] != -1 {
			za.pos++
		} else if za.sb[last-za.pos] != c {
			delete(za.pos)
			za.pos++
		} else {
			break
		}
	}
	if za.pos < last {
		for _, j := range za.memo[last-za.pos] {
			delete(j + za.pos)
		}
	}
	return end
}

// 查询s[i:]与s的最长公共前缀长度.
func (za *ZAlgoOnline[S]) Query(i int32) int32 {
	if za.z[i] == -1 {
		return int32(len(za.sb)) - i
	}
	return za.z[i]
}
