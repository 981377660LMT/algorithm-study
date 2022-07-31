from typing import List

# 在无序数组中能被二分搜索命中的个数
INF = int(1e20)


class Solution:
    """nlogn"""

    def solve(self, nums: List[int]):
        def bisect(left: int, right: int, lower: int, upper: int) -> None:
            nonlocal res
            if left <= right:
                mid = (left + right) // 2
                # If the mid value is between those two, it has correctly binary searched up until this point
                if lower <= nums[mid] <= upper:
                    res += 1
                if mid - 1 >= 0:
                    bisect(left, mid - 1, lower, min(upper, nums[mid]))
                if mid + 1 < len(nums):
                    bisect(mid + 1, right, max(lower, nums[mid]), upper)

        if len(nums) == 1:
            return 1

        res = 0
        bisect(0, len(nums) - 1, -INF, INF)
        return res
