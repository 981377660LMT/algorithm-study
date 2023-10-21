from itertools import accumulate
from typing import List

# 1 <= numArrows <= 1e5
# aliceArrows.length == bobArrows.length == 12
# 我这种写法(模板写法) 遍历了很多不是最优解情形时的容量，真正有意义的背包容量是十二个背包的可能的子集和，最多有(1<<12种），而我这种每次都要遍历numArrows个容量，所以超时


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


# arrows = [2 ** i for i in range(4, 16)]
# print(arrows, sum(arrows))
print(
    Solution().maximumBobPoints(
        numArrows=65520,
        aliceArrows=[16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768],
    )
)
