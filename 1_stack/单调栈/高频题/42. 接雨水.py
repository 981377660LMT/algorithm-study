# 每个柱子顶部可以储水的高度为：
# 该柱子的左右两侧最大高度的较小者减去此柱子的高度。
from itertools import accumulate
from typing import List


class Solution:
    def trap(self, height: List[int]) -> int:
        preMax = list(accumulate(height, max))
        sufMax = list(accumulate(height[::-1], max))[::-1]
        return sum(min(h1, h2) - h0 for h0, h1, h2 in zip(height, preMax, sufMax))

