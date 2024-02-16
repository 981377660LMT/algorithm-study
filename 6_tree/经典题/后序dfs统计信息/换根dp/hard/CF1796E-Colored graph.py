# https://www.luogu.com.cn/problem/CF1796E
# !将一棵树划分成若干条不相交的链（“不相交”指每个结点在且仅在一条链），求所有划分中，最短链长度得最大值
#
# 解：
# !自底向上找链，对于结点x，最优情况是将结点x与它的子结点y中最短的一条链连接，而其他的结点y对应的链不会继续增长了，需要将其记录

from Rerooting import Rerooting

from typing import List, Tuple

INF = int(1e18)
E = Tuple[int, int]  # 当前链的长度, 其余不会增长的链的最短长度


def cf1796E(edges: List[Tuple[int, int]]) -> int:
    n = len(edges) + 1

    def e(root: int) -> E:
        return 0, INF

    def op(childRes1: E, childRes2: E) -> E:
        curLen1, minLen1 = childRes1
        curLen2, minLen2 = childRes2
        if curLen1 == 0 or curLen2 == 0:
            return childRes1 if curLen2 == 0 else childRes2
        if curLen1 > curLen2:
            curLen1, curLen2 = curLen2, curLen1
            minLen1, minLen2 = minLen2, minLen1
        return curLen1, min(curLen2, minLen1, minLen2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        curLen, minLen = fromRes
        return curLen + 1, minLen

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    return max(min(a + 1, b) for a, b in dp)  # +1 是因为dp没有包含顶点


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    T = int(input())
    for _ in range(T):
        n = int(input())
        edges = []
        for i in range(n - 1):
            u, v = map(int, input().split())
            u, v = u - 1, v - 1
            edges.append((u, v))
        print(cf1796E(edges))
