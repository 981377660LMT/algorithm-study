# N 個の頂点と M 本の辺からなる連結かつ単純な無向グラフが与えられます。
# i=1,2,…,M について、i 番目の辺は頂点 u
# i
# ​
#   と頂点 v
# i
# ​
#   を結んでいます。

# 高橋君は、はじめレベルが 0 の状態で頂点 1 におり、下記の行動をちょうど K 回行います。

# まず、いまいる頂点に隣接する頂点の中から、1 つを等確率でランダムに選択し、その頂点に移動する。
# その後、移動後の頂点 v に応じて、下記のイベントが発生します。
# C
# v
# ​
#  =0 のとき : 高橋君のレベルが 1 だけ増加する。
# C
# v
# ​
#  =1 のとき : 高橋君のいまのレベルを X とする。高橋君は X
# 2
#   円のお金を獲得する。
# 上記の K 回の行動の過程で高橋君が獲得するお金の合計金額の期待値を mod998244353 で出力してください（注記参照）。

from functools import lru_cache
import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# !n<=3000 k<=3000 稀疏图
if __name__ == "__main__":
    n, m, k = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)
    values = list(map(int, input().split()))

    @lru_cache(None)
    def dfs(cur: int, pre: int, level: int, k: int) -> Tuple[int, int]:
        if k == 0:
            return 0, 1
        res = 0
        allCount = 0
        for next in adjList[cur]:
            if next == pre:
                continue
            if values[next] == 0:
                a, b = dfs(next, cur, level + 1, k - 1)
                res = (res + a * b) % MOD
                allCount = (allCount + b) % MOD
            else:
                a, b = dfs(next, cur, level, k - 1)
                res = (res + a * b + level * level) % MOD
                allCount = (allCount + b) % MOD

        print("asa", res, allCount)
        return res * pow(allCount, MOD - 2, MOD) % MOD, allCount

    res = dfs(0, -1, 0, k)
    print(res[0])
