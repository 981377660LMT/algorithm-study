# https://atcoder.jp/contests/abc295/tasks/abc295_g
# 给定一棵外向树，每次修改为连一条有向边 (u,v) 保证 v 可以到达 u，
# 查询某一个节点能到达的编号最小的点。
# 发现每次连边是缩点的过程，而且一定是将 u,v 路径上强连通分量合并，显然可以用并查集维护。
#
# 操作一可以看作 u→v 路径上所有点向 v 连边，我们可以优化为每个点向父节点连边。
# 操作二相当于从 x 出发，沿我们添加的边向上走到最高点。
# 我们考虑使用并查集实现，操作一暴力跳就可以。

from UnionFind import UnionFindArray

from typing import List, Tuple


def minimumRechableCity(
    n: int,
    parents: List[int],
    queries: List[Tuple[int, int, int]],
) -> List[int]:
    uf = UnionFindArray(n)  # 并查集维护强连通分量
    res = []
    for op, *args in queries:
        if op == 1:
            u, v = args
            target = uf.find(v)
            while True:
                # 暴力上跳，向父节点连边union
                u = uf.find(u)
                if u == target:
                    break
                uf.unionTo(u, parents[u])
        else:
            x = args[0]
            res.append(uf.find(x))
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    parents = [-1] + [int(x) - 1 for x in input().split()]
    q = int(input())
    queries = []
    for _ in range(q):
        op, *args = map(int, input().split())
        if op == 1:
            queries.append((op, args[0] - 1, args[1] - 1))
        else:
            queries.append((op, args[0] - 1, 0))

    res = minimumRechableCity(n, parents, queries)
    for r in res:
        print(r + 1)
