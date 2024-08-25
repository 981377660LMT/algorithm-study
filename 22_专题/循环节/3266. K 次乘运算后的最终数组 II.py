# 3266. K 次乘运算后的最终数组 II
# https://leetcode.cn/problems/final-array-state-after-k-multiplication-operations-ii/description/
# 给你一个整数数组 nums ，一个整数 k  和一个整数 multiplier 。
# 你需要对 nums 执行 k 次操作，每次操作中：
# 找到 nums 中的 最小 值 x ，如果存在多个最小值，选择最 前面 的一个。
# 将 x 替换为 x * multiplier 。
# k 次操作以后，你需要将 nums 中每一个数值对 109 + 7 取余。
# 请你返回执行完 k 次乘运算以及取余运算之后，最终的 nums 数组。
#
# !注意到操作一定次数后，会进入循环节.


from typing import List
from sortedcontainers import SortedList


MOD = int(1e9 + 7)


class Solution:
    def getFinalState(self, nums: List[int], k: int, multiplier: int) -> List[int]:
        if multiplier == 1:
            return nums

        n = len(nums)
        times = [0] * n
        sl = SortedList((v, i) for i, v in enumerate(nums))
        remain = k
        while remain:
            v, i = sl[0]
            # !最小值*multiplier>最大值后，开始均匀分配
            if v * multiplier > sl[-1][0]:
                break
            sl.pop(0)
            sl.add((v * multiplier, i))
            times[i] += 1
            remain -= 1

        div, mod = divmod(remain, n)
        for _, i in sl:
            times[i] += div
            if mod:
                times[i] += 1
                mod -= 1
        return [(v * pow(multiplier, t, MOD)) % MOD for v, t in zip(nums, times)]
