from typing import List

# 959. 由斜杠划分区域.ts
# 1 表示连接左单元格和右单元格的街道。
# 2 表示连接上单元格和下单元格的街道。
# 3 表示连接左单元格和下单元格的街道。
# 4 表示连接右单元格和下单元格的街道。
# 5 表示连接左单元格和上单元格的街道。
# 6 表示连接右单元格和上单元格的街道。
# https://leetcode.com/problems/check-if-there-is-a-valid-path-in-a-grid/discuss/547229/Python-Union-Find


# 两个点可以互通:
# 下一个点可以返回上一个点
# (-di, -dj) in directions[grid[ni][nj]]  # 下一个点可以原路返回

MAPPING = {
    1: [(0, -1), (0, 1)],
    2: [(-1, 0), (1, 0)],
    3: [(0, -1), (1, 0)],
    4: [(0, 1), (1, 0)],
    5: [(0, -1), (-1, 0)],
    6: [(0, 1), (-1, 0)],
}


class Solution:
    def hasValidPath(self, grid: List[List[int]]) -> bool:
        def dfs(r: int, c: int) -> bool:
            if (r, c) == target:
                return True
            for dr, dc in MAPPING[grid[r][c]]:
                nr, nc = r + dr, c + dc
                if (
                    0 <= nr < ROW
                    and 0 <= nc < COL
                    and (nr, nc) not in visited
                    and (-dr, -dc) in MAPPING[grid[nr][nc]]  # 下一个点可以原路返回
                ):
                    visited.add((nr, nc))
                    if dfs(nr, nc):
                        return True
            return False

        ROW, COL = len(grid), len(grid[0])
        target = (ROW - 1, COL - 1)
        visited = set([(0, 0)])
        return dfs(0, 0)


print(Solution().hasValidPath([[2, 4, 3], [6, 5, 2]]))
# 输出：true
# 解释：如图所示，你可以从 (0, 0) 开始，访问网格中的所有单元格并到达 (m - 1, n - 1) 。
