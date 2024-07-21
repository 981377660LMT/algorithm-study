from itertools import pairwise
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


# 给你两个长度相同的正整数数组 nums 和 target。
# 在一次操作中，你可以选择 nums 的任何子数组，并将该子数组内的每个元素的值增加或减少 1。
# 返回使 nums 数组变为 target 数组所需的 最少 操作次数。


# 等价差分数组于一个数+1，另一个-1
class Solution:
    def minimumOperations(self, nums: List[int], target: List[int]) -> int:
        diff = [0] + [b - a for a, b in zip(nums, target)]
        diff = [b - a for a, b in pairwise(diff)]
        posSum, negSum = 0, 0
        for d in diff:
            if d > 0:
                posSum += d
            elif d < 0:
                negSum += -d
        return max(posSum, negSum)


# nums = [3,5,1,2], target = [4,6,2,4]

print(Solution().minimumOperations([3, 5, 1, 2], [4, 6, 2, 4]))  # 1
print(Solution().minimumOperations(nums=[1, 3, 2], target=[2, 1, 4]))  # 1
