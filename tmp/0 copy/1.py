from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 。

# 如果一个前缀 nums[0..i] 满足对于 1 <= j <= i 的所有元素都有 nums[j] = nums[j - 1] + 1 ，那么我们称这个前缀是一个 顺序前缀 。特殊情况是，只包含 nums[0] 的前缀也是一个 顺序前缀 。


# 请你返回 nums 中没有出现过的 最小 整数 x ，满足 x 大于等于 最长 顺序前缀的和。
class Solution:
    def missingInteger(self, nums: List[int]) -> int:
        dp = 1
        pre = nums[0]
        for i in range(1, len(nums)):
            if nums[i] == pre + 1:
                dp += 1
                pre = nums[i]
            else:
                break
        dp = sum(nums[:dp])
        while dp in nums:
            dp += 1
        return dp
