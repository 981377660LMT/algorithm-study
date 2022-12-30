# 01矩阵中是否存在一条路径使得0和1的个数相等
# m,n<=100


# !位运算加速dp
from typing import List


class Solution:
    def isThereAPath(self, grid: List[List[int]]) -> bool:
        ROW, COL = len(grid), len(grid[0])
        if not (ROW + COL) & 1:
            return False

        # !dp[i][j]记录每个位置可以取到的前缀和
        dp = [[0] * COL for _ in range(ROW)]
        dp[0][0] = 1 << (ROW + COL)
        for r in range(ROW):
            for c in range(COL):
                if r:
                    dp[r][c] |= dp[r - 1][c]
                if c:
                    dp[r][c] |= dp[r][c - 1]
                if grid[r][c]:
                    dp[r][c] <<= 1
                else:
                    dp[r][c] >>= 1

        return not not dp[-1][-1] & (1 << (ROW + COL))


assert Solution().isThereAPath(grid=[[0, 1, 0, 0], [0, 1, 0, 0], [1, 0, 1, 0]])
