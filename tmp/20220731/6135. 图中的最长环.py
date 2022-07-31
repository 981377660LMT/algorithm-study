from typing import DefaultDict, List, Set, Tuple, Optional
from collections import defaultdict, Counter, deque


MOD = int(1e9 + 7)
INF = int(1e20)


def findCycleAndCalDepth(
    n: int, adjMap: DefaultDict[int, Set[int]], degrees: List[int], *, isDirected: bool
) -> Tuple[List[List[int]], List[int]]:
    depth = [0] * n
    queue = deque([(i, 0) for i in range(n) if degrees[i] == (0 if isDirected else 1)])
    visited = [False] * n
    while queue:
        cur, dist = queue.popleft()
        visited[cur] = True
        for next in adjMap[cur]:
            depth[next] = max(depth[next], dist + 1)
            degrees[next] -= 1
            if degrees[next] == (0 if isDirected else 1):
                queue.append((next, dist + 1))

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
    def longestCycle(self, edges: List[int]) -> int:
        """
        每个节点至多有一条出边
        外向基环树最大环
        """
        n = len(edges)
        adjMap = defaultdict(set)
        deg = [0] * n
        for u, v in enumerate(edges):
            if v == -1:
                continue
            adjMap[u].add(v)
            deg[v] += 1
        cycle, _ = findCycleAndCalDepth(n, adjMap, deg, isDirected=True)
        if not cycle:
            return -1
        return max(len(g) for g in cycle)


print(Solution().longestCycle(edges=[3, 3, 4, 2, 3]))
