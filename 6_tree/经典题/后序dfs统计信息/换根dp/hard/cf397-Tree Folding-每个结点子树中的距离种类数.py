# cf397-Tree Folding-每个结点子树中的距离种类数
# 给你一棵树，可以把树上父亲相同的两条长度相同的链合并。
# 问你最后能不能变成一条链，
# 能的话求链的最短长度。 （n<2*10^5）
# https://www.luogu.com.cn/problem/CF765E
# https://blog.csdn.net/baidu_35520981/article/details/55261333

# 以每个点为根,看子树内的距离种类数是否可行
# !`去除根节点`的子树内每个结点处距离种数必须为1
# 答案为根节点处距离(1种或2种)的和

# !注意答案为偶数时,需要除以2直到奇数为止

import sys
from typing import Set, Tuple
from Rerooting import Rerooting

input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(1e9)

if __name__ == "__main__":
    E = Tuple[Set[int], int]  # (子树中的距离种类, 子树中有多少个结点处距离种类为2)

    def e(root: int) -> E:
        return set(), 0

    def op(childRes1: E, childRes2: E) -> E:
        dist1, count1 = childRes1
        dist2, count2 = childRes2
        newDist = dist1 | dist2
        if len(newDist) > 2:
            return (set(), INF)
        return (newDist, count1 + count2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        dist, count = fromRes
        if len(dist) == 0 and count == 0:
            return (set([1]), 0)
        return (set([x + 1 for x in dist]), count + (len(dist) == 2))

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))
    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)
    dp = R.rerooting(e=e, op=op, composition=composition, root=0)

    res = INF
    for i in range(n):
        dists, count2 = dp[i]
        if count2 == 0:
            dist = sum(dists)
            while dist % 2 == 0:
                dist //= 2
            res = min(res, dist)
    print(res if res != INF else -1)
