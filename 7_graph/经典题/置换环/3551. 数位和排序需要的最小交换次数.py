# https://leetcode.cn/problems/minimum-swaps-to-sort-by-digit-sum/description/

from typing import List


class Solution:
    def minSwaps(self, nums: List[int]) -> int:
        idx = {x: i for i, x in enumerate(nums)}
        a = sorted(nums, key=lambda x: (sum(map(int, str(x))), x))
        res = 0
        for i, x in enumerate(a):
            j = idx[x]
            if j == i:
                continue
            # 交换nums的第i位和第j位
            nums[j] = nums[i]
            idx[nums[j]] = j
            res += 1
        return res
