# https://yukicoder.me/problems/no/1295


# 给定一个n个点的树
# 所有的顶点分为"访问过的顶点"和"未访问过的顶点"
# 开始时,将棋子放在根节点i上
# 每个回合,有两种移动方式:
# !1.将棋子移动到相邻的"未访问过的顶点"中编号最小的顶点上
# !2.将棋子移动到相邻的"访问过的顶点"中编号最小的顶点上
# 问是否可以使得所有的顶点都被访问到

# 对每个初始的根节点询问

from Rerooting import Rerooting


from typing import Tuple


INF = int(4e18)

if __name__ == "__main__":
    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

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

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

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
