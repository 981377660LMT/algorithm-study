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

    def lessThanOrEqualTo(self, i: int, j: int, w: int) -> None:
        """f(i) - f(j) <= w"""
        self.addEdge(i, j, w)

    def greaterThanOrEqualTo(self, i: int, j: int, w: int) -> None:
        """f(i) - f(j) >= w"""
        self.lessThanOrEqualTo(j, i, -w)

    def equalTo(self, i: int, j: int, w: int) -> None:
        """f(i) - f(j) == w"""
        self.greaterThanOrEqualTo(i, j, w)
        self.lessThanOrEqualTo(i, j, w)

    def lessThan(self, i: int, j: int, w: int) -> None:
        """f(i) - f(j) < w"""
        self.lessThanOrEqualTo(i, j, w - 1)

    def greaterThan(self, i: int, j: int, w: int) -> None:
        """f(i) - f(j) > w"""
        self.greaterThanOrEqualTo(i, j, w + 1)

    def run(self) -> Tuple[List[int], bool]:
        """求 `f(i) - f(0)` 的最小值/最大值, 并检测是否有负环/正环"""
        if self._min:
            return self._spfaMin()
        if not self._hasNeg:
            return self._dijkMax()
        return self._spfaMax()

    def _spfaMin(self) -> Tuple[List[int], bool]:
        """每个变量的最小值 (带 SLF 优化)"""
        n, g = self._n, self._g
        dist = [0] * n
        queue = deque(list(range(n)))
        count = [1] * n
        inQueue = [True] * n

        while queue:
            cur = queue.popleft()
            inQueue[cur] = False
            for next_, weight in g[cur]:
                cand = dist[cur] + weight
                if cand < dist[next_]:
                    dist[next_] = cand
                    if not inQueue[next_]:
                        count[next_] += 1
                        if count[next_] >= n + 1:
                            return [], False
                        inQueue[next_] = True

                        if queue and cand < dist[queue[0]]:
                            queue.appendleft(next_)
                        else:
                            queue.append(next_)

        return [-num for num in dist], True

    def _spfaMax(self) -> Tuple[List[int], bool]:
        """每个变量的最大值 (带 SLF 优化)"""
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
            for next_, weight in g[cur]:
                cand = dist[cur] + weight
                if cand < dist[next_]:
                    dist[next_] = cand
                    if not inQueue[next_]:
                        count[next_] += 1
                        if count[next_] >= n + 1:
                            return [], False
                        inQueue[next_] = True

                        if queue and cand < dist[queue[0]]:
                            queue.appendleft(next_)
                        else:
                            queue.append(next_)

        return dist, True

    def _dijkMax(self) -> Tuple[List[int], bool]:
        dist = [INF] * self._n
        dist[0] = 0
        pq = [(0, 0)]
        while pq:
            curDist, cur = heappop(pq)
            if curDist > dist[cur]:
                continue
            for next_, weight in self._g[cur]:
                cand = curDist + weight
                if cand < dist[next_]:
                    dist[next_] = cand
                    heappush(pq, (cand, next_))
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


# https://atcoder.jp/contests/awc0005/tasks/awc0005_c
def awc0005():
    N, K = map(int, input().split())
    A = list(map(int, input().split()))

    dsp = DualShortestPath(N + 1, False)
    for i in range(1, N + 1):
        dsp.addEdge(i, 0, -A[i - 1])
    for i in range(1, N):
        dsp.addEdge(i, i + 1, K)
        dsp.addEdge(i + 1, i, K)
    dist, _ = dsp.run()
    res = 0
    for i in range(1, N + 1):
        tmp = -dist[i]
        res += tmp - A[i - 1]

    print(res)


if __name__ == "__main__":
    awc0005()
