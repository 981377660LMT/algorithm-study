package mowithupdate

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
)

// 带修莫队（支持单点修改）
// https://oi-wiki.org/misc/modifiable-mo-algo/
// https://codeforces.com/blog/entry/72690
// 模板题 数颜色 https://www.luogu.com.cn/problem/P1903
// https://codeforces.com/problemset/problem/940/F
// https://codeforces.com/problemset/problem/1476/G
// todo https://www.codechef.com/FEB17/problems/DISTNUM3
// todo 二逼平衡树（树套树）https://www.luogu.com.cn/problem/P3380
// !https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/mo.go#L204
func mowithUpdate(in io.Reader) []int {
	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n+1) // !从 1 开始，方便 debug
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	blockSize := int(math.Round(math.Pow(float64(n), 2.0/3)))
	type query struct{ leftBlock, rightBlock, left, right, time, qid int }
	type modify struct{ pos, val int }
	queries := []query{}
	modifies := []modify{}
	for ; q > 0; q-- {
		var op string
		if fmt.Fscan(in, &op); op[0] == 'Q' {
			var left, right int
			fmt.Fscan(in, &left, &right)
			// 改成左闭右开
			queries = append(queries, query{left / blockSize, (right + 1) / blockSize, left, right + 1, len(modifies), len(queries)})
		} else {
			var pos, val int
			fmt.Fscan(in, &pos, &val)
			modifies = append(modifies, modify{pos, val})
		}
	}
	sort.Slice(queries, func(i, j int) bool {
		a, b := queries[i], queries[j]
		if a.leftBlock != b.leftBlock {
			return a.leftBlock < b.leftBlock
		}
		if a.rightBlock != b.rightBlock {
			if a.leftBlock&1 == 0 {
				return a.rightBlock < b.rightBlock
			}
			return a.rightBlock > b.rightBlock
		}
		if a.rightBlock&1 == 0 {
			return a.time < b.time
		}
		return a.time > b.time
	})

	// !业务代码 (示例：统计区间内不同数字的个数)
	const N int = 1e6 // TODO
	counter, kind := [N + 1]int{}, 0
	left, right, now := 1, 1, 0
	add := func(val int) {
		if counter[val] == 0 {
			kind++
		}
		counter[val]++
	}
	remove := func(val int) {
		counter[val]--
		if counter[val] == 0 {
			kind--
		}
	}
	// 注：由于函数套函数不会内联，直接写到主流程的 for now 循环中会快不少
	timeSlip := func(left, right int) {
		m := modifies[now]
		p, v := m.pos, m.val
		if left <= p && p < right {
			remove(nums[p])
			add(v)
		}
		nums[p], modifies[now].val = v, nums[p]
	}
	getRes := func(q query) int {
		// 提醒：q.r 是加一后的，计算时需要注意
		// sz := q.r - q.l
		// ...
		return kind
	}

	// 主流程
	res := make([]int, len(queries))
	for _, q := range queries {
		for ; right < q.right; right++ {
			add(nums[right])
		}
		for ; left < q.left; left++ {
			remove(nums[left])
		}
		for left > q.left {
			left--
			add(nums[left])
		}
		for right > q.right {
			right--
			remove(nums[right])
		}
		for ; now < q.time; now++ {
			timeSlip(q.left, q.right)
		}
		for now > q.time {
			now--
			timeSlip(q.left, q.right)
		}

		res[q.qid] = getRes(q)
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	res := mowithUpdate(in)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}
