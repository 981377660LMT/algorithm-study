# 铁道建设/公路建设
# 要选择两个点建设一条公路
# 总花费为:每公里的花费(曼哈顿距离)+每个城市所在地的固定花费
# !求grid[i1][j1]+grid[i2][j2]+c*(abs(i1-i2)+abs(j1-j2))的最小值
# ROW,COL<=1000

# 19_数学/计算几何/坐标轴旋转/曼哈顿与切比雪夫距离/1131. 绝对值表达式的最大值 .py

# !曼哈顿距离 去绝对值+公式变形
# !由对称性 假设i1<=i2,j1<=j2 (可以上下翻转变成这样的情形) 从而去绝对值
# 因此只需维护二维前缀的 `grid[i1][j1]-c*(i1+j1)` 的最小值
# 每次加上grid[i2][j2]+c*(i2+j2)即可
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class PreMinMatrix:
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

    def query(self, x1: int, y1: int):
        """查询[0:x1+1, 0:y1+1]的最小值"""
        if x1 < 0 or y1 < 0:
            return INF
        if x1 >= len(self.preMin):
            x1 = len(self.preMin) - 1
        if y1 >= len(self.preMin[0]):
            y1 = len(self.preMin[0]) - 1
        return self.preMin[x1][y1]


def flipTopToBottom(matrix: List[List[int]]) -> List[List[int]]:
    """矩阵上下翻转"""
    newMatrix = [row[:] for row in matrix]
    ROW = len(matrix)
    for i in range(ROW // 2):
        newMatrix[i], newMatrix[~i] = newMatrix[~i], newMatrix[i]
    return newMatrix


if __name__ == "__main__":

    def helper(grid: List[List[int]], cost: int) -> int:
        """
        寻找两个点(i1,j1)与(i2,j2)使得建设花费最小
        i1<=i2,j1<=j2

        用前缀min求出
        """
        ROW, COL = len(grid), len(grid[0])
        newGrid = [[grid[i][j] - cost * (i + j) for j in range(COL)] for i in range(ROW)]
        preMin = PreMinMatrix(newGrid)
        res = INF
        for r in range(ROW):
            for c in range(COL):
                min_ = min(preMin.query(r - 1, c), preMin.query(r, c - 1))
                cand = min_ + grid[r][c] + cost * (r + c)
                res = cand if cand < res else res
        return res

    ROW, COL, cost = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(ROW)]
    print(min(helper(grid, cost), helper(flipTopToBottom(grid), cost)))
