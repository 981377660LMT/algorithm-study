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
# 状态是：dp[leftrow][in][ex][laststate]。

# 第二个方法的思路是，一个格子一个格子自的搜索（而非一行一行的搜索）。每搜索一个格子，它只可能影响它左边的和上边的格子。
# 状态是：dp[x][y][in][ex][laststate]。
# 所以我们只需要知道最近的n个格子的填入情况，就可以统计幸福度的变化。


class Solution:
    def getMaxGridHappiness(
        self, m: int, n: int, introvertsCount: int, extrovertsCount: int
    ) -> int:
        @lru_cache(None)
        def dfs(x, y, intro, extro, state: Tuple[int]) -> int:
            # state:
            # 最近的n个安排，再多不会被当前的影响到
            # 上面的人是state[0]，左边的人是state[-1]

            if y == n:
                return dfs(x + 1, 0, intro, extro, state)
            if x == m:
                return 0

            # 不填
            res = dfs(x, y + 1, intro, extro, state[1:] + (0,))

            if intro:
                score = 120
                if state[0] == 1:
                    score -= 30 * 2
                elif state[0] == 2:
                    # 一个内向一个外向
                    score -= 10
                if y:
                    if state[-1] == 1:
                        score -= 30 * 2
                    elif state[-1] == 2:
                        score -= 10
                res = max(res, score + dfs(x, y + 1, intro - 1, extro, state[1:] + (1,)))

            if extro:
                score = 40
                if state[0] == 1:
                    score -= 10
                elif state[0] == 2:
                    score += 40
                if y:
                    if state[-1] == 1:
                        score -= 10
                    elif state[-1] == 2:
                        score += 40
                res = max(res, score + dfs(x, y + 1, intro, extro - 1, state[1:] + (2,)))

            return res

        return dfs(0, 0, introvertsCount, extrovertsCount, tuple([0] * n))


print(Solution().getMaxGridHappiness(m=2, n=3, introvertsCount=1, extrovertsCount=2))
