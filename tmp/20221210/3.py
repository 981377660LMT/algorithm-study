from itertools import pairwise
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 stones ，数组中的元素 严格递增 ，表示一条河中石头的位置。

# 一只青蛙一开始在第一块石头上，它想到达最后一块石头，然后回到第一块石头。同时每块石头 至多 到达 一次。

# 一次跳跃的 长度 是青蛙跳跃前和跳跃后所在两块石头之间的距离。

# 更正式的，如果青蛙从 stones[i] 跳到 stones[j] ，跳跃的长度为 |stones[i] - stones[j]| 。
# 一条路径的 代价 是这条路径里的 最大跳跃长度 。

# 请你返回这只青蛙的 最小代价 。


class Solution:
    def maxJump(self, stones: List[int]) -> int:
        if len(stones) == 2:
            return stones[1] - stones[0]

        even = stones[::2]
        odd = stones[1::2]
        res1 = 0
        for a, b in pairwise(even):
            res1 = max(res1, b - a)
        res2 = 0
        for a, b in pairwise(odd):
            res2 = max(res2, b - a)
        return max(res1, res2)


print(Solution().maxJump([1, 2, 3, 4, 5]))
