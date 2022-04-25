# 给出一张图（有节点和边），保证它是一棵树。
# 每个节点表示该位置的海拔高度。假设从每个节点位置发洪水的话，
# 水会顺着海拔往严格更低的相邻节点流动。
# 问在那个节点位置发洪水会覆盖最多的区域。
from collections import defaultdict
from functools import lru_cache
from typing import List


class Solution:
    def solve(self, parent: List[int], values: List[int]) -> int:
        """请你找出路径上任意一对相邻节点都没有分配到相同字符的 最长路径 ，并返回该路径的长度。"""

        @lru_cache(None)
        def dfs(cur: int, pre: int) -> int:
            res = 1
            for next in adjMap[cur]:
                if next == pre:
                    continue
                if values[next] >= values[cur]:
                    continue
                res = max(res, dfs(next, cur) + 1)
            return res

        n = len(parent)
        adjMap = defaultdict(set)
        for i in range(n):
            pre, cur = parent[i], i
            if pre == -1:
                continue
            adjMap[pre].add(cur)
            adjMap[cur].add(pre)

        res, max_ = -1, -1
        for i in range(n):
            cur = dfs(i, -1)
            if cur > max_:
                max_ = cur
                res = i

        dfs.cache_clear()
        return res
