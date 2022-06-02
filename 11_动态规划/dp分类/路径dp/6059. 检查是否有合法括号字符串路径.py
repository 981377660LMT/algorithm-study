from typing import List
from functools import lru_cache

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= m, n <= 100

# DFS 的写法相比某些递推的写法要快 10 倍以上，这是因为有很多状态是无法访问到的
# https://leetcode-cn.com/problems/check-if-there-is-a-valid-parentheses-string-path/solution/tian-jia-zhuang-tai-hou-dfscpythonjavago-f287/
class Solution:
    def hasValidPath(self, grid: List[List[str]]) -> bool:
        """一个括号字符串是一个 非空 且只包含 '(' 和 ')' 的字符串。判断这个括号字符串是否是 合法的 。"""

        @lru_cache(None)
        def dfs(row: int, col: int, diff: int) -> bool:
            if diff < 0:  # !不合题意
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

