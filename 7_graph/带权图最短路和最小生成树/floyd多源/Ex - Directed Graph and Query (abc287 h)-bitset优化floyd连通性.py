"""
!bitset优化floyd连通性 n可以到2000 O(n^3/64)

给定一张有向图,回答q组询问。
每组询问包括 s,t,问从 s 到点 t的路径的代价最小值。
一条路径的代价是其经过的点的编号的最大值。

如何根据连通性来得到距离?
!更新一次k后,更新一下所有询问的答案(如果ij在第k轮循环首次连通,那么距离就为k)
总时间杂度 O(n^3/64+n*q)
"""


from typing import List, Tuple


def directedGraphAndQuery(
    n: int, edges: List[Tuple[int, int]], queries: List[Tuple[int, int]]
) -> List[int]:
    dp = [1 << i for i in range(n)]  # 用bitset存储邻接矩阵
    for u, v in edges:
        dp[u] |= 1 << v

    res = [-1] * len(queries)
    for k in range(n):
        for i in range(n):
            if dp[i] & (1 << k):
                dp[i] |= dp[k]
        for qi, (s, t) in enumerate(queries):
            if res[qi] == -1 and dp[s] & (1 << t):
                res[qi] = max(s, t, k) + 1
    return res


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)

    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))

    q = int(input())
    queries = []
    for _ in range(q):
        start, target = map(int, input().split())
        start, target = start - 1, target - 1
        queries.append((start, target))

    res = directedGraphAndQuery(n, edges, queries)
    print(*res, sep="\n")
