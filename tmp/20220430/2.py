from typing import List
from itertools import accumulate


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minimumAverageDifference(self, nums: List[int]) -> int:
        preSum = [0] + list(accumulate(nums))
        sum_ = preSum[-1]
        min_ = int(1e20)
        res = 0

        # 前几个数
        for i in range(1, len(nums) + 1):
            left = preSum[i] // i
            right = (sum_ - preSum[i]) // (len(nums) - i) if i < len(nums) else 0
            diff = abs(left - right)
            if diff < min_:
                min_ = diff
                res = i - 1

        return res


print(Solution().minimumAverageDifference([0, 1, 0, 1, 0, 1]))

# 0
