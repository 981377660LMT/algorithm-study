from typing import List


class Solution:
    def diagonalSum(self, mat: List[List[int]]) -> int:
        n = len(mat)
        res = 0
        for row in range(n):
            res += mat[row][row]
            res += mat[row][n - row - 1]

        return res - ((mat[n >> 1][n >> 1]) if (n & 1) else 0)


print(Solution().diagonalSum([[1, 1, 1, 1], [1, 1, 1, 1], [1, 1, 1, 1], [1, 1, 1, 1]]))

