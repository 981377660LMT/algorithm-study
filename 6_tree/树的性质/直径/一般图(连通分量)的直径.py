# 求连通分量的直径/一般图的直径
# !要暴力枚举每个起点做 bfs 看最大层数
# 时间复杂度`O(V*(V+E))`

from collections import deque
from typing import List


def calDiameter(n: int, adjList: List[List[int]], group: List[int]) -> int:
    """bfs求连通分量 `group` 的直径长度"""
    res = 0
    for start in group:
        visited, queue = set([start]), deque([start])
        diameter = -1
        while queue:
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                for next in adjList[cur]:
                    if next in visited:
                        continue
                    visited.add(next)
                    queue.append(next)
            diameter += 1
        res = max(res, diameter)
    return res


assert calDiameter(5, [[1], [0, 2], [1, 3], [2, 4], [4]], [0, 1, 2, 3, 4]) == 4
