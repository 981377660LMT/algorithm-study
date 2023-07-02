from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 和一个整数 threshold 。

# 请你从 nums 的子数组中找出以下标 l 开头、下标 r 结尾 (0 <= l <= r < nums.length) 且满足以下条件的 最长子数组 ：

# nums[l] % 2 == 0
# 对于范围 [l, r - 1] 内的所有下标 i ，nums[i] % 2 != nums[i + 1] % 2
# 对于范围 [l, r] 内的所有下标 i ，nums[i] <= threshold
# 以整数形式返回满足题目要求的最长子数组的长度。


# 注意：子数组 是数组中的一个连续非空元素序列。
class Solution:
    def longestAlternatingSubarray(self, nums: List[int], threshold: int) -> int:
        res = 0
        # 枚举子数组
        for l in range(len(nums)):
            for r in range(l, len(nums)):
                if nums[l] % 2 == 0:
                    ok = True
                    for i in range(l, r):
                        if nums[i] % 2 == nums[i + 1] % 2:
                            ok = False
                            break
                    for i in range(l, r + 1):
                        if nums[i] > threshold:
                            ok = False
                            break
                    if ok:
                        res = max(res, r - l + 1)
        return res
