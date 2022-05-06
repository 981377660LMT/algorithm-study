# Detecting an Odd Length Cycle
# return whether the graph has an odd length cycle.

# 无向图中是否存在奇数长度环:等价于不是二分图
from collections import defaultdict
from typing import List


def isBipartite(adjList: List[List[int]]) -> bool:
    def dfs(cur: int, color: int) -> bool:
        colors[cur] = color
        for next in adjList[cur]:
            if colors[next] == -1:
                if not dfs(next, color ^ 1):  # type: ignore
                    return False
            elif colors[next] == color:
                return False
        return True

    colors = defaultdict(lambda: -1)
    n = len(adjList)
    for i in range(n):
        if colors[i] == -1:
            if not dfs(i, 0):
                return False
    return True


class Solution:
    def solve(self, graph):
        return not isBipartite(graph)


print(Solution().solve(graph=[[1, 2], [0, 2], [0, 1]]))
print(Solution().solve(graph=[[2, 3], [2, 3], [0, 1, 3], [0, 1, 2]]))

