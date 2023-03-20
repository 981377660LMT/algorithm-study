from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 。你需要将 nums 重新排列成一个新的数组 perm 。

# 定义 nums 的 伟大值 为满足 0 <= i < nums.length 且 perm[i] > nums[i] 的下标数目。

# 请你返回重新排列 nums 后的 最大 伟大值。

# 田忌赛马
class Solution:
    def maximizeGreatness(self, nums: List[int]) -> int:
        pool = SortedList(nums)
        res = 0
        for i in range(len(nums)):
            cur = nums[i]
            # 找到最小的大于cur的数
            upper = pool.bisect_right(cur)
            if upper < len(pool):
                pool.remove(pool[upper])
                res += 1
            else:
                pool.pop(0)
        return res
