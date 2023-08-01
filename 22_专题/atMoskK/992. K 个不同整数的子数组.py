# 给定一个正整数数组 nums和一个整数 k，返回 nums 中 「好子数组」 的数目。
# !如果 nums 的某个子数组中不同整数的个数恰好为 k，则称 nums 的这个连续、不一定不同的子数组为 「好子数组 」。

from collections import defaultdict
from typing import List


class Solution:
    def subarraysWithKDistinct(self, nums: List[int], k: int) -> int:
        if k <= 0:
            return 0

        def cal(ceiling: int) -> int:
            """计算不同整数的个数小于等于 ceiling 的子数组的个数."""
            res, left, n = 0, 0, len(nums)
            counter = defaultdict(int)
            for right in range(n):
                counter[nums[right]] += 1
                while left <= right and len(counter) > ceiling:
                    removed = nums[left]
                    counter[removed] -= 1
                    if counter[removed] == 0:
                        del counter[removed]
                    left += 1
                res += right - left + 1
            return res

        return cal(k) - cal(k - 1)
