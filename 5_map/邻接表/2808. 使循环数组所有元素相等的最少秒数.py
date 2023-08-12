# https://leetcode.cn/problems/minimum-seconds-to-equalize-a-circular-array/description/


# 看成是扩散元素
# 给你一个下标从 0 开始长度为 n 的数组 nums 。
# 每一秒，你可以对数组执行以下操作：
# !对于范围在 [0, n - 1] 内的每一个下标 i ，将 nums[i] 替换成 nums[i] ，nums[(i - 1 + n) % n] 或者 nums[(i + 1) % n] 三者之一。
# !注意，所有元素会被同时替换。
# 请你返回将数组 nums 中所有元素变成相等元素所需要的 最少 秒数。
# 1<=n<=1e5 nums[i]<=1e9


# 1. 枚举最后变成的数
# 2. 这个数扩散到各个位置的时间

from collections import defaultdict
from typing import List


class Solution:
    def minimumSeconds(self, nums: List[int]) -> int:
        mp = defaultdict(list)
        for i, num in enumerate(nums):
            mp[num].append(i)

        n = len(nums)
        res = n
        for pos in mp.values():
            # 因为是环形数组，所以加一个数
            maxDist = max((b - a) // 2 for a, b in zip(pos, pos[1:] + [pos[0] + n]))
            res = min(res, maxDist)
        return res
