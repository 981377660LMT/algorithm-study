# https://www.luogu.com.cn/problem/CF771C
# 有一棵 n 个节点的树，树上住着一只熊。
# 若熊在节点 u，熊跳一次可以从 u 跳到离 u 树上距离 ≤k 的任意节点 v。
# 设 f(s,t) 为熊从 s 到 t 跳的次数的最小值。
# !求∑f(s,t)(0<=s<t<n) n<=2e5 1<=k<=5
# cf405-Bear and Tree Jumps-每个结点子树中距离种类数

from typing import Tuple, List
from Rerooting import Rerooting

import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)

if __name__ == "__main__":

    E = Tuple[List[int], int]  # (每种距离(0-k-1)的个数,距离正好为k时全变成距离为0,记入跳跃次数中).

    def e(root: int) -> E:
        return ([0] * k, 0)

    def op(childRes1: E, childRes2: E) -> E:
        dist1, sum1 = childRes1
        dist2, sum2 = childRes2
        return ([a + b for a, b in zip(dist1, dist2)], sum1 + sum2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        dist, sum_ = fromRes

        if dist == [0] * k:
            if k == 1:
                return [1], sum_ + 1
            newDist = dist[:]
            newDist[1] = 1
            return newDist, sum_

        if k == 1:
            newDist = dist[:]
            newDist[0] += 1
            sum_ += newDist[0]
            return newDist, sum_
        sum_ += dist[-1]
        newDist = [dist[-1]] + dist[:-1]
        newDist[1] += 1
        return newDist, sum_

    n, k = map(int, input().split())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))
    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    res = 0
    for dist, sum_ in dp:
        res += sum(dist[1:]) + sum_
    print(res // 2)
