from itertools import accumulate
from typing import List

# 给定一个数组 A，将其划分为两个连续子数组 left 和 right
# left 中的每个元素都小于或等于 right 中的每个元素。
# left 和 right 都是非空的。
# left 的长度要尽可能小。
# 在完成这样的分组后返回 left 的长度


# 不检验 all(L <= R for L in left for R in right)，而是检验 max(left) <= min(right)
# 辅助数组：扫两遍，记录leftMax和rightMin  => 类似接雨水那题的数组
class Solution:
    def partitionDisjoint(self, nums: List[int]) -> int:
        preMax = list(accumulate(nums, max))
        sufMin = list(accumulate(nums[::-1], min))[::-1]
        return next(i + 1 for i in range(len(nums)) if preMax[i] <= sufMin[i + 1])


print(Solution().partitionDisjoint([1, 1, 1, 0, 6, 12]))
print(Solution().partitionDisjoint([5, 0, 3, 8, 6]))
