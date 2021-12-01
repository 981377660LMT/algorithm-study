# 3 x 3 的幻方是一个填充有从 1 到 9 的不同数字的 3 x 3 矩阵，其中每行，每列以及两条对角线上的各数之和都相等。
from typing import List

# 给定一个由整数组成的 grid，其中有多少个 3 × 3 的 “幻方” 子矩阵？
class Solution(object):
    def numMagicSquaresInside(self, grid):
        R, C = len(grid), len(grid[0])

        def magic(a, b, c, d, e, f, g, h, i):
            return sorted([a, b, c, d, e, f, g, h, i]) == range(1, 10) and (
                a + b + c
                == d + e + f
                == g + h + i
                == a + d + g
                == b + e + h
                == c + f + i
                == a + e + i
                == c + e + g
                == 15
            )

        ans = 0
        for r in range(R - 2):
            for c in range(C - 2):
                if grid[r + 1][c + 1] != 5:
                    continue  # optional skip
                if magic(
                    grid[r][c],
                    grid[r][c + 1],
                    grid[r][c + 2],
                    grid[r + 1][c],
                    grid[r + 1][c + 1],
                    grid[r + 1][c + 2],
                    grid[r + 2][c],
                    grid[r + 2][c + 1],
                    grid[r + 2][c + 2],
                ):
                    ans += 1
        return ans


print(Solution().numMagicSquaresInside([[4, 3, 8, 4], [9, 5, 1, 9], [2, 7, 6, 2]]))
