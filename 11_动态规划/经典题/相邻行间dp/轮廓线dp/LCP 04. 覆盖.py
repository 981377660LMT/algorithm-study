# 棋盘上有一些格子已经坏掉了。你还有无穷块大小为1 * 2的多米诺骨牌，你想把这些骨牌不重叠地覆盖在完好的格子上，请找出你最多能在棋盘上放多少块骨牌
# 多米诺覆盖经典问题 ：匈牙利算法/轮廓线dp

# 1 <= n <= 8
# 1 <= m <= 8
from functools import cache
from typing import List, Tuple


class Solution:
    def domino(self, n: int, m: int, broken: List[List[int]]) -> int:
        """轮廓线DP，复杂度n * m * (2 ^ m)"""

        @cache
        def dfs(r: int, c: int, cState: Tuple[int, ...]) -> int:
            if r == ROW:
                return 0
            if c == COL:
                return dfs(r + 1, 0, cState)

            res = dfs(r, c + 1, cState[1:] + (0,))  # 不选择当前位置

            # 横/竖放置 规定横竖放的方向是右=>左，下=>上
            if (r, c) not in broken:
                if r and (r - 1, c) not in broken and cState[0] == 0:
                    res = max(res, dfs(r, c + 1, cState[1:] + (1,)) + 1)
                if c and (r, c - 1) not in broken and cState[-1] == 0:
                    res = max(res, dfs(r, c + 1, cState[1:-1] + (1, 1)) + 1)

            return res

        ROW, COL = n, m
        if COL > ROW:
            ROW, COL = COL, ROW
            broken = [(c, r) for r, c in broken]
        broken = set((r, c) for r, c in broken)  # type: ignore
        res = dfs(0, 0, tuple([1] * COL))  # 一开始放不了
        dfs.cache_clear()
        return res


print(Solution().domino(n=3, m=3, broken=[]))
