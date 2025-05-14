// E - Forbidden Prefix
// https://atcoder.jp/contests/abc403/tasks/abc403_e
// 一开始有两个空的多重集合X和Y。
// 有Q次操作，操作按顺序进行，第i次操作会给定一个整数Ti和一个字符串Si：
// 如果Ti=1，将Si加入到集合X内如果Ti=2，将Si加入到集合Y内在每次操作结束后，请输出一个整数，表示：
// Y内总共有多少个字符串，满足X中不存在任意一个字符串是它的前缀。

package main

import (
	"bufio"
	"fmt"
	"os"
)

type TrieNode struct {
	children [26]*TrieNode
	isEnd1   bool
	prefix2  []int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	res := 0
	id := 0

	root := &TrieNode{}

	deleted2Set := make(map[int]struct{})
	delete2 := func(s int) {
		if _, d := deleted2Set[s]; d {
			return
		}
		deleted2Set[s] = struct{}{}
		res--
	}

	insert1 := func(s string) {
		id++

		cur := root
		for _, c := range s {
			b := c - 'a'
			if cur.children[b] == nil {
				cur.children[b] = &TrieNode{}
			}
			cur = cur.children[b]
		}

		cur.isEnd1 = true
		for _, id := range cur.prefix2 {
			delete2(id)
		}
		cur.prefix2 = nil
	}

	insert2 := func(s string) {
		id++
		res++
		cur := root
		for _, c := range s {
			b := c - 'a'
			if cur.children[b] == nil {
				cur.children[b] = &TrieNode{}
			}
			cur = cur.children[b]
			if cur.isEnd1 {
				delete2(id)
				return
			}
			cur.prefix2 = append(cur.prefix2, id)
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op int
		var s string
		fmt.Fscan(in, &op, &s)
		if op == 1 {
			insert1(s)
		} else {
			insert2(s)
		}

		fmt.Fprintln(out, res)
	}
}
