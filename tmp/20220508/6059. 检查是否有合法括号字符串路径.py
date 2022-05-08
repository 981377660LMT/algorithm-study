from functools import lru_cache
from typing import List, Tuple
from collections import defaultdict

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= m, n <= 100


class Solution:
    def hasValidPath(self, grid: List[List[str]]) -> bool:
        @lru_cache(None)
        def dfs(row: int, col: int, diff: int) -> bool:
            if diff < 0:
                return False
            if (row, col) == (ROW - 1, COL - 1):
                return diff == 0

            res = False
            for dr, dc in [(0, 1), (1, 0)]:
                nr, nc = row + dr, col + dc
                if 0 <= nr < ROW and 0 <= nc < COL:
                    # 反思：注意不要在函数外改变diff 要把diff改变写在函数内
                    res = res or dfs(nr, nc, diff + (1 if grid[nr][nc] == '(' else -1))
            return res

        ROW, COL = len(grid), len(grid[0])
        # (not (ROW + COL) & 1) 这个可以剪掉很多枝
        if grid[0][0] == ')' or grid[-1][-1] == '(' or (not (ROW + COL) & 1):
            return False
        res = dfs(0, 0, 1)
        dfs.cache_clear()
        return res


print(
    Solution().hasValidPath(
        grid=[["(", "(", "("], [")", "(", ")"], ["(", "(", ")"], ["(", "(", ")"]]
    )
)

print(Solution().hasValidPath(grid=[["(", ")"], ["(", ")"]]))
# False
print(Solution().hasValidPath(grid=[["("]]))
print(
    Solution().hasValidPath(
        [
            ["(", "(", ")", "(", "(", ")", "(", ")", "("],
            [")", "(", "(", "(", ")", ")", ")", "(", "("],
            ["(", ")", ")", "(", "(", ")", "(", "(", "("],
            ["(", "(", ")", "(", ")", "(", "(", ")", "("],
            ["(", ")", "(", ")", ")", ")", "(", ")", ")"],
        ]
    )
)

