// !可持久化队列
// https://judge.yosupo.jp/problem/persistent_queue
// 初始的队列版本为-1
// 你需要维护这样的一个队列，支持如下几种操作
// 0 versioni x => 在版本 versioni 上入队x(append) 版本+1
// 1 versioni => 在版本 versioni 上出队(popleft) 版本+1 输出出队的元素
// q<=5e5 -1<=versioni<i

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
	buf := make([]*PersistentQueue, 1)
	buf[0] = NewPersistentQueue()
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var t, x int
			fmt.Fscan(in, &t, &x)
			t++
			buf = append(buf, buf[t].Append(QueueItem(x)))
		} else {
			var t int
			fmt.Fscan(in, &t)
			t++
			fmt.Fprintln(out, buf[t].Front())
			buf = append(buf, buf[t].PopLeft())
		}
	}

}

// maxMutation < 1<<20
const MUTATION_LOG int = 20

type QueueItem int

type PersistentQueue struct {
	fNode *queueNode
	bNode *queueNode
	size  int
}

type queueNode struct {
	val      QueueItem
	children [MUTATION_LOG]*queueNode
}

func NewPersistentQueue() *PersistentQueue {
	return &PersistentQueue{}
}

func (q *PersistentQueue) Append(v QueueItem) *PersistentQueue {
	t := &queueNode{val: v}
	t.children[0] = q.bNode
	for i := 1; i < MUTATION_LOG; i++ {
		c := t.children[i-1]
		if c != nil {
			t.children[i] = c.children[i-1]
		} else {
			break
		}
	}
	var fNode *queueNode
	if q.fNode != nil {
		fNode = q.fNode
	} else {
		fNode = t
	}
	return &PersistentQueue{fNode: fNode, bNode: t, size: q.size + 1}
}

func (q *PersistentQueue) PopLeft() *PersistentQueue {
	if q.fNode == nil || q.bNode == nil || q.size == 1 {
		return &PersistentQueue{}
	}
	d := q.size - 2
	t := q.bNode
	for d > 0 {
		k := bits.Len32(uint32(d)) - 1
		d -= 1 << k
		t = t.children[k]
	}
	return &PersistentQueue{fNode: t, bNode: q.bNode, size: q.size - 1}
}

func (q *PersistentQueue) Front() QueueItem {
	return q.fNode.val
}

func (q *PersistentQueue) Back() QueueItem {
	return q.bNode.val
}
