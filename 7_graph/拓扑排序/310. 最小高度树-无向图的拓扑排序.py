# 请你找到所有的 最小高度树 并按 任意顺序 返回它们的根节点标签列表。
# 思路：不断删除叶子节点，无向图的拓扑排序

# 310. 最小高度树-无向图的拓扑排序
from collections import deque
from typing import List


class Solution:
    def findMinHeightTrees(self, n: int, edges: List[List[int]]) -> List[int]:
        if n == 1:
            return [0]

        adjList = [[] for _ in range(n)]
        deg = [0] * n
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
            deg[u] += 1
            deg[v] += 1

        queue = deque([i for i in range(n) if deg[i] == 1])
        while n > 2:
            len_ = len(queue)
            n -= len_
            for _ in range(len_):
                cur = queue.popleft()
                for next in adjList[cur]:
                    deg[next] -= 1
                    if deg[next] == 1:
                        queue.append(next)

        return list(queue)
