// !用两个 slice 头对头拼在一起实现(每个slice头部只删除元素，不添加元素，互相弥补劣势)
// 在知道数据量的情况下，也可以直接创建一个两倍数据量大小的 slice，然后用两个下标表示头尾，初始化在 slice 正中
// https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e
// NewDeque

package arraydeque

type Deque[D any] struct{ left, right []D }

func NewDeque[D any](initCapacity int) *Deque[D] {
	return &Deque[D]{make([]D, 0, 1+initCapacity/2), make([]D, 0, 1+initCapacity/2)}
}

func (q *Deque[D]) Empty() bool {
	return len(q.left) == 0 && len(q.right) == 0
}

func (q *Deque[D]) Len() int {
	return len(q.left) + len(q.right)
}

func (q *Deque[D]) AppendLeft(v D) {
	q.left = append(q.left, v)
}

func (q *Deque[D]) Append(v D) {
	q.right = append(q.right, v)
}

func (q *Deque[D]) PopLeft() (v D) {
	if len(q.left) > 0 {
		q.left, v = q.left[:len(q.left)-1], q.left[len(q.left)-1]
	} else {
		v, q.right = q.right[0], q.right[1:]
	}
	return
}

func (q *Deque[D]) Pop() (v D) {
	if len(q.right) > 0 {
		q.right, v = q.right[:len(q.right)-1], q.right[len(q.right)-1]
	} else {
		v, q.left = q.left[0], q.left[1:]
	}
	return
}

func (q *Deque[D]) Front() D {
	if len(q.left) > 0 {
		return q.left[len(q.left)-1]
	}
	return q.right[0]
}

func (q *Deque[D]) Back() D {
	if len(q.right) > 0 {
		return q.right[len(q.right)-1]
	}
	return q.left[0]
}

// 0 <= i < q.Len()
func (q *Deque[D]) At(i int) D {
	if i < len(q.left) {
		return q.left[len(q.left)-1-i]
	}
	return q.right[i-len(q.left)]
}

func (q *Deque[D]) Clear() {
	q.left = q.left[:0]
	q.right = q.right[:0]
}
