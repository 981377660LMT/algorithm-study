from functools import lru_cache
from typing import DefaultDict, List, Set, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList


from collections import defaultdict
from typing import DefaultDict, List, Callable


class UnionFindArray:
    __slots__ = ("n", "part", "_parent", "_rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self._rank = [1] * n

    def find(self, x: int) -> int:
        while x != self._parent[x]:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        return True

    def unionWithCallback(self, x: int, y: int, f: Callable[[int, int], None]) -> bool:
        """
        f: 合并后的回调函数, 入参为 (big, small)
        """
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        f(rootY, rootX)
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self._parent))

    def getSize(self, x: int) -> int:
        return self._rank[self.find(x)]

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


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
    def countVisitedNodes(self, edges: List[int]) -> List[int]:
        n = len(edges)
        adjMap = defaultdict(set)
        degrees = [0] * n
        for i in range(n):
            u, v = i, edges[i]
            adjMap[u].add(v)
            degrees[v] += 1
        cycles, _ = findCycleAndCalDepth(n, adjMap, degrees, isDirected=True)
        uf = UnionFindArray(n)
        isInCycle = [False] * n
        for cycle in cycles:
            for node in cycle:
                isInCycle[node] = True
                uf.union(node, cycle[0])
        groups = uf.getGroups()

        @lru_cache(None)
        def dfs(cur: int) -> int:
            if isInCycle[cur]:
                return len(groups[uf.find(cur)])
            return 1 + dfs(edges[cur])

        return [dfs(i) for i in range(n)]
