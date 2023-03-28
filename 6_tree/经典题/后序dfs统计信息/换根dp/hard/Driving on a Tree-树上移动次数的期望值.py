# https://atcoder.jp/contests/s8pc-4/tasks/s8pc_4_d
# 从根节点出发,等概率转移到相邻的`未访问`过的结点
# !对每个根节点,求移动次数的期望值
# Driving on a Tree-树上移动次数的期望值


import sys
from typing import Tuple
from Rerooting import Rerooting

input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)

if __name__ == "__main__":

    E = Tuple[float, int]  # (期望次数,权重)

    def e(root: int) -> E:
        return (0, 0)

    def op(childRes1: E, childRes2: E) -> E:
        dp1, weight1 = childRes1
        dp2, weight2 = childRes2
        if weight1 + weight2 == 0:
            return (0, 0)
        return (((dp1 * weight1) + (dp2 * weight2)) / (weight1 + weight2)), weight1 + weight2

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        dp, _ = fromRes
        return (dp + 1, 1)

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)
    tree = R.adjList
    dp = R.rerooting(e=e, op=op, composition=composition, root=0)

    for res, _ in dp:
        print(res)
