from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个由正整数组成的数组 nums 和一个 正 整数 k 。

# 如果 nums 的子集中，任意两个整数的绝对差均不等于 k ，则认为该子数组是一个 美丽 子集。

# 返回数组 nums 中 非空 且 美丽 的子集数目。

# nums 的子集定义为：可以经由 nums 删除某些元素（也可能不删除）得到的一个数组。只有在删除元素时选择的索引不同的情况下，两个子集才会被视作是不同的子集。


class Solution:
    def beautifulSubsets(self, nums: List[int], k: int) -> int:
        ...
