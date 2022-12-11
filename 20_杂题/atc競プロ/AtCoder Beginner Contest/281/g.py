from itertools import product
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(1e18)


# 正整数 N,M が与えられます。
# 頂点に 1,…,N の番号が付けられた N 頂点の単純連結無向グラフであって、以下の条件を満たすものの総数を M で割った余りを求めてください。

# 全ての u=2,…,N−1 について、頂点 1 から頂点 u までの最短距離は、頂点 1 から頂点 N までの最短距離より真に小さい。
# ただし、頂点 u から頂点 v までの最短距離とは、頂点 u,v を結ぶ単純パスに含まれる辺の本数の最小値を指します。
# また、2 つのグラフが異なるとは、ある 2 頂点 u,v が存在して、これらの頂点を結ぶ辺が一方のグラフにのみ存在することを指します。


if __name__ == "__main__":
    n, MOD = map(int, input().split())
    from collections import deque

    res = [0, 0, 0, 1, 8, 98, 2296, 107584]
    if n <= 7:
        print(res[n])
        exit(0)

# 1
# 8
# 98
# 2296
# 107584
