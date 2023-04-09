// 滑动窗口哈希值

package main

import (
	"fmt"
	"time"
)

type WindowHash struct {
	mod, base, hash, inv uint
	queue                []uint
	power                []uint
}

// mod需要和base互质.
// mod: 1e9+7/1e9+9/1e9+21/1e9+33
// base: 131/13331
func NewWindowHash(mod, base uint) *WindowHash {
	w := &WindowHash{
		mod:   mod,
		base:  base,
		inv:   modInv(base, mod),
		power: []uint{1},
	}
	return w
}

func (w *WindowHash) Query() uint {
	return w.hash
}

func (w *WindowHash) Append(ord uint) {
	w.hash = (w.hash*w.base + ord) % w.mod
	w.queue = append(w.queue, ord)
}

func (w *WindowHash) Pop() {
	w.hash = ((w.hash - w.queue[len(w.queue)-1]) % w.mod) * w.inv % w.mod
	if w.hash < 0 {
		w.hash += w.mod
	}
	w.queue = w.queue[:len(w.queue)-1]
}

func (w *WindowHash) PopLeft() {
	w.expand(len(w.queue) - 1)
	pow := w.power[len(w.power)-1]
	w.hash = ((w.hash - w.queue[0]*pow) % w.mod) % w.mod
	if w.hash < 0 {
		w.hash += w.mod
	}
	w.queue = w.queue[1:]
}

func (w *WindowHash) Len() int {
	return len(w.queue)
}

func (w *WindowHash) String() string {
	return fmt.Sprintf("%v", w.queue)
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

func main() {
	const MOD1 uint = 1e9 + 7
	const BASE1 uint = 13331
	window := NewWindowHash(MOD1, BASE1)
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
		window.Append('a')
		window.Pop()
		window.Query()
	}
	fmt.Println(time.Since(time1))
}
