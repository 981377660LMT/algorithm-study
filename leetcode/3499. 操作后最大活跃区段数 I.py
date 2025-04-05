# 3502. 到达每个位置的最小费用
from itertools import accumulate


class Solution:
    def minCosts(self, cost: List[int]) -> List[int]:
        return list(accumulate(cost, min))


# 100001110001101110001
