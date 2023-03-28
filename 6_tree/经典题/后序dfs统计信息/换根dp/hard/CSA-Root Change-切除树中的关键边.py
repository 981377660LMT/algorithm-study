# CSA-Root Change-切除树中的关键边
# https://csacademy.com/contest/archive/task/root-change/
# 对每个根节点询问：
# !存在有多少条边,切除后可以不改变原树的高度?
# 切除树中任意一条边后,连接这条边的子树中的所有结点全都会消失.


from typing import Tuple
from Rerooting import Rerooting

import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)

if __name__ == "__main__":

    E = Tuple[int, int]  # (maxHeight,critialEdge)

    def e(root: int) -> E:
        return (0, 0)

    def op(childRes1: E, childRes2: E) -> E:
        h1, _ = childRes1
        h2, _ = childRes2
        if h1 == h2:
            return (h1, 0)  # 子树中所有边不再关键
        return childRes1 if h1 > h2 else childRes2

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        maxHeight, critialEdge = fromRes
        return (maxHeight + 1, critialEdge + 1)

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    for _, critialEdge in dp:
        print(n - 1 - critialEdge)
