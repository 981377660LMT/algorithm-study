from typing import List

# 每一行 中选取一个格子，选中坐标为 (r, c) 的格子会给你的总得分 增加 points[r][c] 。
# 对于相邻行 r 和 r + 1 （其中 0 <= r < m - 1），选中坐标为 (r, c1) 和 (r + 1, c2) 的格子，
# 你的总得分 减少 abs(c1 - c2) 。
# 请你返回你能得到的 最大 得分。
INF = int(1e20)


class Solution:
    def maxPoints(self, points: List[List[int]]) -> int:
        """
        对于所有上一行左边的数，都是加上它的坐标减去自己的坐标；
        对于上一行所有右边的数，都是减去它的坐标加上自己的坐标；

        绝对值=>拆绝对值+讨论
        """
        row, col = len(points), len(points[0])
        dp = points[0][:]

        # 拆去绝对值讨论  左边的点 val[j]+j+val[i]-i 右边的点 val[j]-j+val[i]+i
        for i in range(1, row):
            ndp = [-INF] * col

            # preMax = -INF
            # for i in range(col):
            #     preMax = max(preMax, dp[i] + i)
            #     ndp[i] = max(ndp[i], preMax + points[r][i] - i)

            # preMax = -INF
            # for i in range(col - 1, -1, -1):
            #     preMax = max(preMax, dp[i] - i)
            #     ndp[i] = max(ndp[i], preMax + points[r][i] + i)

            lMax, rMax = -INF, -INF
            for j in range(col):  # 左右一起合并写 i 全部换成 n-1-i
                rj = col - 1 - j
                lMax = max(lMax, dp[j] + j)
                rMax = max(rMax, dp[rj] - rj)
                ndp[j] = max(ndp[j], lMax + points[i][j] - j)
                ndp[rj] = max(ndp[rj], rMax + points[i][rj] + rj)
            dp = ndp

        return max(dp)
