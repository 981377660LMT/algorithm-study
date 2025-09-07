# abc422-F - Eat and Ride-动态边权最短路dp
# https://atcoder.jp/contests/abc422/tasks/abc422_f
# 有一个 N 个顶点、M 条边的连通无向图。顶点编号为 1, 2, ..., N。第 i 条边连接着顶点 u_i 和 v_i。
#
# 对于 i = 1, 2, ..., N，请解决以下问题：
#
# 起初，高桥的体重为 0。
#
# 高桥坐车先访问顶点 1，然后向顶点 i 移动。当高桥访问任何一个顶点 v (1≤v≤N) 时，他的体重会增加 W_v。
#
# !高桥乘坐的汽车可以沿着边移动。当高桥通过一条边时，如果他当时的体重是 X，汽车会消耗 X 的燃料。
#
# 请计算高桥到达顶点 i 所需消耗的最小燃料量。
#
# 限制条件
#
# 1 ≤ N ≤ 5000
# 1 ≤ M ≤ 5000
# 1 ≤ W_i ≤ 10^9 (1≤i≤N)
# 1 ≤ u_i ≤ v_i ≤ N (1≤i≤M)
# 给定的图是连通的。
# 所有输入均为整数。
#
# !dp[i][j] -> 从1触发经过i个点到达j的最小花费

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N, M = map(int, input().split())
    W = list(map(int, input().split()))
    edges = [tuple(map(lambda x: int(x) - 1, input().split())) for _ in range(M)]
    dp, ndp = [INF] * N, [INF] * N
    dp[0] = 0
    for i in range(N - 1, 0, -1):
        dp, ndp = ndp, dp
        dp[0] = 0
        for u, v in edges:
            ndp[u] = min(ndp[u], dp[v] + W[v] * i)
            ndp[v] = min(ndp[v], dp[u] + W[u] * i)
    print(*ndp)
