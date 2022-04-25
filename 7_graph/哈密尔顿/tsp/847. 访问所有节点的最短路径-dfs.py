# from functools import lru_cache
# from typing import List
# from collections import deque

# 这道题dfs是做不出来的


# TSP problem
# 1 <= n <= 12
# 返回能够访问所有节点的最短路径的长度。你可以在任一节点开始和停止，也可以多次重访节点，并且可以重用边。
# class Solution:
#     def shortestPathLength(self, graph: List[List[int]]) -> int:
#         @lru_cache(None)
#         def dfs(cur: int, pre: int, visited: int) -> int:
#             if visited == target:
#                 return 0
#             res = int(1e20)
#             for next in graph[cur]:
#                 if next == pre:
#                     continue
#                 res = min(res, dfs(next, cur, visited | (1 << next)) + 1)
#             return res

#         n = len(graph)
#         target = (1 << n) - 1

#         res = int(1e20)
#         for i in range(n):
#             res = min(res, dfs(i, -1, 1 << i))
#         dfs.cache_clear()
#         return res


# print(Solution().shortestPathLength(graph=[[1, 2, 3], [0], [0], [0]]))
# 输出：4
# 解释：一种可能的路径为 [1,0,2,0,3]
