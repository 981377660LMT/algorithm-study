from typing import List

# 如果两台服务器位于同一行或者同一列，我们就认为它们之间可以进行通信。
# 请你统计并返回能够与至少一台其他服务器进行通信的服务器的数量。


class Solution:
    def countServers(self, grid: List[List[int]]) -> int:
        row, col = len(grid), len(grid[0])
        rowCounter = [0] * row
        colCounter = [0] * col
        for r in range(row):
            for c in range(col):
                if grid[r][c]:
                    rowCounter[r] += 1
                    colCounter[c] += 1

        res = 0
        for r in range(row):
            for c in range(col):
                res += grid[r][c] and (rowCounter[r] > 1 or colCounter[c] > 1)
        return res


print(Solution().countServers(grid=[[1, 1, 0, 0], [0, 0, 1, 0], [0, 0, 1, 0], [0, 0, 0, 1]]))
# 输出：4
# 解释：第一行的两台服务器互相通信，第三列的两台服务器互相通信，但右下角的服务器无法与其他服务器通信。
