from typing import List, Tuple


# 我们`考虑`旋转、翻转操作。

# 1. 怎么比较相同的岛屿
# 大致思路为：
# 首先找出每个岛屿，随后将它进行旋转和翻转操作得到各种不同的情况，
# 接着对这些情况分别进行归一化处理，然后对归一化结果进行哈希得到具有唯一性的签名，
# 最终所有岛屿经所有变换并归一化且哈希后得到的不同签名数量就是满足要求的不同岛屿数量。

# 即：
# 计算出岛屿上每个点的局部坐标值。
# 局部坐标值的计算方法是，对于坐标的每一维，求出岛屿上所有点中的最小值，
# 并把所有点的这一维减去这个最小值。形象地来说，就是用一个最小的矩形框住岛屿，
# 矩形的左上角坐标为 (0, 0)，岛屿上每个点的局部坐标值就是相对于矩形左上角的坐标值。

# 岛屿哈希值
class Solution:
    def _dfs_flood(self, grid: List[List[int]], island: List[Tuple[int, int]], i: int, j: int):
        m, n = len(grid), len(grid[0])
        if i < 0 or i >= m or j < 0 or j >= n:
            return
        if grid[i][j] == 0:
            return
        island.append((i, j))
        grid[i][j] = 0
        self._dfs_flood(grid, island, i - 1, j)
        self._dfs_flood(grid, island, i + 1, j)
        self._dfs_flood(grid, island, i, j - 1)
        self._dfs_flood(grid, island, i, j + 1)

    # 八种情况
    def _transform(
        self, island: List[Tuple[int, int]], transformation: List[Tuple[int, int]]
    ) -> List[Tuple[int, int]]:
        transformed_island = []
        for row, col in island:
            transformed_island.append(
                (
                    transformation[0][0] * row + transformation[1][0] * col,
                    transformation[0][1] * row + transformation[1][1] * col,
                )
            )
        return transformed_island

    def _normalize(self, island: List[Tuple[int, int]]) -> List[Tuple[int, int]]:
        normalized_island = []
        minRow = min([row for row, _ in island])
        minCol = min([col for _, col in island])
        # 基准点一样
        for row, col in island:
            normalized_island.append((row - minRow, col - minCol))
        return normalized_island

    def numDistinctIslands2(self, grid: List[List[int]]) -> int:
        m, n, diff = len(grid), len(grid[0]), 0
        transformations = [
            [(1, 0), (0, 1)],
            [(0, -1), (1, 0)],
            [(-1, 0), (0, -1)],
            [(0, 1), (-1, 0)],
            [(1, 0), (0, -1)],
            [(-1, 0), (0, 1)],
            [(0, 1), (1, 0)],
            [(0, -1), (-1, 0)],
        ]

        distinct_transformations = set()

        for i in range(m):
            for j in range(n):
                if grid[i][j] == 1:
                    island = []
                    self._dfs_flood(grid, island, i, j)
                    size_before = len(distinct_transformations)

                    # 每次都加八个
                    for transformation in transformations:
                        distinct_transformations.add(
                            tuple(sorted(self._normalize(self._transform(island, transformation))))
                        )

                    size_after = len(distinct_transformations)
                    if size_after > size_before:
                        diff += 1
        return diff


# 11000
# 11000
# 00011
# 00011
# 给定上图，返回结果 1 。
print(
    Solution().numDistinctIslands2(
        [[1, 1, 0, 0, 0], [1, 0, 0, 0, 0], [0, 0, 0, 0, 1], [0, 0, 0, 1, 1]]
    )
)
