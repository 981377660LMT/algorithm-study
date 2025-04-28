from bisect import bisect_left
from typing import List


class Solution:
    def singleNonDuplicate(self, nums: List[int]) -> int:
        # 支持带key的二分
        return nums[bisect_left(range(len(nums) - 1), True, key=lambda x: nums[x] != nums[x ^ 1])]


print(Solution().singleNonDuplicate(nums=[1, 1, 2, 3, 3, 4, 4, 8, 8]))
print(Solution().singleNonDuplicate(nums=[1, 1, 2]))
