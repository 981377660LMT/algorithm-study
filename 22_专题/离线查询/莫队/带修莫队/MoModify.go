// 带有时间序列的莫队(带修莫队),时间复杂度O(n^5/3)
// https://maspypy.github.io/library/ds/offline_query/mo_3d.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/mo.go
// 带修莫队（支持单点修改）
// https://oi-wiki.org/misc/modifiable-mo-algo/
// https://codeforces.com/blog/entry/72690
// 模板题 数颜色 https://www.luogu.com.cn/problem/P1903
// https://codeforces.com/problemset/problem/940/F

// 普通莫队是不能带修改的。
// 我们可以强行让它可以修改，就像 DP 一样，可以强行加上一维 时间维, 表示这次操作的时间。
// 时间维表示经历的修改次数。
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	Cf940F()
	Luo1903()
}

// F. Machine Learning
// https://codeforces.com/problemset/problem/940/F
// 1 start end 询问区间[start,end)中每个数字出现次数的mex
// 2 pos val 修改第pos个数字为val
func Cf940F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1], &queries[i][2])
		queries[i][1]--
	}

	D := NewDictionary()
	for i := range nums {
		nums[i] = D.Id(nums[i])
	}
	for i := range queries {
		if queries[i][0] == 2 {
			queries[i][2] = D.Id(queries[i][2])
		}
	}

	tmpNums := append([]int{}, nums...)
	queryCount := 0
	M := NewMoModify()
	change := make([][3]int, q) // [pos, old, new]
	for i := 0; i < q; i++ {
		item := queries[i]
		if item[0] == 1 {
			M.AddQuery(len(change), item[1], item[2])
			queryCount++
		} else {
			pos, newValue := item[1], item[2]
			change = append(change, [3]int{pos, tmpNums[pos], newValue})
			tmpNums[pos] = newValue
		}
	}

	valueCounter := make([]int, D.Size())
	freqCounter := make([]int, n+1)
	freqCounter[0] = 1 << 30
	res := make([]int, queryCount)

	update := func(value int, count int) {
		c := valueCounter[value]
		freqCounter[c]--
		valueCounter[value] += count
		freqCounter[c+count]++
	}

	addLeft := func(index int) { update(nums[index], 1) }
	addRight := func(index int) { update(nums[index], 1) }
	removeLeft := func(index int) { update(nums[index], -1) }
	removeRight := func(index int) { update(nums[index], -1) }
	addChange := func(time int, start, end int) {
		item := change[time]
		pos, old, new := item[0], item[1], item[2]
		if start <= pos && pos < end {
			update(old, -1)
			update(new, 1)
		}
		nums[pos] = new
	}
	removeChange := func(time int, start, end int) {
		item := change[time]
		pos, old, new := item[0], item[1], item[2]
		if start <= pos && pos < end {
			update(new, -1)
			update(old, 1)
		}
		nums[pos] = old
	}
	query := func(qid int) {
		mex := 0
		for freqCounter[mex] > 0 {
			mex++
		}
		res[qid] = mex
	}

	M.Run(addLeft, addRight, removeLeft, removeRight, addChange, removeChange, query, -1)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

func Luo1903() {
	// https://www.luogu.com.cn/problem/P1903
	// Q L R 查询第L支画笔到第R支画笔中共有几种不同颜色的画笔。
	// R P Col 把第P支画笔替换为颜色 Col

	// n,q<=1e5
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	D := NewDictionary()
	for i, v := range nums {
		nums[i] = D.Id(v)
	}

	tmpNums := append([]int{}, nums...)
	queryCount := 0
	M := NewMoModify()
	change := make([][3]int, q) // [pos, old, new]
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "Q" {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			M.AddQuery(len(change), start, end) // 时间点, 左闭右开区间
			queryCount++
		} else {
			var pos, newValue int
			fmt.Fscan(in, &pos, &newValue)
			pos--
			newValue = D.Id(newValue)
			change = append(change, [3]int{pos, tmpNums[pos], newValue})
			tmpNums[pos] = newValue
		}
	}

	res := make([]int, queryCount)
	counter := make([]int, D.Size())
	kind := 0

	add := func(value int) {
		if counter[value] == 0 {
			kind++
		}
		counter[value]++
	}
	remove := func(value int) {
		counter[value]--
		if counter[value] == 0 {
			kind--
		}
	}

	addLeft := func(index int) { add(nums[index]) }
	addRight := func(index int) { add(nums[index]) }
	removeLeft := func(index int) { remove(nums[index]) }
	removeRight := func(index int) { remove(nums[index]) }
	addChange := func(time int, start, end int) {
		item := change[time]
		pos, old, new := item[0], item[1], item[2]
		if start <= pos && pos < end {
			remove(old)
			add(new)
		}
		nums[pos] = new
	}
	removeChange := func(time int, start, end int) {
		item := change[time]
		pos, old, new := item[0], item[1], item[2]
		if start <= pos && pos < end {
			remove(new)
			add(old)
		}
		nums[pos] = old
	}
	query := func(qid int) { res[qid] = kind }
	M.Run(addLeft, addRight, removeLeft, removeRight, addChange, removeChange, query, -1)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 支持单点修改的莫队算法.时间复杂度O(n^5/3).
type MoModify struct {
	query [][3]int // [time, start, end]
}

func NewMoModify() *MoModify {
	return &MoModify{}
}

// 添加一个查询，查询范围为`左闭右开区间` [start, end).
//
//	0 <= start <= end <= n.
//	time: 查询的时间点, 即当前修改的次数`len(modify)`.
func (mm *MoModify) AddQuery(time int, start, end int) {
	mm.query = append(mm.query, [3]int{time, start, end})
}

// 返回每个查询的结果.
//
// addLeft: 将数据添加到窗口左侧.
// addRight: 将数据添加到窗口右侧.
// removeLeft: 将数据从窗口左侧移除.
// removeRight: 将数据从窗口右侧移除.
// addChange: 添加修改. time: 修改的时间点. start, end: 当前窗口的范围.
// removeChange: 移除修改. time: 修改的时间点. start, end: 当前窗口的范围.
// query: 查询窗口内的数据.
// blockSize: 分块大小.-1表示自动计算.
func (mm *MoModify) Run(
	addLeft func(index int),
	addRight func(index int),
	removeLeft func(index int),
	removeRight func(index int),
	addChange func(time int, start, end int),
	removeChange func(time int, start, end int),
	query func(qid int),
	blockSize int,
) {
	if blockSize == -1 {
		q := max(1, len(mm.query))
		blockSize = 1
		for blockSize*blockSize*blockSize < q*q {
			blockSize++
		}
	}
	order := mm._getOrder(blockSize)
	t, l, r := 0, 0, 0
	for _, qid := range order {
		item := mm.query[qid]
		nt, nl, nr := item[0], item[1], item[2]
		for l > nl {
			l--
			addLeft(l)
		}
		for r < nr {
			addRight(r)
			r++
		}
		for l < nl {
			removeLeft(l)
			l++
		}
		for r > nr {
			r--
			removeRight(r)
		}
		for t < nt {
			addChange(t, l, r)
			t++
		}
		for t > nt {
			t--
			removeChange(t, l, r)
		}
		query(qid)
	}
}

func (mm *MoModify) _getOrder(blockSize int) []int {
	k := 1 << 20
	q := len(mm.query)
	key := make([]int, q)
	for i, item := range mm.query {
		t := item[0] / blockSize
		l := item[1] / blockSize
		x := item[2]
		if l&1 == 1 {
			x = -x
		}
		x += l * k
		if t&1 == 1 {
			x = -x
		}
		x += t * k * k
		key[i] = x
	}

	order := make([]int, q)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return key[order[i]] < key[order[j]] })

	cost := func(a, b int) int {
		q1 := mm.query[order[a]]
		q2 := mm.query[order[b]]
		return abs(q1[0]-q2[0]) + abs(q1[1]-q2[1]) + abs(q1[2]-q2[2])
	}
	for k := 0; k < q-5; k++ {
		if cost(k, k+2)+cost(k+1, k+3) < cost(k, k+1)+cost(k+2, k+3) {
			order[k+1], order[k+2] = order[k+2], order[k+1]
		}
		if cost(k, k+3)+cost(k+1, k+4) < cost(k, k+1)+cost(k+3, k+4) {
			order[k+1], order[k+3] = order[k+3], order[k+1]
		}
	}

	return order
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
}
