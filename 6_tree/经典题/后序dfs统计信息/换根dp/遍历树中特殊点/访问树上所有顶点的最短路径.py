# (访问所有顶点的最短路径)
# https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=1595
# 给定一棵树,求从每个点出发访问所有顶点的最短路径.
# !不需要返回到出发点.
# !等于两倍边权之和减去每个点到最远点的距离,其中最好的出发点就是直径的两个端点。

from collections import defaultdict
from typing import List, Tuple
from Rerooting import Rerooting


def solve(n: int, edges: List[Tuple[int, int, int]]) -> List[int]:
    E = int

    def e(root: int) -> E:
        return 0

    def op(childRes1: E, childRes2: E) -> E:
        return max(childRes1, childRes2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        return fromRes + weight[cur][parent]

    R = Rerooting(n)
    weight = [defaultdict(int) for _ in range(n)]
    wSum = 0
    for u, v, w in edges:
        R.addEdge(u, v)
        weight[u][v] = w
        weight[v][u] = w
        wSum += w
    dp = R.rerooting(e, op, composition)
    print(dp)
    return [2 * wSum - x for x in dp]


if __name__ == "__main__":
    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, 1))
    res = solve(n, edges)
    print(*res, sep="\n")
