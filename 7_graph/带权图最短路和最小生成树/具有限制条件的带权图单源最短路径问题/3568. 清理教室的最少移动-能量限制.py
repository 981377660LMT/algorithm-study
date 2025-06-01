# 864. 获取所有钥匙的最短路径 带上了能量限制
# https://leetcode.cn/problems/shortest-path-to-get-all-keys/
# 给你一个 m x n 的网格图 classroom，其中一个学生志愿者负责清理散布在教室里的垃圾。网格图中的每个单元格是以下字符之一：
#
# 'S' ：学生的起始位置
# 'L' ：必须收集的垃圾（收集后，该单元格变为空白）
# 'R' ：重置区域，可以将学生的能量恢复到最大值，无论学生当前的能量是多少（可以多次使用）
# 'X' ：学生无法通过的障碍物
# '.' ：空白空间
# 同时给你一个整数 energy，表示学生的最大能量容量。学生从起始位置 'S' 开始，带着 energy 的能量出发。
#
# 每次移动到相邻的单元格（上、下、左或右）会消耗 1 单位能量。如果能量为 0，学生此时只有处在 'R' 格子时可以继续移动，此区域会将能量恢复到 最大 能量值 energy。
#
# 返回收集所有垃圾所需的 最少 移动次数，如果无法完成，返回 -1。

from typing import List

DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]


class Solution:
    def minMoves(self, classroom: List[str], energy: int) -> int:
        m, n = len(classroom), len(classroom[0])
        idx = [[0] * n for _ in range(m)]
        sx, sy = 0, 0
        numL = 0
        for i, row in enumerate(classroom):
            for j, b in enumerate(row):
                if b == "S":
                    sx, sy = i, j
                elif b == "L":
                    idx[i][j] = 1 << numL
                    numL += 1
        if numL == 0:
            return 0

        bestEnergy = [[[-1] * (1 << numL) for _ in range(n)] for _ in range(m)]
        fullMask = (1 << numL) - 1
        bestEnergy[sx][sy][0] = energy
        queue = [(sx, sy, energy, 0)]

        res = 0
        while queue:
            nextQueue = []
            for x, y, e, visited in queue:
                if visited == fullMask:
                    return res
                if e == 0:
                    continue
                for dx, dy in DIR4:
                    nx, ny = x + dx, y + dy
                    if 0 <= nx < m and 0 <= ny < n:
                        cell = classroom[nx][ny]
                        if cell == "X":
                            continue
                        nextVisited = visited | idx[nx][ny]
                        if cell == "R":
                            nextEnergy = energy
                        else:
                            nextEnergy = e - 1
                        if bestEnergy[nx][ny][nextVisited] >= nextEnergy:
                            continue
                        bestEnergy[nx][ny][nextVisited] = nextEnergy
                        nextQueue.append((nx, ny, nextEnergy, nextVisited))
            queue = nextQueue
            res += 1

        return -1
