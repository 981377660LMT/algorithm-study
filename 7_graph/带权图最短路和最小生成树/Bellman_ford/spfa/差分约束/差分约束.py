"""如果所有边权非负,可以把spfa换成dijkstra"""

from collections import deque
from heapq import heappop, heappush
from typing import List, Tuple

INF = int(1e18)


class DualShortestPath:
    """差分约束求不等式组每个变量的`最优解`"""

    __slots__ = ("_n", "_g", "_min", "_hasNeg")

    def __init__(self, n: int, min: bool) -> None:
        """
        Args:
            n: 变量数量
            min: True求最小值(最长路), False求最大值(最短路)
        """
        self._n = n
        self._g = [[] for _ in range(n)]
        self._min = min
        self._hasNeg = False

    def addEdge(self, i: int, j: int, w: int) -> None:
        """
        添加约束:
        如果 min=False (求最大值/上界): 意味着 val[i] <= val[j] + w
        如果 min=True  (求最小值/下界): 意味着 val[i] >= val[j] + w (通常逻辑)
        """
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


class Solution:
    # 100613. 找到带限制序列的最大值
    # https://leetcode.cn/problems/find-maximum-value-in-a-constrained-sequence/description/
    def findMaxVal(self, n: int, restrictions: List[List[int]], diff: List[int]) -> int:
        dsp = DualShortestPath(n, False)
        for i, d in enumerate(diff):
            dsp.addEdge(i + 1, i, d)
            dsp.addEdge(i, i + 1, d)
        for i, v in restrictions:
            dsp.addEdge(i, 0, v)
        dist, _ = dsp.run()
        return max(dist)

    def findMaxVal2(self, n: int, restrictions: List[List[int]], diff: List[int]) -> int:
        res = [0] + [INF] * (n - 1)
        for i, v in restrictions:
            res[i] = min(res[i], v)
        for i in range(n - 1):
            res[i + 1] = min(res[i + 1], res[i] + diff[i])
        for i in range(n - 2, -1, -1):
            res[i] = min(res[i], res[i + 1] + diff[i])
        return max(res)
