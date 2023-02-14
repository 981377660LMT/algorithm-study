# https://nyaannyaan.github.io/library/data-structure/range-union-find.hpp
# 区间并查集 RangeUnionFind


from bisect import bisect_left, bisect_right


class RangeUnionFind:
    ___slots___ = ("_data", "_left", "_right")

    def __init__(self, n: int):
        self._data = [-1] * n
        self._left = list(range(n))  # 每个组的左边界
        self._right = [i + 1 for i in range(n)]  # 每个组的右边界

    def find(self, x: int) -> int:
        if self._data[x] < 0:
            return x
        self._data[x] = self.find(self._data[x])
        return self._data[x]

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._data[rootX] > self._data[rootY]:
            rootX, rootY = rootY, rootX
        self._data[rootX] += self._data[rootY]
        self._data[rootY] = rootX
        if self._left[rootY] < self._left[rootX]:
            self._left[rootX] = self._left[rootY]
        if self._right[rootY] > self._right[rootX]:
            self._right[rootX] = self._right[rootY]
        return True

    def unionRange(self, start: int, end: int) -> int:
        """合并`左闭右开区间[start, end)`，返回新合并的个数(次数)"""
        if start < 0:
            start = 0
        if end > len(self._data):
            end = len(self._data)
        if start >= end:
            return 0
        m, count = 0, 0
        while True:
            m = self._right[self.find(start)]
            if m >= end:
                break
            self.union(start, m)
            count += 1
        return count

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def size(self, x: int) -> int:
        return -self._data[self.find(x)]


if __name__ == "__main__":
    # No.1170 Never Want to Walk
    # https://yukicoder.me/problems/no/1170
    # 数轴上有n个车站,第i个位置在xi
    # 如果两个车站之间的距离di与dj满足 A<=|di-dj|<=B,则这两个车站可以相互到达,否则不能相互到达
    # 对每个车站,求出从该车站出发,可以到达的车站的数量
    # 1<=n<=2e5 0<=A<=B<=1e9 0<=x1<=x2<...<=xn<=1e9

    # !每个车站向右合并可以到达的车站,把合并分解为单点合并+区间合并
    n, A, B = map(int, input().split())
    pos = list(map(int, input().split()))
    uf = RangeUnionFind(n)
    for i, p in enumerate(pos):
        left = bisect_left(pos, p + A)
        right = bisect_right(pos, p + B)
        if left != right:  # 有可以到达的车站
            uf.union(i, left)
            uf.unionRange(left, right)

    for i in range(n):
        print(uf.size(i))

    # uf = RangeUnionFind(10)
    # print(uf.unionRange(0, 5))
    # print(uf.unionRange(0, 5))
    # print(uf.isConnected(0, 4))
    # print(uf.size(0))

    # 2158. 每天绘制新区域的数量
    # https://leetcode-cn.com/problems/amount-of-new-area-painted-each-day/
    # 1 <= paint.length <= 1e5
    # paint[i].length == 2
    # 0 <= starti < endi <= 5 * 104
    from typing import List

    class Solution:
        def amountPainted(self, paint: List[List[int]]) -> List[int]:
            uf = RangeUnionFind(int(5e4) + 10)
            res = [0] * len(paint)
            for i, (start, end) in enumerate(paint):
                res[i] = uf.unionRange(start, end + 1)
            return res
