from typing import List


class Solution:
    def longestNiceSubarray(self, nums: List[int]) -> int:
        """如果 nums 的子数组中位于 不同 位置的每对元素按位 与（AND）运算的结果等于 0 ，
        则称该子数组为 优雅 子数组。"""
        res, left, n = 1, 0, len(nums)
        counter = [0] * 32
        for right in range(n):
            for bit in range(32):
                if (nums[right] >> bit) & 1:
                    counter[bit] += 1
            while left <= right and any(counter[i] > 1 for i in range(32)):
                for bit in range(32):
                    if (nums[left] >> bit) & 1:
                        counter[bit] -= 1
                left += 1
            res = max(res, right - left + 1)
        return res


print(Solution().longestNiceSubarray([1, 3, 8, 48, 10]))
