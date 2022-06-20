# 如果distance[i][j]<=2米说明它们属于同一`组`
# 请你设计一个算法判断该场地是否可以保证每组摄像头至少有一个负责人管理。
from typing import List
from UnionFind import UnionFindArray


class Solution:
    def isCompliance(self, distance: List[List[int]], k: int) -> bool:
        n = len(distance)
        uf = UnionFindArray(n)
        for i in range(n):
            for j in range(i + 1, n):
                if distance[i][j] <= 2:
                    uf.union(i, j)
        return uf.count <= k

