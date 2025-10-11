# D - The Simple Game
# https://atcoder.jp/contests/abc427/tasks/abc427_d
# 有一个 N 个顶点、M 条边的有向图。顶点编号为 1 到 N。第 i 条边是从顶点 U_i 指向顶点 V_i 的有向边。这里，每个顶点的出度（从该顶点出发的边的数量）至少为 1。

# 此外，每个顶点上都写有一个字符，顶点 i 上写的字符是 S_i（字符串 S 的第 i 个字符）。

# Alice 和 Bob 在这个图上使用一个棋子进行以下游戏：

# 开始时，棋子放在顶点 1 上。
# Alice 先手，Bob 后手，双方交替进行以下操作，每人 K 次（总共 2K 步）。
# 操作：假设当前棋子在顶点 u。选择一个顶点 v，使得存在从 u 到 v 的边，然后将棋子移动到顶点 v。
# 游戏结束后（即总共移动 2K 步后），棋子最终所在的顶点为 v。如果 S_v = 'A'，则 Alice 获胜；如果 S_v = 'B'，则 Bob 获胜。

# 请你判断在双方都采取最优策略的情况下，谁会是赢家。

# 一个输入包含 T 个测试用例，请对每个测试用例作答。
# 1 ≤ T
# 2 ≤ N, M ≤ 2 × 10^5
# 1 ≤ K ≤ 10


def solve():
    n, m, k = map(int, readline().split())

    dp = [int(i == "A") for i in readline().strip()]
    g = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, readline().split())
        u -= 1
        v -= 1
        g[u].append(v)

    for _ in range(2 * k):
        ndp = [0] * n
        for i in range(n):
            for j in g[i]:
                if dp[j] == 0:
                    ndp[i] = 1
                    break

        dp = ndp

    return dp[0]


import sys

readline = sys.stdin.readline

T = int(readline())
for _ in range(T):
    res = solve()
    print("Alice" if res else "Bob")
