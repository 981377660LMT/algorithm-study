# https://leetcode.cn/problems/XiqZWx/
# 多源d点
# 给出一棵 n 个结点的树，树上每条边的边权都是 1，
# 这 n 个结点中有 m 个特殊点，
# !请你求出树上距离这m个特殊点距离均不超过 d 的点的数量，
# 包含特殊点本身。
# 1 <= n, m, d <= 50000

# !树形dp，O(n)的复杂度
# 维护距离特殊点的`最大`距离
# !点X距离所有的特殊点都在d以内 等价表述：距离点X最远的一个特殊点的距离 <= d
# !因此求出每个点距离的特殊点的距离的最大值max，然后统计max<=d的点的数量即可

import sys
from typing import Literal

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    from Rerooting import Rerooting

    n, m, d = map(int, input().split())
    R = Rerooting(n)
    starts = set([int(x) - 1 for x in input().split()])
    parents = list(map(int, input().split()))
    for cur, pre in enumerate(parents, 1):
        pre -= 1
        R.addEdge(pre, cur)

    def op(fromRes: int, parent: int, cur: int, direction: Literal[0, 1]) -> int:
        return fromRes + 1

    def merge(childRes1: int, childRes2: int) -> int:
        return max(childRes1, childRes2)

    def e(root: int) -> int:
        """如果root不在starts中,返回-INF(因为统计的是距离特殊点的最大距离)"""
        if root in starts:
            return 0
        return -INF

    maxDist = R.rerooting(op=op, merge=merge, e=e, root=0)  # maxDist[i]表示i点距离特殊点的最大距离
    print(sum(dist <= d for dist in maxDist))
