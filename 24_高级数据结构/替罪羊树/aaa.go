package main

import (
	"bufio"
	. "fmt"
	"io"
	"os"
)

// github.com/EndlessCheng/codeforces-go
type node struct {
	lr  [2]*node
	key int
	cnt int
	sz  int
	msz int
}

func (o *node) cmp(b int) int {
	switch {
	case b < o.key:
		return 0 // 左儿子
	case b > o.key:
		return 1 // 右儿子
	default:
		return -1
	}
}

func (o *node) size() int {
	if o != nil {
		return o.sz
	}
	return 0
}

func (o *node) mSize() int {
	if o != nil {
		return o.msz
	}
	return 0
}

func (o *node) maintain() {
	o.sz = o.lr[0].size() + o.lr[1].size()
	if o.cnt > 0 {
		o.sz++
	}
	o.msz = o.lr[0].mSize() + o.lr[1].mSize() + o.cnt
}

func (o *node) nodes() []*node {
	nodes := make([]*node, 0, o.size())
	var f func(*node)
	f = func(o *node) {
		if o == nil {
			return
		}
		f(o.lr[0])
		if o.cnt > 0 {
			nodes = append(nodes, o)
		}
		f(o.lr[1])
	}
	f(o)
	return nodes
}

func buildSGT(nodes []*node) *node {
	if len(nodes) == 0 {
		return nil
	}
	m := len(nodes) / 2
	o := nodes[m]
	o.lr[0] = buildSGT(nodes[:m])
	o.lr[1] = buildSGT(nodes[m+1:])
	o.maintain()
	return o
}

func (o *node) rebuild() *node { return buildSGT(o.nodes()) }

type scapegoatTree struct {
	root   *node
	delCnt int
	foo    *node // 最后个需要重构的节点
}

func (t *scapegoatTree) _put(o *node, key int) *node {
	if o == nil {
		return &node{key: key, cnt: 1, sz: 1, msz: 1}
	}
	if d := o.cmp(key); d >= 0 {
		o.lr[d] = t._put(o.lr[d], key)
	} else {
		if o.cnt == 0 {
			t.delCnt--
		}
		o.cnt++
	}
	o.maintain()
	if 4*max(o.lr[0].size(), o.lr[1].size()) > 3*o.size() { // alpha=3/4
		t.foo = o
	}
	return o
}

func (t *scapegoatTree) put(key int) {
	t.root = t._put(t.root, key)

}

func (t *scapegoatTree) _delete(o *node, key int) {
	if o == nil {
		return
	}
	if d := o.cmp(key); d >= 0 {
		t._delete(o.lr[d], key)
	} else if o.cnt > 0 {
		o.cnt--
		t.delCnt++
	}
	o.maintain()
}

func (t *scapegoatTree) delete(key int) {
	t._delete(t.root, key)
	if t.delCnt > t.root.size()/2 {
		t.root = t.root.rebuild()
		t.delCnt = 0
	}
}

func (t *scapegoatTree) rank(key int) (kth int, o *node) {
	for o = t.root; o != nil; {
		switch c := o.cmp(key); {
		case c == 0:
			o = o.lr[0]
		case c > 0:
			kth += o.cnt + o.lr[0].mSize()
			o = o.lr[1]
		default:
			kth += o.lr[0].mSize()
			return
		}
	}
	return
}

func (t *scapegoatTree) kth(k int) (o *node) {
	for o = t.root; o != nil; {
		switch lsz := o.lr[0].mSize(); {
		case k < lsz:
			o = o.lr[0]
		default:
			k -= o.cnt + lsz
			if k < 0 {
				return
			}
			o = o.lr[1]
		}
	}
	return
}

func (t *scapegoatTree) prev(key int) (prev *node) {
	rk, _ := t.rank(key)
	return t.kth(rk - 1)
}

func (t *scapegoatTree) next(key int) (next *node) {
	rk, o := t.rank(key)
	if o != nil {
		rk += o.cnt
	}
	return t.kth(rk)
}

func run(_r io.Reader, _w io.Writer) {
	out := bufio.NewWriter(_w)
	defer out.Flush()
	buf := make([]byte, 4096)
	_i := len(buf)
	rc := func() byte {
		if _i == len(buf) {
			_r.Read(buf)
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}
	r := func() (x int) {
		b := rc()
		for ; '0' > b; b = rc() {
		}
		for ; '0' <= b; b = rc() {
			x = x*10 + int(b&15)
		}
		return
	}

	t := &scapegoatTree{}
	n, q := r(), r()
	for ; n > 0; n-- {
		t.put(r())
	}
	ans := 0
	for last := 0; q > 0; q-- {
		op := rc()
		for ; '0' > op; op = rc() {
		}
		v := r() ^ last
		switch op {
		case '1':
			t.put(v)
		case '2':
			t.delete(v)
		case '3':
			last, _ = t.rank(v)
			last++
			ans ^= last
		case '4':
			last = t.kth(v - 1).key
			ans ^= last
		case '5':
			last = t.prev(v).key
			ans ^= last
		default:
			last = t.next(v).key
			ans ^= last
		}
	}
	Fprintln(out, ans)
}

func main() { run(os.Stdin, os.Stdout) }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
