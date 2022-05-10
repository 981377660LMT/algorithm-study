from typing import List, Tuple
from functools import lru_cache

# 1 <= rows, cols <= 8
# 玩家由字符 'C' （代表猫）和 'M' （代表老鼠）表示。
# 地板由字符 '.' 表示，玩家可以通过这个格子。
# 墙用字符 '#' 表示，玩家不能通过这个格子。
# 食物用字符 'F' 表示，玩家可以通过这个格子。
# 字符 'C' ， 'M' 和 'F' 在 grid 中都只会出现一次。

# 游戏有 4 种方式会结束：

# 如果猫跟老鼠处在相同的位置，那么猫获胜。
# 如果猫先到达食物，那么猫获胜。
# 如果老鼠先到达食物，那么老鼠获胜。
# 如果老鼠不能在 1000 次操作以内到达食物，那么猫获胜。
# https://leetcode.cn/problems/cat-and-mouse-ii/solution/by-ac_oier-gse8/


DIR4 = ((0, 1), (0, -1), (1, 0), (-1, 0))

# 极小极大博弈


class Solution:
    def canMouseWin(self, grid: List[str], catJump: int, mouseJump: int) -> bool:
        ROW, COL = len(grid), len(grid[0])
        walls = set()
        cat = mouse = (0, 0)
        for i in range(ROW):
            for j in range(COL):
                if grid[i][j] == "F":
                    food = (i, j)
                elif grid[i][j] == "C":
                    cat = (i, j)
                elif grid[i][j] == "M":
                    mouse = (i, j)
                elif grid[i][j] == "#":
                    walls.add((i, j))

        @lru_cache(None)
        def dfs(cat: Tuple[int, int], mouse: Tuple[int, int], turn: int) -> bool:
            """老鼠是否赢"""
            if (
                cat == food or cat == mouse or turn >= 150
            ):  # 2*n 不足以判断平局 2*n*n 才可以 题目很贴心调整了规则为 1000 步以内为猫获胜 但是1000会TLE
                return False
            if mouse == food:
                return True

            # 老鼠动
            if not turn & 1:
                r, c = mouse
                for dr, dc in DIR4:
                    for jump in range(mouseJump + 1):
                        nr, nc = r + jump * dr, c + jump * dc
                        if not (0 <= nr < ROW and 0 <= nc < COL) or ((nr, nc) in walls):
                            break
                        if dfs(cat, (nr, nc), turn + 1):
                            return True
                return False

            # 猫动
            else:
                r, c = cat
                for dr, dc in DIR4:
                    for jump in range(catJump + 1):
                        nr, nc = r + jump * dr, c + jump * dc
                        if not (0 <= nr < ROW and 0 <= nc < COL) or ((nr, nc) in walls):
                            break
                        if not dfs((nr, nc), mouse, turn + 1):
                            return False
                return True

        res = dfs(cat, mouse, 0)
        dfs.cache_clear()
        return res


print(Solution().canMouseWin(grid=["####F", "#C...", "M...."], catJump=1, mouseJump=2))
# 输出：true
# 解释：猫无法抓到老鼠，也没法比老鼠先到达食物。
