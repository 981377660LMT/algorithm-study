//
// Regex --> NFA --> DFA
//
// Description:
//   We are given a regex, eg, "a(b|c)d*e".
//   It construct NFA by recursive descent parsing,
//   and then construct DFA by powerset construction.
//
// Algorithm:
//   recursive descent parsing.
//   Robin-Scott's powerset construction.
//
// Complexity:
//   recursive descent parsiong: O(n)
//   powerset construction: O(2^n) in worst case.
//
// Verified:
//   SPOJ 10354: Count Strings
//
// https://github.com/spaghetti-source/algorithm/blob/master/string/NFAtoDFA.cc
//
// Thompson 算法
//
// TODO

package main

const NFA_STATE int32 = 128
const DFA_STATE int32 = 1000
const ALPHA int32 = 256

var nfaFactory = NewNfaFactory()
var dfaFactory = NewDfaFactory()

type NfaFactory struct {
	size int32
	next [NFA_STATE][ALPHA][]int32
}

func NewNfaFactory() *NfaFactory {
	return &NfaFactory{}
}

func (nf *NfaFactory) NewNode() int32 {
	for a := int32(0); a < ALPHA; a++ {
		nf.next[nf.size][a] = nf.next[nf.size][a][:0]
	}
	nf.size++
	return nf.size - 1
}

func (nf *NfaFactory) Symbol(a int32) *Nfa {
	begin := nf.NewNode()
	end := nf.NewNode()
	nf.next[begin][a] = append(nf.next[begin][a], end)
	return &Nfa{begin: begin, end: end}
}

func (nf *NfaFactory) Unite(x, y *Nfa) *Nfa {
	begin := nf.NewNode()
	end := nf.NewNode()
	nf.next[begin][0] = append(nf.next[begin][0], x.begin, y.begin)
	nf.next[x.end][0] = append(nf.next[x.end][0], end)
	nf.next[y.end][0] = append(nf.next[y.end][0], end)
	return &Nfa{begin: begin, end: end}
}

func (nf *NfaFactory) Concat(x, y *Nfa) *Nfa {
	nf.next[x.end][0] = append(nf.next[x.end][0], y.begin)
	return &Nfa{begin: x.begin, end: y.end}
}

func (nf *NfaFactory) Star(x *Nfa) *Nfa {
	begin := nf.NewNode()
	end := nf.NewNode()
	nf.next[begin][0] = append(nf.next[begin][0], x.begin, end)
	nf.next[x.end][0] = append(nf.next[x.end][0], x.begin, end)
	return &Nfa{begin: begin, end: end}
}

type Nfa struct {
	begin, end int32
}

func (nfa *Nfa) Run(s string) bool {
	x := NewBitSetSimple(NFA_STATE, 0)
	nfa.Closure(nfa.begin, x)
	for _, c := range s {
		y := NewBitSetSimple(NFA_STATE, 0)
		for u := int32(0); u < nfaFactory.size; u++ {
			if x.Has(u) {
				for _, v := range nfaFactory.next[u][c] {
					nfa.Closure(v, y)
				}
			}
		}
		x = y
	}
	return x.Has(nfa.end)
}

func (nfa *Nfa) Closure(u int32, x *BitSetSimple) {
	x.Add(u)
	for _, v := range nfaFactory.next[u][0] {
		if !x.Has(v) {
			nfa.Closure(v, x)
		}
	}
}

type DfaFactory struct {
	size int32
	next [DFA_STATE][ALPHA]int32
}

func NewDfaFactory() *DfaFactory {
	return &DfaFactory{}
}

func (df *DfaFactory) NewNode() int32 {
	for i := int32(0); i < ALPHA; i++ {
		df.next[df.size][i] = -1
	}
	df.size++
	return df.size - 1
}

type Dfa struct {
	begin int32
	end   [DFA_STATE]int32
}

func NewDfa() *Dfa {
	return &Dfa{}
}

func (d *Dfa) Run(s string) bool {
	u := d.begin
	for _, c := range s {
		u = dfaFactory.next[u][c]
		if u < 0 {
			return false
		}
	}
	return d.end[u] != 0
}

func Parse(s string) *Nfa {
	var regex, factor, term func() *Nfa

	ptr := int32(0)

	regex = func() *Nfa {
		a := factor()
		if s[ptr] == '|' {
			ptr++
			a = nfaFactory.Unite(a, regex())
		}
		return a
	}

	factor = func() *Nfa {
		a := term()
		if s[ptr] == '*' {
			a = nfaFactory.Star(a)
			ptr++
		}
		if ptr < int32(len(s)) && s[ptr] != '|' && s[ptr] != ')' {
			a = nfaFactory.Concat(a, factor())
		}
		return a
	}

	term = func() *Nfa {
		if s[ptr] == '(' {
			ptr++
			a := regex()
			ptr++
			return a
		}
		a := nfaFactory.Symbol(int32(s[ptr]))
		ptr++
		return a
	}

	return regex()
}

// nfaToDfa
func Convert(x *Nfa) *Dfa {
	z := NewDfa()
	states := make(map[*BitSetSimple]int32)
	process := make([]*BitSetSimple, 1)
	process[0] = NewBitSetSimple(NFA_STATE, 0)
	x.Closure(x.begin, process[0])
	z.begin = dfaFactory.NewNode()
	states[process[0]] = z.begin
	for len(process) > 0 {
		S := process[len(process)-1]
		process = process[:len(process)-1]
		for a := int32(1); a < ALPHA; a++ {
			T := NewBitSetSimple(NFA_STATE, 0)
			for u := int32(0); u < nfaFactory.size; u++ {
				if S.Has(u) {
					for _, v := range nfaFactory.next[u][a] {
						x.Closure(v, T)
					}
				}
			}
			if T.data == nil {
				continue
			}
			if _, ok := states[T]; !ok {
				states[T] = dfaFactory.NewNode()
				if T.Has(x.end) {
					z.end[states[T]] = 1
				} else {
					z.end[states[T]] = 0
				}
				process = append(process, T)
			}
			dfaFactory.next[states[S]][a] = states[T]
		}
	}
	return z
}

type BitSetSimple struct {
	n    int32
	data []uint64
}

func NewBitSetSimple(n int32, filledValue int) *BitSetSimple {
	if !(filledValue == 0 || filledValue == 1) {
		panic("filledValue should be 0 or 1")
	}
	data := make([]uint64, n>>6+1)
	if filledValue == 1 {
		for i := range data {
			data[i] = ^uint64(0)
		}
		if n != 0 {
			data[len(data)-1] >>= (int32(len(data)) << 6) - n
		}
	}
	return &BitSetSimple{n: n, data: data}
}

func (bs *BitSetSimple) Add(i int32) *BitSetSimple {
	bs.data[i>>6] |= 1 << (i & 63)
	return bs
}

func (bs *BitSetSimple) Has(i int32) bool {
	return bs.data[i>>6]>>(i&63)&1 == 1
}

func (bs *BitSetSimple) Discard(i int32) {
	bs.data[i>>6] &^= 1 << (i & 63)
}

func (bs *BitSetSimple) Flip(i int32) {
	bs.data[i>>6] ^= 1 << (i & 63)
}
