// AGC002D Stamp Rally-操作分块在线
// https://hoikoro.hatenablog.com/entry/2017/12/14/040750

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		edges[i] = [2]int{u, v}
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		var x, y, z int
		fmt.Fscan(in, &x, &y, &z)
		x--
		y--
		queries[i] = [3]int{x, y, z}
	}

	res := StampRally(n, edges, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://atcoder.jp/contests/agc002/tasks/agc002_d
// 一张连通图，q 次询问从两个点 x 和 y 出发，
// 希望经过的点数量等于 z（每个点可以重复经过，但是重复经过只计算一次）
// 求经过的边最大编号最小是多少。
//
// 很容易想到离线加边，判断每一个询问是否满足条件，但必须每次检查每一个询问，时间复杂度O(mqlogn)
// 上述过程每一次检验非常浪费时间，考虑分块，每插入O(sqrt(m))次检验一次，
// 对于每一个询问倒着扫一遍,恢复时用undo。
// 时间复杂度O(qsqrt(m)logn).
// 这种方法是整体二分的下位版本，但是实现起来比较简单。
func StampRally(n int, edges [][2]int, queries [][3]int) []int {
	m, q := len(edges), len(queries)

	uf := NewUnionFindArrayWithUndo(n)
	undo := func() {
		uf.Undo()
	}
	mutate := func(id int) {
		u, v := edges[id][0], edges[id][1]
		uf.Union(u, v)
	}
	predicate := func(qid int) bool {
		x, y, z := queries[qid][0], queries[qid][1], queries[qid][2]
		if uf.IsConnected(x, y) {
			return uf.GetSize(x) >= z
		} else {
			return uf.GetSize(x)+uf.GetSize(y) >= z
		}
	}
	res := ParallelSqrtSearchUndo(m, q, mutate, undo, predicate)
	for i := range res {
		res[i]++
	}
	return res
}

//   - 给定一个长度为n的操作序列, 按顺序执行这些操作;
//   - 给定q个查询,每个查询形如:"第一次条件qi为真(满足条件)是在第几次操作?".
//     !要求对条件为真的判定具有单调性，即某个操作后qi为真,后续操作都会满足qi为真.

// 返回:
//   - -1 => 不需要操作就满足条件的查询.
//   - [0, n) => 满足条件的最早的操作的编号(0-based).
//   - n => 执行完所有操作后都不满足条件的查询.
//
// !这种操作分块方法是整体二分的下位版本，但是实现起来比较简单。
func ParallelSqrtSearchUndo(
	n, q int,
	mutate func(mutationId int), // 执行第 mutationId 次操作，一共调用 nsqrtn 次.
	undo func(), // 撤销上一次`mutate`操作，一共调用 nsqrtn 次.
	predicate func(queryId int) bool, // 判断第 queryId 次查询是否满足条件，一共调用 nsqrtn 次.
) []int {
	res := make([]int, q)
	for i := range res {
		res[i] = n
	}

	// 不需要操作就满足条件的查询
	for i := 0; i < q; i++ {
		if predicate(i) {
			res[i] = -1
		}
	}

	sqrt := 1 + int(math.Sqrt(float64(n)))
	for i := 0; i < n; i++ {
		mutate(i)
		if i%sqrt == 0 || i == n-1 {
			for j := 0; j < q; j++ {
				if res[j] != n {
					continue
				}
				ptr := i
				for predicate(j) {
					res[j] = ptr
					ptr--
					undo()
				}
				for ptr < i {
					ptr++
					mutate(ptr)
				}
			}
		}
	}

	return res
}

func NewUnionFindArrayWithUndo(n int) *UnionFindArrayWithUndo {
	data := make([]int, n)
	for i := range data {
		data[i] = -1
	}
	return &UnionFindArrayWithUndo{Part: n, n: n, data: data}
}

type historyItem struct{ a, b int }

type UnionFindArrayWithUndo struct {
	Part      int
	n         int
	innerSnap int
	data      []int
	history   []historyItem // (root,data)
}

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *UnionFindArrayWithUndo) Undo() bool {
	if len(uf.history) == 0 {
		return false
	}
	small, smallData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	big, bigData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	uf.data[small] = smallData
	uf.data[big] = bigData
	if big != small {
		uf.Part++
	}
	return true
}

// 保存并查集当前的状态.
//
//	!Snapshot() 之后可以调用 Rollback(-1) 回滚到这个状态.
func (uf *UnionFindArrayWithUndo) Snapshot() {
	uf.innerSnap = len(uf.history) >> 1
}

// 回滚到指定的状态.
//
//	state 为 -1 表示回滚到上一次 `SnapShot` 时保存的状态.
//	其他值表示回滚到状态id为state时的状态.
func (uf *UnionFindArrayWithUndo) Rollback(state int) bool {
	if state == -1 {
		state = uf.innerSnap
	}
	state <<= 1
	if state < 0 || state > len(uf.history) {
		return false
	}
	for state < len(uf.history) {
		uf.Undo()
	}
	return true
}

// 获取当前并查集的状态id.
//
//	也就是当前合并(Union)被调用的次数.
func (uf *UnionFindArrayWithUndo) GetState() int {
	return len(uf.history) >> 1
}

func (uf *UnionFindArrayWithUndo) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UnionFindArrayWithUndo) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	uf.history = append(uf.history, historyItem{x, uf.data[x]})
	uf.history = append(uf.history, historyItem{y, uf.data[y]})
	if x == y {
		return false
	}
	if uf.data[x] > uf.data[y] {
		x ^= y
		y ^= x
		x ^= y
	}
	uf.data[x] += uf.data[y]
	uf.data[y] = x
	uf.Part--
	return true
}

func (uf *UnionFindArrayWithUndo) Find(x int) int {
	cur := x
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}
func (ufa *UnionFindArrayWithUndo) SetPart(part int) { ufa.Part = part }

func (uf *UnionFindArrayWithUndo) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }

func (uf *UnionFindArrayWithUndo) GetSize(x int) int { return -uf.data[uf.Find(x)] }

func (ufa *UnionFindArrayWithUndo) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArrayWithUndo) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
