# 3552. 网格传送门旅游
# 给你一个大小为 m x n 的二维字符网格 matrix，用字符串数组表示，其中 matrix[i][j] 表示第 i 行和第 j 列处的单元格。每个单元格可以是以下几种字符之一：
#
# '.' 表示一个空单元格。
# '#' 表示一个障碍物。
# 一个大写字母（'A' 到 'Z'）表示一个传送门。
# 你从左上角单元格 (0, 0) 出发，目标是到达右下角单元格 (m - 1, n - 1)。你可以从当前位置移动到相邻的单元格（上、下、左、右），移动后的单元格必须在网格边界内且不是障碍物。
#
# 如果你踏入一个包含传送门字母的单元格，并且你之前没有使用过该传送门字母，你可以立即传送到网格中另一个具有相同字母的单元格。这次传送不计入移动次数，但每个字母对应的传送门在旅程中 最多 只能使用一次。
#
# 返回到达右下角单元格所需的 最少 移动次数。如果无法到达目的地，则返回 -1。
#
# 本质上是计算如下图的最短路：
#
# 所有相同字母之间都有一条边权为 0 的边。
# 所有相邻格子之间都有一条边权为 1 的边。
#
# !一旦我们处理过字符 'A' 的传送，以后再遇到 'A' 就不需要再处理传送了。因为第一次到达 'A' 时，我们已经把所有 'A' 的格子以最优步数加入队列了。后续再通过走路到达某个 'A' 格子，步数肯定不会比第一次更优。

from collections import defaultdict, deque
from typing import List

DIRS = [(0, 1), (0, -1), (1, 0), (-1, 0)]
INF = int(1e18)


class Solution:
    def minMoves(self, matrix: List[str]) -> int:
        if matrix[-1][-1] == "#":
            return -1

        ROW, COL = len(matrix), len(matrix[0])
        portals = defaultdict(list)
        for i, row in enumerate(matrix):
            for j, c in enumerate(row):
                if c.isupper():
                    portals[c].append((i, j))

        dist = [[INF] * COL for _ in range(ROW)]
        dist[0][0] = 0
        queue = deque([(0, 0)])

        while queue:
            x, y = queue.popleft()
            d = dist[x][y]
            if x == ROW - 1 and y == COL - 1:
                return d

            c = matrix[x][y]
            if c in portals:
                for nx, ny in portals[c]:
                    if dist[nx][ny] > d:
                        dist[nx][ny] = d
                        queue.appendleft((nx, ny))
                del portals[c]  # 清空, 只能用一次

            for dx, dy in DIRS:
                nx, ny = x + dx, y + dy
                if 0 <= nx < ROW and 0 <= ny < COL and matrix[nx][ny] != "#":
                    if dist[nx][ny] > d + 1:
                        dist[nx][ny] = d + 1
                        queue.append((nx, ny))

        return -1
