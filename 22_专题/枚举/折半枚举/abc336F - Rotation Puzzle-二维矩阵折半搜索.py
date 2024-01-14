# https://atcoder.jp/contests/abc336/tasks/abc336_f
# 折半搜索/双向bfs
# 给定一个n*m(3<=n,m<=8)的网格,0到n*m-1每个数字恰好出现一次.
# !现在可以以(0,0),(1,0),(0,1)或(1,1)为左上角180度顺时针旋转(n-1)*(m-1)大小的矩阵.
# 最多操作`20`次,问能否使得网格变成有序的,即0-n*m-1从左到右从上到下依次排列.
# 如果能,输出最少操作次数,否则输出-1.
#
# 最多20次 -> 保证操作次数有限
# !状态定义:grid元组化后的状态

from collections import defaultdict
from random import randint
from typing import DefaultDict, List

INF = int(1e20)

CENTER_4 = [(0, 0), (0, 1), (1, 0), (1, 1)]


def rotationPuzzle(grid: List[List[int]]) -> int:
    ROW, COL = len(grid), len(grid[0])
    target = [[i * COL + j for j in range(COL)] for i in range(ROW)]
    hashBase = [[0] * COL for _ in range(ROW)]
    for i in range(ROW):
        for j in range(COL):
            hashBase[i][j] = randint(1, 1 << 61)

    def rotate(mat: List[List[int]], leftX: int, leftY: int) -> List[List[int]]:
        """以(leftX, leftY)为左上角顺时针旋转180度."""
        res = [list(row) for row in mat]
        for i in range(ROW - 1):
            for j in range(COL - 1):
                res[leftX + i][leftY + j] = mat[leftX + ROW - 2 - i][leftY + COL - 2 - j]
        return res

    def getStates(mat: List[List[int]]) -> DefaultDict[int, int]:
        def dfs(curMat: List[List[int]], step: int) -> None:
            matHash = sum(curMat[i][j] * hashBase[i][j] for i in range(ROW) for j in range(COL))
            if matHash in res and res[matHash] <= step:
                return
            res[matHash] = step
            if step == 10:
                return
            for x, y in CENTER_4:
                dfs(rotate(curMat, x, y), step + 1)

        res = defaultdict(int)
        dfs(mat, 0)
        return res

    mp1 = getStates(grid)
    mp2 = getStates(target)
    res = INF
    for k, v in mp1.items():
        if k in mp2:
            res = min(res, v + mp2[k])
    return res if res != INF else -1


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline
    H, W = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(H)]
    for i in range(H):
        for j in range(W):
            grid[i][j] -= 1
    print(rotationPuzzle(grid))
