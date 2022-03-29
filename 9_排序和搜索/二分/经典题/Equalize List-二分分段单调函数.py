from typing import List

# 二分分段单调函数
# mid可变，求函数sum((nums[i]-mid)*cost[i])的最小值


class Solution:
    def solve(self, nums: List[int], costs: List[int]) -> int:
        """you want to make all elements equal in nums, return the minimum total cost required."""

        def calCost(mid: int) -> int:
            res = 0
            for i, num in enumerate(nums):
                res += costs[i] * abs(num - mid)
            return res

        left, right = 0, max(nums)
        while left <= right:
            mid = (left + right) >> 1
            if calCost(mid) < calCost(mid + 1):
                right = mid - 1
            else:
                left = mid + 1
        return calCost(left)


print(Solution().solve(nums=[2, 1, 3], costs=[1, 10, 2]))
