"""两数之和找到所有数对"""


from collections import defaultdict
from typing import List


class Solution:
    def pairSums(self, nums: List[int], target: int) -> List[List[int]]:
        """设计一个算法，找出数组中两数之和为指定值的所有整数对。一个数只能属于一个数对。"""
        res = []
        counter = defaultdict(int)
        for num in nums:
            need = target - num
            if counter[need] > 0:
                counter[need] -= 1
                res.append([need, num])
            else:
                counter[num] += 1
        return res
