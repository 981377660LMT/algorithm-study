// !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/mo.go#L204

package moalgo

import (
	"fmt"
	"io"
	"math"
	"sort"
)

func normalMo(in io.Reader, nums []int, queries int) []int {
	n := len(nums)
	blockSize := int(math.Ceil(float64(n) / math.Sqrt(float64(queries))))
	type query struct{ lb, l, r, qid int }
	qs := make([]query, queries)
	for i := range qs {
		var l, r int
		fmt.Fscan(in, &l, &r) // 从 1 开始，[l,r)
		qs[i] = query{l / blockSize, l, r + 1, i}
	}
	sort.Slice(qs, func(i, j int) bool {
		a, b := qs[i], qs[j]
		if a.lb != b.lb {
			return a.lb < b.lb
		}
		// 奇偶化排序
		if a.lb&1 == 0 {
			return a.r < b.r
		}
		return a.r > b.r
	})

	cnt := 0
	l, r := 1, 1 // 区间从 1 开始，方便 debug
	move := func(idx, delta int) {
		// NOTE: 有些题目在 delta 为 1 和 -1 时逻辑的顺序是严格对称的
		// v := a[idx-1]
		// ...
		// cnt += delta
		if delta > 0 {
			cnt++
		} else {
			cnt--
		}
	}
	getAns := func(q query) int {
		// 提醒：q.r 是加一后的，计算时需要注意
		// sz := q.r - q.l
		// ...
		return cnt
	}
	ans := make([]int, queries)
	for _, q := range qs {
		for ; r < q.r; r++ {
			move(r, 1)
		}
		for ; l < q.l; l++ {
			move(l, -1)
		}
		for l > q.l {
			l--
			move(l, 1)
		}
		for r > q.r {
			r--
			move(r, -1)
		}
		ans[q.qid] = getAns(q)
	}
	return ans
}

// 带修莫队（支持单点修改）
// https://oi-wiki.org/misc/modifiable-mo-algo/
// https://codeforces.com/blog/entry/72690
// 模板题 数颜色 https://www.luogu.com.cn/problem/P1903
// https://codeforces.com/problemset/problem/940/F
// https://codeforces.com/problemset/problem/1476/G
// todo https://www.codechef.com/FEB17/problems/DISTNUM3
// todo 二逼平衡树（树套树）https://www.luogu.com.cn/problem/P3380
func moWithUpdate(in io.Reader) []int {
	var n, q int
	Fscan(in, &n, &q)
	a := make([]int, n+1) // 从 1 开始，方便 debug
	for i := 1; i <= n; i++ {
		Fscan(in, &a[i])
	}
	blockSize := int(math.Round(math.Pow(float64(n), 2.0/3)))
	type query struct{ lb, rb, l, r, t, qid int }
	type modify struct{ pos, val int }
	qs := []query{}
	ms := []modify{}
	for ; q > 0; q-- {
		var op string
		if Fscan(in, &op); op[0] == 'Q' {
			var l, r int
			Fscan(in, &l, &r)
			// 改成左闭右开
			qs = append(qs, query{l / blockSize, (r + 1) / blockSize, l, r + 1, len(ms), len(qs)})
		} else {
			var pos, val int
			Fscan(in, &pos, &val)
			ms = append(ms, modify{pos, val})
		}
	}
	sort.Slice(qs, func(i, j int) bool {
		a, b := qs[i], qs[j]
		if a.lb != b.lb {
			return a.lb < b.lb
		}
		if a.rb != b.rb {
			if a.lb&1 == 0 {
				return a.rb < b.rb
			}
			return a.rb > b.rb
		}
		if a.rb&1 == 0 {
			return a.t < b.t
		}
		return a.t > b.t
	})

	const mx int = 1e6 // TODO
	cnt, cc := [mx + 1]int{}, 0
	l, r, now := 1, 1, 0
	add := func(val int) {
		if cnt[val] == 0 {
			cc++
		}
		cnt[val]++
	}
	del := func(val int) {
		cnt[val]--
		if cnt[val] == 0 {
			cc--
		}
	}
	// 注：由于函数套函数不会内联，直接写到主流程的 for now 循环中会快不少
	timeSlip := func(l, r int) {
		m := ms[now]
		p, v := m.pos, m.val
		if l <= p && p < r {
			del(a[p])
			add(v)
		}
		a[p], ms[now].val = v, a[p]
	}
	getAns := func(q query) int {
		// 提醒：q.r 是加一后的，计算时需要注意
		// sz := q.r - q.l
		// ...
		return cc
	}
	ans := make([]int, len(qs))
	for _, q := range qs {
		for ; r < q.r; r++ {
			add(a[r])
		}
		for ; l < q.l; l++ {
			del(a[l])
		}
		for l > q.l {
			l--
			add(a[l])
		}
		for r > q.r {
			r--
			del(a[r])
		}
		for ; now < q.t; now++ {
			timeSlip(q.l, q.r)
		}
		for now > q.t {
			now--
			timeSlip(q.l, q.r)
		}
		ans[q.qid] = getAns(q)
	}
	return ans
}
