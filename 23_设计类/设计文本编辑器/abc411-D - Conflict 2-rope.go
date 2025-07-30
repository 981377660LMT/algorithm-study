// abc411-D - Conflict 2-rope
// https://atcoder.jp/contests/abc411/tasks/abc411_d
// 有 1 台服务器和 N 台 PC，每台机器各维护一个字符串，初始均为空。
// 依次处理 Q 条指令：
//
// 1 p  把 PC p 的字符串设为当前服务器字符串
// 2 p s 在 PC p 的字符串末尾追加字符串 s
// 3 p  把服务器字符串设为当前 PC p 的字符串
// 在所有指令执行完后输出服务器字符串。
// 限制
// • 1 ≤ N,Q ≤ 2×10⁵
// • 所有追加字符串 s 的总长度 ≤ 10⁶
// 需要 O(Q+|总追加字符|) 解决。
//
// 思路
// 直接拷贝整串会导致 O(|串|) 复杂度，必须共享内存。
// 把每个字符串看成不可变 Rope.
// 所有修改都创建新节点或指针赋值，不触碰旧数据，确保 O(1) 更新

package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	left, right *Node
	str         string
	length      int
}

func leaf(s string) *Node { return &Node{str: s, length: len(s)} }
func concat(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return &Node{left: a, right: b, length: a.length + b.length}
}

func build(n *Node, buf *[]byte) {
	if n == nil {
		return
	}
	if n.left == nil && n.right == nil {
		*buf = append(*buf, n.str...)
		return
	}
	build(n.left, buf)
	build(n.right, buf)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, Q int
	fmt.Fscan(in, &N, &Q)

	pc := make([]*Node, N+1)
	var server *Node

	for ; Q > 0; Q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var p int
			fmt.Fscan(in, &p)
			pc[p] = server
		} else if t == 2 {
			var p int
			var s string
			fmt.Fscan(in, &p, &s)
			pc[p] = concat(pc[p], leaf(s))
		} else {
			var p int
			fmt.Fscan(in, &p)
			server = pc[p]
		}
	}

	buf := make([]byte, 0, 1<<20)
	build(server, &buf)
	fmt.Fprintln(out, string(buf))
}
