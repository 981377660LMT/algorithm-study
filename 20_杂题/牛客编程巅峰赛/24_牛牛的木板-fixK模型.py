from typing import List


class Solution:
    def solve(self, n: int, k: int, nums: List[int]) -> int:
        """fix k 模型"""
        res = 0
        left = 0
        for right in range(len(nums)):
            if nums[right] == 0:
                k -= 1
            while k < 0:
                k += int(nums[left] == 0)
                left += 1
            res = max(res, right - left + 1)
        return res


print(Solution().solve(6, 1, [1, 0, 0, 1, 1, 1]))
