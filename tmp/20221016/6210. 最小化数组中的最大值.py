"""
6210. 最小化数组中的最大值

请你返回可以得到的 nums 数组中 最大值 最小 为多少。
"""

from itertools import accumulate
from math import ceil
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)

# 选择一个满足 1 <= i < n 的整数 i ，且 nums[i] > 0 。
# 将 nums[i] 减 1 。
# 将 nums[i - 1] 加 1 。
# 你可以对数组执行 任意 次上述操作，请你返回可以得到的 nums 数组中 最大值 最小 为多少。


class Solution:
    def minimizeArrayValue2(self, nums: List[int]) -> int:
        """贪心:最小化的话就要尽量每个前缀都持平"""
        preSum = [0] + list(accumulate(nums))
        res = -INF
        for i in range(1, len(preSum)):
            res = max(res, ceil(preSum[i] / i))
        return res

    def minimizeArrayValue(self, nums: List[int]) -> int:
        def check(mid: int) -> bool:
            """每个数可以往左边加 是否能使得最大值最小为mid"""
            over = 0
            for i in range(n):
                if nums[i] > mid:
                    delta = nums[i] - mid
                    if delta > over:
                        return False
                    over -= delta
                else:
                    over += mid - nums[i]
            return True

        n = len(nums)
        left, right = 0, int(1e10)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


print(Solution().minimizeArrayValue([0, 100, 90]))
print(Solution().minimizeArrayValue2([0, 100, 90]))
