from typing import List
from math import gcd
from functools import lru_cache

# 请你返回 n 次操作后你能获得的分数和最大为多少。
# 在第 i 次操作时（操作编号从 1 开始），你需要：

# 选择两个元素 x 和 y 。
# 获得分数 i * gcd(x, y) 。
# 将 x 和 y 从 nums 中删除。

# 1 <= n <= 7
# nums.length == 2 * n


class Solution:
    def maxScore(self, nums: List[int]) -> int:
        n = len(nums)
        target = (1 << n) - 1

        @lru_cache(None)
        def dfs(steps: int, visited: int) -> int:
            """Return maximum score at kth operation with available numbers by mask."""
            if visited == target:
                return 0

            res = -0x7FFFFFFF
            for i in range(n):
                if visited & (1 << i):
                    continue
                for j in range(i + 1, n):
                    if visited & (1 << j):
                        continue
                    nextVisited = visited | (1 << i) | (1 << j)
                    res = max(res, steps * gcd(nums[i], nums[j]) + dfs(steps + 1, nextVisited))

            return res

        return dfs(1, 0)


print(Solution().maxScore(nums=[3, 4, 6, 8]))
# 输出：11
# 解释：最优操作是：
# (1 * gcd(3, 6)) + (2 * gcd(4, 8)) = 3 + 8 = 11
