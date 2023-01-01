# Enumerate Cliques (clique:团)
# 求所有团乘积之和（团就是一个两两之间有边的顶点集合）
# n,m<=100
from typing import List


MOD = 998244353

# https://judge.yosupo.jp/submission/100863
def _solve_gu(graph, S):
    # 頂点集合Sに含まれるクリークを全探索。愚直
    ns = len(S)
    ans = 0
    for bit in range(1, 1 << ns):
        nowc = []
        for i in range(ns):
            if 1 & (bit >> i):
                nowc.append(S[i])
        is_clique = True

        for i in nowc:
            for j in nowc:
                if i == j:
                    continue
                if j not in graph[i]:
                    is_clique = False

        if is_clique:
            cnt = 1
            for v in nowc:
                cnt *= values[v]
                cnt %= MOD
            ans += cnt
            ans %= MOD

    return ans


def solve(graph: List[List[int]], S: List[int]) -> int:
    if len(S) <= 1:
        for i in S:
            return values[i]
        return 0

    # Sに含まれて次数が√(2m) 以下の頂点を1つとって、targetにする
    target = -1
    thr = round((2 * m) ** 0.5)
    for i in range(n):
        if len(graph[i]) < thr and i in S:
            target = i
    # targetが決まらなければ、頂点数が十分少ないので頂点で全探索できる
    if target == -1:
        return _solve_gu(graph, S)

    # targetが決まったら、targetを含むかどうかで場合分け
    # クリークの定義より、targetを含むクリークは必ずtargetとその隣接頂点の集合の部分集合となる。よって、これを全探索する
    ng = len(graph[target])
    res = 0
    for bit in range(1 << ng):
        nowc = [target]
        for i in range(ng):
            if 1 & (bit >> i):
                nowc.append(graph[target][i])
        is_clique = True
        for i in nowc:
            for j in nowc:
                if i == j:
                    continue
                if i not in S or j not in graph[i]:
                    is_clique = False
        if is_clique:
            cnt = 1
            for v in nowc:
                cnt *= values[v]
                cnt %= MOD
            res += cnt
            res %= MOD
    # targetを含まないようなクリークについては、Sからtargetを除いて再帰的に解く。
    S.remove(target)

    return (res + solve(graph, S)) % MOD


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
n, m = map(int, input().split())
values = list(map(int, input().split()))
adjList = [[] for _ in range(n)]
for _ in range(m):
    u, v = map(int, input().split())
    adjList[u].append(v)
    adjList[v].append(u)

print(solve(adjList, list(range(n))))
