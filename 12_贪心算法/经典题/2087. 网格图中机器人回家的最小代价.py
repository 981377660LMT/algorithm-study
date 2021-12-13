from typing import List

# 1 <= m, n <= 105
# 不能dp因为dp是O(n^2)


# 如果机器人往 上 或者往 下 移动到第 r 行 的格子，那么代价为 rowCosts[r] 。
# 如果机器人往 左 或者往 右 移动到第 c 列 的格子，那么代价为 colCosts[c] 。


class Solution:
    def minCost(
        self, startPos: List[int], homePos: List[int], rowCosts: List[int], colCosts: List[int]
    ) -> int:
        res = 0
        sr, sc = startPos[0], startPos[1]
        er, ec = homePos[0], homePos[1]

        # ----sr的数值小
        for r in range(sr + 1, er + 1):
            res += rowCosts[r]
        # ----sr的数值大
        for r in range(sr - 1, er - 1, -1):
            res += rowCosts[r]
        # ----sc的数值小
        for c in range(sc + 1, ec + 1):
            res += colCosts[c]
        # ----sc的数值大
        for c in range(sc - 1, ec - 1, -1):
            res += colCosts[c]

        return res


print(
    Solution().minCost(startPos=[1, 0], homePos=[2, 3], rowCosts=[5, 4, 3], colCosts=[8, 2, 6, 7])
)
