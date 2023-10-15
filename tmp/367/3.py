from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始、长度为 n 的整数数组 nums ，以及整数 indexDifference 和整数 valueDifference 。

# 你的任务是从范围 [0, n - 1] 内找出  2 个满足下述所有条件的下标 i 和 j ：

# abs(i - j) >= indexDifference 且
# abs(nums[i] - nums[j]) >= valueDifference
# 返回整数数组 answer。如果存在满足题目要求的两个下标，则 answer = [i, j] ；否则，answer = [-1, -1] 。如果存在多组可供选择的下标对，只需要返回其中任意一组即可。


# 注意：i 和 j 可能 相等 。


class Solution:
    def findIndices(self, nums: List[int], indexDifference: int, valueDifference: int) -> List[int]:
        if indexDifference >= len(nums):
            return [-1, -1]
        if valueDifference == 0:
            return [0, indexDifference]
        left = SortedList()
        for i in range(indexDifference, len(nums)):
            left.add((nums[i - indexDifference], i - indexDifference))
            min_ = left[0][0]
            if abs(min_ - nums[i]) >= valueDifference:
                return [left[0][1], i]
            max_ = left[-1][0]
            if abs(max_ - nums[i]) >= valueDifference:
                return [left[-1][1], i]

        return [-1, -1]


# [0,1]
#
# [15,8]
# 0
# 6
print(Solution().findIndices([15, 8, 2, 3, 1, 10], 0, 6))
#
# 2
# 4
print(Solution().findIndices([2, 0, 9, 2], 2, 4))
