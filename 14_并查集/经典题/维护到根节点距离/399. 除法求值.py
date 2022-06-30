# https://leetcode.cn/problems/evaluate-division/submissions/

from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


T = TypeVar("T", bound=Hashable)


class UnionFindMapWithDist(Generic[T]):
    """需要手动添加元素 维护距离的并查集"""

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self.parent = dict()
        self.distToRoot = defaultdict(lambda: 1.0)
        for item in iterable or []:
            self.add(item)

    def add(self, key: T) -> "UnionFindMapWithDist[T]":
        if key in self.parent:
            return self
        self.parent[key] = key
        self.part += 1
        return self

    def union(self, son: T, father: T, dist: float) -> bool:
        """
        father 与 son 间的距离为 dist
        围绕着'到根的距离'进行计算
        注意从走两条路到新根节点的距离是一样的
        """
        root1 = self.find(son)
        root2 = self.find(father)
        if (root1 == root2) or (root1 not in self.parent) or (root2 not in self.parent):
            return False

        self.parent[root1] = root2
        # !1. 合并距离 保持两条路到新根节点的距离是一样的
        self.distToRoot[root1] = dist * self.distToRoot[father] / self.distToRoot[son]
        self.part -= 1
        return True

    def find(self, key: T) -> T:
        """此处不自动add"""
        if key not in self.parent:
            return key

        # !2. 从上往下懒更新到根的距离
        if key != self.parent[key]:
            root = self.find(self.parent[key])
            self.distToRoot[key] *= self.distToRoot[self.parent[key]]
            self.parent[key] = root
        return self.parent[key]

    def isConnected(self, key1: T, key2: T) -> bool:
        if (key1 not in self.parent) or (key2 not in self.parent):
            return False
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self.parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self.parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def __str__(self) -> str:
        return "\n".join(
            f"{root}: {member}" for root, member in self.getGroups().items()
        )

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self.parent


class Solution:
    def calcEquation(
        self, equations: List[List[str]], values: List[float], queries: List[List[str]]
    ) -> List[float]:
        """如果存在某个无法确定的答案，则用 -1.0 替代这个答案。
        如果问题中出现了给定的已知条件中没有出现的字符串，也需要用 -1.0 替代这个答案。
        乘积关系取对数就是加法 等价于维护到根节点的距离
        """
        uf = UnionFindMapWithDist[str]()
        for (key1, key2), value in zip(equations, values):
            uf.add(key1).add(key2).union(key2, key1, value)  # !value * key2 = key1

        res = []
        for u, v in queries:
            if u not in uf or v not in uf or not uf.isConnected(u, v):
                res.append(-1.0)
            else:
                res.append(uf.distToRoot[v] / uf.distToRoot[u])

        return res


# print(
#     Solution().calcEquation(
#         [["a", "b"], ["b", "c"]],
#         [2.0, 3.0],
#         [["a", "c"], ["b", "a"], ["a", "e"], ["a", "a"], ["x", "x"]],
#     )
# )
# print(
#     Solution().calcEquation(
#         equations=[["a", "b"], ["b", "c"], ["bc", "cd"]],
#         values=[1.5, 2.5, 5.0],
#         queries=[["a", "c"], ["c", "b"], ["bc", "cd"], ["cd", "bc"]],
#     )
# )
print(
    Solution().calcEquation(
        equations=[["a", "e"], ["b", "e"]],
        values=[4.0, 3.0],
        queries=[["a", "b"], ["e", "e"], ["x", "x"]],
    )
)
# [1.33333,1.0,-1.0]
