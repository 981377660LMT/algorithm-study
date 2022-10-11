from typing import List

# 每个单元格中的整数就表示这一单元格中的黄金数量
# 每个单元格中的整数就表示这一单元格中的黄金数量
# 矿工可以从网格中 任意一个 有黄金的单元格出发或者是停止。
# !不得开采（进入）黄金数目为 0 的单元格。
# 求最大收益
# 1 <= grid.length, grid[i].length <= 15
# !最多 25 个单元格中有黄金。

# 二维回溯


class Solution:
    def getMaximumGold(self, grid: List[List[int]]) -> int:
        def dfs(r: int, c: int, gold: int) -> None:
            nonlocal res
            gold += grid[r][c]
            res = gold if gold > res else res

            tmp = grid[r][c]
            grid[r][c] = 0  # 当前结点标记访问

            for nextR, nextC in ((r - 1, c), (r + 1, c), (r, c - 1), (r, c + 1)):
                if 0 <= nextR < ROW and 0 <= nextC < COL and grid[nextR][nextC] > 0:
                    dfs(nextR, nextC, gold + grid[r][c])

            grid[r][c] = tmp

        ROW, COL = len(grid), len(grid[0])
        res = 0
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] != 0:
                    dfs(r, c, 0)
        return res


print(Solution().getMaximumGold(grid=[[0, 6, 0], [5, 8, 7], [0, 9, 0]]))
