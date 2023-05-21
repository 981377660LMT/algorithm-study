# https://leetcode.cn/problems/check-for-contradictions-in-equations/
# 2307. 检查方程中的矛盾之处
# 检查方程是否存在矛盾
# 如果两个结点不在并查集或者不连通，那么就把他们合并到一组
# 如果两个结点联通，那么就求出距离并检查是否冲突
# !注意精度问题

from collections import defaultdict
from typing import List
from UnionFindWithDist import UnionFindMapWithDist2

EPS = 1e-6


def id(o: object) -> int:
    if o not in _pool:
        _pool[o] = len(_pool)
    return _pool[o]


_pool = dict()


class Solution:
    def checkContradictions(self, equations: List[List[str]], values: List[float]) -> bool:
        """检查方程是否存在矛盾"""
        uf = UnionFindMapWithDist2()
        for (key1, key2), value in zip(equations, values):
            id1, id2 = id(key1), id(key2)
            if not uf.isConnected(id1, id2):
                uf.union(id1, id2, value)
            else:
                dist = uf.dist(id1, id2)
                if abs(dist - value) > EPS:
                    return True

        return False

    def checkContradictions2(self, equations: List[List[str]], values: List[float]) -> bool:
        """dfs或bfs求出每个点到每个组的根的距离,再逐一检验"""
        adjMap = defaultdict(list)
        allVertex = set()
        for (key1, key2), value in zip(equations, values):
            adjMap[key1].append((key2, value))
            adjMap[key2].append((key1, 1 / value))
            allVertex.add(key1)
            allVertex.add(key2)

        def dfs(cur: int, curValue: float) -> None:
            if cur in visited:
                return
            visited.add(cur)
            dist[cur] = curValue
            for next, value in adjMap[cur]:
                dfs(next, curValue * value)

        visited = set()
        dist = defaultdict(lambda: 1.0)
        for cur in allVertex:
            if cur not in visited:
                dfs(cur, 1.0)

        for (key1, key2), value in zip(equations, values):
            if abs((dist[key2] / dist[key1]) - value) > EPS:
                return True  # !有矛盾
        return False
