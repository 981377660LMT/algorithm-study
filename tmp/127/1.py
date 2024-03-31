from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个 非负 整数数组 nums 和一个整数 k 。

# 如果一个数组中所有元素的按位或运算 OR 的值 至少 为 k ，那么我们称这个数组是 特别的 。


# 请你返回 nums 中 最短特别非空 子数组的长度，如果特别子数组不存在，那么返回 -1 。
class Solution:
    def minimumSubarrayLength(self, nums: List[int], k: int) -> int:
        res = INF
        for i in range(len(nums)):
            for j in range(i, len(nums)):
                curXor = 0
                for t in range(i, j + 1):
                    curXor |= nums[t]
                if curXor >= k:
                    res = min(res, j - i + 1)
        return res if res != INF else -1
