# 请你找到所有的 最小高度树 并按 任意顺序 返回它们的根节点标签列表。
# 思路：不断删除叶子节点，无向图的拓扑排序

# 310. 最小高度树-无向图的拓扑排序
from collections import defaultdict, deque
from typing import List


class Solution:
    def findMinHeightTrees(self, n: int, edges: List[List[int]]) -> List[int]:
        if n == 1:
            return [0]

        degree, adjMap = [0] * n, defaultdict(set)
        for u, v in edges:
            degree[u] += 1
            degree[v] += 1
            adjMap[u].add(v)
            adjMap[v].add(u)

        queue = deque([i for i in range(n) if degree[i] == 1])

        while n >= 3:
            len_ = len(queue)
            n -= len_
            for _ in range(len_):
                cur = queue.popleft()
                for next in adjMap[cur]:
                    degree[next] -= 1
                    if degree[next] == 1:
                        queue.append(next)

        return list(queue)
