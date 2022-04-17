from functools import lru_cache
from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= n <= 105
# parent[0] == -1


# 树的特点
class Solution:
    def longestPath(self, parent: List[int], s: str) -> int:
        """请你找出路径上任意一对相邻节点都没有分配到相同字符的 最长路径 ，并返回该路径的长度。"""

        @lru_cache(None)
        def dfs(cur: int, pre: int) -> int:
            res = 1
            for next in adjMap[cur]:
                if next == pre:
                    continue
                if s[next] == s[cur]:
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

        res = 1
        for i in range(n):
            res = max(res, dfs(i, -1))
        dfs.cache_clear()
        return res


print(Solution().longestPath(parent=[-1, 0, 0, 1, 1, 2], s="abacbe"))
print(Solution().longestPath(parent=[-1, 0, 0, 0], s="aabc"))
