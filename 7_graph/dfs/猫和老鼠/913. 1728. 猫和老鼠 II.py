from typing import List
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


class Solution:
    def canMouseWin(self, grid: List[str], catJump: int, mouseJump: int) -> bool:
        m, n = len(grid), len(grid[0])  # dimensions
        walls = set()
        available = m * n  # available steps for mouse and cat
        for i in range(m):
            for j in range(n):
                if grid[i][j] == "F":
                    food = (i, j)
                elif grid[i][j] == "C":
                    cat = (i, j)
                elif grid[i][j] == "M":
                    mouse = (i, j)
                elif grid[i][j] == "#":
                    walls.add((i, j))
                    available -= 1

        @lru_cache(None)
        def fn(cat, mouse, turn):
            """Return True if mouse wins."""
            if cat == food or cat == mouse or turn >= available * 2:
                return False
            if mouse == food:
                return True  # mouse reaching food

            if not turn & 1:  # mouse moving
                x, y = mouse
                for dx, dy in (-1, 0), (0, 1), (1, 0), (0, -1):
                    for jump in range(0, mouseJump + 1):
                        xx, yy = x + jump * dx, y + jump * dy
                        if not (0 <= xx < m and 0 <= yy < n) or (xx, yy) in walls:
                            # Stop extending the jump since we cannot go further
                            break
                        if fn(cat, (xx, yy), turn + 1):
                            return True
                return False

            else:  # cat moving
                x, y = cat
                for dx, dy in (-1, 0), (0, 1), (1, 0), (0, -1):
                    for jump in range(0, catJump + 1):
                        xx, yy = x + jump * dx, y + jump * dy
                        if not (0 <= xx < m and 0 <= yy < n) or (xx, yy) in walls:
                            break
                        if not fn((xx, yy), mouse, turn + 1):
                            return False
                return True

        return fn(cat, mouse, 0)


print(Solution().canMouseWin(grid=["####F", "#C...", "M...."], catJump=1, mouseJump=2))
# 输出：true
# 解释：猫无法抓到老鼠，也没法比老鼠先到达食物。
