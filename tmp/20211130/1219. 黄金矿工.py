from typing import List

# 每个单元格中的整数就表示这一单元格中的黄金数量
# 每个单元格中的整数就表示这一单元格中的黄金数量
# 矿工可以从网格中 任意一个 有黄金的单元格出发或者是停止。
# 求最大收益
# 1 <= grid.length, grid[i].length <= 15
# 最多 25 个单元格中有黄金。

# O(3^k)
# 总结:
# 想象二叉树，需要dfs后序遍历
class Solution:
    def getMaximumGold(self, grid: List[List[int]]) -> int:
        m, n = len(grid), len(grid[0])

        def dfs(r: int, c: int) -> int:
            if not (0 <= r < m and 0 <= c < n) or grid[r][c] == 0:
                return 0
            tmp = grid[r][c]
            grid[r][c] = 0

            # dfs后序即可
            next_sum = -0x3FFFFFFF
            for nr, nc in ((r - 1, c), (r + 1, c), (r, c - 1), (r, c + 1)):
                next_sum = max(next_sum, dfs(nr, nc))

            grid[r][c] = tmp
            return tmp + next_sum

        res = 0
        for r in range(m):
            for c in range(n):
                if grid[r][c] != 0:
                    res = max(res, dfs(r, c))
        return res


print(Solution().getMaximumGold(grid=[[0, 6, 0], [5, 8, 7], [0, 9, 0]]))
