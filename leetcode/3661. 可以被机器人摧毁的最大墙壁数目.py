# 3661. 可以被机器人摧毁的最大墙壁数目
#
# 一条无限长的直线上分布着一些机器人和墙壁。给你整数数组 robots ，distance 和 walls：
# robots[i] 是第 i 个机器人的位置。
# distance[i] 是第 i 个机器人的子弹可以行进的 最大 距离。
# walls[j] 是第 j 堵墙的位置。
# 每个机器人有 一颗 子弹，可以向左或向右发射，最远距离为 distance[i] 米。
#
# 子弹会摧毁其射程内路径上的每一堵墙。机器人是固定的障碍物：如果子弹在到达墙壁前击中另一个机器人，它会 立即 在该机器人处停止，无法继续前进。
#
# 返回机器人可以摧毁墙壁的 最大 数量。
#
# 注意：
#
# 墙壁和机器人可能在同一位置；该位置的墙壁可以被该位置的机器人摧毁。
# 机器人不会被子弹摧毁。
#
# dp(i,j) 看到第i个机器人，当前往哪边射击
# 加哨兵方便边界处理

from bisect import bisect_left
from functools import lru_cache
from typing import List


INF = int(1e18)


class Solution:
    def maxWalls(self, robots: List[int], distance: List[int], walls: List[int]) -> int:
        n = len(robots)
        data = [(0, 0)] + sorted(zip(robots, distance)) + [(INF, 0)]
        walls.sort()

        def countWall(start: int, end: int):
            return bisect_left(walls, end) - bisect_left(walls, start)

        @lru_cache(None)
        def dfs(pos: int, j: int) -> int:
            if pos == 0:
                return 0
            x1, d1 = data[pos]
            leftX = max(x1 - d1, data[pos - 1][0] + 1)
            resLeft = dfs(pos - 1, 0) + countWall(leftX, x1 + 1)
            x2, d2 = data[pos + 1]
            if j == 0:
                x2 -= d2
            rightX = min(x1 + d1, x2 - 1)
            resRight = dfs(pos - 1, 1) + countWall(x1, rightX + 1)
            return max(resLeft, resRight)

        res = dfs(n, 1)  # 右边没有机器人，相当于哨兵往右射击
        dfs.cache_clear()
        return res
