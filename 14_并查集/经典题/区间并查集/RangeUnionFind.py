# 区间并查集 RangeUnionFind/UnionFindRange
# !只使用了路径压缩,每次操作O(logn)

from typing import Callable, DefaultDict, Optional, List, Tuple
from bisect import bisect_left, bisect_right
from collections import defaultdict


class UnionFindRange:
    __slots__ = "part", "_n", "_parent", "_rank"

    def __init__(self, n: int):
        self.part = n
        self._n = n
        self._parent = list(range(n))
        self._rank = [1] * n

    def find(self, x: int) -> int:
        while x != self._parent[x]:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int, f: Optional[Callable[[int, int], None]] = None) -> bool:
        """union后, 大的编号所在的组的指向小的编号所在的组(向左合并)."""
        if x < y:
            x, y = y, x
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        if f is not None:
            f(rootY, rootX)
        return True

    def unionRange(
        self, left: int, right: int, f: Optional[Callable[[int, int], None]] = None
    ) -> int:
        """合并[left,right]区间, 返回合并次数."""
        if left >= right:
            return 0
        leftRoot = self.find(left)
        rightRoot = self.find(right)
        unionCount = 0
        while rightRoot != leftRoot:
            unionCount += 1
            self.union(rightRoot, rightRoot - 1, f)
            rightRoot = self.find(rightRoot - 1)
        return unionCount

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getSize(self, x: int) -> int:
        return self._rank[self.find(x)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        group = defaultdict(list)
        for i in range(self._n):
            group[self.find(i)].append(i)
        return group


class UnionFindRange2:
    """
    维护每个分组左右边界的区间并查集.
    按秩合并.
    """

    ___slots___ = ("groupStart", "groupEnd", "part", "_n", "_data")

    def __init__(self, n: int):
        self.groupStart = list(range(n))  # 每个组的左边界,包含端点
        self.groupEnd = [i + 1 for i in range(n)]  # 每个组的右边界,不包含端点
        self.part = n
        self._n = n
        self._data = [-1] * n

    def unionRange(self, start: int, end: int) -> int:
        """合并[start,end)左闭右开区间, 返回合并次数.
        0<=groupStart<=groupEnd<=n.
        """
        if start < 0:
            start = 0
        if end > self._n:
            end = self._n
        if start >= end:
            return 0
        count = 0
        while True:
            next_ = self.groupEnd[self.find(start)]
            if next_ >= end:
                break
            self.union(start, next_)
            count += 1
        return count

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
        if self.groupStart[rootY] < self.groupStart[rootX]:
            self.groupStart[rootX] = self.groupStart[rootY]
        if self.groupEnd[rootY] > self.groupEnd[rootX]:
            self.groupEnd[rootX] = self.groupEnd[rootY]
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getSize(self, x: int) -> int:
        return -self._data[self.find(x)]

    def getRange(self, x: int) -> Tuple[int, int]:
        """每个点所在的组的左右边界[左边界,右边界)."""
        root = self.find(x)
        return self.groupStart[root], self.groupEnd[root]

    def getGroups(self) -> "DefaultDict[int, List[int]]":
        group = defaultdict(list)
        for i in range(self._n):
            group[self.find(i)].append(i)
        return group


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
    uf = UnionFindRange(n)
    for i, p in enumerate(pos):
        left = bisect_left(pos, p + A)
        right = bisect_right(pos, p + B) - 1
        if right >= left:  # 有可以到达的车站
            uf.union(i, left)
            uf.unionRange(left, right)
    for i in range(n):
        print(uf.getSize(i))

    # 2158. 每天绘制新区域的数量
    # https://leetcode-cn.com/problems/amount-of-new-area-painted-each-day/
    # 1 <= paint.length <= 1e5
    # paint[i].length == 2
    # 0 <= starti < endi <= 5 * 104
    from typing import List

    class Solution:
        def amountPainted(self, paint: List[List[int]]) -> List[int]:
            uf = UnionFindRange2(int(5e4) + 10)
            return [uf.unionRange(start, end + 1) for start, end in paint]
