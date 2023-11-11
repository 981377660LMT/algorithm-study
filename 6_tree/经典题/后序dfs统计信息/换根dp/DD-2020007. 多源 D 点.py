# https://www.luogu.com.cn/problem/CF337D
# https://leetcode.cn/problems/XiqZWx/

# 多源d点
# 给出一棵 n 个结点的树，树上每条边的边权都是 1，
# 这 n 个结点中有 m 个特殊点，
# !请你求出树上距离这m个特殊点距离均不超过 d 的点的数量，
# 包含特殊点本身。
# 1 <= n, m, d <= 50000


# 求树中每个点到特殊点的最远距离
# !树形dp，O(n)的复杂度
# 维护距离特殊点的`最大`距离
# !点X距离所有的特殊点都在d以内 等价表述：距离点X最远的一个特殊点的距离 <= d
# !因此求出每个点距离的特殊点的距离的最大值max，然后统计max<=d的点的数量即可

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    from Rerooting import Rerooting

    E = int

    def e(root: int) -> E:
        """如果root不在starts中,返回-INF(因为统计的是距离特殊点的最大距离)"""
        return 0 if root in monsters else -INF

    def op(childRes1: E, childRes2: E) -> E:
        return max(childRes1, childRes2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        if direction == 0:  # cur -> parent
            return fromRes + 1
        return fromRes + 1  # parent -> cur

    n, m, d = map(int, input().split())
    monsters = [int(x) - 1 for x in input().split()]
    monsters = set(monsters)
    R = Rerooting(n)
    for _ in range(n - 1):
        u, v = map(int, input().split())
        R.addEdge(u - 1, v - 1)
    maxDist = R.rerooting(e=e, op=op, composition=composition, root=0)  # maxDist[i]表示i点距离特殊点的最大距离
    print(sum(dist <= d for dist in maxDist))
