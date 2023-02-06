# 矩阵转无向图邻接表
# 网格转无向图邻接表
# matrixToAdjMap
from collections import defaultdict
from typing import DefaultDict, List, Set


def matrixToAdjMap(matrix: List[List[int]]) -> DefaultDict[int, Set[int]]:
    """二维矩阵转无向图邻接表"""
    row, col = len(matrix), len(matrix[0])
    adjMap = defaultdict(set)
    for r in range(row):
        for c in range(col):
            cur = r * col + c
            if r + 1 < row:
                next = (r + 1) * col + c
                adjMap[cur].add(next)
                adjMap[next].add(cur)
            if c + 1 < col:
                next = r * col + c + 1
                adjMap[cur].add(next)
                adjMap[next].add(cur)
    return adjMap


if __name__ == "__main__":
    matrix = [[0, 0, 0], [0, 0, 0], [0, 0, 0]]
    print(matrixToAdjMap(matrix))
