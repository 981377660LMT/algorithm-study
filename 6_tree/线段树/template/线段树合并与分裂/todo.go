

const stNodeDefaultTodoVal = 0

var lazyRoot = &lazyNode{l: 1, r: 1e9, sum: stNodeDefaultVal}

type lazyNode struct {
	lo, ro *lazyNode
	l, r   int
	sum    int
	todo   int
}

func (o *lazyNode) get() int {
	if o != nil {
		return o.sum
	}
	return stNodeDefaultVal
}

func (lazyNode) op(a, b int) int {
	return a + b // max(a, b)
}

func (o *lazyNode) maintain() {
	o.sum = o.op(o.lo.get(), o.ro.get())
}

func (o *lazyNode) build(a []int, l, r int) {
	o.l, o.r = l, r
	o.todo = stNodeDefaultTodoVal
	if l == r {
		o.sum = a[l-1]
		return
	}
	m := (l + r) >> 1
	o.lo = &lazyNode{}
	o.lo.build(a, l, m)
	o.ro = &lazyNode{}
	o.ro.build(a, m+1, r)
	o.maintain()
}

func (o *lazyNode) do(add int) {
	o.todo += add                  // % mod
	o.sum += (o.r - o.l + 1) * add // % mod
}

func (o *lazyNode) spread() {
	m := (o.l + o.r) >> 1
	if o.lo == nil {
		o.lo = &lazyNode{l: o.l, r: m, sum: stNodeDefaultVal}
	}
	if o.ro == nil {
		o.ro = &lazyNode{l: m + 1, r: o.r, sum: stNodeDefaultVal}
	}
	if todo := o.todo; todo != stNodeDefaultTodoVal {
		o.lo.do(todo)
		o.ro.do(todo)
		o.todo = stNodeDefaultTodoVal
	}
}

func (o *lazyNode) update(l, r int, add int) {
	if l <= o.l && o.r <= r {
		o.do(add)
		return
	}
	o.spread()
	m := (o.l + o.r) >> 1
	if l <= m {
		o.lo.update(l, r, add)
	}
	if m < r {
		o.ro.update(l, r, add)
	}
	o.maintain()
}

func (o *lazyNode) query(l, r int) int {
	if o == nil || l > o.r || r < o.l {
		return stNodeDefaultVal
	}
	if l <= o.l && o.r <= r {
		return o.sum
	}
	o.spread()
	return o.op(o.lo.query(l, r), o.ro.query(l, r))
}

// EXTRA: 线段树合并
// https://www.luogu.com.cn/problem/P5494
// todo 一些题目 https://www.luogu.com.cn/blog/styx-ferryman/xian-duan-shu-ge-bing-zong-ru-men-dao-fang-qi
//   https://codeforces.com/blog/entry/83969
//   https://www.luogu.com.cn/problem/P4556
//   https://www.luogu.com.cn/problem/P5298
//   https://codeforces.com/problemset/problem/600/E
// rt = rt.merge(rt2)

// EXTRA: 线段树合并
// https://www.luogu.com.cn/problem/P5494
// todo 一些题目 https://www.luogu.com.cn/blog/styx-ferryman/xian-duan-shu-ge-bing-zong-ru-men-dao-fang-qi
//
//	 https://zhuanlan.zhihu.com/p/575513452
//		https://codeforces.com/problemset/problem/600/E
//
// rt = rt.merge(rt2)
func (o *stNode) merge(b *stNode) *stNode {
	if o == nil {
		return b
	}
	if b == nil {
		return o
	}
	if o.l == o.r {
		// 按照所需合并，如加法
		o.val += b.val
		return o
	}
	o.lo = o.lo.merge(b.lo)
	o.ro = o.ro.merge(b.ro)
	o.maintain()
	return o
}

// EXTRA: 线段树分裂
// 将区间 [l,r] 从 o 中分离到 b 上
// https://www.luogu.com.cn/blog/cyffff/talk-about-segument-trees-split
// https://www.luogu.com.cn/problem/P5494
// rt, rt2 := rt.split(nil, l, r)
func (o *stNode) split(b *stNode, l, r int) (*stNode, *stNode) {
	if o == emptyStNode || l > o.r || r < o.l {
		return o, emptyStNode
	}
	if l <= o.l && o.r <= r {
		return emptyStNode, o
	}
	if b == emptyStNode {
		b = &stNode{lo: emptyStNode, ro: emptyStNode, l: o.l, r: o.r, val: stNodeDefaultVal}
	}
	o.lo, b.lo = o.lo.split(b.lo, l, r)
	o.ro, b.ro = o.ro.split(b.ro, l, r)
	o.maintain()
	b.maintain()
	return o, b
}

// 权值线段树求第 k 小
// 调用前需保证 1 <= k <= root.val
func (o *stNode) kth(k int) int {
	if o.l == o.r {
		return o.l
	}
	cntL := o.lo.val
	if k <= cntL {
		return o.lo.kth(k)
	}
	return o.ro.kth(k - cntL)
}
