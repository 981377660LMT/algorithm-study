# 输入n，打印出所有骰子朝上一面的点数之和为s的所有可能的值出现的概率。
# 1 <= n <= 11
from collections import defaultdict
from typing import List


class Solution:
    def dicesProbability(self, n: int) -> List[float]:
        dp = defaultdict(int, {i: 1 for i in range(1, 7)})
        for _ in range(n - 1):
            ndp = defaultdict(int)
            for pre in dp:
                for cur in range(1, 7):
                    ndp[pre + cur] += dp[pre]
            dp = ndp
        sum_ = sum(dp.values())
        return [dp[key] / sum_ for key in sorted(dp)]

