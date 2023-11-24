// https://www.luogu.com.cn/problem/P3380

package main

import (
	"bufio"
	. "fmt"
	"io"
	"os"
	"runtime/debug"
	"sort"
)

// https://space.bilibili.com/206214
func init() { debug.SetGCPercent(-1) }

const stNodeDefaultVal = 0

type stNode struct {
	lo, ro *stNode
	val    int32
}

var emptyNode = &stNode{}

func init() {
	emptyNode.lo = emptyNode
	emptyNode.ro = emptyNode
}

func (o *stNode) update(l, r, i int, val int32) {
	if l == r {
		o.val += val
		return
	}
	m := (l + r) >> 1
	if i <= m {
		if o.lo == emptyNode {
			o.lo = &stNode{emptyNode, emptyNode, stNodeDefaultVal}
		}
		o.lo.update(l, m, i, val)
	} else {
		if o.ro == emptyNode {
			o.ro = &stNode{emptyNode, emptyNode, stNodeDefaultVal}
		}
		o.ro.update(m+1, r, i, val)
	}
	o.val = o.lo.val + o.ro.val
}

func (o *stNode) query(l, r, R int) int32 {
	if o == emptyNode || R < l {
		return stNodeDefaultVal
	}
	if r <= R {
		return o.val
	}
	m := (l + r) >> 1
	return o.lo.query(l, m, R) + o.ro.query(m+1, r, R)
}

type fenwickWithSeg []*stNode

var mx int

// 二维单点更新：位置 (i,j) 用 val 更新
func (f fenwickWithSeg) update(i, j int, val int32) {
	for ; i < len(f); i += i & -i {
		f[i].update(0, mx, j, val)
	}
}

// 二维前缀和：累加所有 x <= i 且 y <= j 的值
func (f fenwickWithSeg) pre(i, j int) (res int32) {
	for ; i > 0; i &= i - 1 {
		res += f[i].query(0, mx, j)
	}
	return
}

func (f fenwickWithSeg) rank(l, r, k int) int32 {
	return f.pre(r, k-1) - f.pre(l-1, k-1) + 1
}

var ar, al []*stNode

func (f fenwickWithSeg) kth(i, j int, k int32) int {
	ar = ar[:0]
	al = al[:0]
	for ; j > 0; j &= j - 1 {
		ar = append(ar, f[j])
	}
	for i--; i > 0; i &= i - 1 {
		al = append(al, f[i])
	}
	l, r := 0, mx
	for l < r {
		s := int32(0)
		for _, o := range ar {
			s += o.lo.val
		}
		for _, o := range al {
			s -= o.lo.val
		}
		m := (l + r) >> 1
		if s >= k {
			for i, o := range ar {
				ar[i] = o.lo
			}
			for i, o := range al {
				al[i] = o.lo
			}
			r = m
		} else {
			k -= s
			for i, o := range ar {
				ar[i] = o.ro
			}
			for i, o := range al {
				al[i] = o.ro
			}
			l = m + 1
		}
	}
	return l
}

func run(_r io.Reader, _w io.Writer) {
	in := bufio.NewReader(_r)
	out := bufio.NewWriter(_w)
	defer out.Flush()

	var n, m int
	Fscan(in, &n, &m)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		Fscan(in, &a[i])
	}
	qs := make([]struct{ op, l, r, k int }, m)
	b := make([]int, 0, n+m)
	b = append(b, a[1:]...)
	for i := range qs {
		Fscan(in, &qs[i].op)
		if qs[i].op == 3 {
			Fscan(in, &qs[i].l, &qs[i].k)
		} else {
			Fscan(in, &qs[i].l, &qs[i].r, &qs[i].k)
		}
		if qs[i].op != 2 {
			b = append(b, qs[i].k)
		}
	}
	sort.Ints(b)
	i := 1
	for k := 1; k < len(b); k++ {
		if b[k] != b[k-1] {
			if i != k {
				b[i] = b[k]
			}
			i++
		}
	}
	b = b[:i]
	mx = i

	t := make(fenwickWithSeg, n+1)
	for i := range t {
		t[i] = &stNode{emptyNode, emptyNode, stNodeDefaultVal}
	}
	for i := 1; i <= n; i++ {
		a[i] = sort.SearchInts(b, a[i])
		t.update(i, a[i], 1)
	}
	for _, q := range qs {
		op, l, r, k := q.op, q.l, q.r, q.k
		if op != 2 {
			k = sort.SearchInts(b, k)
		}
		if op == 3 {
			t.update(l, a[l], -1)
			a[l] = k
			t.update(l, k, 1)
			continue
		}
		switch op {
		case 1:
			Fprintln(out, t.rank(l, r, k))
		case 2:
			Fprintln(out, b[t.kth(l, r, int32(k))])
		case 4:
			rk := t.rank(l, r, k)
			if rk == 1 {
				Fprintln(out, "-2147483647")
			} else {
				Fprintln(out, b[t.kth(l, r, rk-1)])
			}
		default:
			rk := t.rank(l, r, k+1)
			if rk > int32(r-l+1) {
				Fprintln(out, "2147483647")
			} else {
				Fprintln(out, b[t.kth(l, r, rk)])
			}
		}
	}
}

func main() { run(os.Stdin, os.Stdout) }
