from collections import deque
from heapq import heappop, heappush
from typing import List, Tuple

INF = int(1e18)


class DualShortestPath:
    """差分约束求不等式组每个变量的`最优解`"""

    __slots__ = ("_n", "_g", "_min", "_hasNeg")

    def __init__(self, n: int, min: bool) -> None:
        self._n = n
        self._g = [[] for _ in range(n)]
        self._min = min
        self._hasNeg = False

    def addEdge(self, i: int, j: int, w: int) -> None:
        """f(j) <= f(i) + w"""
        if self._min:
            self._g[i].append((j, w))
        else:
            self._g[j].append((i, w))
        self._hasNeg |= w < 0

    def run(self) -> Tuple[List[int], bool]:
        """求 `f(i) - f(0)` 的最小值/最大值, 并检测是否有负环/正环"""
        if self._min:
            return self._spfaMin()
        return self._spfaMax() if self._hasNeg else self._dijkMax()

    def _spfaMin(self) -> Tuple[List[int], bool]:
        """每个变量的最小值"""
        n, g = self._n, self._g
        dist = [0] * n
        queue = deque(list(range(n)))
        count = [1] * n
        inQueue = [True] * n

        while queue:
            cur = queue.popleft()
            inQueue[cur] = False
            for next, weight in g[cur]:
                cand = dist[cur] + weight
                if cand < dist[next]:
                    dist[next] = cand
                    if not inQueue[next]:
                        count[next] += 1
                        if count[next] >= n + 1:
                            return [], False
                        inQueue[next] = True
                        queue.appendleft(next)

        return [-num for num in dist], True

    def _spfaMax(self) -> Tuple[List[int], bool]:
        """每个变量的最大值"""
        n, g = self._n, self._g
        dist = [INF] * n
        inQueue = [False] * n
        count = [0] * n

        queue = deque([0])
        dist[0] = 0
        inQueue[0] = True
        count[0] = 1
        while queue:
            cur = queue.popleft()
            inQueue[cur] = False
            for next, weight in g[cur]:
                cand = dist[cur] + weight
                if cand < dist[next]:
                    dist[next] = cand
                    if not inQueue[next]:
                        count[next] += 1
                        if count[next] >= n + 1:
                            return [], False
                        inQueue[next] = True
                        queue.appendleft(next)

        return dist, True

    def _dijkMax(self) -> Tuple[List[int], bool]:
        dist = [INF] * self._n
        dist[0] = 0
        pq = [(0, 0)]
        while pq:
            curDist, cur = heappop(pq)
            if curDist > dist[cur]:
                continue
            for next, weight in self._g[cur]:
                cand = curDist + weight
                if cand < dist[next]:
                    dist[next] = cand
                    heappush(pq, (cand, next))
        return dist, True
