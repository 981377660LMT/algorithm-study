package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://atcoder.jp/contests/abc196/tasks/abc196_e
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	M := MonoidAddChminChmax{}
	f := M.e()
	var n int
	fmt.Fscan(in, &n)
	for i := 0; i < n; i++ {
		var v, kind int
		fmt.Fscan(in, &v, &kind)
		if kind == 1 {
			f = M.op(f, M.add(v))
		} else if kind == 2 {
			f = M.op(f, M.chmax(v))
		} else {
			f = M.op(f, M.chmin(v))
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var x int
		fmt.Fscan(in, &x)
		fmt.Fprintln(out, M.eval(f, x))
	}
}

const INF int = 1e18

// 合成 max(min(x+a,b),c)
type F = [3]int
type MonoidAddChminChmax struct{}

func (*MonoidAddChminChmax) e() F { return [3]int{0, INF, -INF} }
func (*MonoidAddChminChmax) op(f1, f2 F) F {
	a, b, c := f1[0], f1[1], f1[2]
	d, e, f := f2[0], f2[1], f2[2]
	a = a + d
	if b != INF {
		b += d
	}
	if c != -INF {
		c += d
	}
	b = min(b, e)
	c = max(min(c, e), f)
	return [3]int{a, b, c}
}
func (*MonoidAddChminChmax) add(a int) F   { return [3]int{a, INF, -INF} }
func (*MonoidAddChminChmax) chmin(b int) F { return [3]int{0, b, -INF} }
func (*MonoidAddChminChmax) chmax(c int) F { return [3]int{0, INF, c} }
func (*MonoidAddChminChmax) eval(f F, x int) int {
	a, b, c := f[0], f[1], f[2]
	return max(min(x+a, b), c)
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
