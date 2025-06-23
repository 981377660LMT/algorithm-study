from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minIncrease(self, n: int, edges: List[List[int]], cost: List[int]) -> int:
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        res = 0

        def dfs(cur: int, pre: int) -> int:
            nonlocal res
            if len(adjList[cur]) == 1 and adjList[cur][0] == pre:
                return cost[cur]
            childDp = [dfs(next_, cur) for next_ in adjList[cur] if next_ != pre]
            max_ = max(childDp)
            for v in childDp:
                if v < max_:
                    res += 1
            return cost[cur] + max_

        dfs(0, -1)
        return res


# 输入： n = 3, edges = [[0,1],[0,2]], cost = [2,1,3]©leetcode
print(Solution().minIncrease(3, [[0, 1], [0, 2]], [2, 1, 3]))  # 输出：2
#  n = 3, edges = [[0,1],[1,2]], cost = [5,1,4]©leetcode
print(Solution().minIncrease(3, [[0, 1], [1, 2]], [5, 1, 4]))  # 输出：3
# n = 5, edges = [[0,4],[0,1],[1,2],[1,3]], cost = [3,4,1,1,7]©leetcode
print(Solution().minIncrease(5, [[0, 4], [0, 1], [1, 2], [1, 3]], [3, 4, 1, 1, 7]))  # 输出：3
