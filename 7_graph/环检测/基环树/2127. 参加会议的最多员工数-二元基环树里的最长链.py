# 两种情况:1.所有的二元基环树里的最长链之和;2.唯一的最长环的长度

from collections import defaultdict, deque
from typing import DefaultDict, List, Set, Tuple


def findCycleAndCalDepth(
    n: int, adjMap: DefaultDict[int, Set[int]], degrees: List[int], *, isDirected: bool
) -> Tuple[List[List[int]], List[int]]:
    """无/有向基环树找环上的点,并记录每个点在拓扑排序中的最大深度,最外层的点深度为0"""

    def max(a: int, b: int) -> int:
        return a if a > b else b

    depth = [0] * n
    startDeg = 0 if isDirected else 1
    queue = deque([i for i in range(n) if degrees[i] == startDeg])
    visited = [False] * n
    while queue:
        cur = queue.popleft()
        visited[cur] = True
        for next in adjMap[cur]:
            depth[next] = max(depth[next], depth[cur] + 1)
            degrees[next] -= 1
            if degrees[next] == startDeg:
                queue.append(next)

    def dfs(cur: int, path: List[int]) -> None:
        if visited[cur]:
            return
        visited[cur] = True
        path.append(cur)
        for next in adjMap[cur]:
            dfs(next, path)

    cycleGroup = []
    for i in range(n):
        if visited[i]:
            continue
        path = []
        dfs(i, path)
        cycleGroup.append(path)

    return cycleGroup, depth


class Solution:
    def maximumInvitations(self, favorite: List[int]) -> int:
        n = len(favorite)
        adjMap, degrees = defaultdict(set), [0] * n
        for u, v in enumerate(favorite):
            adjMap[u].add(v)
            adjMap[v].add(u)
            degrees[u] += 1
            degrees[v] += 1

        cycleGroup, depth = findCycleAndCalDepth(n, adjMap, degrees, isDirected=False)
        # 两种情况:1.所有的二元基环树里的最长链之和;2.唯一的最长环的长度
        cand1 = sum((1 + depth[i]) for i in range(n) if favorite[favorite[i]] == i)
        cand2 = max(len(cycle) for cycle in cycleGroup)
        return max(cand1, cand2)


print(Solution().maximumInvitations([7, 12, 17, 9, 0, 7, 14, 5, 3, 15, 6, 14, 10, 14, 10, 1, 1, 4]))
