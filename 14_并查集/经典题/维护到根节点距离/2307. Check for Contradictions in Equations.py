# 检查方程是否存在矛盾
# 如果两个结点不在并查集或者不连通，那么就把他们合并到一组
# 如果两个结点联通，那么就求出距离并检查是否冲突

from typing import List
from UnionFindMapWithDist import UnionFindMapWithDist

EPS = 1e-5


class Solution:
    def checkContradictions(
        self, equations: List[List[str]], values: List[float]
    ) -> bool:
        """检查方程是否存在矛盾"""
        uf = UnionFindMapWithDist[str]()
        for (key1, key2), value in zip(equations, values):
            if key1 not in uf or key2 not in uf:
                uf.add(key1).add(key2)
            if not uf.isConnected(key1, key2):
                uf.union(key2, key1, value)
            else:
                dist = uf.distToRoot[key2] / uf.distToRoot[key1]
                if abs(dist - value) > EPS:
                    return True

        return False
