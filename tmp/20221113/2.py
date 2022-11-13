from math import lcm
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums 和一个整数 k ，请你统计并返回 nums 的 子数组 中满足 元素最小公倍数为 k 的子数组数目。

# 子数组 是数组中一个连续非空的元素序列。

# 数组的最小公倍数 是可被所有数组元素整除的最小正整数。

# !注意要break lcm一直算会有溢出风险
class Solution:
    def subarrayLCM(self, nums: List[int], k: int) -> int:
        res = 0
        for i in range(len(nums)):
            lcm_ = 1  # 幺元
            for j in range(i, len(nums)):
                lcm_ = lcm(lcm_, nums[j])
                if lcm_ == k:
                    res += 1
        return res
