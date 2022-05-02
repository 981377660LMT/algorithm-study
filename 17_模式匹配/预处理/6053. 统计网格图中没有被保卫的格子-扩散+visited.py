from typing import List

# 请你返回空格子中，有多少个格子是 没被保卫 的。
# 前缀后缀统计(蜡烛间的盘子)
# 这种方法很差，可行但是代码量很大


# 更好的方法是四个方向扩散+visited去重

DIRS = [(0, 1), (0, -1), (1, 0), (-1, 0)]


class Solution:
    def countUnguarded(
        self, row: int, col: int, guards: List[List[int]], walls: List[List[int]]
    ) -> int:
        """从每个点向四个方向扩散+visited去重,染色
        
        从守卫扩散，碰到墙/守卫就停止
        这个优化的关键是要想清楚每行的每个空格最多会被扫到两次，因此时间复杂度是O(m*n)的
        """
        matrix = [[0] * col for _ in range(row)]
        for r, c in guards:
            matrix[r][c] = 1
        for r, c in walls:
            matrix[r][c] = 2

        for r in range(row):
            for c in range(col):
                if matrix[r][c] != 0:
                    continue
                for dr, dc in DIRS:
                    nr, nc = r + dr, c + dc
                    while 0 <= nr < row and 0 <= nc < col and matrix[nr][nc] not in (1, 2):
                        matrix[nr][nc] = 3
                        nr, nc = nr + dr, nc + dc

        return sum(1 for r in range(row) for c in range(col) if matrix[r][c] == 0)


print(
    Solution().countUnguarded(
        row=4, col=6, guards=[[0, 0], [1, 1], [2, 3]], walls=[[0, 1], [2, 2], [1, 4]]
    )
)
# print(Solution().countUnguarded(m=3, n=3, guards=[[1, 1]], walls=[[0, 1], [1, 0], [2, 1], [1, 2]]))

