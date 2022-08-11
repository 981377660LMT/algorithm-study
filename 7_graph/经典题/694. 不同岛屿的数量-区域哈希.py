from typing import List, Set

# 694. 不同岛屿的数量-区域哈希
# 我们不考虑旋转、翻转操作。

# 1. 怎么比较相同的dfs
# 每次dfs都从左上角开始 最后生成一个dfs的哈希(用每次的转向表示) 如果岛屿一样则dfs必定一样

# 注意此题必须visited数组 不能修改矩阵值来visited  为什么?

# 岛屿哈希值
class Solution:
    def numDistinctIslands(self, grid: List[List[int]]) -> int:
        def getHashOfIslands(grid: List[List[int]]) -> Set[str]:
            """获取矩阵中各个岛屿的哈希值"""

            def dfs(r: int, c: int, path: List[str]) -> None:
                if visited[r][c]:
                    return

                visited[r][c] = True
                dirs = [(r, c + 1), (r + 1, c), (r, c - 1), (r - 1, c)]
                for i in range(4):
                    nr, nc = dirs[i]
                    if 0 <= nr < row and 0 <= nc < col and grid[nr][nc] == 1:
                        path.append(str(i + 1))
                        dfs(nr, nc, path)

            row, col = len(grid), len(grid[0])
            visited = [[False for _ in range(col)] for _ in range(row)]

            res = set()
            for r in range(row):
                for c in range(col):
                    if grid[r][c] == 1 and not visited[r][c]:
                        path = []
                        dfs(r, c, path)
                        # path.extend(['_', str(r), '_', str(c)])
                        res.add(''.join(path))

            return res

        return len(getHashOfIslands(grid))


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

