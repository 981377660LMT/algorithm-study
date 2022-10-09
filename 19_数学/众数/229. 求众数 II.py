# !给定一个大小为 n 的整数数组，找出其中所有出现超过 ⌊ n/3 ⌋ 次的元素。
# 尝试设计时间复杂度为 O(n)、空间复杂度为 O(1)的算法解决此问题
from collections import Counter
from typing import List


class Solution:
    def majorityElement(self, nums: List[int]) -> List[int]:
        """counter肯定是最直接的
        
        O(n)
        """
        counter = Counter(nums)
        res = []
        for num, count in counter.items():
            if count > len(nums) // 3:
                res.append(num)
        return res

    def majorityElement2(self, nums: List[int]) -> List[int]:
        """摩尔投票法"""
        res = []
        count1, count2 = 0, 0
        cand1, cand2 = int(1e20), int(1e30)
        for num in nums:
            if num == cand1:
                count1 += 1
            elif num == cand2:
                count2 += 1
            elif count1 == 0:
                cand1 = num
                count1 = 1
            elif count2 == 0:
                cand2 = num
                count2 = 1
            else:
                count1 -= 1
                count2 -= 1
        if nums.count(cand1) > len(nums) // 3:
            res.append(cand1)
        if nums.count(cand2) > len(nums) // 3:
            res.append(cand2)
        return res
