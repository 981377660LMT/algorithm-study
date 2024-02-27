// (部分)可持久化KMP
// 记字符集大小为m, 模式串长度为n, 则构造复杂度为O(nm).
// 每一步对于 i 的转移我们都只会复制 fail[i]的状态并进行一次单点修改, 时间复杂度为O(m).
// 字符集较小时，直接拷贝数组修改；
// 字符集较大时，需要可持久化线段树等支持复制历史版本和单点修改的数据结构(例如可持久化数组)。

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF1721E()
}

// Prefix Function Queries
// https://www.luogu.com.cn/problem/CF1721E
// 给定字符串s以及q个串ti，求将s分别与每个ti拼接起来后，最靠右的|ti|个前缀的 border 长度。询问间相互独立。。
//
// 由于 KMP 是基于均摊的，所以显然不能每次询问暴力跑一遍 KMP.
// 考虑优化询问时 KMP 跳 next 的过程：
// 预处理时记录每种状态后面加每种字符的 next，其实就是单串的 AC 自动机，
// 当询问 KMP 时跳到原串部分后，直接返回结果。
func CF1721E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)

	K := NewKmpPersistent(int32(len(s)), func(i int32) byte { return s[i] })
	sVersion := K.CurVersion
	fmt.Println(K.Next)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var t string
		fmt.Fscan(in, &t)
		curVersion := sVersion
		for i := 0; i < len(t); i++ {
			curVersion = K.Append(curVersion, t[i])
			fmt.Fprint(out, K.Query(curVersion), " ")
		}
		fmt.Fprintln(out)
	}
}

const SIGMA byte = 26
const OFFSET byte = 97

type KmpPersistent struct {
	CurVersion int32
	Fail       []int32
	Next       [][SIGMA]int32
}

func NewKmpPersistent(n int32, f func(i int32) byte) *KmpPersistent {
	fail := make([]int32, n)
	j := int32(0)
	for i := int32(1); i < n; i++ {
		for j > 0 && f(i) != f(j) {
			j = fail[j-1]
		}
		if f(i) == f(j) {
			j++
		}
		fail[i] = j
	}

	next := make([][SIGMA]int32, n)
	for i := int32(0); i < n; i++ {
		if i > 0 {
			next[i] = next[fail[i-1]]
		}
		next[i][f(i)-OFFSET] = i + 1
	}

	return &KmpPersistent{CurVersion: 0, Fail: fail, Next: next}
}

// 在当前版本字符串末尾追加一个字符.
func (kmp *KmpPersistent) Append(version int32, c byte) int32 {

}

// 查询当前版本字符串的border长度.
func (kmp *KmpPersistent) Query(version int32) int32 {
	return kmp.Fail[kmp.CurVersion+1]
}
