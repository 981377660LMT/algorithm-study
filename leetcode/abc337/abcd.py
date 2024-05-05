from collections import deque
from itertools import groupby
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# H 行
# W 列のグリッドがあります。上から
# i 行目、左から
# j 行目のマスをマス
# (i,j) と呼びます。

# 各マスには o 、x 、. のうちいずれかの文字が書かれています。 各マスに書かれた文字は
# H 個の長さ
# W の文字列
# S
# 1
# ​
#  ,S
# 2
# ​
#  ,…,S
# H
# ​
#   で表され、 マス
# (i,j) に書かれた文字は、文字列
# S
# i
# ​
#   の
# j 文字目と一致します。

# このグリッドに対して、下記の操作を
# 0 回以上好きな回数だけ繰り返します。

# . が書かれているマスを
# 1 個選び、そのマスに書かれた文字を o に変更する。
# その結果、縦方向または横方向に連続した
# K 個のマスであってそれらに書かれた文字がすべて o であるようなものが存在する（ すなわち、下記の
# 2 つの条件のうち少なくとも一方を満たす）ようにすることが可能かを判定し、可能な場合はそのために行う操作回数の最小値を出力してください。

# 1≤i≤H かつ
# 1≤j≤W−K+1 を満たす整数の組
# (i,j) であって、マス
# (i,j),(i,j+1),…,(i,j+K−1) に書かれた文字が o であるものが存在する。
# 1≤i≤H−K+1 かつ
# 1≤j≤W を満たす整数の組
# (i,j) であって、マス
# (i,j),(i+1,j),…,(i+K−1,j) に書かれた文字が o であるものが存在する。


class PreSumMatrix:
    """二维前缀和模板(矩阵不可变)"""

    def __init__(self, A: List[List[int]]):
        m, n = len(A), len(A[0])

        # 前缀和数组
        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                preSum[r + 1][c + 1] = A[r][c] + preSum[r][c + 1] + preSum[r + 1][c] - preSum[r][c]
        self.preSum = preSum

    def sumRegion(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值::

        preSumMatrix.sumRegion(0, 0, 2, 2) # 左上角(0, 0)到右下角(2, 2)的值
        """
        return (
            self.preSum[r2 + 1][c2 + 1]
            - self.preSum[r2 + 1][c1]
            - self.preSum[r1][c2 + 1]
            + self.preSum[r1][c1]
        )


# !枚举左上角
if __name__ == "__main__":
    H, W, K = map(int, input().split())
    grid = [list(input()) for _ in range(H)]

    goodMatrix = [[0] * W for _ in range(H)]
    badMatrix = [[0] * W for _ in range(H)]
    for r in range(H):
        for c in range(W):
            goodMatrix[r][c] = grid[r][c] == "o"
            badMatrix[r][c] = grid[r][c] == "x"

    goodSum = PreSumMatrix(goodMatrix)
    badSum = PreSumMatrix(badMatrix)

    res = INF
    for c in range(W - K + 1):
        for r in range(H):
            row1, col1 = r, c + K - 1
            if badSum.sumRegion(r, c, row1, col1) > 0:
                continue
            res = min(res, K - goodSum.sumRegion(r, c, row1, col1))
    for r in range(H - K + 1):
        for c in range(W):
            row2, col2 = r + K - 1, c
            if badSum.sumRegion(r, c, row2, col2) > 0:
                continue
            res = min(res, K - goodSum.sumRegion(r, c, row2, col2))
    if res == INF:
        print(-1)
    else:
        print(res)
