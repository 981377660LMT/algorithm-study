# 区间并查集 RangeUnionFind/UnionFindRange
# !只使用了路径压缩,每次操作O(logn)
# !上位替代：SegmentSet/线段树.

from typing import Callable, DefaultDict, Optional, List
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

    def union(
        self, x: int, y: int, beforeUnion: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        """union后, 大的编号所在的组的指向小的编号所在的组(向左合并)."""
        if x < y:
            x, y = y, x
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if beforeUnion is not None:
            beforeUnion(rootY, rootX)
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        return True

    def unionRange(
        self, left: int, right: int, beforeUnion: Optional[Callable[[int, int], None]] = None
    ) -> int:
        """合并[left,right]区间, 返回合并次数."""
        if left >= right:
            return 0
        leftRoot = self.find(left)
        rightRoot = self.find(right)
        unionCount = 0
        while rightRoot != leftRoot:
            unionCount += 1
            self.union(rightRoot, rightRoot - 1, beforeUnion)
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

    from typing import List

    class Solution:
        # 2158. 每天绘制新区域的数量
        # https://leetcode-cn.com/problems/amount-of-new-area-painted-each-day/
        # 1 <= paint.length <= 1e5
        # paint[i].length == 2
        # 0 <= starti < endi <= 5 * 104
        def amountPainted(self, paint: List[List[int]]) -> List[int]:
            uf = UnionFindRange(int(5e4) + 10)
            return [uf.unionRange(start, end) for start, end in paint]

        # !100376. 新增道路查询后的最短距离 II (注意这里合并的是边，右端点要减去1)
        # https://leetcode.cn/problems/shortest-distance-after-road-addition-queries-ii/description/
        # 给你一个整数 n 和一个二维整数数组 queries。
        # 有 n 个城市，编号从 0 到 n - 1。初始时，每个城市 i 都有一条单向道路通往城市 i + 1（ 0 <= i < n - 1）。
        # !queries[i] = [ui, vi] 表示新建一条从城市 ui 到城市 vi 的单向道路。每次查询后，你需要找到从城市 0 到城市 n - 1 的最短路径的长度。
        # 所有查询中不会存在两个查询都满足 queries[i][0] < queries[j][0] < queries[i][1] < queries[j][1]。
        # 返回一个数组 answer，对于范围 [0, queries.length - 1] 中的每个 i，answer[i] 是处理完前 i + 1 个查询后，从城市 0 到城市 n - 1 的最短路径的长度。
        def shortestDistanceAfterQueries(self, n: int, queries: List[List[int]]) -> List[int]:
            uf = UnionFindRange(n - 1)
            res = []
            for u, v in queries:
                uf.unionRange(u, v - 1)
                res.append(uf.part)
            return res

        # 3532. 针对图的路径存在性查询 I
        # https://leetcode.cn/problems/path-existence-queries-in-a-graph-i/
        def pathExistenceQueries(
            self, n: int, nums: List[int], maxDiff: int, queries: List[List[int]]
        ) -> List[bool]:
            n = len(nums)
            uf = UnionFindRange(n)
            ptr = 0
            while ptr < len(nums):
                start = ptr
                ptr += 1
                while ptr < len(nums) and nums[ptr] - nums[ptr - 1] <= maxDiff:
                    ptr += 1
                uf.unionRange(start, ptr - 1)
            return [uf.isConnected(u, v) for u, v in queries]
