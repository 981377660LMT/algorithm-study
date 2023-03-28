# https://yukicoder.me/problems/no/1718
# 从每个结点出发,求遍历所有特殊点的最短距离之和
# n<=1e5
# 每个

from collections import defaultdict
from typing import List, Tuple
from Rerooting import Rerooting

INF = int(1e18)


def randomSquirrel(n: int, edges: List[Tuple[int, int]], specials: List[int]) -> List[int]:
    E = Tuple[int, int]  # (maxDist, weightSum)

    def e(root: int) -> E:
        return (0, 0) if isSpecial[root] else (-INF, 0)

    def op(childRes1: E, childRes2: E) -> E:
        dist1, wSum1 = childRes1
        dist2, wSum2 = childRes2
        return (max(dist1, dist2), wSum1 + wSum2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        dist, wSum = fromRes
        w = weights[cur][parent] if direction == 0 else weights[parent][cur]
        # dist>=0 时这条边才开始算入必须经过的路径
        return (dist + w, wSum + 1) if dist >= 0 else (dist, wSum)

    isSpecial = [False] * n
    for v in specials:
        isSpecial[v] = True

    R = Rerooting(n)
    weights = [defaultdict(int) for _ in range(n)]
    for u, v in edges:
        R.addEdge(u, v)
        weights[u][v] = 1
        weights[v][u] = 1

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    res = [wSum * 2 - maxDist for maxDist, wSum in dp]  # !可以不回到出发点,2*边权-最大距离
    return res


if __name__ == "__main__":
    n, k = map(int, input().split())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
    sweets = [int(x) - 1 for x in input().split()]
    print(*randomSquirrel(n, edges, sweets), sep="\n")
