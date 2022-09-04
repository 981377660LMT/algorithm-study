# 2393. 严格递增的子数组个数
from typing import List


class Solution:
    def countSubarrays(self, nums: List[int]) -> int:
        n, dp, res = len(nums), 0, 0
        for i in range(n):  # 以每个位置为右端点的子数组个数
            if i - 1 < 0 or nums[i - 1] < nums[i]:
                dp += 1
            else:
                dp = 1
            res += dp
        return res


print(Solution().countSubarrays([1, 2, 3, 4]))
print(Solution().countSubarrays([1]))
