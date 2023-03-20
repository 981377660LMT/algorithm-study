from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 2 行
# L 列のマス目があります。 上から
# i 行目
# (i∈{1,2})、左から
# j 行目
# (1≤j≤L)のマス目を
# (i,j) で表します。
# (i,j) には整数
# x
# i,j
# ​
#   が書かれています。

# x
# 1,j
# ​
#  =x
# 2,j
# ​
#   であるような整数
# j の個数を求めてください。

# ただし、
# x
# i,j
# ​
#   の情報は
# (x
# 1,1
# ​
#  ,x
# 1,2
# ​
#  ,…,x
# 1,L
# ​
#  ) と
# (x
# 2,1
# ​
#  ,x
# 2,2
# ​
#  ,…,x
# 2,L
# ​
#  ) をそれぞれ連長圧縮した、長さ
# N
# 1
# ​
#   の列
# ((v
# 1,1
# ​
#  ,l
# 1,1
# ​
#  ),…,(v
# 1,N
# 1
# ​

# ​
#  ,l
# 1,N
# 1
# ​

# ​
#  )) と長さ
# N
# 2
# ​
#   の列
# ((v
# 2,1
# ​
#  ,l
# 2,1
# ​
#  ),…,(v
# 2,N
# 2
# ​

# ​
#  ,l
# 2,N
# 2
# ​

# ​
#  )) として与えられます。

# ここで、列
# A の連長圧縮とは、
# A の要素
# v
# i
# ​
#   と正整数
# l
# i
# ​
#   の組
# (v
# i
# ​
#  ,l
# i
# ​
#  ) の列であって、次の操作で得られるものです。

# A を異なる要素が隣り合っている部分で分割する。
# 分割した各列
# B
# 1
# ​
#  ,B
# 2
# ​
#  ,…,B
# k
# ​
#   に対して、
# v
# i
# ​
#   を
# B
# i
# ​
#   の要素、
# l
# i
# ​
#   を
# B
# i
# ​
#   の長さとする。


class BIT2:
    """范围修改"""

    __slots__ = ("size", "_tree1", "_tree2")

    def __init__(self, n: int):
        self.size = n
        self._tree1 = dict()
        self._tree2 = dict()

    def add(self, left: int, right: int, delta: int) -> None:
        """闭区间[left, right]加delta"""
        self._add(left, delta)
        self._add(right + 1, -delta)

    def query(self, left: int, right: int) -> int:
        """闭区间[left, right]的和"""
        return self._query(right) - self._query(left - 1)

    def _add(self, index: int, delta: int) -> None:
        # if index <= 0:
        #    raise ValueError('index 必须是正整数')

        rawIndex = index
        while index <= self.size:
            self._tree1[index] = self._tree1.get(index, 0) + delta
            self._tree2[index] = self._tree2.get(index, 0) + (rawIndex - 1) * delta
            index += index & -index

    def _query(self, index: int) -> int:
        if index > self.size:
            index = self.size

        rawIndex = index
        res = 0
        while index > 0:
            res += rawIndex * self._tree1.get(index, 0) - self._tree2.get(index, 0)
            index -= index & -index
        return res

    def __repr__(self) -> str:
        preSum = []
        for i in range(self.size):
            preSum.append(self._query(i))
        return str(preSum)

    def __len__(self) -> int:
        return self.size


if __name__ == "__main__":
    L, N1, N2 = map(int, input().split())
    # !快速查询区间里某种数有多少个
    bits = defaultdict(lambda: BIT2(int(1e12) + 10))  # 数 => BIT
    start = 1
    for _ in range(N1):
        v, l = map(int, input().split())
        bits[v].add(start, start + l - 1, 1)
        start += l
    res = 0
    start = 1
    for _ in range(N2):
        v, l = map(int, input().split())
        res += bits[v].query(start, start + l - 1)
        start += l
    print(res)
