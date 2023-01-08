// 滑动窗口哈希值

package main

import "fmt"

type WindowHash struct {
	mod, base, hash, inv int
	queue                []int
}

func NewWindowHash(ords []int, mod, base int) *WindowHash {
	w := &WindowHash{
		mod:   mod,
		base:  base,
		inv:   modInv(base, mod),
		queue: make([]int, 0, len(ords)),
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
	w.queue = append(w.queue, ord)
}

func (w *WindowHash) Pop() {
	w.hash = ((w.hash-w.queue[len(w.queue)-1])%w.mod + w.mod) * w.inv % w.mod
	w.queue = w.queue[:len(w.queue)-1]
}

func (w *WindowHash) PopLeft() {
	pow := Pow(w.base, len(w.queue)-1, w.mod)
	w.hash = ((w.hash-w.queue[0]*pow)%w.mod + w.mod) % w.mod
	w.queue = w.queue[1:]
}

func (w *WindowHash) Len() int {
	return len(w.queue)
}

func (w *WindowHash) String() string {
	return fmt.Sprintf("%v", w.queue)
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

func main() {
	const MOD1 int = 1e9 + 7
	const BASE1 int = 13131
	window := NewWindowHash([]int{}, MOD1, BASE1)
	window.Append('a')
	fmt.Println(window.Query())
	window.Append('b')
	fmt.Println(window.Query())
	window.PopLeft()
	fmt.Println(window.Query())
}
