import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 頂点に
# 1 から
# N までの、辺に
# 1 から
# M までの番号がついた
# N 頂点
# M 辺の単純無向グラフがあります。 辺
# i は頂点
# u
# i
# ​
#   と頂点
# v
# i
# ​
#   を結んでいます。
# また、全ての頂点は赤か青のいずれか一方で塗られています。頂点
# i の色は
# C
# i
# ​
#   で表されて、
# C
# i
# ​
#   が
# 0 ならば頂点
# i は赤く、
# 1 ならば頂点
# i は青く塗られています。

# 今、高橋君が頂点
# 1 に、青木君が頂点
# N にいます。
# 2 人は次の行動を
# 0 回以上好きな回数繰り返します。

# 2 人が同時に、今いる頂点に隣接している頂点のいずれか 1 個に移動する。
# ただし、高橋君の移動先の頂点の色と、青木君の移動先の頂点の色は異なる必要がある。
# 上記の行動を繰り返すことで、高橋君が頂点
# N に、青木君が頂点
# 1 にいる状態にできますか？
# 可能である場合は必要な行動回数の最小値を答えてください。不可能である場合は -1 を出力してください。

# 入力のはじめに
# T が与えられるので、
# T 個のテストケースについて問題を解いてください。

from collections import deque


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        n, m = map(int, input().split())
        C = list(map(int, input().split()))
        adjList = [[] for _ in range(n)]
        for _ in range(m):
            u, v = map(int, input().split())
            adjList[u - 1].append(v - 1)
            adjList[v - 1].append(u - 1)

        p1, p2 = 0, n - 1
        # 每个时候的状态是O(n^2) 记录两个人的位置
        queue = deque([(p1, p2, 0)])
        visited = set([(p1, p2)])
        ok = False
        while queue:
            curX, curY, curDist = queue.popleft()
            if curX == n - 1 and curY == 0:
                print(curDist)
                ok = True
                break
            for nextX in adjList[curX]:
                for nextY in adjList[curY]:
                    if C[nextX] != C[nextY] and (nextX, nextY) not in visited:
                        queue.append((nextX, nextY, curDist + 1))
                        visited.add((nextX, nextY))
        if not ok:
            print(-1)
