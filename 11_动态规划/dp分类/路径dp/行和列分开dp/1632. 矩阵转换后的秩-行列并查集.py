# https://leetcode.cn/problems/rank-transform-of-a-matrix/
# 1632. 矩阵转换后的秩
# 给你一个 m x n 的矩阵 matrix ，请你返回一个新的矩阵 answer ，
# 其中 answer[row][col] 是 matrix[row][col] 的秩。
# 每个元素的 秩 是一个整数，表示这个元素相对于其他元素的大小关系，它按照如下规则计算：
# 秩是从 1 开始的一个整数。
# 如果两个元素 p 和 q 在 同一行 或者 同一列 ，那么：
# 如果 p < q ，那么 rank(p) < rank(q)
# 如果 p == q ，那么 rank(p) == rank(q)
# 如果 p > q ，那么 rank(p) > rank(q)
# 秩 需要越 小 越好。
# !ROW*COL<=1e5


# !注意存在相等元素
# !1.如果不存在相等元素:从小到大遍历元素，并维护每行/列的最大秩，该元素的秩即为同行/列的最大秩加 1
# !2.处理相等元素:假设两个相同元素同行/列，那么就要考虑到两个元素分别对应的 行/列 的最大秩.还可能出现联动.
#    需要行列并查集维护

from collections import defaultdict
from typing import DefaultDict, List


class Solution:
    def matrixRankTransform(self, matrix: List[List[int]]) -> List[List[int]]:
        ROW, COL = len(matrix), len(matrix[0])
        mp = defaultdict(list)
        for i, row in enumerate(matrix):
            for j, num in enumerate(row):
                mp[num].append((i, j))  # 相同值的元素一起处理

        dp = defaultdict(int)  # (x,y)处的值
        rowDp = [0] * ROW  # 当前行最大值
        colDp = [0] * COL  # 当前列最大值

        keys = sorted(mp)

        for key in keys:
            pos = mp[key]
            uf = UnionFindMap()
            groupMax = defaultdict(int)
            for x, y in pos:  # !处理相等元素,序号考虑相等元素对应行/列的最大值
                uf.union(x, y + ROW)  # !行列并查集,0-ROW-1:行,ROW-ROW+COL-1:列
            for x, y in pos:
                root = uf.find(x)
                groupMax[root] = max(groupMax[root], rowDp[x], colDp[y])

            for x, y in pos:
                dp[(x, y)] = 1 + groupMax[uf.find(x)]
            for x, y in pos:
                rowDp[x] = max(rowDp[x], dp[(x, y)])
                colDp[y] = max(colDp[y], dp[(x, y)])

        return [[dp[(x, y)] for y in range(COL)] for x in range(ROW)]


from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


T = TypeVar("T", bound=Hashable)


class UnionFindMap(Generic[T]):
    """当元素不是数组index时(例如字符串),更加通用的并查集写法,支持动态添加"""

    __slots__ = ("part", "_parent", "_rank")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self._parent = dict()
        self._rank = dict()
        for item in iterable or []:
            self.add(item)

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            return False
        if self._rank[root1] > self._rank[root2]:
            root1, root2 = root2, root1
        self._parent[root1] = root2
        self._rank[root2] += self._rank[root1]
        self.part -= 1
        return True

    def find(self, key: T) -> T:
        if key not in self._parent:
            self.add(key)
            return key

        while self._parent.get(key, key) != key:
            self._parent[key] = self._parent[self._parent[key]]
            key = self._parent[key]
        return key

    def isConnected(self, key1: T, key2: T) -> bool:
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self._parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self._parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def add(self, key: T) -> bool:
        if key in self._parent:
            return False
        self._parent[key] = key
        self._rank[key] = 1
        self.part += 1
        return True

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self._parent
