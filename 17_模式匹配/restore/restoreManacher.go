package main

import (
	"bufio"
	"fmt"
	"os"
)

// G - Palindrome Construction
// https://atcoder.jp/contests/abc349/tasks/abc349_g
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)

	r := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &r[i])
		r[i] = 2*r[i] + 1
	}

	res, ok := RestoreManacher(n, func(i int32) int32 { return r[i] })
	if !ok {
		fmt.Fprintln(out, "No")
		return
	}

	fmt.Fprintln(out, "Yes")
	for i := int32(0); i < n; i++ {
		fmt.Fprint(out, res[i]+1, " ")
	}
}

// 给定以每个字符为中心的最长回文串长度，返回字典序最小的原字符串.
func RestoreManacher(n int32, f func(i int32) int32) (res []int32, ok bool) {
	for i := int32(0); i < n; i++ {
		if f(i)&1 == 0 {
			panic("f(i) must be odd")
		}
	}
	g := func(i int32) int32 {
		return (f(i) + 1) >> 1
	}

	uf := NewUnionFindArraySimple32(n)
	diff := make([][]int32, n)
	i, j := int32(0), int32(0)
	for i < n {
		for i >= j && i+j < n {
			if g(i) != j {
				if j > 0 {
					uf.Union(i+j, i-j, nil)
				}
				j++
			} else {
				diff[i+j] = append(diff[i+j], i-j)
				diff[i-j] = append(diff[i-j], i+j)
				break
			}
		}
		if g(i) != j {
			return
		}
		k := int32(1)
		for i >= k && i+k < n && k+g(i-k) < j {
			if g(i+k) != g(i-k) {
				return
			}
			k++
		}
		i += k
		j -= k
	}

	vs := make([][]int32, n)
	for v := int32(0); v < n; v++ {
		root := uf.Find(v)
		vs[root] = append(vs[root], v)
	}
	res = make([]int32, n)
	for i := int32(0); i < n; i++ {
		res[i] = -1
	}
	for i := int32(0); i < n; i++ {
		root := uf.Find(i)
		if res[root] != -1 {
			continue
		}
		var tmp []int32
		for _, w := range vs[root] {
			for _, to := range diff[w] {
				if res[to] != -1 {
					tmp = append(tmp, res[to])
				}
			}
		}
		x := mex(int32(len(tmp)), func(i int32) int32 { return tmp[i] })
		for _, w := range vs[root] {
			res[w] = x
		}
	}

	ok = true
	return
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32, beforeMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if beforeMerge != nil {
		beforeMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	root := key
	for u.data[root] >= 0 {
		root = u.data[root]
	}
	for key != root {
		key, u.data[key] = u.data[key], root
	}
	return root
}

func (u *UnionFindArraySimple32) Size(key int32) int32 {
	return -u.data[u.Find(key)]
}

// 最小的不在集合中的非负整数.
func mex(n int32, f func(i int32) int32) int32 {
	aru := make([]bool, n+1)
	for i := int32(0); i < n; i++ {
		if v := f(i); v < n {
			aru[v] = true
		}
	}
	mex := int32(0)
	for aru[mex] {
		mex++
	}
	return mex
}
