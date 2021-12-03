from typing import List
from functools import lru_cache

# 你挑选 任意 一块披萨。
# Alice 将会挑选你所选择的披萨逆时针方向的下一块披萨。
# Bob 将会挑选你所选择的披萨顺时针方向的下一块披萨。
# 重复上述过程直到没有披萨剩下。
# 每一块披萨的大小按顺时针方向由循环数组 slices 表示。

# 请你返回你可以获得的披萨大小总和的最大值。
# 化简为取n/3个不相邻数字的最大和


class Solution:
    def maxSizeSlices(self, slices: List[int]) -> int:
        count = len(slices) // 3

        @lru_cache(None)
        def dfs(start: int, end: int, remain: int) -> int:
            if remain == 0 or start > end:
                return 0
            # 每一轮取不取start
            return max(dfs(start + 2, end, remain - 1) + slices[start], dfs(start + 1, end, remain))

        return max(dfs(2, len(slices) - 2, count - 1) + slices[0], dfs(1, len(slices) - 1, count))


print(Solution().maxSizeSlices(slices=[1, 2, 3, 4, 5, 6]))
# 输出：10
# 解释：选择大小为 4 的披萨，Alice 和 Bob 分别挑选大小为 3 和 5 的披萨。
# 然后你选择大小为 6 的披萨，Alice 和 Bob 分别挑选大小为 2 和 1 的披萨。
# 你获得的披萨总大小为 4 + 6 = 10 。

