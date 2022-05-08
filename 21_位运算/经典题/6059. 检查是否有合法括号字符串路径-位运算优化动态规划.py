from functools import lru_cache
from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= m, n <= 100

# 6059. 检查是否有合法括号字符串路径-位运算优化动态规划
class Solution:
    def hasValidPath(self, grid: List[List[str]]) -> bool:
        # 用一个数记录每个点所有可能的状态
        ROW, COL = len(grid), len(grid[0])
        dp = [0] * COL
        dp[0] = 1

        for row in grid:
            for c in range(COL):
                # preSum.append(preSum[-1]+nums[i])
                if c:
                    dp[c] |= dp[c - 1]  # preSum[-1]
                if row[c] == '(':  # nums[i]
                    dp[c] <<= 1
                else:
                    dp[c] >>= 1

        return not not dp[-1] & 1


print(
    Solution().hasValidPath(
        grid=[["(", "(", "("], [")", "(", ")"], ["(", "(", ")"], ["(", "(", ")"]]
    )
)

