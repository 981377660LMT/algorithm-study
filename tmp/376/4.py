from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 和一个整数 k 。

# 你可以对数组执行 至多 k 次操作：

# 从数组中选择一个下标 i ，将 nums[i] 增加 或者 减少 1 。
# 最终数组的频率分数定义为数组中众数的 频率 。

# 请你返回你可以得到的 最大 频率分数。


# 众数指的是数组中出现次数最多的数。一个元素的频率指的是数组中这个元素的出现次数。
class Solution:
    def maxFrequencyScore(self, nums: List[int], k: int) -> int:
        if k == 0:
            counter = Counter(nums)
            return max(counter.values())

        nums.sort()
        cands = [(a + b) // 2 for a, b in zip(nums, nums[1:])]
