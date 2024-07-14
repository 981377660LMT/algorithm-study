# 1502. 判断能否形成等差数列
# https://leetcode.cn/problems/can-make-arithmetic-progression-from-sequence/description/
# 如果可以重新排列数组形成等差数列，请返回 true ；否则，返回 false

from itertools import pairwise
from typing import List


class Solution:
    def canMakeArithmeticProgression(self, arr: List[int]) -> bool:
        if len(arr) < 3:
            return True
        return len(set(y - x for x, y in pairwise(sorted(arr)))) == 1
