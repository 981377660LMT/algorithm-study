from functools import lru_cache
from typing import Tuple

# 内向的人 开始 时有 120 个幸福感，但每存在一个邻居（内向的或外向的）他都会 失去  30 个幸福感。
# 外向的人 开始 时有 40 个幸福感，每存在一个邻居（内向的或外向的）他都会 得到  20 个幸福感。
# 请你决定网格中应当居住多少人，并为每个人分配一个网格单元。 注意，不必 让所有人都生活在网格中。
# 网格幸福感 是每个人幸福感的 总和 。 返回 最大可能的网格幸福感 。
# 1 <= m, n <= 5
# 0 <= introvertsCount, extrovertsCount <= min(m * n, 6)

# https://leetcode-cn.com/problems/maximize-grid-happiness/solution/dong-tai-gui-hua-shi-jian-qia-de-you-dian-er-jin-y/
# 第一个思路是，一行一行的进行安排。
# 状态是：dp[rowIndex][preState][in][ex]。

# 第二个方法的思路是，一个格子一个格子自的搜索（而非一行一行的搜索）。
# !每搜索一个格子，它只可能影响它左边的和上边的格子。
# !所以我们只需要知道最近的n个格子的填入情况，就可以统计幸福度的变化。
# 状态是：dp[r][c][laststate][in][ex]。


class Solution:
    def getMaxGridHappiness(
        self, ROW: int, COL: int, introvertsCount: int, extrovertsCount: int
    ) -> int:
        @lru_cache(None)
        def dfs(r: int, c: int, cState: Tuple[int], intro: int, extro: int) -> int:
            # state:最近的COL个安排，再多不会被当前的影响到
            # 上面的人是state[0]，左边的人是state[-1]
            # !轮廓线dp+用元组来保存三进制状态

            if r == ROW:
                return 0
            if c == COL:
                return dfs(r + 1, 0, cState, intro, extro)

            res = dfs(r, c + 1, cState[1:] + (0,), intro, extro)  # 不选择当前位置

            if intro:
                score = 120
                if cState[0] == 1:
                    score -= 30 * 2
                elif cState[0] == 2:
                    # 一个内向一个外向
                    score -= 10
                if c:
                    if cState[-1] == 1:
                        score -= 30 * 2
                    elif cState[-1] == 2:
                        score -= 10
                res = max(res, score + dfs(r, c + 1, cState[1:] + (1,), intro - 1, extro))

            if extro:
                score = 40
                if cState[0] == 1:
                    score -= 10
                elif cState[0] == 2:
                    score += 40
                if c:
                    if cState[-1] == 1:
                        score -= 10
                    elif cState[-1] == 2:
                        score += 40
                res = max(res, score + dfs(r, c + 1, cState[1:] + (2,), intro, extro - 1,))

            return res

        ROW, COL = sorted((ROW, COL), reverse=True)
        return dfs(0, 0, tuple([0] * COL), introvertsCount, extrovertsCount)


print(Solution().getMaxGridHappiness(ROW=2, COL=3, introvertsCount=1, extrovertsCount=2))
