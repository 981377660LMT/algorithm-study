from typing import List

# 炸弹人游戏

# 'W' 表示一堵墙
# 'E' 表示一个敌人
# '0'（数字 0）表示一个空位

# 请你计算放一个炸弹最多能炸多少敌人。
# 好像999. 可以被一步捕获的棋子数
# 暴力:O(n^3) 可以看到暴力时有很多重复的计算 优化:相邻两个格子可以dpO(n^2)
class Solution:
    def maxKilledEnemies(self, grid: List[List[str]]) -> int:
        if not any(grid):
            return 0

        row = len(grid)
        col = len(grid[0])
        dp = [[0] * col for _ in range(row)]

        # 炸右边
        for r in range(row):
            cur = 0
            for c in range(col):
                if grid[r][c] == 'W':
                    cur = 0
                elif grid[r][c] == 'E':
                    cur += 1
                dp[r][c] += cur

        # 炸左边
        for r in range(row):
            cur = 0
            for c in range(col - 1, -1, -1):
                if grid[r][c] == 'W':
                    cur = 0
                elif grid[r][c] == 'E':
                    cur += 1
                dp[r][c] += cur

        # 炸下边
        for c in range(col):
            cur = 0
            for r in range(row):
                if grid[r][c] == 'W':
                    cur = 0
                elif grid[r][c] == 'E':
                    cur += 1
                dp[r][c] += cur

        # 炸上边
        for c in range(col):
            cur = 0
            for r in range(row - 1, -1, -1):
                if grid[r][c] == 'W':
                    cur = 0
                elif grid[r][c] == 'E':
                    cur += 1
                dp[r][c] += cur

        res = 0
        for r in range(row):
            for c in range(col):
                if grid[r][c] == "0":
                    res = max(res, dp[r][c])

        return res


print(
    Solution().maxKilledEnemies([["0", "E", "0", "0"], ["E", "0", "W", "E"], ["0", "E", "0", "0"]])
)

# 解释: 对于如下网格

# 0 E 0 0
# E 0 W E
# 0 E 0 0

# 假如在位置 (1,1) 放置炸弹的话，可以炸到 3 个敌人
