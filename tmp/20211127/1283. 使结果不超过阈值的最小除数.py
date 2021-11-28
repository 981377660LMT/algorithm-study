from typing import List
from math import ceil

# 你需要选择一个正整数作为除数，然后将数组里每个数都除以它，并对除法结果求和。
# 请你找出能够使上述结果小于等于阈值 threshold 的除数中 最小 的那个。
# 1 <= nums.length <= 5 * 10^4
# 1 <= nums[i] <= 10^6

# 一眼二分答案


class Solution:
    def smallestDivisor(self, nums: List[int], threshold: int) -> int:
        left, right = 1, max(nums)

        while left <= right:
            mid = (left + right) >> 1
            if sum(ceil(x / mid) for x in nums) <= threshold:
                right = mid - 1
            else:
                left = mid + 1

        return left


print(Solution().smallestDivisor(nums=[1, 2, 5, 9], threshold=6))
# 输出：5
# 解释：如果除数为 1 ，我们可以得到和为 17 （1+2+5+9）。
# 如果除数为 4 ，我们可以得到和为 7 (1+1+2+3) 。如果除数为 5 ，和为 5 (1+1+1+2)。

