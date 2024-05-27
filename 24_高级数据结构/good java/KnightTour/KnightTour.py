# 骑士路线/马的路线
# 找到一个骑士的移动顺序，使得 board 中每个单元格都 恰好 被访问一次（起始单元格已被访问，不应 再次访问）
# 输入的数据保证在给定条件下至少存在一种访问所有单元格的移动顺序。

import random
from typing import List


DIR8 = [(1, 2), (-1, 2), (1, -2), (-1, -2), (2, 1), (-2, 1), (-2, -1), (2, -1)]
INF = int(1e18)


class KnightTour:
    __slots__ = "order", "_ids", "_n", "_m", "_step"

    def __init__(self, n: int, m: int, r: int, c: int):
        self.order = [[-1] * m for _ in range(n)]
        self._ids = list(range(n * m))
        self._n, self._m = n, m
        while True:
            for row in self.order:
                for j in range(m):
                    row[j] = -1
            random.shuffle(self._ids)
            self._step = 0
            if self._dfs(r, c):
                break

    def _possible(self, i: int, j: int) -> bool:
        return 0 <= i < self._n and 0 <= j < self._m and self.order[i][j] == -1

    def _degree(self, i: int, j: int) -> int:
        return sum(self._possible(i + di, j + dj) for di, dj in DIR8)

    def _dfs(self, i: int, j: int) -> bool:
        self.order[i][j] = self._step
        self._step += 1
        bestX, bestY = -1, -1
        bestDeg = INF
        bestId = -1
        for di, dj in DIR8:
            x, y = i + di, j + dj
            if not self._possible(x, y):
                continue
            d = self._degree(x, y)
            if d < bestDeg:
                bestDeg = d
                bestId = -1
            if d == bestDeg and bestId < self._ids[x * self._m + y]:
                bestId = self._ids[x * self._m + y]
                bestX, bestY = x, y
        if bestDeg == INF:
            return self._step == self._n * self._m
        return self._dfs(bestX, bestY)


if __name__ == "__main__":
    # 2664. 巡逻的骑士
    # https://leetcode.cn/problems/the-knights-tour/description/
    class Solution:
        def tourOfKnight(self, m: int, n: int, r: int, c: int) -> List[List[int]]:
            K = KnightTour(m, n, r, c)
            return K.order
