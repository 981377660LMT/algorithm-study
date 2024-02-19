# 非常像 1218. 最长定差子序列
# 100205. 修改数组后最大化数组中的连续元素数目
# https://leetcode.cn/problems/maximize-consecutive-elements-in-an-array-after-modification/description/
# 给你一个下标从 0 开始只包含 正 整数的数组 nums 。
# 一开始，你可以将数组中 任意数量 元素增加 至多 1 。
# 修改后，你可以从最终数组中选择 一个或者更多 元素，并确保这些元素升序排序后是 连续 的。比方说，[3, 4, 5] 是连续的，但是 [3, 4, 6] 和 [1, 1, 2, 3] 不是连续的。
# 请你返回 最多 可以选出的元素数目。


from collections import defaultdict
from typing import List


class Solution:
    def maxSelectedElements(self, nums: List[int]) -> int:
        nums.sort()  # 排序消除后效性
        dp = defaultdict(int)  # dp[x] 表示子序列最后一个元素为 x 时，最长连续子序列的长度
        for x in nums:
            dp[x], dp[x + 1] = dp[x - 1] + 1, dp[x] + 1
        return max(dp.values())
