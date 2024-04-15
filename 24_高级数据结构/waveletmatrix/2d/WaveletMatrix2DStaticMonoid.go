// #include "ds/bit_vector.hpp"
// #include "ds/segtree/segtree.hpp"
// #include "alg/monoid/add.hpp"
// #include "ds/static_range_product.hpp"

// template <typename Monoid, typename ST, typename XY, bool SMALL_X, bool SMALL_Y>
// struct Wavelet_Matrix_2D_Range_Static_Monoid {
//   // 点群を Y 昇順に並べる.
//   // X を整数になおして binary trie みたいに振り分ける
//   using MX = Monoid;
//   using X = typename MX::value_type;
//   static_assert(MX::commute);

//   template <bool SMALL>
//   struct TO_IDX {
//     vc<XY> key;
//     XY mi, ma;
//     vc<int> dat;

//     void build(vc<XY>& X) {
//       if constexpr (SMALL) {
//         mi = (X.empty() ? 0 : MIN(X));
//         ma = (X.empty() ? 0 : MAX(X));
//         dat.assign(ma - mi + 2, 0);
//         for (auto& x: X) { dat[x - mi + 1]++; }
//         FOR(i, len(dat) - 1) dat[i + 1] += dat[i];
//       } else {
//         key = X;
//         sort(all(key));
//       }
//     }
//     int operator()(XY x) {
//       if constexpr (SMALL) {
//         return dat[clamp<XY>(x - mi, 0, ma - mi + 1)];
//       } else {
//         return LB(key, x);
//       }
//     }
//   };

//   TO_IDX<SMALL_X> XtoI;
//   TO_IDX<SMALL_Y> YtoI;

//   int N, lg;
//   vector<int> mid;
//   vector<Bit_Vector> bv;
//   vc<int> new_idx;
//   vc<int> A;
//   using SEG = Static_Range_Product<Monoid, ST, 4>;
//   vc<SEG> dat;

//   template <typename F>
//   Wavelet_Matrix_2D_Range_Static_Monoid(int N, F f) {
//     build(N, f);
//   }

//   template <typename F>
//   void build(int N_, F f) {
//     N = N_;
//     if (N == 0) {
//       lg = 0;
//       return;
//     }
//     vc<XY> tmp(N), Y(N);
//     vc<X> S(N);
//     FOR(i, N) tie(tmp[i], Y[i], S[i]) = f(i);
//     auto I = argsort(Y);
//     tmp = rearrange(tmp, I), Y = rearrange(Y, I), S = rearrange(S, I);
//     XtoI.build(tmp), YtoI.build(Y);
//     new_idx.resize(N);
//     FOR(i, N) new_idx[I[i]] = i;

//     // あとは普通に
//     lg = tmp.empty() ? 0 : __lg(XtoI(MAX(tmp) + 1)) + 1;
//     mid.resize(lg), bv.assign(lg, Bit_Vector(N));
//     dat.resize(lg);
//     A.resize(N);
//     FOR(i, N) A[i] = XtoI(tmp[i]);

//     vc<XY> A0(N), A1(N);
//     vc<X> S0(N), S1(N);
//     FOR_R(d, lg) {
//       int p0 = 0, p1 = 0;
//       FOR(i, N) {
//         bool f = (A[i] >> d & 1);
//         if (!f) { S0[p0] = S[i], A0[p0] = A[i], p0++; }
//         if (f) { S1[p1] = S[i], A1[p1] = A[i], bv[d].set(i), p1++; }
//       }
//       mid[d] = p0;
//       bv[d].build();
//       swap(A, A0), swap(S, S0);
//       FOR(i, p1) A[p0 + i] = A1[i], S[p0 + i] = S1[i];
//       dat[d].build(N, [&](int i) -> X { return S[i]; });
//     }
//     FOR(i, N) A[i] = XtoI(tmp[i]);
//   }

//   int count(XY x1, XY x2, XY y1, XY y2) {
//     if (N == 0) return 0;
//     x1 = XtoI(x1), x2 = XtoI(x2);
//     y1 = YtoI(y1), y2 = YtoI(y2);
//     return prefix_count(y1, y2, x2) - prefix_count(y1, y2, x1);
//   }

//   X prod(XY x1, XY x2, XY y1, XY y2) {
//     if (N == 0) return MX::unit();
//     assert(x1 <= x2 && y1 <= y2);
//     x1 = XtoI(x1), x2 = XtoI(x2);
//     y1 = YtoI(y1), y2 = YtoI(y2);
//     X res = MX::unit();
//     prod_dfs(y1, y2, x1, x2, lg - 1, res);
//     return res;
//   }

// private:
//   int prefix_count(int L, int R, int x) {
//     int cnt = 0;
//     FOR_R(d, lg) {
//       int l0 = bv[d].rank(L, 0), r0 = bv[d].rank(R, 0);
//       if (x >> d & 1) {
//         cnt += r0 - l0, L += mid[d] - l0, R += mid[d] - r0;
//       } else {
//         L = l0, R = r0;
//       }
//     }
//     return cnt;
//   }

//   void prod_dfs(int L, int R, int x1, int x2, int d, X& res) {
//     chmax(x1, 0), chmin(x2, 1 << (d + 1));
//     if (x1 >= x2) { return; }
//     assert(0 <= x1 && x1 < x2 && x2 <= (1 << (d + 1)));
//     if (x1 == 0 && x2 == (1 << (d + 1))) {
//       res = MX::op(res, dat[d + 1].prod(L, R));
//       return;
//     }
//     int l0 = bv[d].rank(L, 0), r0 = bv[d].rank(R, 0);
//     prod_dfs(l0, r0, x1, x2, d - 1, res);
//     prod_dfs(L + mid[d] - l0, R + mid[d] - r0, x1 - (1 << d), x2 - (1 << d),
//              d - 1, res);
//   }
// };

package main

import (
	"math/bits"
	"sort"
)

func main() {

}

type XY = int32
type Weight = int32
type Monoid = int

func e() Monoid             { return 0 }
func op(a, b Monoid) Monoid { return a + b }

// y维度升序排列，x维度按照BinaryTrie构建.
type WaveletMatrix2DStaticMonoid struct {
	smallX, smallY bool
	xToIdx, yToIdx *ToIdx
	n, lg          int32
	mid            []int32
	bv             []*bitVector
	dat            [][]Monoid
}

func NewWaveletMatrix2DStaticMonoid(
	n int32, f func(i int32) (x, y XY, w Monoid),
	smallX, smallY bool,
) *WaveletMatrix2DStaticMonoid {
	res := &WaveletMatrix2DStaticMonoid{smallX: smallX, smallY: smallY, xToIdx: NewToIdx(), yToIdx: NewToIdx()}
	res._build(n, f)
	return res
}

func (wm *WaveletMatrix2DStaticMonoid) Count(x1, x2, y1, y2 XY) int32 {
	if wm.n == 0 {
		return 0
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	return wm._prefixCount(y1, y2, x2) - wm._prefixCount(y1, y2, x1)
}

func (wm *WaveletMatrix2DStaticMonoid) Query(x1, x2, y1, y2 XY) Monoid {
	if wm.n == 0 {
		return e()
	}
	if x1 > x2 || y1 > y2 {
		return e()
	}
	x1, x2 = wm.xToIdx.Get(x1, wm.smallX), wm.xToIdx.Get(x2, wm.smallX)
	y1, y2 = wm.yToIdx.Get(y1, wm.smallY), wm.yToIdx.Get(y2, wm.smallY)
	add := wm._prefixSum(y1, y2, x2)
	sub := wm._prefixSum(y1, y2, x1)
	return op(add, inv(sub))
}

func (wm *WaveletMatrix2DStaticMonoid) _build(n int32, f func(i int32) (x, y XY, w Monoid)) {
	wm.n = n
	if n == 0 {
		wm.lg = 0
		return
	}
	A, Y, S := make([]XY, n), make([]XY, n), make([]Monoid, n)
	for i := int32(0); i < n; i++ {
		A[i], Y[i], S[i] = f(i)
	}
	order := argSort(Y)
	A = reArrange(A, order)
	Y = reArrange(Y, order)
	S = reArrange(S, order)
	wm.xToIdx.Build(A, wm.smallX)
	wm.yToIdx.Build(Y, wm.smallY)

	tmp := wm.xToIdx.Get(maxs(A, 0)+1, wm.smallX)
	wm.lg = int32(bits.Len32(uint32(tmp)))
	wm.mid = make([]int32, wm.lg)
	wm.bv = make([]*bitVector, wm.lg)
	for i := range wm.bv {
		wm.bv[i] = newBitVector(n)
	}
	wm.dat = make([][]Monoid, wm.lg+1)
	for i := range wm.dat {
		nums := make([]Monoid, n+1)
		for j := range nums {
			nums[j] = e()
		}
		wm.dat[i] = nums
	}
	for i := int32(0); i < n; i++ {
		A[i] = wm.xToIdx.Get(A[i], wm.smallX)
	}

	A0, A1 := make([]XY, n), make([]XY, n)
	S0, S1 := make([]Monoid, n), make([]Monoid, n)
	for d := wm.lg - 1; d >= -1; d-- {
		p0, p1 := int32(0), int32(0)
		tmp := wm.dat[d+1]
		for i := int32(0); i < n; i++ {
			tmp[i+1] = op(tmp[i], S[i])
		}
		if d == -1 {
			break
		}
		for i := int32(0); i < n; i++ {
			f := (A[i]>>d)&1 == 1
			if !f {
				S0[p0], A0[p0] = S[i], A[i]
				p0++
			} else {
				S1[p1], A1[p1] = S[i], A[i]
				wm.bv[d].Set(i)
				p1++
			}
		}
		wm.mid[d] = p0
		wm.bv[d].Build()
		A, A0 = A0, A
		S, S0 = S0, S
		for i := int32(0); i < p1; i++ {
			A[p0+i], S[p0+i] = A1[i], S1[i]
		}
	}
}

func (wm *WaveletMatrix2DStaticMonoid) _prefixCount(L, R, x int32) int32 {
	cnt := int32(0)
	for d := wm.lg - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(L, false), wm.bv[d].Rank(R, false)
		if (x>>d)&1 == 1 {
			cnt += r0 - l0
			L += wm.mid[d] - l0
			R += wm.mid[d] - r0
		} else {
			L, R = l0, r0
		}
	}
	return cnt
}

func (wm *WaveletMatrix2DStaticMonoid) _prefixSum(L, R, x int32) Monoid {
	add, sub := e(), e()
	for d := wm.lg - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(L, false), wm.bv[d].Rank(R, false)
		if (x>>d)&1 == 1 {
			add = op(add, wm.dat[d][r0])
			sub = op(sub, wm.dat[d][l0])
			L += wm.mid[d] - l0
			R += wm.mid[d] - r0
		} else {
			L, R = l0, r0
		}
	}
	return op(add, inv(sub))
}

type ToIdx struct {
	key    []XY
	mi, ma XY
	dat    []int32
}

func NewToIdx() *ToIdx { return &ToIdx{} }
func (ti *ToIdx) Build(X []XY, small bool) {
	if small {
		ti.mi, ti.ma = mins(X, 0), maxs(X, 0)
		ti.dat = make([]int32, ti.ma-ti.mi+2)
		for _, x := range X {
			ti.dat[x-ti.mi+1]++
		}
		for i := 1; i < len(ti.dat); i++ {
			ti.dat[i] += ti.dat[i-1]
		}
	} else {
		ti.key = make([]XY, len(X))
		copy(ti.key, X)
		sort.Slice(ti.key, func(i, j int) bool { return ti.key[i] < ti.key[j] })
	}
}

func (ti *ToIdx) Get(x XY, small bool) int32 {
	if small {
		return ti.dat[clamp(x-ti.mi, 0, ti.ma-ti.mi+1)]
	} else {
		return lb(ti.key, x)
	}
}

type bitVector struct {
	bits   []uint64
	preSum []int32
}

func newBitVector(n int32) *bitVector {
	return &bitVector{bits: make([]uint64, n>>6+1), preSum: make([]int32, n>>6+1)}
}

func (bv *bitVector) Set(i int32) {
	bv.bits[i>>6] |= 1 << (i & 63)
}

func (bv *bitVector) Build() {
	for i := 0; i < len(bv.bits)-1; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bits[i]))
	}
}

func (bv *bitVector) Rank(k int32, f bool) int32 {
	m, s := bv.bits[k>>6], bv.preSum[k>>6]
	res := s + int32(bits.OnesCount64(m&((1<<(k&63))-1)))
	if f {
		return res
	}
	return k - res
}

func mins(nums []XY, defaultValue XY) XY {
	if len(nums) == 0 {
		return defaultValue
	}
	res := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < res {
			res = nums[i]
		}
	}
	return res
}

func maxs(nums []XY, defaultValue XY) XY {
	if len(nums) == 0 {
		return defaultValue
	}
	res := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > res {
			res = nums[i]
		}
	}
	return res
}

func clamp(x, mi, ma XY) XY {
	if x < mi {
		return mi
	}
	if x > ma {
		return ma
	}
	return x
}

func lb(nums []XY, x XY) int32 {
	left, right := int32(0), int32(len(nums)-1)
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func argSort(nums []XY) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func reArrange[T any](nums []T, order []int32) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}
