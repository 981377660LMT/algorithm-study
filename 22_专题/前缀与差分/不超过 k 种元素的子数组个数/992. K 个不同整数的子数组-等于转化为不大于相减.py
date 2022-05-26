# 求子数组中有k中元素的子数组个数

from collections import defaultdict
from typing import List


def atMostK(nums: List[int], k: int) -> int:
    """元素个数<=k的子数组个数"""
    left, res, counter = 0, 0, defaultdict(int)
    for right in range(len(nums)):
        counter[nums[right]] += 1
        while len(counter) > k:
            counter[nums[left]] -= 1
            if counter[nums[left]] == 0:
                del counter[nums[left]]
            left += 1
        res += right - left + 1
    return res


class Solution:
    def subarraysWithKDistinct(self, nums: List[int], k: int) -> int:
        return atMostK(nums, k) - atMostK(nums, k - 1)

