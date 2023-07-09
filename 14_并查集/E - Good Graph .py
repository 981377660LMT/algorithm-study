# https://atcoder.jp/contests/abc304/tasks/abc304_e
# 给定一张无向图，有k个限制，第 i个限制表示 点xi和 点yi不能相互到达。原图满足这k条限制。
# 依次回答q个`独立`的询问，每个询问添加一条边(u,v)后，是否还满足这 k个限制。
# n,m,k,q<=2e5


from typing import List, Tuple
from UnionFind import UnionFindArray


def goodGraph(
    n: int, edges: List[Tuple[int, int]], ban: List[Tuple[int, int]], queries: List[Tuple[int, int]]
) -> List[bool]:
    uf = UnionFindArray(n)
    for u, v in edges:
        uf.union(u, v)
    bad = set()
    for u, v in ban:
        a, b = uf.find(u), uf.find(v)
        if a > b:
            a, b = b, a
        bad.add((a, b))

    res = [True] * len(queries)
    for i, (u, v) in enumerate(queries):
        a, b = uf.find(u), uf.find(v)
        if a > b:
            a, b = b, a
        res[i] = (a, b) not in bad
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
    k = int(input())
    ban = []
    for _ in range(k):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        ban.append((u, v))
    q = int(input())
    queries = []
    for _ in range(q):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        queries.append((u, v))

    res = goodGraph(n, edges, ban, queries)
    for v in res:
        print("Yes" if v else "No")
