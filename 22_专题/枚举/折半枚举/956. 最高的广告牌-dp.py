from typing import List
from collections import defaultdict

# 对于每一根钢筋 x，我们会写下 +x，-x 或者 0。
# 我们的目标是最终得到结果 0 并让正数之和最大。
# 我们记所有写下的正数之和为 score。例如，+1 +2 +3 -6 的 score 为 6。


class Solution:
    def tallestBillboard(self, rods: List[int]) -> int:
        dp = defaultdict(int, {0: 0})
        for num in rods:
            for preSum, maxSum in list(dp.items()):
                dp[preSum + num] = max(dp[preSum + num], maxSum + num)
                dp[preSum - num] = max(dp[preSum - num], maxSum)
        return dp[0]
