# cf302-Road Improvement-改造不良道路的方案数
# https://codeforces.com/contest/543/problem/D

# 给定一棵树
# 所有的道路最初都是不良的，但是政府想要改善一些路的状况。
# !我们认为如果从首都x城到其他城市的道路最多包含一条不良道路，
# 市民会对此感到满意。
# 你的任务是对于每一个可能的x，
# 求出所有能够满足市民条件的改良道路的方式MOD 1e9+7。

# 和
# !6_tree/经典题/后序dfs统计信息/换根dp/hard/EDPC-Subtree-所有的黑色节点组成一个联通块的染色方案数.py
# 转移函数一样


import sys
from Rerooting import Rerooting

input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)
MOD = int(1e9 + 7)

if __name__ == "__main__":

    E = int

    def e(root: int) -> E:
        return 1

    def op(childRes1: E, childRes2: E) -> E:
        return childRes1 * childRes2 % MOD

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        return fromRes + 1  # 加一是子树内全都修的方案数

    n = int(input())
    edges = []
    parents = list(map(int, input().split()))
    for i in range(1, n):
        p = parents[i - 1] - 1
        edges.append((p, i))

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    print(*dp)
