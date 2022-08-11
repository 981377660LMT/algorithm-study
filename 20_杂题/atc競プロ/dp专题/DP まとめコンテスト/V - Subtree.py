# 在树上选取若干个点，这些点之间都是连通的
# 对v从0到n-1 输出`选取顶点v的组合`有多少种
# n<=1e5


import sys
from typing import Any, Callable, List


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)


# n, MOD = map(int, input().split())
# adjList = [[] for _ in range(n)]
# for _ in range(n - 1):
#     a, b = map(int, input().split())
#     adjList[a - 1].append(b - 1)
#     adjList[b - 1].append(a - 1)


# TLE TODO
# @lru_cache(None)
# def dfs(cur: int, pre: int) -> int:
#     res = 1
#     for next in adjList[cur]:
#         if next == pre:
#             continue
#         nextRes = dfs(next, cur)
#         res *= nextRes + 1  # 加1表示子节点不选
#         res %= MOD
#     return res


# for cur in range(n):
#     print(dfs(cur, -1))

###############################################################
# !正解似乎是换根dp
def rerooting(
    n: int,
    adjList: List[List[int]],
    init: Callable[[], int],
    child: Callable[[int], int],
    merge: Callable[[int, int], int],
):
    """换根dp框架 https://atcoder.jp/contests/dp/submissions/33369486

    以每个结点作为根 求某种性质

    Args:
        n: 顶点数
        adjList: 树 编号从1-n 0表示虚拟根节点
        init: 初始化根节点的值
        node: 子结点贡献值
        merge: 子节点与父节点的合并方式
    """

    stack = [(1, 0)]
    order = []
    while stack:
        pair = stack.pop()
        cur, pre = pair
        order.append(pair)
        for next in adjList[cur]:
            if next == pre:
                continue
            stack.append((next, cur))

    dp1 = [0] * (n + 1)
    ls = [0] * (n + 1)
    rs = [0] * (n + 1)

    for cur, pre in reversed(order):  # dfs序
        nexts = adjList[cur]
        v = init()
        for next in nexts:
            if next == pre:
                continue
            ls[next] = v
            v = merge(v, child(dp1[next]))
        v = init()
        for next in reversed(nexts):
            if next == pre:
                continue
            rs[next] = v
            v = merge(v, child(dp1[next]))
        dp1[cur] = v

    dp2 = [0] * (n + 1)
    res = [0] * (n + 1)
    res[1] = dp1[1]
    for i in range(1, n):
        cur, pre = order[i]
        dp2[cur] = v = merge(merge(ls[cur], rs[cur]), child(dp2[pre]))
        res[cur] = merge(child(v), dp1[cur])
    return res


n, MOD = map(int, input().split())
if n == 1:
    print(1)
    exit(0)

adjList = [[] for _ in range(n + 1)]
for _ in range(n - 1):
    a, b = map(int, input().split())
    adjList[a].append(b)
    adjList[b].append(a)


res = rerooting(
    n, adjList, init=lambda: 1, child=lambda x: (x + 1) % MOD, merge=lambda x, y: x * y % MOD
)
print(*res[1:], sep="\n")
