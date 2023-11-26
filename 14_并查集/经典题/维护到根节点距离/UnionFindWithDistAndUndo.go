// 由于每次Union不一定会修改成功,从而不记录修改
// (实际上这种设计并不好，但是出于性能考虑，这里还是这么做了)
// 因此不提供Undo操作,只提供GetTime/Rollback操作

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	Arc90B()
	// demo()
}

func demo() {
	uf := NewUnionFindWithDistAndUndo(10)
	fmt.Println(uf.Union(1, 2, 1))
	fmt.Println(uf.Union(2, 3, 1))
	fmt.Println(uf.Union(1, 3, 2))
	fmt.Println(uf.DistToRoot(1))
	fmt.Println(uf.DistToRoot(2))
	fmt.Println(uf.Dist(1, 2))
	uf.Union(2, 3, 2)
	fmt.Println(uf.GetSize(1))

	fmt.Println(uf.GetSize(1))
}

// https://atcoder.jp/contests/abc087/tasks/arc090_b
func Arc90B() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	uf := NewUnionFindWithDistAndUndo(n + 10)
	for i := 0; i < m; i++ {
		var left, right, weight int
		fmt.Fscan(in, &left, &right, &weight)
		if !uf.Union(left, right, weight) {
			fmt.Fprintln(out, "No")
			return
		}
	}
	fmt.Fprintln(out, "Yes")
}

type T = int

func e() T        { return 0 }
func op(x, y T) T { return x + y }
func inv(x T) T   { return -x }

// 维护到每个组根节点距离的可撤销并查集.
// 用于维护环的权值，树上的距离等.
type UnionFindWithDistAndUndo struct {
	data *RollbackArray
}

func NewUnionFindWithDistAndUndo(n int) *UnionFindWithDistAndUndo {
	return &UnionFindWithDistAndUndo{
		data: NewRollbackArray(n, func(index int) arrayItem { return arrayItem{parent: -1, dist: e()} }),
	}
}

// distToRoot(parent) + dist = distToRoot(child).
func (uf *UnionFindWithDistAndUndo) Union(parent int, child int, dist T) bool {
	v1, x1 := uf.Find(parent)
	v2, x2 := uf.Find(child)
	if v1 == v2 {
		return dist == op(x2, inv(x1))
	}
	s1, s2 := -uf.data.Get(v1).parent, -uf.data.Get(v2).parent
	if s1 < s2 {
		s1, s2 = s2, s1
		v1, v2 = v2, v1
		x1, x2 = x2, x1
		dist = inv(dist)
	}
	// v1 <- v2
	dist = op(x1, dist)
	dist = op(dist, inv(x2))
	uf.data.Set(v2, arrayItem{parent: v1, dist: dist})
	uf.data.Set(v1, arrayItem{parent: -(s1 + s2), dist: e()})
	return true
}

// 返回v所在组的根节点和到v到根节点的距离.
func (uf *UnionFindWithDistAndUndo) Find(v int) (root int, distToRoot T) {
	root, distToRoot = v, e()
	for {
		item := uf.data.Get(root)
		if item.parent < 0 {
			break
		}
		distToRoot = op(distToRoot, item.dist)
		root = item.parent
	}
	return
}

// Dist(x, y) = DistToRoot(x) - DistToRoot(y).
// 如果x和y不在同一个集合,抛出错误.
func (uf *UnionFindWithDistAndUndo) Dist(x int, y int) T {
	vx, dx := uf.Find(x)
	vy, dy := uf.Find(y)
	if vx != vy {
		panic("x and y are not in the same set")
	}
	return op(dx, inv(dy))
}

func (uf *UnionFindWithDistAndUndo) DistToRoot(x int) T {
	_, dx := uf.Find(x)
	return dx
}

func (uf *UnionFindWithDistAndUndo) GetTime() int {
	return uf.data.GetTime()
}

func (uf *UnionFindWithDistAndUndo) Rollback(time int) {
	uf.data.Rollback(time)
}

func (uf *UnionFindWithDistAndUndo) GetSize(x int) int {
	root, _ := uf.Find(x)
	return -uf.data.Get(root).parent
}

func (uf *UnionFindWithDistAndUndo) GetGroups() map[int][]int {
	res := make(map[int][]int)
	for i := 0; i < uf.data.Len(); i++ {
		root, _ := uf.Find(i)
		res[root] = append(res[root], i)
	}
	return res
}

type arrayItem struct {
	parent int
	dist   T
}

type historyItem struct {
	index int
	value arrayItem
}

type RollbackArray struct {
	n       int
	data    []arrayItem
	history []historyItem
}

func NewRollbackArray(n int, f func(index int) arrayItem) *RollbackArray {
	data := make([]arrayItem, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray{
		n:    n,
		data: data,
	}
}

func (r *RollbackArray) GetTime() int {
	return len(r.history)
}

func (r *RollbackArray) Rollback(time int) {
	for len(r.history) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair.index] = pair.value
	}
}

func (r *RollbackArray) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *RollbackArray) Get(index int) arrayItem {
	return r.data[index]
}

func (r *RollbackArray) Set(index int, value arrayItem) {
	r.history = append(r.history, historyItem{index: index, value: r.data[index]})
	r.data[index] = value
}

func (r *RollbackArray) GetAll() []arrayItem {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray) Len() int {
	return r.n
}
