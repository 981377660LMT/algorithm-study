from collections import deque
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 縦
# N 行 横
# N 列のマスからなるグリッドがあります。ここで、
# N は
# 45 以下の奇数です。

# 上から
# i 行目、左から
# j 列目のマスをマス
# (i,j) と表します。

# このグリッドに、以下の条件を満たすように高橋君と
# 1 から
# N
# 2
#  −1 までの番号がついた
# N
# 2
#  −1 個のパーツからなる
# 1 匹の龍を配置します。

# 高橋君はグリッドの中央、すなわちマス
# (
# 2
# N+1
# ​
#  ,
# 2
# N+1
# ​
#  ) に配置しなければならない。
# 高橋君がいるマスを除く各マスには龍のパーツをちょうど
# 1 つ配置しなければならない。
# 2≤x≤N
# 2
#  −1 を満たす全ての整数
# x について、龍のパーツ
# x はパーツ
# x−1 があるマスと辺で隣接するマスに配置しなければならない。
# マス
# (i,j) とマス
# (k,l) が辺で隣接するとは、
# ∣i−k∣+∣j−l∣=1 であることを意味します。
# 条件を満たす配置方法を
# 1 つ出力してください。なお、条件を満たす配置は必ず存在します。

from typing import List


DIR4 = ((0, 1), (1, 0), (0, -1), (-1, 0))


def generateMatrix(n: int) -> List[List[int]]:
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


if __name__ == "__main__":
    n = int(input())
    grid = generateMatrix(n)
    grid[n // 2][n // 2] = "T"
    for i in range(n):
        print(" ".join(map(str, grid[i])))
