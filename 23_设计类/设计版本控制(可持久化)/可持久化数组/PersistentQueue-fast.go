// !可持久化队列
// https://judge.yosupo.jp/problem/persistent_queue
// 初始的队列版本为-1
// 你需要维护这样的一个队列，支持如下几种操作
// 0 versioni x => 在版本 versioni 上入队x(append) 版本+1
// 1 versioni => 在版本 versioni 上出队(popleft) 版本+1 输出出队的元素
// q<=5e5 -1<=versioni<i
// https://hitonanode.github.io/cplib-cpp/data_structure/persistent_queue.hpp

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	queue := NewPersistentQueue(q)
	for i := 0; i < q; i++ {
		var op, version, x int
		fmt.Fscan(in, &op, &version)
		version++
		if op == 0 {
			fmt.Fscan(in, &x)
			queue.Append(version, x)
		} else {
			_, x := queue.PopLeft(version)
			fmt.Fprintln(out, x)
		}
	}
}

type E = int

type PersistentQueue struct {
	CurVersion int // 初始版本为0
	log        int
	data       []E     // Elements on each node of tree
	par        [][]int // binary-lifted parents
	backId     []int   // backId[t] = leaf id of the tree at time t
	size       []int   // size[t] = size of the queue at time t
}

func NewPersistentQueue(maxQuery int) *PersistentQueue {
	log := bits.Len(uint(maxQuery))
	data := make([]E, 0, maxQuery)
	par := make([][]int, 0, maxQuery)
	backId := make([]int, 1, maxQuery)
	size := make([]int, 1, maxQuery)
	return &PersistentQueue{
		log:    log,
		data:   data,
		par:    par,
		backId: backId,
		size:   size,
	}
}

// version>=0
func (q *PersistentQueue) Append(version int, value E) (newVersion int) {
	q.CurVersion++
	newId := len(q.data)
	q.data = append(q.data, value)
	q.par = append(q.par, make([]int, q.log))
	q.par[newId][0] = q.backId[version]
	q.backId = append(q.backId, newId)
	q.size = append(q.size, q.size[version]+1)
	for d := 1; d < q.log; d++ {
		q.par[newId][d] = q.par[q.par[newId][d-1]][d-1]
	}
	return q.CurVersion
}

// version>=0
func (q *PersistentQueue) PopLeft(version int) (newVersion int, popped E) {
	q.CurVersion++
	r := q.backId[version]
	len_ := q.size[version] - 1
	q.backId = append(q.backId, r)
	q.size = append(q.size, len_)
	for d := 0; d < q.log; d++ {
		if len_>>d&1 == 1 {
			r = q.par[r][d]
		}
	}
	newVersion = q.CurVersion
	popped = q.data[r]
	return
}
