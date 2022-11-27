# 0 ≤ n ≤ 100,000

# 795. 区间子数组个数
# 最大值在范围内的子数组个数
# atMostK
from typing import List


MOD = int(1e9 + 7)


class Solution:
    def numSubarrayBoundedMax(self, nums: List[int], left: int, right: int) -> int:
        """解法1:容斥"""

        def atMostK(k: int) -> int:
            """子数组的最大值不大于k的子数组个数"""
            res, dp = 0, 0
            for num in nums:
                if num <= k:
                    dp += 1
                else:
                    dp = 0
                res += dp
            return res

        return (atMostK(right) - atMostK(left - 1)) % MOD

    def numSubarrayBoundedMax2(self, nums: List[int], lower: int, upper: int) -> int:
        """解法2:定界子数组"""

        n = len(nums)
        res = 0
        pos1, pos2 = -1, -1
        for right in range(n):
            if lower <= nums[right] <= upper:
                pos1 = right
            if nums[right] > upper:
                pos2 = right
                pos1 = -1
            if pos1 != -1:
                res += pos1 - pos2
        return res


print(Solution().numSubarrayBoundedMax(nums=[1, 5, 3, 2], left=1, right=4))
print(Solution().numSubarrayBoundedMax2(nums=[1, 5, 3, 2], lower=1, upper=4))

# We have the following sublists where 1 ≤ max(A) ≤ 4

# [1]
# [3]
# [3, 2]
# [2]
