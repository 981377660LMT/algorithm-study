from collections import deque
from typing import List


def bfsDepth(adjList: List[List[int]], start: int, dist: int) -> List[int]:
    """返回距离start为dist的结点"""
    if dist < 0:
        return []
    if dist == 0:
        return [start]
    queue = deque([start])
    visited = set([start])
    todo = dist
    while queue and todo > 0:
        len_ = len(queue)
        for _ in range(len_):
            cur = queue.popleft()
            for next in adjList[cur]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)
        todo -= 1
    return list(queue)
