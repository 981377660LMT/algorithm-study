# https://yukicoder.me/problems/no/1295


# 给定一个n个点的树,树上有m只猫位于不同的点上
# 初始时,先手选一个位置,将猫移动到`这只猫`没有访问过的空位置
# 交替进行,直到无法移动
# !问先手必胜的方案,输出(选择第i只猫,第一步移动到v)的 (i,v) 对
# 无解输出(-1,-1)


# 树上的Grundy数
# TODO

from Rerooting import Rerooting


from typing import Tuple


INF = int(4e18)

if __name__ == "__main__":

    E = Tuple[int, bool]

    def e(root: int) -> E:
        return (0, True)

    def op(childRes1: E, childRes2: E) -> E:
        dist1, ok1 = childRes1
        dist2, ok2 = childRes2
        return (max(dist1, dist2), dist1 == dist2 or (ok1 and ok2))

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        dist, ok = fromRes
        return (dist + 1, ok)

    n, m = map(int, input().split())
    positions = [int(input()) - 1 for _ in range(m)]  # 开始时每个猫咪的位置
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    def bitMex(x: int) -> int:
        ...

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)

    res, ok = INF, False
    for i in range(n):
        dist, curOk = dp[i]
        if curOk:
            ok = True
            res = min(res, dist)
    if ok:
        print(res + 1)
    else:
        print(-1)
