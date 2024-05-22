# 一棵树上每个点都有一个字母， 问树中路径同时存在 you 三个字母的路径条数
#
# !如果是笔试，先写一个记忆化搜索dfs(cur,pre,mask)
# 记忆化搜索可以过大部分用例，只会被菊花图卡掉
#
# 换根dp，以每个点为根，统计以每个点为根的路径上同时存在you的路径数
# !以每个节点为根，作为路径的一个端点，子树内每个点到根的路径有8种状态，向上统计状态个数就可以了

from typing import List
from Rerooting import Rerooting


if __name__ == "__main__":
    E = List[int]  # [0,0,0,0,0,0,0,0] y->1 o->2 u->4

    def e(_: int) -> E:
        return [0] * 8

    def op(childRes1: E, childRes2: E) -> E:
        a1, b1, c1, d1, e1, f1, g1, h1 = childRes1
        a2, b2, c2, d2, e2, f2, g2, h2 = childRes2
        return [a1 + a2, b1 + b2, c1 + c2, d1 + d2, e1 + e2, f1 + f2, g1 + g2, h1 + h2]

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        from_ = cur if direction == 0 else parent
        char = chars[from_]
        charMask = getMask(char)
        res = [0] * 8
        for preMask, v in enumerate(fromRes):
            nextMask = preMask | charMask
            if preMask == nextMask:
                res[preMask] += 1
            else:
                res[preMask] = 0
                res[nextMask] += v + 1
        return res

    def getMask(c: str) -> int:
        return 1 if c == "y" else 2 if c == "o" else 4

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))
    chars = input()
    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)
    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    res = 0
    for i, states in enumerate(dp):
        curMask = getMask(chars[i])
        for s, v in enumerate(states):
            if s | curMask == 7:
                res += v
    print(res // 2)

    print(dp.count(0))  # 不在最大匹配中的点的个数
