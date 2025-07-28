# Q3. 灯光调整
# https://leetcode.cn/contest/2025_pudong_ai/problems/MWmcf7/description/
# 灯光师需要调整灯光的亮度。brightness[i] 表示第 i 盏照明灯的初始亮度。
# 灯光师有一个控制器，每次操作可以选择任意 一段连续的照明灯，将其中所有灯的亮度 调高 1 或 降低 1。
# 请返回灯光师使所有灯的亮度一致所需要的 最小操作次数。
#
# 差分数组
# !给很多数，每次可以“选一个数加 1 ，选一个数减 1 ”或者“只选一个数加减 1 ”，求把所有数变成 0 的最少操作次数。

from itertools import pairwise
from typing import List


class Solution:
    def lightAdjustment(self, brightness: List[int]) -> int:
        pos, neg = 0, 0
        for a, b in pairwise(brightness):
            diff = b - a
            if diff > 0:
                pos += diff
            elif diff < 0:
                neg += -diff
        return max(pos, neg)
