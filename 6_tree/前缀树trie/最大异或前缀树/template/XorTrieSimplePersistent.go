// TODO: 优化写法.补充remove.补充enumerate方法.以及若何扩展Node属性.

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

// 注：由于用的是指针写法，必要时禁止 GC，能加速不少
func init() { debug.SetGCPercent(-1) }

type Trie01Node struct {
	lastIndex int
	chidlren  [2]*Trie01Node
}

var BITLEN int = 24
var preXor []int

// 持久化01字典树
// 注意为了拷贝一份 trie01Node，这里的接收器不是指针
// usage:
//
//	roots := make([]*trie01Node, n+1)
//	roots[0] = trie01Node{}.put(0, trieBitLen-1)
//	roots[i+1] = roots[i].put(v, trieBitLen-1)
//
// !https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/trie01.go#L19
// 在第k位插入一个数，返回新的根节点
func (o Trie01Node) Insert(index, value, k int) *Trie01Node {
	o.lastIndex = index
	if k < 0 {
		return &o
	}
	bit := (value >> k) & 1
	if o.chidlren[bit] == nil {
		o.chidlren[bit] = &Trie01Node{}
	}
	o.chidlren[bit] = o.chidlren[bit].Insert(index, value, k-1)
	return &o
}

func (o *Trie01Node) Query(value, lowerIndex int) *Trie01Node {
	for k := BITLEN - 1; k >= 0; k-- {
		bit := (value >> k) & 1
		if o.chidlren[bit^1] != nil && o.chidlren[bit^1].lastIndex >= lowerIndex { // guard
			bit ^= 1
		}
		o = o.chidlren[bit]
	}
	return o
}

// 最大异或和
// https://www.luogu.com.cn/problem/P4735
// 给定一个非负整数序列nums，初始长度为n。
// 有q个操作，有以下两种操作类型:
//  - A x: 添加操作，表示在序列末尾添加一个数x，序列的长度+1。
//  - Q left right x: 询问操作，你需要找到一个位置p，
//    !满足left≤p≤right，使得x与后缀的异或和 x^(nums[pos]^nums[pos+1]^...^nums[n]) 最大，输出最大是多少。
// !n,m<=3e5 0<=x<=1e7

// !解法
// 维护前缀异或 查询变为 preXor[pos]^preXor[n]^x
// 即求 前缀异或与(x^preXor[n])的最大值
// 用持久化trie01来维护前缀异或，第i个版本为插入了nums[i]后的trie树
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	maxVersion := n + q + 1
	preXor = make([]int, maxVersion)
	roots := make([]*Trie01Node, maxVersion)
	roots[0] = Trie01Node{}.Insert(0, 0, BITLEN-1)
	for i := 0; i < n; i++ {
		var num int
		fmt.Fscan(in, &num)
		preXor[i+1] = preXor[i] ^ num
		roots[i+1] = roots[i].Insert(i+1, preXor[i+1], BITLEN-1)
	}

	for i, cur := 0, n; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "A" {
			var num int
			fmt.Fscan(in, &num)
			preXor[cur+1] = preXor[cur] ^ num
			roots[cur+1] = roots[cur].Insert(cur+1, preXor[cur+1], BITLEN-1)
			cur++
		} else {
			var left, right, x int
			fmt.Fscan(in, &left, &right, &x)
			left, right = left-1, right-1
			x ^= preXor[cur]
			node := roots[right].Query(x, left)
			fmt.Fprintln(out, preXor[node.lastIndex]^x)
		}
	}
}
