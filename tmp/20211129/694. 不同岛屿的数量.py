from typing import List

# 我们不考虑旋转、翻转操作。

# 1. 怎么比较相同的dfs
# 每次dfs都从左上角开始 最后生成一个dfs的哈希(用每次的转向表示) 如果岛屿一样则dfs必定一样

# 注意此题必须visited数组 不能修改矩阵值来visited  为什么?
class Solution:
    def numDistinctIslands(self, grid: List[List[int]]) -> int:
        Row, Col = len(grid), len(grid[0])
        visited = [[False for _ in range(Col)] for _ in range(Row)]

        def dfs(r: int, c: int) -> None:
            if visited[r][c]:
                return
            visited[r][c] = True
            dir = [(r, c + 1), (r + 1, c), (r, c - 1), (r - 1, c)]
            for i in range(4):
                nr, nc = dir[i]
                if 0 <= nr < Row and 0 <= nc < Col and grid[nr][nc] == 1:
                    grid[nr][nc] = -1
                    self.path += str(i + 1)
                    dfs(nr, nc)

        hash = set()
        for r in range(Row):
            for c in range(Col):
                if grid[r][c] == 1 and not visited[r][c]:
                    self.path = "0"
                    dfs(r, c)
                    hash.add(self.path)

        return len(hash)


# 11000
# 11000
# 00011
# 00011
# 给定上图，返回结果 1 。
print(
    Solution().numDistinctIslands(
        [
            [
                0,
                0,
                1,
                0,
                1,
                0,
                1,
                1,
                1,
                0,
                0,
                0,
                0,
                1,
                0,
                0,
                1,
                0,
                0,
                1,
                1,
                1,
                0,
                1,
                1,
                1,
                0,
                0,
                0,
                1,
                1,
                0,
                1,
                1,
                0,
                1,
                0,
                1,
                0,
                1,
                0,
                0,
                0,
                0,
                0,
                1,
                1,
                1,
                1,
                0,
            ],
            [
                0,
                0,
                1,
                0,
                0,
                1,
                1,
                1,
                0,
                0,
                1,
                0,
                1,
                0,
                0,
                1,
                1,
                0,
                0,
                1,
                0,
                0,
                0,
                1,
                0,
                1,
                1,
                1,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                1,
                1,
                1,
                0,
                0,
                0,
                1,
                0,
                1,
                1,
                0,
                1,
                0,
                0,
                0,
            ],
            [
                0,
                1,
                0,
                1,
                0,
                1,
                1,
                1,
                0,
                0,
                1,
                1,
                0,
                0,
                0,
                0,
                1,
                0,
                1,
                0,
                1,
                1,
                1,
                0,
                1,
                1,
                1,
                0,
                0,
                0,
                1,
                0,
                1,
                0,
                1,
                0,
                0,
                0,
                1,
                1,
                1,
                1,
                1,
                0,
                0,
                1,
                0,
                0,
                1,
                0,
            ],
            [
                1,
                0,
                1,
                0,
                0,
                1,
                0,
                1,
                0,
                0,
                1,
                0,
                0,
                1,
                1,
                1,
                0,
                1,
                0,
                0,
                0,
                0,
                1,
                0,
                1,
                0,
                0,
                1,
                0,
                1,
                1,
                1,
                0,
                1,
                0,
                0,
                0,
                1,
                1,
                1,
                0,
                0,
                0,
                0,
                1,
                1,
                1,
                1,
                1,
                1,
            ],
        ]
    )
)

