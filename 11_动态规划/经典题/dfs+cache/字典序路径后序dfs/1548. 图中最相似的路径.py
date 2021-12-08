from typing import List, Tuple
from collections import defaultdict
from functools import lru_cache

# 给定一个字符串数组 targetPath，你需要找出图中与 targetPath 的 长度相同 且 编辑距离最小 的路径。
# 你需要返回 编辑距离最小的路径中节点的顺序

# 2 <= n <= 100

INF = 0x7FFFFFFF


class Solution:
    def mostSimilar(
        self, n: int, roads: List[List[int]], names: List[str], targetPath: List[str]
    ) -> List[int]:
        target = len(targetPath)
        adjMap = defaultdict(list)
        for a, b in roads:
            adjMap[a].append(b)
            adjMap[b].append(a)

        @lru_cache(None)
        def dfs(curId: int, visitedLen: int) -> Tuple[int, List[int]]:
            if visitedLen == target:
                return 0, []

            cost, path = INF, []
            for next in adjMap[curId]:
                nextCost, nextPath = dfs(next, visitedLen + 1)
                costCand, pathCand = (
                    nextCost + (names[curId] != targetPath[visitedLen]),
                    [curId] + nextPath,
                )
                if costCand < cost:
                    cost = costCand
                    path = pathCand

            return cost, path

        return min((dfs(i, 0) for i in range(n)), key=lambda x: x[0])[1]

    def mostSimilar2(
        self, n: int, roads: List[List[int]], names: List[str], targetPath: List[str]
    ) -> List[int]:
        d = defaultdict(list)
        for a, b in roads:
            d[a].append(b)
            d[b].append(a)

        @lru_cache(None)
        def dp(roadIdx, targetIdx):
            if targetIdx == len(targetPath):
                return 0, []
            mincost, minpath = min(
                [dp(nxt, targetIdx + 1) for nxt in d[roadIdx]], key=lambda x: x[0]
            )
            return mincost + (names[roadIdx] != targetPath[targetIdx]), [roadIdx] + minpath

        return min([dp(i, 0) for i in range(n)], key=lambda x: x[0])[1]


print(
    Solution().mostSimilar(
        n=5,
        roads=[[0, 2], [0, 3], [1, 2], [1, 3], [1, 4], [2, 4]],
        names=["ATL", "PEK", "LAX", "DXB", "HND"],
        targetPath=["ATL", "DXB", "HND", "LAX"],
    )
)

