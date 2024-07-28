# 1526. 形成目标数组的子数组最少增加次数
# https://leetcode.cn/problems/minimum-number-of-increments-on-subarrays-to-form-a-target-array/
# 给你一个整数数组 target 和一个数组 initial ，initial 数组与 target  数组有同样的维度，且一开始全部为 0 。
# 请你返回从 initial 得到  target 的最少操作次数，每次操作需遵循以下规则：
# !在 initial 中选择 任意 子数组，并将子数组中每个元素增加 1 。

from itertools import pairwise
from typing import List


class Solution:
    def minNumberOperations(self, target: List[int]) -> int:
        diff = [0] + [b - a for a, b in pairwise([0] + target)]
        return sum(max(0, d) for d in diff)
