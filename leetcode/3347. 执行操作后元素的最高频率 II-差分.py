# 3347. 执行操作后元素的最高频率 II-差分
# https://leetcode.cn/problems/maximum-frequency-of-an-element-after-performing-operations-ii/description/
# 给你一个整数数组 nums 和两个整数 k 和 numOperations 。
# 你必须对 nums 执行 操作  numOperations 次。每次操作中，你可以：
# 选择一个下标 i ，它在之前的操作中 没有 被选择过。
# 将 nums[i] 增加范围 [-k, k] 中的一个整数。
# 在执行完所有操作以后，请你返回 nums 中出现 频率最高 元素的出现次数。
# 一个元素 x 的 频率 指的是它在数组中出现的次数。

from collections import defaultdict
from typing import List


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxFrequency(self, nums: List[int], k: int, numOperations: int) -> int:
        counter = defaultdict(int)
        diff = defaultdict(int)
        visited = set()

        for v in nums:
            counter[v] += 1
            diff[v - k] += 1
            diff[v + k + 1] -= 1
            visited |= {v - k, v, v + k + 1}

        curSum = 0
        res = 0
        for v in sorted(visited):
            curSum += diff[v]
            count = counter[v]
            delta = curSum - count
            res = max2(res, count + min2(numOperations, delta))
        return res
