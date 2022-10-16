# 二维前缀最小值/最大值
from typing import List

INF = int(1e18)


class PreminMatrix:
    def __init__(self, matrix: List[List[int]]):
        ROW, COL = len(matrix), len(matrix[0])
        preMin = [row[:] for row in matrix]
        for r in range(ROW):
            for c in range(COL):
                if r - 1 >= 0:
                    cand = preMin[r - 1][c]
                    preMin[r][c] = cand if cand < preMin[r][c] else preMin[r][c]
                if c - 1 >= 0:
                    cand = preMin[r][c - 1]
                    preMin[r][c] = cand if cand < preMin[r][c] else preMin[r][c]
        self.preMin = preMin

    def query(self, x: int, y: int):
        """查询[0:x+1, 0:y+1]的最小值"""
        if x < 0 or y < 0:
            return INF
        if x >= len(self.preMin):
            x = len(self.preMin) - 1
        if y >= len(self.preMin[0]):
            y = len(self.preMin[0]) - 1
        return self.preMin[x][y]


class PremaxMatrix:
    def __init__(self, matrix: List[List[int]]):
        ROW, COL = len(matrix), len(matrix[0])
        preMax = [row[:] for row in matrix]
        for r in range(ROW):
            for c in range(COL):
                if r - 1 >= 0:
                    cand = preMax[r - 1][c]
                    preMax[r][c] = cand if cand > preMax[r][c] else preMax[r][c]
                if c - 1 >= 0:
                    cand = preMax[r][c - 1]
                    preMax[r][c] = cand if cand > preMax[r][c] else preMax[r][c]
        self.preMax = preMax

    def query(self, x: int, y: int):
        """查询[0:x+1, 0:y+1]的最大值"""
        if x < 0 or y < 0:
            return -INF
        if x >= len(self.preMax):
            x = len(self.preMax) - 1
        if y >= len(self.preMax[0]):
            y = len(self.preMax[0]) - 1
        return self.preMax[x][y]


if __name__ == "__main__":
    preMinMatrix = PreminMatrix([[1, 2, 3], [4, 5, 6], [7, 8, 9]])
    print(preMinMatrix.query(0, 0))  # 1

    preMaxMatrix = PremaxMatrix([[1, 2, 3], [4, 5, 6], [7, 8, 9]])
    print(preMaxMatrix.query(0, 0))  # 1
    print(preMaxMatrix.query(1, 1))  # 5
