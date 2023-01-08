// 滑动窗口哈希值

package main

import (
	"fmt"
	"strings"
)

type WindowHash struct {
	mod, base, hash, inv int
	deque                *ArrayDeque
}

func NewWindowHashDeque(ords []int, mod, base int) *WindowHash {
	w := &WindowHash{
		mod:   mod,
		base:  base,
		inv:   modInv(base, mod),
		deque: NewArrayDeque(len(ords)),
	}
	for _, ord := range ords {
		w.Append(ord)
	}
	return w
}

func (w *WindowHash) Query() int {
	return w.hash
}

func (w *WindowHash) Append(ord int) {
	w.hash = (w.hash*w.base + ord) % w.mod
	w.deque.Append(ord)
}

func (w *WindowHash) Pop() {
	w.hash = ((w.hash-w.deque.At(w.deque.Len()-1))%w.mod + w.mod) * w.inv % w.mod
	w.deque.Pop()
}

func (w *WindowHash) AppendLeft(ord int) {
	pow := Pow(w.base, w.deque.Len(), w.mod)
	w.hash = (w.hash + ord*pow) % w.mod
	w.deque.AppendLeft(ord)
}

func (w *WindowHash) PopLeft() {
	pow := Pow(w.base, w.deque.Len()-1, w.mod)
	w.hash = ((w.hash-w.deque.At(0)*pow)%w.mod + w.mod) % w.mod
	w.deque.PopLeft()
}

func (w *WindowHash) Len() int {
	return w.deque.Len()
}

func (w *WindowHash) String() string {
	return fmt.Sprintf("%v", w.deque)
}

func Pow(base, exp, mod int) int {
	if exp == -1 {
		return modInv(base, mod)
	}

	base %= mod
	res := 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}

func exgcd(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

func modInv(a, mod int) int {
	gcd, x, _ := exgcd(a, mod)
	if gcd != 1 {
		panic(fmt.Sprintf("no inverse element for %d", a))
	}
	return (x%mod + mod) % mod
}

type E = int

func NewArrayDeque(numElements int) *ArrayDeque {
	half := numElements / 2
	return &ArrayDeque{
		left:  make([]E, 0, half+1),
		right: make([]E, 0, half+1),
	}
}

type ArrayDeque struct {
	left, right []E
}

func (queue *ArrayDeque) Append(value E) {
	queue.right = append(queue.right, value)
}

func (queue *ArrayDeque) AppendLeft(value E) {
	queue.left = append(queue.left, value)
}

func (queue *ArrayDeque) Pop() (value E) {
	if queue.Len() == 0 {
		return
	}

	if len(queue.right) > 0 {
		queue.right, value = queue.right[:len(queue.right)-1], queue.right[len(queue.right)-1]
	} else {
		value, queue.left = queue.left[0], queue.left[1:]
	}

	return
}

func (queue *ArrayDeque) PopLeft() (value E) {
	if queue.Len() == 0 {
		return
	}

	if len(queue.left) > 0 {
		queue.left, value = queue.left[:len(queue.left)-1], queue.left[len(queue.left)-1]
	} else {
		value, queue.right = queue.right[0], queue.right[1:]
	}

	return
}

func (queue *ArrayDeque) At(index int) (value E) {
	n := queue.Len()
	if index < 0 {
		index += n
	}

	if index < 0 || index >= n {
		return
	}

	if index < len(queue.left) {
		value = queue.left[len(queue.left)-1-index]
	} else {
		value = queue.right[index-len(queue.left)]
	}

	return
}

func (queue *ArrayDeque) ForEach(f func(value E, index int)) {
	leftLen := len(queue.left)
	for i := 0; i < leftLen; i++ {
		f(queue.left[i], i)
	}

	for i := 0; i < len(queue.right); i++ {
		f(queue.right[i], leftLen+i)
	}
}

func (queue *ArrayDeque) Len() int {
	return len(queue.left) + len(queue.right)
}

func (queue *ArrayDeque) String() string {
	res := []string{"ArrayDeque{"}
	values := []string{}
	queue.ForEach(func(value E, index int) {
		values = append(values, fmt.Sprintf("%v", value))
	})

	res = append(res, strings.Join(values, ", "), "}")
	return strings.Join(res, "")
}
