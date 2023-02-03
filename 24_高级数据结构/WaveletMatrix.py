from typing import List, Tuple


class BitVector:
    __slots__ = ("block", "bit", "cnt")

    def __init__(self, size: int) -> None:
        self.block = (size + 31) >> 5
        self.bit = [0] * self.block
        self.cnt = [0] * self.block

    def set(self, i: int) -> None:
        self.bit[i >> 5] |= 1 << (i & 31)

    def build(self) -> None:
        for i in range(self.block - 1):
            self.cnt[i + 1] = self.cnt[i] + self.popcount(self.bit[i])

    def popcount(self, x: int) -> int:
        x = x - ((x >> 1) & 0x55555555)
        x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
        x = (x + (x >> 4)) & 0x0F0F0F0F
        x = x + (x >> 8)
        x = x + (x >> 16)
        return x & 0x0000007F

    def count1(self, r: int) -> int:
        msk = (1 << (r & 31)) - 1
        return self.cnt[r >> 5] + self.popcount(self.bit[r >> 5] & msk)

    def rank0(self, r: int) -> int:
        return r - self.count1(r)


class WaveletMatrix:
    def __init__(self, array: List[int], log: int = 32) -> None:
        self.n = len(array)
        self.mat = []
        self.zs = []
        self.log = log
        for d in range(self.log)[::-1]:
            ls, rs = [], []
            BV = BitVector(self.n + 1)
            for ind, val in enumerate(array):
                if val & (1 << d):
                    rs.append(val)
                    BV.set(ind)
                else:
                    ls.append(val)
            BV.build()
            self.mat.append(BV)
            self.zs.append(len(ls))
            array = ls + rs

    def access(self, i: int) -> int:
        res = 0
        for d in range(self.log):
            res <<= 1
            if self.mat[d][i]:
                res |= 1
                i = self.mat[d].rank1(i) + self.zs[d]
            else:
                i = self.mat[d].rank0(i)
        return res

    def rank(self, val: int, l: int, r: int) -> int:
        """[l,r)中val出现次数"""
        for d in range(self.log):
            if val >> (self.log - d - 1) & 1:
                l = self.mat[d].rank1(l) + self.zs[d]
                r = self.mat[d].rank1(r) + self.zs[d]
            else:
                l = self.mat[d].rank0(l)
                r = self.mat[d].rank0(r)
        return r - l

    def quantile(self, l: int, r: int, k: int) -> int:
        res = 0
        for d in range(self.log):
            res <<= 1
            cntl, cntr = self.mat[d].rank1(l), self.mat[d].rank1(r)
            if cntr - cntl >= k:
                l = cntl + self.zs[d]
                r = cntr + self.zs[d]
                res |= 1
            else:
                l -= cntl
                r -= cntr
                k -= cntr - cntl
        return res

    def kth_smallest(self, l: int, r: int, k: int) -> int:
        """区间[l,r)第k小 0<=l<=r<=n"""
        return self.quantile(l, r, r - l - k)


# 区间某个数出现次数  rank (更好的方法是邻接表+二分)
# 区间第k小 kth_smallest

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    # # 区间第k小
    n, q = map(int, input().split())
    nums = list(map(int, input().split()))
    wm = WaveletMatrix(nums)
    for _ in range(q):
        left, right, k = map(int, input().split())
        print(wm.kth_smallest(left, right, k))

    # # 区间频率
    # n, q = map(int, input().split())
    # nums = list(map(int, input().split()))
    # wm = WaveletMatrix(nums)
    # for _ in range(q):
    #     left, right, x = map(int, input().split())
    #     print(wm.rank(x, left, right))
