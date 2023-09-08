// #pragma once

// // https://codeforces.com/contest/914/problem/F
// // https://yukicoder.me/problems/no/142
// // わずかに普通の bitset より遅いときもあるようだが，
// // 固定長にしたくないときや slice 操作が必要なときに使う
// struct My_Bitset {
//   using T = My_Bitset;
//   int N;
//   vc<u64> dat;

//   // x で埋める
//   My_Bitset(int N = 0, int x = 0) : N(N) {
//     assert(x == 0 || x == 1);
//     u64 v = (x == 0 ? 0 : -1);
//     dat.assign((N + 63) >> 6, v);
//     if (N) dat.back() >>= (64 * len(dat) - N);
//   }

//   int size() { return N; }

//   void resize(int size) {
//     dat.resize((size + 63) >> 6);
//     int remainingBits = size & 63;
//     if (remainingBits != 0) {
//       u64 mask = (u64(1) << remainingBits) - 1;
//       dat.back() &= mask;
//     }
//     N = size;
//   }

//   // thanks to chatgpt!
//   class Proxy {
//   public:
//     Proxy(vc<u64> &d, int i) : dat(d), index(i) {}
//     operator bool() const { return (dat[index >> 6] >> (index & 63)) & 1; }

//     Proxy &operator=(u64 value) {
//       dat[index >> 6] &= ~(u64(1) << (index & 63));
//       dat[index >> 6] |= (value & 1) << (index & 63);
//       return *this;
//     }
//     void flip() {
//       dat[index >> 6] ^= (u64(1) << (index & 63)); // XOR to flip the bit
//     }

//   private:
//     vc<u64> &dat;
//     int index;
//   };

//   Proxy operator[](int i) { return Proxy(dat, i); }

//   T &operator&=(const T &p) {
//     assert(N == p.N);
//     FOR(i, len(dat)) dat[i] &= p.dat[i];
//     return *this;
//   }
//   T &operator|=(const T &p) {
//     assert(N == p.N);
//     FOR(i, len(dat)) dat[i] |= p.dat[i];
//     return *this;
//   }
//   T &operator^=(const T &p) {
//     assert(N == p.N);
//     FOR(i, len(dat)) dat[i] ^= p.dat[i];
//     return *this;
//   }
//   T operator&(const T &p) const { return T(*this) &= p; }
//   T operator|(const T &p) const { return T(*this) |= p; }
//   T operator^(const T &p) const { return T(*this) ^= p; }

//   int count() {
//     int ans = 0;
//     for (u64 val: dat) ans += popcnt(val);
//     return ans;
//   }

//   int next(int i) {
//     chmax(i, 0);
//     if (i >= N) return N;
//     int k = i >> 6;
//     {
//       u64 x = dat[k];
//       int s = i & 63;
//       x = (x >> s) << s;
//       if (x) return (k << 6) | lowbit(x);
//     }
//     FOR(idx, k + 1, len(dat)) {
//       if (dat[idx] == 0) continue;
//       return (idx << 6) | lowbit(dat[idx]);
//     }
//     return N;
//   }

//   int prev(int i) {
//     chmin(i, N - 1);
//     if (i <= -1) return -1;
//     int k = i >> 6;
//     if ((i & 63) < 63) {
//       u64 x = dat[k];
//       x &= (u64(1) << ((i & 63) + 1)) - 1;
//       if (x) return (k << 6) | topbit(x);
//       --k;
//     }
//     FOR_R(idx, k + 1) {
//       if (dat[idx] == 0) continue;
//       return (idx << 6) | topbit(dat[idx]);
//     }
//     return -1;
//   }

//   My_Bitset range(int L, int R) {
//     assert(L <= R);
//     My_Bitset p(R - L);
//     int rm = (R - L) & 63;
//     FOR(rm) {
//       p[R - L - 1] = bool((*this)[R - 1]);
//       --R;
//     }
//     int n = (R - L) >> 6;
//     int hi = L & 63;
//     int lo = 64 - hi;
//     int s = L >> 6;
//     if (hi == 0) {
//       FOR(i, n) { p.dat[i] ^= dat[s + i]; }
//     } else {
//       FOR(i, n) { p.dat[i] ^= (dat[s + i] >> hi) ^ (dat[s + i + 1] << lo); }
//     }
//     return p;
//   }

//   int count_range(int L, int R) {
//     assert(L <= R);
//     int cnt = 0;
//     while ((L < R) && (L & 63)) cnt += (*this)[L++];
//     while ((L < R) && (R & 63)) cnt += (*this)[--R];
//     int l = L >> 6, r = R >> 6;
//     FOR(i, l, r) cnt += popcnt(dat[i]);
//     return cnt;
//   }

//   // [L,R) に p を代入
//   void assign_to_range(int L, int R, My_Bitset &p) {
//     assert(p.N == R - L);
//     int a = 0, b = p.N;
//     while (L < R && (L & 63)) { (*this)[L++] = bool(p[a++]); }
//     while (L < R && (R & 63)) { (*this)[--R] = bool(p[--b]); }
//     // p[a:b] を [L:R] に
//     int l = L >> 6, r = R >> 6;
//     int s = a >> 6, t = b >> t;
//     int n = r - l;
//     if (!(a & 63)) {
//       FOR(i, n) dat[l + i] = p.dat[s + i];
//     } else {
//       int hi = a & 63;
//       int lo = 64 - hi;
//       FOR(i, n) dat[l + i] = (p.dat[s + i] >> hi) | (p.dat[1 + s + i] << lo);
//     }
//   }

//   // [L,R) に p を xor
//   void xor_to_range(int L, int R, My_Bitset &p) {
//     assert(p.N == R - L);
//     int a = 0, b = p.N;
//     while (L < R && (L & 63)) {
//       dat[L >> 6] ^= u64(p[a]) << (L & 63);
//       ++a, ++L;
//     }
//     while (L < R && (R & 63)) {
//       --b, --R;
//       dat[R >> 6] ^= u64(p[b]) << (R & 63);
//     }
//     // p[a:b] を [L:R] に
//     int l = L >> 6, r = R >> 6;
//     int s = a >> 6, t = b >> t;
//     int n = r - l;
//     if (!(a & 63)) {
//       FOR(i, n) dat[l + i] ^= p.dat[s + i];
//     } else {
//       int hi = a & 63;
//       int lo = 64 - hi;
//       FOR(i, n) dat[l + i] ^= (p.dat[s + i] >> hi) | (p.dat[1 + s + i] << lo);
//     }
//   }

//   // [L,R) に p を and
//   void and_to_range(int L, int R, My_Bitset &p) {
//     assert(p.N == R - L);
//     int a = 0, b = p.N;
//     while (L < R && (L & 63)) {
//       if (!p[a++]) (*this)[L++] = 0;
//     }
//     while (L < R && (R & 63)) {
//       if (!p[--b]) (*this)[--R] = 0;
//     }
//     // p[a:b] を [L:R] に
//     int l = L >> 6, r = R >> 6;
//     int s = a >> 6, t = b >> t;
//     int n = r - l;
//     if (!(a & 63)) {
//       FOR(i, n) dat[l + i] &= p.dat[s + i];
//     } else {
//       int hi = a & 63;
//       int lo = 64 - hi;
//       FOR(i, n) dat[l + i] &= (p.dat[s + i] >> hi) | (p.dat[1 + s + i] << lo);
//     }
//   }

//   // [L,R) に p を or
//   void or_to_range(int L, int R, My_Bitset &p) {
//     assert(p.N == R - L);
//     int a = 0, b = p.N;
//     while (L < R && (L & 63)) {
//       dat[L >> 6] |= u64(p[a]) << (L & 63);
//       ++a, ++L;
//     }
//     while (L < R && (R & 63)) {
//       --b, --R;
//       dat[R >> 6] |= u64(p[b]) << (R & 63);
//     }
//     // p[a:b] を [L:R] に
//     int l = L >> 6, r = R >> 6;
//     int s = a >> 6, t = b >> t;
//     int n = r - l;
//     if (!(a & 63)) {
//       FOR(i, n) dat[l + i] |= p.dat[s + i];
//     } else {
//       int hi = a & 63;
//       int lo = 64 - hi;
//       FOR(i, n) dat[l + i] |= (p.dat[s + i] >> hi) | (p.dat[1 + s + i] << lo);
//     }
//   }

//   string to_string() const {
//     string S;
//     FOR(i, N) S += '0' + (dat[i >> 6] >> (i & 63) & 1);
//     return S;
//   }

//   // bitset に仕様を合わせる
//   void set(int i) { (*this)[i] = 1; }
//   void reset(int i) { (*this)[i] = 0; }
//   void flip(int i) { (*this)[i].flip(); }
//   void set() { fill(all(dat), 0); }
//   void reset() {
//     fill(all(dat), u64(-1));
//     resize(N);
//   }

//   int _Find_first() { return next(0); }
//   int _Find_next(int p) { return next(p + 1); }
// };

package main

import (
	"fmt"
	"math/bits"
	"strings"
	"time"
)

func main() {
	bs := NewBitsetDynamic(100, 0).Add(0).Add(1).Add(39)
	// IXorRange
	other := NewBitsetDynamic(40, 0).Add(0).Add(1).Add(39)
	bs.Set(other, 1)
	fmt.Println(bs, other)
	fmt.Println(bs.OnesCount(0, bs.n), bs)
	bs.Add(80)
	fmt.Println(bs.OnesCount(0, 79), bs)

	bs.AddRange(0, 80)
	fmt.Println(bs.OnesCount(0, 100), bs)
	bs.DiscardRange(0, 80)
	fmt.Println(bs.OnesCount(0, 100), bs)
	bs.FlipRange(30, 79)
	fmt.Println(bs.OnesCount(0, 100), bs)

	bs.Fill(1)
	fmt.Println(bs.OnesCount(0, 100), bs)
	bs.Fill(0)
	fmt.Println(bs.OnesCount(0, 100), bs)
	bs.AddRange(20, 90)
	fmt.Println(bs.AllOne(19, 90), bs)
	fmt.Println(bs.AllZero(90, 100))

	fmt.Println(bs.IndexOfOne(78))

	time1 := time.Now()
	// resize 1e5
	for i := 0; i < 1e5; i++ {
		bs.Resize(i)
	}
	fmt.Println(time.Since(time1))

}

// 动态bitset，支持切片操作.
type BitSetDynamic struct {
	n    int
	data []uint64
}

// 建立一个大小为 n 的 bitset，初始值为 filledValue.
func NewBitsetDynamic(n int, filledValue int) *BitSetDynamic {
	if !(filledValue == 0 || filledValue == 1) {
		panic("filledValue should be 0 or 1")
	}
	data := make([]uint64, (n+63)>>6)
	if filledValue == 1 {
		for i := range data {
			data[i] = ^uint64(0)
		}
		if n != 0 {
			data[len(data)-1] >>= 64 - n&63
		}
	}
	return &BitSetDynamic{n: n, data: data}
}

func (bs *BitSetDynamic) Add(i int) *BitSetDynamic {
	bs.data[i>>6] |= 1 << (i & 63)
	return bs
}

func (bs *BitSetDynamic) Has(i int) bool {
	return bs.data[i>>6]>>(i&63)&1 == 1
}

func (bs *BitSetDynamic) Discard(i int) {
	bs.data[i>>6] &^= 1 << (i & 63)
}

func (bs *BitSetDynamic) Flip(i int) {
	bs.data[i>>6] ^= 1 << (i & 63)
}

func (bs *BitSetDynamic) AddRange(start, end int) {
	maskL := ^uint64(0) << (start & 63)
	maskR := ^uint64(0) << (end & 63)
	i := start >> 6
	if i == end>>6 {
		bs.data[i] |= maskL ^ maskR
		return
	}
	bs.data[i] |= maskL
	for i++; i < end>>6; i++ {
		bs.data[i] = ^uint64(0)
	}
	bs.data[i] |= ^maskR
}

func (bs *BitSetDynamic) DiscardRange(start, end int) {
	maskL := ^uint64(0) << (start & 63)
	maskR := ^uint64(0) << (end & 63)
	i := start >> 6
	if i == end>>6 {
		bs.data[i] &= ^maskL | maskR
		return
	}
	bs.data[i] &= ^maskL
	for i++; i < end>>6; i++ {
		bs.data[i] = 0
	}
	bs.data[i] &= maskR
}

func (bs *BitSetDynamic) FlipRange(start, end int) {
	maskL := ^uint64(0) << (start & 63)
	maskR := ^uint64(0) << (end & 63)
	i := start >> 6
	if i == end>>6 {
		bs.data[i] ^= maskL ^ maskR
		return
	}
	bs.data[i] ^= maskL
	for i++; i < end>>6; i++ {
		bs.data[i] = ^bs.data[i]
	}
	bs.data[i] ^= ^maskR
}

func (bs *BitSetDynamic) Fill(zeroOrOne int) {
	if zeroOrOne == 0 {
		for i := range bs.data {
			bs.data[i] = 0
		}
	} else {
		for i := range bs.data {
			bs.data[i] = ^uint64(0)
		}
		if bs.n != 0 {
			bs.data[len(bs.data)-1] >>= 64 - bs.n&63
		}
	}
}

func (bs *BitSetDynamic) Clear() {
	for i := range bs.data {
		bs.data[i] = 0
	}
}

func (bs *BitSetDynamic) OnesCount(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > bs.n {
		end = bs.n
	}
	if start == 0 && end == bs.n {
		res := 0
		for _, v := range bs.data {
			res += bits.OnesCount64(v)
		}
		return res
	}
	pos1 := start >> 6
	pos2 := end >> 6
	if pos1 == pos2 {
		return bits.OnesCount64(bs.data[pos1] & (^uint64(0) << (start & 63)) & ((1 << (end & 63)) - 1))
	}
	count := 0
	if (start & 63) > 0 {
		count += bits.OnesCount64(bs.data[pos1] & (^uint64(0) << (start & 63)))
		pos1++
	}
	for i := pos1; i < pos2; i++ {
		count += bits.OnesCount64(bs.data[i])
	}
	if (end & 63) > 0 {
		count += bits.OnesCount64(bs.data[pos2] & ((1 << (end & 63)) - 1))
	}
	return count
}

func (bs *BitSetDynamic) AllOne(start, end int) bool {
	i := start >> 6
	if i == end>>6 {
		mask := ^uint64(0)<<(start&63) ^ ^uint64(0)<<(end&63)
		return (bs.data[i] & mask) == mask
	}
	mask := ^uint64(0) << (start & 63)
	if (bs.data[i] & mask) != mask {
		return false
	}
	for i++; i < end>>6; i++ {
		if bs.data[i] != ^uint64(0) {
			return false
		}
	}
	mask = ^uint64(0) << (end & 63)
	return ^(bs.data[end>>6] | mask) == 0
}

func (bs *BitSetDynamic) AllZero(start, end int) bool {
	i := start >> 6
	if i == end>>6 {
		mask := ^uint64(0)<<(start&63) ^ ^uint64(0)<<(end&63)
		return (bs.data[i] & mask) == 0
	}

	if (bs.data[i] >> (start & 63)) != 0 {
		return false
	}

	for i++; i < end>>6; i++ {
		if bs.data[i] != 0 {
			return false
		}
	}

	mask := ^uint64(0) << (end & 63)
	return (bs.data[end>>6] & ^mask) == 0
}

// 返回第一个 1 的下标，若不存在则返回-1.
func (bs *BitSetDynamic) IndexOfOne(position int) int {
	if position == 0 {
		for i, v := range bs.data {
			if v != 0 {
				return i<<6 | bs._lowbit(v)
			}
		}
		return -1
	}
	for i := position >> 6; i < len(bs.data); i++ {
		v := bs.data[i] & (^uint64(0) << (position & 63))
		if v != 0 {
			return i<<6 | bs._lowbit(v)
		}
		for i++; i < len(bs.data); i++ {
			if bs.data[i] != 0 {
				return i<<6 | bs._lowbit(bs.data[i])
			}
		}
	}
	return -1
}

// 返回第一个 0 的下标，若不存在则返回-1。
func (bs *BitSetDynamic) IndexOfZero(position int) int {
	if position == 0 {
		for i, v := range bs.data {
			if v != ^uint64(0) {
				return i<<6 | bs._lowbit(^v)
			}
		}
		return -1
	}

	i := position >> 6
	if i < len(bs.data) {
		v := bs.data[i]
		if position&63 != 0 {
			v |= ^((^uint64(0)) << (position & 63))
		}
		if ^v != 0 {
			res := i<<6 | bs._lowbit(^v)
			if res < bs.n {
				return res
			}
			return -1
		}
		for i++; i < len(bs.data); i++ {
			if ^bs.data[i] != 0 {
				res := i<<6 | bs._lowbit(^bs.data[i])
				if res < bs.n {
					return res
				}
				return -1
			}
		}
	}
	return -1
}

// 返回右侧第一个 1 的位置(`包含`当前位置).
//
//	如果不存在, 返回 n.
func (bs *BitSetDynamic) Next(index int) int {
	if index < 0 {
		index = 0
	}
	if index >= bs.n {
		return bs.n
	}
	k := index >> 6
	x := bs.data[k]
	s := index & 63
	x = (x >> s) << s
	if x != 0 {
		return (k << 6) | bs._lowbit(x)
	}
	for i := k + 1; i < len(bs.data); i++ {
		if bs.data[i] == 0 {
			continue
		}
		return (i << 6) | bs._lowbit(bs.data[i])
	}
	return bs.n
}

// 返回左侧第一个 1 的位置(`包含`当前位置).
//
//	如果不存在, 返回 -1.
func (bs *BitSetDynamic) Prev(index int) int {
	if index >= bs.n-1 {
		index = bs.n - 1
	}
	if index < 0 {
		return -1
	}
	k := index >> 6
	if (index & 63) < 63 {
		x := bs.data[k]
		x &= (1 << ((index & 63) + 1)) - 1
		if x != 0 {
			return (k << 6) | bs._topbit(x)
		}
		k--
	}
	for i := k; i >= 0; i-- {
		if bs.data[i] == 0 {
			continue
		}
		return (i << 6) | bs._topbit(bs.data[i])
	}
	return -1
}

func (bs *BitSetDynamic) Equals(other *BitSetDynamic) bool {
	if len(bs.data) != len(other.data) {
		return false
	}
	for i := range bs.data {
		if bs.data[i] != other.data[i] {
			return false
		}
	}
	return true
}

func (bs *BitSetDynamic) IsSubset(other *BitSetDynamic) bool {
	if bs.n > other.n {
		return false
	}
	for i, v := range bs.data {
		if (v & other.data[i]) != v {
			return false
		}
	}
	return true
}

func (bs *BitSetDynamic) IsSuperset(other *BitSetDynamic) bool {
	if bs.n < other.n {
		return false
	}
	for i, v := range other.data {
		if (v & bs.data[i]) != v {
			return false
		}
	}
	return true
}

func (bs *BitSetDynamic) Ior(other *BitSetDynamic) *BitSetDynamic {
	for i, v := range other.data {
		bs.data[i] |= v
	}
	return bs
}

func (bs *BitSetDynamic) Iand(other *BitSetDynamic) *BitSetDynamic {
	for i, v := range other.data {
		bs.data[i] &= v
	}
	return bs
}

func (bs *BitSetDynamic) Ixor(other *BitSetDynamic) *BitSetDynamic {
	for i, v := range other.data {
		bs.data[i] ^= v
	}
	return bs
}

func (bs *BitSetDynamic) Or(other *BitSetDynamic) *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] | v
	}
	return res
}

func (bs *BitSetDynamic) And(other *BitSetDynamic) *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] & v
	}
	return res
}

func (bs *BitSetDynamic) Xor(other *BitSetDynamic) *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	for i, v := range other.data {
		res.data[i] = bs.data[i] ^ v
	}
	return res
}

func (bs *BitSetDynamic) IOrRange(start, end int, other *BitSetDynamic) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := 0, other.n
	for start < end && (start&63) != 0 {
		bs.data[start>>6] |= other._get(a) << (start & 63)
		a++
		start++
	}
	for start < end && (end&63) != 0 {
		end--
		b--
		bs.data[end>>6] |= other._get(b) << (end & 63)
	}

	// other[a:b] -> this[start:end]
	l, r := start>>6, end>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] |= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] |= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

func (bs *BitSetDynamic) IAndRange(start, end int, other *BitSetDynamic) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := 0, other.n
	for start < end && (start&63) != 0 {
		if other._get(a) == 0 {
			bs.data[start>>6] &^= 1 << (start & 63)
		}
		a++
		start++
	}
	for start < end && (end&63) != 0 {
		end--
		b--
		if other._get(b) == 0 {
			bs.data[end>>6] &^= 1 << (end & 63)
		}
	}

	// other[a:b] -> this[start:end]
	l, r := start>>6, end>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] &= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] &= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}

}

func (bs *BitSetDynamic) IXorRange(start, end int, other *BitSetDynamic) {
	if other.n != end-start {
		panic("length of other must equal to end-start")
	}
	a, b := 0, other.n
	for start < end && (start&63) != 0 {
		bs.data[start>>6] ^= other._get(a) << (start & 63)
		a++
		start++
	}
	for start < end && (end&63) != 0 {
		end--
		b--
		bs.data[end>>6] ^= other._get(b) << (end & 63)
	}

	// other[a:b] -> this[start:end]
	l, r := start>>6, end>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] ^= other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] ^= (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

// 类似js中类型数组的set操作.如果超出赋值范围，抛出异常.
//
//	other: 要赋值的bitset.
//	offset: 赋值的起始元素下标.
func (bs *BitSetDynamic) Set(other *BitSetDynamic, offset int) {
	left, right := offset, offset+other.n
	if right > bs.n {
		panic("out of range")
	}
	a, b := 0, other.n
	for left < right && (left&63) != 0 {
		if other.Has(a) {
			bs.Add(left)
		} else {
			bs.Discard(left)
		}
		a++
		left++
	}
	for left < right && (right&63) != 0 {
		right--
		b--
		if other.Has(b) {
			bs.Add(right)
		} else {
			bs.Discard(right)
		}
	}

	// other[a:b] -> this[start:end]
	l, r := left>>6, right>>6
	s := a >> 6
	n := r - l
	if (a & 63) == 0 {
		for i := 0; i < n; i++ {
			bs.data[l+i] = other.data[s+i]
		}
	} else {
		hi := a & 63
		lo := 64 - hi
		for i := 0; i < n; i++ {
			bs.data[l+i] = (other.data[s+i] >> hi) | (other.data[s+i+1] << lo)
		}
	}
}

func (bs *BitSetDynamic) Slice(start, end int) *BitSetDynamic {
	if start < 0 {
		start += bs.n
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += bs.n
	}
	if end > bs.n {
		end = bs.n
	}
	if start >= end {
		return NewBitsetDynamic(0, 0)
	}
	if start == 0 && end == bs.n {
		return bs.Copy()
	}

	res := NewBitsetDynamic(end-start, 0)
	remain := (end - start) & 63
	for i := 0; i < remain; i++ {
		if bs.Has(end - 1) {
			res.Add(end - start - 1)
		}
		end--
	}

	n := (end - start) >> 6
	hi := start & 63
	lo := 64 - hi
	s := start >> 6
	if hi == 0 {
		for i := 0; i < n; i++ {
			res.data[i] = bs.data[s+i]
		}
	} else {
		for i := 0; i < n; i++ {
			res.data[i] = (bs.data[s+i] >> hi) ^ (bs.data[s+i+1] << lo)
		}
	}

	return res
}

func (bs *BitSetDynamic) Copy() *BitSetDynamic {
	res := NewBitsetDynamic(bs.n, 0)
	copy(res.data, bs.data)
	return res
}

func (bs *BitSetDynamic) Resize(size int) {
	newBits := make([]uint64, (size+63)>>6)
	copy(newBits, bs.data[:min(len(bs.data), len(newBits))])
	remainingBits := size & 63
	if remainingBits != 0 {
		mask := (1 << remainingBits) - 1
		newBits[len(newBits)-1] &= uint64(mask)
	}
	bs.data = newBits
	bs.n = size
}

func (bs *BitSetDynamic) Expand(size int) {
	if size <= bs.n {
		return
	}
	bs.Resize(size)
}

func (bs *BitSetDynamic) BitLength() int {
	return bs._lastIndexOfOne() + 1
}

// 遍历所有 1 的位置.
func (bs *BitSetDynamic) ForEach(f func(pos int)) {
	for i, v := range bs.data {
		for v != 0 {
			j := (i << 6) | bs._lowbit(v)
			f(j)
			v &= v - 1
		}
	}
}

func (bs *BitSetDynamic) Size() int {
	return bs.n
}

func (bs *BitSetDynamic) String() string {
	sb := strings.Builder{}
	sb.WriteString("BitSetDynamic{")
	nums := []string{}
	bs.ForEach(func(pos int) {
		nums = append(nums, fmt.Sprintf("%d", pos))
	})
	sb.WriteString(strings.Join(nums, ","))
	sb.WriteString("}")
	return sb.String()
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 1, 2)
func (bs *BitSetDynamic) _topbit(x uint64) int {
	if x == 0 {
		return -1
	}
	return 63 - bits.LeadingZeros64(x)
}

// (0, 1, 2, 3, 4) -> (-1, 0, 1, 0, 2)
func (bs *BitSetDynamic) _lowbit(x uint64) int {
	if x == 0 {
		return -1
	}
	return bits.TrailingZeros64(x)
}

func (bs *BitSetDynamic) _get(i int) uint64 {
	return bs.data[i>>6] >> (i & 63) & 1
}

func (bs *BitSetDynamic) _lastIndexOfOne() int {
	for i := len(bs.data) - 1; i >= 0; i-- {
		x := bs.data[i]
		if x != 0 {
			return (i << 6) | (bs._topbit(x))
		}
	}
	return -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
