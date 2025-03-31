# !3489. 零数组变换 IV - bitset优化可行性01背包
# https://leetcode.cn/problems/zero-array-transformation-iv/description/
#
# 给你一个长度为 n 的整数数组 nums 和一个二维数组 queries ，其中 queries[i] = [li, ri, vali]。
# 每个 queries[i] 表示以下操作在 nums 上执行：
# 从数组 nums 中选择范围 [li, ri] 内的一个下标子集。
# 将每个选中下标处的值减去 正好 vali。
# 零数组 是指所有元素都等于 0 的数组。
# 返回使得经过前 k 个查询（按顺序执行）后，nums 转变为 零数组 的最小可能 非负 值 k。如果不存在这样的 k，返回 -1。

from typing import List


class Solution:
    def minZeroArray(self, nums: List[int], queries: List[List[int]]) -> int:
        res = 0
        for i, x in enumerate(nums):
            if x == 0:
                continue
            dp = 1
            for j, (l, r, v) in enumerate(queries):
                if not l <= i <= r:
                    continue
                dp |= dp << v
                if (dp >> x) & 1:
                    res = max(res, j + 1)
                    break
            else:
                return -1
        return res
