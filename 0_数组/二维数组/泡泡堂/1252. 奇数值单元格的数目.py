# 对 indices[i] 所指向的每个位置，应同时执行下述增量操作：

# ri 行上的所有单元格，加 1 。
# ci 列上的所有单元格，加 1 。
from typing import List

# 矩阵中位于 (x, y) 位置的数为奇数，当且仅当 rows[x] 和 cols[y] 中恰好有一个为奇数。
# 时间复杂度为 O(n + m + indices.length) 且仅用 O(n + m) 额外空间的算法来解决此问题吗？
class Solution:
    def oddCells(self, n: int, m: int, indices: List[List[int]]) -> int:
        rows = [0] * n
        cols = [0] * m
        for x, y in indices:
            rows[x] += 1
            cols[y] += 1

        odd_rows = sum(x & 1 for x in rows)
        odd_cols = sum(y & 1 for y in cols)
        return odd_rows * (m - odd_cols) + odd_cols * (n - odd_rows)

