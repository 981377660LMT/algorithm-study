from typing import List


DIR4 = ((0, 1), (1, 0), (0, -1), (-1, 0))


class Solution:
    def spiralOrder(self, matrix: List[List[int]]) -> List[int]:
        """请按照 顺时针螺旋顺序 ，返回矩阵中的所有元素。"""
        ROW, COL = len(matrix), len(matrix[0])
        r, c, di = 0, 0, 0
        res, visited = [], set()

        while len(res) < ROW * COL:
            res.append(matrix[r][c])
            visited.add((r, c))
            nr, nc = r + DIR4[di][0], c + DIR4[di][1]
            if nr < 0 or nr >= ROW or nc < 0 or nc >= COL or (nr, nc) in visited:
                di = (di + 1) % 4
                nr, nc = r + DIR4[di][0], c + DIR4[di][1]
            r, c = nr, nc

        return res
