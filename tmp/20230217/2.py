from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 。

# nums 的 最小 得分是满足 0 <= i < j < nums.length 的 |nums[i] - nums[j]| 的最小值。
# nums的 最大 得分是满足 0 <= i < j < nums.length 的 |nums[i] - nums[j]| 的最大值。
# nums 的分数是 最大 得分与 最小 得分的和。
# 我们的目标是最小化 nums 的分数。你 最多 可以修改 nums 中 2 个元素的值。

# 请你返回修改 nums 中 至多两个 元素的值后，可以得到的 最小分数 。

# |x| 表示 x 的绝对值。
class Solution:
    def minimizeSum(self, nums: List[int]) -> int:
        nums = sorted(nums)
        if len(nums) <= 3:
            return 0
        return min(nums[-1] - nums[2], nums[-2] - nums[1], nums[-3] - nums[0])


print(Solution().minimizeSum([31, 25, 72, 79, 74, 65]))
