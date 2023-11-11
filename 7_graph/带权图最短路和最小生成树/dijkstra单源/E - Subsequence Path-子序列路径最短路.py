"""
现在有 n 个点, m 条带权有向边, 现在我们给定一个序列 E , 我们称从1 到 n 的一条路径为好的, 
当且仅当他经过的边的编号, 是 E 的子序列, 现在请找出最短的好的路径, 如果没有则输出-1

按照指定顺序松弛边的dijk过程
"""

# 子序列路径最短路/子序列最短路
# 最短路经过的边需要为nums的子序列
# 求0到n-1的最短路
# n,m<=2e5

# !dp[i][pos] 表示E的前i个元素时，到达pos的最短距离
# !每个元素选或不选 注意到可以in-place更新

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m, k = map(int, input().split())
    edges = []
    for i in range(m):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, w))

    nums = [int(x) - 1 for x in input().split()]  # !路径必须为E的子序列

    dist = [INF] * n
    dist[0] = 0
    for ei in nums:
        u, v, w = edges[ei]
        dist[v] = min(dist[v], dist[u] + w)  # !规定松弛边顺序下的dijk过程

    print(dist[n - 1] if dist[n - 1] != INF else -1)
