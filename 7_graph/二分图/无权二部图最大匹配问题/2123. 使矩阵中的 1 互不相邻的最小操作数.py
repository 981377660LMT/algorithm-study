#  相邻两个1组成一条边，每条边都要去掉一个端点，
#  !其实是找最小点覆盖，即求二分图的最大匹配，跑匈牙利算法
#  2123. 使矩阵中的 1 互不相邻的最小操作数

from typing import List
from 匈牙利算法 import Hungarian

DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]


class Solution:
    def minimumOperations(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        H = Hungarian(ROW * COL, ROW * COL)
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == 0 or (r + c) & 1:
                    continue
                cur = r * COL + c
                for dr, dc in DIR4:
                    nr, nc = r + dr, c + dc
                    if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == 1:
                        H.addEdge(cur, nr * COL + nc)
        return len(H.work())


print(Solution().minimumOperations(grid=[[1, 1, 0], [0, 1, 1], [1, 1, 1]]))
