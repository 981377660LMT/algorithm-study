from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个整数数组 nums ，数组由 不同正整数 组成，请你找出并返回数组中 任一 既不是 最小值 也不是 最大值 的数字，如果不存在这样的数字，返回 -1 。


# 返回所选整数。
class Solution:
    def findNonMinOrMax(self, nums: List[int]) -> int:
        min_, max_ = min(nums), max(nums)
        for num in nums:
            if num != min_ and num != max_:
                return num
        return -1
