from typing import List

# 在无序数组中能被二分搜索命中的个数


class Solution:
    """ nlogn"""

    def solve1(self, nums):

        res = 0
        for i in range(len(nums)):
            low = 0
            high = len(nums) - 1
            target = nums[i]
            while low <= high:
                mid = (low + high) // 2
                if nums[mid] == target:
                    res += 1
                    break
                elif nums[mid] > target:
                    high = mid - 1
                else:
                    low = mid + 1

        return res

    def solve(self, nums: List[int]):
        if len(nums) == 1:
            return 1

        res = 0

        def bisect(left: int, right: int, lower: int, upper: int) -> None:
            nonlocal res
            if left <= right:
                mid = left + right >> 1
                # If the mid value is between those two, it has correctly binary searched up until this point
                if lower <= nums[mid] <= upper:
                    res += 1
                if mid - 1 >= 0:
                    bisect(left, mid - 1, lower, min(upper, nums[mid]))
                if mid + 1 < len(nums):
                    bisect(mid + 1, right, max(lower, nums[mid]), upper)

        bisect(0, len(nums) - 1, -int(1e20), int(1e20))
        return res

