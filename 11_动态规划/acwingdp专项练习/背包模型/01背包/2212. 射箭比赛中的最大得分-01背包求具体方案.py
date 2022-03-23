from typing import List

# 1 <= numArrows <= 1e5
# aliceArrows.length == bobArrows.length == 12
# 背包直接超时


class Solution:
    def maximumBobPoints(self, numArrows: int, aliceArrows: List[int]) -> List[int]:
        goods = []
        for score, cost in enumerate(aliceArrows):
            goods.append((score, cost + 1))

        dp = [0] * (numArrows + 1)
        dp[0] = 0
        assign = [0 for _ in range(numArrows + 1)]

        for i, (score, cost) in enumerate(goods):
            for j in range(numArrows, cost - 1, -1):
                if dp[j - cost] + score > dp[j]:
                    dp[j] = dp[j - cost] + score
                    assign[j] = assign[j - cost] | (1 << i)

        res = [0] * 12
        last = assign[-1]
        for i in range(12):
            if last & (1 << i):
                res[i] += goods[i][1]

        res[0] += numArrows - sum(res)
        return res


print(Solution().maximumBobPoints(numArrows=9, aliceArrows=[1, 1, 0, 1, 0, 0, 2, 1, 0, 1, 2, 0]))
