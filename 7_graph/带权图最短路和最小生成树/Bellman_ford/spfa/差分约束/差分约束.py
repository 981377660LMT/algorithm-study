from collections import deque
from typing import List, Tuple

INF = int(1e18)


class DualShortestPath:
    """差分约束求不等式组每个变量的`最优解`"""

    __slots__ = ("_n", "_g", "_min")

    def __init__(self, n: int, min: bool) -> None:
        self._n = n
        self._g = [[] for _ in range(n)]
        self._min = min

    def addEdge(self, i: int, j: int, w: int) -> None:
        """f(j) <= f(i) + w"""
        if self._min:
            self._g[i].append((j, w))
        else:
            self._g[j].append((i, w))

    def run(self) -> Tuple[List[int], bool]:
        """求 `f(i) - f(0)` 的最小值/最大值, 并检测是否有负环/正环"""
        return self._spfaMin() if self._min else self._spfaMax()

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
                        if count[next] >= n:
                            return [], False
                        inQueue[next] = True
                        queue.appendleft(next)  # !栈优化

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
                        if count[next] >= n:
                            return [], False
                        inQueue[next] = True
                        queue.appendleft(next)  # !栈优化
        return dist, True
