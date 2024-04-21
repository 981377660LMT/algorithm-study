// WaveletMatrixStatic/StaticWaveletMatrix

// class BitVectorInterface(ABC):
//     @abstractmethod
//     def access(self, k: int) -> int:
//         raise NotImplementedError

//     @abstractmethod
//     def __getitem__(self, k: int) -> int:
//         raise NotImplementedError

//     @abstractmethod
//     def rank0(self, r: int) -> int:
//         raise NotImplementedError

//     @abstractmethod
//     def rank1(self, r: int) -> int:
//         raise NotImplementedError

//     @abstractmethod
//     def rank(self, r: int, v: int) -> int:
//         raise NotImplementedError

//     @abstractmethod
//     def select0(self, k: int) -> int:
//         raise NotImplementedError

//     @abstractmethod
//     def select1(self, k: int) -> int:
//         raise NotImplementedError

//     @abstractmethod
//     def select(self, k: int, v: int) -> int:
//         raise NotImplementedError

//     @abstractmethod
//     def __len__(self) -> int:
//         raise NotImplementedError

//     @abstractmethod
//     def __str__(self) -> str:
//         raise NotImplementedError

//     @abstractmethod
//     def __repr__(self) -> str:
//         raise NotImplementedError

// from array import array

// class BitVector(BitVectorInterface):
//     """コンパクトな bit vector です。"""

//     def __init__(self, n: int):
//         """長さ ``n`` の ``BitVector`` です。

//         bit を保持するのに ``array[I]`` を使用します。
//         ``block_size= n / 32`` として、使用bitは ``32*block_size=2n bit`` です。

//         累積和を保持するのに同様の ``array[I]`` を使用します。
//         32bitごとの和を保存しています。同様に使用bitは ``2n bit`` です。
//         """
//         assert 0 <= n < 4294967295
//         self.N = n
//         self.block_size = (n + 31) >> 5
//         b = bytes(4 * (self.block_size + 1))
//         self.bit = array("I", b)
//         self.acc = array("I", b)

//     @staticmethod
//     def _popcount(x: int) -> int:
//         x = x - ((x >> 1) & 0x55555555)
//         x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
//         x = x + (x >> 4) & 0x0F0F0F0F
//         x += x >> 8
//         x += x >> 16
//         return x & 0x0000007F

//     def set(self, k: int) -> None:
//         """``k`` 番目の bit を ``1`` にします。
//         :math:`O(1)` です。

//         Args:
//           k (int): インデックスです。
//         """
//         self.bit[k >> 5] |= 1 << (k & 31)

//     def build(self) -> None:
//         """構築します。
//         **これ以降 ``set`` メソッドを使用してはいけません。**
//         :math:`O(n)` です。
//         """
//         acc, bit = self.acc, self.bit
//         for i in range(self.block_size):
//             acc[i + 1] = acc[i] + BitVector._popcount(bit[i])

//     def access(self, k: int) -> int:
//         """``k`` 番目の bit を返します。
//         :math:`O(1)` です。
//         """
//         return (self.bit[k >> 5] >> (k & 31)) & 1

//     def __getitem__(self, k: int) -> int:
//         return (self.bit[k >> 5] >> (k & 31)) & 1

//     def rank0(self, r: int) -> int:
//         """``a[0, r)`` に含まれる ``0`` の個数を返します。
//         :math:`O(1)` です。
//         """
//         return r - (
//             self.acc[r >> 5] + BitVector._popcount(self.bit[r >> 5] & ((1 << (r & 31)) - 1))
//         )

//     def rank1(self, r: int) -> int:
//         """``a[0, r)`` に含まれる ``1`` の個数を返します。
//         :math:`O(1)` です。
//         """
//         return self.acc[r >> 5] + BitVector._popcount(self.bit[r >> 5] & ((1 << (r & 31)) - 1))

//     def rank(self, r: int, v: int) -> int:
//         """``a[0, r)`` に含まれる ``v`` の個数を返します。
//         :math:`O(1)` です。
//         """
//         return self.rank1(r) if v else self.rank0(r)

//     def select0(self, k: int) -> int:
//         """``k`` 番目の ``0`` のインデックスを返します。
//         :math:`O(\\log{n})` です。
//         """
//         if k < 0 or self.rank0(self.N) <= k:
//             return -1
//         l, r = 0, self.block_size + 1
//         while r - l > 1:
//             m = (l + r) >> 1
//             if m * 32 - self.acc[m] > k:
//                 r = m
//             else:
//                 l = m
//         indx = 32 * l
//         k = k - (l * 32 - self.acc[l]) + self.rank0(indx)
//         l, r = indx, indx + 32
//         while r - l > 1:
//             m = (l + r) >> 1
//             if self.rank0(m) > k:
//                 r = m
//             else:
//                 l = m
//         return l

//     def select1(self, k: int) -> int:
//         """``k`` 番目の ``1`` のインデックスを返します。
//         :math:`O(\\log{n})` です。
//         """
//         if k < 0 or self.rank1(self.N) <= k:
//             return -1
//         l, r = 0, self.block_size + 1
//         while r - l > 1:
//             m = (l + r) >> 1
//             if self.acc[m] > k:
//                 r = m
//             else:
//                 l = m
//         indx = 32 * l
//         k = k - self.acc[l] + self.rank1(indx)
//         l, r = indx, indx + 32
//         while r - l > 1:
//             m = (l + r) >> 1
//             if self.rank1(m) > k:
//                 r = m
//             else:
//                 l = m
//         return l

//     def select(self, k: int, v: int) -> int:
//         """``k`` 番目の ``v`` のインデックスを返します。
//         :math:`O(\\log{n})` です。
//         """
//         return self.select1(k) if v else self.select0(k)

//     def __len__(self):
//         return self.N

//     def __str__(self):
//         return str([self.access(i) for i in range(self.N)])

//     def __repr__(self):
//         return f"{self.__class__.__name__}({self})"

// from typing import Sequence, List, Tuple
// from heapq import heappush, heappop
// from array import array

// class WaveletMatrix:
//     """``WaveletMatrix`` です。
//     静的であることに注意してください。

//     以下の仕様の計算量には嘘があるかもしれません。import 元の ``BitVector`` の計算量も参考にしてください。

//     - 参考:
//       - `https://miti-7.hatenablog.com/entry/2018/04/28/152259 <https://miti-7.hatenablog.com/entry/2018/04/28/152259>`
//       - `https://www.slideshare.net/pfi/ss-15916040 <https://www.slideshare.net/pfi/ss-15916040>`
//       - `デwiki <https://scrapbox.io/data-structures/Wavelet_Matrix>`
//     """

//     def __init__(self, sigma: int, a: Sequence[int] = []):
//         """``[0, sigma)`` の整数列を管理する ``WaveletMatrix`` を構築します。
//         :math:`O(n\\log{\\sigma})` です。

//         Args:
//           sigma (int): 扱う整数の上限です。
//           a (Sequence[int], optional): 構築する配列です。
//         """
//         self.sigma: int = sigma
//         self.log: int = (sigma - 1).bit_length()
//         self.mid: array[int] = array("I", bytes(4 * self.log))
//         self.size: int = len(a)
//         self.v: List[BitVector] = [BitVector(self.size) for _ in range(self.log)]
//         self._build(a)

//     def _build(self, a: Sequence[int]) -> None:
//         # 列 a から wm を構築する
//         for bit in range(self.log - 1, -1, -1):
//             # bit目の0/1に応じてvを構築 + aを安定ソート
//             v = self.v[bit]
//             zero, one = [], []
//             for i, e in enumerate(a):
//                 if e >> bit & 1:
//                     v.set(i)
//                     one.append(e)
//                 else:
//                     zero.append(e)
//             v.build()
//             self.mid[bit] = len(zero)  # 境界をmid[bit]に保持
//             a = zero + one

//     def access(self, k: int) -> int:
//         """k番目の値を返します。
//         :math:`O(\\log{\\sigma})` です。

//         Args:
//           k (int): インデックスです。
//         """
//         assert (
//             -self.size <= k < self.size
//         ), f"IndexError: {self.__class__.__name__}.access({k}), size={self.size}"
//         if k < 0:
//             k += self.size
//         s = 0  # 答え
//         for bit in range(self.log - 1, -1, -1):
//             if self.v[bit].access(k):
//                 # k番目が立ってたら、
//                 # kまでの1とすべての0が次のk
//                 s |= 1 << bit
//                 k = self.v[bit].rank1(k) + self.mid[bit]
//             else:
//                 # kまでの0が次のk
//                 k = self.v[bit].rank0(k)
//         return s

//     def __getitem__(self, k: int) -> int:
//         assert (
//             -self.size <= k < self.size
//         ), f"IndexError: {self.__class__.__name__}[{k}], size={self.size}"
//         return self.access(k)

//     def rank(self, r: int, x: int) -> int:
//         """``a[0, r)`` に含まれる ``x`` の個数を返します。
//         :math:`O(\\log{\\sigma})` です。
//         """
//         assert (
//             0 <= r <= self.size
//         ), f"IndexError: {self.__class__.__name__}.rank(), r={r}, size={self.size}"
//         assert (
//             0 <= x < 1 << self.log
//         ), f"ValueError: {self.__class__.__name__}.rank(), x={x}, LIM={1<<self.log}"
//         l = 0
//         mid = self.mid
//         for bit in range(self.log - 1, -1, -1):
//             # 位置 r より左に x が何個あるか
//             # x の bit 目で場合分け
//             if x >> bit & 1:
//                 # 立ってたら、次のl, rは以下
//                 l = self.v[bit].rank1(l) + mid[bit]
//                 r = self.v[bit].rank1(r) + mid[bit]
//             else:
//                 # そうでなければ次のl, rは以下
//                 l = self.v[bit].rank0(l)
//                 r = self.v[bit].rank0(r)
//         return r - l

//     def select(self, k: int, x: int) -> int:
//         """``k`` 番目の ``v`` のインデックスを返します。
//         :math:`O(\\log{\\sigma})` です。
//         """
//         assert (
//             0 <= k < self.size
//         ), f"IndexError: {self.__class__.__name__}.select({k}, {x}), k={k}, size={self.size}"
//         assert (
//             0 <= x < 1 << self.log
//         ), f"ValueError: {self.__class__.__name__}.select({k}, {x}), x={x}, LIM={1<<self.log}"
//         # x の開始位置 s を探す
//         s = 0
//         for bit in range(self.log - 1, -1, -1):
//             if x >> bit & 1:
//                 s = self.v[bit].rank0(self.size) + self.v[bit].rank1(s)
//             else:
//                 s = self.v[bit].rank0(s)
//         s += k  # s から k 進んだ位置が、元の列で何番目か調べる
//         for bit in range(self.log):
//             if x >> bit & 1:
//                 s = self.v[bit].select1(s - self.v[bit].rank0(self.size))
//             else:
//                 s = self.v[bit].select0(s)
//         return s

//     def kth_smallest(self, l: int, r: int, k: int) -> int:
//         """``a[l, r)`` の中で k 番目に **小さい** 値を返します。
//         :math:`O(\\log{\\sigma})` です。
//         """
//         assert (
//             0 <= l <= r <= self.size
//         ), f"IndexError: {self.__class__.__name__}.kth_smallest({l}, {r}, {k}), size={self.size}"
//         assert (
//             0 <= k < r - l
//         ), f"IndexError: {self.__class__.__name__}.kth_smallest({l}, {r}, {k}), wrong k"
//         s = 0
//         mid = self.mid
//         for bit in range(self.log - 1, -1, -1):
//             r0, l0 = self.v[bit].rank0(r), self.v[bit].rank0(l)
//             cnt = r0 - l0  # 区間内の 0 の個数
//             if cnt <= k:  # 0 が k 以下のとき、 k 番目は 1
//                 s |= 1 << bit
//                 k -= cnt
//                 # この 1 が次の bit 列でどこに行くか
//                 l = l - l0 + mid[bit]
//                 r = r - r0 + mid[bit]
//             else:
//                 # この 0 が次の bit 列でどこに行くか
//                 l = l0
//                 r = r0
//         return s

//     quantile = kth_smallest

//     def kth_largest(self, l: int, r: int, k: int) -> int:
//         """``a[l, r)`` の中で k 番目に **大きい値** を返します。
//         :math:`O(\\log{\\sigma})` です。
//         """
//         assert (
//             0 <= l <= r <= self.size
//         ), f"IndexError: {self.__class__.__name__}.kth_largest({l}, {r}, {k}), size={self.size}"
//         assert (
//             0 <= k < r - l
//         ), f"IndexError: {self.__class__.__name__}.kth_largest({l}, {r}, {k}), wrong k"
//         return self.kth_smallest(l, r, r - l - k - 1)

//     def topk(self, l: int, r: int, k: int) -> List[Tuple[int, int]]:
//         """``a[l, r)`` の中で、要素を出現回数が多い順にその頻度とともに ``k`` 個返します。
//         :math:`O(\\min(r-l, \\sigam) \\log(\\sigam))` です。

//         Note:
//           :math:`\\sigma` が大きい場合、計算量に注意です。

//         Returns:
//           List[Tuple[int, int]]: ``(要素, 頻度)`` を要素とする配列です。
//         """
//         assert (
//             0 <= l <= r <= self.size
//         ), f"IndexError: {self.__class__.__name__}.topk({l}, {r}, {k}), size={self.size}"
//         assert 0 <= k < r - l, f"IndexError: {self.__class__.__name__}.topk({l}, {r}, {k}), wrong k"
//         # heap[-length, x, l, bit]
//         hq: List[Tuple[int, int, int, int]] = [(-(r - l), 0, l, self.log - 1)]
//         ans = []
//         while hq:
//             length, x, l, bit = heappop(hq)
//             length = -length
//             if bit == -1:
//                 ans.append((x, length))
//                 k -= 1
//                 if k == 0:
//                     break
//             else:
//                 r = l + length
//                 l0 = self.v[bit].rank0(l)
//                 r0 = self.v[bit].rank0(r)
//                 if l0 < r0:
//                     heappush(hq, (-(r0 - l0), x, l0, bit - 1))
//                 l1 = self.v[bit].rank1(l) + self.mid[bit]
//                 r1 = self.v[bit].rank1(r) + self.mid[bit]
//                 if l1 < r1:
//                     heappush(hq, (-(r1 - l1), x | (1 << bit), l1, bit - 1))
//         return ans

//     def sum(self, l: int, r: int) -> int:
//         """``topk`` メソッドを用いて ``a[l, r)`` の総和を返します。

//         計算量に注意です。
//         """
//         assert False, "Yabai Keisanryo Error"
//         return sum(k * v for k, v in self.topk(l, r, r - l))

//     def _range_freq(self, l: int, r: int, x: int) -> int:
//         """a[l, r) で x 未満の要素の数を返す"""
//         ans = 0
//         for bit in range(self.log - 1, -1, -1):
//             l0, r0 = self.v[bit].rank0(l), self.v[bit].rank0(r)
//             if x >> bit & 1:
//                 # bit が立ってたら、区間の 0 の個数を答えに加算し、新たな区間は 1 のみ
//                 ans += r0 - l0
//                 # 1 が次の bit 列でどこに行くか
//                 l += self.mid[bit] - l0
//                 r += self.mid[bit] - r0
//             else:
//                 # 0 が次の bit 列でどこに行くか
//                 l, r = l0, r0
//         return ans

//     def range_freq(self, l: int, r: int, x: int, y: int) -> int:
//         """``a[l, r)`` に含まれる、 ``x`` 以上 ``y`` 未満である要素の個数を返します。
//         :math:`O(\\log{\\sigma})` です。
//         """
//         assert (
//             0 <= l <= r <= self.size
//         ), f"IndexError: {self.__class__.__name__}.range_freq({l}, {r}, {x}, {y})"
//         return self._range_freq(l, r, y) - self._range_freq(l, r, x)

//     def prev_value(self, l: int, r: int, x: int) -> int:
//         """``a[l, r)`` で、``x`` 以上 ``y`` 未満であるような要素のうち最大の要素を返します。
//         :math:`O(\\log{\\sigma})` です。
//         """
//         assert (
//             0 <= l <= r <= self.size
//         ), f"IndexError: {self.__class__.__name__}.prev_value({l}, {r}, {x})"
//         return self.kth_smallest(l, r, self._range_freq(l, r, x) - 1)

//     def next_value(self, l: int, r: int, x: int) -> int:
//         """``a[l, r)`` で、``x`` 以上 ``y`` 未満であるような要素のうち最小の要素を返します。
//         :math:`O(\\log{\\sigma})` です。
//         """
//         assert (
//             0 <= l <= r <= self.size
//         ), f"IndexError: {self.__class__.__name__}.next_value({l}, {r}, {x})"
//         return self.kth_smallest(l, r, self._range_freq(l, r, x))

//     def range_count(self, l: int, r: int, x: int) -> int:
//         """``a[l, r)`` に含まれる ``x`` の個数を返します。
//         ``wm.rank(r, x) - wm.rank(l, x)`` と等価です。
//         :math:`O(\\log{\\sigma})` です。
//         """
//         assert (
//             0 <= l <= r <= self.size
//         ), f"IndexError: {self.__class__.__name__}.range_count({l}, {r}, {x})"
//         return self.rank(r, x) - self.rank(l, x)

//     def __len__(self):
//         return self.size

//     def __str__(self):
//         return f"{self.__class__.__name__}({[self.access(i) for i in range(self.size)]})"

//     __repr__ = __str__

package main

import "math/bits"

func main() {

}

type WaveletMatrixStatic struct {
}

type bitVector struct {
	n      int32
	size   int32
	bit    []uint64
	preSum []int32
}

func newBitVector(n int32) *bitVector {
	size := (n + 63) >> 6
	bit := make([]uint64, size+1)
	preSum := make([]int32, size+1)
	return &bitVector{n: n, size: size, bit: bit, preSum: preSum}
}

func (bv *bitVector) Set(i int32) {
	bv.bit[i>>6] |= 1 << (i & 63)
}

func (bv *bitVector) Build() {
	for i := int32(0); i < bv.size; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bit[i]))
	}
}

func (bv *bitVector) Get(i int32) int32 {
	return int32(bv.bit[i>>6] >> (i & 63) & 1)
}

func (bv *bitVector) Count0(end int32) int32 {
	return end - (bv.preSum[end>>6] + int32(bits.OnesCount64(bv.bit[end>>6]&(1<<(end&63)-1))))
}

func (bv *bitVector) Count1(end int32) int32 {
	return bv.preSum[end>>6] + int32(bits.OnesCount64(bv.bit[end>>6]&(1<<(end&63)-1)))
}

func (bv *bitVector) Count(end int32, value int32) int32 {
	if value == 1 {
		return bv.Count1(end)
	}
	return end - bv.Count1(end)
}

func (bv *bitVector) Kth0(k int32) int32 {
	if k < 0 || bv.Count0(bv.n) <= k {
		return -1
	}
	l, r := int32(0), bv.size+1
	for r-l > 1 {
		m := (l + r) >> 1
		if m<<6-bv.preSum[m] > k {
			r = m
		} else {
			l = m
		}
	}
	indx := l << 6
	k -= (l<<6 - bv.preSum[l]) - bv.Count0(indx)
	l, r = indx, indx+64
	for r-l > 1 {
		m := (l + r) >> 1
		if bv.Count0(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}

// k>=0.
func (bv *bitVector) Kth1(k int32) int32 {
	if k < 0 || bv.Count1(bv.n) <= k {
		return -1
	}
	l, r := int32(0), bv.size+1
	for r-l > 1 {
		m := (l + r) >> 1
		if bv.preSum[m] > k {
			r = m
		} else {
			l = m
		}
	}
	indx := l << 6
	k -= bv.preSum[l] - bv.Count1(indx)
	l, r = indx, indx+64
	for r-l > 1 {
		m := (l + r) >> 1
		if bv.Count1(m) > k {
			r = m
		} else {
			l = m
		}
	}
	return l
}

func (bv *bitVector) Kth(k int32, v int32) int32 {
	if v == 1 {
		return bv.Kth1(k)
	}
	return bv.Kth0(k)
}

func (bv *bitVector) GetAll() []int32 {
	res := make([]int32, 0, bv.n)
	for i := int32(0); i < bv.n; i++ {
		res = append(res, bv.Get(i))
	}
	return res
}
