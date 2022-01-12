from typing import List
from functools import lru_cache
from math import comb

# 桌面上有 2n 个颜色不完全相同的球，球上的颜色共有 k 种。
# 其中 balls[i] 是颜色为 i 的球的数量。
class Solution:
    def getProbability(self, balls: List[int]) -> float:
        n = sum(balls) // 2

        @lru_cache(None)
        def dfs(index, ball1, ball2, color1, color2) -> int:
            """[summary]

            Args:
                index ([type]): [description]
                ball1 ([type]): [description]
                ball2 ([type]): [description]
                color1 ([type]): [颜色种类数]
                color2 ([type]): [颜色种类数]

            Returns:
                int: [description]
            """
            if ball1 > n or ball2 > n:
                return 0
            if index == len(balls):
                return int(color1 == color2)

            res = 0
            for b in range(balls[index] + 1):
                res += dfs(
                    index + 1,
                    ball1 + b,
                    ball2 + balls[index] - b,
                    color1 + (b > 0),
                    color2 + (b < balls[index]),
                ) * comb(balls[index], b)

            return res

        return dfs(0, 0, 0, 0, 0) / comb(2 * n, n)


print(Solution().getProbability([2, 1, 1]))
# 输出：0.66667
# 解释：球的列表为 [1, 1, 2, 3]
# 随机打乱，得到 12 种等概率的不同打乱方案，每种方案概率为 1/12 ：
# [1,1 / 2,3], [1,1 / 3,2], [1,2 / 1,3], [1,2 / 3,1], [1,3 / 1,2], [1,3 / 2,1], [2,1 / 1,3], [2,1 / 3,1], [2,3 / 1,1], [3,1 / 1,2], [3,1 / 2,1], [3,2 / 1,1]
# 然后，我们将前两个球放入第一个盒子，后两个球放入第二个盒子。
# 这 12 种可能的随机打乱方式中的 8 种满足「两个盒子中球的颜色数相同」。
# 概率 = 8/12 = 0.66667
