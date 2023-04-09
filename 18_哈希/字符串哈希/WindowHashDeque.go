// 滑动窗口哈希值

package main

import (
	"fmt"
	"strings"
	"time"
)

type WindowHash struct {
	mod, base, hash, inv uint
	power                []uint
	deque                *ArrayDeque
}

// mod需要和base互质.
// mod: 1e9+7/1e9+9/1e9+21/1e9+33
// base: 131/13331
func NewWindowHashDeque(mod, base uint) *WindowHash {
	w := &WindowHash{
		mod:   mod,
		base:  base,
		inv:   modInv(base, mod),
		deque: NewArrayDeque(16),
		power: []uint{1},
	}
	return w
}

func (w *WindowHash) Query() uint {
	return w.hash
}

func (w *WindowHash) Append(ord uint) {
	w.hash = (w.hash*w.base + ord) % w.mod
	w.deque.Append(ord)
}

func (w *WindowHash) Pop() {
	w.hash = ((w.hash - w.deque.At(w.deque.Len()-1)) % w.mod) * w.inv % w.mod
	if w.hash < 0 {
		w.hash += w.mod
	}
	w.deque.Pop()
}

func (w *WindowHash) AppendLeft(ord uint) {
	w.expand(w.deque.Len())
	pow := w.power[w.deque.Len()]
	w.hash = (w.hash + ord*pow) % w.mod
	w.deque.AppendLeft(ord)
}

func (w *WindowHash) PopLeft() {
	w.expand(w.deque.Len() - 1)
	pow := w.power[w.deque.Len()-1]
	w.hash = (w.hash - w.deque.At(0)*pow) % w.mod
	if w.hash < 0 {
		w.hash += w.mod
	}
	w.deque.PopLeft()
}

func (w *WindowHash) Len() int {
	return w.deque.Len()
}

func (w *WindowHash) String() string {
	return fmt.Sprintf("%v", w.deque)
}

func (w *WindowHash) expand(size int) {
	if len(w.power) < size+1 {
		preSz := len(w.power)
		w.power = append(w.power, make([]uint, size+1-preSz)...)
		for i := preSz - 1; i < size; i++ {
			w.power[i+1] = w.power[i] * w.base
		}
	}
}

func exgcd(a, b uint) (gcd, x, y uint) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

// 模逆元,注意模为1时不存在逆元
func modInv(a, mod uint) uint {
	gcd, x, _ := exgcd(a, mod)
	if gcd != 1 {
		panic(fmt.Sprintf("no inverse element for %d", a))
	}
	return (x%mod + mod) % mod
}

type E = uint

func NewArrayDeque(numElements uint) *ArrayDeque {
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
	queue.ForEach(func(value E, _ int) {
		values = append(values, fmt.Sprintf("%v", value))
	})

	res = append(res, strings.Join(values, ", "), "}")
	return strings.Join(res, "")
}

func main() {
	const MOD1 uint = 1e9 + 7
	const BASE1 uint = 13331
	window := NewWindowHashDeque(MOD1, BASE1)
	window.Append('a')
	fmt.Println(window.Query())
	window.Append('b')
	fmt.Println(window.Query())
	window.PopLeft()
	fmt.Println(window.Query())

	time1 := time.Now()
	for i := 0; i < 1e6; i++ {
		window.Append('a')
		window.PopLeft()
		window.Append('a')
		window.AppendLeft('c')
		window.Pop()
	}
	fmt.Println(time.Since(time1))
}
