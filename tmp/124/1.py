from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums ，如果 nums 至少 包含 2 个元素，你可以执行以下操作：

# 选择 nums 中的前两个元素并将它们删除。
# 一次操作的 分数 是被删除元素的和。

# 在确保 所有操作分数相同 的前提下，请你求出 最多 能进行多少次操作。


# 请你返回按照上述要求 最多 可以进行的操作次数。
class Solution:
    def maxOperations(self, nums: List[int]) -> int:
        res = 1
        sum_ = nums[0] + nums[1]
        nums = nums[2:]
        while len(nums) >= 2:
            curSum = nums[0] + nums[1]
            nums = nums[2:]
            if curSum == sum_:
                res += 1
            else:
                break
        return res
