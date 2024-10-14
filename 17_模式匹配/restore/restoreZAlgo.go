package main

import (
	"bufio"
	"fmt"
	"os"
)

// D - LCP(prefix,suffix)
// https://atcoder.jp/contests/yahoo-procon2018-final-open/tasks/yahoo_procon2018_final_d
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int32
	fmt.Fscan(in, &N)
	A := make([]int32, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &A[i])
	}

	_, ok := RestoreZAlgo(N, func(i int32) int32 { return A[N-i-1] })
	if ok {
		fmt.Fprintln(out, "Yes")
	} else {
		fmt.Fprintln(out, "No")
	}
}

// 给定每个后缀 s[i:] 与 s 的最长公共前缀长度，返回字典序最小的原字符串.
// !f(0) = n.
func RestoreZAlgo(n int32, f func(i int32) int32) (res []int32, ok bool) {
	if n == 0 {
		ok = true
		return
	}
	if f(0) != n {
		return
	}
	uf := NewUnionFindArraySimple32(n)
	neq := make([][]int32, n)
	i, j := int32(1), int32(0)
	for i < n {
		for i+j < n {
			if f(i) != j {
				uf.Union(j, i+j, nil)
				j++
			} else {
				neq[j] = append(neq[j], i+j)
				neq[i+j] = append(neq[i+j], j)
				break
			}
		}
		if f(i) != j {
			return
		}
		if j == 0 {
			i++
			continue
		}
		k := int32(1)
		for i+k < n && k+f(k) < j {
			if f(i+k) != f(k) {
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
			for _, to := range neq[w] {
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
