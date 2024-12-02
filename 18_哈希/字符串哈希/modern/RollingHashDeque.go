// 滑动窗口字符串哈希(SlidingWindowRollingHash)

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	demo()
}

func demo() {
	H := NewRollingHashDeque(0)
	H.Append(1)
	H.Append(2)
	H.Append(1)
	H.AddLeft(2)
	fmt.Println(H.Get(0, 2)) // 2<<61 + 1
	fmt.Println(H.Get(2, 4)) // 2<<61 + 1<<61 + 2
	fmt.Println(H.PopLeft())
	fmt.Println(H.PopLeft())
	fmt.Println(H.PopLeft())
	fmt.Println(H.PopLeft())
}

const (
	mod61  uint64 = (1 << 61) - 1
	mask30 uint64 = (1 << 30) - 1
	mask31 uint64 = (1 << 31) - 1
	mask61 uint64 = mod61
)

type RollingHashDeque struct {
	a           uint64 // prefix hash(n) = dat[n]+a*pow[n]
	base, ibase uint64
	powTable    []uint64
	s           *Deque[uint64]
	data        *Deque[uint64]
}

// base: 0 表示随机生成.
func NewRollingHashDeque(base uint64) *RollingHashDeque {
	for base == 0 {
		base = rand.Uint64() % mod61 // rng61
	}
	ibase := uint64(modInv(int(base), int(mod61)))
	powTable := []uint64{1, base}
	s, data := NewDeque[uint64](0), NewDeque[uint64](0)
	data.Append(0)
	return &RollingHashDeque{base: base, ibase: ibase, powTable: powTable, s: s, data: data}
}

func (rh *RollingHashDeque) Size() int {
	return rh.s.Size()
}

func (rh *RollingHashDeque) PopLeft() uint64 {
	rh.data.PopLeft()
	ch := rh.s.PopLeft()
	rh.a = sub(mul(rh.a, rh.base), ch)
	return ch
}

func (rh *RollingHashDeque) Pop() uint64 {
	rh.data.Pop()
	return rh.s.Pop()
}

func (rh *RollingHashDeque) AddLeft(ch uint64) {
	rh.s.AppendLeft(ch)
	rh.a = mul(add(rh.a, ch), rh.ibase)
	rh.data.AppendLeft(sub(0, rh.a))
}

func (rh *RollingHashDeque) Append(ch uint64) {
	rh.s.Append(ch)
	rh.data.Append(add(mul(rh.data.Back(), rh.base), ch))
}

func (rh *RollingHashDeque) Get(l, r int) uint64 {
	if l < 0 {
		l = 0
	}
	if n := rh.data.Size(); r > n {
		r = n
	}
	if l >= r {
		return 0
	}
	return sub(rh.data.At(r), mul(rh.data.At(l), rh.pow(r-l)))
}

func (rh *RollingHashDeque) pow(i int) uint64 {
	base := rh.base
	for i >= len(rh.powTable) {
		rh.powTable = append(rh.powTable, mul(rh.powTable[len(rh.powTable)-1], base))
	}
	return rh.powTable[i]
}

// x % (2^61-1)
func mod(x uint64) uint64 {
	xu := x >> 61
	xd := x & mask61
	res := xu + xd
	if res >= mod61 {
		res -= mod61
	}
	return res
}

// a*b % (2^61-1)
func mul(a, b uint64) uint64 {
	au := a >> 31
	ad := a & mask31
	bu := b >> 31
	bd := b & mask31
	mid := ad*bu + au*bd
	midu := mid >> 30
	midd := mid & mask30
	return mod(au*bu<<1 + midu + (midd << 31) + ad*bd)
}

// a,b: modint61
func add(a, b uint64) uint64 {
	res := a + b
	if res >= mod61 {
		res -= mod61
	}
	return res
}

// a,b: modint61
func sub(a, b uint64) uint64 {
	res := a - b
	if res >= mod61 {
		res += mod61
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

type Deque[D any] struct{ l, r []D }

func NewDeque[D any](cap int32) *Deque[D] {
	return &Deque[D]{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)}
}

func (q Deque[D]) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque[D]) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque[D]) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque[D]) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque[D]) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque[D]) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque[D]) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque[D]) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque[D]) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
