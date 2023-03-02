import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 個の都市があり、都市
# 1, 都市
# 2,
# …, 都市
# N と番号づけられています。
# いくつかの異なる都市の間は一方通行のテレポーターによって移動できます。 都市
# i
# (1≤i≤N) からテレポーターによって直接移動できる都市は 0 と 1 からなる長さ
# M の文字列
# S
# i
# ​
#   によって表されます。具体的には、
# 1≤j≤N に対して、

# 1≤j−i≤M かつ
# S
# i
# ​
#   の
# (j−i) 文字目が 1 ならば、都市
# i から都市
# j に直接移動できる。
# そうでない時、都市
# i から都市
# j へは直接移動できない。
# k=2,3,…,N−1 に対して次の問題を解いてください。

# テレポータを繰り返し使用することによって、都市
# k を通らずに都市
# 1 から 都市
# N へ移動できるか判定し、 できるならばそのために必要なテレポーターの使用回数の最小値を、 できないならば
# −1 を出力せよ。
if __name__ == "__main__":
    n, m = map(int, input().split())
    words = [input() for _ in range(n)]
    adjList = [set() for _ in range(n)]
    rAdjList = [set() for _ in range(n)]
    for i, w in enumerate(words):
        for j, c in enumerate(w):
            if c == "1":
                adjList[i].add(i + j + 1)
                rAdjList[i + j + 1].add(i)

    from collections import deque

    def bfs(start: int, adjList: List[List[int]]) -> List[int]:
        n = len(adjList)
        dist = [INF] * n
        dist[start] = 0
        queue = deque([start])
        while queue:
            cur = queue.popleft()
            for next in adjList[cur]:
                cand = dist[cur] + 1
                if cand < dist[next]:
                    dist[next] = cand
                    queue.append(next)
        return dist

    # 枚举他左右范围即可
    dist1 = bfs(0, adjList)
    dist2 = bfs(n - 1, rAdjList)

    for i in range(1, n - 1):
        # 枚举左右落脚点(10个*10个)
        res = INF
        for left in range(max(0, i - m), i):
            for right in range(i + 1, min(n, i + m + 1)):
                if right in adjList[left]:
                    cand = dist1[left] + dist2[right] + 1
                    res = min(res, cand)
        print(-1 if res == INF else res)
