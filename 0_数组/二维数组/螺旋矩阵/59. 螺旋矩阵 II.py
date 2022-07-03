# 生成包含 1-n^2 所有元素的螺旋矩阵

from typing import List


DIR4 = ((0, 1), (1, 0), (0, -1), (-1, 0))


class Solution:
    def generateMatrix(self, n: int) -> List[List[int]]:
        res = [[0] * n for _ in range(n)]
        r, c, di = 0, 0, 0
        visited = set()
        cur = 1

        while cur <= n * n:
            res[r][c] = cur
            visited.add((r, c))
            cur += 1
            nr, nc = r + DIR4[di][0], c + DIR4[di][1]
            if nr < 0 or nr >= n or nc < 0 or nc >= n or (nr, nc) in visited:
                di = (di + 1) % 4
                nr, nc = r + DIR4[di][0], c + DIR4[di][1]
            r, c = nr, nc

        return res
