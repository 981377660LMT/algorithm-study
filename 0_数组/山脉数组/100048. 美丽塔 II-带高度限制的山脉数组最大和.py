# 100048. 美丽塔 II
# https://leetcode.cn/problems/beautiful-towers-ii/description/
# 每个建筑物都有一个高度限制 maxHeights[i]，这意味着它的高度不能超过这个值。
# 现在要建造一个山脉数组.
# 对于所有 0 < j <= i ，都有 heights[j - 1] <= heights[j].
# 对于所有 i <= k < n - 1 ，都有 heights[k + 1] <= heights[k].
# !求建筑物高度之和的最大值。
# n<=1e5
# maxHeights[i]<=1e9

# 带高度限制的山脉数组最大和
# !单调栈+前缀和
# !前后缀问题尽量用makeDp

from typing import List
from typing import Any
from 每个元素作为最值的影响范围 import getRange


class Solution:
    def maximumSumOfHeights(self, maxHeights: List[int]) -> int:
        def makeDp(seq: List[Any], rev=True) -> List[int]:
            n = len(seq)
            minRange = getRange(seq)
            dp = [0] * (n + 1)
            for i in range(n):
                cur = seq[i]
                j = minRange[i][0]
                dp[i + 1] = dp[j] + (i - j + 1) * cur
            return dp

        pre, suf = makeDp(maxHeights[:], False), makeDp(maxHeights[::-1])[::-1]
        res = 0
        for i in range(len(pre)):  # 枚举分割点
            res = max(res, pre[i] + suf[i])
        return res


assert Solution().maximumSumOfHeights([5, 3, 4, 1, 1]) == 13
assert Solution().maximumSumOfHeights([6, 5, 3, 9, 2, 7]) == 22
assert Solution().maximumSumOfHeights([3, 2, 5, 5, 2, 3]) == 18
